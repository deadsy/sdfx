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
		log.Println("[DevRenderer] Error on remote call (RendererService.Dimensions):", err)
	}
	return out
}

func (d *rendererClient) BoundingBox() sdf.Box3 {
	var out sdf.Box3
	err := d.cl.Call("RendererService.BoundingBox", &out, &out)
	if err != nil {
		log.Println("[DevRenderer] Error on remote call (RendererService.BoundingBox):", err)
	}
	return out
}

func (d *rendererClient) ColorModes() int {
	var out int
	err := d.cl.Call("RendererService.ColorModes", &out, &out)
	if err != nil {
		log.Println("[DevRenderer] Error on remote call (RendererService.ColorModes):", err)
	}
	return out
}

func (d *rendererClient) Render(args *renderArgs) error {
	fullRenderSize := args.fullRender.Bounds().Size()
	args.stateLock.RLock() // Clone the state to avoid locking while the rendering is happening
	argsRemote := &RemoteRenderArgs{
		RenderSize: sdf.V2i{fullRenderSize.X, fullRenderSize.Y},
		State:      deepcopy.MustAnything(args.state).(*RendererState),
	}
	args.stateLock.RUnlock()
	var ignoreMe int
	err := d.cl.Call("RendererService.RenderStart", argsRemote, &ignoreMe)
	if err != nil {
		return err
	}
	for {
		var res RemoteRenderResults
		err = d.cl.Call("RendererService.RenderGet", ignoreMe, &res)
		if err != nil {
			return err
		}
		select {
		case <-args.ctx.Done(): // Cancel remote renderer also
			err = d.cl.Call("RendererService.RenderCancel", ignoreMe, &ignoreMe)
			if err != nil {
				log.Println("[DevRenderer] Error on remote call (RendererService.RenderCancel):", err)
			}
			return args.ctx.Err()
		default:
		}
		if res.NewState != nil {
			args.stateLock.Lock() // Clone back the new state to avoid locking while the rendering is happening
			*args.state = *res.NewState
			args.stateLock.Unlock()
		}
		if res.IsPartial {
			if args.partialRenders != nil {
				args.partialRenders <- &*res.RenderedImg // Read-only shallow-copy is enough
			}
		} else { // Final render
			if args.partialRenders != nil {
				close(args.partialRenders)
			}
			args.cachedRenderLock.Lock()
			*args.fullRender = *res.RenderedImg
			args.cachedRenderLock.Unlock()
			break
		}
	}
	return err
}

func (d *rendererClient) Shutdown(timeout time.Duration) error {
	var out int
	return d.cl.Call("RendererService.Shutdown", &timeout, &out)
}

// RendererService is the server counter-part to rendererClient.
// It provides remote access to a devRendererImpl.
type RendererService struct {
	impl                        devRendererImpl
	prevRenderCancel            func()
	renderCtx                   context.Context
	stateLock, cachedRenderLock *sync.RWMutex
	renders                     chan *RemoteRenderResults
	done                        chan os.Signal
}

