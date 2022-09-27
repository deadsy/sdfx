//-----------------------------------------------------------------------------
/*

Extrusions

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func hex() (sdf.SDF2, error) {
	return sdf.Polygon2D(sdf.Nagon(6, 20))
}

func extrude1() (sdf.SDF3, error) {

	h, err := hex()
	if err != nil {
		return nil, err
	}

	// make the extrusions
	sLinear := sdf.Extrude3D(sdf.Offset2D(h, 8), 100)
	sFwd := sdf.TwistExtrude3D(sdf.Offset2D(h, 8), 100, sdf.Tau)
	sRev := sdf.TwistExtrude3D(sdf.Offset2D(h, 8), 100, -sdf.Tau)
	sCombo := sdf.Union3D(sFwd, sRev)

	// position them on the y-axis
	d := 60.0
	sLinear = sdf.Transform3D(sLinear, sdf.Translate3d(v3.Vec{0, -1.5 * d, 0}))
	sFwd = sdf.Transform3D(sFwd, sdf.Translate3d(v3.Vec{0, -0.5 * d, 0}))
	sRev = sdf.Transform3D(sRev, sdf.Translate3d(v3.Vec{0, 0.5 * d, 0}))
	sCombo = sdf.Transform3D(sCombo, sdf.Translate3d(v3.Vec{0, 1.5 * d, 0}))

	// return a union of them all
	return sdf.Union3D(sLinear, sFwd, sRev, sCombo), nil
}

func extrude2() (sdf.SDF3, error) {

	h, err := hex()
	if err != nil {
		return nil, err
	}

	s0 := sdf.ScaleExtrude3D(sdf.Offset2D(h, 8), 80, v2.Vec{0.25, 0.5})
	s1 := sdf.ScaleTwistExtrude3D(sdf.Offset2D(h, 8), 80, sdf.Pi, v2.Vec{0.25, 0.5})

	// position them on the y-axis
	d := 30.0
	s0 = sdf.Transform3D(s0, sdf.Translate3d(v3.Vec{0, -d, 0}))
	s1 = sdf.Transform3D(s1, sdf.Translate3d(v3.Vec{0, d, 0}))

	return sdf.Union3D(s0, s1), nil
}

//-----------------------------------------------------------------------------

func main() {
	ex1, err := extrude1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(ex1, "extrude1.stl", render.NewMarchingCubesUniform(200))

	ex2, err := extrude2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(ex2, "extrude2.stl", render.NewMarchingCubesUniform(200))
}

//-----------------------------------------------------------------------------
