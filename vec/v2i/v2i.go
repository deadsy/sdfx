//-----------------------------------------------------------------------------
/*

Integer 2D Vectors

*/
//-----------------------------------------------------------------------------

package v2i

//-----------------------------------------------------------------------------

// Vec is a 2D integer vector.
type Vec struct {
	X, Y int
}

//-----------------------------------------------------------------------------

// AddScalar adds a scalar to each component of the vector.
func (a Vec) AddScalar(b int) Vec {
	return Vec{a.X + b, a.Y + b}
}

// SubScalar subtracts a scalar from each component of the vector.
func (a Vec) SubScalar(b int) Vec {
	return Vec{a.X - b, a.Y - b}
}

// Add adds two vectors. Return v = a + b.
func (a Vec) Add(b Vec) Vec {
	return Vec{a.X + b.X, a.Y + b.Y}
}

//-----------------------------------------------------------------------------
