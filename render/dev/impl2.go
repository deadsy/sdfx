package dev

import (
	"context"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"image/color"
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

func (r *devRenderer2) Render(ctx context.Context, screenSize sdf.V2i, state *devRendererState, stateLock, cachedRenderLock *sync.RWMutex, partialImages chan<- *ebiten.Image) (*ebiten.Image, error) {
	if r.evalMin == 0 && r.evalMax == 0 { // First render (ignoring external cache)
		// Compute minimum and maximum evaluate values for a shared color scale for all blocks
		r.evalMin, r.evalMax = utilSdf2MinMax(r.s, r.s.BoundingBox(), sdf.V2i{64, 64} /* TODO: Configurable? */)
		//log.Println("MIN:", r.evalMin, "MAX:", r.evalMax)
	}
	stateLock.Lock()
	if state.resolution == 0 { // First render as a SDF2 (considering external cache)
		if state.bb.Size().Length2() == 0 {
			state.bb = r.s.BoundingBox() // 100% zoom (will fix aspect ratio later)
		}
		// TODO: Guess a resolution based on rendering performance
		state.resolution = 8
	}

	// Maintain bb aspect ratio on resolution change, increasing the size as needed
	bbAspectRatio := state.bb.Size().X / state.bb.Size().Y
	screenAspectRatio := float64(screenSize[0]) / float64(screenSize[1])
	if math.Abs(bbAspectRatio-screenAspectRatio) > 1e-12 {
		if bbAspectRatio > screenAspectRatio {
			scaleYBy := bbAspectRatio / screenAspectRatio
			state.bb = sdf.NewBox2(state.bb.Center(), state.bb.Size().Mul(sdf.V2{X: 1, Y: scaleYBy}))
		} else {
			scaleXBy := screenAspectRatio / bbAspectRatio
			state.bb = sdf.NewBox2(state.bb.Center(), state.bb.Size().Mul(sdf.V2{X: scaleXBy, Y: 1}))
		}
	}
	stateLock.Unlock()

	// Create the new full GPU image
	fullImgSize := screenSize
	fullImg, err := ebiten.NewImage(fullImgSize[0], fullImgSize[1], ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	err = fullImg.Fill(color.Transparent) // Set to transparent for non-rendered blocks (to easily display over previous render)
	if err != nil {
		return nil, err
	}

	// Render each blockIndex of the image individually to allow cancelling the render
	pixelsPerBlock := sdf.V2i{128, 128} /* TODO: Configurable? */
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
	for i := 0; i < jobCount; i++ {
		err = <-errors
		if err != nil {
			return fullImg, err // Quick exit on error
		}
	}
	return fullImg, err
}

func (r *devRenderer2) renderBlock(ctx context.Context, fullImg *ebiten.Image, blockIndex sdf.V2i,
	pixelsPerBlock sdf.V2i, state *devRendererState, stateLock, cachedRenderLock *sync.RWMutex, numBlocks sdf.V2i, fullImgSize sdf.V2i,
	partialImages chan<- *ebiten.Image) error {
	select {
	case <-ctx.Done(): // Render cancelled
		return ctx.Err()
	default: // Render not cancelled
	}

	stateLock.RLock()
	defer stateLock.RUnlock()

	// Compute positions and sizes
	blockStartPixel := blockIndex.ToV2().Mul(pixelsPerBlock.ToV2()).ToV2i()
	blockSizePixels := pixelsPerBlock.AddScalar(nextPowerOf2(state.resolution))
	if blockIndex[0] == numBlocks[0]-1 {
		blockSizePixels[0] = fullImgSize[0] - blockStartPixel[0] + nextPowerOf2(state.resolution)
	}
	if blockIndex[1] == numBlocks[1]-1 { // Inverted Y
		blockSizePixels[1] = fullImgSize[1] - blockStartPixel[1] + nextPowerOf2(state.resolution)
	}
	blockSizePixels = blockSizePixels.ToV2().DivScalar(float64(state.resolution)).ToV2i()
	if blockSizePixels[0] == 0 || blockSizePixels[1] == 0 {
		return nil // Empty block ignored
	}
	blockBb := sdf.Box2{
		Min: state.bb.Min.Add(state.bb.Size().Mul(blockStartPixel.ToV2().Div(fullImgSize.ToV2()))),
		Max: state.bb.Min.Add(state.bb.Size().Mul(
			blockStartPixel.Add(blockSizePixels.ToV2().MulScalar(float64(state.resolution)).ToV2i()).ToV2().Div(fullImgSize.ToV2()))),
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

	// Move to GPU and draw over the other GPU image (ebiten optimizes several draws over the same image)
	blockImgGpu, err := ebiten.NewImageFromImage(blockImg, ebiten.FilterDefault)
	drawOptions := &ebiten.DrawImageOptions{Filter: ebiten.FilterDefault /* Linear (blurry) drawing causes black lines to merge */}
	translateY := fullImgSize[1] - blockStartPixel[1] - blockSizePixels[1]*state.resolution // Y is inverted
	drawOptions.GeoM.Scale(float64(state.resolution), float64(state.resolution))            // Apply resolution
	drawOptions.GeoM.Translate(float64(blockStartPixel[0]), float64(translateY+nextPowerOf2(state.resolution)))
	if err != nil {
		return err
	}
	cachedRenderLock.Lock()
	err = fullImg.DrawImage(blockImgGpu, drawOptions)
	cachedRenderLock.Unlock()
	select {
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
