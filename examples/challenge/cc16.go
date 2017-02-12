//-----------------------------------------------------------------------------

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------
// CAD Challenge #16 Part A
// https://www.reddit.com/r/cad/comments/5t5z31/cad_challenge_16/

func cc16a() SDF3 {

	base_w := 4.5
	base_d := 2.0
	base_h := 0.62
	base_radius := 0.5

	slot_l := 0.5 * 2
	slot_r := 0.38 / 2.0

	base_2d := NewBoxSDF2(V2{base_w, base_d}, base_radius)
	slot := NewLineSDF2(slot_l, slot_r)
	slot0 := NewTransformSDF2(slot, Translate2d(V2{base_w / 2, 0}))
	slot1 := NewTransformSDF2(slot, Translate2d(V2{-base_w / 2, 0}))
	slots := NewUnionSDF2(slot0, slot1)
	base_2d = NewDifferenceSDF2(base_2d, slots)
	base_3d := NewExtrudeSDF3(base_2d, base_h)

	hole_h := 0.75
	block_radius := 1.0
	block_w := 0.62
	block_l := base_h + 2.0*hole_h
	y_ofs := (base_d - block_w) / 2

	hole_radius := 0.625 / 2.0
	cb_radius := 1.25 / 2.0
	cb_depth := 0.12

	block_2d := NewLineSDF2(block_l, block_radius)
	block_2d = NewCutSDF2(block_2d, V2{0, 0}, V2{0, 1})
	block_3d := NewExtrudeSDF3(block_2d, block_w)

	cb_3d := CounterBored_Hole3d(block_w, hole_radius, cb_radius, cb_depth)
	cb_3d = NewTransformSDF3(cb_3d, Translate3d(V3{block_l / 2, 0, 0}))
	block_3d = NewDifferenceSDF3(block_3d, cb_3d)

	m := RotateX(DtoR(-90))
	m = RotateY(DtoR(-90)).Mul(m)
	m = Translate3d(V3{0, y_ofs, 0}).Mul(m)
	block_3d = NewTransformSDF3(block_3d, m)

	return NewUnionSDF3(base_3d, block_3d)
}

//-----------------------------------------------------------------------------
// CAD Challenge #16 Part B
// https://www.reddit.com/r/cad/comments/5t5z31/cad_challenge_16/

