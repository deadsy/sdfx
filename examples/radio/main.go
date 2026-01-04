//-----------------------------------------------------------------------------
/*

Radio Parts

Variable Capacitor: https://a.co/d/hFRjz4D

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

const screwHoleRadius = 3.7 * 0.5
const shaftRadius = 8.0 * 0.5

//-----------------------------------------------------------------------------

func vcapMountHole(length float64) (sdf.SDF3, error) {
	// screw holes for mounting
	const screwOffset = 14.0 * 0.5
	sh, err := sdf.Circle2D(screwHoleRadius)
	if err != nil {
		return nil, err
	}
	h0 := sdf.Transform2D(sh, sdf.Translate2d(v2.Vec{screwOffset, 0}))
	h1 := sdf.Transform2D(sh, sdf.Translate2d(v2.Vec{-screwOffset, 0}))
	// shaft hole
	h2, err := sdf.Circle2D(shaftRadius + 0.4)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(sdf.Union2D(h0, h1, h2), length), nil
}

func vcapShaftHole(length float64) (sdf.SDF3, error) {

	// tip for variable cpacitor shaft
	const tipRadius = 6.4 * 0.5
	const tipFlatToFlat = 4.4
	const tipLength = 2.5
	tip, err := sdf.Cylinder3D(tipLength, tipRadius, 0)
	xOfs := tipFlatToFlat * 0.5
	tip = sdf.Cut3D(tip, v3.Vec{xOfs, 0, 0}, v3.Vec{-1, 0, 0})
	tip = sdf.Cut3D(tip, v3.Vec{-xOfs, 0, 0}, v3.Vec{1, 0, 0})
	zOfs := 0.5 * (length - tipLength)
	tip = sdf.Transform3D(tip, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	// countersink
	const csRadius = 8.0 * 0.5
	const csLength = 3.0
	cs, err := sdf.Cylinder3D(csLength, csRadius, 0)
	zOfs = 0.5 * (length - csLength)
	cs = sdf.Transform3D(cs, sdf.Translate3d(v3.Vec{0, 0, -zOfs}))

	// screw hole
	hole, err := sdf.Cylinder3D(length, screwHoleRadius, 0)
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
	shaft, err := sdf.Cylinder3D(shaftLength, shaftRadius, 0)
	if err != nil {
		return nil, err
	}
	zOfs := 0.5 * (shaftLength - knobHeight)
	shaft = sdf.Transform3D(shaft, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	hole, err := vcapShaftHole(shaftLength)
	if err != nil {
		return nil, err
	}
	hole = sdf.Transform3D(hole, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	return sdf.Difference3D(sdf.Union3D(knob, shaft), hole), nil
}

func vcapMount() (sdf.SDF3, error) {
	const length = 40.0
	const thickness = 3.2

	mount, err := sdf.Box3D(v3.Vec{length, length, thickness}, 0)
	if err != nil {
		return nil, err
	}

	holes, err := vcapMountHole(thickness)
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(mount, holes), nil
}

//-----------------------------------------------------------------------------

func main() {

	knob, err := vcapKnob()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(knob, "knob.stl", render.NewMarchingCubesUniform(500))

	mount, err := vcapMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(mount, "mount.stl", render.NewMarchingCubesOctree(500))
}

//-----------------------------------------------------------------------------
