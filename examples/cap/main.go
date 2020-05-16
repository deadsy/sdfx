//-----------------------------------------------------------------------------
/*

Tube Cap

This is a simple round cap that fits onto the outside of a tube.

*/
//-----------------------------------------------------------------------------

package main

import "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

const wallThickness = 2.0
const innerDiameter = 75.5
const innerHeight = 15.0

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func tubeCap() sdf.SDF3 {

	h := innerHeight + wallThickness
	r := (innerDiameter * 0.5) + wallThickness
	outer := sdf.Cylinder3D(h, r, 1.0)

	h = innerHeight
	r = innerDiameter * 0.5
	inner := sdf.Cylinder3D(h, r, 1.0)
	inner = sdf.Transform3D(inner, sdf.Translate3d(sdf.V3{0, 0, wallThickness * 0.5}))

	return sdf.Difference3D(outer, inner)
}

//-----------------------------------------------------------------------------

func main() {
	sdf.RenderSTL(sdf.ScaleUniform3D(tubeCap(), shrink), 120, "cap.stl")
}

//-----------------------------------------------------------------------------
