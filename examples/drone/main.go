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

// https://www.flashhobby.com/d2830-fixed-wing-motor.html
var kArm = obj.DroneArmParms{
	MotorSize:     v2.Vec{28, 30},      // motor diameter/height
	MotorMount:    v3.Vec{16, 19, 3.4}, // motor mount l0, l1, diameter
	RotorCavity:   v2.Vec{9, 1.5},      // cavity for bottom of rotor
	WallThickness: 3.0,                 // wall thickness
	SideClearance: 1.5,                 // wall to motor clearance
	MountHeight:   0.7,                 // height of motor mount wrt motor height
	ArmHeight:     0.9,                 // height of arm wrt motor mount height
	ArmLength:     70.0,                // length of rotor arm
}

var kSocket = obj.DroneArmSocketParms{
	Arm:       &kArm,              // drone arm parameters
	Size:      v3.Vec{40, 30, 30}, // body size for socket
	Clearance: 0.5,                // clearance between arm and socket
	Stop:      35,                 // depth of arm stop
}

//-----------------------------------------------------------------------------

func main() {

	arm, err := obj.DroneMotorArm(&kArm)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(arm, shrink), "arm.stl", render.NewMarchingCubesOctree(300))

	socket, err := obj.DroneMotorArmSocket(&kSocket)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(socket, shrink), "socket.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
