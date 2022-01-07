package dev

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"github.com/subchen/go-trylock/v2"
	"image"
	"math"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Renderer is a SDF2/SDF3 renderer intended for fast development iterations that renders directly to a window.
// The first time, it starts the renderer process. It also starts listening for code changes.
// When a code change is detected, the app is recompiled (taking advantage of go's fast compilation times)
// by the renderer and communicates directly to the renderer, providing the new surface data to the previous window.
//
// It allows very fast SDF updates (saving camera position) whenever the code changes, speeding up the modelling process.
// The renderer is mainly CPU-based (with a resolution parameter to control speed vs quality), as sdfx is also CPU-based.
// The scene is only rendered when something changes, as rendering SDFs with good quality is not real-time.
//
// The SDF2 renderer is based on the PNG renderer, showing the image directly on screen (without creating the PNG file).
// The camera can be moved and scaled (using the mouse), rendering only the interesting part of the SDF.
//
// SDF3s are raycasted from a perspective arc-ball camera that can be rotated around a pivot point, move its pivot and
// move closer or farther away from the pivot (using Blender-like mouse controls).
// Note that only the shown surface is actually rendered thanks to raycasting from the camera.
// This also means that the resulting surface can be much more detailed (depending on chosen resolution)
// than the triangle meshes generated by standard renderers.
//
// TODO: Once merged, use max-resolution runtime-computed VoxelSdf2 and VoxelSdf3 to accelerate camera movements
//
// It uses [ebiten](https://github.com/hajimehoshi/ebiten) for rendering, which is cross-platform, so it could also
// be used to showcase a surface (without automatic updates) creating an application for desktop, web or mobile.
type Renderer struct {
	impl                devRendererImpl   // the implementation to use SDF2/SDF3/remote process.
	implDimCache        int               // the number of dimensions of impl (cached to avoid remote calls every frame)
	implLock            *sync.RWMutex     // the implementation lock
	implState           *RendererState    // the renderer's state, so impl can be swapped while keeping the state.
	implStateLock       *sync.RWMutex     // the renderer's state lock
	cachedRender        *ebiten.Image     // the latest cached render (to avoid rendering every frame, or frame parts even if nothing changed)
	cachedRenderCpu     *image.RGBA       // the latest cached render (to avoid rendering every frame, or frame parts even if nothing changed)
	cachedRenderBb2     sdf.Box2          // what part of the SDF2 the latest cached render represents (not implemented, and no equivalent optimization available for SDF3s)
	cachedPartialRender *ebiten.Image     // the latest partial render (to display render progress visually)
	cachedRenderLock    *sync.RWMutex     // the lock over tha partial render
	screenSize          sdf.V2i           // the screen ResInv
	renderingCtxCancel  func()            // non-nil if we are currently rendering
	renderingLock       trylock.TryLocker // locked when we are rendering, use renderingCtx to cancel the previous render
	translateFrom       sdf.V2i           // Translate/rotate (for 3D) screen space start
	translateFromStop   sdf.V2i           // Translate/rotate (for 3D) screen space end (recorded while processing the new frame)
	// Static configuration
	runCmd             func() *exec.Cmd // generates a new command to compile and run the code for the new SDF
	watchFiles         []string         // the files to watch for recompilation of new code
	backOff            backoff.BackOff  // the backoff to connect to the new process after recompilation
	partialRenderEvery time.Duration    // how much time to wait between partial render updates to screen
	zoomFactor         float64          // how much to scale the SDF2/SDF3 on each zoom operation (> 1)
}

// NewDevRenderer see DevRenderer
func NewDevRenderer(anySDF interface{}, opts ...Option) *Renderer {
	r := &Renderer{
		implLock:          &sync.RWMutex{},
		implStateLock:     &sync.RWMutex{},
		cachedRenderLock:  &sync.RWMutex{},
		renderingLock:     trylock.New(),
		translateFrom:     sdf.V2i{math.MaxInt, math.MaxInt},
		translateFromStop: sdf.V2i{math.MaxInt, math.MaxInt},
		// Configuration
		runCmd: func() *exec.Cmd {
			return exec.Command("go", "run", "-v", ".")
		},
		watchFiles:         []string{"."},
		backOff:            backoff.NewExponentialBackOff(),
		partialRenderEvery: time.Second,
		zoomFactor:         1.25,
	}
	r.backOff.(*backoff.ExponentialBackOff).InitialInterval = 10 * time.Millisecond
	switch s := anySDF.(type) {
	case sdf.SDF2:
		r.impl = newDevRenderer2(s)
		r.cachedRenderBb2 = s.BoundingBox()
	case sdf.SDF3:
		r.impl = newDevRenderer3(s)
	default:
		panic("anySDF must be either a SDF2 or a SDF3")
	}
	r.implDimCache = r.impl.Dimensions()
	r.implState = r.newRendererState()
	// Apply all configuration options
	for _, opt := range opts {
		opt(r)
	}
	return r
}

const RequestedAddressEnvKey = "SDFX_DEV_RENDERER_CHILD"

func (r *Renderer) Run() error {
	requestedAddress := os.Getenv(RequestedAddressEnvKey)
	if requestedAddress != "" { // Found a parent renderer (environment variable)
		return r.runChild(requestedAddress)
	} else { // Otherwise, listen for code changes to spawn a child renderer and create the local renderer
		return r.runRenderer(r.runCmd, r.watchFiles)
	}
}
