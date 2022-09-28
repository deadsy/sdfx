//-----------------------------------------------------------------------------
/*

Wallaby Cylinder Head

No draft version for 3d printing and lost-PLA investment casting.

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
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

const shrink = desired_scale * al_shrink * pla_shrink

const general_round = 0.1

//-----------------------------------------------------------------------------
// exhaust bosses

const eb_side_radius = 5.0 / 32.0
const eb_main_radius = 5.0 / 16.0
const eb_hole_radius = 3.0 / 16.0
const eb_c2c_distance = 13.0 / 16.0
const eb_distance = eb_c2c_distance / 2.0

const eb_x_offset = 0.5*(head_length+eb_height) - eb_height0
const eb_y_offset = (head_width / 2.0) - eb_distance - eb_side_radius
const eb_z_offset = 1.0 / 16.0

const eb_height0 = 1.0 / 16.0
const eb_height1 = 1.0 / 8.0
const eb_height = eb_height0 + eb_height1

func exhaust_boss(mode string, x_ofs float64) sdf.SDF3 {

	var s0 sdf.SDF2

	if mode == "body" {
		s0 = sdf.NewFlange1(eb_distance, eb_main_radius, eb_side_radius)
	} else if mode == "hole" {
		s0, _ = sdf.Circle2D(eb_hole_radius)
	} else {
		panic("bad mode")
	}

	s1 := sdf.Extrude3D(s0, eb_height)
	m := sdf.RotateZ(sdf.DtoR(90))
	m = sdf.RotateY(sdf.DtoR(90)).Mul(m)
	m = sdf.Translate3d(v3.Vec{x_ofs, eb_y_offset, eb_z_offset}).Mul(m)
	s1 = sdf.Transform3D(s1, m)
	return s1
}

func exhaust_bosses(mode string) sdf.SDF3 {
	return sdf.Union3D(exhaust_boss(mode, eb_x_offset), exhaust_boss(mode, -eb_x_offset))
}

//-----------------------------------------------------------------------------
// spark plug bosses

const sp2sp_distance = 1.0 + (5.0 / 8.0)

var sp_theta = sdf.DtoR(30)

const sp_boss_r1 = 21.0 / 64.0
const sp_boss_r2 = 15.0 / 32.0
const sp_boss_h1 = 0.79
const sp_boss_h2 = 0.94
const sp_boss_h3 = 2

const sp_hole_d = 21.0 / 64.0
const sp_hole_r = sp_hole_d / 2.0
const sp_hole_h = 1.0

const sp_cb_h1 = 1.0
const sp_cb_h2 = 2.0
const sp_cb_r = 5.0 / 16.0

var sp_hyp = sp_hole_h + sp_cb_r*math.Tan(sp_theta)
var sp_y_ofs = sp_hyp*math.Cos(sp_theta) - head_width/2
var sp_z_ofs = -sp_hyp * math.Sin(sp_theta)

func sparkplug(mode string, x_ofs float64) sdf.SDF3 {
	var vlist []v2.Vec
	if mode == "boss" {
		boss := sdf.NewPolygon()
		boss.Add(0, 0)
		boss.Add(sp_boss_r1, 0)
		boss.Add(sp_boss_r1, sp_boss_h1).Smooth(sp_boss_r1*0.3, 3)
		boss.Add(sp_boss_r2, sp_boss_h2).Smooth(sp_boss_r2*0.3, 3)
		boss.Add(sp_boss_r2, sp_boss_h3)
		boss.Add(0, sp_boss_h3)
		vlist = boss.Vertices()
	} else if mode == "hole" {
		vlist = []v2.Vec{
			{0, 0},
			{sp_hole_r, 0},
			{sp_hole_r, sp_hole_h},
			{0, sp_hole_h},
		}
	} else if mode == "counterbore" {
		p := sdf.NewPolygon()
		p.Add(0, sp_cb_h1)
		p.Add(sp_cb_r, sp_cb_h1).Smooth(sp_cb_r/6.0, 3)
		p.Add(sp_cb_r, sp_cb_h2)
		p.Add(0, sp_cb_h2)
		vlist = p.Vertices()
	} else {
		panic("bad mode")
	}
	s0, _ := sdf.Polygon2D(vlist)
	s, _ := sdf.Revolve3D(s0)
	m := sdf.RotateX(sdf.Pi/2 - sp_theta)
	m = sdf.Translate3d(v3.Vec{x_ofs, sp_y_ofs, sp_z_ofs}).Mul(m)
	s = sdf.Transform3D(s, m)
	return s
}

func sparkplugs(mode string) sdf.SDF3 {
	x_ofs := 0.5 * sp2sp_distance
	return sdf.Union3D(sparkplug(mode, x_ofs), sparkplug(mode, -x_ofs))
}

//-----------------------------------------------------------------------------
// valve bosses

const valve_diameter = 1.0 / 4.0
const valve_radius = valve_diameter / 2.0
const valve_y_offset = 1.0 / 8.0
const valve_wall = 5.0 / 32.0
const v2v_distance = 1.0 / 2.0

var valve_draft = sdf.DtoR(5)

func valve(d float64, mode string) sdf.SDF3 {

	var s sdf.SDF3
	h := head_height - cylinder_height

	if mode == "boss" {
		delta := h * math.Tan(valve_draft)
		r1 := valve_radius + valve_wall
		r0 := r1 + delta
		s, _ = sdf.Cone3D(h, r0, r1, 0)
	} else if mode == "hole" {
		s, _ = sdf.Cylinder3D(h, valve_radius, 0)
	} else {
		panic("bad mode")
	}

	z_ofs := cylinder_height / 2
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{d, valve_y_offset, z_ofs}))
}

func valve_set(d float64, mode string) sdf.SDF3 {
	delta := v2v_distance / 2
	s := sdf.Union3D(valve(-delta, mode), valve(delta, mode))
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(general_round))
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{d, 0, 0}))
}

func valve_sets(mode string) sdf.SDF3 {
	delta := c2c_distance / 2
	return sdf.Union3D(valve_set(-delta, mode), valve_set(delta, mode))
}

//-----------------------------------------------------------------------------
// cylinder domes (or full base)

const cylinder_height = 3.0 / 16.0
const cylinder_diameter = 1.0 + (1.0 / 8.0)
const cylinder_wall = 1.0 / 4.0
const cylinder_radius = cylinder_diameter / 2.0

const dome_radius = cylinder_wall + cylinder_radius
const dome_height = cylinder_wall + cylinder_height

var dome_draft = sdf.DtoR(5)

const c2c_distance = 1.0 + (3.0 / 8.0)

func cylinder_head(d float64, mode string) sdf.SDF3 {
	var s sdf.SDF3

	if mode == "dome" {
		z_ofs := (head_height - dome_height) / 2
		extra_z := general_round * 2
		s, _ = sdf.Cylinder3D(dome_height+extra_z, dome_radius, general_round)
		s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{d, 0, -z_ofs - extra_z}))
	} else if mode == "chamber" {
		z_ofs := (head_height - cylinder_height) / 2
		s, _ = sdf.Cylinder3D(cylinder_height, cylinder_radius, 0)
		s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{d, 0, -z_ofs}))
	} else {
		panic("bad mode")
	}
	return s
}

func cylinder_heads(mode string) sdf.SDF3 {
	x_ofs := c2c_distance / 2
	s := sdf.Union3D(cylinder_head(-x_ofs, mode), cylinder_head(x_ofs, mode))
	if mode == "dome" {
		s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(general_round))
	}
	return s
}

func head_base() sdf.SDF3 {
	z_ofs := (head_height - dome_height) / 2
	s := sdf.Extrude3D(head_wall_inner_2d(), dome_height)
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, -z_ofs}))
}

//-----------------------------------------------------------------------------
// cylinder studs: location, bosses and holes

const stud_hole_radius = 1.0 / 16.0
const stud_boss_radius = 3.0 / 16.0
const stud_hole_dy = 11.0 / 16.0
const stud_hole_dx0 = 7.0 / 16.0
const stud_hole_dx1 = 1.066

var stud_locations = []v2.Vec{
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

func head_stud_holes() sdf.SDF3 {
	c, _ := sdf.Circle2D(stud_hole_radius)
	s := sdf.Multi2D(c, stud_locations)
	return sdf.Extrude3D(s, head_height)
}

//-----------------------------------------------------------------------------
// head walls

const head_length = 4.30 / 1.25
const head_width = 2.33 / 1.25
const head_height = 7.0 / 8.0
const head_corner_round = (5.0 / 32.0) / 1.25
const head_wall_thickness = 0.154

func head_wall_outer_2d() sdf.SDF2 {
	return sdf.Box2D(v2.Vec{head_length, head_width}, head_corner_round)
}

func head_wall_inner_2d() sdf.SDF2 {
	l := head_length - (2 * head_wall_thickness)
	w := head_width - (2 * head_wall_thickness)
	s0 := sdf.Box2D(v2.Vec{l, w}, 0)
	c, _ := sdf.Circle2D(stud_boss_radius)
	s1 := sdf.Multi2D(c, stud_locations)
	s := sdf.Difference2D(s0, s1)
	s.(*sdf.DifferenceSDF2).SetMax(sdf.PolyMax(general_round))
	return s
}

func head_envelope() sdf.SDF3 {
	s0 := sdf.Box2D(v2.Vec{head_length + 2*eb_height1, head_width}, 0)
	return sdf.Extrude3D(s0, head_height)
}

func head_wall() sdf.SDF3 {
	s := head_wall_outer_2d()
	s = sdf.Difference2D(s, head_wall_inner_2d())
	return sdf.Extrude3D(s, head_height)
}

//-----------------------------------------------------------------------------
// manifolds

const manifold_radius = 4.5 / 16.0
const manifold_hole_radius = 1.0 / 8.0
const inlet_theta = 30.2564
const exhaust_theta = 270.0 + 13.9736
const exhaust_x_offset = (c2c_distance / 2) + (v2v_distance / 2)
const inlet_x_offset = (c2c_distance / 2) - (v2v_distance / 2)

func manifold_set(r float64) sdf.SDF3 {

	const h = 2

	s_ex, _ := sdf.Cylinder3D(h, r, 0)
	m := sdf.Translate3d(v3.Vec{0, 0, h / 2})
	m = sdf.RotateX(sdf.DtoR(-90)).Mul(m)
	m = sdf.RotateZ(sdf.DtoR(exhaust_theta)).Mul(m)
	m = sdf.Translate3d(v3.Vec{exhaust_x_offset, valve_y_offset, eb_z_offset}).Mul(m)
	s_ex = sdf.Transform3D(s_ex, m)

	s_in, _ := sdf.Cylinder3D(h, r, 0)
	m = sdf.Translate3d(v3.Vec{0, 0, h / 2})
	m = sdf.RotateX(sdf.DtoR(-90)).Mul(m)
	m = sdf.RotateZ(sdf.DtoR(inlet_theta)).Mul(m)
	m = sdf.Translate3d(v3.Vec{inlet_x_offset, valve_y_offset, eb_z_offset}).Mul(m)
	s_in = sdf.Transform3D(s_in, m)

	return sdf.Union3D(s_ex, s_in)
}

func manifolds(mode string) sdf.SDF3 {
	var r float64
	if mode == "body" {
		r = manifold_radius
	} else if mode == "hole" {
		r = manifold_hole_radius
	} else {
		panic("bad mode")
	}
	s0 := manifold_set(r)
	s1 := sdf.Transform3D(s0, sdf.MirrorYZ())
	s := sdf.Union3D(s0, s1)
	if mode == "body" {
		s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(general_round))
	}
	return s
}

//-----------------------------------------------------------------------------

func allowances(s sdf.SDF3) sdf.SDF3 {
	//eb0_2d := Slice2D(s, v3.Vec{eb_x_offset, 0, 0}, v3.Vec{1, 0, 0})
	//return Extrude3D(eb0_2d, 10.0)
	return nil
}

//-----------------------------------------------------------------------------

func additive() sdf.SDF3 {
	s := sdf.Union3D(
		head_wall(),
		//head_base(),
		cylinder_heads("dome"),
		valve_sets("boss"),
		sparkplugs("boss"),
		manifolds("body"),
		exhaust_bosses("body"),
	)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(general_round))

	s = sdf.Difference3D(s, sparkplugs("counterbore"))

	// cleanup the blending artifacts on the outside
	s = sdf.Intersect3D(s, head_envelope())

	if casting == true {
		s = sdf.Union3D(s, allowances(s))
	}

	return s
}

//-----------------------------------------------------------------------------

func subtractive() sdf.SDF3 {
	var s sdf.SDF3
	if casting == false {
		s = sdf.Union3D(cylinder_heads("chamber"),
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
	s := sdf.Difference3D(additive(), subtractive())
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "head.stl", render.NewMarchingCubesOctree(400))
}

//-----------------------------------------------------------------------------
