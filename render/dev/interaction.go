package dev

import (
	"context"
	"fmt"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"
)

// onUpdateInputs handles inputs
func (r *Renderer) onUpdateInputs() {
	// SHARED CONTROLS
	if inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) {
		r.implStateLock.Lock()
		r.implState.resolution *= 2
		if r.implState.resolution > 64 {
			r.implState.resolution = 64
		}
		r.implStateLock.Unlock()
		r.rerender()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract) {
		r.implStateLock.Lock()
		r.implState.resolution /= 2
		if r.implState.resolution < 1 {
			r.implState.resolution = 1
		}
		r.implStateLock.Unlock()
		r.rerender()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		r.implStateLock.Lock()
		r.implState.drawBbs = !r.implState.drawBbs
		r.implStateLock.Unlock()
		r.rerender()
	}
	// SDF2/SDF3-SPECIFIC CONTROLS
	switch r.impl.Dimensions() {
	case 2:
		// Zooming
		_, wheelUpDown := ebiten.Wheel()
		if wheelUpDown != 0 {
			r.implStateLock.Lock()
			scale := 1 - wheelUpDown*r.implState.bb.Size().Length2()*0.02 // Scale depending on current scale
			scale = math.Max(0.5, math.Min(2, scale))                     // Reasonable limits
			r.implState.bb = r.implState.bb.ScaleAboutCenter(scale)
			r.implStateLock.Unlock()
			r.rerender() // Reuse previous render if we are zooming out!
		}
		// Translation
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Save the cursor's position for previsualization and applying the final translation
			cx, cy := ebiten.CursorPosition()
			r.implStateLock.Lock()
			r.translateFrom = sdf.V2i{cx, cy}
			r.implStateLock.Unlock()
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			// Actually apply the translation and force a rerender
			cx, cy := ebiten.CursorPosition()
			r.implStateLock.Lock()
			clone := r.implState.bb
			r.implState.bb = r.implState.bb.Translate(r.translateFrom.ToV2().Sub(sdf.V2i{cx, cy}.ToV2()).Mul(sdf.V2{X: 1, Y: -1}).
				Div(r.screenSize.ToV2()).Mul(r.implState.bb.Size()))
			changed := clone != r.implState.bb
			r.translateFrom = sdf.V2i{math.MaxInt, math.MaxInt}
			r.implStateLock.Unlock()
			if changed {
				r.rerender()
			}
		}
		// Reset transform (100% of surface)
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			r.implStateLock.Lock()
			r.implState.bb = toBox2(r.impl.BoundingBox()) // 100% zoom (impl2 will fix aspect ratio)
			r.implStateLock.Unlock()
			r.rerender()
		}
	//case 3: TODO
	default:
		panic("devRendererState.onUpdateInputs not implemented for " + strconv.Itoa(r.impl.Dimensions()) + " dimensions")
	}
}

// ControlsText returns the help text
func (r *Renderer) drawUI(screen *ebiten.Image) {
	// Notify when rendering
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancelFunc()
	if r.renderingLock.RTryLock(ctx) {
		r.renderingLock.RUnlock()
	} else {
		drawDefaultTextWithShadow(screen, "Rendering...", 5, 5+12, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	}

	// Draw current state and controls
	r.implStateLock.Lock()
	defer r.implStateLock.Unlock()
	msg := fmt.Sprintf("TPS: %0.2f/%d\nResolution: %d [+/-]\nShow boxes: %t [B]\nReset transform [R]",
		ebiten.CurrentTPS(), ebiten.MaxTPS(), r.implState.resolution, r.implState.drawBbs)
	drawDefaultTextWithShadow(screen, msg, 5, r.screenSize[1]-5-16*strings.Count(msg, "\n"),
		color.RGBA{R: 0, G: 200, B: 0, A: 255})
}
