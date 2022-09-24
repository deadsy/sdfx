//-----------------------------------------------------------------------------
/*

nut cover

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
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

func cover() (sdf.SDF3, error) {
	r := (hexRadius(nutFlat2Flat) * nutFit) + wallThickness
	h := recessHeight + wallThickness
	return sdf.Cylinder3D(2*h, r, 0.1*r)
}

func recess() (sdf.SDF3, error) {
	r := hexRadius(nutFlat2Flat) * nutFit
	h := recessHeight
	return obj.HexHead3D(r, 2*h, "")
}

func counterbore() (sdf.SDF3, error) {
	r := counterBoreDiameter * 0.5
	h := counterBoreDepth
	return sdf.Cylinder3D(2*h, r, 0)
}

func nutcover() (sdf.SDF3, error) {
	cover, err := cover()
	if err != nil {
		return nil, err
	}
	recess, err := recess()
	if err != nil {
		return nil, err
	}
	counterbore, err := counterbore()
	if err != nil {
		return nil, err
	}
	cover = sdf.Difference3D(cover, sdf.Union3D(recess, counterbore))
	return sdf.Cut3D(cover, v3.Vec{0, 0, 0}, v3.Vec{0, 0, 1}), nil
}

func main() {
	s, err := nutcover()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	// un-comment for a cut-away view
	//s = sdf.Cut3D(s, v3.Vec{0, 0, 0}, v3.Vec{1, 0, 0})
	render.RenderSTL(s, 150, "cover.stl")
}

//-----------------------------------------------------------------------------
