//-----------------------------------------------------------------------------
/*

Text Operations

Convert a string and font specification into an SDF2

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

type align int

const (
	lAlign align = iota // left hand side x = 0
	rAlign              // right hand side x = 0
	cAlign              // center x = 0
)

// Text stores a UTF8 string and it's rendering parameters.
type Text struct {
	s      string
	halign align
}

//-----------------------------------------------------------------------------

// pToV2 converts a truetype point to a V2
func pToV2(p truetype.Point) V2 {
	return V2{float64(p.X), float64(p.Y)}
}

//-----------------------------------------------------------------------------

// glyphCurve returns the SDF2 for the n-th curve of the glyph
func glyphCurve(g *truetype.GlyphBuf, n int) (SDF2, bool) {
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
	offPrev := false
	vPrev := pToV2(g.Points[end])

	for i := start; i <= end; i++ {
		p := g.Points[i]
		v := pToV2(p)
		// is the point off/on the curve?
		off := p.Flags&1 == 0
		// do we have an implicit on point?
		if off && offPrev {
			// implicit on point at the midpoint of the 2 off points
			b.AddV2(v.Add(vPrev).MulScalar(0.5))
		}
		// add the point
		x := b.AddV2(v)
		if off {
			x.Mid()
		}
		// accumulate the cw/ccw direction
		sum += (v.X - vPrev.X) * (v.Y + vPrev.Y)
		// next point...
		vPrev = v
		offPrev = off
	}
	b.Close()

	return Polygon2D(b.Polygon().Vertices()), sum > 0
}

// glyphConvert returns the SDF2 for a glyph
func glyphConvert(g *truetype.GlyphBuf) SDF2 {
	var s0 SDF2
	for n := 0; n < len(g.Ends); n++ {
		s1, cw := glyphCurve(g, n)
		if cw {
			s0 = Union2D(s0, s1)
		} else {
			s0 = Difference2D(s0, s1)
		}
	}
	return s0
}

//-----------------------------------------------------------------------------

// lineSDF2 returns an SDF2 slice for a line of text
func lineSDF2(f *truetype.Font, l string) ([]SDF2, float64, error) {
	iPrev := truetype.Index(0)
	scale := fixed.Int26_6(f.FUnitsPerEm())
	xOfs := 0.0

	var ss []SDF2

	for _, r := range l {
		i := f.Index(r)

		// get the glyph metrics
		hm := f.HMetric(scale, i)

		// apply kerning
		k := f.Kern(scale, iPrev, i)
		xOfs += float64(k)
		iPrev = i

		// load the glyph
		g := &truetype.GlyphBuf{}
		err := g.Load(f, scale, i, font.HintingNone)
		if err != nil {
			return nil, 0, err
		}

		s := glyphConvert(g)
		if s != nil {
			s = Transform2D(s, Translate2d(V2{xOfs, 0}))
			ss = append(ss, s)
		}

		xOfs += float64(hm.AdvanceWidth)
	}

	return ss, xOfs, nil
}

//-----------------------------------------------------------------------------
// public api

// NewText returns a text object (text and alignment).
func NewText(s string) *Text {
	return &Text{
		s:      s,
		halign: cAlign,
	}
}

// LoadFont loads a truetype (*.ttf) font file.
func LoadFont(fname string) (*truetype.Font, error) {
	// read the font file
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(b)
}

// TextSDF2 returns a sized SDF2 for a text object.
func TextSDF2(f *truetype.Font, t *Text, h float64) (SDF2, error) {
	scale := fixed.Int26_6(f.FUnitsPerEm())
	lines := strings.Split(t.s, "\n")
	yOfs := 0.0
	vm := f.VMetric(scale, f.Index('\n'))
	ah := float64(vm.AdvanceHeight)

	var ss []SDF2

	for i := range lines {
		ssLine, hlen, err := lineSDF2(f, lines[i])
		if err != nil {
			return nil, err
		}
		xOfs := 0.0
		if t.halign == rAlign {
			xOfs = -hlen
		} else if t.halign == cAlign {
			xOfs = -hlen / 2.0
		}
		for i := range ssLine {
			ssLine[i] = Transform2D(ssLine[i], Translate2d(V2{xOfs, yOfs}))
		}
		ss = append(ss, ssLine...)
		yOfs -= ah
	}

	return CenterAndScale2D(Union2D(ss...), h/ah), nil
}

//-----------------------------------------------------------------------------
