//-----------------------------------------------------------------------------
/*

Bushing for the Box Joint Jig

https://woodgears.ca/box_joint/jig.html

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// R6-2RS 3/8 x 7/8 x 9/32 bearing
var bearing_outer_od = (7.0 / 8.0) * MM_PER_INCH   // outer diameter of outer race
var bearing_outer_id = 19.0                        // inner diameter of outer race
var bearing_inner_id = (3.0 / 8.0) * MM_PER_INCH   // inner diameter of inner race
var bearing_inner_od = 12.0                        // outer diameter of inner race
var bearing_thickness = (9.0 / 32.0) * MM_PER_INCH // bearing thickness

// Adjust clearance to give good interference fits for the bearing
var clearance = 0.0

//-----------------------------------------------------------------------------

func bushing() SDF3 {

	r0 := 2.3 // radius of central screw
	r1 := (bearing_outer_od + bearing_inner_od) / 2.0
	r2 := (bearing_inner_id / 2) - clearance

	h0 := 3.0 // height of cap
	h1 := h0 + bearing_thickness + 1.0

	p := NewPolygon()
	p.Add(r0, 0)
	p.Add(r1, 0)
	p.Add(r1, h0)
	p.Add(r2, h0)
	p.Add(r2, h1)
	p.Add(r0, h1)
	return Revolve3D(Polygon2D(p.Vertices()))
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(bushing(), 300, "bushing.stl")
}

//-----------------------------------------------------------------------------
