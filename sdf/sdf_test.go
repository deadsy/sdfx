//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v2i"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/stretchr/testify/assert"
)

//-----------------------------------------------------------------------------

func Test_Determinant(t *testing.T) {
	m := M33{2, 3, 1, -1, -6, 7, 4, 5, -1}
	if m.Determinant() != 42 {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_Inverse(t *testing.T) {
	a := M33{2, 1, 1, 3, 2, 1, 2, 1, 2}
	aInv := M33{3, -1, -1, -4, 2, 1, -1, 0, 1}
	if a.Inverse().Equals(aInv, tolerance) == false {
		t.Error("FAIL")
	}
	if a.Mul(aInv).Equals(Identity2d(), tolerance) == false {
		t.Error("FAIL")
	}

	for i := 0; i < 100; i++ {
		a = RandomM33(-5, 5)
		aInv = a.Inverse()
		if a.Mul(aInv).Equals(Identity2d(), tolerance) == false {
			t.Error("FAIL")
		}
	}

	for i := 0; i < 100; i++ {
		b := RandomM44(-1, 1)
		bInv := b.Inverse()
		if b.Mul(bInv).Equals(Identity3d(), tolerance) == false {
			t.Error("FAIL")
		}
	}

	for i := 0; i < 100; i++ {
		c := RandomM22(-7, 7)
		cInv := c.Inverse()
		if c.Mul(cInv).Equals(Identity(), tolerance) == false {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_MulBox(t *testing.T) {

	// 2D boxes
	b2d := Box2{v2.Vec{-1, -1}, v2.Vec{1, 1}}
	for i := 0; i < 100; i++ {
		b := NewBox2(v2.Vec{0, 0}, v2.Vec{10, 10})
		v := b.Random()
		// translating
		m0 := Translate2d(v)
		m1 := Translate2d(v.Neg())
		b1 := m0.MulBox(b2d)
		b2 := m1.MulBox(b1)
		if b2d.Equals(b2, tolerance) == false {
			t.Error("FAIL")
		}
		// scaling
		m0 = Scale2d(v)
		m1 = Scale2d(v2.Vec{1 / v.X, 1 / v.Y})
		b1 = m0.MulBox(b2d)
		b2 = m1.MulBox(b1)
		if b2d.Equals(b2, tolerance) == false {
			t.Error("FAIL")
		}
	}

	// 3D boxes
	b3d := Box3{v3.Vec{-1, -1, -1}, v3.Vec{1, 1, 1}}
	for i := 0; i < 100; i++ {
		b := NewBox3(v3.Vec{0, 0, 0}, v3.Vec{10, 10, 10})
		v := b.Random()
		// translating
		m0 := Translate3d(v)
		m1 := Translate3d(v.Neg())
		b1 := m0.MulBox(b3d)
		b2 := m1.MulBox(b1)
		if b3d.Equals(b2, tolerance) == false {
			t.Error("FAIL")
		}
		// scaling
		m0 = Scale3d(v)
		m1 = Scale3d(v3.Vec{1 / v.X, 1 / v.Y, 1 / v.Z})
		b1 = m0.MulBox(b3d)
		b2 = m1.MulBox(b1)
		if b3d.Equals(b2, tolerance) == false {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_ScaleBox(t *testing.T) {
	b0 := Box3{v3.Vec{-1, -1, -1}, v3.Vec{1, 1, 1}}
	b1 := Box3{v3.Vec{-2, -2, -2}, v3.Vec{2, 2, 2}}
	b2 := NewBox3(b0.Center(), b0.Size().MulScalar(2))
	if b1.Equals(b2, tolerance) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_Line(t *testing.T) {

	l := newLinePP(v2.Vec{0, 1}, v2.Vec{0, 2})
	points := []struct {
		p v2.Vec
		d float64
	}{
		{v2.Vec{0, 1}, 0},
		{v2.Vec{0, 2}, 0},
		{v2.Vec{0, 1.5}, 0},
		{v2.Vec{0, 0}, 1},
		{v2.Vec{0, 3}, 1},
		{v2.Vec{0.5, 1.1}, 0.5},
		{v2.Vec{-0.5, 1.1}, -0.5},
		{v2.Vec{0.1, 1.98}, 0.1},
		{v2.Vec{-0.1, 1.98}, -0.1},
		{v2.Vec{3, 6}, 5},
		{v2.Vec{-3, 6}, -5},
		{v2.Vec{3, -3}, 5},
		{v2.Vec{-3, -3}, -5},
	}
	for _, p := range points {
		d := l.Distance(p.p)
		if math.Abs(d-p.d) > tolerance {
			fmt.Printf("%+v %f (expected) %f (actual)\n", p.p, p.d, d)
			t.Error("FAIL")
		}
	}

	lineTests := []struct {
		p0, v0 v2.Vec
		p1, v1 v2.Vec
		t0     float64
		t1     float64
		err    string
	}{
		{v2.Vec{0, 0}, v2.Vec{0, 1}, v2.Vec{1, 0}, v2.Vec{-1, 0}, 0, 1, ""},
		{v2.Vec{0, 0}, v2.Vec{0, 1}, v2.Vec{1, 1}, v2.Vec{0, 1}, 0, 1, "zero/many"},
		{v2.Vec{0, 0}, v2.Vec{0, 1}, v2.Vec{0, 0}, v2.Vec{0, 1}, 0, 1, "zero/many"},
		{v2.Vec{0, 0}, v2.Vec{1, 1}, v2.Vec{0, 10}, v2.Vec{1, -1}, 5 * math.Sqrt(2), 5 * math.Sqrt(2), ""},
		{v2.Vec{0, 0}, v2.Vec{1, 1}, v2.Vec{10, 0}, v2.Vec{0, 1}, 10 * math.Sqrt(2), 10, ""},
	}
	for _, test := range lineTests {
		l0 := newLinePV(test.p0, test.v0)
		l1 := newLinePV(test.p1, test.v1)
		t0, t1, err := l0.Intersect(l1)
		if err != nil {
			if err.Error() != test.err {
				fmt.Printf("l0: %+v\n", l0)
				fmt.Printf("l1: %+v\n", l1)
				fmt.Printf("error: %s\n", err)
				t.Error("FAIL")
			}
		} else {
			if math.Abs(test.t0-t0) > tolerance || math.Abs(test.t1-t1) > tolerance {
				fmt.Printf("l0: %+v\n", l0)
				fmt.Printf("l1: %+v\n", l1)
				fmt.Printf("%f %f (expected) %f %f (actual)\n", test.t0, test.t1, t0, t1)
				t.Error("FAIL")
			}
		}
	}

	for i := 0; i < 10000; i++ {
		b := NewBox2(v2.Vec{0, 0}, v2.Vec{20, 20})
		l0 := newLinePV(b.Random(), b.Random())
		l1 := newLinePP(b.Random(), b.Random())
		t0, t1, err := l0.Intersect(l1)
		if err != nil {
			continue
		}
		i0 := l0.Position(t0)
		i1 := l1.Position(t1)
		if !i0.Equals(i1, tolerance) {
			fmt.Printf("l0: %+v\n", l0)
			fmt.Printf("l1: %+v\n", l1)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Polygon1(t *testing.T) {
	s, _ := Polygon2D([]v2.Vec{{0, 0}, {1, 0}, {0, 1}})
	b := s.BoundingBox()
	b0 := Box2{v2.Vec{0, 0}, v2.Vec{1, 1}}
	if b.Equals(b0, tolerance) == false {
		t.Error("FAIL")
	}

	s, _ = Polygon2D([]v2.Vec{{0, -2}, {1, 1}, {-2, 2}})
	b = s.BoundingBox()
	b0 = Box2{v2.Vec{-2, -2}, v2.Vec{1, 2}}
	if b.Equals(b0, tolerance) == false {
		t.Error("FAIL")
	}

	points := []v2.Vec{
		{0, -1},
		{1, 1},
		{-1, 1},
	}

	s, _ = Polygon2D(points)

	b = s.BoundingBox()
	b0 = Box2{v2.Vec{-1, -1}, v2.Vec{1, 1}}
	if b.Equals(b0, tolerance) == false {
		t.Error("FAIL")
	}

	testPoints := []struct {
		p v2.Vec
		d float64
	}{
		{v2.Vec{0, -1}, 0},
		{v2.Vec{1, 1}, 0},
		{v2.Vec{-1, 1}, 0},
		{v2.Vec{0, 1}, 0},
		{v2.Vec{0, 2}, 1},
		{v2.Vec{0, -2}, 1},
		{v2.Vec{1, 0}, 1 / math.Sqrt(5)},
		{v2.Vec{-1, 0}, 1 / math.Sqrt(5)},
		{v2.Vec{0, 0}, -1 / math.Sqrt(5)},
		{v2.Vec{3, 0}, math.Sqrt(5)},
		{v2.Vec{-3, 0}, math.Sqrt(5)},
	}

	for _, p := range testPoints {
		d := s.Evaluate(p.p)
		if d != p.d {
			fmt.Printf("%+v %f (expected) %f (actual)\n", p.p, p.d, d)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Polygon2(t *testing.T) {
	k := 1.2

	s0, _ := Polygon2D([]v2.Vec{{k, -k}, {k, k}, {-k, k}, {-k, -k}})
	s0 = Transform2D(s0, Translate2d(v2.Vec{0.8, 0}))

	s1 := Box2D(v2.Vec{2 * k, 2 * k}, 0)
	s1 = Transform2D(s1, Translate2d(v2.Vec{0.8, 0}))

	for i := 0; i < 10000; i++ {
		b := NewBox2(v2.Vec{0, 0}, v2.Vec{20 * k, 20 * k})
		p := b.Random()
		if math.Abs(s0.Evaluate(p)-s1.Evaluate(p)) > tolerance {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Polygon3(t *testing.T) {

	// size
	a := 1.4
	b := 2.2
	// rotation
	theta := -15.0
	c := math.Cos(DtoR(theta))
	s := math.Sin(DtoR(theta))
	// translate
	j := -1.0
	k := 2.0

	s1 := Box2D(v2.Vec{2 * a, 2 * b}, 0)
	s1 = Transform2D(s1, Rotate2d(DtoR(theta)))
	s1 = Transform2D(s1, Translate2d(v2.Vec{j, k}))

	points := []v2.Vec{
		{j + c*a - s*b, k + s*a + c*b},
		{j - c*a - s*b, k - s*a + c*b},
		{j - c*a + s*b, k - s*a - c*b},
		{j + c*a + s*b, k + s*a - c*b},
	}

	s0, _ := Polygon2D(points)

	for i := 0; i < 1000; i++ {
		b := NewBox2(v2.Vec{0, 0}, v2.Vec{10 * b, 10 * b})
		p := b.Random()
		if math.Abs(s0.Evaluate(p)-s1.Evaluate(p)) > tolerance {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_ArraySDF2(t *testing.T) {
	r := 0.5
	s, _ := Circle2D(r)
	bb := s.BoundingBox()
	if bb.Min.Equals(v2.Vec{-r, -r}, tolerance) == false {
		t.Error("FAIL")
	}
	if bb.Max.Equals(v2.Vec{r, r}, tolerance) == false {
		t.Error("FAIL")
	}

	j := 3
	k := 4
	dx := 2.0
	dy := 7.0
	sa := Array2D(s, v2i.Vec{j, k}, v2.Vec{dx, dy})
	saBox := sa.BoundingBox()
	if saBox.Min.Equals(v2.Vec{-r, -r}, tolerance) == false {
		t.Error("FAIL")
	}
	if saBox.Max.Equals(v2.Vec{r + (float64(j-1) * dx), r + (float64(k-1) * dy)}, tolerance) == false {
		t.Error("FAIL")
	}

	j = 7
	k = 4
	dx = -3.0
	dy = -5.0
	sa = Array2D(s, v2i.Vec{j, k}, v2.Vec{dx, dy})
	saBox = sa.BoundingBox()
	if saBox.Min.Equals(v2.Vec{-r + (float64(j-1) * dx), -r + (float64(k-1) * dy)}, tolerance) == false {
		t.Error("FAIL")
	}
	if saBox.Max.Equals(v2.Vec{r, r}, tolerance) == false {
		t.Error("FAIL")
	}

	j = 6
	k = 8
	dx = 5.0
	dy = -3.0
	sa = Array2D(s, v2i.Vec{j, k}, v2.Vec{dx, dy})
	saBox = sa.BoundingBox()
	if saBox.Min.Equals(v2.Vec{-r, -r + (float64(k-1) * dy)}, tolerance) == false {
		t.Error("FAIL")
	}
	if saBox.Max.Equals(v2.Vec{r + (float64(j-1) * dx), r}, tolerance) == false {
		t.Error("FAIL")
	}

	j = 9
	k = 1
	dx = -0.5
	dy = 6.5
	sa = Array2D(s, v2i.Vec{j, k}, v2.Vec{dx, dy})
	saBox = sa.BoundingBox()
	if saBox.Min.Equals(v2.Vec{-r + (float64(j-1) * dx), -r}, tolerance) == false {
		t.Error("FAIL")
	}
	if saBox.Max.Equals(v2.Vec{r, r + (float64(k-1) * dy)}, tolerance) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_Rotation2d(t *testing.T) {
	r := Rotate2d(DtoR(90))
	v := v2.Vec{1, 0}
	v = r.MulPosition(v)
	if v.Equals(v2.Vec{0, 1}, tolerance) == false {
		t.Error("FAIL")
	}
}

func Test_Rotation3d(t *testing.T) {
	r := Rotate3d(v3.Vec{0, 0, 1}, DtoR(90))
	v := v3.Vec{1, 0, 0}
	v = r.MulPosition(v)
	if v.Equals(v3.Vec{0, 1, 0}, tolerance) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_TriDiagonal(t *testing.T) {
	n := 5
	m := make([]v3.Vec, n)
	for i := 0; i < n; i++ {
		m[i].X = 1
		m[i].Y = 4
		m[i].Z = 1
	}
	m[0].X = 0
	m[0].Y = 2
	m[n-1].Y = 2
	m[n-1].Z = 0

	d := []float64{0, 1, 2, 3, 4}
	x0 := []float64{-1.0 / 12.0, 1.0 / 6.0, 5.0 / 12.0, 1.0 / 6.0, 23.0 / 12.0}
	x, err := triDiagonal(m, d)
	if err != nil {
		t.Error("FAIL")
	}
	for i := 0; i < n; i++ {
		if math.Abs(x[i]-x0[i]) > tolerance {
			t.Error("FAIL")
		}
	}

	d = []float64{10, 20, 30, 40, 50}
	x0 = []float64{15.0 / 4.0, 5.0 / 2.0, 25.0 / 4.0, 5.0 / 2.0, 95.0 / 4.0}
	x, err = triDiagonal(m, d)
	if err != nil {
		t.Error("FAIL")
	}
	for i := 0; i < n; i++ {
		if math.Abs(x[i]-x0[i]) > tolerance {
			t.Error("FAIL")
		}
	}

	m[0] = v3.Vec{0, 1, 2}
	m[1] = v3.Vec{3, 4, 5}
	m[2] = v3.Vec{6, 7, 8}
	m[3] = v3.Vec{9, 10, 11}
	m[4] = v3.Vec{12, 13, 0}
	d = []float64{-10, -20, -30, 40, 50}
	x0 = []float64{60.0 / 49.0, -275.0 / 49.0, -12.0 / 49.0, 33.0 / 49.0, 158.0 / 49.0}
	x, err = triDiagonal(m, d)
	if err != nil {
		t.Error("FAIL")
	}
	for i := 0; i < n; i++ {
		if math.Abs(x[i]-x0[i]) > tolerance {
			t.Error("FAIL")
		}
	}

}

//-----------------------------------------------------------------------------

func Test_CubicSpline(t *testing.T) {

	knot := []v2.Vec{
		{-1.5, -1.2},
		{-0.2, 0},
		{1, 0.5},
		{5, 1},
		{10, 2.2},
		{12, 3.2},
		{-16, -1.2},
		{-18, -3.2},
	}
	s0, err := CubicSpline2D(knot)
	if err != nil {
		t.Error("FAIL")
	}
	s := s0.(*CubicSplineSDF2)
	n := len(s.spline)
	if n != len(knot)-1 {
		t.Error("FAIL")
	}
	// check interpolation of the knots
	for i, k := range knot {
		p := s.f0(float64(i))
		if !k.Equals(p, tolerance) {
			t.Error("FAIL")
		}
	}
	// check 1st and 2nd derivatives
	for i := 0; i < n; i++ {
		cs := s.spline[i]
		if i == 0 {
			// 2nd derivative at start == 0
			f2 := cs.f2(0)
			if !f2.Equals(v2.Vec{0, 0}, tolerance) {
				t.Error("FAIL")
			}
		}
		if i == n-1 {
			// 2nd derivative at end == 0
			f2 := cs.f2(1)
			if !f2.Equals(v2.Vec{0, 0}, tolerance) {
				t.Error("FAIL")
			}
		} else {
			csNext := s.spline[i+1]
			// check continuity of 1st derivative
			if !cs.f1(1).Equals(csNext.f1(0), tolerance) {
				t.Error("FAIL")
			}
			// check continuity of 2nd derivative
			if !cs.f2(1).Equals(csNext.f2(0), tolerance) {
				t.Error("FAIL")
			}
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Quadratic(t *testing.T) {

	x, rc := quadratic(4, 2, 1)
	if x != nil || rc != zeroSoln {
		t.Error("FAIL")
	}

	x, rc = quadratic(0, 0, 1)
	if x != nil || rc != zeroSoln {
		t.Error("FAIL")
	}

	x, rc = quadratic(0, 2, -4)
	if x[0] != 2 || rc != oneSoln {
		t.Error("FAIL")
	}

	x, rc = quadratic(1, -5, 6)
	if x[0] != 3 || x[1] != 2 || rc != twoSoln {
		t.Error("FAIL")
	}

	x, rc = quadratic(0, 0, 0)
	if x != nil || rc != infSoln {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_Colinear_Fast(t *testing.T) {
	test := []struct {
		a, b, c v2.Vec
		result  bool
	}{
		{v2.Vec{}, v2.Vec{}, v2.Vec{}, true},
		{v2.Vec{}, v2.Vec{}, v2.Vec{10, 17}, true},
		{v2.Vec{1, 1}, v2.Vec{1, 2}, v2.Vec{1, 3}, true},
		{v2.Vec{1, 1}, v2.Vec{2, 2}, v2.Vec{3, 3}, true},
		{v2.Vec{1, 1}, v2.Vec{2, 1}, v2.Vec{3, 1}, true},
		{v2.Vec{1, 1}, v2.Vec{2, 5}, v2.Vec{-2, 4}, false},
		{v2.Vec{1, 1}, v2.Vec{1, 1}, v2.Vec{-2, 4}, true},
		{v2.Vec{1, 1}, v2.Vec{-1, 1}, v2.Vec{0, -1}, false},
	}
	for _, v := range test {
		if colinearFast(v.a, v.b, v.c, epsilon) != v.result {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Float_Comparison(t *testing.T) {

	nInf := math.Inf(-1)
	pInf := math.Inf(1)
	nan := math.NaN()
	maxValue := math.MaxFloat64
	minValue := math.SmallestNonzeroFloat64

	epsilon := 0.00001

	test0 := []struct {
		a, b   float64
		equals bool
	}{
		// Regular large numbers - generally not problematic
		{1000000, 1000001, true},
		{1000001, 1000000, true},
		{10000, 10001, false},
		{10001, 10000, false},
		// Negative large numbers
		{-1000000, -1000001, true},
		{-1000001, -1000000, true},
		{-10000, -10001, false},
		{-10001, -10000, false},
		// Numbers around 1
		{1.0000001, 1.0000002, true},
		{1.0000002, 1.0000001, true},
		{1.0002, 1.0001, false},
		{1.0001, 1.0002, false},
		// Numbers around -1
		{-1.0000001, -1.0000002, true},
		{-1.0000002, -1.0000001, true},
		{-1.0002, -1.0001, false},
		{-1.0001, -1.0002, false},
		// Numbers between 1 and 0
		{0.000000001000001, 0.000000001000002, true},
		{0.000000001000002, 0.000000001000001, true},
		{0.000000000001002, 0.000000000001001, false},
		{0.000000000001001, 0.000000000001002, false},
		// Numbers between -1 and 0
		{-0.000000001000001, -0.000000001000002, true},
		{-0.000000001000002, -0.000000001000001, true},
		{-0.000000000001002, -0.000000000001001, false},
		{-0.000000000001001, -0.000000000001002, false},
		// Comparisons involving zero
		{0, 0, true},
		{0, -0, true},
		{-0, -0, true},
		{0.00000001, 0, false},
		{0, 0.00000001, false},
		{-0.00000001, 0, false},
		{0, -0.00000001, false},
		// Comparisons of numbers on opposite sides of 0
		{1.000000001, -1.0, false},
		{-1.0, 1.000000001, false},
		{-1.000000001, 1.0, false},
		{1.0, -1.000000001, false},
		//{10 * minValue, 10 * -minValue, true},        // problem
		//{10000 * minValue, 10000 * -minValue, false}, // problem
		// The really tricky part - comparisons of numbers very close to zero.
		{minValue, minValue, true},
		{minValue, -minValue, true},
		{-minValue, minValue, true},
		{minValue, 0, true},
		{0, minValue, true},
		{-minValue, 0, true},
		{0, -minValue, true},
		{0.000000001, -minValue, false},
		{0.000000001, minValue, false},
		{minValue, 0.000000001, false},
		{-minValue, 0.000000001, false},
		// Comparisons involving NaN values
		{nan, nan, false},
		{nan, 0.0, false},
		{-0.0, nan, false},
		{nan, -0.0, false},
		{0.0, nan, false},
		{nan, pInf, false},
		{pInf, nan, false},
		{nan, nInf, false},
		{nInf, nan, false},
		{nan, maxValue, false},
		{maxValue, nan, false},
		{nan, -maxValue, false},
		{-maxValue, nan, false},
		{nan, minValue, false},
		{minValue, nan, false},
		{nan, -minValue, false},
		{-minValue, nan, false},
		// Comparisons involving extreme values (overflow potential)
		{maxValue, maxValue, true},
		{maxValue, -maxValue, false},
		{-maxValue, maxValue, false},
		{maxValue, maxValue / 2, false},
		{maxValue, -maxValue / 2, false},
		{-maxValue, maxValue / 2, false},
		// Comparisons involving infinities
		{pInf, pInf, true},
		{nInf, nInf, true},
		{nInf, pInf, false},
		{pInf, maxValue, false},
		{nInf, -maxValue, false},
	}

	for _, v := range test0 {
		if EqualFloat64(v.a, v.b, epsilon) != v.equals {
			t.Error("FAIL")
		}
	}

	test1 := []struct {
		a, b, e float64
		equals  bool
	}{
		// Comparisons involving zero
		//{0.0, 1e-40, 0.01, true},
		//{1e-40, 0.0, 0.01, true},
		//{1e-40, 0.0, 0.000001, false},
		//{0.0, 1e-40, 0.000001, false},
		//{0.0, -1e-40, 0.1, true},
		//{-1e-40, 0.0, 0.1, true},
		//{-1e-40, 0.0, 0.00000001, false},
		//{0.0, -1e-40, 0.00000001, false},
	}

	for _, v := range test1 {
		if EqualFloat64(v.a, v.b, v.e) != v.equals {
			t.Error("FAIL")
		}
	}

}

//-----------------------------------------------------------------------------

func Test_Box2_Distances(t *testing.T) {
	b0 := NewBox2(v2.Vec{0, 0}, v2.Vec{10, 10})
	b1 := NewBox2(v2.Vec{10, 20}, v2.Vec{30, 40})
	tests := []struct {
		b      Box2
		p      v2.Vec
		result Interval
	}{
		{b0, v2.Vec{0, 0}, Interval{0, 50}},
		{b0, v2.Vec{5, 5}, Interval{0, 200}},
		{b0, v2.Vec{20, 0}, Interval{225, 650}},
		{b1, v2.Vec{0, 0}, Interval{0, 2225}},
		{b1, v2.Vec{10, 20}, Interval{0, 625}},
		{b1, v2.Vec{0, -10}, Interval{100, 3125}},
		{b1, v2.Vec{0, 5}, Interval{0, 1850}},
	}
	for _, v := range tests {
		x := v.b.MinMaxDist2(v.p)
		if !x.Equals(v.result, tolerance) {
			t.Logf("expected %v, actual %v\n", v.result, x)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Rotate_To_Vector(t *testing.T) {
	tests := []struct {
		a, b   v3.Vec
		result M44
	}{
		{v3.Vec{0, 0, 1}, v3.Vec{0, 0, 1}, M44{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}},
		{v3.Vec{0, 0, 1}, v3.Vec{0, 0, -1}, M44{-1, 0, 0, 0, 0, -1, 0, 0, 0, 0, -1, 0, 0, 0, 0, 1}},
		{v3.Vec{1, 0, 0}, v3.Vec{0, 0, 1}, M44{0, 0, -1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1}},
		{v3.Vec{1, 0, 1}, v3.Vec{-1, 0, 1}, M44{0, 0, -1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1}},
	}
	for _, v := range tests {
		x := RotateToVector(v.a, v.b)
		if !x.Equals(v.result, tolerance) {
			t.Logf("expected %v, actual %v\n", v.result, x)
			t.Error("FAIL")
		}
	}

	box := NewBox3(v3.Vec{}, v3.Vec{100, 100, 100})
	for i := 0; i < 1000; i++ {
		a := box.Random()
		b := box.Random()
		ax := RotateToVector(a, b).MulPosition(a)
		// ax should have the same magnitude as a
		if math.Abs(ax.Length()-a.Length()) > 1e-10 {
			t.Error("FAIL")
		}
		// ax should have the same direction as b
		ax = ax.Normalize()
		b = b.Normalize()
		if !ax.Equals(b, epsilon) {
			t.Error("FAIL")
		}
	}

}

//-----------------------------------------------------------------------------

func TestColinearity(t *testing.T) {
	a := v2.Vec{37.4, 88.8}
	m := v2.Vec{3.0, 5.0}
	b := a.Add(m.MulScalar(16.0))
	c := a.Sub(m.MulScalar(7.0))
	d := v2.Vec{55.5, 66.6}

	assert.True(t, colinearFast(a, b, c, 0.0001), "ABC are colinear fast")
	assert.True(t, colinearFast(a, c, b, 0.0001), "ACB are colienar fast")
	assert.True(t, colinearFast(b, a, c, 0.0001), "BAC are colinear fast")
	assert.True(t, colinearFast(b, c, a, 0.0001), "BCA are colienar fast")
	assert.True(t, colinearFast(c, a, b, 0.0001), "CAB are colinear fast")
	assert.True(t, colinearFast(c, b, a, 0.0001), "CBA are colinear fast")

	assert.False(t, colinearFast(a, b, d, 0.0001), "ABD are not colinear fast")
	assert.False(t, colinearFast(a, c, d, 0.0001), "ACD are not colienar fast")
	assert.False(t, colinearFast(b, a, d, 0.0001), "BAD are not colinear fast")
	assert.False(t, colinearFast(b, c, d, 0.0001), "BCD are not colienar fast")
	assert.False(t, colinearFast(c, a, d, 0.0001), "CAD are not colinear fast")
	assert.False(t, colinearFast(c, b, d, 0.0001), "CBD are not colinear fast")

	assert.True(t, colinearSlow(a, b, c, 0.0001), "ABC are colinear slow")
	assert.True(t, colinearSlow(a, c, b, 0.0001), "ACB are colienar slow")
	assert.True(t, colinearSlow(b, a, c, 0.0001), "BAC are colinear slow")
	assert.True(t, colinearSlow(b, c, a, 0.0001), "BCA are colienar slow")
	assert.True(t, colinearSlow(c, a, b, 0.0001), "CAB are colinear slow")
	assert.True(t, colinearSlow(c, b, a, 0.0001), "CBA are colinear slow")

	assert.False(t, colinearSlow(a, b, d, 0.0001), "ABD are not colinear slow")
	assert.False(t, colinearSlow(a, c, d, 0.0001), "ACD are not colienar slow")
	assert.False(t, colinearSlow(b, a, d, 0.0001), "BAD are not colinear slow")
	assert.False(t, colinearSlow(b, c, d, 0.0001), "BCD are not colienar slow")
	assert.False(t, colinearSlow(c, a, d, 0.0001), "CAD are not colinear slow")
	assert.False(t, colinearSlow(c, b, d, 0.0001), "CBD are not colinear slow")
}

//-----------------------------------------------------------------------------

func Test_Raycast(t *testing.T) {
	testSdf := Box2D(v2.Vec{1, 1}, 0.2)
	eps := 1e-10
	side, td, steps := Raycast2(testSdf, v2.Vec{-1.32442, 0}, v2.Vec{1, 0}, 0, 1, eps, 5, 10)
	if math.Abs(side.X-(-0.5)) > eps {
		t.Fatal("Should have collided with the side of the cube at -0.5, but got", side, td, steps)
	}
	side, td, steps = Raycast2(testSdf, v2.Vec{-1.32442, 0}, v2.Vec{1, 0}, 1, 0.1, eps, 5, 500)
	if math.Abs(side.X-(-0.5)) > eps {
		t.Fatal("Should have collided with the side of the cube at -0.5, but got", side, td, steps)
	}
	side, td, steps = Raycast2(testSdf, v2.Vec{-1.32442, 0}, v2.Vec{1, 0}, 0, 0.1, eps, 5, 1000)
	if math.Abs(side.X-(-0.5)) > eps {
		t.Fatal("Should have collided with the side of the cube at -0.5, but got", side, td, steps)
	}
	side, td, steps = Raycast2(testSdf, v2.Vec{-1.32442, 0}, v2.Vec{1, 0}, 0, 1, eps, 0.1, 10)
	if td >= 0 {
		t.Fatal("Should have reached maxDist and not collided", side, td, steps)
	}
	side, td, steps = Raycast2(testSdf, v2.Vec{-1.32442, -1.32442}, v2.Vec{1, 1}, 0, 0.1, eps, 5, 10)
	if side.X < -0.45 {
		t.Fatal("Should have collided with the ROUND side of the cube at <-0.45, but got", side, td, steps)
	}
	side, td, steps = Raycast2(testSdf, v2.Vec{0, 0}, v2.Vec{1, 0}, 0, 1, eps, 5, 10)
	if math.Abs(side.X-(0.5)) > eps {
		t.Fatal("Should have returned surface of the SDF from the inside", side, td, steps)
	}
}

//-----------------------------------------------------------------------------

func Test_Normal(t *testing.T) {
	testSdf := Box2D(v2.Vec{1, 1}, 0.2)
	eps := 1e-10
	n := Normal2(testSdf, v2.Vec{0.5, 0}, eps)
	if !reflect.DeepEqual(n, v2.Vec{1, 0}) {
		t.Fatal("Bad normal for box's right side: expected {1 0} but got", n)
	}
	n = Normal2(testSdf, v2.Vec{0.25, 0}, eps)
	if !reflect.DeepEqual(n, v2.Vec{1, 0}) {
		t.Fatal("Bad normal for box's right side (inside): expected {1 0} but got", n)
	}
	n = Normal2(testSdf, v2.Vec{0.5, 0.25}, eps)
	if !reflect.DeepEqual(n, v2.Vec{1, 0}) {
		t.Fatal("Bad normal for box's right side (displaced): expected {1 0} but got", n)
	}
	n = Normal2(testSdf, v2.Vec{0.45, 0.45}, eps)
	if !(n.X > 0 && n.Y > 0) {
		t.Fatal("Bad normal for box's right side (corner): expected {>0 >0} but got", n)
	}
}

//-----------------------------------------------------------------------------
