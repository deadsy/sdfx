//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
	"testing"
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
	b2d := Box2{V2{-1, -1}, V2{1, 1}}
	for i := 0; i < 100; i++ {
		b := NewBox2(V2{0, 0}, V2{10, 10})
		v := b.Random()
		// translating
		m0 := Translate2d(v)
		m1 := Translate2d(v.Negate())
		b1 := m0.MulBox(b2d)
		b2 := m1.MulBox(b1)
		if b2d.Equals(b2, tolerance) == false {
			t.Error("FAIL")
		}
		// scaling
		m0 = Scale2d(v)
		m1 = Scale2d(V2{1 / v.X, 1 / v.Y})
		b1 = m0.MulBox(b2d)
		b2 = m1.MulBox(b1)
		if b2d.Equals(b2, tolerance) == false {
			t.Error("FAIL")
		}
	}

	// 3D boxes
	b3d := Box3{V3{-1, -1, -1}, V3{1, 1, 1}}
	for i := 0; i < 100; i++ {
		b := NewBox3(V3{0, 0, 0}, V3{10, 10, 10})
		v := b.Random()
		// translating
		m0 := Translate3d(v)
		m1 := Translate3d(v.Negate())
		b1 := m0.MulBox(b3d)
		b2 := m1.MulBox(b1)
		if b3d.Equals(b2, tolerance) == false {
			t.Error("FAIL")
		}
		// scaling
		m0 = Scale3d(v)
		m1 = Scale3d(V3{1 / v.X, 1 / v.Y, 1 / v.Z})
		b1 = m0.MulBox(b3d)
		b2 = m1.MulBox(b1)
		if b3d.Equals(b2, tolerance) == false {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_ScaleBox(t *testing.T) {
	b0 := Box3{V3{-1, -1, -1}, V3{1, 1, 1}}
	b1 := Box3{V3{-2, -2, -2}, V3{2, 2, 2}}
	b2 := NewBox3(b0.Center(), b0.Size().MulScalar(2))
	if b1.Equals(b2, tolerance) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_Line(t *testing.T) {

	l := newLinePP(V2{0, 1}, V2{0, 2})
	points := []struct {
		p V2
		d float64
	}{
		{V2{0, 1}, 0},
		{V2{0, 2}, 0},
		{V2{0, 1.5}, 0},
		{V2{0, 0}, 1},
		{V2{0, 3}, 1},
		{V2{0.5, 1.1}, 0.5},
		{V2{-0.5, 1.1}, -0.5},
		{V2{0.1, 1.98}, 0.1},
		{V2{-0.1, 1.98}, -0.1},
		{V2{3, 6}, 5},
		{V2{-3, 6}, -5},
		{V2{3, -3}, 5},
		{V2{-3, -3}, -5},
	}
	for _, p := range points {
		d := l.Distance(p.p)
		if Abs(d-p.d) > tolerance {
			fmt.Printf("%+v %f (expected) %f (actual)\n", p.p, p.d, d)
			t.Error("FAIL")
		}
	}

	lineTests := []struct {
		p0, v0 V2
		p1, v1 V2
		t0     float64
		t1     float64
		err    string
	}{
		{V2{0, 0}, V2{0, 1}, V2{1, 0}, V2{-1, 0}, 0, 1, ""},
		{V2{0, 0}, V2{0, 1}, V2{1, 1}, V2{0, 1}, 0, 1, "zero/many"},
		{V2{0, 0}, V2{0, 1}, V2{0, 0}, V2{0, 1}, 0, 1, "zero/many"},
		{V2{0, 0}, V2{1, 1}, V2{0, 10}, V2{1, -1}, 5 * math.Sqrt(2), 5 * math.Sqrt(2), ""},
		{V2{0, 0}, V2{1, 1}, V2{10, 0}, V2{0, 1}, 10 * math.Sqrt(2), 10, ""},
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
			if Abs(test.t0-t0) > tolerance || Abs(test.t1-t1) > tolerance {
				fmt.Printf("l0: %+v\n", l0)
				fmt.Printf("l1: %+v\n", l1)
				fmt.Printf("%f %f (expected) %f %f (actual)\n", test.t0, test.t1, t0, t1)
				t.Error("FAIL")
			}
		}
	}

	for i := 0; i < 10000; i++ {
		b := NewBox2(V2{0, 0}, V2{20, 20})
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
	s := Polygon2D([]V2{{0, 0}, {1, 0}, {0, 1}})
	b := s.BoundingBox()
	b0 := Box2{V2{0, 0}, V2{1, 1}}
	if b.Equals(b0, tolerance) == false {
		t.Error("FAIL")
	}

	s = Polygon2D([]V2{{0, -2}, {1, 1}, {-2, 2}})
	b = s.BoundingBox()
	b0 = Box2{V2{-2, -2}, V2{1, 2}}
	if b.Equals(b0, tolerance) == false {
		t.Error("FAIL")
	}

	points := []V2{
		{0, -1},
		{1, 1},
		{-1, 1},
	}

	s = Polygon2D(points)

	b = s.BoundingBox()
	b0 = Box2{V2{-1, -1}, V2{1, 1}}
	if b.Equals(b0, tolerance) == false {
		t.Error("FAIL")
	}

	testPoints := []struct {
		p V2
		d float64
	}{
		{V2{0, -1}, 0},
		{V2{1, 1}, 0},
		{V2{-1, 1}, 0},
		{V2{0, 1}, 0},
		{V2{0, 2}, 1},
		{V2{0, -2}, 1},
		{V2{1, 0}, 1 / math.Sqrt(5)},
		{V2{-1, 0}, 1 / math.Sqrt(5)},
		{V2{0, 0}, -1 / math.Sqrt(5)},
		{V2{3, 0}, math.Sqrt(5)},
		{V2{-3, 0}, math.Sqrt(5)},
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

	s0 := Polygon2D([]V2{{k, -k}, {k, k}, {-k, k}, {-k, -k}})
	s0 = Transform2D(s0, Translate2d(V2{0.8, 0}))

	s1 := Box2D(V2{2 * k, 2 * k}, 0)
	s1 = Transform2D(s1, Translate2d(V2{0.8, 0}))

	for i := 0; i < 10000; i++ {
		b := NewBox2(V2{0, 0}, V2{20 * k, 20 * k})
		p := b.Random()
		if Abs(s0.Evaluate(p)-s1.Evaluate(p)) > tolerance {
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

	s1 := Box2D(V2{2 * a, 2 * b}, 0)
	s1 = Transform2D(s1, Rotate2d(DtoR(theta)))
	s1 = Transform2D(s1, Translate2d(V2{j, k}))

	points := []V2{
		{j + c*a - s*b, k + s*a + c*b},
		{j - c*a - s*b, k - s*a + c*b},
		{j - c*a + s*b, k - s*a - c*b},
		{j + c*a + s*b, k + s*a - c*b},
	}

	s0 := Polygon2D(points)

	for i := 0; i < 1000; i++ {
		b := NewBox2(V2{0, 0}, V2{10 * b, 10 * b})
		p := b.Random()
		if Abs(s0.Evaluate(p)-s1.Evaluate(p)) > tolerance {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_ArraySDF2(t *testing.T) {
	r := 0.5
	s := Circle2D(r)
	bb := s.BoundingBox()
	if bb.Min.Equals(V2{-r, -r}, tolerance) == false {
		t.Error("FAIL")
	}
	if bb.Max.Equals(V2{r, r}, tolerance) == false {
		t.Error("FAIL")
	}

	j := 3
	k := 4
	dx := 2.0
	dy := 7.0
	sa := Array2D(s, V2i{j, k}, V2{dx, dy})
	saBox := sa.BoundingBox()
	if saBox.Min.Equals(V2{-r, -r}, tolerance) == false {
		t.Error("FAIL")
	}
	if saBox.Max.Equals(V2{r + (float64(j-1) * dx), r + (float64(k-1) * dy)}, tolerance) == false {
		t.Error("FAIL")
	}

	j = 7
	k = 4
	dx = -3.0
	dy = -5.0
	sa = Array2D(s, V2i{j, k}, V2{dx, dy})
	saBox = sa.BoundingBox()
	if saBox.Min.Equals(V2{-r + (float64(j-1) * dx), -r + (float64(k-1) * dy)}, tolerance) == false {
		t.Error("FAIL")
	}
	if saBox.Max.Equals(V2{r, r}, tolerance) == false {
		t.Error("FAIL")
	}

	j = 6
	k = 8
	dx = 5.0
	dy = -3.0
	sa = Array2D(s, V2i{j, k}, V2{dx, dy})
	saBox = sa.BoundingBox()
	if saBox.Min.Equals(V2{-r, -r + (float64(k-1) * dy)}, tolerance) == false {
		t.Error("FAIL")
	}
	if saBox.Max.Equals(V2{r + (float64(j-1) * dx), r}, tolerance) == false {
		t.Error("FAIL")
	}

	j = 9
	k = 1
	dx = -0.5
	dy = 6.5
	sa = Array2D(s, V2i{j, k}, V2{dx, dy})
	saBox = sa.BoundingBox()
	if saBox.Min.Equals(V2{-r + (float64(j-1) * dx), -r}, tolerance) == false {
		t.Error("FAIL")
	}
	if saBox.Max.Equals(V2{r, r + (float64(k-1) * dy)}, tolerance) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_Rotation2d(t *testing.T) {
	r := Rotate2d(DtoR(90))
	v := V2{1, 0}
	v = r.MulPosition(v)
	if v.Equals(V2{0, 1}, tolerance) == false {
		t.Error("FAIL")
	}
}

func Test_Rotation3d(t *testing.T) {
	r := Rotate3d(V3{0, 0, 1}, DtoR(90))
	v := V3{1, 0, 0}
	v = r.MulPosition(v)
	if v.Equals(V3{0, 1, 0}, tolerance) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_TriDiagonal(t *testing.T) {
	n := 5
	m := make([]V3, n)
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
	x := TriDiagonal(m, d)
	for i := 0; i < n; i++ {
		if Abs(x[i]-x0[i]) > tolerance {
			t.Error("FAIL")
		}
	}

	d = []float64{10, 20, 30, 40, 50}
	x0 = []float64{15.0 / 4.0, 5.0 / 2.0, 25.0 / 4.0, 5.0 / 2.0, 95.0 / 4.0}
	x = TriDiagonal(m, d)
	for i := 0; i < n; i++ {
		if Abs(x[i]-x0[i]) > tolerance {
			t.Error("FAIL")
		}
	}

	m[0] = V3{0, 1, 2}
	m[1] = V3{3, 4, 5}
	m[2] = V3{6, 7, 8}
	m[3] = V3{9, 10, 11}
	m[4] = V3{12, 13, 0}
	d = []float64{-10, -20, -30, 40, 50}
	x0 = []float64{60.0 / 49.0, -275.0 / 49.0, -12.0 / 49.0, 33.0 / 49.0, 158.0 / 49.0}
	x = TriDiagonal(m, d)
	for i := 0; i < n; i++ {
		if Abs(x[i]-x0[i]) > tolerance {
			t.Error("FAIL")
		}
	}

}

//-----------------------------------------------------------------------------

func Test_CubicSpline(t *testing.T) {

	knot := []V2{
		{-1.5, -1.2},
		{-0.2, 0},
		{1, 0.5},
		{5, 1},
		{10, 2.2},
		{12, 3.2},
		{-16, -1.2},
		{-18, -3.2},
	}
	s := CubicSpline2D(knot).(*CubicSplineSDF2)
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
			if !f2.Equals(V2{0, 0}, tolerance) {
				t.Error("FAIL")
			}
		}
		if i == n-1 {
			// 2nd derivative at end == 0
			f2 := cs.f2(1)
			if !f2.Equals(V2{0, 0}, tolerance) {
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
		a, b, c V2
		result  bool
	}{
		{V2{}, V2{}, V2{}, true},
		{V2{}, V2{}, V2{10, 17}, true},
		{V2{1, 1}, V2{1, 2}, V2{1, 3}, true},
		{V2{1, 1}, V2{2, 2}, V2{3, 3}, true},
		{V2{1, 1}, V2{2, 1}, V2{3, 1}, true},
		{V2{1, 1}, V2{2, 5}, V2{-2, 4}, false},
		{V2{1, 1}, V2{1, 1}, V2{-2, 4}, true},
		{V2{1, 1}, V2{-1, 1}, V2{0, -1}, false},
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
	b0 := NewBox2(V2{0, 0}, V2{10, 10})
	b1 := NewBox2(V2{10, 20}, V2{30, 40})
	tests := []struct {
		b      Box2
		p      V2
		result V2
	}{
		{b0, V2{0, 0}, V2{0, 50}},
		{b0, V2{5, 5}, V2{0, 200}},
		{b0, V2{20, 0}, V2{225, 650}},
		{b1, V2{0, 0}, V2{0, 2225}},
		{b1, V2{10, 20}, V2{0, 625}},
		{b1, V2{0, -10}, V2{100, 3125}},
		{b1, V2{0, 5}, V2{0, 1850}},
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
