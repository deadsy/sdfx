//-----------------------------------------------------------------------------
/*

Drone Parts

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// material shrinkage

const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func arm() (sdf.SDF3, error) {

	// https://www.flashhobby.com/d2830-fixed-wing-motor.html
	k := obj.DroneArmParms{
		MotorSize:     v2.Vec{28, 30},      // motor diameter/height
		MotorMount:    v3.Vec{16, 19, 3.4}, // motor mount l0, l1, diameter
		RotorCavity:   v2.Vec{9, 1.5},      // cavity for bottom of rotor
		WallThickness: 3.0,                 // wall thickness
		SideClearance: 1.5,                 // wall to motor clearance
		MountHeight:   0.7,                 // height of motor mount wrt motor height
		ArmHeight:     0.9,                 // height of arm wrt motor mount height
		ArmLength:     70.0,                // length of rotor arm
	}

	return obj.DroneMotorArm(&k)
}

//-----------------------------------------------------------------------------

func main() {

	s, err := arm()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "arm.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
