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

const desired_scale = 1.25
const al_shrink = 1.0 / 0.99   // ~1%
const pla_shrink = 1.0 / 0.998 //~0.2%
const abs_shrink = 1.0 / 0.995 //~0.5%

// dimension scaling
func dim(x float64) float64 {
	return x * desired_scale * MillimetresPerInch * al_shrink * pla_shrink
}

var general_round = dim(0.1)

//-----------------------------------------------------------------------------
// exhaust bosses

var eb_side_radius = dim(5.0 / 32.0)
var eb_main_radius = dim(5.0 / 16.0)
var eb_hole_radius = dim(3.0 / 16.0)
var eb_c2c_distance = dim(13.0 / 16.0)
var eb_distance = eb_c2c_distance / 2.0

var eb_x_offset = 0.5*(head_length+eb_height) - eb_height0
var eb_y_offset = (head_width / 2.0) - eb_distance - eb_side_radius
var eb_z_offset = dim(1.0 / 16.0)

var eb_height0 = dim(1.0 / 16.0)
var eb_height1 = dim(1.0 / 8.0)
var eb_height = eb_height0 + eb_height1

func exhaust_boss(mode string, x_ofs float64) SDF3 {

	var s0 SDF2

	if mode == "body" {
		s0 = NewFlange1(eb_distance, eb_main_radius, eb_side_radius)
	} else if mode == "hole" {
		s0 = Circle2D(eb_hole_radius)
	} else {
		panic("bad mode")
	}

	s1 := Extrude3D(s0, eb_height)
	m := RotateZ(DtoR(90))
	m = RotateY(DtoR(90)).Mul(m)
	m = Translate3d(V3{x_ofs, eb_y_offset, eb_z_offset}).Mul(m)
	s1 = Transform3D(s1, m)
	return s1
}

func exhaust_bosses(mode string) SDF3 {
	return Union3D(exhaust_boss(mode, eb_x_offset), exhaust_boss(mode, -eb_x_offset))
}

//-----------------------------------------------------------------------------
// spark plug bosses

var sp2sp_distance = dim(1.0 + (5.0 / 8.0))
var sp_theta = DtoR(30)

var sp_boss_r1 = dim(21.0 / 64.0)
var sp_boss_r2 = dim(15.0 / 32.0)
var sp_boss_h1 = dim(0.79)
var sp_boss_h2 = dim(0.94)
var sp_boss_h3 = dim(2)

var sp_hole_d = dim(21.0 / 64.0)
var sp_hole_r = sp_hole_d / 2.0
var sp_hole_h = dim(1.0)

var sp_cb_h1 = dim(1.0)
var sp_cb_h2 = dim(2.0)
var sp_cb_r = dim(5.0 / 16.0)

var sp_hyp = sp_hole_h + sp_cb_r*math.Tan(sp_theta)
var sp_y_ofs = sp_hyp*math.Cos(sp_theta) - head_width/2
var sp_z_ofs = -sp_hyp * math.Sin(sp_theta)

func sparkplug(mode string, x_ofs float64) SDF3 {
	var vlist []V2
	if mode == "boss" {
		boss := NewPolygon()
		boss.Add(0, 0)
		boss.Add(sp_boss_r1, 0)
		boss.Add(sp_boss_r1, sp_boss_h1).Smooth(sp_boss_r1*0.3, 3)
		boss.Add(sp_boss_r2, sp_boss_h2).Smooth(sp_boss_r2*0.3, 3)
		boss.Add(sp_boss_r2, sp_boss_h3)
		boss.Add(0, sp_boss_h3)
		vlist = boss.Vertices()
	} else if mode == "hole" {
		vlist = []V2{
			{0, 0},
			{sp_hole_r, 0},
			{sp_hole_r, sp_hole_h},
			{0, sp_hole_h},
		}
	} else if mode == "counterbore" {
		p := NewPolygon()
		p.Add(0, sp_cb_h1)
		p.Add(sp_cb_r, sp_cb_h1).Smooth(sp_cb_r/6.0, 3)
		p.Add(sp_cb_r, sp_cb_h2)
		p.Add(0, sp_cb_h2)
		vlist = p.Vertices()
	} else {
		panic("bad mode")
	}
	s0 := Polygon2D(vlist)
	s := Revolve3D(s0)
	m := RotateX(Pi/2 - sp_theta)
	m = Translate3d(V3{x_ofs, sp_y_ofs, sp_z_ofs}).Mul(m)
	s = Transform3D(s, m)
	return s
}

func sparkplugs(mode string) SDF3 {
	x_ofs := 0.5 * sp2sp_distance
	return Union3D(sparkplug(mode, x_ofs), sparkplug(mode, -x_ofs))
}

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
		s = Cone3D(h, r0, r1, 0)
	} else if mode == "hole" {
		s = Cylinder3D(h, valve_radius, 0)
	} else {
		panic("bad mode")
	}

	z_ofs := cylinder_height / 2
	return Transform3D(s, Translate3d(V3{d, valve_y_offset, z_ofs}))
}

func valve_set(d float64, mode string) SDF3 {
	delta := v2v_distance / 2
	s := Union3D(valve(-delta, mode), valve(delta, mode))
	s.(*UnionSDF3).SetMin(PolyMin(general_round))
	return Transform3D(s, Translate3d(V3{d, 0, 0}))
}

func valve_sets(mode string) SDF3 {
	delta := c2c_distance / 2
	return Union3D(valve_set(-delta, mode), valve_set(delta, mode))
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
		s = Cylinder3D(dome_height+extra_z, dome_radius, general_round)
		s = Transform3D(s, Translate3d(V3{d, 0, -z_ofs - extra_z}))
	} else if mode == "chamber" {
		z_ofs := (head_height - cylinder_height) / 2
		s = Cylinder3D(cylinder_height, cylinder_radius, 0)
		s = Transform3D(s, Translate3d(V3{d, 0, -z_ofs}))
	} else {
		panic("bad mode")
	}
	return s
}