// newDevRendererService see RendererService
func newDevRendererService(impl devRendererImpl, done chan os.Signal) *rpc.Server {
	server := rpc.NewServer()
	srv := RendererService{
		impl:             impl,
		prevRenderCancel: func() {},
		renderCtx:        context.Background(),
		renders:          make(chan *RemoteRenderResults),
		done:             done,
	}
	close(srv.renders) // Mark the previous render as finished
	err := server.Register(&srv)
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

// RemoteRenderArgs is an internal structure, exported for (de)serialization reasons
type RemoteRenderArgs struct {
	RenderSize sdf.V2i
	State      *RendererState
}

// RemoteRenderResults is an internal structure, exported for (de)serialization reasons
type RemoteRenderResults struct {
	IsPartial   bool
	RenderedImg *image.RGBA
	NewState    *RendererState
}

// RenderStart starts a new render (cancelling the previous one)
func (d *RendererService) RenderStart(args RemoteRenderArgs, _ *int) error {
	d.prevRenderCancel() // Cancel previous render always (no concurrent renderings, although each rendering is parallel by itself)
	var newCtx context.Context
	newCtx, d.prevRenderCancel = context.WithCancel(context.Background())
loop: // Wait for previous renders to be properly completed/cancelled before continuing
	for {
		select {
		case <-newCtx.Done(): // End before started
			return newCtx.Err()
		case _, ok := <-d.renders:
			if !ok {
				break loop
			}
		}
	}
	d.stateLock = &sync.RWMutex{}
	d.cachedRenderLock = &sync.RWMutex{}
	d.cachedRenderLock.Lock()
	d.renderCtx = newCtx
	d.renders = make(chan *RemoteRenderResults)
	d.cachedRenderLock.Unlock()
	partialRenders := make(chan *image.RGBA)
	partialRendersFinish := make(chan struct{})
	go func() { // Start processing partial renders as requested (will silently drop it if not requested)
	loop:
		for partialRender := range partialRenders {
			select {
			case <-d.renderCtx.Done():
				log.Println("[DevRenderer] partialRender cancel")
				break loop
			case d.renders <- &RemoteRenderResults{
				IsPartial:   true,
				RenderedImg: partialRender,
				NewState:    args.State,
			}:
			default:
			}
		}
		close(partialRendersFinish)
	}()
	go func() { // spawn the blocking render in a different goroutine
		fullRender := image.NewRGBA(image.Rect(0, 0, args.RenderSize[0], args.RenderSize[1]))
		err := d.impl.Render(&renderArgs{
			ctx:              d.renderCtx,
			state:            args.State,
			stateLock:        d.stateLock,
			cachedRenderLock: d.cachedRenderLock,
			partialRenders:   partialRenders,
			fullRender:       fullRender,
		})
		if err != nil {
			log.Println("[DevRenderer] RendererService.Render error:", err)
		}
		<-partialRendersFinish // Make sure all partial renders are sent before the full render
		if err == nil {        // Now we can send the full render
			select {
			case d.renders <- &RemoteRenderResults{
				IsPartial:   false,
				RenderedImg: fullRender,
				NewState:    args.State,
			}:
			case <-d.renderCtx.Done():
			}
		}
		close(d.renders)
	}()
	return nil
}

var errNoRenderRunning = errors.New("no render currently running")

// RenderGet gets the next partial or full render available (partial renders might be lost if not called, but not the full render).
// It will return an error if no render is running (or it was cancelled before returning the next result)
func (d *RendererService) RenderGet(_ int, out *RemoteRenderResults) error {
	//d.renderMu.Lock()
	//defer d.renderMu.Unlock()
	select {
	case read, ok := <-d.renders:
		if !ok {
			return errNoRenderRunning
		}
		out.IsPartial = read.IsPartial
		d.cachedRenderLock.RLock() // Need to perform a copy of the image to avoid races with the encoder task
		out.RenderedImg = image.NewRGBA(read.RenderedImg.Rect)
		copy(out.RenderedImg.Pix, read.RenderedImg.Pix)
		d.cachedRenderLock.RUnlock()
		d.stateLock.RLock()
		out.NewState = deepcopy.MustAnything(read.NewState).(*RendererState)
		d.stateLock.RUnlock()
		return nil
	case <-d.renderCtx.Done():
		return errNoRenderRunning // It was cancelled after get was called
	}
}

// RenderCancel cancels the current rendering. It will always succeed with no error.
func (d *RendererService) RenderCancel(_ int, _ *int) error {
	//d.renderMu.Lock()
	//defer d.renderMu.Unlock()
	d.prevRenderCancel() // Cancel previous render
	return nil
}

// Shutdown sends a signal on the configured channel (with a timeout)
func (d *RendererService) Shutdown(t time.Duration, _ *int) error {
	select {
	case d.done <- os.Kill:
		return nil
	case <-time.After(t):
		return errors.New("shutdown timeout")
	}
}
