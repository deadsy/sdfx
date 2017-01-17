//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"testing"
)

//-----------------------------------------------------------------------------

const TOLERANCE = 0.0000001

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
	if a.Mul(a_inv).Equals(IdentityM33(), TOLERANCE) == false {
		t.Error("FAIL")
	}

	for i := 0; i < 100; i++ {
		a = RandomM33(-5, 5)
		a_inv = a.Inverse()
		if a.Mul(a_inv).Equals(IdentityM33(), TOLERANCE) == false {
			t.Error("FAIL")
		}
	}

	for i := 0; i < 100; i++ {
		b := RandomM44(-1, 1)
		b_inv := b.Inverse()
		if b.Mul(b_inv).Equals(IdentityM44(), TOLERANCE) == false {
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
		m0 := TranslateM33(v)
		m1 := TranslateM33(v.Negate())
		b1 := m0.MulBox(b2d)
		b2 := m1.MulBox(b1)
		if b2d.Equals(b2, TOLERANCE) == false {
			t.Error("FAIL")
		}
		// scaling
		m0 = ScaleM33(v)
		m1 = ScaleM33(V2{1 / v.X, 1 / v.Y})
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
		m0 := TranslateM44(v)
		m1 := TranslateM44(v.Negate())
		b1 := m0.MulBox(b3d)
		b2 := m1.MulBox(b1)
		if b3d.Equals(b2, TOLERANCE) == false {
			t.Error("FAIL")
		}
		// scaling
		m0 = ScaleM44(v)
		m1 = ScaleM44(V3{1 / v.X, 1 / v.Y, 1 / v.Z})
		b1 = m0.MulBox(b3d)
		b2 = m1.MulBox(b1)
		if b3d.Equals(b2, TOLERANCE) == false {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------
