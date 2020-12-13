//-----------------------------------------------------------------------------
/*

Involute Gear and Gear Rack

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

var module = (5.0 / 8.0) / 20.0
var pa = sdf.DtoR(20.0)
var h = 0.15
var numberTeeth = 20

//-----------------------------------------------------------------------------

func gear() (sdf.SDF3, error) {
	gear2d, err := obj.InvoluteGear(
		numberTeeth, // number of gear teeth
		module,      // gear_module
		pa,          // pressure_angle
		0.0,         // backlash
		0.0,         // clearance
		0.05,        // ring_width
		7,           // facets
	)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(gear2d, h), nil
}

//-----------------------------------------------------------------------------

func rack() (sdf.SDF3, error) {
	rack2d, err := sdf.GearRack2D(
		11,     // number of rack teeth
		module, // gear_module
		pa,     // pressure_angle
		0.00,   // backlash
		0.025,  // base_height
	)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(rack2d, h), nil
}

//-----------------------------------------------------------------------------

func main() {
	gear, err := gear()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	rack, err := rack()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	m := sdf.Rotate3d(sdf.V3{0, 0, 1}, sdf.DtoR(180.0/float64(numberTeeth)))
	m = sdf.Translate3d(sdf.V3{0, 0.39, 0}).Mul(m)
	gear = sdf.Transform3D(gear, m)

	render.RenderSTL(sdf.Union3D(rack, gear), 200, "gear.stl")
}

//-----------------------------------------------------------------------------
