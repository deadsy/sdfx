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
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

var module = (5.0 / 8.0) / 20.0
var pa = sdf.DtoR(20.0)
var h = 0.15
var numberTeeth = 20

//-----------------------------------------------------------------------------

func gear() (sdf.SDF3, error) {
	k := obj.InvoluteGearParms{
		NumberTeeth:   numberTeeth,
		Module:        module,
		PressureAngle: pa,
		RingWidth:     0.05,
		Facets:        7,
	}
	gear2d, err := obj.InvoluteGear(&k)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(gear2d, h), nil
}

//-----------------------------------------------------------------------------

func rack() (sdf.SDF3, error) {
	k := sdf.GearRackParms{
		NumberTeeth:   11,
		Module:        module,
		PressureAngle: pa,
		BaseHeight:    0.025,
	}
	rack2d, err := sdf.GearRack2D(&k)
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

	m := sdf.Rotate3d(v3.Vec{0, 0, 1}, sdf.DtoR(180.0/float64(numberTeeth)))
	m = sdf.Translate3d(v3.Vec{0, 0.39, 0}).Mul(m)
	gear = sdf.Transform3D(gear, m)

	render.ToSTL(sdf.Union3D(rack, gear), "gear.stl", render.NewMarchingCubesOctree(200))
}

//-----------------------------------------------------------------------------
