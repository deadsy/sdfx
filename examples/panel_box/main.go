//-----------------------------------------------------------------------------
/*

Demonstration for Parametric Box/Case

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func box() error {

	bp := sdf.PanelBoxParms{
		Size:       sdf.V3{50.0, 40.0, 60.0}, // width, height, length
		Wall:       2.5,                      // wall thickness
		Panel:      3.0,                      // panel thickness
		Rounding:   5.0,                      // outer corner rounding
		FrontInset: 2.0,                      // inset for front panel
		BackInset:  2.0,                      // inset for pack panel
		Hole:       3.4,                      // #6 screw
		SideTabs:   "TbtbT",                  // tab pattern
	}

	box, err := sdf.PanelBox3D(&bp)
	if err != nil {
		return err
	}

	sdf.RenderSTL(box[0], 300, "panel.stl")
	sdf.RenderSTL(box[1], 300, "top.stl")
	sdf.RenderSTL(box[2], 300, "bottom.stl")
	return nil
}

//-----------------------------------------------------------------------------

func main() {
	err := box()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

//-----------------------------------------------------------------------------
