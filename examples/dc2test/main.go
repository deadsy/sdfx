//-----------------------------------------------------------------------------
/*

Text Example

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

	f, err := sdf.LoadFont("../../files/cmr10.ttf")
	//f, err := sdf.LoadFont("Times_New_Roman.ttf")
	//f, err := sdf.LoadFont("wt064.ttf")

	if err != nil {
		log.Fatalf("can't read font file %s\n", err)
	}

	t := sdf.NewText("hi!")
	//t := sdf.NewText("相同的不同")

	s0, err := sdf.Text2D(f, t, 10.0)
	if err != nil {
		log.Fatalf("can't generate text %s\n", err)
	}

	//render.ToDXF(s2d, "output.dxf", render.NewMarchingSquaresQuadtree(600))
	render.ToDXF(s0, "output.dxf", render.NewDualContouring2D(50))
}

//-----------------------------------------------------------------------------
