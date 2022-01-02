package dev

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
)

var _ ebiten.Game = &Renderer{}

func (r *Renderer) Update(_ *ebiten.Image) error {
	var err error
	r.cachedRenderLock.RLock()
	if r.cachedRender == nil { // This always runs before the first frame
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
	r.implStateLock.RLock()
	newScreenSize := sdf.V2i{outsideWidth, outsideHeight}
	if r.screenSize != newScreenSize {
		// Reuse previous render if we are zooming out (making the resolution smaller)
		r.screenSize = newScreenSize
		r.implStateLock.RUnlock()
		r.rerender()
	} else {
		r.implStateLock.RUnlock()
	}
	return outsideWidth, outsideHeight // Use all available pixels, no re-scaling (unless resolution is modified)
}
