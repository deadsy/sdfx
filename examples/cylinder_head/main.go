//-----------------------------------------------------------------------------
/*

Wallaby Cylinder Head

No draft version for 3d printing and lost-PLA investment casting.

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// overall build controls
const casting = false // add allowances, remove machined features

//-----------------------------------------------------------------------------
// scaling

const MM_PER_INCH = 25.4
const desired_scale = 1.25
const al_shrink = 1.0 / 0.99   // ~1%
const pla_shrink = 1.0 / 0.998 //~0.2%
const abs_shrink = 1.0 / 0.995 //~0.5%

// dimension scaling
func dim(x float64) float64 {
	return x * desired_scale * MM_PER_INCH * al_shrink * pla_shrink
}

var general_round = dim(0.125)

//-----------------------------------------------------------------------------
// cylinder domes (or full base)

var cylinder_height = dim(3.0 / 16.0)
var cylinder_diameter = dim(1.0 + (1.0 / 8.0))
var cylinder_wall = dim(1.0 / 4.0)
var cylinder_radius = cylinder_diameter / 2.0

var dome_radius = cylinder_wall + cylinder_radius
var dome_height = cylinder_wall + cylinder_height
var dome_draft = DtoR(12)

var c2c_distance = dim(1.0 + (3.0 / 8.0))

func cylinder_head() SDF3 {
	delta := dome_height * math.Tan(dome_draft)
	// build the cylinder dome cross section
	s := NewSmoother(false)
	s.Add(V2{0, 0})
	s.Add(V2{dome_radius + delta, 0})
	s.AddSmooth(V2{dome_radius, dome_height}, 4, general_round)
	s.Add(V2{0, dome_height})
	s.Smooth()
	return NewSorSDF3(NewPolySDF2(s.Vertices()))
}

func cylinder_chamber() SDF3 {
	return NewCylinderSDF3(cylinder_height, cylinder_radius, 0.0)
}

func cylinder_heads() SDF3 {
	s0 := cylinder_head()
	d := c2c_distance / 2
	s1 := NewTransformSDF3(s0, Translate3d(V3{0, d, 0}))
	s2 := NewTransformSDF3(s0, Translate3d(V3{0, -d, 0}))
	s := NewUnionSDF3(s1, s2)
	s.(*UnionSDF3).SetMin(PolyMin, general_round)
	return s
}

//-----------------------------------------------------------------------------
// cylinder studs: location, bosses and holes

var stud_hole_radius = dim(1.0 / 16.0)
var stud_boss_radius = dim(3.0 / 16.0)
var stud_hole_dy = dim(11.0 / 16.0)
var stud_hole_dx0 = dim(7.0 / 16.0)
var stud_hole_dx1 = dim(1.066)

var stud_locations = []V2{
	V2{stud_hole_dx0 + stud_hole_dx1, 0},
	V2{stud_hole_dx0 + stud_hole_dx1, stud_hole_dy},
	V2{stud_hole_dx0 + stud_hole_dx1, -stud_hole_dy},
	V2{stud_hole_dx0, stud_hole_dy},
	V2{stud_hole_dx0, -stud_hole_dy},
	V2{-stud_hole_dx0 - stud_hole_dx1, 0},
	V2{-stud_hole_dx0 - stud_hole_dx1, stud_hole_dy},
	V2{-stud_hole_dx0 - stud_hole_dx1, -stud_hole_dy},
	V2{-stud_hole_dx0, stud_hole_dy},
	V2{-stud_hole_dx0, -stud_hole_dy},
}

func head_stud_holes() SDF3 {
	s := NewMultiCircleSDF2(stud_hole_radius, stud_locations)
	return NewExtrudeSDF3(s, head_height)
}

//-----------------------------------------------------------------------------
// head walls

var head_length = dim(4.30 / 1.25)
var head_width = dim(2.33 / 1.25)
var head_height = dim(7.0 / 8.0)
var head_corner_round = dim((5.0 / 32.0) / 1.25)
var head_internal_round = dim(0.125)
var head_wall_thickness = dim(0.154)

func head_wall_outer_2d() SDF2 {
	return NewBoxSDF2(V2{head_length, head_width}, head_corner_round)
}

func head_wall_inner_2d() SDF2 {
	l := head_length - (2 * head_wall_thickness)
	w := head_width - (2 * head_wall_thickness)
	s0 := NewBoxSDF2(V2{l, w}, 0)
	s1 := NewMultiCircleSDF2(stud_boss_radius, stud_locations)
	s := NewDifferenceSDF2(s0, s1)
	s.(*DifferenceSDF2).SetMax(PolyMax, head_internal_round)
	return s
}

func head_envelope() SDF3 {
	return NewExtrudeSDF3(head_wall_outer_2d(), head_height)
}

func head_wall() SDF3 {
	s := head_wall_outer_2d()
	s = NewDifferenceSDF2(s, head_wall_inner_2d())
	return NewExtrudeSDF3(s, head_height)
}

func head_base() SDF3 {
	delta := head_height - dome_height
	s := NewExtrudeSDF3(head_wall_inner_2d(), dome_height)
	m := Translate3d(V3{0, 0, -delta / 2})
	s = NewTransformSDF3(s, m)
	return s
}

//-----------------------------------------------------------------------------

func additive() SDF3 {
	s := NewUnionSDF3(head_wall(), head_base())
	s.(*UnionSDF3).SetMin(PolyMin, head_internal_round)
	// cleanup the blending artifacts on the outside
	s = NewIntersectionSDF3(s, head_envelope())

	return s
}

//-----------------------------------------------------------------------------

func subtractive(s SDF3) SDF3 {
	s = NewDifferenceSDF3(s, head_stud_holes())
	return s
}

//-----------------------------------------------------------------------------

func main() {
	s := additive()
	if !casting {
		s = subtractive(s)
	}
	RenderSTL(s, "head.stl")
	//RenderSTL(cylinder_heads(), "head2.stl")

}

//-----------------------------------------------------------------------------
