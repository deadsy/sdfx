//-----------------------------------------------------------------------------
/*

Replacement Cap for Plastic Gas/Oil Can

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var cap_radius = 56.0 / 2.0
var cap_height = 28.0
var cap_thickness = 4.0
var thread_radius = 45.0 / 2.0
var thread_pitch = 6.0

//-----------------------------------------------------------------------------

func gas_cap() SDF3 {
	tp := ANSIButtressThread(thread_radius, thread_pitch)
	screw := Screw3D(tp, cap_height, thread_pitch, 1)
	screw = Transform3D(screw, Translate3d(V3{0, 0, -cap_thickness}))
	outer := Cylinder3D(cap_height, cap_radius, 0.0)
	return Difference3D(outer, screw)
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(gas_cap(), 300, "cap.stl")
}

//-----------------------------------------------------------------------------
