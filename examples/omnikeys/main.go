//-----------------------------------------------------------------------------
/*

OmniKeys 3 x 12 chord keys

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var button_r1 = (7.0 / 2.0)  // minor radius
var button_r2 = (12.0 / 2.0) // major radius

var panel_h1 = 1.5 // minor thickness
var panel_h2 = 6.5 // major thickness

//-----------------------------------------------------------------------------

func button_cavity() SDF3 {
	p := NewPolygon()
	p.Add(0, 0)
	p.Add(button_r1, 0)
	p.Add(button_r1, panel_h1)
	p.Add(button_r2, panel_h1)
	p.Add(button_r2, panel_h2)
	p.Add(0, panel_h2)
	return Revolve3D(Polygon2D(p.Vertices()))
}

//-----------------------------------------------------------------------------

func panel() SDF3 {
	pp := &PanelParms{
		Size:         V2{30.0, 30.0},
		CornerRadius: 5.0,
		HoleDiameter: 0,
	}
	// extrude to 3d
	p := Extrude3D(Panel2D(pp), panel_h2)
	p = Transform3D(p, Translate3d(V3{0, 0, panel_h2 * 0.5}))
	return p
}

//-----------------------------------------------------------------------------

func keys() SDF3 {
	return Difference3D(panel(), button_cavity())
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL_New(keys(), 150, "keys.stl")
}

//-----------------------------------------------------------------------------
