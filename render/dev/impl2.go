package dev

import (
	"context"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"image/draw"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
)

type Renderer2 struct {
	s                sdf.SDF2 // The SDF to render
	evalMin, evalMax float64  // The pre-computed minimum and maximum of the whole surface (for stable colors and speed)
}

func newDevRenderer2(s sdf.SDF2) devRendererImpl {
	r := &Renderer2{
		s: s,
	}
	return r
}

func (r *Renderer2) Dimensions() int {
	return 2
}

func (r *Renderer2) BoundingBox() sdf.Box3 {
	bb := r.s.BoundingBox()
	return sdf.Box3{Min: bb.Min.ToV3(0), Max: bb.Max.ToV3(0)}
}

func (r *Renderer2) Render(ctx context.Context, state *RendererState, stateLock,
	cachedRenderLock *sync.RWMutex, partialImages chan<- *image.RGBA, fullRender *image.RGBA) error {
	if r.evalMin == 0 && r.evalMax == 0 { // First render (ignoring external cache)
		// Compute minimum and maximum evaluate values for a shared color scale for all blocks
		r.evalMin, r.evalMax = utilSdf2MinMax(r.s, r.s.BoundingBox(), sdf.V2i{128, 128} /* TODO: Configurable? */)
		//log.Println("MIN:", r.evalMin, "MAX:", r.evalMax)
	}

	fullRenderSize := fullRender.Bounds().Size()
	bbAspectRatio := state.Bb.Size().X / state.Bb.Size().Y
	stateLock.Lock() // Maintain Bb aspect ratio on ResInv change, increasing the size as needed
	screenAspectRatio := float64(fullRenderSize.X) / float64(fullRenderSize.Y)
	if math.Abs(bbAspectRatio-screenAspectRatio) > 1e-12 {
		if bbAspectRatio > screenAspectRatio {
			scaleYBy := bbAspectRatio / screenAspectRatio
			state.Bb = sdf.NewBox2(state.Bb.Center(), state.Bb.Size().Mul(sdf.V2{X: 1, Y: scaleYBy}))
		} else {
			scaleXBy := screenAspectRatio / bbAspectRatio
			state.Bb = sdf.NewBox2(state.Bb.Center(), state.Bb.Size().Mul(sdf.V2{X: scaleXBy, Y: 1}))
		}
	}
	stateLock.Unlock()

	// Create the new full CPU image (downscaled by resolution)
	fullImgSize := sdf.V2i{fullRenderSize.X, fullRenderSize.Y} // screenSize.ToV2().DivScalar(float64(resolution)).ToV2i()
	fullImg := fullRender                                      //image.NewRGBA(image.Rect(0, 0, fullImgSize[0], fullImgSize[1]))
	for i := 3; i < len(fullImg.Pix); i += 4 {
		fullImg.Pix[i] = 255 // Set all pixels to transparent initially
	}

	// Render each blockIndex of the image individually to allow cancelling the render
	pixelsPerBlock := sdf.V2i{128, 128} /* TODO: Configurable? */
	numBlocks := fullImgSize.ToV2().Div(pixelsPerBlock.ToV2()).Ceil().ToV2i()

	// Parallelize: spawn workers
	jobCount := numBlocks[0] * numBlocks[1]
	blockIndexJobs := make(chan sdf.V2i)
	errors := make(chan error, jobCount) // Buffer them to avoid deadlocks
	errCount := uint32(0)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				blockIndexTask, ok := <-blockIndexJobs
				if !ok {
					break
				}
				errors <- r.renderBlock(ctx, fullImg, blockIndexTask, pixelsPerBlock, state, stateLock, cachedRenderLock, numBlocks, fullImgSize, partialImages)
				if atomic.AddUint32(&errCount, 1) == uint32(jobCount) {
					close(errors)
				}
			}
		}()
	}

	// Generate jobs forming a spiral
	blockIndex := numBlocks.ToV2().DivScalar(2).ToV2i()
	blockIndexJobs <- blockIndex
	jobN := 1
	for n := 0; jobN < jobCount; n++ {
		stepSize := n/2 + 1            // 1, 1, 2, 2, 3, 3...
		stepDir := dirs2[n%len(dirs2)] // Up, Right, Down, Left...
		for step := 0; step < stepSize; step++ {
			blockIndex = blockIndex.Add(stepDir)
			if blockIndex[0] >= 0 && blockIndex[1] >= 0 && blockIndex[0] < numBlocks[0] && blockIndex[1] < numBlocks[1] {
				// This will avoid spawning new tasks if any of them failed previously (racy: solved locking state)
				select {
				case err := <-errors:
					if err != nil {
						// The first block that renders with an error closes the partial image channel
						close(blockIndexJobs)
						if partialImages != nil {
							close(partialImages)
						}
						return err // Quick exit on error
					}
				default:
				}
				blockIndexJobs <- blockIndex
				jobN++
				if jobN == jobCount {
					break // Will also break parent loop
				}
			}
		}
	}
	close(blockIndexJobs)
	// Wait for the full image (only final blocks remaining, as jobs are created as they are ready to be processed)
	for err := range errors {
		err = <-errors
		if err != nil {
			// The first block that renders with an error closes the partial image channel
			if partialImages != nil {
				close(partialImages)
			}
			return err // Quick exit on error
		}
	}
	// If no block threw an error, close partialImages now
	if partialImages != nil {
		close(partialImages)
	}
	// TODO: Draw bounding boxes over the image
	return nil
}

