//-----------------------------------------------------------------------------
/*

Floating Point 2D/3D Vectors

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
	"math/rand"
)

//-----------------------------------------------------------------------------

// V3 is a 3d float64 cartesian vector.
type V3 struct {
	X, Y, Z float64
}

// V2 is a 2d float64 cartesian vector.
type V2 struct {
	X, Y float64
}

// P2 is a 2d float64 polar vector.
type P2 struct {
	R, Theta float64
}

// V2Set is a set of 2d float64 vectors.
type V2Set []V2

// V3Set is a set of 3d float64 vectors.
type V3Set []V3

//-----------------------------------------------------------------------------

// Equals returns true if a == b within the tolerance limit.
func (a V3) Equals(b V3, tolerance float64) bool {
	return (Abs(a.X-b.X) < tolerance &&
		Abs(a.Y-b.Y) < tolerance &&
		Abs(a.Z-b.Z) < tolerance)
}

// Equals returns true if a == b within the tolerance limit.
func (a V2) Equals(b V2, tolerance float64) bool {
	return (Abs(a.X-b.X) < tolerance &&
		Abs(a.Y-b.Y) < tolerance)
}

//-----------------------------------------------------------------------------

// randomRange returns a random float64 [a,b)
func randomRange(a, b float64) float64 {
	return a + (b-a)*rand.Float64()
}

// Random returns a random point within a bounding box.
func (b *Box2) Random() V2 {
	return V2{
		randomRange(b.Min.X, b.Max.X),
		randomRange(b.Min.Y, b.Max.Y),
	}
}

// Random returns a random point within a bounding box.
func (b *Box3) Random() V3 {
	return V3{
		randomRange(b.Min.X, b.Max.X),
		randomRange(b.Min.Y, b.Max.Y),
		randomRange(b.Min.Z, b.Max.Z),
	}
}

// RandomSet returns a set of random points from within a bounding box.
func (b *Box2) RandomSet(n int) V2Set {
	s := make([]V2, n)
	for i := range s {
		s[i] = b.Random()
	}
	return s
}

// RandomSet returns a set of random points from within a bounding box.
func (b *Box3) RandomSet(n int) V3Set {
	s := make([]V3, n)
	for i := range s {
		s[i] = b.Random()
	}
	return s
}

//-----------------------------------------------------------------------------

// Dot returns the dot product of a and b.
func (a V3) Dot(b V3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Dot returns the dot product of a and b.
func (a V2) Dot(b V2) float64 {
	return a.X*b.X + a.Y*b.Y
}

// Cross returns the cross product of a and b.
func (a V3) Cross(b V3) V3 {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return V3{x, y, z}
}

// Cross returns the cross product of a and b.
func (a V2) Cross(b V2) float64 {
	return (a.X * b.Y) - (a.Y * b.X)
}

// colinearSlow return true if 3 points are colinear (slow test).
func colinearSlow(a, b, c V2, tolerance float64) bool {
	// use the cross product as a measure of colinearity
	pa := a.Sub(c).Normalize()
	pb := b.Sub(c).Normalize()
	return Abs(pa.Cross(pb)) < tolerance
}

// colinearFast return true if 3 points are colinear (fast test).
func colinearFast(a, b, c V2, tolerance float64) bool {
	// use the cross product as a measure of colinearity
	ac := a.Sub(b)
	bc := b.Sub(c)
	return Abs(ac.Cross(bc)) < tolerance
}

//-----------------------------------------------------------------------------

// AddScalar adds a scalar to each vector component.
func (a V3) AddScalar(b float64) V3 {
	return V3{a.X + b, a.Y + b, a.Z + b}
}

// AddScalar adds a scalar to each vector component.
func (a V2) AddScalar(b float64) V2 {
	return V2{a.X + b, a.Y + b}
}

// SubScalar subtracts a scalar from each vector component.
func (a V3) SubScalar(b float64) V3 {
	return V3{a.X - b, a.Y - b, a.Z - b}
}

// SubScalar subtracts a scalar from each vector component.
func (a V2) SubScalar(b float64) V2 {
	return V2{a.X - b, a.Y - b}
}

// MulScalar multiplies each vector component by a scalar.
func (a V3) MulScalar(b float64) V3 {
	return V3{a.X * b, a.Y * b, a.Z * b}
}

// MulScalar multiplies each vector component by a scalar.
func (a V2) MulScalar(b float64) V2 {
	return V2{a.X * b, a.Y * b}
}

// DivScalar divides each vector component by a scalar.
func (a V3) DivScalar(b float64) V3 {
	return V3{a.X / b, a.Y / b, a.Z / b}
}

// DivScalar divides each vector component by a scalar.
func (a V2) DivScalar(b float64) V2 {
	return V2{a.X / b, a.Y / b}
}

//-----------------------------------------------------------------------------

// Negate negates each vector component.
func (a V3) Negate() V3 {
	return V3{-a.X, -a.Y, -a.Z}
}

// Negate negates each vector component.
func (a V2) Negate() V2 {
	return V2{-a.X, -a.Y}
}

// Abs takes the absolute value of each vector component.
func (a V3) Abs() V3 {
	return V3{Abs(a.X), Abs(a.Y), Abs(a.Z)}
}

// Abs takes the absolute value of each vector component.
func (a V2) Abs() V2 {
	return V2{Abs(a.X), Abs(a.Y)}
}

// Ceil takes the ceiling value of each vector component.
func (a V3) Ceil() V3 {
	return V3{math.Ceil(a.X), math.Ceil(a.Y), math.Ceil(a.Z)}
}

// Ceil takes the ceiling value of each vector component.
func (a V2) Ceil() V2 {
	return V2{math.Ceil(a.X), math.Ceil(a.Y)}
}

//-----------------------------------------------------------------------------

// Min return a vector with the minimum components of two vectors.
func (a V3) Min(b V3) V3 {
	return V3{Min(a.X, b.X), Min(a.Y, b.Y), Min(a.Z, b.Z)}
}

// Min return a vector with the minimum components of two vectors.
func (a V2) Min(b V2) V2 {
	return V2{Min(a.X, b.X), Min(a.Y, b.Y)}
}

// Max return a vector with the maximum components of two vectors.
func (a V3) Max(b V3) V3 {
	return V3{Max(a.X, b.X), Max(a.Y, b.Y), Max(a.Z, b.Z)}
}

// Max return a vector with the maximum components of two vectors.
func (a V2) Max(b V2) V2 {
	return V2{Max(a.X, b.X), Max(a.Y, b.Y)}
}

// Add adds two vectors. Returns a + b.
func (a V3) Add(b V3) V3 {
	return V3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Add adds two vectors. Returns a + b.
func (a V2) Add(b V2) V2 {
	return V2{a.X + b.X, a.Y + b.Y}
}

// Sub subtracts two vectors. Returns a - b.
func (a V3) Sub(b V3) V3 {
	return V3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Sub subtracts two vectors. Returns a - b.
func (a V2) Sub(b V2) V2 {
	return V2{a.X - b.X, a.Y - b.Y}
}

// Mul multiplies two vectors by component.
func (a V3) Mul(b V3) V3 {
	return V3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Mul multiplies two vectors by component.
func (a V2) Mul(b V2) V2 {
	return V2{a.X * b.X, a.Y * b.Y}
}

// Div divides two vectors by component.
func (a V3) Div(b V3) V3 {
	return V3{a.X / b.X, a.Y / b.Y, a.Z / b.Z}
}

// Div divides two vectors by component.
func (a V2) Div(b V2) V2 {
	return V2{a.X / b.X, a.Y / b.Y}
}

// Neg negates a vector.
func (a V2) Neg() V2 {
	return V2{-a.X, -a.Y}
}

// Neg negates a vector.
func (a V3) Neg() V3 {
	return V3{-a.X, -a.Y, -a.Z}
}

//-----------------------------------------------------------------------------

// Min return the minimum components of a set of vectors.
func (a V3Set) Min() V3 {
	vmin := a[0]
	for _, v := range a {
		vmin = vmin.Min(v)
	}
	return vmin
}

// Min return the minimum components of a set of vectors.
func (a V2Set) Min() V2 {
	vmin := a[0]
	for _, v := range a {
		vmin = vmin.Min(v)
	}
	return vmin
}

// Max return the maximum components of a set of vectors.
func (a V3Set) Max() V3 {
	vmax := a[0]
	for _, v := range a {
		vmax = vmax.Max(v)
	}
	return vmax
}

// Max return the maximum components of a set of vectors.
func (a V2Set) Max() V2 {
	vmax := a[0]
	for _, v := range a {
		vmax = vmax.Max(v)
	}
	return vmax
}

//-----------------------------------------------------------------------------

// Length returns the vector length.
func (a V3) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

// Length2 returns the vector length * length.
func (a V3) Length2() float64 {
	return a.X*a.X + a.Y*a.Y + a.Z*a.Z
}

// Length returns the vector length.
func (a V2) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

// Length2 returns the vector length * length.
func (a V2) Length2() float64 {
	return a.X*a.X + a.Y*a.Y
}

// MinComponent returns the minimum component of the vector.
func (a V3) MinComponent() float64 {
	return Min(Min(a.X, a.Y), a.Z)
}

// MinComponent returns the minimum component of the vector.
func (a V2) MinComponent() float64 {
	return Min(a.X, a.Y)
}

// MaxComponent returns the maximum component of the vector.
func (a V3) MaxComponent() float64 {
	return Max(Max(a.X, a.Y), a.Z)
}

// MaxComponent returns the maximum component of the vector.
func (a V2) MaxComponent() float64 {
	return Max(a.X, a.Y)
}

//-----------------------------------------------------------------------------

// Normalize scales a vector to unit length.
func (a V3) Normalize() V3 {
	d := a.Length()
	return V3{a.X / d, a.Y / d, a.Z / d}
}

// Normalize scales a vector to unit length.
func (a V2) Normalize() V2 {
	d := a.Length()
	return V2{a.X / d, a.Y / d}
}

//-----------------------------------------------------------------------------

// ToV3 converts a V2 to a V3 with a specified Z value.
func (a V2) ToV3(z float64) V3 {
	return V3{a.X, a.Y, z}
}

//-----------------------------------------------------------------------------

// Overlap returns true if 1D line segments a and b overlap.
func (a V2) Overlap(b V2) bool {
	return a.Y >= b.X && b.Y >= a.X
}

//-----------------------------------------------------------------------------

// PolarToCartesian converts a polar to a cartesian coordinate.
func (a P2) PolarToCartesian() V2 {
	return V2{a.R * math.Cos(a.Theta), a.R * math.Sin(a.Theta)}
}

// CartesianToPolar converts a cartesian to a polar coordinate.
func (a V2) CartesianToPolar() P2 {
	return P2{a.Length(), math.Atan2(a.Y, a.X)}
}

// PolarToXY converts polar to cartesian coordinates. (TODO remove)
func PolarToXY(r, theta float64) V2 {
	return P2{r, theta}.PolarToCartesian()
}

//-----------------------------------------------------------------------------
// Sort By X for a V2Set

// V2SetByX used to sort V2Set by X-value.
type V2SetByX V2Set

func (a V2SetByX) Len() int           { return len(a) }
func (a V2SetByX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a V2SetByX) Less(i, j int) bool { return a[i].X < a[j].X }

//-----------------------------------------------------------------------------
