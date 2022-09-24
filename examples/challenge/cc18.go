//-----------------------------------------------------------------------------
/*
CAD Challenge #18
https://www.reddit.com/r/cad/comments/5vwdnc/cad_challenge_18/
*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// Part A

func cc18a() {
	p := sdf.NewPolygon()
	// start at the top left corner
	p.Add(0, 0)
	p.Add(175, sdf.DtoR(-15)).Polar().Rel()
	p.Add(130, 0).Rel()
	p.Add(0, -25).Rel()
	p.Add(80, 0).Rel()
	p.Add(0, 25).Rel()
	p.Add(75, 0).Rel()
	p.Add(0, -75).Rel()
	p.Add(115, sdf.DtoR(-105)).Polar().Rel()
	p.Add(-50, 0).Rel()
	p.Add(150, sdf.DtoR(-195)).Polar().Rel().Arc(-120, 15)
	p.Add(100, sdf.DtoR(-150)).Polar().Rel()
	p.Add(-60, 0).Rel()
	p.Add(-10, 0).Rel()
	p.Add(-30, 0).Rel()
	p.Add(0, 135).Rel()
	p.Add(-60, 0).Rel()
	// back to the the start with a closed polygon
	p.Close()
	render.Poly(p, "cc18a.dxf")
}

//-----------------------------------------------------------------------------
// Part B

func cc18b() (sdf.SDF3, error) {

	// build the vertical pipe
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(6, 0)
	p.Add(6, 19).Smooth(0.5, 5)
	p.Add(8, 19)
	p.Add(8, 21)
	p.Add(6, 21)
	p.Add(6, 20)
	p.Add(0, 20)
	s2, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	vpipe_3d, err := sdf.Revolve3D(s2)
	if err != nil {
		return nil, err
	}

	// bolt circle for the top flange
	top_holes_3d, err := obj.BoltCircle3D(
		2.0,       // hole_depth
		0.5/2.0,   // hole_radius
		14.50/2.0, // circle_radius
		6,         // num_holes
	)
	if err != nil {
		return nil, err
	}
	m := sdf.RotateZ(sdf.DtoR(30))
	m = sdf.Translate3d(v3.Vec{0, 0, 1.0 + 19.0}).Mul(m)
	top_holes_3d = sdf.Transform3D(top_holes_3d, m)
	vpipe_3d = sdf.Difference3D(vpipe_3d, top_holes_3d)

	// build the horizontal pipe
	p = sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(5, 0)
	p.Add(5, 12).Smooth(0.5, 5)
	p.Add(8, 12)
	p.Add(8, 14)
	p.Add(6, 14)
	p.Add(6, 14.35)
	p.Add(0, 14.35)
	s2, err = sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	hpipe_3d, err := sdf.Revolve3D(s2)
	if err != nil {
		return nil, err
	}
	// bolt circle for the side flanges
	side_holes_3d, err := obj.BoltCircle3D(
		2.0,      // hole_depth
		1.0/2.0,  // hole_radius
		14.0/2.0, // circle_radius
		4,        // num_holes
	)
	if err != nil {
		return nil, err
	}
	m = sdf.RotateZ(sdf.DtoR(45))
	m = sdf.Translate3d(v3.Vec{0, 0, 1.0 + 12.0}).Mul(m)
	side_holes_3d = sdf.Transform3D(side_holes_3d, m)
	hpipe_3d = sdf.Difference3D(hpipe_3d, side_holes_3d)
	hpipe_3d = sdf.Union3D(sdf.Transform3D(hpipe_3d, sdf.RotateY(sdf.DtoR(90))), sdf.Transform3D(hpipe_3d, sdf.RotateY(sdf.DtoR(-90))))
	hpipe_3d = sdf.Transform3D(hpipe_3d, sdf.Translate3d(v3.Vec{0, 0, 9}))

	s := sdf.Union3D(hpipe_3d, vpipe_3d)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(1.0))

	// vertical blind hole
	vertical_hole_3d, _ := sdf.Cylinder3D(
		19.0,    // height
		9.0/2.0, // radius
		0.0,     // round
	)
	m = sdf.Translate3d(v3.Vec{0, 0, 19.0/2.0 + 1})
	vertical_hole_3d = sdf.Transform3D(vertical_hole_3d, m)

	// horizontal through hole
	horizontal_hole_3d, _ := sdf.Cylinder3D(
		28.70,   // height
		9.0/2.0, // radius
		0.0,     // round
	)
	m = sdf.RotateY(sdf.DtoR(90))
	m = sdf.Translate3d(v3.Vec{0, 0, 9}).Mul(m)
	horizontal_hole_3d = sdf.Transform3D(horizontal_hole_3d, m)

	return sdf.Difference3D(s, sdf.Union3D(vertical_hole_3d, horizontal_hole_3d)), nil
}

//-----------------------------------------------------------------------------
// Part C

func cc18c() (sdf.SDF3, error) {

	// build the tabs
	tab_3d, err := sdf.Box3D(v3.Vec{43, 12, 20}, 0)
	if err != nil {
		return nil, err
	}
	tab_3d = sdf.Transform3D(tab_3d, sdf.Translate3d(v3.Vec{43.0 / 2.0, 0, 0}))
	// tab hole
	tab_hole_3d, err := sdf.Cylinder3D(12, 7.0/2.0, 0)
	if err != nil {
		return nil, err
	}
	m := sdf.RotateX(sdf.DtoR(90))
	m = sdf.Translate3d(v3.Vec{35, 0, 0}).Mul(m)
	tab_hole_3d = sdf.Transform3D(tab_hole_3d, m)
	tab_3d = sdf.Difference3D(tab_3d, tab_hole_3d)
	// rotate and copy 3 times
	tab_3d = sdf.RotateCopy3D(tab_3d, 3)

	// Build the ecntral body
	body_3d, err := sdf.Cylinder3D(20, 26.3, 0)
	if err != nil {
		return nil, err
	}
	body_3d = sdf.Union3D(body_3d, tab_3d)
	body_3d.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(2.0))
	// clean up the top and bottom face
	body_3d = sdf.Cut3D(body_3d, v3.Vec{0, 0, -10}, v3.Vec{0, 0, 1})
	body_3d = sdf.Cut3D(body_3d, v3.Vec{0, 0, 10}, v3.Vec{0, 0, -1})

	// build the central sleeve
	r_outer := 42.3 / 2.0
	p := []v2.Vec{
		{0, 0},
		{r_outer, 0},
		{r_outer, 29},
		{r_outer - 1.0, 30},
		{0, 30},
	}
	s, err := sdf.Polygon2D(p)
	if err != nil {
		return nil, err
	}
	sleeve_3d, err := sdf.Revolve3D(s)
	if err != nil {
		return nil, err
	}
	sleeve_3d = sdf.Transform3D(sleeve_3d, sdf.Translate3d(v3.Vec{0, 0, -10}))
	body_3d = sdf.Union3D(body_3d, sleeve_3d)

	// Remove the central hole
	sleeve_hole_3d, err := sdf.Cylinder3D(30, 36.5/2.0, 0)
	if err != nil {
		return nil, err
	}
	sleeve_hole_3d = sdf.Transform3D(sleeve_hole_3d, sdf.Translate3d(v3.Vec{0, 0, 5}))
	body_3d = sdf.Difference3D(body_3d, sleeve_hole_3d)

	return body_3d, nil
}

//-----------------------------------------------------------------------------
