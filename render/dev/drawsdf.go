package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"log"
	"math"
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
		tr = sdf.V2i{cx, cy}.ToV2().Sub(r.translateFrom.ToV2())
		// TODO: Place SDF2 render at the right location during special renders (zooming, changing resolution),
		// TODO: Also, skip rendering unneeded parts (or blocks) of the image
	}
	drawOpts.GeoM.Translate(tr.X, tr.Y)
	cachedRenderWidth, cachedRenderHeight := r.cachedRender.Size()
	drawOpts.GeoM.Scale(float64(r.screenSize[0])/float64(cachedRenderWidth), float64(r.screenSize[1])/float64(cachedRenderHeight))
	err := screen.DrawImage(r.cachedRender, drawOpts)
	if err != nil {
		panic(err) // Can this happen?
	}
	drawOpts.GeoM.Reset()
	drawOpts.GeoM.Translate(tr.X, tr.Y)
	err = screen.DrawImage(r.cachedPartialRender, drawOpts)
	if err != nil {
		panic(err) // Can this happen?
	}
}

// rerender will discard any current rendering and start a new render (use it when something changes).
// It does not lock execution (renders in background).
func (r *Renderer) rerender() {
	if r.cachedRender == nil {
		log.Println("Trying to render too soon (before first Update()). FIXME!")
	}
	go func() {
		if !r.renderingLock.TryLock(nil) {
			// This is OK because the previous render may have finished between the previous and next instruction,
			// but calling cancel on an unused context is still OK, and the next lock will succeed.
			r.renderingCtxCancel()
			r.renderingLock.Lock() // Wait for previous render to finish (should be very fast)
		}
		defer r.renderingLock.Unlock()
		var renderCtx context.Context
		renderCtx, r.renderingCtxCancel = context.WithCancel(context.Background())
		renderStartTime := time.Now()
		partialRenders := make(chan *ebiten.Image)
		go func() {
			for partialRender := range partialRenders {
				r.cachedRenderLock.Lock()
				r.cachedPartialRender = partialRender
				r.cachedRenderLock.Unlock()
			}
		}()
		render, err := r.impl.Render(renderCtx, r.screenSize, r.implState, r.implStateLock, r.cachedRenderLock, partialRenders)
		close(partialRenders)
		if err != nil {
			if err != context.Canceled {
				log.Println("[DevRenderer] Error rendering:", err)
			}
			return
		}
		log.Println("[DevRenderer] Render took", time.Since(renderStartTime))
		r.implStateLock.RLock() // WARNING: Locking order (to avoid deadlocks)
		r.cachedRenderLock.Lock()
		if r.impl.Dimensions() == 2 {
			r.cachedRenderBb2 = r.implState.bb
		} else {
			r.cachedRenderBb2 = sdf.Box2{}
		}
		// Reuse the previous render for the parts that did not change
		if sX, sY := r.cachedRender.Size(); sX == 1 && sY == 1 {
			r.cachedRender = render
		} else {
			err = r.cachedRender.DrawImage(render, nil)
			if err != nil {
				log.Println("[DevRenderer] Error rendering (DrawImage):", err)
			}
		}
		r.cachedRenderLock.Unlock()
		r.implStateLock.RUnlock()
	}()
}
