//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
)

//-----------------------------------------------------------------------------
// constants

// Pi (3.14159...)
const Pi = math.Pi

// Tau (2 * Pi).
const Tau = 2 * math.Pi

// MillimetresPerInch is millimetres per inch (25.4)
const MillimetresPerInch = 25.4

// Mil is millimetres per 1/1000 of an inch
const Mil = MillimetresPerInch / 1000.0

const sqrtHalf = 0.7071067811865476
const tolerance = 1e-9
const epsilon = 1e-12

//-----------------------------------------------------------------------------

// DtoR converts degrees to radians
func DtoR(degrees float64) float64 {
	return (Pi / 180) * degrees
}

// RtoD converts radians to degrees
func RtoD(radians float64) float64 {
	return (180 / Pi) * radians
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

// Mix does a linear interpolation from x to y, a = [0,1]
func Mix(x, y, a float64) float64 {
	return x + (a * (y - x))
}

//-----------------------------------------------------------------------------
// Max/Min functions
// Note: math.Max/math.Min don't inline

// Max returns the maximum of a and b
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of a and b
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

// Abs returns the absolute value of x
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}

// Sign returns the sign of x
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

// SawTooth generates a sawtooth function. Returns [-period/2, period/2)
func SawTooth(x, period float64) float64 {
	x += period / 2
	t := x / period
	return period*(t-math.Floor(t)) - period/2
}

//-----------------------------------------------------------------------------

// MinFunc is a minimum functions for SDF blending.
type MinFunc func(a, b float64) float64

// RoundMin returns a minimum function that uses a quarter-circle to join the two objects smoothly.
func RoundMin(k float64) MinFunc {
	return func(a, b float64) float64 {
		u := V2{k - a, k - b}.Max(V2{0, 0})
		return Max(k, Min(a, b)) - u.Length()
	}
}

// ChamferMin returns a minimum function that makes a 45-degree chamfered edge (the diagonal of a square of size <r>).
// TODO: why the holes in the rendering?
func ChamferMin(k float64) MinFunc {
	return func(a, b float64) float64 {
		return Min(Min(a, b), (a-k+b)*sqrtHalf)
	}
}

// ExpMin returns a minimum function with exponential smoothing (k = 32).
func ExpMin(k float64) MinFunc {
	return func(a, b float64) float64 {
		return -math.Log(math.Exp(-k*a)+math.Exp(-k*b)) / k
	}
}

// PowMin returns  a minimum function (k = 8).
// TODO - weird results, is this correct?
func PowMin(k float64) MinFunc {
	return func(a, b float64) float64 {
		a = math.Pow(a, k)
		b = math.Pow(b, k)
		return math.Pow((a*b)/(a+b), 1/k)
	}
}

func poly(a, b, k float64) float64 {
	h := Clamp(0.5+0.5*(b-a)/k, 0.0, 1.0)
	return Mix(b, a, h) - k*h*(1.0-h)
}

// PolyMin returns a minimum function (Try k = 0.1, a bigger k gives a bigger fillet).
func PolyMin(k float64) MinFunc {
	return func(a, b float64) float64 {
		return poly(a, b, k)
	}
}

//-----------------------------------------------------------------------------

// MaxFunc is a maximum function for SDF blending.
type MaxFunc func(a, b float64) float64

// PolyMax returns a maximum function (Try k = 0.1, a bigger k gives a bigger fillet).
func PolyMax(k float64) MaxFunc {
	return func(a, b float64) float64 {
		return -poly(-a, -b, k)
	}
}

//-----------------------------------------------------------------------------

// ExtrudeFunc maps V3 to V2 - the point used to evaluate the SDF2.
type ExtrudeFunc func(p V3) V2

// NormalExtrude returns an extrusion function.
func NormalExtrude(p V3) V2 {
	return V2{p.X, p.Y}
}

