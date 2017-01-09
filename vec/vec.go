//-----------------------------------------------------------------------------
/*
Vector Math
*/
//-----------------------------------------------------------------------------

package vec

import (
	"math"
)

//-----------------------------------------------------------------------------

type V3 [3]float64
type V2 [2]float64

// Return the Euclidean length of a
func (a V3) Length() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1] + a[2]*a[2])
}

func (a V2) Length() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1])
}

// Return a * k
func (a V3) Scale(k float64) V3 {
	return V3{
		a[0] * k,
		a[1] * k,
		a[2] * k,
	}
}

func (a V2) Scale(k float64) V2 {
	return V2{
		a[0] * k,
		a[1] * k,
	}
}

// Return a + b
func (a V3) Sum(b V3) V3 {
	return V3{
		a[0] + b[0],
		a[1] + b[1],
		a[2] + b[2],
	}
}

func (a V2) Sum(b V2) V2 {
	return V2{
		a[0] + b[0],
		a[1] + b[1],
	}
}

// Return a - b
func (a V3) Sub(b V3) V3 {
	return V3{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
	}
}

func (a V2) Sub(b V2) V2 {
	return V2{
		a[0] - b[0],
		a[1] - b[1],
	}
}

// Return a x b
func (a V3) Cross(b V3) V3 {
	return V3{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

// Return a.b
func (a V3) Dot(b V3) float64 {
	return (a[0] * b[0]) +
		(a[1] * b[1]) +
		(a[2] * b[2])
}

func (a V2) Dot(b V2) float64 {
	return (a[0] * b[0]) +
		(a[1] * b[1])
}

// Normalize a
func (a V3) Normalize() V3 {
	l := a.Length()
	if l == 0 {
		return a
	} else {
		return V3{
			a[0] / l,
			a[1] / l,
			a[2] / l,
		}
	}
}

func (a V2) Normalize() V2 {
	l := a.Length()
	if l == 0 {
		return a
	} else {
		return V2{
			a[0] / l,
			a[1] / l,
		}
	}
}

// Return Absolute value of a
func (a V3) Abs() V3 {
	return V3{
		math.Abs(a[0]),
		math.Abs(a[1]),
		math.Abs(a[2]),
	}
}

func (a V2) Abs() V2 {
	return V2{
		math.Abs(a[0]),
		math.Abs(a[1]),
	}
}

// Return maximum vector of a and b
func (a V3) Max(b V3) V3 {
	return V3{
		math.Max(a[0], b[0]),
		math.Max(a[1], b[1]),
		math.Max(a[2], b[2]),
	}
}

// Return maximum vector of a and b
func (a V2) Max(b V2) V2 {
	return V2{
		math.Max(a[0], b[0]),
		math.Max(a[1], b[1]),
	}
}

// Return minimum vector of a and b
func (a V3) Min(b V3) V3 {
	return V3{
		math.Min(a[0], b[0]),
		math.Min(a[1], b[1]),
		math.Min(a[2], b[2]),
	}
}

func (a V2) Min(b V2) V2 {
	return V2{
		math.Min(a[0], b[0]),
		math.Min(a[1], b[1]),
	}
}

// Return maximum component of a
func (a V3) Vmax() float64 {
	return math.Max(math.Max(a[0], a[1]), a[2])
}

func (a V2) Vmax() float64 {
	return math.Max(a[0], a[1])
}

// Return minimum component of a
func (a V3) Vmin() float64 {
	return math.Min(math.Min(a[0], a[1]), a[2])
}

func (a V2) Vmin() float64 {
	return math.Min(a[0], a[1])
}

//-----------------------------------------------------------------------------
// Scalar Functions (similar to GLSL counterparts)

// Return 0 if x < edge, else 1
func Step(edge, x float64) float64 {
	if x < edge {
		return 0
	}
	return 1
}

// Linear Interpolation
func Mix(x, y, a float64) float64 {
	return (x * (1 - a)) + (y * a)
}

func Clamp(x, a, b float64) float64 {
	return math.Min(math.Max(x, a), b)
}

func Saturate(x float64) float64 {
	return Clamp(x, 0, 1)
}

//-----------------------------------------------------------------------------
