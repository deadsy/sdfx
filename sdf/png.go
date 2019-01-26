//-----------------------------------------------------------------------------
/*

2D Rendering Code

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
)

//-----------------------------------------------------------------------------

// PNG is a png image object.
type PNG struct {
	name   string
	bb     Box2
	pixels V2i
	m      *Map2
	img    *image.RGBA
}

// NewPNG returns an empty PNG object.
func NewPNG(name string, bb Box2, pixels V2i) (*PNG, error) {
	d := PNG{}
	d.name = name
	d.bb = bb
	d.pixels = pixels
	m, err := NewMap2(bb, pixels, true)
	if err != nil {
		return nil, err
	}
	d.m = m
	d.img = image.NewRGBA(image.Rect(0, 0, pixels[0]-1, pixels[1]-1))
	return &d, nil
}

// RenderSDF2 renders a 2d signed distance field as gray scale.
func (d *PNG) RenderSDF2(s SDF2) {
	// sample the distance field
	var dmax, dmin float64
	distance := make([]float64, d.pixels[0]*d.pixels[1])
	xofs := 0
	for x := 0; x < d.pixels[0]; x++ {
		for y := 0; y < d.pixels[1]; y++ {
			d := s.Evaluate(d.m.ToV2(V2i{x, y}))
			dmax = Max(dmax, d)
			dmin = Min(dmin, d)
			distance[xofs+y] = d
		}
		xofs += d.pixels[1]
	}
	// scale and set the pixel values
	xofs = 0
	for x := 0; x < d.pixels[0]; x++ {
		for y := 0; y < d.pixels[1]; y++ {
			val := 255.0 * ((distance[xofs+y] - dmin) / (dmax - dmin))
			d.img.Set(x, y, color.Gray{uint8(val)})
		}
		xofs += d.pixels[1]
	}
}

// Line adds a line to a png object.
func (d *PNG) Line(p0, p1 V2) {
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

// Lines adds a a set of lines line to a png object.
func (d *PNG) Lines(s V2Set) {
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
	d.Lines([]V2{t[0], t[1], t[2], t[0]})
}

// Save saves a png object to a file.
func (d *PNG) Save() error {
	f, err := os.Create(d.name)
	if err != nil {
		return err
	}
	defer f.Close()
	png.Encode(f, d.img)
	return nil
}

//-----------------------------------------------------------------------------
