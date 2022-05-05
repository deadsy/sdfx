//-----------------------------------------------------------------------------
/*

Floating Point 2D Vectors

*/
//-----------------------------------------------------------------------------

package v2

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

// Vec is a 2D float64 vector.
type Vec struct {
	X, Y float64
}

//-----------------------------------------------------------------------------

// Equals returns true if a == b within the tolerance limit.
func (a Vec) Equals(b Vec, tolerance float64) bool {
	return (math.Abs(a.X-b.X) <= tolerance &&
		math.Abs(a.Y-b.Y) <= tolerance)
}

// LTZero returns true if any vector components are < 0.
func (a Vec) LTZero() bool {
	return (a.X < 0) || (a.Y < 0)
}

// LTEZero returns true if any vector components are < 0.
func (a Vec) LTEZero() bool {
	return (a.X <= 0) || (a.Y <= 0)
}

//-----------------------------------------------------------------------------

// Dot returns the dot product of a and b.
func (a Vec) Dot(b Vec) float64 {
	return a.X*b.X + a.Y*b.Y
}

// Cross returns the cross product of a and b.
func (a Vec) Cross(b Vec) float64 {
	return (a.X * b.Y) - (a.Y * b.X)
}

// AddScalar adds a scalar to each vector component.
func (a Vec) AddScalar(b float64) Vec {
	return Vec{a.X + b, a.Y + b}
}

// SubScalar subtracts a scalar from each vector component.
func (a Vec) SubScalar(b float64) Vec {
	return Vec{a.X - b, a.Y - b}
}

// MulScalar multiplies each vector component by a scalar.
func (a Vec) MulScalar(b float64) Vec {
	return Vec{a.X * b, a.Y * b}
}

// DivScalar divides each vector component by a scalar.
func (a Vec) DivScalar(b float64) Vec {
	return a.MulScalar(1 / b)
}

// Abs takes the absolute value of each vector component.
func (a Vec) Abs() Vec {
	return Vec{math.Abs(a.X), math.Abs(a.Y)}
}

// Ceil takes the ceiling value of each vector component.
func (a Vec) Ceil() Vec {
	return Vec{math.Ceil(a.X), math.Ceil(a.Y)}
}

// Clamp clamps a vector between 2 other vectors.
func (a Vec) Clamp(b, c Vec) Vec {
	return Vec{clamp(a.X, b.X, c.X), clamp(a.Y, b.Y, c.Y)}
}

// Min return a vector with the minimum components of two vectors.
func (a Vec) Min(b Vec) Vec {
	return Vec{math.Min(a.X, b.X), math.Min(a.Y, b.Y)}
}

// Max return a vector with the maximum components of two vectors.
func (a Vec) Max(b Vec) Vec {
	return Vec{math.Max(a.X, b.X), math.Max(a.Y, b.Y)}
}

// Add adds two vectors. Returns a + b.
func (a Vec) Add(b Vec) Vec {
	return Vec{a.X + b.X, a.Y + b.Y}
}

// Sub subtracts two vectors. Returns a - b.
func (a Vec) Sub(b Vec) Vec {
	return Vec{a.X - b.X, a.Y - b.Y}
}

// Mul multiplies two vectors by component.
func (a Vec) Mul(b Vec) Vec {
	return Vec{a.X * b.X, a.Y * b.Y}
}

// Div divides two vectors by component.
func (a Vec) Div(b Vec) Vec {
	return Vec{a.X / b.X, a.Y / b.Y}
}

// Neg negates a vector.
func (a Vec) Neg() Vec {
	return Vec{-a.X, -a.Y}
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
	return math.Min(a.X, a.Y)
}

// MaxComponent returns the maximum component of the vector.
func (a Vec) MaxComponent() float64 {
	return math.Max(a.X, a.Y)
}

// Overlap returns true if 1D line segments a and b overlap.
func (a Vec) Overlap(b Vec) bool {
	return a.Y >= b.X && b.Y >= a.X
}

//-----------------------------------------------------------------------------

// VecSet is a set of 2D float64 vectors.
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

// VecSetByX sorts the VecSet by X value
type VecSetByX VecSet

func (a VecSetByX) Len() int           { return len(a) }
func (a VecSetByX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VecSetByX) Less(i, j int) bool { return a[i].X < a[j].X }

//-----------------------------------------------------------------------------
