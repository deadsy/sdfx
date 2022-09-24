//-----------------------------------------------------------------------------

package main

import (
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// CAD Challenge #16 Part A
// https://www.reddit.com/r/cad/comments/5t5z31/cad_challenge_16/

func cc16a() (sdf.SDF3, error) {

	base_w := 4.5
	base_d := 2.0
	base_h := 0.62
	base_radius := 0.5

	slot_l := 0.5 * 2
	slot_r := 0.38 / 2.0

	base_2d := sdf.Box2D(v2.Vec{base_w, base_d}, base_radius)
	slot := sdf.Line2D(slot_l, slot_r)
	slot0 := sdf.Transform2D(slot, sdf.Translate2d(v2.Vec{base_w / 2, 0}))
	slot1 := sdf.Transform2D(slot, sdf.Translate2d(v2.Vec{-base_w / 2, 0}))
	slots := sdf.Union2D(slot0, slot1)
	base_2d = sdf.Difference2D(base_2d, slots)
	base_3d := sdf.Extrude3D(base_2d, base_h)

	hole_h := 0.75
	block_radius := 1.0
	block_w := 0.62
	block_l := base_h + 2.0*hole_h
	y_ofs := (base_d - block_w) / 2

	hole_radius := 0.625 / 2.0
	cb_radius := 1.25 / 2.0
	cb_depth := 0.12

	block_2d := sdf.Line2D(block_l, block_radius)
	block_2d = sdf.Cut2D(block_2d, v2.Vec{0, 0}, v2.Vec{0, 1})
	block_3d := sdf.Extrude3D(block_2d, block_w)

	cb_3d, err := obj.CounterBoredHole3D(block_w, hole_radius, cb_radius, cb_depth)
	if err != nil {
		return nil, err
	}
	cb_3d = sdf.Transform3D(cb_3d, sdf.Translate3d(v3.Vec{block_l / 2, 0, 0}))
	block_3d = sdf.Difference3D(block_3d, cb_3d)

	m := sdf.RotateX(sdf.DtoR(-90))
	m = sdf.RotateY(sdf.DtoR(-90)).Mul(m)
	m = sdf.Translate3d(v3.Vec{0, y_ofs, 0}).Mul(m)
	block_3d = sdf.Transform3D(block_3d, m)

	return sdf.Union3D(base_3d, block_3d), nil
}

//-----------------------------------------------------------------------------
// CAD Challenge #16 Part B
// https://www.reddit.com/r/cad/comments/5t5z31/cad_challenge_16/

func cc16b() (sdf.SDF3, error) {

	// Base
	base_w := 120.0
	base_d := 80.0
	base_h := 24.0
	base_radius := 25.0

	// 2d rounded box - larger so we can remove an edge
	base_2d := sdf.Box2D(v2.Vec{base_w, 2 * base_d}, base_radius)

	// remove the edge and re-center on y-axis
	base_2d = sdf.Cut2D(base_2d, v2.Vec{0, 0}, v2.Vec{-1, 0})
	base_2d = sdf.Transform2D(base_2d, sdf.Translate2d(v2.Vec{0, -base_d / 2}))

	// cut out the base holes
	base_hole_r := 14.0 / 2.0
	base_hole_yofs := (base_d / 2.0) - 25.0
	base_hole_xofs := (base_w / 2.0) - 25.0
	holes := []v2.Vec{
		{base_hole_xofs, base_hole_yofs},
		{-base_hole_xofs, base_hole_yofs},
	}

	c, err := sdf.Circle2D(base_hole_r)
	if err != nil {
		return nil, err
	}
	holes_2d := sdf.Multi2D(c, holes)
	base_2d = sdf.Difference2D(base_2d, holes_2d)

	// cut out the slotted hole
	slot_l := 20.0
	slot_r := 9.0
	slot_2d := sdf.Line2D(slot_l, slot_r)
	m := sdf.Rotate2d(sdf.DtoR(90))
	m = sdf.Translate2d(v2.Vec{0, slot_l / 2}).Mul(m)
	slot_2d = sdf.Transform2D(slot_2d, m)
	base_2d = sdf.Difference2D(base_2d, slot_2d)

	// Extrude the base to 3d
	base_3d := sdf.Extrude3D(base_2d, base_h)

	// cut out the rails
	rail_w := 15.0 // rails have square cross section
	rail_zofs := (base_h - rail_w) / 2.0
	rail_3d, err := sdf.Box3D(v3.Vec{rail_w, base_d, rail_w}, 0)
	if err != nil {
		return nil, err
	}

	rail0_3d := sdf.Transform3D(rail_3d, sdf.Translate3d(v3.Vec{base_hole_xofs, 0, -rail_zofs}))
	rail1_3d := sdf.Transform3D(rail_3d, sdf.Translate3d(v3.Vec{-base_hole_xofs, 0, -rail_zofs}))
	base_3d = sdf.Difference3D(base_3d, rail0_3d)
	base_3d = sdf.Difference3D(base_3d, rail1_3d)

	// cut out the surface recess
	recess_w := 40.0
	recess_h := 2.0
	recess_zofs := (base_h / 2.0) - recess_h
	recess := []v2.Vec{
		{0, 0},
		{recess_w, 0},
		{recess_w + recess_h, recess_h},
		{0, recess_h},
	}
	recess_2d, err := sdf.Polygon2D(recess)
	if err != nil {
		return nil, err
	}
	recess_3d := sdf.Extrude3D(recess_2d, base_w)
	q := sdf.RotateX(sdf.DtoR(90))
	q = sdf.RotateZ(sdf.DtoR(-90)).Mul(q)
	q = sdf.Translate3d(v3.Vec{0, recess_w, recess_zofs}).Mul(q)
	recess_3d = sdf.Transform3D(recess_3d, q)
	base_3d = sdf.Difference3D(base_3d, recess_3d)

	// Tool Support
	support_h := 109.0 - base_h
	support_w := 24.0
	support_base_w := 14.0
	support_theta := math.Atan(support_h / (support_w - support_base_w)) // 83d 17m 25s
	support_xofs := support_h / math.Tan(support_theta)

	// make a polygon for the support profile
	facets := 5
	support := sdf.NewPolygon()
	support.Add(base_w/2, -1)
	support.Add(base_w/2, 0)
	support.Add(base_hole_xofs, 0).Smooth(5.0, facets)
	support.Add(base_hole_xofs+support_xofs, support_h).Smooth(25.0, 3*facets)
	support.Add(-base_hole_xofs-support_xofs, support_h).Smooth(25.0, 3*facets)
	support.Add(-base_hole_xofs, 0).Smooth(5.0, facets)
	support.Add(-base_w/2, 0)
	support.Add(-base_w/2, -1)
	//support.Render("support.dxf")
	support_2d, err := sdf.Polygon2D(support.Vertices())
	if err != nil {
		return nil, err
	}
	// extrude the support to 3d
	support_3d := sdf.Extrude3D(support_2d, support_w)

	// remove the chamfered hole
	hole_h := 84.0 - base_h
	hole_r := 35.0 / 2.0
	chamfer_d := 2.0
	hole_3d, err := obj.ChamferedHole3D(support_w, hole_r, chamfer_d)
	if err != nil {
		return nil, err
	}

	q = sdf.Translate3d(v3.Vec{0, hole_h, 0})
	hole_3d = sdf.Transform3D(hole_3d, q)
	support_3d = sdf.Difference3D(support_3d, hole_3d)

	// cut the sloped face of the support
	support_3d = sdf.Cut3D(support_3d, v3.Vec{0, support_h, -support_w / 2}, v3.Vec{0, math.Cos(support_theta), math.Sin(support_theta)})

	// position the support
	support_yofs := (base_d - support_w) / 2.0
	q = sdf.RotateX(sdf.DtoR(90))
	q = sdf.Translate3d(v3.Vec{0, -support_yofs, base_h / 2}).Mul(q)
	support_3d = sdf.Transform3D(support_3d, q)

	// Gussets
	gusset_l := 20.0
	gusset_w := 3.0
	gusset_xofs := 37.0 / 2.0
	gusset_h := 12.53

	gusset_yofs := base_d / 2.0
	gusset_yofs -= support_base_w
	gusset_yofs -= gusset_h / math.Tan(support_theta)
	gusset_yofs -= gusset_h

	gusset := sdf.NewPolygon()
	gusset.Add(gusset_l, 0)
	gusset.Add(0, 0).Smooth(20.0, facets)
	gusset.Add(-gusset_l, gusset_l)
	gusset.Add(-gusset_l, 0)
	//gusset.Render("gusset.dxf")
	gusset_2d, err := sdf.Polygon2D(gusset.Vertices())
	if err != nil {
		return nil, err
	}

	// extrude the gusset to 3d
	gusset_3d := sdf.Extrude3D(gusset_2d, gusset_w)

	// orient the gusset
	q = sdf.RotateX(sdf.DtoR(90))
	q = sdf.RotateZ(sdf.DtoR(90)).Mul(q)
	q = sdf.Translate3d(v3.Vec{0, -gusset_yofs, base_h / 2}).Mul(q)
	gusset_3d = sdf.Transform3D(gusset_3d, q)

	gusset0_3d := sdf.Transform3D(gusset_3d, sdf.Translate3d(v3.Vec{gusset_xofs, 0, 0}))
	gusset1_3d := sdf.Transform3D(gusset_3d, sdf.Translate3d(v3.Vec{-gusset_xofs, 0, 0}))
	gusset_3d = sdf.Union3D(gusset0_3d, gusset1_3d)

	return sdf.Union3D(base_3d, support_3d, gusset_3d), nil
}

//-----------------------------------------------------------------------------
