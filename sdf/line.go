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
	p       V2      // line start point
	v       V2      // normalized line vector
}

// Create a new line given a point and vector
func NewLine2_PV(p, v V2) Line2 {
	l := Line2{}
	l.segment = false
	l.length = 0.0
	l.p = p
	l.v = v.Normalize()
	return l
}

// Return the position given the t value
func (l *Line2) Position(t float64) V2 {
	return l.p.Add(l.v.MulScalar(t))
}

// Return the ta and tb parameters for the intersection between lines a and b
func (a Line2) Intersect(b Line2) (V2, error) {
	m := M22{a.v.X, -b.v.X, a.v.Y, -b.v.Y}

	if m.Determinant() == 0 {
		return V2{}, fmt.Errorf("no intersection")
	}

	p := b.p.Sub(a.p)
	return m.Inverse().MulPosition(p), nil
}

//-----------------------------------------------------------------------------
