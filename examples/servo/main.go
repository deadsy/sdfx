//-----------------------------------------------------------------------------
/*

Servo Models

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
)

//-----------------------------------------------------------------------------

func main() {

	k, err := obj.ServoLookup("standard")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	s, err := obj.Servo3D(k)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	render.ToSTL(s, 300, "servo.stl", &render.MarchingCubesOctree{})
}

//-----------------------------------------------------------------------------