func cc16b() SDF3 {

	// Base
	base_w := 120.0
	base_d := 80.0
	base_h := 24.0
	base_radius := 25.0

	// 2d rounded box - larger so we can remove an edge
	base_2d := NewBoxSDF2(V2{base_w, 2 * base_d}, base_radius)

	// remove the edge and re-center on y-axis
	base_2d = NewCutSDF2(base_2d, V2{0, 0}, V2{-1, 0})
	base_2d = NewTransformSDF2(base_2d, Translate2d(V2{0, -base_d / 2}))

	// cut out the base holes
	base_hole_r := 14.0 / 2.0
	base_hole_yofs := (base_d / 2.0) - 25.0
	base_hole_xofs := (base_w / 2.0) - 25.0
	holes := []V2{
		V2{base_hole_xofs, base_hole_yofs},
		V2{-base_hole_xofs, base_hole_yofs},
	}
	holes_2d := NewMultiCircleSDF2(base_hole_r, holes)
	base_2d = NewDifferenceSDF2(base_2d, holes_2d)

	// cut out the slotted hole
	slot_l := 20.0
	slot_r := 9.0
	slot_2d := NewLineSDF2(slot_l, slot_r)
	m := Rotate2d(DtoR(90))
	m = Translate2d(V2{0, slot_l / 2}).Mul(m)
	slot_2d = NewTransformSDF2(slot_2d, m)
	base_2d = NewDifferenceSDF2(base_2d, slot_2d)

	// Extrude the base to 3d
	base_3d := NewExtrudeSDF3(base_2d, base_h)

	// cut out the rails
	rail_w := 15.0 // rails have square cross section
	rail_zofs := (base_h - rail_w) / 2.0
	rail_3d := NewBoxSDF3(V3{rail_w, base_d, rail_w}, 0)
	rail0_3d := NewTransformSDF3(rail_3d, Translate3d(V3{base_hole_xofs, 0, -rail_zofs}))
	rail1_3d := NewTransformSDF3(rail_3d, Translate3d(V3{-base_hole_xofs, 0, -rail_zofs}))
	base_3d = NewDifferenceSDF3(base_3d, rail0_3d)
	base_3d = NewDifferenceSDF3(base_3d, rail1_3d)

	// cut out the surface recess
	recess_w := 40.0
	recess_h := 2.0
	recess_zofs := (base_h / 2.0) - recess_h
	recess := []V2{
		V2{0, 0},
		V2{recess_w, 0},
		V2{recess_w + recess_h, recess_h},
		V2{0, recess_h},
	}
	recess_2d := NewPolySDF2(recess)
	recess_3d := NewExtrudeSDF3(recess_2d, base_w)
	q := RotateX(DtoR(90))
	q = RotateZ(DtoR(-90)).Mul(q)
	q = Translate3d(V3{0, recess_w, recess_zofs}).Mul(q)
	recess_3d = NewTransformSDF3(recess_3d, q)
	base_3d = NewDifferenceSDF3(base_3d, recess_3d)

	// Tool Support
	support_h := 109.0 - base_h
	support_w := 24.0
	support_base_w := 14.0
	support_theta := math.Atan(support_h / (support_w - support_base_w)) // 83d 17m 25s
	support_xofs := support_h / math.Tan(support_theta)

	// make a polygon for the support profile
	facets := 5
	support := NewSmoother(false)
	support.Add(V2{base_w / 2, -1})
	support.Add(V2{base_w / 2, 0})
	support.AddSmooth(V2{base_hole_xofs, 0}, facets, 5.0)
	support.AddSmooth(V2{base_hole_xofs + support_xofs, support_h}, 3*facets, 25.0)
	support.AddSmooth(V2{-base_hole_xofs - support_xofs, support_h}, 3*facets, 25.0)
	support.AddSmooth(V2{-base_hole_xofs, 0}, facets, 5.0)
	support.Add(V2{-base_w / 2, 0})
	support.Add(V2{-base_w / 2, -1})
	support.Smooth()
	//RenderDXF("support.dxf", support.Vertices())
	support_2d := NewPolySDF2(support.Vertices())

	// extrude the support to 3d
	support_3d := NewExtrudeSDF3(support_2d, support_w)

	// remove the chamfered hole
	hole_h := 84.0 - base_h
	hole_r := 35.0 / 2.0
	chamfer_d := 2.0
	hole_3d := Chamfered_Hole3d(support_w, hole_r, chamfer_d)
	q = Translate3d(V3{0, hole_h, 0})
	hole_3d = NewTransformSDF3(hole_3d, q)
	support_3d = NewDifferenceSDF3(support_3d, hole_3d)

	// cut the sloped face of the support
	support_3d = NewCutSDF3(support_3d, V3{0, support_h, -support_w / 2}, V3{0, math.Cos(support_theta), math.Sin(support_theta)})

	// position the support
	support_yofs := (base_d - support_w) / 2.0
	q = RotateX(DtoR(90))
	q = Translate3d(V3{0, -support_yofs, base_h / 2}).Mul(q)
	support_3d = NewTransformSDF3(support_3d, q)

	// Gussets
	gusset_l := 20.0
	gusset_w := 3.0
	gusset_xofs := 37.0 / 2.0
	gusset_h := 12.53

	gusset_yofs := base_d / 2.0
	gusset_yofs -= support_base_w
	gusset_yofs -= gusset_h / math.Tan(support_theta)
	gusset_yofs -= gusset_h

	gusset := NewSmoother(false)
	gusset.Add(V2{gusset_l, 0})
	gusset.AddSmooth(V2{0, 0}, facets, 20.0)
	gusset.Add(V2{-gusset_l, gusset_l})
	gusset.Add(V2{-gusset_l, 0})
	gusset.Smooth()
	//RenderDXF("gusset.dxf", gusset.Vertices())
	gusset_2d := NewPolySDF2(gusset.Vertices())

	// extrude the gusset to 3d
	gusset_3d := NewExtrudeSDF3(gusset_2d, gusset_w)

	// orient the gusset
	q = RotateX(DtoR(90))
	q = RotateZ(DtoR(90)).Mul(q)
	q = Translate3d(V3{0, -gusset_yofs, base_h / 2}).Mul(q)
	gusset_3d = NewTransformSDF3(gusset_3d, q)

	gusset0_3d := NewTransformSDF3(gusset_3d, Translate3d(V3{gusset_xofs, 0, 0}))
	gusset1_3d := NewTransformSDF3(gusset_3d, Translate3d(V3{-gusset_xofs, 0, 0}))
	gusset_3d = NewUnionSDF3(gusset0_3d, gusset1_3d)

	return NewUnionSDF3(base_3d, NewUnionSDF3(support_3d, gusset_3d))
}

//-----------------------------------------------------------------------------
