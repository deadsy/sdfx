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

func Test_Inverse(t *testing.T) {
	a := M33{2, 1, 1, 3, 2, 1, 2, 1, 2}
	a_inv := M33{3, -1, -1, -4, 2, 1, -1, 0, 1}
	if a.Inverse().Equals(a_inv, TOLERANCE) == false {
		t.Error("FAIL")
	}
	if a.Mul(a_inv).Equals(IdentityM33(), TOLERANCE) == false {
		t.Error("FAIL")
	}

	a = RandomM33(-5, 5)
	a_inv = a.Inverse()
	if a.Mul(a_inv).Equals(IdentityM33(), TOLERANCE) == false {
		t.Error("FAIL")
	}

	b := RandomM44(-1, 1)
	b_inv := b.Inverse()
	if b.Mul(b_inv).Equals(IdentityM44(), TOLERANCE) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------
