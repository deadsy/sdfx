package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"sync"
)

// devRendererImpl is the interface implemented by the SDF2 and SDF3 renderers
type devRendererImpl interface {
	// Dimensions are 2 for SDF2 and 3 for SDF3: determines how to update the DevRendererState on input
	Dimensions() int
	// BoundingBox returns the full bounding box of the surface (Z is ignored for SDF2)
	BoundingBox() sdf.Box3
	// Render performs a full render, given the screen size (it may be cancelled using the given context).
	// Returns partially rendered images as progress is made through partialImages (if non-nil, channel not closed).
	Render(ctx context.Context, screenSize sdf.V2i, state *DevRendererState,
		stateLock, cachedRenderLock *sync.RWMutex, partialRender chan<- *image.RGBA) (*image.RGBA, error)
	// TODO: Click to identify SDFs affecting the touched surface (using reflection on devRendererImpl and some way to identify SDFs in source code)
}

// DevRendererState is an internal structure, exported for (de)serialization reasons
type DevRendererState struct {
	// SHARED
	Resolution int  // How detailed is the image: number screen pixels for each pixel rendered (SDF2: use a power of two)
	DrawBbs    bool // Whether to show all bounding boxes (useful for debugging subtraction of SDFs) TODO
	// SDF2
	Bb sdf.Box2 // Controls the scale and displacement
	// SDF3
	CamCenter                     sdf.V3  // Arc-Ball camera center
	CamThetaX, CamThetaY, CamDist float64 // Arc-Ball camera angles and distance
}

func newDevRendererState() *DevRendererState {
	return &DevRendererState{
		// TODO: Guess a Resolution based on rendering performance
		Resolution: 1,
	}
}