func (r *Renderer2) renderBlock(ctx context.Context, fullImg *image.RGBA, blockIndex sdf.V2i,
	pixelsPerBlock sdf.V2i, state *RendererState, stateLock, cachedRenderLock *sync.RWMutex, numBlocks sdf.V2i, fullImgSize sdf.V2i,
	partialImages chan<- *image.RGBA) (err error) {
	select {
	case <-ctx.Done(): // Render cancelled
		return ctx.Err()
	default: // Render not cancelled
	}

	stateLock.RLock()
	defer stateLock.RUnlock()

	// Compute positions and sizes
	blockStartPixel := blockIndex.ToV2().Mul(pixelsPerBlock.ToV2()).ToV2i()
	blockSizePixels := pixelsPerBlock.AddScalar(1) // nextPowerOf2(state.ResInv)
	if blockIndex[0] == numBlocks[0]-1 {
		blockSizePixels[0] = fullImgSize[0] - blockStartPixel[0] + 1 //+ nextPowerOf2(state.ResInv)
	}
	if blockIndex[1] == numBlocks[1]-1 { // Inverted Y
		blockSizePixels[1] = fullImgSize[1] - blockStartPixel[1] + 1 //+ nextPowerOf2(state.ResInv)
	}
	//blockSizePixels = blockSizePixels.ToV2().DivScalar(float64(state.ResInv)).ToV2i()
	if blockSizePixels[0] == 0 || blockSizePixels[1] == 0 {
		return nil // Empty block ignored
	}
	blockBb := sdf.Box2{
		Min: state.Bb.Min.Add(state.Bb.Size().Mul(blockStartPixel.ToV2().Div(fullImgSize.ToV2()))),
		Max: state.Bb.Min.Add(state.Bb.Size().Mul(
			blockStartPixel.Add(blockSizePixels).ToV2().Div(fullImgSize.ToV2()))),
	}
	if blockBb.Size().Length2() <= 1e-12 || blockSizePixels.ToV2().Length2() < 1e-12 { // SANITY CHECK that skips the block
		log.Println("SANITY CHECK FAILED: PIXELS: start:", blockStartPixel, "size:", blockSizePixels, "| BOUNDING BOX:", blockBb)
		return nil
	}
	//if sdfSkip.Contains(blockBb.Min) && sdfSkip.Contains(blockBb.Max) {
	//	return nil // Block is fully contained in the skipped section of the screen, ignore
	//}
	//log.Println("PIXELS: start:", blockStartPixel, "size:", blockSizePixels, "| BOUNDING BOX:", blockBb)

	// Render the current block to a CPU image
	png, err := render.NewPNG("unused", blockBb, blockSizePixels)
	if err != nil {
		return err
	}
	evalMin, evalMax := r.evalMin, r.evalMax
	if state.blackAndWhite {
		evalMin, evalMax = -1e-12, 1e-12
	}
	png.RenderSDF2MinMax(r.s, evalMin, evalMax)
	blockImg := png.Image()

	// Merge blocks to full render image (CPU, downscaled)
	cachedRenderLock.Lock()
	translateY := fullImgSize[1] - (blockStartPixel[1] + blockSizePixels[1]) + 1 // block Y is inverted
	draw.Draw(fullImg, image.Rect(blockStartPixel[0], translateY,
		blockStartPixel[0]+blockSizePixels[0], translateY+blockSizePixels[1]),
		blockImg, image.Point{}, draw.Over)
	cachedRenderLock.Unlock()

	// Notify of partial image progress
	if err != nil {
		return err
	}
	if partialImages != nil {
		partialImages <- fullImg
	}
	return nil
}
