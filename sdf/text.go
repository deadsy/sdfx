//-----------------------------------------------------------------------------
/*

Text Operations

Convert a string and a font specification into an SDF2

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//-----------------------------------------------------------------------------

const POINT_PER_INCH = 72.0

//-----------------------------------------------------------------------------

func printBounds(b fixed.Rectangle26_6) {
	fmt.Printf("Min.X:%d Min.Y:%d Max.X:%d Max.Y:%d\n", b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
}

func printGlyph(g *truetype.GlyphBuf) {
	printBounds(g.Bounds)
	fmt.Print("Points:\n---\n")
	e := 0
	for i, p := range g.Points {
		fmt.Printf("%4d, %4d", p.X, p.Y)
		if p.Flags&0x01 != 0 {
			fmt.Print("  on\n")
		} else {
			fmt.Print("  off\n")
		}
		if i+1 == int(g.Ends[e]) {
			fmt.Print("---\n")
			e++
		}
	}

	fmt.Printf("points: %v\n", g.Points)
	fmt.Printf("ends: %v\n", g.Ends)

}

//-----------------------------------------------------------------------------

// return the SDF2 for the n-th curve of the glyph
func glyph_curve(g *truetype.GlyphBuf, n int) (SDF2, bool) {
	// get the start and end point
	start := 0
	if n != 0 {
		start = g.Ends[n-1]
	}
	end := g.Ends[n] - 1
	// build a bezier curve from the points
	b := NewBezier()
	prev_off := false
	for i := start; i <= end; i++ {
		p := g.Points[i]
		// is the point off/on the curve?
		off := p.Flags&1 == 0
		// do we have an implicit on point?
		if off && prev_off {
			// implicit on point at the midpoint of the 2 off points
			pp := g.Points[i-1]
			v := V2{(float64(pp.X) + float64(p.X)) * 0.5, (float64(pp.Y) + float64(p.Y)) * 0.5}
			b.AddV2(v)
		}
		// add the point
		v := V2{float64(p.X), float64(p.Y)}
		x := b.AddV2(v)
		if off {
			x.Mid()
		}
		prev_off = off
	}
	b.Close()
	// determine the cw/ccw direction
	// TODO - write this better
	cw := false
	sum := 0.0
	for i := start; i <= end; i++ {
		prev := i - 1
		if i == start {
			prev = end
		}
		next := i + 1
		if i == end {
			next = start
		}
		p := g.Points[i]
		pn := g.Points[next]
		pp := g.Points[prev]
		v := V2{float64(p.X), float64(p.Y)}
		vn := V2{float64(pn.X), float64(pn.Y)}
		vp := V2{float64(pp.X), float64(pp.Y)}
		v0 := v.Sub(vp)
		v1 := vn.Sub(v)
		sum += v0.Cross(v1)
	}
	if sum < 0 {
		cw = true
	}
	return Polygon2D(b.Polygon().Vertices()), cw
}

// return the SDF2 for a glyph
func glyph_convert(g *truetype.GlyphBuf) SDF2 {
	var s0 SDF2
	for n := 0; n < len(g.Ends); n++ {
		s1, cw := glyph_curve(g, n)
		if cw {
			s0 = Union2D(s0, s1)
		} else {
			s0 = Difference2D(s0, s1)
		}
	}
	return s0
}

//-----------------------------------------------------------------------------

func Test_Text() error {

	// get the font data
	fontfile := "/usr/share/fonts/truetype/msttcorefonts/Arial_Black.ttf"
	b, err := ioutil.ReadFile(fontfile)
	if err != nil {
		return err
	}

	f, err := truetype.Parse(b)
	if err != nil {
		return err
	}

	fupe := fixed.Int26_6(f.FUnitsPerEm())
	printBounds(f.Bounds(fupe))
	fmt.Printf("FUnitsPerEm:%d\n\n", fupe)

	c := 'Q'
	i := f.Index(c)
	g := &truetype.GlyphBuf{}
	err = g.Load(f, fupe, i, font.HintingNone)
	if err != nil {
		return err
	}

	//hm := f.HMetric(fupe, i)
	//fmt.Printf("'%c' glyph\n", c)
	//fmt.Printf("AdvanceWidth:%d LeftSideBearing:%d\n", hm.AdvanceWidth, hm.LeftSideBearing)
	//printGlyph(g)

	s2d := glyph_convert(g)
	RenderDXF(s2d, 200, "shape.dxf")

	s3d := ExtrudeRounded3D(s2d, 200, 20)
	RenderSTL(s3d, 300, "shape.stl")

	//a := truetype.NewFace(f, &truetype.Options{
	//	Size: 12,
	//	DPI:  72,
	//})
	//fmt.Printf("%#v\n", a.Metrics())

	return nil
}

//-----------------------------------------------------------------------------
