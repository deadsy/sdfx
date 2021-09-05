//-----------------------------------------------------------------------------
/*

Nuts and Bolts

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

const hexRadius = 40.0
const hexHeight = 20.0
const screwRadius = hexRadius * 0.7
const threadPitch = screwRadius / 5.0
const screwLength = 40.0
const tolerance = 0.5

const baseThickness = 4.0

//-----------------------------------------------------------------------------

func boltContainer() (sdf.SDF3, error) {

	// build hex head
	hex, err := obj.HexHead3D(hexRadius, hexHeight, "tb")
	if err != nil {
		return nil, err
	}
	// build the screw portion
	r := screwRadius - tolerance
	l := screwLength
	isoThread, err := sdf.ISOThread(r, threadPitch, true)
	if err != nil {
		return nil, err
	}
	screw, err := sdf.Screw3D(isoThread, l, 0, threadPitch, 1)
	if err != nil {
		return nil, err
	}
	// chamfer the thread
	screw, err = obj.ChamferedCylinder(screw, 0, 0.25)
	if err != nil {
		return nil, err
	}
	screw = sdf.Transform3D(screw, sdf.Translate3d(sdf.V3{0, 0, l / 2}))

	// build the internal cavity
	r = screwRadius * 0.75
	l = screwLength + hexHeight
	round := screwRadius * 0.1
	ofs := (l / 2) - (hexHeight / 2) + baseThickness
	cavity, err := sdf.Cylinder3D(l, r, round)
	if err != nil {
		return nil, err
	}
	cavity = sdf.Transform3D(cavity, sdf.Translate3d(sdf.V3{0, 0, ofs}))

	return sdf.Difference3D(sdf.Union3D(hex, screw), cavity), nil
}

//-----------------------------------------------------------------------------

func nutTop() sdf.SDF3 {
	return nil
}

//-----------------------------------------------------------------------------

func main() {
	bc, err := boltContainer()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(bc, 200, "container.stl")
}

//-----------------------------------------------------------------------------
