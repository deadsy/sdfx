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

type V3 [3]float64
type V2 [2]float64

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

func RandomV3(a, b V3) V3 {
	x := random_range(a[0], b[0])
	y := random_range(a[0], b[0])
	z := random_range(a[0], b[0])
	return V3{x, y, z}
}

func RandomV2(a, b V2) V2 {
	x := random_range(a[0], b[0])
	y := random_range(a[0], b[0])
	return V2{x, y}
}

//-----------------------------------------------------------------------------

func (a V3) AddScalar(b float64) V3 {
	return V3{a[0] + b, a[1] + b, a[2] + b}
}
func (a V2) AddScalar(b float64) V2 {
	return V2{a[0] + b, a[1] + b}
}

func (a V3) SubScalar(b float64) V3 {
	return V3{a[0] - b, a[1] - b, a[2] - b}
}
func (a V2) SubScalar(b float64) V2 {
	return V2{a[0] - b, a[1] - b}
}

func (a V3) MulScalar(b float64) V3 {
	return V3{a[0] * b, a[1] * b, a[2] * b}
}
func (a V2) MulScalar(b float64) V2 {
	return V2{a[0] * b, a[1] * b}
}

func (a V3) DivScalar(b float64) V3 {
	return V3{a[0] / b, a[1] / b, a[2] / b}
}
func (a V2) DivScalar(b float64) V2 {
	return V2{a[0] / b, a[1] / b}
}

//-----------------------------------------------------------------------------

func (a V3) Negate() V3 {
	return V3{-a[0], -a[1], -a[2]}
}
func (a V2) Negate() V2 {
	return V2{-a[0], -a[1]}
}

func (a V3) Abs() V3 {
	return V3{math.Abs(a[0]), math.Abs(a[1]), math.Abs(a[2])}
}
func (a V2) Abs() V2 {
	return V2{math.Abs(a[0]), math.Abs(a[1])}
}

//-----------------------------------------------------------------------------

func (a V3) Min(b V3) V3 {
	return V3{math.Min(a[0], b[0]), math.Min(a[1], b[1]), math.Min(a[2], b[2])}
}
func (a V2) Min(b V2) V2 {
	return V2{math.Min(a[0], b[0]), math.Min(a[1], b[1])}
}

func (a V3) Max(b V3) V3 {
	return V3{math.Max(a[0], b[0]), math.Max(a[1], b[1]), math.Max(a[2], b[2])}
}
func (a V2) Max(b V2) V2 {
	return V2{math.Max(a[0], b[0]), math.Max(a[1], b[1])}
}

func (a V3) Add(b V3) V3 {
	return V3{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}
func (a V2) Add(b V2) V2 {
	return V2{a[0] + b[0], a[1] + b[1]}
}

func (a V3) Sub(b V3) V3 {
	return V3{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}
func (a V2) Sub(b V2) V2 {
	return V2{a[0] - b[0], a[1] - b[1]}
}

func (a V3) Mul(b V3) V3 {
	return V3{a[0] * b[0], a[1] * b[1], a[2] * b[2]}
}
func (a V2) Mul(b V2) V2 {
	return V2{a[0] * b[0], a[1] * b[1]}
}

func (a V3) Div(b V3) V3 {
	return V3{a[0] / b[0], a[1] / b[1], a[2] / b[2]}
}
func (a V2) Div(b V2) V2 {
	return V2{a[0] / b[0], a[1] / b[1]}
}

//-----------------------------------------------------------------------------

func (a V3) Length() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1] + a[2]*a[2])
}
func (a V2) Length() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1])
}

func (a V3) MinComponent() float64 {
	return math.Min(math.Min(a[0], a[1]), a[2])
}
func (a V2) MinComponent() float64 {
	return math.Min(a[0], a[1])
}

func (a V3) MaxComponent() float64 {
	return math.Max(math.Max(a[0], a[1]), a[2])
}
func (a V2) MaxComponent() float64 {
	return math.Max(a[0], a[1])
}

//-----------------------------------------------------------------------------

func (a V3) Normalize() V3 {
	d := a.Length()
	return V3{a[0] / d, a[1] / d, a[2] / d}
}

func (a V2) Normalize() V2 {
	d := a.Length()
	return V2{a[0] / d, a[1] / d}
}

//-----------------------------------------------------------------------------
