//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------
// CAD Challenge #18 Part B
// https://www.reddit.com/r/cad/comments/5vwdnc/cad_challenge_18/

func cc18b() SDF3 {

	// build the vertical pipe
	p := NewSmoother(false)
	p.Add(V2{0, 0})
	p.Add(V2{6, 0})
	p.AddSmooth(V2{6, 19}, 5, 0.5)
	p.Add(V2{8, 19})
	p.Add(V2{8, 21})
	p.Add(V2{6, 21})
	p.Add(V2{6, 20})
	p.Add(V2{0, 20})
	p.Smooth()
	vpipe_3d := Revolve3D(Polygon2D(p.Vertices()))
	// bolt circle for the top flange
	top_holes_3d := MakeBoltCircle3D(
		2.0,       // hole_depth
		0.5/2.0,   // hole_radius
		14.50/2.0, // circle_radius
		6,         // num_holes
	)
	m := RotateZ(DtoR(30))
	m = Translate3d(V3{0, 0, 1.0 + 19.0}).Mul(m)
	top_holes_3d = Transform3D(top_holes_3d, m)
	vpipe_3d = Difference3D(vpipe_3d, top_holes_3d)

	// build the horizontal pipe
	p = NewSmoother(false)
	p.Add(V2{0, 0})
	p.Add(V2{5, 0})
	p.AddSmooth(V2{5, 12}, 5, 0.5)
	p.Add(V2{8, 12})
	p.Add(V2{8, 14})
	p.Add(V2{6, 14})
	p.Add(V2{6, 14.35})
	p.Add(V2{0, 14.35})
	p.Smooth()
	hpipe_3d := Revolve3D(Polygon2D(p.Vertices()))
	// bolt circle for the side flanges
	side_holes_3d := MakeBoltCircle3D(
		2.0,      // hole_depth
		1.0/2.0,  // hole_radius
		14.0/2.0, // circle_radius
		4,        // num_holes
	)
	m = RotateZ(DtoR(45))
	m = Translate3d(V3{0, 0, 1.0 + 12.0}).Mul(m)
	side_holes_3d = Transform3D(side_holes_3d, m)
	hpipe_3d = Difference3D(hpipe_3d, side_holes_3d)
	hpipe_3d = Union3D(Transform3D(hpipe_3d, RotateY(DtoR(90))), Transform3D(hpipe_3d, RotateY(DtoR(-90))))
	hpipe_3d = Transform3D(hpipe_3d, Translate3d(V3{0, 0, 9}))

	s := Union3D(hpipe_3d, vpipe_3d)
	s.(*UnionSDF3).SetMin(PolyMin, 1.0)

	// vertical blind hole
	vertical_hole_3d := Cylinder3D(
		19.0,    // height
		9.0/2.0, // radius
		0.0,     // round
	)
	m = Translate3d(V3{0, 0, 19.0/2.0 + 1})
	vertical_hole_3d = Transform3D(vertical_hole_3d, m)

	// horizontal through hole
	horizontal_hole_3d := Cylinder3D(
		28.70,   // height
		9.0/2.0, // radius
		0.0,     // round
	)
	m = RotateY(DtoR(90))
	m = Translate3d(V3{0, 0, 9}).Mul(m)
	horizontal_hole_3d = Transform3D(horizontal_hole_3d, m)

	return Difference3D(s, Union3D(vertical_hole_3d, horizontal_hole_3d))
}

//-----------------------------------------------------------------------------
