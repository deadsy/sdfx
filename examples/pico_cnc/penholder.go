//-----------------------------------------------------------------------------
/*

Pen Holder for Path Testing

Inspired by: https://www.thingiverse.com/thing:2625750)

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// pen holder

func penHolder() (sdf.SDF3, error) {

	const holderHeight = 20.0
	const holderWidth = 25.0
	const shaftRadius = 8.0 * 0.5
	const penRadius = 13.0 * 0.5
	const bossDiameter = 6.0

	// spring
	k := &obj.SpringParms{
		Width:         holderWidth,                       // width of spring
		Height:        holderHeight,                      // height of spring (3d only)
		WallThickness: 1,                                 // thickness of wall
		Diameter:      5,                                 // diameter of spring turn
		NumSections:   3,                                 // number of spring sections
		Boss:          [2]float64{2.0 * bossDiameter, 8}, // boss sizes
	}
	spring, err := k.Spring3D()
	if err != nil {
		return nil, err
	}

	// shaft hole
	shaft, err := sdf.Cylinder3D(k.SpringLength(), shaftRadius, 0)
	if err != nil {
		return nil, err
	}
	shaft = sdf.Transform3D(shaft, sdf.RotateY(sdf.DtoR(90)))

	// shaft screw boss
	bossParms := &obj.ThreadedCylinderParms{
		Height:    0.5 * holderHeight,
		Diameter:  bossDiameter,
		Thread:    "unc_8_32",
		Tolerance: 0,
	}
	boss, err := bossParms.Object()
	if err != nil {
		return nil, err
	}
	boss = sdf.Transform3D(boss, sdf.Translate3d(v3.Vec{0, 0, 30}))

	return sdf.Difference3D(sdf.Union3D(spring, boss), shaft), nil
}

//-----------------------------------------------------------------------------
