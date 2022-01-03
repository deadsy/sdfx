package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"log"
	"net/rpc"
	"sync"
)

// RendererClient implements devRendererImpl by calling a remote implementation (using Go's net/rpc)
type RendererClient struct {
	cl *rpc.Client
}

// newDevRendererClient see RendererClient
func newDevRendererClient(client *rpc.Client) devRendererImpl {
	return &RendererClient{cl: client}
}

func (d *RendererClient) Dimensions() int {
	var out int
	err := d.cl.Call("RendererService.Dimensions", &out, &out)
	if err != nil {
		log.Println("Error on remote call:", err)
	}
	return out
}

func (d *RendererClient) BoundingBox() sdf.Box3 {
	var out sdf.Box3
	err := d.cl.Call("RendererService.BoundingBox", &out, &out)
	if err != nil {
		log.Println("Error on remote call:", err)
	}
	return out
}

func (d *RendererClient) Render(ctx context.Context, state *RendererState, stateLock, cachedRenderLock *sync.RWMutex, partialRender chan<- *image.RGBA, fullRender *image.RGBA) error {
	fullRenderSize := fullRender.Bounds().Size()
	argsOut := &RemoteRenderArgsAndResults{
		RenderSize: sdf.V2i{fullRenderSize.X, fullRenderSize.Y},
		State:      state,
	}
	err := d.cl.Call("RendererService.Render", argsOut, &argsOut)
	if err != nil {
		log.Println("Error on remote call:", err)
		return nil
	}
	*fullRender = *argsOut.RenderedImg
	return nil
}

// RendererService is the server counter-part to RendererClient.
// It provides remote access to a devRendererImpl.
// It will block until
type RendererService struct {
	impl devRendererImpl
}

// newDevRendererService see RendererService
func newDevRendererService(impl devRendererImpl) *rpc.Server {
	server := rpc.NewServer()
	err := server.Register(&RendererService{impl: impl})
	if err != nil {
		panic(err) // Shouldn't happen (only on bad implementation)
	}
	return server
}

func (d *RendererService) Dimensions(_ int, out *int) error {
	*out = d.impl.Dimensions()
	return nil
}

func (d *RendererService) BoundingBox(_ sdf.Box3, out *sdf.Box3) error {
	*out = d.impl.BoundingBox()
	return nil
}

// RemoteRenderArgsAndResults is an internal structure, exported for (de)serialization reasons
type RemoteRenderArgsAndResults struct {
	RenderSize  sdf.V2i
	State       *RendererState
	RenderedImg *image.RGBA
}

func (d *RendererService) Render(args RemoteRenderArgsAndResults, out *RemoteRenderArgsAndResults) error {
	// TODO: Cancelling!
	// TODO: Publish partial renders!
	img := image.NewRGBA(image.Rect(0, 0, args.RenderSize[0], args.RenderSize[1]))
	err := d.impl.Render(context.Background(), args.State, &sync.RWMutex{}, &sync.RWMutex{}, make(chan *image.RGBA, 512), img)
	if err != nil {
		return err
	}
	// State attributes that Render might change
	out.State = args.State
	out.RenderedImg = img // The output image
	return nil
}
