//-----------------------------------------------------------------------------
/*

Involute Gear and Gear Rack

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {

	module := (5.0 / 8.0) / 20.0
	pa := sdf.DtoR(20.0)
	h := 0.15
	numberTeeth := 20

	gear2d := obj.InvoluteGear(
		numberTeeth, // number of gear teeth
		module,      // gear_module
		pa,          // pressure_angle
		0.0,         // backlash
		0.0,         // clearance
		0.05,        // ring_width
		7,           // facets
	)
	gear3d := sdf.Extrude3D(gear2d, h)
	m := sdf.Rotate3d(sdf.V3{0, 0, 1}, sdf.DtoR(180.0/float64(numberTeeth)))
	m = sdf.Translate3d(sdf.V3{0, 0.39, 0}).Mul(m)
	gear3d = sdf.Transform3D(gear3d, m)

	rack2d := sdf.GearRack2D(
		11,     // number of rack teeth
		module, // gear_module
		pa,     // pressure_angle
		0.00,   // backlash
		0.025,  // base_height
	)
	rack3d := sdf.Extrude3D(rack2d, h)

	render.RenderSTL(sdf.Union3D(rack3d, gear3d), 200, "gear.stl")
}

//-----------------------------------------------------------------------------
