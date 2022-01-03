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
)

type devRenderer2 struct {
	s                sdf.SDF2 // The SDF to render
	evalMin, evalMax float64  // The pre-computed minimum and maximum of the whole surface (for stable colors and speed)
}

func newDevRenderer2(s sdf.SDF2) devRendererImpl {
	r := &devRenderer2{
		s: s,
	}
	return r
}

func (r *devRenderer2) Dimensions() int {
	return 2
}

func (r *devRenderer2) BoundingBox() sdf.Box3 {
	bb := r.s.BoundingBox()
	return sdf.Box3{Min: bb.Min.ToV3(0), Max: bb.Max.ToV3(0)}
}

func (r *devRenderer2) Render(ctx context.Context, screenSize sdf.V2i, state *DevRendererState, stateLock, cachedRenderLock *sync.RWMutex, partialImages chan<- *image.RGBA) (*image.RGBA, error) {
	if r.evalMin == 0 && r.evalMax == 0 { // First render (ignoring external cache)
		// Compute minimum and maximum evaluate values for a shared color scale for all blocks
		r.evalMin, r.evalMax = utilSdf2MinMax(r.s, r.s.BoundingBox(), sdf.V2i{64, 64} /* TODO: Configurable? */)
		//log.Println("MIN:", r.evalMin, "MAX:", r.evalMax)
	}
	stateLock.Lock()
	if state.Bb.Size().Length2() == 0 {
		state.Bb = r.s.BoundingBox() // 100% zoom (will fix aspect ratio later)
	}
	// Maintain Bb aspect ratio on Resolution change, increasing the size as needed
	bbAspectRatio := state.Bb.Size().X / state.Bb.Size().Y
	screenAspectRatio := float64(screenSize[0]) / float64(screenSize[1])
	if math.Abs(bbAspectRatio-screenAspectRatio) > 1e-12 {
		if bbAspectRatio > screenAspectRatio {
			scaleYBy := bbAspectRatio / screenAspectRatio
			state.Bb = sdf.NewBox2(state.Bb.Center(), state.Bb.Size().Mul(sdf.V2{X: 1, Y: scaleYBy}))
		} else {
			scaleXBy := screenAspectRatio / bbAspectRatio
			state.Bb = sdf.NewBox2(state.Bb.Center(), state.Bb.Size().Mul(sdf.V2{X: scaleXBy, Y: 1}))
		}
	}
	resolution := state.Resolution
	stateLock.Unlock()

	// Create the new full CPU image (downscaled by resolution)
	fullImgSize := screenSize.ToV2().DivScalar(float64(resolution)).ToV2i()
	fullImg := image.NewRGBA(image.Rect(0, 0, fullImgSize[0], fullImgSize[1]))
	for i := 3; i < len(fullImg.Pix); i += 4 {
		fullImg.Pix[i] = 255 // Set all pixels to transparent initially
	}

	// Render each blockIndex of the image individually to allow cancelling the render
	pixelsPerBlock := sdf.V2i{128, 128} /* TODO: Configurable? */
	// FIXME: Resolution >= 8 causes infinite loop???
	numBlocks := fullImgSize.ToV2().Div(pixelsPerBlock.ToV2()).Ceil().ToV2i()
	// Parallelize: spawn workers
	jobCount := numBlocks[0] * numBlocks[1]
	blockIndexJobs := make(chan sdf.V2i, jobCount)
	errors := make(chan error)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				blockIndexTask, ok := <-blockIndexJobs
				if !ok {
					break
				}
				err := r.renderBlock(ctx, fullImg, blockIndexTask, pixelsPerBlock, state, stateLock, cachedRenderLock, numBlocks, fullImgSize, partialImages)
				errors <- err // Probably nil
			}
		}()
	}
	// Generate jobs forming a spiral
	blockIndex := numBlocks.ToV2().DivScalar(2).ToV2i()
	blockIndexJobs <- blockIndex
	jobN := 1
jobSpiral:
	for n := 0; ; n++ {
		stepSize := n/2 + 1            // 1, 1, 2, 2, 3, 3...
		stepDir := dirs2[n%len(dirs2)] // Up, Right, Down, Left...
		for step := 0; step < stepSize; step++ {
			blockIndex = blockIndex.Add(stepDir)
			if blockIndex[0] >= 0 && blockIndex[1] >= 0 && blockIndex[0] < numBlocks[0] && blockIndex[1] < numBlocks[1] {
				blockIndexJobs <- blockIndex
				jobN++
				if jobN == jobCount {
					break jobSpiral
				}
			}
		}
	}
	//blockIndex := sdf.V2i{}
	//for blockIndex[0] = 0; blockIndex[0] < numBlocks[0]; blockIndex[0]++ {
	//	for blockIndex[1] = 0; blockIndex[1] < numBlocks[1]; blockIndex[1]++ {
	//		blockIndexJobs <- blockIndex
	//	}
	//}
	close(blockIndexJobs)
	// Return the full image
	var err error
	for i := 0; i < jobCount; i++ {
		err = <-errors
		if err != nil {
			return fullImg, err // Quick exit on error
		}
	}
	return fullImg, err
}

func (r *devRenderer2) renderBlock(ctx context.Context, fullImg *image.RGBA, blockIndex sdf.V2i,
	pixelsPerBlock sdf.V2i, state *DevRendererState, stateLock, cachedRenderLock *sync.RWMutex, numBlocks sdf.V2i, fullImgSize sdf.V2i,
	partialImages chan<- *image.RGBA) error {
	select {
	case <-ctx.Done(): // Render cancelled
		return ctx.Err()
	default: // Render not cancelled
	}

	stateLock.RLock()
	defer stateLock.RUnlock()

	// Compute positions and sizes
	blockStartPixel := blockIndex.ToV2().Mul(pixelsPerBlock.ToV2()).ToV2i()
	blockSizePixels := pixelsPerBlock.AddScalar(1) // nextPowerOf2(state.Resolution)
	if blockIndex[0] == numBlocks[0]-1 {
		blockSizePixels[0] = fullImgSize[0] - blockStartPixel[0] + 1 //+ nextPowerOf2(state.Resolution)
	}
	if blockIndex[1] == numBlocks[1]-1 { // Inverted Y
		blockSizePixels[1] = fullImgSize[1] - blockStartPixel[1] + 1 //+ nextPowerOf2(state.Resolution)
	}
	blockSizePixels = blockSizePixels //.ToV2().DivScalar(float64(state.Resolution)).ToV2i()
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
	png.RenderSDF2MinMax(r.s, r.evalMin, r.evalMax)
	blockImg := png.Image()

	// Merge blocks to full render image (CPU, downscaled)
	cachedRenderLock.Lock()
	translateY := fullImgSize[1] - (blockStartPixel[1] + blockSizePixels[1]) + 1 // block Y is inverted
	draw.Draw(fullImg, image.Rect(blockStartPixel[0], translateY,
		blockStartPixel[0]+blockSizePixels[0], translateY+blockSizePixels[1]),
		blockImg, image.Point{}, draw.Over)
	cachedRenderLock.Unlock()

	// Notify of partial image progress
	select { // FIXME: Close partialImages from sender to avoid races here!
	case <-ctx.Done(): // Render cancelled
		return ctx.Err()
	default: // Render not cancelled
	}
	if err != nil {
		return err
	}
	if partialImages != nil {
		partialImages <- fullImg
	}
	return nil
}
