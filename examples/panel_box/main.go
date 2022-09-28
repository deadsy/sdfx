//-----------------------------------------------------------------------------
/*

Demonstration for Parametric Box/Case

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func box() ([]sdf.SDF3, error) {
	k := obj.PanelBoxParms{
		Size:       v3.Vec{50.0, 40.0, 60.0}, // width, height, length
		Wall:       2.5,                      // wall thickness
		Panel:      3.0,                      // panel thickness
		Rounding:   5.0,                      // outer corner rounding
		FrontInset: 2.0,                      // inset for front panel
		BackInset:  2.0,                      // inset for pack panel
		Hole:       3.4,                      // #6 screw
		SideTabs:   "TbtbT",                  // tab pattern
	}
	return obj.PanelBox3D(&k)
}

//-----------------------------------------------------------------------------

func main() {
	s, err := box()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	render.ToSTL(s[0], "panel.stl", render.NewMarchingCubesOctree(300))
	render.ToSTL(s[1], "top.stl", render.NewMarchingCubesOctree(300))
	render.ToSTL(s[2], "bottom.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
