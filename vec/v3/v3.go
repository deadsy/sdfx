//-----------------------------------------------------------------------------
/*

Floating Point 3D Vectors

*/
//-----------------------------------------------------------------------------

package v3

import "math"

//-----------------------------------------------------------------------------

// clamp x between a and b, assume a <= b
func clamp(x, a, b float64) float64 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

//-----------------------------------------------------------------------------

// Vec is a 3D float64 vector.
type Vec struct {
	X, Y, Z float64
}

//-----------------------------------------------------------------------------

// Equals returns true if a == b within the tolerance limit.
func (a Vec) Equals(b Vec, tolerance float64) bool {
	return (math.Abs(a.X-b.X) <= tolerance &&
		math.Abs(a.Y-b.Y) <= tolerance &&
		math.Abs(a.Z-b.Z) <= tolerance)
}

// LTZero returns true if any vector components are < 0.
func (a Vec) LTZero() bool {
	return (a.X < 0) || (a.Y < 0) || (a.Z < 0)
}

// LTEZero returns true if any vector components are <= 0.
func (a Vec) LTEZero() bool {
	return (a.X <= 0) || (a.Y <= 0) || (a.Z <= 0)
}

//-----------------------------------------------------------------------------

// Dot returns the dot product of a and b.
func (a Vec) Dot(b Vec) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross returns the cross product of a and b.
func (a Vec) Cross(b Vec) Vec {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return Vec{x, y, z}
}

// AddScalar adds a scalar to each vector component.
func (a Vec) AddScalar(b float64) Vec {
	return Vec{a.X + b, a.Y + b, a.Z + b}
}

// SubScalar subtracts a scalar from each vector component.
func (a Vec) SubScalar(b float64) Vec {
	return Vec{a.X - b, a.Y - b, a.Z - b}
}

// MulScalar multiplies each vector component by a scalar.
func (a Vec) MulScalar(b float64) Vec {
	return Vec{a.X * b, a.Y * b, a.Z * b}
}

// DivScalar divides each vector component by a scalar.
func (a Vec) DivScalar(b float64) Vec {
	return a.MulScalar(1 / b)
}

// Abs takes the absolute value of each vector component.
func (a Vec) Abs() Vec {
	return Vec{math.Abs(a.X), math.Abs(a.Y), math.Abs(a.Z)}
}

// Ceil takes the ceiling value of each vector component.
func (a Vec) Ceil() Vec {
	return Vec{math.Ceil(a.X), math.Ceil(a.Y), math.Ceil(a.Z)}
}

// Clamp clamps a vector between 2 other vectors.
func (a Vec) Clamp(b, c Vec) Vec {
	return Vec{clamp(a.X, b.X, c.X), clamp(a.Y, b.Y, c.Y), clamp(a.Z, b.Z, c.Z)}
}

// Min return a vector with the minimum components of two vectors.
func (a Vec) Min(b Vec) Vec {
	return Vec{math.Min(a.X, b.X), math.Min(a.Y, b.Y), math.Min(a.Z, b.Z)}
}

// Max return a vector with the maximum components of two vectors.
func (a Vec) Max(b Vec) Vec {
	return Vec{math.Max(a.X, b.X), math.Max(a.Y, b.Y), math.Max(a.Z, b.Z)}
}

// Add adds two vectors. Returns a + b.
func (a Vec) Add(b Vec) Vec {
	return Vec{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Sub subtracts two vectors. Returns a - b.
func (a Vec) Sub(b Vec) Vec {
	return Vec{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Mul multiplies two vectors by component.
func (a Vec) Mul(b Vec) Vec {
	return Vec{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Div divides two vectors by component.
func (a Vec) Div(b Vec) Vec {
	return Vec{a.X / b.X, a.Y / b.Y, a.Z / b.Z}
}

// Neg negates a vector.
func (a Vec) Neg() Vec {
	return Vec{-a.X, -a.Y, -a.Z}
}

// Length returns the vector length.
func (a Vec) Length() float64 {
	return math.Sqrt(a.Length2())
}

// Length2 returns the vector length * length.
func (a Vec) Length2() float64 {
	return a.Dot(a)
}

// Normalize scales a vector to unit length.
func (a Vec) Normalize() Vec {
	return a.MulScalar(1 / a.Length())
}

// MinComponent returns the minimum component of the vector.
func (a Vec) MinComponent() float64 {
	return math.Min(math.Min(a.X, a.Y), a.Z)
}

// MaxComponent returns the maximum component of the vector.
func (a Vec) MaxComponent() float64 {
	return math.Max(math.Max(a.X, a.Y), a.Z)
}

// Sin takes the sine of each vector component.
func (a Vec) Sin() Vec {
	return Vec{math.Sin(a.X), math.Sin(a.Y), math.Sin(a.Z)}
}

// Cos takes the cosine of each vector component.
func (a Vec) Cos() Vec {
	return Vec{math.Cos(a.X), math.Cos(a.Y), math.Cos(a.Z)}
}

//-----------------------------------------------------------------------------

// Get the n-th component of the vector.
func (a Vec) Get(i int) float64 {
	switch i {
	case 0:
		return a.X
	case 1:
		return a.Y
	case 2:
		return a.Z
	}
	panic("bad vector component")
}

// Set the n-th component of the vector.
func (a *Vec) Set(i int, val float64) {
	switch i {
	case 0:
		a.X = val
	case 1:
		a.Y = val
	case 2:
		a.Z = val
	default:
		panic("bad vector component")
	}
}

//-----------------------------------------------------------------------------

// VecSet is a set of 3D float64 vectors.
type VecSet []Vec

// Min return the minimum components of a set of vectors.
func (a VecSet) Min() Vec {
	vmin := a[0]
	for _, v := range a {
		vmin = vmin.Min(v)
	}
	return vmin
}

// Max return the maximum components of a set of vectors.
func (a VecSet) Max() Vec {
	vmax := a[0]
	for _, v := range a {
		vmax = vmax.Max(v)
	}
	return vmax
}

//-----------------------------------------------------------------------------
