package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"sync"
)

// CONFIGURATION

// Opt3Cam sets the default transform for the camera (pivot center, angles and distance).
// WARNING: Need to run again the main renderer to apply a change of this option.
func Opt3Cam(camCenter sdf.V3, pitch, yaw, dist float64) Option {
	return func(r *Renderer) {
		r.implState.CamCenter = camCenter
		r.implState.CamPitch = pitch
		r.implState.CamYaw = yaw
		r.implState.CamDist = dist
	}
}

// Opt3CamFov sets the default Field Of View for the camera (default 90º, in radians).
func Opt3CamFov(fov float64) Option {
	return func(r *Renderer) {
		if r3, ok := r.impl.(*renderer3); ok {
			r3.camFOV = fov
		}
	}
}

// Opt3RayConfig sets the configuration for the raycast (balancing performance and quality).
// Rendering a pink pixel means that the ray reached maxSteps without hitting the surface or reaching the limit
// (consider increasing maxSteps (reduce performance), increasing epsilon or increasing stepScale (both reduce quality)).
func Opt3RayConfig(scaleAndSigmoid, stepScale, epsilon float64, maxSteps int) Option {
	return func(r *Renderer) {
		if r3, ok := r.impl.(*renderer3); ok {
			r3.rayScaleAndSigmoid = scaleAndSigmoid
			r3.rayStepScale = stepScale
			r3.rayEpsilon = epsilon
			r3.rayMaxSteps = maxSteps
		}
	}
}

// Opt3Colors changes rendering colors.
func Opt3Colors(surface, background, error color.RGBA) Option {
	return func(r *Renderer) {
		if r3, ok := r.impl.(*renderer3); ok {
			r3.surfaceColor = surface
			r3.backgroundColor = background
			r3.errorColor = error
		}
	}
}

// Opt3LightDir sets the light direction for basic lighting simulation (set when Color: true).
// Actually, two lights are simulated (the given one and the opposite one), as part of the surface would be hard to see otherwise
func Opt3LightDir(lightDir sdf.V3) Option {
	return func(r *Renderer) {
		if r3, ok := r.impl.(*renderer3); ok {
			r3.lightDir = lightDir.Normalize()
		}
	}
}

// RENDERER: Z is UP

type renderer3 struct {
	s                                         sdf.SDF3 // The SDF to render
	pixelsRand                                []int    // Cached set of pixels in random order to avoid shuffling (reset on recompilation and resolution changes)
	camFOV                                    float64  // The Field Of View (X axis) for the camera
	surfaceColor, backgroundColor, errorColor color.RGBA
	lightDir                                  sdf.V3 // The light's direction for ColorMode: true (simple simulation based on normals)
	// Raycast configuration
	rayScaleAndSigmoid, rayStepScale, rayEpsilon float64
	rayMaxSteps                                  int
}

func newDevRenderer3(s sdf.SDF3) devRendererImpl {
	r := &renderer3{
		s:                  s,
		camFOV:             math.Pi / 2, // 90º FOV-X
		surfaceColor:       color.RGBA{R: 255 - 20, G: 255 - 40, B: 255 - 80, A: 255},
		backgroundColor:    color.RGBA{B: 50, A: 255},
		errorColor:         color.RGBA{R: 255, B: 255, A: 255},
		lightDir:           sdf.V3{X: 1, Y: 1, Z: -1}.Normalize(), // Same as default camera
		rayScaleAndSigmoid: 0,
		rayStepScale:       1,
		rayEpsilon:         0.1,
		rayMaxSteps:        100,
	}
	return r
}

func (r *renderer3) Dimensions() int {
	return 3
}

func (r *renderer3) BoundingBox() sdf.Box3 {
	return r.s.BoundingBox()
}

func (r *renderer3) ColorModes() int {
	// 0: Constant color with basic shading (2 lights and no shadows)
	// 1: Normal XYZ as RGB
	return 2
}

func (r *renderer3) Render(ctx context.Context, state *RendererState, stateLock, cachedRenderLock *sync.RWMutex, partialRender chan<- *image.RGBA, fullRender *image.RGBA) error {
	// Set all pixels to transparent initially (for partial renderings to work)
	for i := 3; i < len(fullRender.Pix); i += 4 {
		fullRender.Pix[i] = 255
	}

	// TODO: Fix blocked Render after reload

	// Update random pixels if needed
	bounds := fullRender.Bounds()
	boundsSize := sdf.V2i{bounds.Size().X, bounds.Size().Y}
	pixelCount := boundsSize[0] * boundsSize[1]
	if pixelCount != len(r.pixelsRand) {
		r.pixelsRand = make([]int, pixelCount)
		for i := 0; i < pixelCount; i++ {
			r.pixelsRand[i] = i
		}
		rand.Shuffle(len(r.pixelsRand), func(i, j int) {
			r.pixelsRand[i], r.pixelsRand[j] = r.pixelsRand[j], r.pixelsRand[i]
		})
	}

	// Spawn the workers that will render 1 pixel at a time
	jobs := make(chan *pixelRender)
	jobResults := make(chan *pixelRender)
	workerWg := &sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		workerWg.Add(1)
		go func() {
			for job := range jobs {
				job.rendered = r.samplePixel(job)
				jobResults <- job
			}
			workerWg.Done()
		}()
	}
	go func() {
		workerWg.Wait()
		close(jobResults)
	}()

	// Compute camera position and main direction
	aspectRatio := float64(boundsSize[0]) / float64(boundsSize[1])
	camViewMatrix := state.Cam3MatrixNoTranslation()
	camPos := state.CamCenter.Add(camViewMatrix.MulPosition(sdf.V3{Y: -state.CamDist}))
	camDir := state.CamCenter.Sub(camPos).Normalize()
	camFovX := r.camFOV
	camFovY := 2 * math.Atan(math.Tan(camFovX/2)*aspectRatio)
	//log.Println("cam:", camPos, "->", camDir)

	// Spawn the work generator
	go func() { // TODO: Races by reusing variables (like i in for loop)?
		// Sample each pixel on the image separately (and in random order to see the image faster)
		for _, randPixelIndex := range r.pixelsRand {
			// Sample a random pixel in the image
			sampledPixel := sdf.V2i{randPixelIndex % boundsSize[0], randPixelIndex / boundsSize[0]}
			// Queue the job for parallel processing
			jobs <- &pixelRender{
				pixel:         sampledPixel,
				bounds:        boundsSize,
				camPos:        camPos,
				camDir:        camDir,
				camViewMatrix: camViewMatrix,
				camFov:        sdf.V2{X: camFovX, Y: camFovY},
				color:         state.ColorMode,
				rendered:      color.RGBA{},
			}
		}
		close(jobs) // Close the jobs channel to mark the end
	}()

	// Listen for all job results and update the image, freeing locks and sending a partial image update every batch of pixels
	const pixelBatch = 100
	pixelNum := 0
	cachedRenderLock.Lock()
	var err error
