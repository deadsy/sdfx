//-----------------------------------------------------------------------------
/*

Line Testing

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"sort"
	"testing"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

func Test_Interval(t *testing.T) {

	testSet := []struct {
		a, b   Interval  // intervals
		result *Interval // intersection
	}{
		// non-intersecting
		{Interval{0, 1}, Interval{2, 3}, nil},
		{Interval{1, 1}, Interval{-5, -5}, nil},

		// intersecting
		{Interval{0, 1}, Interval{0, 1}, &Interval{0, 1}},
		{Interval{0, 1}, Interval{1, 2}, &Interval{1, 1}},
		{Interval{-5, 7}, Interval{-2, 10}, &Interval{-2, 7}},
		{Interval{1, 1}, Interval{1, 1}, &Interval{1, 1}},
	}

	for i, test := range testSet {

		a := test.a.Sort()
		b := test.b.Sort()
		r0 := a.Intersect(b)
		r1 := b.Intersect(a)

		for _, result := range []*Interval{r0, r1} {
			if test.result == nil {
				if result != nil {
					t.Errorf("test %d: expected nil intersection, got %v", i, result)
				}
				continue
			}
			if result == nil && test.result != nil {
				t.Errorf("test %d: expected %v intersection, got nil", i, test.result)
				continue
			}
			for j := range test.result {
				if test.result[j] != result[j] {
					t.Errorf("test %d: expected %v (got %v)", i, test.result, result)
					break
				}
			}
		}
	}
}

//-----------------------------------------------------------------------------

func Test_IntersectLine(t *testing.T) {

	testSet := []struct {
		a, b  *Line2   // line segments
		point []v2.Vec // intersection point(s)
	}{
		// parallel, non-intersecting
		{&Line2{v2.Vec{0, 0}, v2.Vec{10, 10}}, &Line2{v2.Vec{20, 0}, v2.Vec{30, 10}}, nil},

		// collinear, intersecting
		{&Line2{v2.Vec{0, 0}, v2.Vec{10, 0}}, &Line2{v2.Vec{-5, 0}, v2.Vec{5, 0}}, []v2.Vec{{0, 0}, {5, 0}}},
		{&Line2{v2.Vec{0, 0}, v2.Vec{10, 0}}, &Line2{v2.Vec{-5, 0}, v2.Vec{0, 0}}, []v2.Vec{{0, 0}}},
		{&Line2{v2.Vec{0, 5}, v2.Vec{0, 15}}, &Line2{v2.Vec{0, 7}, v2.Vec{0, 18}}, []v2.Vec{{0, 7}, {0, 15}}},

		// collinear, non-intersecting
		{&Line2{v2.Vec{5, 0}, v2.Vec{10, 0}}, &Line2{v2.Vec{-5, 0}, v2.Vec{0, 0}}, nil},
		{&Line2{v2.Vec{0, 5}, v2.Vec{0, 10}}, &Line2{v2.Vec{0, 1}, v2.Vec{0, 2}}, nil},
		{&Line2{v2.Vec{5, 5}, v2.Vec{10, 10}}, &Line2{v2.Vec{-10, -10}, v2.Vec{0, 0}}, nil},
		{&Line2{v2.Vec{-5, 5}, v2.Vec{-10, 10}}, &Line2{v2.Vec{10, -10}, v2.Vec{0, 0}}, nil},

		// non-parallel, intersecting
		{&Line2{v2.Vec{-1, 0}, v2.Vec{1, 0}}, &Line2{v2.Vec{0, -1}, v2.Vec{0, 1}}, []v2.Vec{{0, 0}}},
		{&Line2{v2.Vec{0, 10}, v2.Vec{0, 0}}, &Line2{v2.Vec{20, 10}, v2.Vec{0, 10}}, []v2.Vec{{0, 10}}},
		{&Line2{v2.Vec{0, 1}, v2.Vec{15, 17}}, &Line2{v2.Vec{15, 17}, v2.Vec{-123, 47}}, []v2.Vec{{15, 17}}},
		{&Line2{v2.Vec{0, 1}, v2.Vec{30, 33}}, &Line2{v2.Vec{15, 17}, v2.Vec{-123, 47}}, []v2.Vec{{15, 17}}},

		// non-parallel, non-intersecting
		{&Line2{v2.Vec{5, 5}, v2.Vec{10, 10}}, &Line2{v2.Vec{-10, 10}, v2.Vec{-1, 1}}, nil},
		{&Line2{v2.Vec{0, 1}, v2.Vec{14, 17}}, &Line2{v2.Vec{15, 17}, v2.Vec{-123, 47}}, nil},
	}

	for i, test := range testSet {

		a := test.a
		ra := a.Reverse()
		b := test.b
		rb := b.Reverse()

		r0 := a.IntersectLine(b)
		r1 := a.IntersectLine(rb)
		r2 := ra.IntersectLine(b)
		r3 := ra.IntersectLine(rb)

		sort.Sort(v2.VecSetByXY(r0))
		sort.Sort(v2.VecSetByXY(r1))
		sort.Sort(v2.VecSetByXY(r2))
		sort.Sort(v2.VecSetByXY(r3))

		for _, result := range [][]v2.Vec{r0, r1, r2, r3} {
			if len(result) != len(test.point) {
				t.Errorf("test %d: expected %d intersection(s) (got %d)", i, len(test.point), len(result))
				continue
			}
			for j := range test.point {
				if !result[j].Equals(test.point[j], tolerance) {
					t.Errorf("test %d: expected %v (got %v)", i, test.point[j], result[j])
				}
			}
		}

	}

}

//-----------------------------------------------------------------------------
