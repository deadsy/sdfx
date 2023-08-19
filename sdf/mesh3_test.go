//-----------------------------------------------------------------------------
/*

Mesh 3D Testing and Benchmarking

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"testing"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func Test_Mesh3_minDistance2(t *testing.T) {

	testSet := []struct {
		t  Triangle3
		p  []v3.Vec
		d2 []float64
	}{
		{
			Triangle3{{1, 2, 1}, {-4, -5, 1}, {17, -3, 1}},
			[]v3.Vec{{1, 2, 3}, {-4, -5, 6}, {17, -3, -2}},
			[]float64{4, 25, 9},
		},
		{
			Triangle3{{10, 0, 10}, {0, 0, -10}, {-10, 0, 10}},
			[]v3.Vec{{0, 4, 0}, {0, 0, 0}, {11, 4, 11}, {0, 3, -11}, {-11, 7, 11}},
			[]float64{16, 0, 18, 10, 51},
		},
		{
			Triangle3{{0, 0, 4}, {0, 6, -2}, {0, -6, -2}},
			[]v3.Vec{{0, 0, 5}, {0, 2, 2}, {4, 3, 3}, {3, -3, 3}, {-2, 0, -3}},
			[]float64{1, 0, 18, 11, 5},
		},
	}

	for i, test := range testSet {
		if len(test.p) != len(test.d2) {
			t.Errorf("test %d: len(p) != len(d2)", i)
		}
		triangle := test.t
		for j := 0; j < 3; j++ {
			triangle = triangle.rotateVertex()
			ti := newTriangleInfo(&triangle)
			for k, p := range test.p {
				d2 := ti.minDistance2(p)
				if !EqualFloat64(d2, test.d2[k], tolerance) {
					t.Errorf("test %d.%d: expected %f, got %f", i, k, test.d2[k], d2)
				}
			}
		}
	}

	// sanity test with random triangles
	const boxSize = 100.0
	const d2Max = 3.0 * (boxSize * boxSize)
	b := NewBox3(v3.Vec{0, 0, 0}, v3.Vec{100, 100, 100})
	for i := 0; i < 10000; i++ {
		x := b.RandomTriangle()
		ti := newTriangleInfo(&x)
		p := b.Random()
		d2 := ti.minDistance2(p)
		if d2 < 0 {
			t.Errorf("test %d: expected >= 0, got %f", i, d2)
		}
		if d2 > d2Max {
			t.Errorf("test %d: expected <= %f, got %f", i, d2Max, d2)
		}
	}

}

//-----------------------------------------------------------------------------