func cylinder_heads(mode string) SDF3 {
	x_ofs := c2c_distance / 2
	s := Union3D(cylinder_head(-x_ofs, mode), cylinder_head(x_ofs, mode))
	if mode == "dome" {
		s.(*UnionSDF3).SetMin(PolyMin(general_round))
	}
	return s
}

func head_base() SDF3 {
	z_ofs := (head_height - dome_height) / 2
	s := Extrude3D(head_wall_inner_2d(), dome_height)
	return Transform3D(s, Translate3d(V3{0, 0, -z_ofs}))
}

//-----------------------------------------------------------------------------
// cylinder studs: location, bosses and holes

var stud_hole_radius = dim(1.0 / 16.0)
var stud_boss_radius = dim(3.0 / 16.0)
var stud_hole_dy = dim(11.0 / 16.0)
var stud_hole_dx0 = dim(7.0 / 16.0)
var stud_hole_dx1 = dim(1.066)

var stud_locations = []V2{
	{stud_hole_dx0 + stud_hole_dx1, 0},
	{stud_hole_dx0 + stud_hole_dx1, stud_hole_dy},
	{stud_hole_dx0 + stud_hole_dx1, -stud_hole_dy},
	{stud_hole_dx0, stud_hole_dy},
	{stud_hole_dx0, -stud_hole_dy},
	{-stud_hole_dx0 - stud_hole_dx1, 0},
	{-stud_hole_dx0 - stud_hole_dx1, stud_hole_dy},
	{-stud_hole_dx0 - stud_hole_dx1, -stud_hole_dy},
	{-stud_hole_dx0, stud_hole_dy},
	{-stud_hole_dx0, -stud_hole_dy},
}

func head_stud_holes() SDF3 {
	s := MultiCircle2D(stud_hole_radius, stud_locations)
	return Extrude3D(s, head_height)
}

//-----------------------------------------------------------------------------
// head walls

var head_length = dim(4.30 / 1.25)
var head_width = dim(2.33 / 1.25)
var head_height = dim(7.0 / 8.0)
var head_corner_round = dim((5.0 / 32.0) / 1.25)
var head_wall_thickness = dim(0.154)

func head_wall_outer_2d() SDF2 {
	return Box2D(V2{head_length, head_width}, head_corner_round)
}

func head_wall_inner_2d() SDF2 {
	l := head_length - (2 * head_wall_thickness)
	w := head_width - (2 * head_wall_thickness)
	s0 := Box2D(V2{l, w}, 0)
	s1 := MultiCircle2D(stud_boss_radius, stud_locations)
	s := Difference2D(s0, s1)
	s.(*DifferenceSDF2).SetMax(PolyMax(general_round))
	return s
}

func head_envelope() SDF3 {
	s0 := Box2D(V2{head_length + 2*eb_height1, head_width}, 0)
	return Extrude3D(s0, head_height)
}

func head_wall() SDF3 {
	s := head_wall_outer_2d()
	s = Difference2D(s, head_wall_inner_2d())
	return Extrude3D(s, head_height)
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

	s_ex := Cylinder3D(h, r, 0)
	m := Translate3d(V3{0, 0, h / 2})
	m = RotateX(DtoR(-90)).Mul(m)
	m = RotateZ(DtoR(exhaust_theta)).Mul(m)
	m = Translate3d(V3{exhaust_x_offset, valve_y_offset, eb_z_offset}).Mul(m)
	s_ex = Transform3D(s_ex, m)

	s_in := Cylinder3D(h, r, 0)
	m = Translate3d(V3{0, 0, h / 2})
	m = RotateX(DtoR(-90)).Mul(m)
	m = RotateZ(DtoR(inlet_theta)).Mul(m)
	m = Translate3d(V3{inlet_x_offset, valve_y_offset, eb_z_offset}).Mul(m)
	s_in = Transform3D(s_in, m)

	return Union3D(s_ex, s_in)
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
	s1 := Transform3D(s0, MirrorYZ())
	s := Union3D(s0, s1)
	if mode == "body" {
		s.(*UnionSDF3).SetMin(PolyMin(general_round))
	}
	return s
}

//-----------------------------------------------------------------------------

func allowances(s SDF3) SDF3 {
	//eb0_2d := Slice2D(s, V3{eb_x_offset, 0, 0}, V3{1, 0, 0})
	//return Extrude3D(eb0_2d, 10.0)
	return nil
}

//-----------------------------------------------------------------------------

func additive() SDF3 {
	s := Union3D(
		head_wall(),
		//head_base(),
		cylinder_heads("dome"),
		valve_sets("boss"),
		sparkplugs("boss"),
		manifolds("body"),
		exhaust_bosses("body"),
	)
	s.(*UnionSDF3).SetMin(PolyMin(general_round))

	s = Difference3D(s, sparkplugs("counterbore"))

	// cleanup the blending artifacts on the outside
	s = Intersect3D(s, head_envelope())

	if casting == true {
		s = Union3D(s, allowances(s))
	}

	return s
}

//-----------------------------------------------------------------------------

func subtractive() SDF3 {
	var s SDF3
	if casting == false {
		s = Union3D(cylinder_heads("chamber"),
			head_stud_holes(),
			valve_sets("hole"),
			sparkplugs("hole"),
			manifolds("hole"),
			exhaust_bosses("hole"),
		)
	}
	return s
}

//-----------------------------------------------------------------------------

func main() {
	s := Difference3D(additive(), subtractive())
	RenderSTL(s, 400, "head.stl")
}

//-----------------------------------------------------------------------------
