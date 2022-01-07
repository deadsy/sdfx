package dev

import (
	"github.com/deadsy/sdfx/sdf"
	"image/color"
	"math/rand"
	"runtime"
	"sync"
)

type jobInternal struct {
	pixel   sdf.V2i
	pixel01 sdf.V2
	data    interface{}
}

type jobResult struct {
	pixel sdf.V2i
	color color.RGBA
}

func implCommonRender(genJob func(pixel sdf.V2i, pixel01 sdf.V2) interface{},
	processJob func(pixel sdf.V2i, pixel01 sdf.V2, job interface{}) *jobResult,
	args *renderArgs, pixelsRand *[]int) error {

	// Set all pixels to transparent initially (for partial renderings to work)
	args.cachedRenderLock.Lock()
	for i := 3; i < len(args.fullRender.Pix); i += 4 {
		args.fullRender.Pix[i] = 255
	}
	args.cachedRenderLock.Unlock()

	// Update random pixels if needed
	bounds := args.fullRender.Bounds()
	boundsSize := sdf.V2i{bounds.Size().X, bounds.Size().Y}
	pixelCount := boundsSize[0] * boundsSize[1]
	if pixelCount != len(*pixelsRand) {
		// Random seed shouldn't matter, just make pixel coloring seem random for partial renders
		*pixelsRand = rand.Perm(pixelCount)
	}

	// Spawn the workers that will render 1 pixel at a time
	jobs := make(chan *jobInternal)
	jobResults := make(chan *jobResult)
	workerWg := &sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
		loop:
			for {
				job, ok := <-jobs
				if !ok { // Cancelled or finished render (stopped generating jobs)
					break loop
				}
				jobResults <- processJob(job.pixel, job.pixel01, job.data)
			}
		}()
	}
	go func() { // Make sure job results are closed after all jobs are processed
		workerWg.Wait()
		close(jobResults)
	}()

	// Spawn the work generator
	go func() {
	loop: // Sample each pixel on the image separately (and in random order to see the image faster)
		for _, randPixelIndex := range *pixelsRand {
			// Sample a random pixel in the image
			sampledPixel := sdf.V2i{randPixelIndex % boundsSize[0], randPixelIndex / boundsSize[0]}
			sampledPixel01 := sampledPixel.ToV2().Div(boundsSize.ToV2())
			// Queue the job for parallel processing
			select {
			case <-args.ctx.Done():
				break loop
			case jobs <- &jobInternal{
				pixel:   sampledPixel,
				pixel01: sampledPixel01,
				data:    genJob(sampledPixel, sampledPixel01),
			}:
			}
		}
		close(jobs) // Close the jobs channel to mark the end
	}()

	// Listen for all job results and update the image, freeing locks and sending a partial image update every batch of pixels
	const pixelBatch = 1000 // Configurable? Shouldn't matter much as you can already configure time between partial renders.
	pixelNum := 0
	args.cachedRenderLock.Lock()
	var err error
pixelLoop:
	for renderedPixel := range jobResults {
		args.fullRender.SetRGBA(renderedPixel.pixel[0], renderedPixel.pixel[1], renderedPixel.color)
		pixelNum++
		if pixelNum%pixelBatch == 0 {
			args.cachedRenderLock.Unlock()
			runtime.Gosched() // Breathe (let renderer do something, best-effort)
			select {          // Check if this render is cancelled (could also check every pixel...)
			case <-args.ctx.Done():
				err = args.ctx.Err()
				break pixelLoop
			default:
			}
			if args.partialRender != nil { // Send the partial render update
				// TODO: Use a shader to fill transparent pixel with nearest neighbors to make it look better while rendering (losing previous background render)
				args.partialRender <- args.fullRender
			}
			args.cachedRenderLock.Lock()
		}
	}
	if err == nil {
		args.cachedRenderLock.Unlock()
	}
	if args.partialRender != nil {
		close(args.partialRender)
	}
	return err
}