pixelLoop:
	for renderedPixel := range jobResults {
		fullRender.SetRGBA(renderedPixel.pixel[0], renderedPixel.pixel[1], renderedPixel.rendered)
		pixelNum++
		if pixelNum%pixelBatch == 0 {
			cachedRenderLock.Unlock()
			runtime.Gosched() // Breathe (let renderer do something, best-effort)
			// Check if this render is cancelled (could also check every pixel...)
			select {
			case <-ctx.Done():
				err = ctx.Err()
				break pixelLoop
			default:
			}
			// Send the partial render update
			//log.Println("Sending partial render with", pixelNum, "pixels")
			//tempFile, _ := ioutil.TempFile("", "fullRender-"+strconv.Itoa(pixelNum)+"-*.png")
			//_ = png.Encode(tempFile, fullRender)
			//log.Println("Written PNG to", tempFile.Name())
			if partialRender != nil {
				// TODO: Use a shader to fill transparent pixel with nearest neighbors to make it look better while rendering
				partialRender <- fullRender
			}
			//time.Sleep(time.Second)
			cachedRenderLock.Lock()
		}
	}
	if err == nil {
		cachedRenderLock.Unlock()
	}
	close(partialRender)
	// TODO: Draw bounding boxes over the image
	return err
}

type pixelRender struct {
	// CAMERA RELATED
	pixel, bounds  sdf.V2i // Pixel and bounds for pixel
	camPos, camDir sdf.V3  // Camera parameters
	camViewMatrix  sdf.M44 // The world to camera matrix
	camFov         sdf.V2  // Camera's field of view
	// MISC
	color int
	// OUTPUT
	rendered color.RGBA
}

func (r *renderer3) samplePixel(job *pixelRender) color.RGBA {
	// Generate the ray for this pixel using the given camera parameters
	rayFrom := job.camPos
	// Get pixel inside of ([-1, 1], [-1, 1])
	rayDirXZBase := job.pixel.ToV2().Div(job.bounds.ToV2()).MulScalar(2).SubScalar(1)
	rayDirXZBase.X *= float64(job.bounds[0]) / float64(job.bounds[1]) // Apply aspect ratio again
	// Convert to the projection over a displacement of 1
	rayDirXZBase = rayDirXZBase.Mul(sdf.V2{X: math.Tan(job.camFov.DivScalar(2).X), Y: math.Tan(job.camFov.DivScalar(2).Y)})
	rayDir := sdf.V3{X: rayDirXZBase.X, Y: 1, Z: rayDirXZBase.Y}
	// Apply the camera matrix to the default ray
	rayDir = job.camViewMatrix.MulPosition(rayDir).Normalize()
	// TODO: Orthogonal camera

	// Query the surface with the given ray
	maxRay := 10000. // TODO: Compute the actual value
	hit, t, steps := sdf.Raycast3(r.s, rayFrom, rayDir, r.rayScaleAndSigmoid, r.rayStepScale, r.rayEpsilon, maxRay, r.rayMaxSteps)
	//if job.pixel[0] == job.bounds[0]/2 {
	//	log.Println("ray dir:", rayDir, "T:", t, "HIT:", hit, "STEPS:", steps)
	//}

	// Convert the possible hit to a color
	if t >= 0 { // Hit the surface
		normal := sdf.Normal3(r.s, hit, 1e-3)
		if job.color == 0 { // Basic lighting + constant color
			lightIntensity := math.Abs(normal.Dot(r.lightDir)) // Actually also simulating the opposite light
			// If this was a performant ray-tracer, we could bounce the light
			return color.RGBA{
				R: uint8(float64(r.surfaceColor.R) * lightIntensity),
				G: uint8(float64(r.surfaceColor.G) * lightIntensity),
				B: uint8(float64(r.surfaceColor.B) * lightIntensity),
				A: r.surfaceColor.A,
			}
		} else { // Color == normal
			return color.RGBA{R: uint8(normal.X * 255), G: uint8(normal.Y * 255), B: uint8(normal.Z * 255), A: 255}
		}
	} else {
		if steps == r.rayMaxSteps {
			// Reached the maximum amount of steps (should change parameters)
			return r.errorColor
		}
		// The void
		return r.backgroundColor
	}
}
