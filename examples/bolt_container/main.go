//-----------------------------------------------------------------------------
/*

Nuts and Bolts

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

const hexRadius = 40.0
const hexHeight = 20.0
const screwRadius = hexRadius * 0.7
const threadPitch = screwRadius / 5.0
const screwLength = 40.0
const tolerance = 0.5

const baseThickness = 4.0

//-----------------------------------------------------------------------------

func boltContainer() sdf.SDF3 {

	// build hex head
	hex := sdf.HexHead3D(hexRadius, hexHeight, "tb")

	// build the screw portion
	r := screwRadius - tolerance
	l := screwLength
	isoThread := sdf.ISOThread(r, threadPitch, true)
	screw := sdf.Screw3D(isoThread, l, threadPitch, 1)
	// chamfer the thread
	screw = sdf.ChamferedCylinder(screw, 0, 0.25)
	screw = sdf.Transform3D(screw, sdf.Translate3d(sdf.V3{0, 0, l / 2}))

	// build the internal cavity
	r = screwRadius * 0.75
	l = screwLength + hexHeight
	round := screwRadius * 0.1
	ofs := (l / 2) - (hexHeight / 2) + baseThickness
	cavity := sdf.Cylinder3D(l, r, round)
	cavity = sdf.Transform3D(cavity, sdf.Translate3d(sdf.V3{0, 0, ofs}))

	return sdf.Difference3D(sdf.Union3D(hex, screw), cavity)
}

//-----------------------------------------------------------------------------

func nutTop() sdf.SDF3 {
	return nil
}

//-----------------------------------------------------------------------------

func main() {
	sdf.RenderSTL(boltContainer(), 200, "container.stl")
}

//-----------------------------------------------------------------------------
