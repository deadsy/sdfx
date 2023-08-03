//-----------------------------------------------------------------------------
/*

2D Boxes

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"math"

	"github.com/deadsy/sdfx/vec/conv"
	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v2i"
)

//-----------------------------------------------------------------------------

// Box2 is a 2d bounding box.
type Box2 struct {
	Min, Max v2.Vec
}

// NewBox2 creates a 2d box with a given center and size.
func NewBox2(center, size v2.Vec) Box2 {
	half := size.MulScalar(0.5)
	return Box2{center.Sub(half), center.Add(half)}
}

// Equals test the equality of 2d boxes.
func (a Box2) Equals(b Box2, tolerance float64) bool {
	return (a.Min.Equals(b.Min, tolerance) && a.Max.Equals(b.Max, tolerance))
}

// Extend returns a box enclosing two 2d boxes.
func (a Box2) Extend(b Box2) Box2 {
	return Box2{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

// Include enlarges a 2d box to include a point.
func (a Box2) Include(v v2.Vec) Box2 {
	return Box2{a.Min.Min(v), a.Max.Max(v)}
}

// Translate translates a 2d box.
func (a Box2) Translate(v v2.Vec) Box2 {
	return Box2{a.Min.Add(v), a.Max.Add(v)}
}

// Size returns the size of a 2d box.
func (a Box2) Size() v2.Vec {
	return a.Max.Sub(a.Min)
}

// Center returns the center of a 2d box.
func (a Box2) Center() v2.Vec {
	return a.Min.Add(a.Size().MulScalar(0.5))
}

// ScaleAboutCenter returns a new 2d box scaled about the center of a box.
func (a Box2) ScaleAboutCenter(k float64) Box2 {
	return NewBox2(a.Center(), a.Size().MulScalar(k))
}

// Enlarge returns a new 2d box enlarged by a size vector.
func (a Box2) Enlarge(v v2.Vec) Box2 {
	v = v.MulScalar(0.5)
	return Box2{a.Min.Sub(v), a.Max.Add(v)}
}

// Contains checks if the 2d box contains the vector.
// Note: Min boundary is in, Max boundary is out.
func (a Box2) Contains(v v2.Vec) bool {
	return a.Min.X <= v.X && a.Min.Y <= v.Y &&
		v.X < a.Max.X && v.Y < a.Max.Y
}

//-----------------------------------------------------------------------------

// Quad0 returns the 0th quadtree box of a box (lower-left/south-west).
func (a Box2) Quad0() Box2 {
	delta := a.Size().MulScalar(0.5)
	ll := a.Min
	return Box2{ll, ll.Add(delta)}
}

// Quad1 returns the 1st quadtree box of a box (lower-right/south-east).
func (a Box2) Quad1() Box2 {
	delta := a.Size().MulScalar(0.5)
	ll := v2.Vec{a.Min.X + delta.X, a.Min.Y}
	return Box2{ll, ll.Add(delta)}
}

// Quad2 returns the 2nd quadtree box of a box (top-left/north-west).
func (a Box2) Quad2() Box2 {
	delta := a.Size().MulScalar(0.5)
	ll := v2.Vec{a.Min.X, a.Min.Y + delta.Y}
	return Box2{ll, ll.Add(delta)}
}

// Quad3 returns the 3rd quadtree box of a box (top-right/north-east).
func (a Box2) Quad3() Box2 {
	delta := a.Size().MulScalar(0.5)
	ll := a.Min.Add(delta)
	return Box2{ll, ll.Add(delta)}
}

//-----------------------------------------------------------------------------

// Vertices returns a slice of 2d box corner vertices.
func (a Box2) Vertices() v2.VecSet {
	v := make([]v2.Vec, 4)
	v[0] = a.Min                    // bl
	v[1] = v2.Vec{a.Max.X, a.Min.Y} // br
	v[2] = v2.Vec{a.Min.X, a.Max.Y} // tl
	v[3] = a.Max                    // tr
	return v
}

// BottomLeft returns the bottom left corner of a 2d bounding box.
func (a Box2) BottomLeft() v2.Vec {
	return a.Min
}

// TopLeft returns the top left corner of a 2d bounding box.
func (a Box2) TopLeft() v2.Vec {
	return v2.Vec{a.Min.X, a.Max.Y}
}

//-----------------------------------------------------------------------------

// Map2 maps a 2d region to integer grid coordinates.
type Map2 struct {
	bb    Box2    // bounding box
	grid  v2i.Vec // integral dimension
	delta v2.Vec
	flipy bool // flip the y-axis
}

// NewMap2 returns a 2d region to grid coordinates map.
func NewMap2(bb Box2, grid v2i.Vec, flipy bool) (*Map2, error) {
	// sanity check the bounding box
	bbSize := bb.Size()
	if bbSize.X <= 0 || bbSize.Y <= 0 {
		return nil, errors.New("bad bounding box")
	}
	// sanity check the integer dimensions
	if grid.X <= 0 || grid.Y <= 0 {
		return nil, errors.New("bad grid dimensions")
	}
	m := Map2{}
	m.bb = bb
	m.grid = grid
	m.flipy = flipy
	m.delta = bbSize.Div(conv.V2iToV2(grid))
	return &m, nil
}

// ToV2 converts grid integer coordinates to 2d region float coordinates.
func (m *Map2) ToV2(p v2i.Vec) v2.Vec {
	ofs := conv.V2iToV2(p).AddScalar(0.5).Mul(m.delta)
	var origin v2.Vec
	if m.flipy {
		origin = m.bb.TopLeft()
		ofs.Y = -ofs.Y
	} else {
		origin = m.bb.BottomLeft()
	}
	return origin.Add(ofs)
}

// ToV2i converts 2d region float coordinates to grid integer coordinates.
func (m *Map2) ToV2i(p v2.Vec) v2i.Vec {
	var v v2.Vec
	if m.flipy {
		v = p.Sub(m.bb.TopLeft())
		v.Y = -v.Y
	} else {
		v = p.Sub(m.bb.BottomLeft())
	}
	return conv.V2ToV2i(v.Div(m.delta))
}

//-----------------------------------------------------------------------------
// Minimum/Maximum distances from a point to a box

// MinMaxDist2 returns the minimum and maximum dist * dist from a point to a box.
// Points within the box have minimum distance = 0.
func (a Box2) MinMaxDist2(p v2.Vec) v2.Vec {
	maxDist2 := 0.0
	minDist2 := 0.0

	// translate the box so p is at the origin
	a = a.Translate(p.Neg())

	// consider the vertices
	vs := a.Vertices()

	for i := range vs {
		d2 := vs[i].Length2()
		if i == 0 {
			minDist2 = d2
		} else {
			minDist2 = math.Min(minDist2, d2)
		}
		maxDist2 = math.Max(maxDist2, d2)
	}

	// consider the sides (for the minimum)
	withinX := a.Min.X < 0 && a.Max.X > 0
	withinY := a.Min.Y < 0 && a.Max.Y > 0

	if withinX && withinY {
		minDist2 = 0
	} else {
		if withinX {
			d := math.Min(math.Abs(a.Max.Y), math.Abs(a.Min.Y))
			minDist2 = math.Min(minDist2, d*d)
		}
		if withinY {
			d := math.Min(math.Abs(a.Max.X), math.Abs(a.Min.X))
			minDist2 = math.Min(minDist2, d*d)
		}
	}

	return v2.Vec{minDist2, maxDist2}
}

//-----------------------------------------------------------------------------

// Random returns a random point within a 2d box.
func (a *Box2) Random() v2.Vec {
	return v2.Vec{
		randomRange(a.Min.X, a.Max.X),
		randomRange(a.Min.Y, a.Max.Y),
	}
}

// RandomSet returns a set of random points from within a 2d box.
func (a *Box2) RandomSet(n int) v2.VecSet {
	s := make([]v2.Vec, n)
	for i := range s {
		s[i] = a.Random()
	}
	return s
}

//-----------------------------------------------------------------------------