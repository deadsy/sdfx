package dev

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
)

var _ ebiten.Game = &Renderer{}

func (r *Renderer) Update(_ *ebiten.Image) error {
	var err error
	// FIXME: Debug TPS reduction when rendering (probably locking + forced queued main thread loading of images)
	r.cachedRenderLock.RLock()
	firstFrame := r.cachedRender == nil
	if firstFrame { // This always runs before the first frame
		r.cachedRender, err = ebiten.NewImage(1, 1, ebiten.FilterDefault)
		if err != nil {
			return err
		}
		r.cachedPartialRender = r.cachedRender
	}
	r.cachedRenderLock.RUnlock()
	r.onUpdateInputs()
	return nil
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	r.drawSDF(screen)
	r.drawUI(screen)
}

func (r *Renderer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	r.cachedRenderLock.RLock()
	firstFrame := r.cachedRender == nil
	r.cachedRenderLock.RUnlock()
	r.implStateLock.RLock()
	defer r.implStateLock.RUnlock()
	if !firstFrame { // Layout is called before Update(), but don't render in this case
		newScreenSize := sdf.V2i{outsideWidth, outsideHeight}
		if r.screenSize != newScreenSize {
			r.screenSize = newScreenSize
			r.rerender()
		}
	}
	return outsideWidth, outsideHeight // Use all available pixels, no re-scaling (unless ResInv is modified)
}
