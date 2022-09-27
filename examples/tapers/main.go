//-----------------------------------------------------------------------------
/*

Tapered Threads

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func taper1() (sdf.SDF3, error) {

	pitch := 0.50
	radius := 2.0
	length := 5.0
	taper := sdf.DtoR(20)

	isoThread, err := sdf.ISOThread(radius, pitch, true)
	if err != nil {
		return nil, err
	}

	s0, _ := sdf.Screw3D(isoThread, length, taper, pitch, 7)
	s1, _ := sdf.Screw3D(isoThread, length, taper, pitch, -7)

	return sdf.Union3D(s0, s1), nil
}

//-----------------------------------------------------------------------------

func taper2() (sdf.SDF3, error) {

	pitch := 0.50
	radius := 2.0
	length := 10.0
	taper := sdf.DtoR(3)

	isoThread, err := sdf.ISOThread(radius, pitch, true)
	if err != nil {
		return nil, err
	}

	return sdf.Screw3D(isoThread, length, taper, pitch, 1)
}

//-----------------------------------------------------------------------------

func main() {
	s1, err := taper1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s1, "taper1.stl", render.NewMarchingCubesUniform(300))

	s2, err := taper2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s2, "taper2.stl", render.NewMarchingCubesUniform(300))
}

//-----------------------------------------------------------------------------
