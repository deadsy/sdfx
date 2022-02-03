//-----------------------------------------------------------------------------
/*

Fuel Pump Ring Nut Tool

Many cars have a fuel pump in the tank held in place by a plastic ringnut.
This is a tool for removing them.

This design is for the Mazda 2006 RX-8 (Series1)
Other ring nuts are similar, so feel free to modify.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

const innerDiameter = 132.0
const ringWidth = 20.0
const outerDiameter = innerDiameter + (2.0 * ringWidth)
const ringHeight = 20.0

const numTabs = 20
const tabDepth = 3.0
const tabWidth = 3.0

const sideThickness = 2.5 * tabDepth
const topThickness = 2.0 * tabDepth

// The rx-8 puts an additional tab on the ring
const extraTab = true

//-----------------------------------------------------------------------------

func outerBody() (sdf.SDF3, error) {
	h := (ringHeight + topThickness) * 2.0
	r := (outerDiameter * 0.5) + sideThickness
	round := topThickness * 0.5
	return sdf.Cylinder3D(h, r, round)
}

func innerCavity() (sdf.SDF3, error) {
	h := ringHeight * 2.0
	r := outerDiameter * 0.5
	round := ringHeight * 0.1
	s0, err := sdf.Cylinder3D(h, r, round)
	if err != nil {
		return nil, err
	}
	// central bore
	h = (ringHeight + topThickness) * 2.0
	r = innerDiameter * 0.5
	s1, err := sdf.Cylinder3D(h, r, 0)
	if err != nil {
		return nil, err
	}
	return sdf.Union3D(s0, s1), nil
}

func tab() (sdf.SDF3, error) {
	size := sdf.V3{
		X: tabWidth,
		Y: ringWidth + tabDepth,
		Z: (ringHeight + tabDepth) * 2.0,
	}
	s, err := sdf.Box3D(size, 0)
	if err != nil {
		return nil, err
	}
	yofs := (size.Y + innerDiameter) * 0.5
	s = sdf.Transform3D(s, sdf.Translate3d(sdf.V3{0, yofs, 0}))
	return s, nil
}

func tabs() (sdf.SDF3, error) {
	t, err := tab()
	if err != nil {
		return nil, err
	}

	theta := sdf.Tau / numTabs
	s := sdf.RotateUnion3D(t, numTabs, sdf.Rotate3d(sdf.V3{0, 0, 1}, theta))

	if extraTab {
		et := sdf.Transform3D(t, sdf.Rotate3d(sdf.V3{0, 0, 1}, theta*0.5))
		s = sdf.Union3D(s, et)
	}

	return s, nil
}

func tool() (sdf.SDF3, error) {

	body, err := outerBody()
	if err != nil {
		return nil, err
	}

	cavity, err := innerCavity()
	if err != nil {
		return nil, err
	}

	// add the tabs
	t, err := tabs()
	if err != nil {
		return nil, err
	}
	cavity = sdf.Union3D(cavity, t)

	s := sdf.Difference3D(body, cavity)

	// cut it on the xy plane
	s = sdf.Cut3D(s, sdf.V3{0, 0, 0}, sdf.V3{0, 0, 1})
	return s, nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := tool()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, 300, "tool.stl", &render.MarchingCubesOctree{})
}

//-----------------------------------------------------------------------------
