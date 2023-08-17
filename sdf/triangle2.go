//-----------------------------------------------------------------------------
/*

2D Triangles

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// Triangle2 is a 2D triangle
type Triangle2 [3]v2.Vec

// Circumcenter returns the circumcenter of a triangle.
func (t Triangle2) Circumcenter() (v2.Vec, error) {

	var m1, m2, mx1, mx2, my1, my2 float64
	var xc, yc float64

	x1 := t[0].X
	x2 := t[1].X
	x3 := t[2].X

	y1 := t[0].Y
	y2 := t[1].Y
	y3 := t[2].Y

	fabsy1y2 := math.Abs(y1 - y2)
	fabsy2y3 := math.Abs(y2 - y3)

	// Check for coincident points
	if fabsy1y2 < epsilon && fabsy2y3 < epsilon {
		return v2.Vec{}, ErrMsg("coincident points")
	}

	if fabsy1y2 < epsilon {
		m2 = -(x3 - x2) / (y3 - y2)
		mx2 = (x2 + x3) / 2.0
		my2 = (y2 + y3) / 2.0
		xc = (x2 + x1) / 2.0
		yc = m2*(xc-mx2) + my2
	} else if fabsy2y3 < epsilon {
		m1 = -(x2 - x1) / (y2 - y1)
		mx1 = (x1 + x2) / 2.0
		my1 = (y1 + y2) / 2.0
		xc = (x3 + x2) / 2.0
		yc = m1*(xc-mx1) + my1
	} else {
		m1 = -(x2 - x1) / (y2 - y1)
		m2 = -(x3 - x2) / (y3 - y2)
		mx1 = (x1 + x2) / 2.0
		mx2 = (x2 + x3) / 2.0
		my1 = (y1 + y2) / 2.0
		my2 = (y2 + y3) / 2.0
		xc = (m1*mx1 - m2*mx2 + my2 - my1) / (m1 - m2)
		if fabsy1y2 > fabsy2y3 {
			yc = m1*(xc-mx1) + my1
		} else {
			yc = m2*(xc-mx2) + my2
		}
	}

	return v2.Vec{xc, yc}, nil
}

// InCircumcircle return inside == true if the point is inside the circumcircle of the triangle.
// Returns done == true if the vertex and the subsequent x-ordered vertices are outside the circumcircle.
func (t Triangle2) InCircumcircle(p v2.Vec) (inside, done bool) {
	c, err := t.Circumcenter()
	if err != nil {
		inside = false
		done = true
		return
	}

	// radius squared of circumcircle
	dx := t[0].X - c.X
	dy := t[0].Y - c.Y
	r2 := dx*dx + dy*dy

	// distance squared from circumcenter to point
	dx = p.X - c.X
	dy = p.Y - c.Y
	d2 := dx*dx + dy*dy

	// is the point within the circumcircle?
	inside = d2-r2 <= epsilon

	// If this vertex has an x-value beyond the circumcenter and the distance based on the x-delta
	// is greater than the circumradius, then this triangle is done for this and all subsequent vertices
	// since the vertex list has been sorted by x-value.
	done = (dx > 0) && (dx*dx > r2)

	return
}

//-----------------------------------------------------------------------------
