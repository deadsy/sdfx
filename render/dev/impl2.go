package dev

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"image/color"
	"math"
)

//-----------------------------------------------------------------------------
// CONFIGURATION
//-----------------------------------------------------------------------------

// Opt2Cam sets the default camera for SDF2 (may grow to follow the aspect ratio of the screen).
// WARNING: Need to run again the main renderer to apply a change of this option.
func Opt2Cam(bb sdf.Box2) Option {
	return func(r *Renderer) {
		r.implState.Bb = bb
	}
}

// Opt2EvalRange skips the initial scan of the SDF2 to find the minimum and maximum value, and can also be used to
// make the surface easier to see by setting them to a value close to 0.
func Opt2EvalRange(min, max float64) Option {
	return func(r *Renderer) {
		if r2, ok := r.impl.(*renderer2); ok {
			r2.evalMin = min
			r2.evalMax = max
		}
	}
}

// Opt2EvalScanCells configures the initial scan of the SDF2 to find minimum and maximum values (defaults to 128x128 cells).
func Opt2EvalScanCells(cells sdf.V2i) Option {
	return func(r *Renderer) {
		if r2, ok := r.impl.(*renderer2); ok {
			r2.evalScanCells = cells
		}
	}
}

//-----------------------------------------------------------------------------
// RENDERER
//-----------------------------------------------------------------------------

type renderer2 struct {
	s                sdf.SDF2 // The SDF to render
	pixelsRand       []int    // Cached set of pixels in random order to avoid shuffling (reset on recompilation and resolution changes)
	evalMin, evalMax float64  // The pre-computed minimum and maximum of the whole surface (for stable colors and speed)
	evalScanCells    sdf.V2i
}

func newDevRenderer2(s sdf.SDF2) devRendererImpl {
	r := &renderer2{
		s:             s,
		evalScanCells: sdf.V2i{128, 128},
	}
	return r
}

func (r *renderer2) Dimensions() int {
	return 2
}

func (r *renderer2) BoundingBox() sdf.Box3 {
	bb := r.s.BoundingBox()
	return sdf.Box3{Min: bb.Min.ToV3(0), Max: bb.Max.ToV3(0)}
}

func (r *renderer2) ColorModes() int {
	// 0: Gradient (useful for debugging sides)
	// 1: Black/white (clearer surface boundary)
	return 2
}

func (r *renderer2) Render(args *renderArgs) error {
	if r.evalMin == 0 && r.evalMax == 0 { // First render (ignoring external cache)
		// Compute minimum and maximum evaluate values for a shared color scale for all blocks
		r.evalMin, r.evalMax = utilSdf2MinMax(r.s, r.s.BoundingBox(), r.evalScanCells)
		//log.Println("MIN:", r.evalMin, "MAX:", r.evalMax)
	}

	// Maintain Bb aspect ratio on ResInv change, increasing the size as needed
	args.stateLock.Lock()
	fullRenderSize := args.fullRender.Bounds().Size()
	bbAspectRatio := args.state.Bb.Size().X / args.state.Bb.Size().Y
	screenAspectRatio := float64(fullRenderSize.X) / float64(fullRenderSize.Y)
	if math.Abs(bbAspectRatio-screenAspectRatio) > 1e-12 {
		if bbAspectRatio > screenAspectRatio {
			scaleYBy := bbAspectRatio / screenAspectRatio
			args.state.Bb = sdf.NewBox2(args.state.Bb.Center(), args.state.Bb.Size().Mul(sdf.V2{X: 1, Y: scaleYBy}))
		} else {
			scaleXBy := screenAspectRatio / bbAspectRatio
			args.state.Bb = sdf.NewBox2(args.state.Bb.Center(), args.state.Bb.Size().Mul(sdf.V2{X: scaleXBy, Y: 1}))
		}
	}
	args.stateLock.Unlock()

	// Apply color mode
	evalMin, evalMax := r.evalMin, r.evalMax
	if args.state.ColorMode == 1 { // Force black and white to see the surface better
		evalMin, evalMax = -1e-12, 1e-12
	}

	// Perform the actual render
	return implCommonRender(func(pixel sdf.V2i, pixel01 sdf.V2) interface{} { return nil },
		func(pixel sdf.V2i, pixel01 sdf.V2, job interface{}) *jobResult {
			pixel01.Y = 1 - pixel01.Y // Inverted Y
			args.stateLock.RLock()
			pos := args.state.Bb.Min.Add(pixel01.Mul(args.state.Bb.Size()))
			args.stateLock.RUnlock()
			grayVal := render.ImageColor2(r.s.Evaluate(pos), evalMin, evalMax)
			return &jobResult{
				pixel: pixel,
				color: color.RGBA{R: uint8(grayVal * 255), G: uint8(grayVal * 255), B: uint8(grayVal * 255), A: 255},
			}
		}, args, &r.pixelsRand)

	// TODO: Draw bounding boxes over the image
}
