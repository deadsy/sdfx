//-----------------------------------------------------------------------------
/*

Integer 2D/3D Vectors

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

type V2i [2]int
type V3i [3]int

//-----------------------------------------------------------------------------

// Subtract a scalar from each component of the vector.
func (a V2i) SubScalar(b int) V2i {
	return V2i{a[0] - b, a[1] - b}
}

// Subtract a scalar from each component of the vector.
func (a V3i) SubScalar(b int) V3i {
	return V3i{a[0] - b, a[1] - b, a[2] - b}
}

//-----------------------------------------------------------------------------

// Convert V2i to V2.
func (a V2i) ToV2() V2 {
	return V2{float64(a[0]), float64(a[1])}
}

// Convert V2 to V2i.
func (a V2) ToV2i() V2i {
	return V2i{int(a.X), int(a.Y)}
}

// Convert V3i to V3.
func (a V3i) ToV3() V3 {
	return V3{float64(a[0]), float64(a[1]), float64(a[2])}
}

// Convert V3 to V3i.
func (a V3) ToV3i() V3i {
	return V3i{int(a.X), int(a.Y), int(a.Z)}
}

//-----------------------------------------------------------------------------

// Add two vectors. Return v = a + b.
func (a V3i) Add(b V3i) V3i {
	return V3i{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

//-----------------------------------------------------------------------------
