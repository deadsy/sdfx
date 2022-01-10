//-----------------------------------------------------------------------------
/*

2D Rendering Code

*/
//-----------------------------------------------------------------------------

package render

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/deadsy/sdfx/sdf"
	"github.com/llgcode/draw2d/draw2dimg"
)

//-----------------------------------------------------------------------------

// PNG is a png image object.
type PNG struct {
	name   string
	bb     sdf.Box2
	pixels sdf.V2i
	m      *sdf.Map2
	img    *image.RGBA
}

// NewPNG returns an empty PNG object.
func NewPNG(name string, bb sdf.Box2, pixels sdf.V2i) (*PNG, error) {
	d := PNG{}
	d.name = name
	d.bb = bb
	d.pixels = pixels
	m, err := sdf.NewMap2(bb, pixels, true)
	if err != nil {
		return nil, err
	}
	d.m = m
	d.img = image.NewRGBA(image.Rect(0, 0, pixels[0]-1, pixels[1]-1))
	return &d, nil
}

// RenderSDF2 renders a 2d signed distance field as gray scale.
func (d *PNG) RenderSDF2(s sdf.SDF2) {
	d.RenderSDF2MinMax(s, 0, 0)
}

// RenderSDF2MinMax renders a 2d signed distance field as gray scale (with defined min/max levels).
func (d *PNG) RenderSDF2MinMax(s sdf.SDF2, dmin, dmax float64) {
	// sample the distance field
	minMaxSet := dmin != 0 && dmax != 0
	if !minMaxSet {
		//distance := make([]float64, d.pixels[0]*d.pixels[1]) // Less allocations: faster (70ms -> 60ms), use cache in SDF if needed
		for x := 0; x < d.pixels[0]; x++ {
			for y := 0; y < d.pixels[1]; y++ {
				d := s.Evaluate(d.m.ToV2(sdf.V2i{x, y}))
				dmax = math.Max(dmax, d)
				dmin = math.Min(dmin, d)
			}
		}
	}
	// scale and set the pixel values
	for x := 0; x < d.pixels[0]; x++ {
		for y := 0; y < d.pixels[1]; y++ {
			dist := s.Evaluate(d.m.ToV2(sdf.V2i{x, y}))
			d.img.Set(x, y, color.Gray{Y: uint8(255 * imageColor2(dist, dmin, dmax))})
		}
	}
}

// imageColor2 returns the grayscale color for the returned SDF2.Evaluate value, given the reference minimum and maximum
// SDF2.Evaluate values. The returned value is in the range [0, 1].
func imageColor2(dist, dmin, dmax float64) float64 {
	// Clamp due to possibly forced min and max
	var val float64
	// NOTE: This condition forces the surface to be close to 0.5 gray value, otherwise dmax >>> dmin or viceversa
	// could cause the surface to be displaced visually
	if dist >= 0 {
		val = math.Max(0.5, math.Min(1, 0.5+0.5*((dist)/(dmax))))
	} else { // Force lower scale for inside surface
		val = math.Max(0, math.Min(0.5, 0.5*((dist-dmin)/(-dmin))))
	}
	return val
}

// Line adds a line to a png object.
func (d *PNG) Line(p0, p1 sdf.V2) {
	gc := draw2dimg.NewGraphicContext(d.img)
	gc.SetFillColor(color.RGBA{0xff, 0, 0, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0, 0, 0xff})
	gc.SetLineWidth(1)

	p := d.m.ToV2i(p0)
	gc.MoveTo(float64(p[0]), float64(p[1]))
	p = d.m.ToV2i(p1)
	gc.LineTo(float64(p[0]), float64(p[1]))
	gc.Stroke()
}

// Lines adds a set of lines line to a png object.
func (d *PNG) Lines(s sdf.V2Set) {
	gc := draw2dimg.NewGraphicContext(d.img)
	gc.SetFillColor(color.RGBA{0xff, 0, 0, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0, 0, 0xff})
	gc.SetLineWidth(1)

	p := d.m.ToV2i(s[0])
	gc.MoveTo(float64(p[0]), float64(p[1]))
	for i := 1; i < len(s); i++ {
		p := d.m.ToV2i(s[i])
		gc.LineTo(float64(p[0]), float64(p[1]))
	}
	gc.Stroke()
}

// Triangle adds a triangle to a png object.
func (d *PNG) Triangle(t Triangle2) {
	d.Lines([]sdf.V2{t[0], t[1], t[2], t[0]})
}

// Save saves a png object to a file.
func (d *PNG) Save() error {
	f, err := os.Create(d.name)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, d.img)
}

// Image returns the rendered image instead of writing it to a file
func (d *PNG) Image() *image.RGBA {
	return d.img
}

//-----------------------------------------------------------------------------
