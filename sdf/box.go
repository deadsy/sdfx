//-----------------------------------------------------------------------------
/*

 */
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"math"
	"math/rand"

	"github.com/deadsy/sdfx/vec/conv"
	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v2i"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// Box3 is a 3d bounding box.
type Box3 struct {
	Min, Max v3.Vec
}

// Box2 is a 2d bounding box.
type Box2 struct {
	Min, Max v2.Vec
}

//-----------------------------------------------------------------------------

// NewBox3 creates a 3d box with a given center and size.
func NewBox3(center, size v3.Vec) Box3 {
	half := size.MulScalar(0.5)
	return Box3{center.Sub(half), center.Add(half)}
}

// NewBox2 creates a 2d box with a given center and size.
func NewBox2(center, size v2.Vec) Box2 {
	half := size.MulScalar(0.5)
	return Box2{center.Sub(half), center.Add(half)}
}

//-----------------------------------------------------------------------------

// Equals test the equality of 3d boxes.
func (a Box3) Equals(b Box3, tolerance float64) bool {
	return (a.Min.Equals(b.Min, tolerance) && a.Max.Equals(b.Max, tolerance))
}

// Equals test the equality of 2d boxes.
func (a Box2) Equals(b Box2, tolerance float64) bool {
	return (a.Min.Equals(b.Min, tolerance) && a.Max.Equals(b.Max, tolerance))
}

//-----------------------------------------------------------------------------

