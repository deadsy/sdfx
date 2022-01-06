package dev

import (
	"context"
	"fmt"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"log"
	"math"
	"strconv"
	"time"
)

// onUpdateInputs handles inputs
func (r *Renderer) onUpdateInputs() {
	r.implLock.RLock()
	defer r.implLock.RUnlock()
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
	r.implStateLock.RLock()
	implDimCache := r.implDimCache
	r.implStateLock.RUnlock()
	switch implDimCache {
	case 2:
		// Zooming
		_, wheelUpDown := ebiten.Wheel()
		if wheelUpDown != 0 {
			r.implStateLock.Lock()
			scale := 1 - wheelUpDown*r.implState.Bb.Size().Length2()*0.02   // Scale depending on current scale
			scale = math.Max(1/r.zoomFactor, math.Min(r.zoomFactor, scale)) // Apply zoom limits
			r.implState.Bb = r.implState.Bb.ScaleAboutCenter(scale)
			r.implStateLock.Unlock()
			r.rerender()
		}
		// Translation
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
			// Save the cursor's position for previsualization and applying the final translation
			cx, cy := ebiten.CursorPosition()
			r.implStateLock.Lock()
			if r.translateFrom[0] == math.MaxInt { // Only if not already moving...
				r.translateFrom = sdf.V2i{cx, cy}
			}
			r.implStateLock.Unlock()
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
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
		// Reset camera transform (100% of surface)
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			r.implStateLock.Lock()
			r.implState.Bb = toBox2(r.impl.BoundingBox()) // 100% zoom (impl2 will fix aspect ratio)
			r.implStateLock.Unlock()
			r.rerender()
		}
		// Color
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			r.implStateLock.Lock()
			r.implState.ColorMode = !r.implState.ColorMode
			r.implStateLock.Unlock()
			r.rerender()
		}
	case 3:
		// Zooming
		_, wheelUpDown := ebiten.Wheel()
		if wheelUpDown != 0 {
			r.implStateLock.Lock()
			scale := 1 - wheelUpDown*100
			scale = math.Max(1/r.zoomFactor, math.Min(r.zoomFactor, scale)) // Apply zoom limits
			r.implState.CamDist *= scale
			r.implStateLock.Unlock()
			r.rerender()
		}
		// Rotation
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
			// Save the cursor's position for previsualization and applying the final translation
			cx, cy := ebiten.CursorPosition()
			r.implStateLock.Lock()
			if r.translateFrom[0] == math.MaxInt { // Only if not already moving...
				r.translateFrom = sdf.V2i{cx, cy}
			}
			r.implStateLock.Unlock()
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
			// Actually apply the translation and force a rerender
			cx, cy := ebiten.CursorPosition()
			r.implStateLock.Lock()
			changed := false
			if r.translateFrom[0] < math.MaxInt { // Only if already moving...
				r.implState.CamYaw += float64(cx-r.translateFrom[0]) / 100 // TODO: Proper delta computation
				if r.implState.CamYaw < -math.Pi {
					r.implState.CamYaw += 2 * math.Pi // Limits (wrap around)
				} else if r.implState.CamYaw > math.Pi {
					r.implState.CamYaw -= 2 * math.Pi // Limits (wrap around)
				}
				r.implState.CamPitch += -float64(cy-r.translateFrom[1]) / 100
				r.implState.CamPitch = math.Max(-math.Pi/2, math.Min(math.Pi/2, r.implState.CamPitch))
				log.Println("New camera rotation (pitch", r.implState.CamPitch, "yaw", r.implState.CamYaw, ")")
				// Keep displacement until rerender is complete (avoid jump) using callback below + extra variable set here
				r.translateFromStop = sdf.V2i{cx, cy}
				changed = true
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
		// TODO
	default:
		panic("RendererState.onUpdateInputs not implemented for " + strconv.Itoa(r.implDimCache) + " dimensions")
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
		drawDefaultTextWithShadow(screen, "Rendering...", 5, 5+12, color.RGBA{R: 255, A: 255})
	}

	// Draw current state and controls
	r.implStateLock.RLock()
	defer r.implStateLock.RUnlock()
	msgFmt := "TPS: %0.2f/%d\nResolution: %.2f [+/-]\nColor: %t [C]\nShow boxes: %t [B]\nReset camera [R]"
	msgValues := []interface{}{ebiten.CurrentTPS(), ebiten.MaxTPS(), 1 / float64(r.implState.ResInv), r.implState.ColorMode, r.implState.DrawBbs}
	switch r.implDimCache {
	case 2:
		msgFmt = "SDF2 Renderer\n=============\n" + msgFmt + "\nTranslate cam [MiddleMouse]\nZoom cam [MouseWheel]"
	case 3:
		msgFmt = "SDF3 Renderer\n=============\n" + msgFmt + "\nRotate cam [MiddleMouse]\nTranslate cam [Shift+MiddleMouse]\nZoom cam [MouseWheel]"
	}
	msg := fmt.Sprintf(msgFmt, msgValues...)
	boundString := text.BoundString(defaultFont, msg)
	drawDefaultTextWithShadow(screen, msg, 5, r.screenSize[1]-boundString.Size().Y+10, color.RGBA{G: 255, A: 255})
}
