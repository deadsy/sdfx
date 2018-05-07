//-----------------------------------------------------------------------------
/*

Demonstration for Parametric Box/Case

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func box1() {

	bp := PanelBoxParms{
		Size:       V3{50.0, 40.0, 60.0}, // width, height, length
		Wall:       2.5,                  // wall thickness
		Panel:      3.0,                  // panel thickness
		Rounding:   5.0,                  // outer corner rounding
		FrontInset: 2.0,                  // inset for front panel
		BackInset:  2.0,                  // inset for pack panel
		Hole:       3.4,                  // #6 screw
		SideTabs:   "TbtbT",              // tab pattern
	}

	box := PanelBox3D(&bp)

	RenderSTL(box[0], 300, "panel.stl")
	RenderSTL(box[1], 300, "top.stl")
	RenderSTL(box[2], 300, "bottom.stl")
}

//-----------------------------------------------------------------------------

func main() {
	box1()
}

//-----------------------------------------------------------------------------
