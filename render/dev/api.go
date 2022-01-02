package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"sync"
)

// devRendererImpl is the interface implemented by the SDF2 and SDF3 renderers
type devRendererImpl interface {
	// Dimensions are 2 for SDF2 and 3 for SDF3: determines how to update the devRendererState on input
	Dimensions() int
	// BoundingBox returns the full bounding box of the surface (Z is ignored for SDF2)
	BoundingBox() sdf.Box3
	// Render performs a full render, given the screen size (it may be cancelled using the given context).
	// Returns partially rendered images as progress is made through partialImages (if non-nil, channel not closed).
	Render(ctx context.Context, screenSize sdf.V2i, state *devRendererState,
		stateLock, cachedRenderLock *sync.RWMutex, partialRender chan<- *ebiten.Image) (*ebiten.Image, error)
}

type devRendererState struct {
	// SHARED
	resolution int  // How detailed is the image: number screen pixels for each pixel rendered (SDF2: use a power of two)
	drawBbs    bool // Whether to show all bounding boxes (useful for debugging subtraction of SDFs) TODO
	// SDF2
	bb sdf.Box2 // Controls the scale and displacement
	// SDF3
	camCenter                     sdf.V3  // Arc-Ball camera center
	camThetaX, camThetaY, camDist float64 // Arc-Ball camera angles and distance
}

func newDevRendererState() *devRendererState {
	return &devRendererState{}
}
