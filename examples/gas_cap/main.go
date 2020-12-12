//-----------------------------------------------------------------------------
/*

Replacement Cap for Plastic Gas/Oil Can

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

const capRadius = 56.0 / 2.0
const capHeight = 28.0
const capThickness = 4.0
const threadPitch = 6.0
const holeRadius = 0.0 // 33.0 / 2.0

//var threadDiameter = 48.0 // tight
const threadDiameter = 48.5 // just right
//var threadDiameter = 49.0 // loose
const threadRadius = threadDiameter / 2.0

//-----------------------------------------------------------------------------

func capOuter() sdf.SDF3 {
	s, _ := obj.KnurledHead3D(capRadius, capHeight, capRadius*0.25)
	return s
}

func capInner() sdf.SDF3 {
	tp := sdf.PlasticButtressThread(threadRadius, threadPitch)
	screw := sdf.Screw3D(tp, capHeight, threadPitch, 1)
	return sdf.Transform3D(screw, sdf.Translate3d(sdf.V3{0, 0, -capThickness}))
}

func capHole() (sdf.SDF3, error) {
	if holeRadius == 0 {
		// no hole
		return nil, nil
	}
	return sdf.Cylinder3D(capHeight, holeRadius, 0)
}

func gasCap() sdf.SDF3 {
	hole, _ := capHole()
	inner := sdf.Union3D(capInner(), hole)
	return sdf.Difference3D(capOuter(), inner)
}

//-----------------------------------------------------------------------------

func main() {
	render.RenderSTLSlow(gasCap(), 200, "cap.stl")
}

//-----------------------------------------------------------------------------
