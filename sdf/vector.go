//-----------------------------------------------------------------------------
/*

3D/2D Vector Operations

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
	"math/rand"
)

//-----------------------------------------------------------------------------

type V3 struct {
	X, Y, Z float64
}
type V2 struct {
	X, Y float64
}

//-----------------------------------------------------------------------------

func (a V3) Equals(b V3, tolerance float64) bool {
	return (Abs(a.X-b.X) < tolerance &&
		Abs(a.Y-b.Y) < tolerance &&
		Abs(a.Z-b.Z) < tolerance)
}

func (a V2) Equals(b V2, tolerance float64) bool {
	return (Abs(a.X-b.X) < tolerance &&
		Abs(a.Y-b.Y) < tolerance)
}

//-----------------------------------------------------------------------------

func RandomUnitV3(rnd *rand.Rand) V3 {
	for {
		var x, y, z float64
		if rnd == nil {
			x = rand.Float64()*2 - 1
			y = rand.Float64()*2 - 1
			z = rand.Float64()*2 - 1
		} else {
			x = rnd.Float64()*2 - 1
			y = rnd.Float64()*2 - 1
			z = rnd.Float64()*2 - 1
		}
		if x*x+y*y+z*z > 1 {
			continue
		}
		return V3{x, y, z}.Normalize()
	}
}

func RandomUnitV2(rnd *rand.Rand) V2 {
	for {
		var x, y float64
		if rnd == nil {
			x = rand.Float64()*2 - 1
			y = rand.Float64()*2 - 1
		} else {
			x = rnd.Float64()*2 - 1
			y = rnd.Float64()*2 - 1
		}
		if x*x+y*y > 1 {
			continue
		}
		return V2{x, y}.Normalize()
	}
}

func random_range(a, b float64) float64 {
	return a + (b-a)*rand.Float64()
}

func RandomV3(a, b float64) V3 {
	return V3{random_range(a, b),
		random_range(a, b),
		random_range(a, b)}
}

func RandomV2(a, b float64) V2 {
	return V2{random_range(a, b),
		random_range(a, b)}
}

//-----------------------------------------------------------------------------

func (a V3) Dot(b V3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a V2) Dot(b V2) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a V3) Cross(b V3) V3 {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return V3{x, y, z}
}

//-----------------------------------------------------------------------------

func (a V3) AddScalar(b float64) V3 {
	return V3{a.X + b, a.Y + b, a.Z + b}
}
func (a V2) AddScalar(b float64) V2 {
	return V2{a.X + b, a.Y + b}
}

func (a V3) SubScalar(b float64) V3 {
	return V3{a.X - b, a.Y - b, a.Z - b}
}
func (a V2) SubScalar(b float64) V2 {
	return V2{a.X - b, a.Y - b}
}

func (a V3) MulScalar(b float64) V3 {
	return V3{a.X * b, a.Y * b, a.Z * b}
}
func (a V2) MulScalar(b float64) V2 {
	return V2{a.X * b, a.Y * b}
}

func (a V3) DivScalar(b float64) V3 {
	return V3{a.X / b, a.Y / b, a.Z / b}
}
func (a V2) DivScalar(b float64) V2 {
	return V2{a.X / b, a.Y / b}
}

//-----------------------------------------------------------------------------

func (a V3) Negate() V3 {
	return V3{-a.X, -a.Y, -a.Z}
}
func (a V2) Negate() V2 {
	return V2{-a.X, -a.Y}
}

func (a V3) Abs() V3 {
	return V3{Abs(a.X), Abs(a.Y), Abs(a.Z)}
}
func (a V2) Abs() V2 {
	return V2{Abs(a.X), Abs(a.Y)}
}

//-----------------------------------------------------------------------------

func (a V3) Min(b V3) V3 {
	return V3{Min(a.X, b.X), Min(a.Y, b.Y), Min(a.Z, b.Z)}
}
func (a V2) Min(b V2) V2 {
	return V2{Min(a.X, b.X), Min(a.Y, b.Y)}
}

func (a V3) Max(b V3) V3 {
	return V3{Max(a.X, b.X), Max(a.Y, b.Y), Max(a.Z, b.Z)}
}
func (a V2) Max(b V2) V2 {
	return V2{Max(a.X, b.X), Max(a.Y, b.Y)}
}

func (a V3) Add(b V3) V3 {
	return V3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}
func (a V2) Add(b V2) V2 {
	return V2{a.X + b.X, a.Y + b.Y}
}

func (a V3) Sub(b V3) V3 {
	return V3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}
func (a V2) Sub(b V2) V2 {
	return V2{a.X - b.X, a.Y - b.Y}
}

func (a V3) Mul(b V3) V3 {
	return V3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}
func (a V2) Mul(b V2) V2 {
	return V2{a.X * b.X, a.Y * b.Y}
}

func (a V3) Div(b V3) V3 {
	return V3{a.X / b.X, a.Y / b.Y, a.Z / b.Z}
}
func (a V2) Div(b V2) V2 {
	return V2{a.X / b.X, a.Y / b.Y}
}

//-----------------------------------------------------------------------------

func (a V3) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}
func (a V2) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

func (a V3) MinComponent() float64 {
	return Min(Min(a.X, a.Y), a.Z)
}
func (a V2) MinComponent() float64 {
	return Min(a.X, a.Y)
}

func (a V3) MaxComponent() float64 {
	return Max(Max(a.X, a.Y), a.Z)
}
func (a V2) MaxComponent() float64 {
	return Max(a.X, a.Y)
}

//-----------------------------------------------------------------------------

func (a V3) Normalize() V3 {
	d := a.Length()
	return V3{a.X / d, a.Y / d, a.Z / d}
}

func (a V2) Normalize() V2 {
	d := a.Length()
	return V2{a.X / d, a.Y / d}
}

//-----------------------------------------------------------------------------
