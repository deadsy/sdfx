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

var general_round = dim(0.1)

//-----------------------------------------------------------------------------
// exhaust bosses

var eb_side_radius = dim(5.0 / 32.0)
var eb_main_radius = dim(5.0 / 16.0)
var eb_hole_radius = dim(3.0 / 16.0)
var eb_c2c_distance = dim(13.0 / 16.0)
var eb_distance = eb_c2c_distance / 2.0

var eb_y_offset = (head_width / 2.0) - eb_distance - eb_side_radius
var eb_z_offset = dim(1.0 / 16.0)
var eb_height = dim(1.0 / 8.0)

//-----------------------------------------------------------------------------
// valve bosses

var valve_diameter = dim(1.0 / 4.0)
var valve_radius = valve_diameter / 2.0
var valve_y_offset = dim(1.0 / 8.0)
var valve_wall = dim(5.0 / 32.0)
var v2v_distance = dim(1.0 / 2.0)
var valve_draft = DtoR(5)

func valve(d float64, mode string) SDF3 {

	var s SDF3
	h := head_height - cylinder_height

	if mode == "boss" {
		delta := h * math.Tan(valve_draft)
		r1 := valve_radius + valve_wall
		r0 := r1 + delta
		s = NewConeSDF3(h, r0, r1, 0)
	} else if mode == "hole" {
		s = NewCylinderSDF3(h, valve_radius, 0)
	} else {
		panic("bad mode")
	}

	z_ofs := cylinder_height / 2
	return NewTransformSDF3(s, Translate3d(V3{d, valve_y_offset, z_ofs}))
}

func valve_set(d float64, mode string) SDF3 {
	delta := v2v_distance / 2
	s := NewUnionSDF3(valve(-delta, mode), valve(delta, mode))
	s.(*UnionSDF3).SetMin(PolyMin, general_round)
	return NewTransformSDF3(s, Translate3d(V3{d, 0, 0}))
}

func valve_sets(mode string) SDF3 {
	delta := c2c_distance / 2
	return NewUnionSDF3(valve_set(-delta, mode), valve_set(delta, mode))
}

//-----------------------------------------------------------------------------
// cylinder domes (or full base)

var cylinder_height = dim(3.0 / 16.0)
var cylinder_diameter = dim(1.0 + (1.0 / 8.0))
var cylinder_wall = dim(1.0 / 4.0)
var cylinder_radius = cylinder_diameter / 2.0

var dome_radius = cylinder_wall + cylinder_radius
var dome_height = cylinder_wall + cylinder_height
var dome_draft = DtoR(5)

var c2c_distance = dim(1.0 + (3.0 / 8.0))

func cylinder_head(d float64, mode string) SDF3 {
	var s SDF3

	if mode == "dome" {
		z_ofs := (head_height - dome_height) / 2
		extra_z := general_round * 2
		s = NewCylinderSDF3(dome_height+extra_z, dome_radius, general_round)
		s = NewTransformSDF3(s, Translate3d(V3{d, 0, -z_ofs - extra_z}))
	} else if mode == "chamber" {
		z_ofs := (head_height - cylinder_height) / 2
		s = NewCylinderSDF3(cylinder_height, cylinder_radius, 0)
		s = NewTransformSDF3(s, Translate3d(V3{d, 0, -z_ofs}))
	} else {
		panic("bad mode")
	}
	return s
}

func cylinder_heads(mode string) SDF3 {
	x_ofs := c2c_distance / 2
	s := NewUnionSDF3(cylinder_head(-x_ofs, mode), cylinder_head(x_ofs, mode))
	if mode == "dome" {
		s.(*UnionSDF3).SetMin(PolyMin, general_round)
	}
	return s
}

func head_base() SDF3 {
	z_ofs := (head_height - dome_height) / 2
	s := NewExtrudeSDF3(head_wall_inner_2d(), dome_height)
	return NewTransformSDF3(s, Translate3d(V3{0, 0, -z_ofs}))
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
	s.(*DifferenceSDF2).SetMax(PolyMax, general_round)
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

//-----------------------------------------------------------------------------
// manifolds

var manifold_radius = dim(4.5 / 16.0)
var manifold_hole_radius = dim(1.0 / 8.0)
var inlet_theta = 30.2564
var exhaust_theta = 270.0 + 13.9736
var exhaust_x_offset = (c2c_distance / 2) + (v2v_distance / 2)
var inlet_x_offset = (c2c_distance / 2) - (v2v_distance / 2)

func manifold_set(r float64) SDF3 {

	h := dim(2)

	s_ex := NewCylinderSDF3(h, r, 0)
	m := Translate3d(V3{0, 0, h / 2})
	m = RotateX(DtoR(-90)).Mul(m)
	m = RotateZ(DtoR(exhaust_theta)).Mul(m)
	m = Translate3d(V3{exhaust_x_offset, valve_y_offset, eb_z_offset}).Mul(m)
	s_ex = NewTransformSDF3(s_ex, m)

	s_in := NewCylinderSDF3(h, r, 0)
	m = Translate3d(V3{0, 0, h / 2})
	m = RotateX(DtoR(-90)).Mul(m)
	m = RotateZ(DtoR(inlet_theta)).Mul(m)
	m = Translate3d(V3{inlet_x_offset, valve_y_offset, eb_z_offset}).Mul(m)
	s_in = NewTransformSDF3(s_in, m)

	return NewUnionSDF3(s_ex, s_in)
}

func manifolds(mode string) SDF3 {

	var r float64

	if mode == "body" {
		r = manifold_radius
	} else if mode == "hole" {
		r = manifold_hole_radius
	} else {
		panic("bad mode")
	}

	s0 := manifold_set(r)
	s1 := NewTransformSDF3(s0, MirrorYZ())
	s := NewUnionSDF3(s0, s1)
	if mode == "body" {
		s.(*UnionSDF3).SetMin(PolyMin, general_round)
	}

	return s
}

//-----------------------------------------------------------------------------

func additive() SDF3 {
	var s SDF3
	s = NewUnionSDF3(s, head_wall())
	//s = NewUnionSDF3(s, head_base())

	s = NewUnionSDF3(s, cylinder_heads("dome"))
	s.(*UnionSDF3).SetMin(PolyMin, general_round)

	s = NewUnionSDF3(s, valve_sets("boss"))
	s.(*UnionSDF3).SetMin(PolyMin, general_round)

	s = NewUnionSDF3(s, manifolds("body"))
	s.(*UnionSDF3).SetMin(PolyMin, general_round)

	// cleanup the blending artifacts on the outside
	s = NewIntersectionSDF3(s, head_envelope())
	return s
}

//-----------------------------------------------------------------------------

func subtractive() SDF3 {
	var s SDF3
	if casting {
	} else {
		s = NewUnionSDF3(s, cylinder_heads("chamber"))
		s = NewUnionSDF3(s, head_stud_holes())
		s = NewUnionSDF3(s, valve_sets("hole"))
		s = NewUnionSDF3(s, manifolds("hole"))
	}
	return s
}

//-----------------------------------------------------------------------------

func main() {
	s := NewDifferenceSDF3(additive(), subtractive())
	RenderSTL(s, 200, "head.stl")
}

//-----------------------------------------------------------------------------
