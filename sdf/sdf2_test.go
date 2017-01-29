//-----------------------------------------------------------------------------
/*

SDF2 Testing

*/
//-----------------------------------------------------------------------------

package sdf

import "testing"

func Test_ArraySDF2(t *testing.T) {
	r := 0.5
	s := NewCircleSDF2(r)
	bb := s.BoundingBox()
	if bb.Min.Equals(V2{-r, -r}, TOLERANCE) == false {
		t.Error("FAIL")
	}
	if bb.Max.Equals(V2{r, r}, TOLERANCE) == false {
		t.Error("FAIL")
	}

	j := 3
	k := 4
	dx := 2.0
	dy := 7.0
	sa := NewArraySDF2(s, V2i{j, k}, V2{dx, dy})
	sa_bb := sa.BoundingBox()
	if sa_bb.Min.Equals(V2{-r, -r}, TOLERANCE) == false {
		t.Error("FAIL")
	}
	if sa_bb.Max.Equals(V2{r + (float64(j-1) * dx), r + (float64(k-1) * dy)}, TOLERANCE) == false {
		t.Error("FAIL")
	}

	j = 7
	k = 4
	dx = -3.0
	dy = -5.0
	sa = NewArraySDF2(s, V2i{j, k}, V2{dx, dy})
	sa_bb = sa.BoundingBox()
	if sa_bb.Min.Equals(V2{-r + (float64(j-1) * dx), -r + (float64(k-1) * dy)}, TOLERANCE) == false {
		t.Error("FAIL")
	}
	if sa_bb.Max.Equals(V2{r, r}, TOLERANCE) == false {
		t.Error("FAIL")
	}

	j = 6
	k = 8
	dx = 5.0
	dy = -3.0
	sa = NewArraySDF2(s, V2i{j, k}, V2{dx, dy})
	sa_bb = sa.BoundingBox()
	if sa_bb.Min.Equals(V2{-r, -r + (float64(k-1) * dy)}, TOLERANCE) == false {
		t.Error("FAIL")
	}
	if sa_bb.Max.Equals(V2{r + (float64(j-1) * dx), r}, TOLERANCE) == false {
		t.Error("FAIL")
	}

	j = 9
	k = 1
	dx = -0.5
	dy = 6.5
	sa = NewArraySDF2(s, V2i{j, k}, V2{dx, dy})
	sa_bb = sa.BoundingBox()
	if sa_bb.Min.Equals(V2{-r + (float64(j-1) * dx), -r}, TOLERANCE) == false {
		t.Error("FAIL")
	}
	if sa_bb.Max.Equals(V2{r, r + (float64(k-1) * dy)}, TOLERANCE) == false {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------