// TwistExtrude returns an extrusion function that twists with z.
func TwistExtrude(height, twist float64) ExtrudeFunc {
	k := twist / height
	return func(p V3) V2 {
		m := Rotate(p.Z * k)
		return m.MulPosition(V2{p.X, p.Y})
	}
}

// ScaleExtrude returns an extrusion functions that scales with z.
func ScaleExtrude(height float64, scale V2) ExtrudeFunc {
	inv := V2{1 / scale.X, 1 / scale.Y}
	m := inv.Sub(V2{1, 1}).DivScalar(height) // slope
	b := inv.DivScalar(2).AddScalar(0.5)     // intercept
	return func(p V3) V2 {
		return V2{p.X, p.Y}.Mul(m.MulScalar(p.Z).Add(b))
	}
}

// ScaleTwistExtrude returns an extrusion function that scales and twists with z.
func ScaleTwistExtrude(height, twist float64, scale V2) ExtrudeFunc {
	k := twist / height
	inv := V2{1 / scale.X, 1 / scale.Y}
	m := inv.Sub(V2{1, 1}).DivScalar(height) // slope
	b := inv.DivScalar(2).AddScalar(0.5)     // intercept
	return func(p V3) V2 {
		// Scale and then Twist
		pnew := V2{p.X, p.Y}.Mul(m.MulScalar(p.Z).Add(b)) // Scale
		return Rotate(p.Z * k).MulPosition(pnew)          // Twist

		// Twist and then scale
		//pnew := Rotate(p.Z * k).MulPosition(V2{p.X, p.Y})
		//return pnew.Mul(m.MulScalar(p.Z).Add(b))
	}
}

//-----------------------------------------------------------------------------

// FloatDecode returns a string that decodes the float64 bitfields.
func FloatDecode(x float64) string {
	i := math.Float64bits(x)
	s := int((i >> 63) & 1)
	f := i & ((1 << 52) - 1)
	e := int((i>>52)&((1<<11)-1)) - 1023
	return fmt.Sprintf("s %d f 0x%013x e %d", s, f, e)
}

// FloatEncode encodes a float64 from sign, fraction and exponent values.
func FloatEncode(s int, f uint64, e int) float64 {
	s &= 1
	exp := uint64(e+1023) & ((1 << 11) - 1)
	f &= (1 << 52) - 1
	return math.Float64frombits(uint64(s)<<63 | exp<<52 | f)
}

//-----------------------------------------------------------------------------
// Floating Point Comparisons
// See: http://floating-point-gui.de/errors/NearlyEqualsTest.java

const minNormal = 2.2250738585072014E-308 // 2**-1022

// EqualFloat64 compares two float64 values for equality.
func EqualFloat64(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)
	if a == 0 || b == 0 || diff < minNormal {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diff < (epsilon * minNormal)
	}
	// use relative error
	return diff/math.Min((absA+absB), math.MaxFloat64) < epsilon
}

//-----------------------------------------------------------------------------

// ZeroSmall zeroes out values that are small relative to a quantity.
func ZeroSmall(x, y, epsilon float64) float64 {
	if Abs(x)/y < epsilon {
		return 0
	}
	return x
}

//-----------------------------------------------------------------------------

// NextCombination generates the next k-length combination of 0 to n-1. (returns false when done).
func NextCombination(n int, a []int) bool {
	k := len(a)
	m := 0
	i := 0
	for {
		i++
		if i > k {
			return false
		}
		if a[k-i] < n-i {
			m = a[k-i]
			for j := i; j >= 1; j-- {
				m++
				a[k-j] = m
			}
			return true
		}
	}
}

// MapCombinations applies a function f to each k-length combination from 0 to n-1.
func MapCombinations(n, k int, f func([]int)) {
	if k >= 0 && n >= k {
		a := make([]int, k)
		for i := range a {
			a[i] = i
		}
		for {
			f(a)
			if NextCombination(n, a) == false {
				break
			}
		}
	}
}

//-----------------------------------------------------------------------------
