//-----------------------------------------------------------------------------
/*

Pipe Flange with a Square base

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var pipe_clearance = 1.02                 // ID of pipe holder slightly larger to allow a slip fit
var pipe_diameter = 48.5 * pipe_clearance // OD of pipe to be fitted
var base_size = V2{60.0, 70.0}            // size of rectangular base
var base_thickness = 3.0                  // base thickness
var pipe_wall = 3.0                       // pipe holder wall thickness
var pipe_length = 20.0                    // length of pipe holder (from bottom)
var pipe_offset = V2{0, 4.0}              // offset of pipe holder from base center

var pipe_radius = pipe_diameter / 2.0
var pipe_fillet = pipe_wall * 0.85

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func flange() SDF3 {

	// base
	pp := &PanelParms{
		Size:         base_size,
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	base := Extrude3D(Panel2D(pp), 2.0*base_thickness)

	// outer pipe
	outer_pipe := Cylinder3D(2.0*pipe_length, pipe_radius+pipe_wall, 0.0)
	outer_pipe = Transform3D(outer_pipe, Translate3d(pipe_offset.ToV3(0)))

	// inner pipe
	inner_pipe := Cylinder3D(2.0*pipe_length, pipe_radius, 0.0)
	inner_pipe = Transform3D(inner_pipe, Translate3d(pipe_offset.ToV3(0)))

	// combine the outer pipe and base (with a fillet)
	s0 := Union3D(base, outer_pipe)
	s0.(*UnionSDF3).SetMin(PolyMin(pipe_fillet))

	// cut the through hole
	s := Difference3D(s0, inner_pipe)

	// return the upper half
	return Cut3D(s, V3{0, 0, 0}, V3{0, 0, 1})
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(Scale3D(flange(), shrink), 300, "flange.stl")
}

//-----------------------------------------------------------------------------
