//-----------------------------------------------------------------------------
/*

Nuts and Bolts

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var hex_radius = 40.0
var hex_height = 20.0
var screw_radius = hex_radius * 0.7
var thread_pitch = screw_radius / 5.0
var screw_length = 40.0
var tolerance = 0.5

var base_thickness = 4.0

//-----------------------------------------------------------------------------

func bolt_container() SDF3 {

	// build hex head
	hex := HexHead3D(hex_radius, hex_height, "tb")

	// build the screw portion
	r := screw_radius - tolerance
	l := screw_length
	screw := Screw3D(ISOThread(r, thread_pitch, "external"), l, thread_pitch, 1)
	// chamfer the thread
	screw = ChamferedCylinder(screw, 0, 0.25)
	screw = Transform3D(screw, Translate3d(V3{0, 0, l / 2}))

	// build the internal cavity
	r = screw_radius * 0.75
	l = screw_length + hex_height
	round := screw_radius * 0.1
	ofs := (l / 2) - (hex_height / 2) + base_thickness
	cavity := Cylinder3D(l, r, round)
	cavity = Transform3D(cavity, Translate3d(V3{0, 0, ofs}))

	return Difference3D(Union3D(hex, screw), cavity)
}

//-----------------------------------------------------------------------------

func nut_top() SDF3 {
	return nil
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(bolt_container(), 200, "container.stl")
}

//-----------------------------------------------------------------------------
