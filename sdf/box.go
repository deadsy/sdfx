//-----------------------------------------------------------------------------
/*

 */
//-----------------------------------------------------------------------------

package sdf

import "errors"

//-----------------------------------------------------------------------------

type Box3 struct {
	Min, Max V3
}

type Box2 struct {
	Min, Max V2
}

//-----------------------------------------------------------------------------

// create a new Box with a given center and size
func NewBox3(center, size V3) Box3 {
	half := size.MulScalar(0.5)
	return Box3{center.Sub(half), center.Add(half)}
}

// create a new Box with a given center and size
func NewBox2(center, size V2) Box2 {
	half := size.MulScalar(0.5)
	return Box2{center.Sub(half), center.Add(half)}
}

//-----------------------------------------------------------------------------

// are the boxes equal?
func (a Box3) Equals(b Box3, tolerance float64) bool {
	return (a.Min.Equals(b.Min, tolerance) && a.Max.Equals(b.Max, tolerance))
}

// are the boxes equal?
func (a Box2) Equals(b Box2, tolerance float64) bool {
	return (a.Min.Equals(b.Min, tolerance) && a.Max.Equals(b.Max, tolerance))
}

//-----------------------------------------------------------------------------

// return a box that encloses two boxes
func (a Box3) Extend(b Box3) Box3 {
	return Box3{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

// return a box that encloses two boxes
func (a Box2) Extend(b Box2) Box2 {
	return Box2{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

//-----------------------------------------------------------------------------

// translate a box a distance v
func (a Box3) Translate(v V3) Box3 {
	return Box3{a.Min.Add(v), a.Max.Add(v)}
}

// translate a box a distance v
func (a Box2) Translate(v V2) Box2 {
	return Box2{a.Min.Add(v), a.Max.Add(v)}
}

//-----------------------------------------------------------------------------

// return the size of the box
func (a Box3) Size() V3 {
	return a.Max.Sub(a.Min)
}

// return the size of the box
func (a Box2) Size() V2 {
	return a.Max.Sub(a.Min)
}

// return the center of the box
func (a Box3) Center() V3 {
	return a.Min.Add(a.Size().MulScalar(0.5))
}

// return the center of the box
func (a Box2) Center() V2 {
	return a.Min.Add(a.Size().MulScalar(0.5))
}

//-----------------------------------------------------------------------------

// ScaleAboutCenter returns a new box scaled about the center of a box.
func (a Box2) ScaleAboutCenter(k float64) Box2 {
	return NewBox2(a.Center(), a.Size().MulScalar(k))
}

// ScaleAboutCenter returns a new box scaled about the center of a box.
func (a Box3) ScaleAboutCenter(k float64) Box3 {
	return NewBox3(a.Center(), a.Size().MulScalar(k))
}

//-----------------------------------------------------------------------------

// Return a slice of box vertices
func (a Box2) Vertices() V2Set {
	v := make([]V2, 4)
	v[0] = a.Min                // bl
	v[1] = V2{a.Max.X, a.Min.Y} // br
	v[2] = V2{a.Min.X, a.Max.Y} // tl
	v[3] = a.Max                // tr
	return v
}

// Return a slice of box vertices
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

// return the bottom left of a 2d bounding box
func (a Box2) BottomLeft() V2 {
	return a.Min
}

// return the top left of a 2d bounding box
func (a Box2) TopLeft() V2 {
	return V2{a.Min.X, a.Max.Y}
}

//-----------------------------------------------------------------------------
// Map a 2d region to integer grid coordinates

type Map2 struct {
	bb    Box2 // bounding box
	grid  V2i  // integral dimension
	delta V2
	flipy bool // flip the y-axis
}

func NewMap2(bb Box2, grid V2i, flipy bool) (*Map2, error) {
	// sanity check the bounding box
	bb_size := bb.Size()
	if bb_size.X <= 0 || bb_size.Y <= 0 {
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
	m.delta = bb_size.Div(grid.ToV2())
	return &m, nil
}

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
	max_d2 := 0.0
	min_d2 := 0.0

	// translate the box so p is at the origin
	a = a.Translate(p.Neg())

	// consider the vertices
	vs := a.Vertices()

	for i := range vs {
		d2 := vs[i].Length2()
		if i == 0 {
			min_d2 = d2
		} else {
			min_d2 = Min(min_d2, d2)
		}
		max_d2 = Max(max_d2, d2)
	}

	// consider the sides (for the minimum)
	within_x := a.Min.X < 0 && a.Max.X > 0
	within_y := a.Min.Y < 0 && a.Max.Y > 0

	if within_x && within_y {
		min_d2 = 0
	} else {
		if within_x {
			d := Min(Abs(a.Max.Y), Abs(a.Min.Y))
			min_d2 = Min(min_d2, d*d)
		}
		if within_y {
			d := Min(Abs(a.Max.X), Abs(a.Min.X))
			min_d2 = Min(min_d2, d*d)
		}
	}

	return V2{min_d2, max_d2}
}

// MinMaxDist2 returns the minimum and maximum dist * dist from a point to a box.
// Points within the box have minimum distance = 0.
func (a Box3) MinMaxDist2(p V3) V2 {
	max_d2 := 0.0
	min_d2 := 0.0

	// translate the box so p is at the origin
	a = a.Translate(p.Neg())

	// consider the vertices
	vs := a.Vertices()
	for i := range vs {
		d2 := vs[i].Length2()
		if i == 0 {
			min_d2 = d2
		} else {
			min_d2 = Min(min_d2, d2)
		}
		max_d2 = Max(max_d2, d2)
	}

	// consider the faces (for the minimum)
	within_x := a.Min.X < 0 && a.Max.X > 0
	within_y := a.Min.Y < 0 && a.Max.Y > 0
	within_z := a.Min.Z < 0 && a.Max.Z > 0

	if within_x && within_y && within_z {
		min_d2 = 0
	} else {
		if within_x && within_y {
			d := Min(Abs(a.Max.Z), Abs(a.Min.Z))
			min_d2 = Min(min_d2, d*d)
		}
		if within_x && within_z {
			d := Min(Abs(a.Max.Y), Abs(a.Min.Y))
			min_d2 = Min(min_d2, d*d)
		}
		if within_y && within_z {
			d := Min(Abs(a.Max.X), Abs(a.Min.X))
			min_d2 = Min(min_d2, d*d)
		}
	}

	return V2{min_d2, max_d2}
}

//-----------------------------------------------------------------------------
