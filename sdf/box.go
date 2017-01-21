//-----------------------------------------------------------------------------
/*

 */
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

type Box3 struct {
	Min, Max V3
}

type Box2 struct {
	Min, Max V2
}

//-----------------------------------------------------------------------------

func (a Box3) Equals(b Box3, tolerance float64) bool {
	return (a.Min.Equals(b.Min, tolerance) && a.Max.Equals(b.Max, tolerance))
}

func (a Box2) Equals(b Box2, tolerance float64) bool {
	return (a.Min.Equals(b.Min, tolerance) && a.Max.Equals(b.Max, tolerance))
}

//-----------------------------------------------------------------------------
// Extend - Return a box that encloses two boxes

func (a Box3) Extend(b Box3) Box3 {
	return Box3{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

func (a Box2) Extend(b Box2) Box2 {
	return Box2{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

//-----------------------------------------------------------------------------
// Size - Return the size of the box

func (a Box3) Size() V3 {
	return a.Max.Sub(a.Min)
}

func (a Box2) Size() V2 {
	return a.Max.Sub(a.Min)
}

//-----------------------------------------------------------------------------

func (a Box3) Anchor(anchor V3) V3 {
	return a.Min.Add(a.Size().Mul(anchor))
}

func (a Box3) Center() V3 {
	return a.Anchor(V3{0.5, 0.5, 0.5})
}

//-----------------------------------------------------------------------------

func (a Box3) Scale(k float64) Box3 {
	k = k / 2
	c := a.Center()
	s := a.Size().Mul(V3{k, k, k})
	return Box3{c.Sub(s), c.Add(s)}
}

//-----------------------------------------------------------------------------
