package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"image"
	"log"
	"math"
	"runtime/debug"
	"time"
)

func (r *Renderer) drawSDF(screen *ebiten.Image) {
	// Draw latest SDF render (and overlay the latest partial render)
	r.implStateLock.RLock()
	defer r.implStateLock.RUnlock()
	r.cachedRenderLock.RLock()
	defer r.cachedRenderLock.RUnlock()
	drawOpts := &ebiten.DrawImageOptions{}
	var tr sdf.V2
	if r.translateFrom[0] < math.MaxInt { // SDF2 special case: preview translations without rendering (until mouse release)
		cx, cy := ebiten.CursorPosition()
		if r.translateFromStop[0] < math.MaxInt {
			cx, cy = r.translateFromStop[0], r.translateFromStop[1]
		}
		tr = sdf.V2i{cx, cy}.ToV2().Sub(r.translateFrom.ToV2()).DivScalar(float64(r.implState.ResInv))
		// TODO: Place SDF2 render at the right location during special renders (zooming, changing resolution)
	}
	drawOpts.GeoM.Translate(tr.X, tr.Y)
	cachedRenderWidth, cachedRenderHeight := r.cachedRender.Size()
	drawOpts.GeoM.Scale(float64(r.screenSize[0])/float64(cachedRenderWidth), float64(r.screenSize[1])/float64(cachedRenderHeight))
	err := screen.DrawImage(r.cachedRender, drawOpts)
	if err != nil {
		panic(err) // Can this happen?
	}
	err = screen.DrawImage(r.cachedPartialRender, drawOpts)
	if err != nil {
		panic(err) // Can this happen?
	}
}

// rerender will discard any current rendering and start a new render (use it when something changes).
// It does not lock execution (renders in background).
func (r *Renderer) rerender(callbacks ...func(err error)) {
	r.cachedRenderLock.RLock()
	if r.cachedRender == nil {
		log.Println("Trying to render too soon (before first Update()). FIXME!")
		debug.PrintStack()
	}
	r.cachedRenderLock.RUnlock()
	go func(callbacks ...func(err error)) {
		var err error
		defer func() {
			for _, callback := range callbacks {
				callback(err)
			}
		}()
		if !r.renderingLock.TryLock(nil) {
			r.implStateLock.RLock() // Avoid race condition with creating a new context
			r.renderingCtxCancel()
			r.implStateLock.RUnlock()
			r.renderingLock.Lock() // Wait for previous render to finish (should be very fast)
		}
		defer r.renderingLock.Unlock()
		var renderCtx context.Context
		r.implStateLock.Lock()
		renderCtx, r.renderingCtxCancel = context.WithCancel(context.Background())
		r.implStateLock.Unlock()
		//partialRenders := make(chan *image.RGBA)
		//go func() {
		//	//imageBufCopy := make([]uint8, r.screenSize[0]*r.screenSize[1]*4) // RGBA
		//	for _ /*partialRender :*/ = range partialRenders {
		//		//r.cachedRenderLock.RLock()
		//		//copy(imageBufCopy, partialRender.Pix) // Is this locking rendering TPS?
		//		//r.cachedRenderLock.RUnlock()
		//		//gpuImg, err := ebiten.NewImageFromImage(&image.RGBA{ // Very slow! Keep outside locks!
		//		//	Pix:    imageBufCopy,
		//		//	Stride: r.screenSize[0] * 4,
		//		//	Rect:   image.Rect(0, 0, r.screenSize[0], r.screenSize[1]),
		//		//}, ebiten.FilterDefault)
		//		//if err != nil {
		//		//	log.Println("Error sending image to GPU:", err)
		//		//	continue
		//		//}
		//		//r.cachedRenderLock.Lock()
		//		//r.cachedPartialRender = gpuImg
		//		//r.cachedRenderLock.Unlock()
		//	}
		//}()
		renderStartTime := time.Now()
		r.implStateLock.RLock()
		renderSize := r.screenSize.ToV2().DivScalar(float64(r.implState.ResInv)).ToV2i()
		sameSize := r.cachedRenderCpu != nil && (sdf.V2i{r.cachedRenderCpu.Rect.Max.X, r.cachedRenderCpu.Rect.Max.Y} == renderSize)
		if !sameSize {
			r.cachedRenderCpu = image.NewRGBA(image.Rect(0, 0, renderSize[0], renderSize[1]))
		}
		r.implStateLock.RUnlock()
		r.implLock.RLock()
		err = r.impl.Render(renderCtx, r.implState, r.implStateLock, r.cachedRenderLock, nil, r.cachedRenderCpu)
		if err != nil {
			if err != context.Canceled {
				log.Println("[DevRenderer] Error rendering:", err)
			}
			return
		}
		r.implLock.RUnlock()
		renderGPUStartTime := time.Now()
		renderGpuImg, err := ebiten.NewImageFromImage(r.cachedRenderCpu, ebiten.FilterDefault)
		if err != nil {
			log.Println("Error sending image to GPU:", err)
			return
		}
		log.Println("[DevRenderer] CPU Render took:", renderGPUStartTime.Sub(renderStartTime), "- Sending to GPU took:", time.Since(renderGPUStartTime))
		r.implLock.RLock()
		r.implStateLock.RLock() // WARNING: Locking order (to avoid deadlocks)
		if r.impl.Dimensions() == 2 {
			r.cachedRenderBb2 = r.implState.Bb
		} else {
			r.cachedRenderBb2 = sdf.Box2{}
		}
		r.implStateLock.RUnlock()
		r.implLock.RUnlock()
		r.cachedRenderLock.Lock()
		// Reuse the previous render for the parts that did not change
		if !sameSize {
			// Need to resize the rendering result: overwrite
			r.cachedRender = renderGpuImg
		} else {
			// No need to resize render result: draw over it in case we implement skipping unneeded parts of the image in the future
			err = r.cachedRender.DrawImage(renderGpuImg, &ebiten.DrawImageOptions{})
			if err != nil {
				log.Println("Error sending image to GPU:", err)
				return
			}
		}
		r.cachedRenderLock.Unlock()
	}(callbacks...)
}
