//-----------------------------------------------------------------------------
/*

Text Operations

Convert a string and a font specification into an SDF2

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"io/ioutil"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//-----------------------------------------------------------------------------

const POINT_PER_INCH = 72.0

//-----------------------------------------------------------------------------

type align int

const (
	L_ALIGN align = iota // left hand side x = 0
	R_ALIGN              // right hand side x = 0
	C_ALIGN              // center x = 0
)

type Text struct {
	s      string
	halign align
}

func NewText(s string) *Text {
	return &Text{
		s:      s,
		halign: C_ALIGN,
	}
}

//-----------------------------------------------------------------------------

// convert a truetype point to a V2
func p_to_V2(p truetype.Point) V2 {
	return V2{float64(p.X), float64(p.Y)}
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
	// work out the cw/ccw direction
	b := NewBezier()
	sum := 0.0
	off_prev := false
	v_prev := p_to_V2(g.Points[end])

	for i := start; i <= end; i++ {
		p := g.Points[i]
		v := p_to_V2(p)
		// is the point off/on the curve?
		off := p.Flags&1 == 0
		// do we have an implicit on point?
		if off && off_prev {
			// implicit on point at the midpoint of the 2 off points
			b.AddV2(v.Add(v_prev).MulScalar(0.5))
		}
		// add the point
		x := b.AddV2(v)
		if off {
			x.Mid()
		}
		// accumulate the cw/ccw direction
		sum += (v.X - v_prev.X) * (v.Y + v_prev.Y)
		// next point...
		v_prev = v
		off_prev = off
	}
	b.Close()

	return Polygon2D(b.Polygon().Vertices()), sum > 0
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

// return the SDF2 for a line of text
func lineSDF2(f *truetype.Font, l string) (SDF2, float64, error) {

	var s0 SDF2
	i_prev := truetype.Index(0)
	scale := fixed.Int26_6(f.FUnitsPerEm())
	x_ofs := 0.0

	for _, r := range l {
		i := f.Index(r)

		// get the glyph metrics
		hm := f.HMetric(scale, i)

		// apply kerning
		k := f.Kern(scale, i_prev, i)
		x_ofs += float64(k)
		i_prev = i

		// load the glyph
		g := &truetype.GlyphBuf{}
		err := g.Load(f, scale, i, font.HintingNone)
		if err != nil {
			return nil, 0, err
		}

		s1 := glyph_convert(g)
		if s1 != nil {
			s1 = Transform2D(s1, Translate2d(V2{x_ofs, 0}))
			s0 = Union2D(s0, s1)
		}

		x_ofs += float64(hm.AdvanceWidth)
	}

	return s0, x_ofs, nil
}

//-----------------------------------------------------------------------------
// public api

// load a truetype (*.ttf) font file
func LoadFont(fname string) (*truetype.Font, error) {
	// read the font file
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(b)
}

// return an SDF2 for the text string
func TextSDF2(f *truetype.Font, t *Text) (SDF2, error) {

	scale := fixed.Int26_6(f.FUnitsPerEm())
	lines := strings.Split(t.s, "\n")
	var s0 SDF2
	y_ofs := 0.0
	vm := f.VMetric(scale, f.Index('\n'))
	ah := float64(vm.AdvanceHeight)

	for i := range lines {
		s1, hlen, err := lineSDF2(f, lines[i])
		if err != nil {
			return nil, err
		}
		x_ofs := 0.0
		if t.halign == R_ALIGN {
			x_ofs = -hlen
		} else if t.halign == C_ALIGN {
			x_ofs = -hlen / 2.0
		}
		s1 = Transform2D(s1, Translate2d(V2{x_ofs, y_ofs}))
		s0 = Union2D(s0, s1)
		y_ofs -= ah
	}

	s0 = Center2D(s0)
	return s0, nil
}

//-----------------------------------------------------------------------------
