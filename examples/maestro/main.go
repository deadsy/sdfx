//-----------------------------------------------------------------------------
/*

PCB Mounting Board for a Pololu Mini Maestro 18 Servo Controller

https://www.pololu.com/product/1354

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

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func servoControllerMount() (sdf.SDF3, error) {

	// standoff
	k0 := obj.StandoffParms{
		PillarHeight:   0.5 * sdf.MillimetresPerInch,
		PillarDiameter: 5,
		HoleDepth:      10,
		HoleDiameter:   2.4, // #4 screw
	}
	standoff, err := obj.Standoff3D(&k0)
	if err != nil {
		return nil, err
	}

	// standoffs
	h0 := sdf.V3{-0.45, -0.8, 0.25}.MulScalar(sdf.MillimetresPerInch)
	h1 := sdf.V3{0.05, 0.8, 0.25}.MulScalar(sdf.MillimetresPerInch)
	standoffs := sdf.Multi3D(standoff, []sdf.V3{h0, h1})

	// base
	k1 := obj.PanelParms{
		Size:         sdf.V2{1.1, 1.8}.MulScalar(sdf.MillimetresPerInch),
		CornerRadius: 2,
		HoleDiameter: 2.4, // #4 screw
		HoleMargin:   [4]float64{4, 4, 4, 4},
		HolePattern:  [4]string{"x", "x", ".x", ""},
		Thickness:    3,
	}
	base, err := obj.Panel3D(&k1)
	if err != nil {
		return nil, err
	}

	return sdf.Union3D(base, standoffs), nil
}

//-----------------------------------------------------------------------------

func main() {

	s, err := servoControllerMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, 300, "mm18.stl", &render.MarchingCubesOctree{})

}

//-----------------------------------------------------------------------------
