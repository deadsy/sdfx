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

func glyph_convert(g *truetype.GlyphBuf) []V2 {

	k := 1.0

	b := NewBezier()

	e := 0 // index into the endpoint table

	prev_off := false

	for i, p := range g.Points {
		// is the point off/on the curve?
		off := p.Flags&0x01 == 0

		if off && prev_off {
			// implicit on point at the midpoint of the 2 off points
			pp := g.Points[i-1]
			v := V2{(float64(pp.X) + float64(p.X)) * 0.5, (float64(pp.Y) + float64(p.Y)) * 0.5}
			v.MulScalar(k)
			b.AddV2(v)
		}

		// add the point
		v := V2{float64(p.X), float64(p.Y)}
		v.MulScalar(k)
		x := b.AddV2(v)
		if off {
			x.Mid()
		}

		prev_off = off

		if i+1 == int(g.Ends[e]) {
			// TODO
			break
			e++
		}
	}

	b.Close()

	return b.Polygon().Vertices()
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

	c := 'm'
	i := f.Index(c)
	hm := f.HMetric(fupe, i)
	g := &truetype.GlyphBuf{}
	err = g.Load(f, fupe, i, font.HintingNone)
	if err != nil {
		return err
	}
	fmt.Printf("'%c' glyph\n", c)
	fmt.Printf("AdvanceWidth:%d LeftSideBearing:%d\n", hm.AdvanceWidth, hm.LeftSideBearing)
	printGlyph(g)

	s2d := Polygon2D(glyph_convert(g))
	RenderDXF(s2d, 200, "shape.dxf")

	a := truetype.NewFace(f, &truetype.Options{
		Size: 12,
		DPI:  72,
	})
	fmt.Printf("%#v\n", a.Metrics())

	return nil
}

//-----------------------------------------------------------------------------
