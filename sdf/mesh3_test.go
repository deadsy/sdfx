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
	}

	for i, test := range testSet {
		if len(test.p) != len(test.d2) {
			t.Errorf("test %d: len(p) != len(d2)", i)
		}
		ti := newTriangleInfo(&test.t)
		for j, p := range test.p {
			d2 := ti.minDistance2(p)
			if !EqualFloat64(d2, test.d2[j], tolerance) {
				t.Errorf("test %d.%d: expected %f, got %f", i, j, test.d2[j], d2)
			}
		}
	}
}

//-----------------------------------------------------------------------------
