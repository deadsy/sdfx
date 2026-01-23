//-----------------------------------------------------------------------------
/*

Radio Parts

Variable Capacitor: https://a.co/d/hFRjz4D

Ferrite Rod: https://a.co/d/c1uaYZN

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

func ferriteMount() (sdf.SDF3, error) {

	const rodRadius = 10.4 * 0.5
	const baseSize = 20.0
	const rodHeight = 25.0
	const WallThickness = 3.0
	const holderDepth = 6.0
	const holderRadius = WallThickness + rodRadius
	const holderLength = holderDepth + WallThickness

	// support wall
	wall2d, err := obj.IsocelesTriangle2D(baseSize, rodHeight)
	if err != nil {
		return nil, err
	}
	wall2d = sdf.Offset2D(wall2d, holderRadius)
	wall := sdf.Extrude3D(wall2d, WallThickness)

	// base
	const baseX = baseSize + 2.0*holderRadius
	const baseY = holderRadius
	const baseZ = 20.0
	base, err := sdf.Box3D(v3.Vec{baseX, baseY, baseZ}, 0)
	if err != nil {
		return nil, err
	}
	yOfs := -0.5 * (baseY + rodHeight)
	zOfs := 0.5 * (baseZ - WallThickness)
	base = sdf.Transform3D(base, sdf.Translate3d(v3.Vec{0, yOfs, zOfs}))

	// holder
	holder, err := sdf.Cylinder3D(holderLength, holderRadius, 0)
	if err != nil {
		return nil, err
	}
	rodHole, err := sdf.Cylinder3D(holderDepth, rodRadius, 0)
	if err != nil {
		return nil, err
	}
	zOfs = 0.5 * (holderLength - holderDepth)
	rodHole = sdf.Transform3D(rodHole, sdf.Translate3d(v3.Vec{0, 0, zOfs}))
	holder = sdf.Difference3D(holder, rodHole)

	// move the holder
	yOfs = rodHeight * 0.5
	zOfs = 0.5 * (holderLength - WallThickness)
	holder = sdf.Transform3D(holder, sdf.Translate3d(v3.Vec{0, yOfs, zOfs}))

	// cut off the excess base
	fm := sdf.Union3D(base, wall, holder)
	yOfs = -0.5*rodHeight - WallThickness
	fm = sdf.Cut3D(fm, v3.Vec{0, yOfs, 0}, v3.Vec{0, 1, 0})

	return fm, nil
}

//-----------------------------------------------------------------------------

const screwHoleRadius = 3.7 * 0.5
const shaftRadius = 8.0 * 0.5

func vcapMountHole(length float64) (sdf.SDF3, error) {
	// screw holes for mounting
	const screwOffset = 14.0 * 0.5
	sh, err := obj.ChamferedHole3D(length, screwHoleRadius, 0.5)
	if err != nil {
		return nil, err
	}
	h0 := sdf.Transform3D(sh, sdf.Translate3d(v3.Vec{screwOffset, 0, 0}))
	h1 := sdf.Transform3D(sh, sdf.Translate3d(v3.Vec{-screwOffset, 0, 0}))
	// shaft hole
	h2, err := sdf.Cylinder3D(length, shaftRadius+0.4, 0)
	if err != nil {
		return nil, err
	}
	return sdf.Union3D(h0, h1, h2), nil
}

func vcapShaftHole(length float64) (sdf.SDF3, error) {

	// tip for variable cpacitor shaft
	const tipRadius = 6.7 * 0.5
	const tipFlatToFlat = 4.6
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

const mountThickness = 5.0

func vcapKnob() (sdf.SDF3, error) {

	const knobRadius = 40.0 * 0.5
	const knobWidth = 15.0
	const shaftLength = mountThickness - 1.3

	knob, err := sdf.Cylinder3D(knobWidth, knobRadius, 2.0)
	if err != nil {
		return nil, err
	}

	knurl, err := obj.KnurledHead3D(knobRadius, knobWidth*0.67, 3.0)
	if err != nil {
		return nil, err
	}

	knob = sdf.Union3D(knob, knurl)

	totalLength := knobWidth + shaftLength
	shaft, err := sdf.Cylinder3D(totalLength, shaftRadius, 0)
	if err != nil {
		return nil, err
	}
	zOfs := 0.5 * shaftLength
	shaft = sdf.Transform3D(shaft, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	hole, err := vcapShaftHole(totalLength)
	if err != nil {
		return nil, err
	}
	hole = sdf.Transform3D(hole, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	return sdf.Difference3D(sdf.Union3D(knob, shaft), hole), nil
}

func vcapMount() (sdf.SDF3, error) {
	const length = 45.0

	mount, err := sdf.Box3D(v3.Vec{length, length, mountThickness}, 0)
	if err != nil {
		return nil, err
	}

	holes, err := vcapMountHole(mountThickness)
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(mount, holes), nil
}

//-----------------------------------------------------------------------------

func main() {

	vcapKnob, err := vcapKnob()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(vcapKnob, "vc_knob.stl", render.NewMarchingCubesUniform(500))

	vcapMount, err := vcapMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(vcapMount, "vc_mount.stl", render.NewMarchingCubesOctree(500))

	ferriteMount, err := ferriteMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(ferriteMount, "fr_mount.stl", render.NewMarchingCubesOctree(500))

}

//-----------------------------------------------------------------------------
