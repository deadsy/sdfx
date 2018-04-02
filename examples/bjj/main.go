//-----------------------------------------------------------------------------
/*

Bushing for the Box Joint Jig

https://woodgears.ca/box_joint/jig.html

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func bushing() SDF3 {

	// R6-2RS 3/8 x 7/8 x 9/32 bearing
	bearing_outer_od := (7.0 / 8.0) * MM_PER_INCH // outer diameter of outer race
	//bearing_outer_id := 19.0                        // inner diameter of outer race
	bearing_inner_id := (3.0 / 8.0) * MM_PER_INCH   // inner diameter of inner race
	bearing_inner_od := 12.0                        // outer diameter of inner race
	bearing_thickness := (9.0 / 32.0) * MM_PER_INCH // bearing thickness

	// Adjust clearance to give good interference fits for the bearing
	clearance := 0.0

	r0 := 2.3 // radius of central screw
	r1 := (bearing_outer_od + bearing_inner_od) / 4.0
	r2 := (bearing_inner_id / 2.0) - clearance

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

var gear_module = 80.0 / 16.0
var pressure_angle = 20.0
var involute_facets = 10

func stacked_gears() SDF3 {

	sg_height := 10.0
	sg0_teeth := 12
	sg1_teeth := 16

	g0_2d := InvoluteGear(
		sg0_teeth,
		gear_module,
		DtoR(pressure_angle),
		0.0,
		0.0,
		20.0, // width of ring wall (from root circle)
		involute_facets,
	)

	g1_2d := InvoluteGear(
		sg1_teeth,
		gear_module,
		DtoR(pressure_angle),
		0.0,
		0.0,
		20.0, // width of ring wall (from root circle)
		involute_facets,
	)

	g0 := Extrude3D(g0_2d, sg_height)
	g1 := Extrude3D(g1_2d, sg_height)

	g0 = Transform3D(g0, Translate3d(V3{0, 0, sg_height / 2.0}))
	g1 = Transform3D(g1, Translate3d(V3{0, 0, -sg_height / 2.0}))

	return Union3D(g0, g1)
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(bushing(), 300, "bushing.stl")
	RenderSTL(stacked_gears(), 300, "gear.stl")
}

//-----------------------------------------------------------------------------
