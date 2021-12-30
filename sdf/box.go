//-----------------------------------------------------------------------------
/*

 */
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"math"
)

//-----------------------------------------------------------------------------

// Box3 is a 3d bounding box.
type Box3 struct {
	Min, Max V3
}

// Box2 is a 2d bounding box.
type Box2 struct {
	Min, Max V2
}

//-----------------------------------------------------------------------------

// NewBox3 creates a 3d box with a given center and size.
func NewBox3(center, size V3) Box3 {
	half := size.MulScalar(0.5)
	return Box3{center.Sub(half), center.Add(half)}
}

// NewBox2 creates a 2d box with a given center and size.
func NewBox2(center, size V2) Box2 {
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

// Extend returns a box that encloses two 3d boxes.
func (a Box3) Extend(b Box3) Box3 {
	return Box3{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

// Extend returns a box that encloses two 2d boxes.
func (a Box2) Extend(b Box2) Box2 {
	return Box2{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

//-----------------------------------------------------------------------------

// Translate translates a 3d box.
func (a Box3) Translate(v V3) Box3 {
	return Box3{a.Min.Add(v), a.Max.Add(v)}
}

// Translate translates a 2d box.
func (a Box2) Translate(v V2) Box2 {
	return Box2{a.Min.Add(v), a.Max.Add(v)}
}

//-----------------------------------------------------------------------------

// Size returns the size of a 3d box.
func (a Box3) Size() V3 {
	return a.Max.Sub(a.Min)
}

// Size returns the size of a 2d box.
func (a Box2) Size() V2 {
	return a.Max.Sub(a.Min)
}

// Center returns the center of a 3d box.
func (a Box3) Center() V3 {
	return a.Min.Add(a.Size().MulScalar(0.5))
}

// Center returns the center of a 2d box.
func (a Box2) Center() V2 {
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
func (a Box3) Enlarge(v V3) Box3 {
	v = v.MulScalar(0.5)
	return Box3{a.Min.Sub(v), a.Max.Add(v)}
}

// Enlarge returns a new 2d box enlarged by a size vector.
func (a Box2) Enlarge(v V2) Box2 {
	v = v.MulScalar(0.5)
	return Box2{a.Min.Sub(v), a.Max.Add(v)}
}

//-----------------------------------------------------------------------------

// Include makes sure that the box includes the given point by extending it if necessary, returning the new Box3
func (a Box3) Include(vertex V3) Box3 {
	if vertex.X < a.Min.X {
		a.Min.X = vertex.X
	}
	if vertex.Y < a.Min.Y {
		a.Min.Y = vertex.Y
	}
	if vertex.Z < a.Min.Z {
		a.Min.Z = vertex.Z
	}
	if vertex.X > a.Max.X {
		a.Max.X = vertex.X
	}
	if vertex.Y > a.Max.Y {
		a.Max.Y = vertex.Y
	}
	if vertex.Z > a.Max.Z {
		a.Max.Z = vertex.Z
	}
	return a // It is a copy (not passed by reference)
}

// Include makes sure that the box includes the given point by extending it if necessary, returning the new Box2
func (a Box2) Include(vertex V2) Box2 {
	if vertex.X < a.Min.X {
		a.Min.X = vertex.X
	}
	if vertex.Y < a.Min.Y {
		a.Min.Y = vertex.Y
	}
	if vertex.X > a.Max.X {
		a.Max.X = vertex.X
	}
	if vertex.Y > a.Max.Y {
		a.Max.Y = vertex.Y
	}
	return a // It is a copy (not passed by reference)
}

//-----------------------------------------------------------------------------

// Contains checks if the 3d box contains the given vector (considering bounds as inside).
func (a Box3) Contains(v V3) bool {
	return a.Min.X <= v.X && a.Min.Y <= v.Y && a.Min.Z <= v.Z &&
		v.X <= a.Max.X && v.Y <= a.Max.Y && v.Z <= a.Max.Z
}

// Contains checks if the 2d box contains the given vector (considering bounds as inside).
func (a Box2) Contains(v V2) bool {
	return a.Min.X <= v.X && a.Min.Y <= v.Y &&
		v.X <= a.Max.X && v.Y <= a.Max.Y
}

//-----------------------------------------------------------------------------

// Vertices returns a slice of 2d box corner vertices.
func (a Box2) Vertices() V2Set {
	v := make([]V2, 4)
	v[0] = a.Min                // bl
	v[1] = V2{a.Max.X, a.Min.Y} // br
	v[2] = V2{a.Min.X, a.Max.Y} // tl
	v[3] = a.Max                // tr
	return v
}

// Vertices returns a slice of 3d box corner vertices.
func (a Box3) Vertices() V3Set {
	v := make([]V3, 8)
	v[0] = a.Min
	v[1] = V3{a.Min.X, a.Min.Y, a.Max.Z}
	v[2] = V3{a.Min.X, a.Max.Y, a.Min.Z}
	v[3] = V3{a.Min.X, a.Max.Y, a.Max.Z}
	v[4] = V3{a.Max.X, a.Min.Y, a.Min.Z}
	v[5] = V3{a.Max.X, a.Min.Y, a.Max.Z}
	v[6] = V3{a.Max.X, a.Max.Y, a.Min.Z}
	v[7] = a.Max
	return v
}

// BottomLeft returns the bottom left corner of a 2d bounding box.
func (a Box2) BottomLeft() V2 {
	return a.Min
}

// TopLeft returns the top left corner of a 2d bounding box.
func (a Box2) TopLeft() V2 {
	return V2{a.Min.X, a.Max.Y}
}

//-----------------------------------------------------------------------------

// Map2 maps a 2d region to integer grid coordinates.
type Map2 struct {
	bb    Box2 // bounding box
	grid  V2i  // integral dimension
	delta V2
	flipy bool // flip the y-axis
}

// NewMap2 returns a 2d region to grid coordinates map.
func NewMap2(bb Box2, grid V2i, flipy bool) (*Map2, error) {
	// sanity check the bounding box
	bbSize := bb.Size()
	if bbSize.X <= 0 || bbSize.Y <= 0 {
		return nil, errors.New("bad bounding box")
	}
	// sanity check the integer dimensions
	if grid[0] <= 0 || grid[1] <= 0 {
		return nil, errors.New("bad grid dimensions")
	}
	m := Map2{}
	m.bb = bb
	m.grid = grid
	m.flipy = flipy
	m.delta = bbSize.Div(grid.ToV2())
	return &m, nil
}

// ToV2 converts grid integer coordinates to 2d region float coordinates.
func (m *Map2) ToV2(p V2i) V2 {
	ofs := p.ToV2().AddScalar(0.5).Mul(m.delta)
	var origin V2
	if m.flipy {
		origin = m.bb.TopLeft()
		ofs.Y = -ofs.Y
	} else {
		origin = m.bb.BottomLeft()
	}
	return origin.Add(ofs)
}

// ToV2i converts 2d region float coordinates to grid integer coordinates.
func (m *Map2) ToV2i(p V2) V2i {
	var v V2
	if m.flipy {
		v = p.Sub(m.bb.TopLeft())
		v.Y = -v.Y
	} else {
		v = p.Sub(m.bb.BottomLeft())
	}
	return v.Div(m.delta).ToV2i()
}

//-----------------------------------------------------------------------------
// Minimum/Maximum distances from a point to a box

// MinMaxDist2 returns the minimum and maximum dist * dist from a point to a box.
// Points within the box have minimum distance = 0.
func (a Box2) MinMaxDist2(p V2) V2 {
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

	return V2{minDist2, maxDist2}
}

// MinMaxDist2 returns the minimum and maximum dist * dist from a point to a box.
// Points within the box have minimum distance = 0.
func (a Box3) MinMaxDist2(p V3) V2 {
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

	return V2{minDist2, maxDist2}
}

//-----------------------------------------------------------------------------
