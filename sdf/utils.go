//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"math"
)

//-----------------------------------------------------------------------------

const PI = math.Pi
const TAU = 2 * math.Pi

//-----------------------------------------------------------------------------

// Degrees to radians
func DtoR(degrees float64) float64 {
	return math.Pi * degrees / 180
}

// Radians to degrees
func RtoD(radians float64) float64 {
	return 180 * radians / math.Pi
}

//-----------------------------------------------------------------------------

// Clamp value between a and b
func Clamp(x, a, b float64) float64 {
	return math.Min(math.Max(x, a), b)
}

// Linear Interpolation
func Mix(x, y, a float64) float64 {
	return (x * (1 - a)) + (y * a)
}

//-----------------------------------------------------------------------------
