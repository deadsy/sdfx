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
		r.implState.ResInv /= 2
		if r.implState.ResInv < 1 {
			r.implState.ResInv = 1
		}
		r.implStateLock.Unlock()
		r.rerender()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract) {
		r.implStateLock.Lock()
		r.implState.ResInv *= 2
		if r.implState.ResInv > 64 {
			r.implState.ResInv = 64
		}
		r.implStateLock.Unlock()
		r.rerender()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		r.implStateLock.Lock()
		r.implState.DrawBbs = !r.implState.DrawBbs
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
			scale := 1 - wheelUpDown*r.implState.Bb.Size().Length2()*0.02 // Scale depending on current scale
			scale = math.Max(0.5, math.Min(2, scale))                     // Reasonable limits
			r.implState.Bb = r.implState.Bb.ScaleAboutCenter(scale)
			r.implStateLock.Unlock()
			r.rerender()
		}
		// Translation
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Save the cursor's position for previsualization and applying the final translation
			cx, cy := ebiten.CursorPosition()
			r.implStateLock.Lock()
			if r.translateFrom[0] == math.MaxInt { // Only if not already moving...
				r.translateFrom = sdf.V2i{cx, cy}
			}
			r.implStateLock.Unlock()
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			// Actually apply the translation and force a rerender
			cx, cy := ebiten.CursorPosition()
			r.implStateLock.Lock()
			changed := false
			if r.translateFrom[0] < math.MaxInt { // Only if already moving...
				clone := r.implState.Bb
				r.implState.Bb = r.implState.Bb.Translate(
					r.translateFrom.ToV2().Sub(sdf.V2i{cx, cy}.ToV2()).Mul(sdf.V2{X: 1, Y: -1}).
						Div(r.screenSize.ToV2()).Mul(r.implState.Bb.Size()))
				// Keep displacement until rerender is complete (avoid jump) using callback below + extra variable set here
				r.translateFromStop = sdf.V2i{cx, cy}
				changed = clone != r.implState.Bb
			}
			r.implStateLock.Unlock()
			if changed {
				r.rerender(func(err error) {
					r.implStateLock.Lock()
					r.translateFrom = sdf.V2i{math.MaxInt, math.MaxInt}
					r.translateFromStop = sdf.V2i{math.MaxInt, math.MaxInt}
					r.implStateLock.Unlock()
				})
			}
		}
		// Reset transform (100% of surface)
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			r.implStateLock.Lock()
			r.implState.Bb = toBox2(r.impl.BoundingBox()) // 100% zoom (impl2 will fix aspect ratio)
			r.implStateLock.Unlock()
			r.rerender()
		}
		// Color
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			r.implStateLock.Lock()
			r.implState.blackAndWhite = !r.implState.blackAndWhite
			r.implStateLock.Unlock()
			r.rerender()
		}
	//case 3: TODO
	default:
		panic("RendererState.onUpdateInputs not implemented for " + strconv.Itoa(r.impl.Dimensions()) + " dimensions")
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
	msg := fmt.Sprintf("TPS: %0.2f/%d\nResolution: %.2f [+/-]\nBlack/White: %t [C]\nShow boxes: %t [B]\nReset transform [R]",
		ebiten.CurrentTPS(), ebiten.MaxTPS(), 1/float64(r.implState.ResInv), r.implState.blackAndWhite, r.implState.DrawBbs)
	drawDefaultTextWithShadow(screen, msg, 5, r.screenSize[1]-5-16*strings.Count(msg, "\n"),
		color.RGBA{R: 0, G: 200, B: 0, A: 255})
}
