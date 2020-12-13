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
	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------
// Part A

func cc18a() {
	p := NewPolygon()
	// start at the top left corner
	p.Add(0, 0)
	p.Add(175, DtoR(-15)).Polar().Rel()
	p.Add(130, 0).Rel()
	p.Add(0, -25).Rel()
	p.Add(80, 0).Rel()
	p.Add(0, 25).Rel()
	p.Add(75, 0).Rel()
	p.Add(0, -75).Rel()
	p.Add(115, DtoR(-105)).Polar().Rel()
	p.Add(-50, 0).Rel()
	p.Add(150, DtoR(-195)).Polar().Rel().Arc(-120, 15)
	p.Add(100, DtoR(-150)).Polar().Rel()
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

func cc18b() (SDF3, error) {

	// build the vertical pipe
	p := NewPolygon()
	p.Add(0, 0)
	p.Add(6, 0)
	p.Add(6, 19).Smooth(0.5, 5)
	p.Add(8, 19)
	p.Add(8, 21)
	p.Add(6, 21)
	p.Add(6, 20)
	p.Add(0, 20)
	vpipe_3d, err := Revolve3D(Polygon2D(p.Vertices()))
	if err != nil {
		return nil, err
	}
	// bolt circle for the top flange
	top_holes_3d := obj.BoltCircle3D(
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
	p = NewPolygon()
	p.Add(0, 0)
	p.Add(5, 0)
	p.Add(5, 12).Smooth(0.5, 5)
	p.Add(8, 12)
	p.Add(8, 14)
	p.Add(6, 14)
	p.Add(6, 14.35)
	p.Add(0, 14.35)
	hpipe_3d, err := Revolve3D(Polygon2D(p.Vertices()))
	if err != nil {
		return nil, err
	}
	// bolt circle for the side flanges
	side_holes_3d := obj.BoltCircle3D(
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
	s.(*UnionSDF3).SetMin(PolyMin(1.0))

	// vertical blind hole
	vertical_hole_3d, _ := Cylinder3D(
		19.0,    // height
		9.0/2.0, // radius
		0.0,     // round
	)
	m = Translate3d(V3{0, 0, 19.0/2.0 + 1})
	vertical_hole_3d = Transform3D(vertical_hole_3d, m)

	// horizontal through hole
	horizontal_hole_3d, _ := Cylinder3D(
		28.70,   // height
		9.0/2.0, // radius
		0.0,     // round
	)
	m = RotateY(DtoR(90))
	m = Translate3d(V3{0, 0, 9}).Mul(m)
	horizontal_hole_3d = Transform3D(horizontal_hole_3d, m)

	return Difference3D(s, Union3D(vertical_hole_3d, horizontal_hole_3d)), nil
}

//-----------------------------------------------------------------------------
// Part C

func cc18c() (SDF3, error) {

	// build the tabs
	tab_3d := Box3D(V3{43, 12, 20}, 0)
	tab_3d = Transform3D(tab_3d, Translate3d(V3{43.0 / 2.0, 0, 0}))
	// tab hole
	tab_hole_3d, err := Cylinder3D(12, 7.0/2.0, 0)
	if err != nil {
		return nil, err
	}
	m := RotateX(DtoR(90))
	m = Translate3d(V3{35, 0, 0}).Mul(m)
	tab_hole_3d = Transform3D(tab_hole_3d, m)
	tab_3d = Difference3D(tab_3d, tab_hole_3d)
	// rotate and copy 3 times
	tab_3d = RotateCopy3D(tab_3d, 3)

	// Build the ecntral body
	body_3d, err := Cylinder3D(20, 26.3, 0)
	if err != nil {
		return nil, err
	}
	body_3d = Union3D(body_3d, tab_3d)
	body_3d.(*UnionSDF3).SetMin(PolyMin(2.0))
	// clean up the top and bottom face
	body_3d = Cut3D(body_3d, V3{0, 0, -10}, V3{0, 0, 1})
	body_3d = Cut3D(body_3d, V3{0, 0, 10}, V3{0, 0, -1})

	// build the central sleeve
	r_outer := 42.3 / 2.0
	p := []V2{
		{0, 0},
		{r_outer, 0},
		{r_outer, 29},
		{r_outer - 1.0, 30},
		{0, 30},
	}
	sleeve_3d, err := Revolve3D(Polygon2D(p))
	if err != nil {
		return nil, err
	}
	sleeve_3d = Transform3D(sleeve_3d, Translate3d(V3{0, 0, -10}))
	body_3d = Union3D(body_3d, sleeve_3d)

	// Remove the central hole
	sleeve_hole_3d, err := Cylinder3D(30, 36.5/2.0, 0)
	if err != nil {
		return nil, err
	}
	sleeve_hole_3d = Transform3D(sleeve_hole_3d, Translate3d(V3{0, 0, 5}))
	body_3d = Difference3D(body_3d, sleeve_hole_3d)

	return body_3d, nil
}

//-----------------------------------------------------------------------------
