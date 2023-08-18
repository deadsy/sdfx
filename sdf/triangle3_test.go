//-----------------------------------------------------------------------------
/*

3D Triangle Testing

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"testing"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func Test_Triangle3_rotateToXY(t *testing.T) {

	testSet := []struct {
		t      *Triangle3
		result v3.Vec // {0,0,0}, {a,0,0}, {b,c,0}
	}{
		{
			&Triangle3{{1, 1, 1}, {2, 2, 2}, {3, 4, 5}},
			v3.Vec{1.7320508075688776, 5.196152422706633, 1.4142135623730954},
		},
		{
			&Triangle3{{1, 1, 1}, {2, 2, 2}, {5, 4, 3}},
			v3.Vec{1.7320508075688776, 5.196152422706633, 1.4142135623730954},
		},
		{
			&Triangle3{{0, 0, 5}, {1, 0, 5}, {4, 8, 5}},
			v3.Vec{1, 4, 8},
		},
		{
			&Triangle3{{0, 0, 0}, {1, 0, 0}, {4, -8, 0}},
			v3.Vec{1, 4, 8},
		},
		{
			&Triangle3{{0, 0, -3}, {10, 0, -3}, {4, 7, -3}},
			v3.Vec{10, 4, 7},
		},
		{
			&Triangle3{{0, -3, 0}, {10, -3, 0}, {4, -3, 7}},
			v3.Vec{10, 4, 7},
		},
		{
			&Triangle3{{-8.61, 1.80, 19.31}, {-5.99, -0.72, 21.51}, {-8.88, 5.05, 20.67}},
			v3.Vec{4.249094021082613, -1.389802148575515, 3.2486073920704683},
		},
		{
			&Triangle3{{0, 0, -1}, {10, 0, -1}, {-10, -8, -1}},
			v3.Vec{10, -10, 8},
		},
		{
			&Triangle3{{0, 0, 3}, {10, 0, 3}, {-10, 8, 3}},
			v3.Vec{10, -10, 8},
		},
	}

	for i, test := range testSet {

		m := test.t.rotateToXY()
		x0 := m.MulPosition(test.t[0])
		x1 := m.MulPosition(test.t[1])
		x2 := m.MulPosition(test.t[2])

		r := test.result
		r0 := v3.Vec{0, 0, 0}
		r1 := v3.Vec{r.X, 0, 0}
		r2 := v3.Vec{r.Y, r.Z, 0}

		if !x0.Equals(r0, tolerance) {
			t.Errorf("test %d: expected %v, got %v", i, r0, x0)
		}
		if !x1.Equals(r1, tolerance) {
			t.Errorf("test %d: expected %v, got %v", i, r1, x1)
		}
		if !x2.Equals(r2, tolerance) {
			t.Errorf("test %d: expected %v, got %v", i, r2, x2)
		}
	}

	// sanity test with random triangles
	b := NewBox3(v3.Vec{-1, 3, -7}, v3.Vec{7, 5, 3})
	for i := 0; i < 10000; i++ {
		x := b.RandomTriangle()
		m := x.rotateToXY()

		x0 := m.MulPosition(x[0])
		x1 := m.MulPosition(x[1])
		x2 := m.MulPosition(x[2])

		// x0 should be {0,0,0}
		if !x0.Equals(v3.Vec{0, 0, 0}, tolerance) {
			t.Errorf("test %d: expected {0,0,0}, got %v", i, x0)
		}
		// x1 should be {a,0,0}, where a > 0
		if x1.X <= 0 || !EqualFloat64(x1.Y, 0, tolerance) || !EqualFloat64(x1.Z, 0, tolerance) {
			t.Errorf("test %d: expected {a,0,0}, got %v", i, x1)
		}
		// x2 should be {b,c,0}, where c > 0
		if x2.Y <= 0 || !EqualFloat64(x2.Z, 0, tolerance) {
			t.Errorf("test %d: expected {b,c,0}, got %v", i, x2)
		}
	}

}

//-----------------------------------------------------------------------------
