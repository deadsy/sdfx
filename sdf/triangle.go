//-----------------------------------------------------------------------------
/*

Triangles

*/
//-----------------------------------------------------------------------------

package sdf

import "errors"

//-----------------------------------------------------------------------------

// 3d triangle
type Triangle3 struct {
	V [3]V3
}

// 2d triangle
type Triangle2 [3]V2

//-----------------------------------------------------------------------------

func NewTriangle3(a, b, c V3) *Triangle3 {
	t := Triangle3{}
	t.V[0] = a
	t.V[1] = b
	t.V[2] = c
	return &t
}

//-----------------------------------------------------------------------------

// return the normal vector to the plane defined by the triangle
func (t *Triangle3) Normal() V3 {
	e1 := t.V[1].Sub(t.V[0])
	e2 := t.V[2].Sub(t.V[0])
	return e1.Cross(e2).Normalize()
}

//-----------------------------------------------------------------------------

// Return true if the point is within the circumcircle of the triangle.
// See: http://www.mathopenref.com/trianglecircumcircle.html
// See: http://paulbourke.net/papers/triangulate/
func (t Triangle2) InCircumcircle(p V2) bool {

	var m1, m2, mx1, mx2, my1, my2 float64
	var dx, dy, drsqr float64
	var xc, yc, rsqr float64

	x1 := t[0].X
	x2 := t[1].X
	x3 := t[2].X

	y1 := t[0].Y
	y2 := t[1].Y
	y3 := t[2].Y

	xp := p.X
	yp := p.Y

	fabsy1y2 := Abs(y1 - y2)
	fabsy2y3 := Abs(y2 - y3)

	// Check for coincident points
	if fabsy1y2 < EPSILON && fabsy2y3 < EPSILON {
		return false
	}

	if fabsy1y2 < EPSILON {
		m2 = -(x3 - x2) / (y3 - y2)
		mx2 = (x2 + x3) / 2.0
		my2 = (y2 + y3) / 2.0
		xc = (x2 + x1) / 2.0
		yc = m2*(xc-mx2) + my2
	} else if fabsy2y3 < EPSILON {
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

	dx = x2 - xc
	dy = y2 - yc
	rsqr = dx*dx + dy*dy

	dx = xp - xc
	dy = yp - yc
	drsqr = dx*dx + dy*dy

	return (drsqr - rsqr) <= EPSILON
}

//-----------------------------------------------------------------------------

func (t Triangle2) Circumcenter() (V2, error) {

	var m1, m2, mx1, mx2, my1, my2 float64
	var xc, yc float64

	x1 := t[0].X
	x2 := t[1].X
	x3 := t[2].X

	y1 := t[0].Y
	y2 := t[1].Y
	y3 := t[2].Y

	fabsy1y2 := Abs(y1 - y2)
	fabsy2y3 := Abs(y2 - y3)

	// Check for coincident points
	if fabsy1y2 < EPSILON && fabsy2y3 < EPSILON {
		return V2{}, errors.New("coincident points")
	}

	if fabsy1y2 < EPSILON {
		m2 = -(x3 - x2) / (y3 - y2)
		mx2 = (x2 + x3) / 2.0
		my2 = (y2 + y3) / 2.0
		xc = (x2 + x1) / 2.0
		yc = m2*(xc-mx2) + my2
	} else if fabsy2y3 < EPSILON {
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

	return V2{xc, yc}, nil
}

//-----------------------------------------------------------------------------
