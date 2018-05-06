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
		Size:       V3{50.0, 40.0, 60.0},
		Wall:       2.5,
		Panel:      3.0,
		Rounding:   5.0,
		FrontInset: 2.0,
		BackInset:  2.0,
		Clearance:  0.05,
		Hole:       2.0,
		SideTabs:   "TbtbT",
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
