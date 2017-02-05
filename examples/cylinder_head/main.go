//-----------------------------------------------------------------------------
/*

Wallaby Cylinder Head

No draft version for 3d printing and lost-PLA investment casting.

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------
// overall build controls
const casting = false // add allowances, remove machined features

//-----------------------------------------------------------------------------
// scaling

const MM_PER_INCH = 25.4
const desired_scale = 1.25
const al_shrink = 1.0 / 0.99   // ~1%
const pla_shrink = 1.0 / 0.998 //~0.2%
const abs_shrink = 1.0 / 0.995 //~0.5%

// dimension scaling
func dim(x float64) float64 {
	return x * desired_scale * MM_PER_INCH * al_shrink * pla_shrink
}

//-----------------------------------------------------------------------------
// cylinder studs: location, bosses and holes

var stud_hole_radius = dim(1.0 / 16.0)
var stud_boss_radius = dim(3.0 / 16.0)
var stud_hole_dy = dim(11.0 / 16.0)
var stud_hole_dx0 = dim(7.0 / 16.0)
var stud_hole_dx1 = dim(1.066)

var stud_locations = []V2{
	V2{stud_hole_dx0 + stud_hole_dx1, 0},
	V2{stud_hole_dx0 + stud_hole_dx1, stud_hole_dy},
	V2{stud_hole_dx0 + stud_hole_dx1, -stud_hole_dy},
	V2{stud_hole_dx0, stud_hole_dy},
	V2{stud_hole_dx0, -stud_hole_dy},
	V2{-stud_hole_dx0 - stud_hole_dx1, 0},
	V2{-stud_hole_dx0 - stud_hole_dx1, stud_hole_dy},
	V2{-stud_hole_dx0 - stud_hole_dx1, -stud_hole_dy},
	V2{-stud_hole_dx0, stud_hole_dy},
	V2{-stud_hole_dx0, -stud_hole_dy},
}

//-----------------------------------------------------------------------------
// head walls

var head_length = dim(4.30 / 1.25)
var head_width = dim(2.33 / 1.25)
var head_height = dim(7.0 / 8.0)
var head_corner_rounding = dim((5.0 / 32.0) / 1.25)
var head_wall_thickness = dim(0.154)

func head_wall_outer_2d() SDF2 {
	return NewBoxSDF2(V2{head_length, head_width}, head_corner_rounding)
}

func head_wall_inner_2d() SDF2 {
	l := head_length - (2 * head_wall_thickness)
	w := head_width - (2 * head_wall_thickness)
	// TODO add studs
	return NewBoxSDF2(V2{l, w}, 0)
}

func head_wall() SDF3 {
	wall_2d := head_wall_outer_2d()
	//  wall_2d = NewDifferenceSDF2(wall_2d, head_wall_inner_2d())
	return NewExtrudeSDF3(wall_2d, head_height)
}

//-----------------------------------------------------------------------------

func main() {
	head := head_wall()
	RenderSTL(head, "head.stl")
}

//-----------------------------------------------------------------------------
