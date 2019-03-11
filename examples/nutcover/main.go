//-----------------------------------------------------------------------------
/*

nut cover

*/
//-----------------------------------------------------------------------------

package main

import (
	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

const nutFit = 1.01        // press fit on nut
const wallThickness = 1.15 // wall thickness wrt nut radius

func nutcover(name string, h float64) SDF3 {
	// nut
	r0 := ThreadLookup(name).HexRadius() * nutFit
	h0 := h
	nut := HexHead3D(r0, 2*h0, "tb")

	// cover
	r1 := wallThickness * r0
	h1 := h * wallThickness
	cover := Cylinder3D(2*h1, r1, 0.1*r1)

	return Cut3D(Difference3D(cover, nut), V3{0, 0, 0}, V3{0, 0, 1})
}

func main() {
	RenderSTL(nutcover("M64x4", 60.0), 300, "cover.stl")
}

//-----------------------------------------------------------------------------
