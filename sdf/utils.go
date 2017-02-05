//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"math"
)

//-----------------------------------------------------------------------------

const PI = math.Pi
const TAU = 2 * math.Pi
const SQRT_HALF = 0.7071067811865476
const EPS = 1e-9

//-----------------------------------------------------------------------------

// Degrees to radians
func DtoR(degrees float64) float64 {
	return (PI / 180) * degrees
}

// Radians to degrees
func RtoD(radians float64) float64 {
	return (180 / PI) * radians
}

//-----------------------------------------------------------------------------

// Clamp x between a and b, assume a <= b
func Clamp(x, a, b float64) float64 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

// Linear Interpolation from x to y, a = [0,1]
func Mix(x, y, a float64) float64 {
	return x + (a * (y - x))
}

//-----------------------------------------------------------------------------
// Max/Min functions
// Note: math.Max/Min don't inline because they do NaN/Inf checking.

// maximum of a and b
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// minimum of a and b
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

// absolute value of x
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}

// sign of x
func Sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

//-----------------------------------------------------------------------------
// Minimum Functions for SDF blending

type MinFunc func(a, b, k float64) float64

// Normal Minimum, no blending.
func NormalMin(a, b, k float64) float64 {
	return Min(a, b)
}

// Round Minimum, uses a quarter-circle to join the two objects smoothly.
func RoundMin(a, b, k float64) float64 {
	u := V2{k - a, k - b}.Max(V2{0, 0})
	return Max(k, Min(a, b)) - u.Length()
}

// Chamfer Minimum, makes a 45-degree chamfered edge (the diagonal of a square of size <r>).
// TODO: why the holes in the rendering?
func ChamferMin(a, b, k float64) float64 {
	return Min(Min(a, b), (a-k+b)*SQRT_HALF)
}

// Exponential Smooth Minimum (k = 32).
func ExpMin(a, b, k float64) float64 {
	return -math.Log(math.Exp(-k*a)+math.Exp(-k*b)) / k
}

// Power Smooth Minimum (k = 8).
// TODO - weird results, is this correct?
func PowMin(a, b, k float64) float64 {
	a = math.Pow(a, k)
	b = math.Pow(b, k)
	return math.Pow((a*b)/(a+b), 1/k)
}

// Polynomial Smooth Minimum (Try k = 0.1, a bigger k gives a bigger fillet).
func PolyMin(a, b, k float64) float64 {
	h := Clamp(0.5+0.5*(b-a)/k, 0.0, 1.0)
	return Mix(b, a, h) - k*h*(1.0-h)
}

//-----------------------------------------------------------------------------
// Maximum Functions for SDF blending

type MaxFunc func(a, b, k float64) float64

// Normal maximum, no blending.
func NormalMax(a, b, k float64) float64 {
	return Max(a, b)
}

// Polynomial Smooth Maximum (Try k = 0.1, a bigger k gives a bigger fillet).
func PolyMax(a, b, k float64) float64 {
	return -PolyMin(-a, -b, k)
}

//-----------------------------------------------------------------------------
