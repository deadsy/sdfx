//-----------------------------------------------------------------------------
/*

2D line segments

Used for building 2D polygons SDFs.

*/
//-----------------------------------------------------------------------------

package sdf

import "fmt"

//-----------------------------------------------------------------------------
// 2D Line Segment

type Line2 struct {
	segment bool    // is this a line segment
	length  float64 // segment length
	a       V2      // line start point
	b       V2      // line end point point (if segment)
	v       V2      // normalized line vector
}

// Create a new 2d line given a point and vector
func NewLine2_PV(p, v V2) Line2 {
	l := Line2{}
	l.segment = false
	l.length = 0.0
	l.a = p
	l.v = v.Normalize()
	return l
}

// Create a new 2d line given a point and a point
func NewLine2_PP(a, b V2) Line2 {
	l := Line2{}
	v := b.Sub(a)
	l.segment = true
	l.length = v.Length()
	l.a = a
	l.b = b
	l.v = v.Normalize()
	return l
}

// Return the position given the t value
func (l Line2) Position(t float64) V2 {
	return l.a.Add(l.v.MulScalar(t))
}

// Return the t parameters for the intersection between lines l0 and l1
func (l0 Line2) Intersect(l1 Line2) (float64, float64, error) {
	m := M22{l0.v.X, -l1.v.X, l0.v.Y, -l1.v.Y}
	if m.Determinant() == 0 {
		return 0, 0, fmt.Errorf("zero/many")
	}
	p := l1.a.Sub(l0.a)
	t := m.Inverse().MulPosition(p)
	return t.X, t.Y, nil
}

// return the distance to the line, > 0 implies to the right of the line vector
func (l Line2) Distance(p V2) float64 {

	n := V2{l.v.Y, -l.v.X} // normal to line
	ap := p.Sub(l.a)       // line from a to p
	dn := ap.Dot(n)        // normal distance to line

	var d float64
	if l.segment {
		// this is a line segment - consider endpoints
		t := ap.Dot(l.v) // t-parameter of projection onto line
		if t < 0 {
			d = ap.Length()
		} else if t > l.length {
			bp := p.Sub(l.b) // line from b to p
			d = bp.Length()
		} else {
			// return the normal distance
			return dn
		}
	} else {
		// not a line segment - just return the normal distance
		return dn
	}

	if dn < 0 {
		d = -d
	}
	return d
}

//-----------------------------------------------------------------------------
