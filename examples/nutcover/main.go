//-----------------------------------------------------------------------------
/*

nut cover

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

const nutFlat2Flat = 19.0        // measured flat 2 flat nut size
const recessHeight = 20.0        // recess within cover
const wallThickness = 2.0        // wall thickness
const counterBoreDiameter = 23.0 // diameter of washer counterbore
const counterBoreDepth = 2.0     // depth of washer counterbore
const nutFit = 1.01              // press fit on nut

//-----------------------------------------------------------------------------

func hexRadius(f2f float64) float64 {
	return f2f / (2.0 * math.Cos(sdf.DtoR(30)))
}

func cover() sdf.SDF3 {
	r := (hexRadius(nutFlat2Flat) * nutFit) + wallThickness
	h := recessHeight + wallThickness
	return sdf.Cylinder3D(2*h, r, 0.1*r)
}

func recess() sdf.SDF3 {
	r := hexRadius(nutFlat2Flat) * nutFit
	h := recessHeight
	s, _ := obj.HexHead3D(r, 2*h, "")
	return s
}

func counterbore() sdf.SDF3 {
	r := counterBoreDiameter * 0.5
	h := counterBoreDepth
	return sdf.Cylinder3D(2*h, r, 0)
}

func nutcover() sdf.SDF3 {
	s0 := cover()
	s1 := sdf.Union3D(recess(), counterbore())
	return sdf.Cut3D(sdf.Difference3D(s0, s1), sdf.V3{0, 0, 0}, sdf.V3{0, 0, 1})
}

func main() {
	s := nutcover()
	// un-comment for a cut-away view
	//s = sdf.Cut3D(s, sdf.V3{0, 0, 0}, sdf.V3{1, 0, 0})
	sdf.RenderSTL(s, 150, "cover.stl")
}

//-----------------------------------------------------------------------------
