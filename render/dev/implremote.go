package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"log"
	"net/rpc"
	"sync"
)

// devRendererClient implements devRendererImpl by calling a remote implementation (using Go's net/rpc)
type devRendererClient struct {
	cl *rpc.Client
}

// newDevRendererClient see devRendererClient
func newDevRendererClient(client *rpc.Client) devRendererImpl {
	return &devRendererClient{cl: client}
}

func (d *devRendererClient) Dimensions() int {
	var out int
	err := d.cl.Call("DevRendererService.Dimensions", &out, &out)
	if err != nil {
		log.Println("Error on remote call:", err)
	}
	return out
}

func (d *devRendererClient) BoundingBox() sdf.Box3 {
	var out sdf.Box3
	err := d.cl.Call("DevRendererService.BoundingBox", &out, &out)
	if err != nil {
		log.Println("Error on remote call:", err)
	}
	return out
}

func (d *devRendererClient) Render(ctx context.Context, screenSize sdf.V2i, state *DevRendererState, stateLock, cachedRenderLock *sync.RWMutex, partialRender chan<- *image.RGBA) (*image.RGBA, error) {
	argsOut := &RemoteRenderArgsAndResults{
		ScreenSize: screenSize,
		State:      state,
	}
	err := d.cl.Call("DevRendererService.Render", argsOut, &argsOut)
	if err != nil {
		log.Println("Error on remote call:", err)
		return nil, err
	}
	return argsOut.RenderedImg, nil
}

// DevRendererService is the server counter-part to devRendererClient.
// It provides remote access to a devRendererImpl.
// It will block until
type DevRendererService struct {
	impl devRendererImpl
}

// newDevRendererService see DevRendererService
func newDevRendererService(impl devRendererImpl) *rpc.Server {
	server := rpc.NewServer()
	err := server.Register(&DevRendererService{impl: impl})
	if err != nil {
		panic(err) // Shouldn't happen (only on bad implementation)
	}
	return server
}

func (d *DevRendererService) Dimensions(_ int, out *int) error {
	*out = d.impl.Dimensions()
	return nil
}

func (d *DevRendererService) BoundingBox(_ sdf.Box3, out *sdf.Box3) error {
	*out = d.impl.BoundingBox()
	return nil
}

// RemoteRenderArgsAndResults is an internal structure, exported for (de)serialization reasons
type RemoteRenderArgsAndResults struct {
	ScreenSize  sdf.V2i
	State       *DevRendererState
	RenderedImg *image.RGBA
}

func (d *DevRendererService) Render(args RemoteRenderArgsAndResults, out *RemoteRenderArgsAndResults) error {
	// TODO: Cancelling!
	// TODO: Publish partial renders!
	img, err := d.impl.Render(context.Background(), args.ScreenSize, args.State, &sync.RWMutex{}, &sync.RWMutex{}, make(chan *image.RGBA, 512))
	if err != nil {
		return err
	}
	// State attributes that Render might change
	out.State = args.State
	out.RenderedImg = img // The output image
	return nil
}