// Extend returns a box enclosing two 3d boxes.
func (a Box3) Extend(b Box3) Box3 {
	return Box3{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

// Extend returns a box enclosing two 2d boxes.
func (a Box2) Extend(b Box2) Box2 {
	return Box2{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

// Include enlarges a 3d box to include a point.
func (a Box3) Include(v v3.Vec) Box3 {
	return Box3{a.Min.Min(v), a.Max.Max(v)}
}

// Include enlarges a 2d box to include a point.
func (a Box2) Include(v v2.Vec) Box2 {
	return Box2{a.Min.Min(v), a.Max.Max(v)}
}

//-----------------------------------------------------------------------------

// Translate translates a 3d box.
func (a Box3) Translate(v v3.Vec) Box3 {
	return Box3{a.Min.Add(v), a.Max.Add(v)}
}

// Translate translates a 2d box.
func (a Box2) Translate(v v2.Vec) Box2 {
	return Box2{a.Min.Add(v), a.Max.Add(v)}
}

//-----------------------------------------------------------------------------

// Size returns the size of a 3d box.
func (a Box3) Size() v3.Vec {
	return a.Max.Sub(a.Min)
}

// Size returns the size of a 2d box.
func (a Box2) Size() v2.Vec {
	return a.Max.Sub(a.Min)
}

// Center returns the center of a 3d box.
func (a Box3) Center() v3.Vec {
	return a.Min.Add(a.Size().MulScalar(0.5))
}

// Center returns the center of a 2d box.
func (a Box2) Center() v2.Vec {
	return a.Min.Add(a.Size().MulScalar(0.5))
}

//-----------------------------------------------------------------------------

// ScaleAboutCenter returns a new 2d box scaled about the center of a box.
func (a Box2) ScaleAboutCenter(k float64) Box2 {
	return NewBox2(a.Center(), a.Size().MulScalar(k))
}

// ScaleAboutCenter returns a new 3d box scaled about the center of a box.
func (a Box3) ScaleAboutCenter(k float64) Box3 {
	return NewBox3(a.Center(), a.Size().MulScalar(k))
}

//-----------------------------------------------------------------------------

// Enlarge returns a new 3d box enlarged by a size vector.
func (a Box3) Enlarge(v v3.Vec) Box3 {
	v = v.MulScalar(0.5)
	return Box3{a.Min.Sub(v), a.Max.Add(v)}
}

// Enlarge returns a new 2d box enlarged by a size vector.
func (a Box2) Enlarge(v v2.Vec) Box2 {
	v = v.MulScalar(0.5)
	return Box2{a.Min.Sub(v), a.Max.Add(v)}
}

//-----------------------------------------------------------------------------

// Contains checks if the 3d box contains the given vector (considering bounds as inside).
func (a Box3) Contains(v v3.Vec) bool {
	return a.Min.X <= v.X && a.Min.Y <= v.Y && a.Min.Z <= v.Z &&
		v.X <= a.Max.X && v.Y <= a.Max.Y && v.Z <= a.Max.Z
}

// Contains checks if the 2d box contains the given vector (considering bounds as inside).
func (a Box2) Contains(v v2.Vec) bool {
	return a.Min.X <= v.X && a.Min.Y <= v.Y &&
		v.X <= a.Max.X && v.Y <= a.Max.Y
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

// Vertices returns a slice of 3d box corner vertices.
func (a Box3) Vertices() v3.VecSet {
	v := make([]v3.Vec, 8)
	v[0] = a.Min
	v[1] = v3.Vec{a.Min.X, a.Min.Y, a.Max.Z}
	v[2] = v3.Vec{a.Min.X, a.Max.Y, a.Min.Z}
	v[3] = v3.Vec{a.Min.X, a.Max.Y, a.Max.Z}
	v[4] = v3.Vec{a.Max.X, a.Min.Y, a.Min.Z}
	v[5] = v3.Vec{a.Max.X, a.Min.Y, a.Max.Z}
	v[6] = v3.Vec{a.Max.X, a.Max.Y, a.Min.Z}
	v[7] = a.Max
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

// MinMaxDist2 returns the minimum and maximum dist * dist from a point to a box.
// Points within the box have minimum distance = 0.
func (a Box3) MinMaxDist2(p v3.Vec) v2.Vec {
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

	// consider the faces (for the minimum)
	withinX := a.Min.X < 0 && a.Max.X > 0
	withinY := a.Min.Y < 0 && a.Max.Y > 0
	withinZ := a.Min.Z < 0 && a.Max.Z > 0

	if withinX && withinY && withinZ {
		minDist2 = 0
	} else {
		if withinX && withinY {
			d := math.Min(math.Abs(a.Max.Z), math.Abs(a.Min.Z))
			minDist2 = math.Min(minDist2, d*d)
		}
		if withinX && withinZ {
			d := math.Min(math.Abs(a.Max.Y), math.Abs(a.Min.Y))
			minDist2 = math.Min(minDist2, d*d)
		}
		if withinY && withinZ {
			d := math.Min(math.Abs(a.Max.X), math.Abs(a.Min.X))
			minDist2 = math.Min(minDist2, d*d)
		}
	}

	return v2.Vec{minDist2, maxDist2}
}

//-----------------------------------------------------------------------------

// randomRange returns a random float64 [a,b)
func randomRange(a, b float64) float64 {
	return a + (b-a)*rand.Float64()
}

// Random returns a random point within a bounding box.
func (a *Box2) Random() v2.Vec {
	return v2.Vec{
		randomRange(a.Min.X, a.Max.X),
		randomRange(a.Min.Y, a.Max.Y),
	}
}

// Random returns a random point within a bounding box.
func (a *Box3) Random() v3.Vec {
	return v3.Vec{
		randomRange(a.Min.X, a.Max.X),
		randomRange(a.Min.Y, a.Max.Y),
		randomRange(a.Min.Z, a.Max.Z),
	}
}

// RandomSet returns a set of random points from within a bounding box.
func (a *Box2) RandomSet(n int) v2.VecSet {
	s := make([]v2.Vec, n)
	for i := range s {
		s[i] = a.Random()
	}
	return s
}

// RandomSet returns a set of random points from within a bounding box.
func (a *Box3) RandomSet(n int) v3.VecSet {
	s := make([]v3.Vec, n)
	for i := range s {
		s[i] = a.Random()
	}
	return s
}

//-----------------------------------------------------------------------------
