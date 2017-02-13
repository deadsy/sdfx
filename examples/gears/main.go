package main

import . "github.com/deadsy/sdfx/sdf"

func main() {
	g := InvoluteGear(
		20,             // number_teeth
		(5.0/8.0)/20.0, // gear_module
		DtoR(20.0),     // pressure_angle
		0.0,            // backlash
		0.0,            // clearance
		0.05,           // ring_width
		7,              // facets
	)
	RenderSTL(NewExtrudeSDF3(g, 0.1), "gear.stl")
}
