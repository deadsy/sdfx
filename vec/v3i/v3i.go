//-----------------------------------------------------------------------------
/*

Integer 3D Vectors

*/
//-----------------------------------------------------------------------------

package v3i

//-----------------------------------------------------------------------------

// Vec is a 3D integer vector.
type Vec struct {
	X, Y, Z int
}

//-----------------------------------------------------------------------------

// AddScalar adds a scalar to each component of the vector.
func (a Vec) AddScalar(b int) Vec {
	return Vec{a.X + b, a.Y + b, a.Z + b}
}

// SubScalar subtracts a scalar from each component of the vector.
func (a Vec) SubScalar(b int) Vec {
	return Vec{a.X - b, a.Y - b, a.Z - b}
}

// Add adds two vectors. Return v = a + b.
func (a Vec) Add(b Vec) Vec {
	return Vec{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

//-----------------------------------------------------------------------------
