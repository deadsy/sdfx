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
