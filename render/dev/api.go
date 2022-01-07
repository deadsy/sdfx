package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"math"
	"sync"
)

// devRendererImpl is the interface implemented by the SDF2 and SDF3 renderers.
// Note that the implementation is independent of the graphics backend used and renders CPU images.
type devRendererImpl interface {
	// Dimensions are 2 for SDF2 and 3 for SDF3
	Dimensions() int
	// BoundingBox returns the full bounding box of the surface (Z is ignored for SDF2)
	BoundingBox() sdf.Box3
	// ColorModes returns the number of color modes supported
	ColorModes() int
	// Render performs a full render, given the screen size (it may be cancelled using the given context).
	// Returns partially rendered images as progress is made through partialImages (if non-nil, channel closed).
	Render(ctx context.Context, state *RendererState, stateLock, cachedRenderLock *sync.RWMutex,
		partialRender chan<- *image.RGBA, fullRender *image.RGBA) error
	// TODO: Map clicks to source code? (using reflection on the SDF and profiling/code generation?)
}

// RendererState is an internal structure, exported for (de)serialization reasons
type RendererState struct {
	// SHARED
	ResInv    int  // How detailed is the image: number screen pixels for each pixel rendered (SDF2: use a power of two)
	DrawBbs   bool // Whether to show all bounding boxes (useful for debugging subtraction/intersection of SDFs) TODO
	ColorMode int  // The color mode (each render may support multiple modes)
	// SDF2
	Bb sdf.Box2 // Controls the scale and displacement
	// SDF3
	CamCenter                 sdf.V3  // Arc-Ball camera center (the point we are looking at)
	CamYaw, CamPitch, CamDist float64 // Arc-Ball rotation angles (around CamCenter) and distance from CamCenter
}

func (r *Renderer) newRendererState() *RendererState {
	r.implLock.RLock()
	defer r.implLock.RUnlock()
	s := &RendererState{
		// TODO: Guess a ResInv based on rendering performance
		ResInv: 8,

		Bb: toBox2(r.impl.BoundingBox()), // 100% zoom (will fix aspect ratio later)

	}
	s.ResetCam3(r)
	return s
}

func (s *RendererState) Cam3MatrixNoTranslation() sdf.M44 {
	return sdf.RotateZ(-s.CamYaw).Mul(sdf.RotateX(-s.CamPitch))
}

func (s *RendererState) ResetCam3(r *Renderer) {
	s.CamCenter = r.impl.BoundingBox().Center()
	s.CamDist = r.impl.BoundingBox().Size().Length() / 2
	s.CamPitch = -math.Pi / 4 // Look from 45º up
	s.CamYaw = -math.Pi / 4   // Look from 45º right
}
