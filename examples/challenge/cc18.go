//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------
// CAD Challenge #18 Part B
// https://www.reddit.com/r/cad/comments/5vwdnc/cad_challenge_18/

func cc18b() SDF3 {

	// build the vertical pipe
	// V2{0,0}
	// V2{6,0}
	// V2{0,19}
	// V2{2,0}
	// V2{0,2}
	// V2{-2,0}
	// V2{0,-1}
	// V2{-6,0}

	// bolt circle for the top flange
	top_holes_3d := MakeBoltCircle3D(
		2.0,       // hole_depth
		0.5/2.0,   // hole_radius
		14.50/2.0, // circle_radius
		6,         // num_holes
	)

	// build the horizontal pipe

	// bolt circle for the side flanges
	side_holes_3d := MakeBoltCircle3D(
		2.0,      // hole_depth
		1.0/2.0,  // hole_radius
		14.0/2.0, // circle_radius
		4,        // num_holes
	)

	// vertical blind hole
	vertical_hole_3d := Cylinder3D(
		19.0,    // height
		9.0/2.0, // radius
		0.0,     // round
	)

	// horizontal through hole
	horizontal_hole_3d := Cylinder3D(
		28.70,   // height
		9.0/2.0, // radius
		0.0,     // round
	)

	_ = top_holes_3d
	_ = side_holes_3d
	_ = vertical_hole_3d
	_ = horizontal_hole_3d

	return nil
}

//-----------------------------------------------------------------------------
