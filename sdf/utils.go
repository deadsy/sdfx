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
// These RonCo-style versions will inline.

func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

// Return abs(x)
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}

//-----------------------------------------------------------------------------
