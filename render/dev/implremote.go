package dev

import (
	"context"
	"errors"
	"github.com/barkimedes/go-deepcopy"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"log"
	"net/rpc"
	"os"
	"sync"
	"time"
)

// rendererClient implements devRendererImpl by calling a remote implementation (using Go's net/rpc)
type rendererClient struct {
	cl *rpc.Client
}

// newDevRendererClient see rendererClient
func newDevRendererClient(client *rpc.Client) devRendererImpl {
	return &rendererClient{cl: client}
}

func (d *rendererClient) Dimensions() int {
	var out int
	err := d.cl.Call("RendererService.Dimensions", &out, &out)
	if err != nil {
		log.Println("Error on remote call:", err)
	}
	return out
}

func (d *rendererClient) BoundingBox() sdf.Box3 {
	var out sdf.Box3
	err := d.cl.Call("RendererService.BoundingBox", &out, &out)
	if err != nil {
		log.Println("Error on remote call:", err)
	}
	return out
}

func (d *rendererClient) ColorModes() int {
	var out int
	err := d.cl.Call("RendererService.ColorModes", &out, &out)
	if err != nil {
		log.Println("Error on remote call:", err)
	}
	return out
}

func (d *rendererClient) Render(args *renderArgs) error {
	fullRenderSize := args.fullRender.Bounds().Size()
	args.stateLock.RLock() // Clone the state to avoid locking while the rendering is happening
	argsOut := &RemoteRenderArgsAndResults{
		RenderSize: sdf.V2i{fullRenderSize.X, fullRenderSize.Y},
		State:      deepcopy.MustAnything(args.state).(*RendererState),
	}
	args.stateLock.RUnlock()

	call := d.cl.Go("RendererService.Render", argsOut, &argsOut, nil)
	select {
	case <-args.ctx.Done(): // Cancelled (call still running on service unless we call render again, which will cancel it)
		return args.ctx.Err()
	case call := <-call.Done:
		if call.Error != nil {
			log.Println("Error on remote call:", call.Error)
			return nil
		}
		args.stateLock.Lock() // Clone back the new state to avoid locking while the rendering is happening
		*args.state = *argsOut.State
		args.stateLock.Unlock()
		args.cachedRenderLock.Lock()
		*args.fullRender = *(*call.Reply.(**RemoteRenderArgsAndResults)).RenderedImg
		args.cachedRenderLock.Unlock()
		return nil
	}
}

func (d *rendererClient) Shutdown(timeout time.Duration) error {
	var out int
	return d.cl.Call("RendererService.Shutdown", &timeout, &out)
}

// RendererService is the server counter-part to rendererClient.
// It provides remote access to a devRendererImpl.
type RendererService struct {
	impl                 devRendererImpl
	prevRenderCancel     func()
	prevRenderCancelLock *sync.Mutex
	done                 chan os.Signal
}

// newDevRendererService see RendererService
func newDevRendererService(impl devRendererImpl, done chan os.Signal) *rpc.Server {
	server := rpc.NewServer()
	err := server.Register(&RendererService{impl: impl, prevRenderCancel: func() {}, prevRenderCancelLock: &sync.Mutex{}, done: done})
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

func (d *RendererService) ColorModes(_ int, out *int) error {
	*out = d.impl.ColorModes()
	return nil
}

// RemoteRenderArgsAndResults is an internal structure, exported for (de)serialization reasons
type RemoteRenderArgsAndResults struct {
	RenderSize  sdf.V2i
	State       *RendererState
	RenderedImg *image.RGBA
}

func (d *RendererService) Render(args RemoteRenderArgsAndResults, out *RemoteRenderArgsAndResults) error {
	// TODO: Publish partial renders!
	img := image.NewRGBA(image.Rect(0, 0, args.RenderSize[0], args.RenderSize[1]))
	var ctx context.Context
	d.prevRenderCancelLock.Lock()
	d.prevRenderCancel() // Cancel previous render
	ctx, d.prevRenderCancel = context.WithCancel(context.Background())
	d.prevRenderCancelLock.Unlock()
	err := d.impl.Render(&renderArgs{
		ctx:              ctx,
		state:            args.State,
		stateLock:        &sync.RWMutex{},
		cachedRenderLock: &sync.RWMutex{},
		partialRender:    nil,
		fullRender:       img,
	})
	if err != nil {
		return err
	}
	// State attributes that Render might change
	out.State = args.State
	out.RenderedImg = img // The output image
	return nil
}

func (d *RendererService) Shutdown(t time.Duration, out *int) error {
	select {
	case d.done <- os.Kill:
		return nil
	case <-time.After(t):
		return errors.New("shutdown timeout")
	}
}
