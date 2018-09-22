//-----------------------------------------------------------------------------
/*

Integer 2D/3D Vectors

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

// V2i is a 2D integer vector.
type V2i [2]int

// V3i is a 3D integer vector.
type V3i [3]int

//-----------------------------------------------------------------------------

// SubScalar subtracts a scalar from each component of the vector.
func (a V2i) SubScalar(b int) V2i {
	return V2i{a[0] - b, a[1] - b}
}

// SubScalar subtracts a scalar from each component of the vector.
func (a V3i) SubScalar(b int) V3i {
	return V3i{a[0] - b, a[1] - b, a[2] - b}
}

// AddScalar adds a scalar to each component of the vector.
func (a V2i) AddScalar(b int) V2i {
	return V2i{a[0] + b, a[1] + b}
}

// AddScalar adds a scalar to each component of the vector.
func (a V3i) AddScalar(b int) V3i {
	return V3i{a[0] + b, a[1] + b, a[2] + b}
}

//-----------------------------------------------------------------------------

// ToV2 converts V2i (integer) to V2 (float).
func (a V2i) ToV2() V2 {
	return V2{float64(a[0]), float64(a[1])}
}

// ToV2i convert V2 (float) to V2i (integer).
func (a V2) ToV2i() V2i {
	return V2i{int(a.X), int(a.Y)}
}

// ToV3 converts V3i (integer) to V3 (float).
func (a V3i) ToV3() V3 {
	return V3{float64(a[0]), float64(a[1]), float64(a[2])}
}

// ToV3i convert V3 (float) to V3i (integer).
func (a V3) ToV3i() V3i {
	return V3i{int(a.X), int(a.Y), int(a.Z)}
}

//-----------------------------------------------------------------------------

// Add adds two vectors. Return v = a + b.
func (a V2i) Add(b V2i) V2i {
	return V2i{a[0] + b[0], a[1] + b[1]}
}

// Add adds two vectors. Return v = a + b.
func (a V3i) Add(b V3i) V3i {
	return V3i{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

//-----------------------------------------------------------------------------
