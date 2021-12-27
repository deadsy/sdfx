//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
	"runtime"
)

//-----------------------------------------------------------------------------
// constants

// Pi (3.14159...)
const Pi = math.Pi

// Tau (2 * Pi).
const Tau = 2 * math.Pi

// MillimetresPerInch is millimetres per inch (25.4)
const MillimetresPerInch = 25.4

// InchesPerMillimetre is inches per millimetre
const InchesPerMillimetre = 1.0 / MillimetresPerInch

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
		return math.Max(k, math.Min(a, b)) - u.Length()
	}
}

// ChamferMin returns a minimum function that makes a 45-degree chamfered edge (the diagonal of a square of size <r>).
// TODO: why the holes in the rendering?
func ChamferMin(k float64) MinFunc {
	return func(a, b float64) float64 {
		return math.Min(math.Min(a, b), (a-k+b)*sqrtHalf)
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
// Raycasting
//-----------------------------------------------------------------------------

func sigmoidScaled(x float64) float64 {
	return 2/(1+math.Exp(-x)) - 1
}

// Raycast3 collides a ray (with an origin point from and a direction dir) with an SDF3.
// sigmoid is useful for fixing bad distance functions (those that do not accurately represent the distance to the
// closest surface, but will probably imply more evaluations)
// stepScale controls precision (less stepSize, more precision, but more SDF evaluations): use 1 if SDF indicates
// distance to the closest surface.
// It returns the collision point, how many normalized distances to reach it (t), and the number of steps performed
// If no surface is found (in maxDist and maxSteps), t is < 0
func Raycast3(s SDF3, from, dir V3, scaleAndSigmoid, stepScale, epsilon, maxDist float64, maxSteps int) (collision V3, t float64, steps int) {
	t = 0
	dirN := dir.Normalize()
	pos := from
	for {
		val := math.Abs(s.Evaluate(pos))
		//log.Print("Raycast step #", steps, " at ", pos, " with value ", val, "\n")
		if val < epsilon {
			collision = pos // Success
			break
		}
		steps++
		if steps == maxSteps {
			t = -1 // Failure
			break
		}
		if scaleAndSigmoid > 0 {
			val = sigmoidScaled(val * 10)
		}
		delta := val * stepScale
		t += delta
		pos = pos.Add(dirN.MulScalar(delta))
		if t < 0 || t > maxDist {
			t = -1 // Failure
			break
		}
	}
	//log.Println("Raycast did", steps, "steps")
	return
}

// Raycast2 see Raycast3. NOTE: implementation using Raycast3 (inefficient?)
func Raycast2(s SDF2, from, dir V2, scaleAndSigmoid, stepScale, epsilon, maxDist float64, maxSteps int) (V2, float64, int) {
	collision, t, steps := Raycast3(Extrude3D(s, 1), from.ToV3(0), dir.ToV3(0), scaleAndSigmoid, stepScale, epsilon, maxDist, maxSteps)
	return V2{collision.X, collision.Y}, t, steps
}

//-----------------------------------------------------------------------------
// Normals
//-----------------------------------------------------------------------------

// Normal3 returns the normal of an SDF3 at a point (doesn't need to be on the surface).
// Computed by sampling it several times inside a box of side 2*eps centered on p.
func Normal3(s SDF3, p V3, eps float64) V3 {
	return V3{
		X: s.Evaluate(p.Add(V3{X: eps})) - s.Evaluate(p.Add(V3{X: -eps})),
		Y: s.Evaluate(p.Add(V3{Y: eps})) - s.Evaluate(p.Add(V3{Y: -eps})),
		Z: s.Evaluate(p.Add(V3{Z: eps})) - s.Evaluate(p.Add(V3{Z: -eps})),
	}.Normalize()
}

// Normal2 returns the normal of an SDF3 at a point (doesn't need to be on the surface).
// Computed by sampling it several times inside a box of side 2*eps centered on p.
func Normal2(s SDF2, p V2, eps float64) V2 {
	return V2{
		X: s.Evaluate(p.Add(V2{X: eps})) - s.Evaluate(p.Add(V2{X: -eps})),
		Y: s.Evaluate(p.Add(V2{Y: eps})) - s.Evaluate(p.Add(V2{Y: -eps})),
	}.Normalize()
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

const minNormal = 2.2250738585072014e-308 // 2**-1022

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
	if math.Abs(x)/y < epsilon {
		return 0
	}
	return x
}

//-----------------------------------------------------------------------------

// ErrMsg returns an error with a message function name and line number.
func ErrMsg(msg string) error {
	pc, _, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("?: %s", msg)
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Errorf("%s line %d: %s", fn.Name(), line, msg)
}

//-----------------------------------------------------------------------------
