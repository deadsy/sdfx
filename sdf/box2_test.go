//-----------------------------------------------------------------------------
/*

2D Box Testing

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"testing"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

func Test_Box2_lineIntersect(t *testing.T) {

	testSet := []struct {
		box    Box2
		line   *Line2
		result *Line2
	}{

		// inside the box, no intersect
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{0, 0}, {1, 1}}, &Line2{{0, 0}, {1, 1}}},
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-5, -7}, {7, 5}}, &Line2{{-5, -7}, {7, 5}}},

		// inside the box, intersects
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{0, 0}, {10, 10}}, &Line2{{0, 0}, {10, 10}}},
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{0, 0}, {20, 20}}, &Line2{{0, 0}, {10, 10}}},
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{7, 2}, {13, 8}}, &Line2{{7, 2}, {10, 5}}},
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{1, -7}, {9, -13}}, &Line2{{1, -7}, {5, -10}}},

		// outside the box, no intersect
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{20, 20}, {30, 30}}, nil},
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-12, 3}, {-8, 17}}, nil},

		// outside the box, intersects
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-20, 1}, {20, 1}}, &Line2{{-10, 1}, {10, 1}}},
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-1, -11}, {-1, 15}}, &Line2{{-1, -10}, {-1, 10}}},
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-15, 0}, {5, 20}}, &Line2{{-10, 5}, {-5, 10}}},

		// horizontal bottom edge
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-20, -10}, {5, -10}}, &Line2{{-10, -10}, {5, -10}}},
		// horizontal top edge
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-20, 10}, {5, 10}}, nil},
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-20, 10}, {20, 10}}, nil},

		// vertical left edge
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{-10, -20}, {-10, 5}}, &Line2{{-10, -10}, {-10, 5}}},
		// vertical right edge
		{Box2{v2.Vec{-10, -10}, v2.Vec{10, 10}}, &Line2{{10, -20}, {10, 5}}, nil},
	}

	for i, test := range testSet {

		box := test.box
		line := test.line
		result := box.lineIntersect(line)

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

		if !test.result.Equals(result, tolerance) {
			t.Errorf("test %d: expected %v, got %v", i, test.result, result)
			continue
		}

	}

}

//-----------------------------------------------------------------------------
