//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
	"testing"
)

//-----------------------------------------------------------------------------

const TOLERANCE = 1e-9

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
	a_inv := M33{3, -1, -1, -4, 2, 1, -1, 0, 1}
	if a.Inverse().Equals(a_inv, TOLERANCE) == false {
		t.Error("FAIL")
	}
	if a.Mul(a_inv).Equals(Identity2d(), TOLERANCE) == false {
		t.Error("FAIL")
	}

	for i := 0; i < 100; i++ {
		a = RandomM33(-5, 5)
		a_inv = a.Inverse()
		if a.Mul(a_inv).Equals(Identity2d(), TOLERANCE) == false {
			t.Error("FAIL")
		}
	}

	for i := 0; i < 100; i++ {
		b := RandomM44(-1, 1)
		b_inv := b.Inverse()
		if b.Mul(b_inv).Equals(Identity3d(), TOLERANCE) == false {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_MulBox(t *testing.T) {

	// 2D boxes
	b2d := Box2{V2{-1, -1}, V2{1, 1}}
	for i := 0; i < 100; i++ {
		v := RandomV2(-5, 5)
		// translating
		m0 := Translate2d(v)
		m1 := Translate2d(v.Negate())
		b1 := m0.MulBox(b2d)
		b2 := m1.MulBox(b1)
		if b2d.Equals(b2, TOLERANCE) == false {
			t.Error("FAIL")
		}
		// scaling
		m0 = Scale2d(v)
		m1 = Scale2d(V2{1 / v.X, 1 / v.Y})
		b1 = m0.MulBox(b2d)
		b2 = m1.MulBox(b1)
		if b2d.Equals(b2, TOLERANCE) == false {
			t.Error("FAIL")
		}
	}

	// 3D boxes
	b3d := Box3{V3{-1, -1, -1}, V3{1, 1, 1}}
	for i := 0; i < 100; i++ {
		v := RandomV3(-5, 5)
		// translating
		m0 := Translate3d(v)
		m1 := Translate3d(v.Negate())
		b1 := m0.MulBox(b3d)
		b2 := m1.MulBox(b1)
		if b3d.Equals(b2, TOLERANCE) == false {
			t.Error("FAIL")
		}
		// scaling
		m0 = Scale3d(v)
		m1 = Scale3d(V3{1 / v.X, 1 / v.Y, 1 / v.Z})
		b1 = m0.MulBox(b3d)
		b2 = m1.MulBox(b1)
		if b3d.Equals(b2, TOLERANCE) == false {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_ScaleBox(t *testing.T) {
	b0 := Box3{V3{-1, -1, -1}, V3{1, 1, 1}}

	b1 := Box3{V3{-2, -2, -2}, V3{2, 2, 2}}
	b2 := b0.Scale(2)
	if b1.Equals(b2, TOLERANCE) == false {
		t.Error("FAIL")
	}

	b1 = Box3{V3{-1.1, -1.1, -1.1}, V3{1.1, 1.1, 1.1}}
	b2 = b0.Scale(1.1)
	if b1.Equals(b2, TOLERANCE) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

func Test_Line(t *testing.T) {
	l := NewLine2(&V2{0, 1}, &V2{0, 2})
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
		d := l.Distance(&p.p)
		if d != p.d {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Polygon1(t *testing.T) {
	s := NewPolySDF2([]*V2{&V2{0, 0}, &V2{1, 0}, &V2{0, 1}})
	b := s.BoundingBox()
	b0 := Box2{V2{0, 0}, V2{1, 1}}
	if b.Equals(b0, TOLERANCE) == false {
		t.Error("FAIL")
	}

	s = NewPolySDF2([]*V2{&V2{0, -2}, &V2{1, 1}, &V2{-2, 2}})
	b = s.BoundingBox()
	b0 = Box2{V2{-2, -2}, V2{1, 2}}
	if b.Equals(b0, TOLERANCE) == false {
		t.Error("FAIL")
	}

	points := []*V2{
		&V2{0, -1},
		&V2{1, 1},
		&V2{-1, 1},
	}

	s = NewPolySDF2(points)
	b = s.BoundingBox()
	b0 = Box2{V2{-1, -1}, V2{1, 1}}
	if b.Equals(b0, TOLERANCE) == false {
		t.Error("FAIL")
	}

	test_points := []struct {
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
	}

	for _, p := range test_points {
		d := s.Evaluate(p.p)
		if d != p.d {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Polygon2(t *testing.T) {
	k := 1.2

	s0 := NewPolySDF2([]*V2{&V2{k, -k}, &V2{k, k}, &V2{-k, k}, &V2{-k, -k}})
	s0 = NewTransformSDF2(s0, Translate2d(V2{0.8, 0}))

	s1 := NewBoxSDF2(V2{2 * k, 2 * k})
	s1 = NewTransformSDF2(s1, Translate2d(V2{0.8, 0}))

	for i := 0; i < 10000; i++ {
		p := RandomV2(-10*k, 10*k)
		if Abs(s0.Evaluate(p)-s1.Evaluate(p)) > TOLERANCE {
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

	s1 := NewBoxSDF2(V2{2 * a, 2 * b})
	s1 = NewTransformSDF2(s1, Rotate2d(DtoR(theta)))
	s1 = NewTransformSDF2(s1, Translate2d(V2{j, k}))

	points := []*V2{
		&V2{j + c*a - s*b, k + s*a + c*b},
		&V2{j - c*a - s*b, k - s*a + c*b},
		&V2{j - c*a + s*b, k - s*a - c*b},
		&V2{j + c*a + s*b, k + s*a - c*b},
	}

	s0 := NewPolySDF2(points)

	for i := 0; i < 1000; i++ {
		p := RandomV2(-5*b, 5*b)
		if Abs(s0.Evaluate(p)-s1.Evaluate(p)) > TOLERANCE {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func Test_KdTree2(t *testing.T) {

	points := []V2{
		V2{2, 3},
		V2{5, 4},
		V2{9, 6},
		V2{4, 7},
		V2{8, 1},
		V2{7, 2},
	}

	fmt.Printf("before %+v\n", points)
	k := NewKdTree2(points)
	fmt.Printf("%+v\n", k)
	fmt.Printf("after %+v\n", points)
}

//-----------------------------------------------------------------------------
