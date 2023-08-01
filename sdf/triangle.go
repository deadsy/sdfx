//-----------------------------------------------------------------------------
/*

Triangles

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"math"
	"sync"

	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/dhconnelly/rtreego"
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
		return v2.Vec{}, errors.New("coincident points")
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

// Triangle3 is a 3D triangle
type Triangle3 [3]v3.Vec

// Normal returns the normal vector to the plane defined by the 3D triangle.
func (t *Triangle3) Normal() v3.Vec {
	e1 := t[1].Sub(t[0])
	e2 := t[2].Sub(t[0])
	return e1.Cross(e2).Normalize()
}

// Degenerate returns true if the triangle is degenerate.
func (t *Triangle3) Degenerate(tolerance float64) bool {
	// check for identical vertices
	if t[0].Equals(t[1], tolerance) {
		return true
	}
	if t[1].Equals(t[2], tolerance) {
		return true
	}
	if t[2].Equals(t[0], tolerance) {
		return true
	}
	// TODO more tests needed
	return false
}

func v3ToPoint(v v3.Vec) rtreego.Point {
	return rtreego.Point{v.X, v.Y, v.Z}
}

// BoundingBox returns a bounding box for the triangle.
func (t *Triangle3) BoundingBox() Box3 {
	return Box3{Min: t[0], Max: t[0]}.Include(t[1]).Include(t[2])
}

// Bounds returns a r-tree bounding rectangle for the triangle.
func (t *Triangle3) Bounds() *rtreego.Rect {
	b := t.BoundingBox()
	r, _ := rtreego.NewRectFromPoints(v3ToPoint(b.Min), v3ToPoint(b.Max))
	return r
}

//-----------------------------------------------------------------------------

// WriteTriangles writes a stream of triangles to a slice.
func WriteTriangles(wg *sync.WaitGroup, triangles *[]Triangle3) chan<- []*Triangle3 {
	// External code writes triangles to this channel.
	// This goroutine reads the channel and appends the triangles to a slice.
	c := make(chan []*Triangle3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// read triangles from the channel and append them to the slice
		for ts := range c {
			for _, t := range ts {
				*triangles = append(*triangles, *t)
			}
		}
	}()

	return c
}

//-----------------------------------------------------------------------------
