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
var thread_pitch = 6.0

//var thread_diameter = 48.0 // tight
var thread_diameter = 48.5 // just right
//var thread_diameter = 49.0 // loose
var thread_radius = thread_diameter / 2.0

//-----------------------------------------------------------------------------

func cap_outer() SDF3 {
	return KnurledHead3D(cap_radius, cap_height, cap_radius*0.25)
}

func cap_inner() SDF3 {
	tp := PlasticButtressThread(thread_radius, thread_pitch)
	screw := Screw3D(tp, cap_height, thread_pitch, 1)
	return Transform3D(screw, Translate3d(V3{0, 0, -cap_thickness}))
}

func gas_cap() SDF3 {
	return Difference3D(cap_outer(), cap_inner())
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTLSlow(gas_cap(), 300, "cap.stl")
}

//-----------------------------------------------------------------------------
