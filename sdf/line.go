//-----------------------------------------------------------------------------
/*

2D lines

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// Interval is a closed interval on real numbers.
type Interval [2]float64

// Sort sorts the interval endpoints lowest to highest.
func (a Interval) Sort() Interval {
	if a[0] <= a[1] {
		return a
	}
	return Interval{a[1], a[0]}
}

// Equals returns true if a == b within the tolerance limit.
func (a Interval) Equals(b Interval, tolerance float64) bool {
	return math.Abs(a[0]-b[0]) <= tolerance && math.Abs(a[1]-b[1]) <= tolerance
}

// Overlap returns true if two intervals overlap.
func (a Interval) Overlap(b Interval) bool {
	return b[0] <= a[1] && a[0] <= b[1]
}

// Intersect returns the intersection of two intervals.
func (a Interval) Intersect(b Interval) *Interval {
	if a.Overlap(b) {
		return &Interval{math.Max(a[0], b[0]), math.Min(a[1], b[1])}
	}
	return nil
}

//-----------------------------------------------------------------------------

// Line2 is a 2d line defined with end-points.
type Line2 [2]v2.Vec

// BoundingBox returns a bounding box for the line.
func (a *Line2) BoundingBox() Box2 {
	return Box2{Min: a[0], Max: a[0]}.Include(a[1])
}

// Reverse the direction of a line segment.
func (a *Line2) Reverse() *Line2 {
	return &Line2{a[1], a[0]}
}

// Equals returns true if the lines are the same (within tolerance).
func (a *Line2) Equals(b *Line2, tolerance float64) bool {
	return a[0].Equals(b[0], tolerance) && a[1].Equals(b[1], tolerance)
}

// Degenerate returns true if the line is degenerate.
func (a Line2) Degenerate(tolerance float64) bool {
	// check for identical vertices
	return a[0].Equals(a[1], tolerance)
}

// IntersectLine intersects 2 line segments.
// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
func (a *Line2) IntersectLine(b *Line2) []v2.Vec {

	p := a[0]
	r := a[1].Sub(a[0])
	q := b[0]
	s := b[1].Sub(b[0])

	k0 := r.Cross(s)        // r x s
	k1 := q.Sub(p).Cross(r) // (q - p) x r

	if k0 == 0 {
		if k1 != 0 {
			// parallel, non-intersecting
			return nil
		}

		// collinear lines
		k2 := 1.0 / r.Dot(r)
		t0 := q.Sub(p).Dot(r) * k2
		t1 := t0 + s.Dot(r)*k2

		t := Interval{t0, t1}.Sort()
		x := t.Intersect(Interval{0, 1})
		if x != nil {
			// collinear, intersecting
			p0 := p.Add(r.MulScalar(x[0]))
			if x[0] == x[1] {
				return []v2.Vec{p0}
			}
			p1 := p.Add(r.MulScalar(x[1]))
			return []v2.Vec{p0, p1}
		}

		// collinear, non-intersecting
		return nil
	}
	// non-parallel
	u := k1 / k0
	t := q.Sub(p).Cross(s) / k0
	if u >= 0 && u <= 1 && t >= 0 && t <= 1 {
		p0 := p.Add(r.MulScalar(t))
		return []v2.Vec{p0}
	}
	// non-parallel, non-intersecting
	return nil
}

//-----------------------------------------------------------------------------

// geometryLine is a 2d line defined as either point/point or point/vector.
type geometryLine struct {
	segment bool    // is this a line segment?
	length  float64 // segment length
	a       v2.Vec  // line start point
	b       v2.Vec  // line end point point (if segment)
	v       v2.Vec  // normalized line vector
}

// NewLinePV returns a 2d line defined by a point and vector.
func newLinePV(p, v v2.Vec) geometryLine {
	l := geometryLine{}
	l.segment = false
	l.length = 0.0
	l.a = p
	l.v = v.Normalize()
	return l
}

// NewLinePP returns a 2d line segment defined by 2 points.
func newLinePP(a, b v2.Vec) geometryLine {
	l := geometryLine{}
	v := b.Sub(a)
	l.segment = true
	l.length = v.Length()
	l.a = a
	l.b = b
	l.v = v.Normalize()
	return l
}

// Position returns the position on the line given the t value.
func (l geometryLine) Position(t float64) v2.Vec {
	return l.a.Add(l.v.MulScalar(t))
}

// Intersect returns the t parameters for the intersection between lines l and lx
func (l geometryLine) Intersect(lx geometryLine) (float64, float64, error) {
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
func (l geometryLine) Distance(p v2.Vec) float64 {

	n := v2.Vec{l.v.Y, -l.v.X} // normal to line
	ap := p.Sub(l.a)           // line from a to p
	dn := ap.Dot(n)            // normal distance to line

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
