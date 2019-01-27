//-----------------------------------------------------------------------------
/*

2D line segments

Used for building 2D polygons SDFs.

*/
//-----------------------------------------------------------------------------

package sdf

import "fmt"

//-----------------------------------------------------------------------------

// Line is a 2d line segment defined with 2 points.
type Line [2]V2

//-----------------------------------------------------------------------------
// 2D Lines

// Line2 is a 2d line defined as either point/point or point/vector.
type Line2 struct {
	segment bool    // is this a line segment?
	length  float64 // segment length
	a       V2      // line start point
	b       V2      // line end point point (if segment)
	v       V2      // normalized line vector
}

// NewLinePV returns a 2d line defined by a point and vector.
func newLinePV(p, v V2) Line2 {
	l := Line2{}
	l.segment = false
	l.length = 0.0
	l.a = p
	l.v = v.Normalize()
	return l
}

// NewLinePP returns a 2d line segment defined by 2 points.
func newLinePP(a, b V2) Line2 {
	l := Line2{}
	v := b.Sub(a)
	l.segment = true
	l.length = v.Length()
	l.a = a
	l.b = b
	l.v = v.Normalize()
	return l
}

// Position returns the position on the line given the t value.
func (l Line2) Position(t float64) V2 {
	return l.a.Add(l.v.MulScalar(t))
}

// Intersect returns the t parameters for the intersection between lines l and lx
func (l Line2) Intersect(lx Line2) (float64, float64, error) {
	m := M22{l.v.X, -lx.v.X, l.v.Y, -lx.v.Y}
	if m.Determinant() == 0 {
		return 0, 0, fmt.Errorf("zero/many")
	}
	p := lx.a.Sub(l.a)
	t := m.Inverse().MulPosition(p)
	return t.X, t.Y, nil
}

// Distance returns the distance to the line.
// Greater than 0 implies to the right of the line vector.
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
