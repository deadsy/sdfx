//-----------------------------------------------------------------------------
/*

Radio Parts

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

func vcapHole(length float64) (sdf.SDF3, error) {

	// tip for variable cpacitor shaft
	const tipRadius = 6.3 * 0.5
	const tipFlatToFlat = 4.0
	const tipLength = 2.5
	tip, err := sdf.Cylinder3D(tipLength, tipRadius, 0)
	xOfs := tipFlatToFlat * 0.5
	tip = sdf.Cut3D(tip, v3.Vec{xOfs, 0, 0}, v3.Vec{-1, 0, 0})
	tip = sdf.Cut3D(tip, v3.Vec{-xOfs, 0, 0}, v3.Vec{1, 0, 0})
	zOfs := 0.5 * (length - tipLength)
	tip = sdf.Transform3D(tip, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	// countersink
	const csRadius = 7.4 * 0.5
	const csLength = 3.0
	cs, err := sdf.Cylinder3D(csLength, csRadius, 0)
	zOfs = 0.5 * (length - csLength)
	cs = sdf.Transform3D(cs, sdf.Translate3d(v3.Vec{0, 0, -zOfs}))

	// screw hole
	const holeRadius = 3.7 * 0.5
	hole, err := sdf.Cylinder3D(length, holeRadius, 0)
	if err != nil {
		return nil, err
	}

	return sdf.Union3D(hole, tip, cs), nil
}

func vcapKnob() (sdf.SDF3, error) {

	const knobRadius = 40.0 * 0.5
	const knobHeight = 22.0

	knob, err := obj.KnurledHead3D(knobRadius, knobHeight, 3.0)
	if err != nil {
		return nil, err
	}

	const shaftLength = knobHeight + 8.0
	const shaftRadius = 8.0
	shaft, err := sdf.Cylinder3D(shaftLength, shaftRadius, 0)
	if err != nil {
		return nil, err
	}
	zOfs := 0.5 * (shaftLength - knobHeight)
	shaft = sdf.Transform3D(shaft, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	hole, err := vcapHole(shaftLength)
	if err != nil {
		return nil, err
	}
	hole = sdf.Transform3D(hole, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	return sdf.Difference3D(sdf.Union3D(knob, shaft), hole), nil
}

//-----------------------------------------------------------------------------

func main() {

	knob, err := vcapKnob()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(knob, "knob.stl", render.NewMarchingCubesUniform(500))

}

//-----------------------------------------------------------------------------
