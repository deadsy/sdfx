//-----------------------------------------------------------------------------
/*

Spirals

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {
	s, err := sdf.ArcSpiral2D(1.0, 20.0, 0.25*sdf.Pi, 8*sdf.Tau, 1.0)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	s = sdf.NewVoxelSDF2(s, 400, false)
	//s.(*sdf.VoxelSDF2).Populate(nil)
	render.RenderDXF(s, 400, "spiral.dxf")
}

//-----------------------------------------------------------------------------
