//-----------------------------------------------------------------------------
/*

2D lines

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"

	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/dhconnelly/rtreego"
)

//-----------------------------------------------------------------------------

// Line2 is a 2d line defined with end-points.
type Line2 [2]v2.Vec

func v2ToPoint(v v2.Vec) rtreego.Point {
	return rtreego.Point{v.X, v.Y}
}

// BoundingBox returns a bounding box for the line.
func (l *Line2) BoundingBox() Box2 {
	return Box2{Min: l[0], Max: l[0]}.Include(l[1])
}

// Bounds returns a r-tree bounding rectangle for the line.
func (l *Line2) Bounds() *rtreego.Rect {
	b := l.BoundingBox()
	r, _ := rtreego.NewRectFromPoints(v2ToPoint(b.Min), v2ToPoint(b.Max))
	return r
}

// Degenerate returns true if the line is degenerate.
func (l Line2) Degenerate(tolerance float64) bool {
	// check for identical vertices
	return l[0].Equals(l[1], tolerance)
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
