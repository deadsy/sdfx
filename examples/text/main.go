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

	t := sdf.NewText("SDFX!\nHello,\nWorld!")
	//t := sdf.NewText("相同的不同")

	s0, err := sdf.Text2D(f, t, 10.0)
	if err != nil {
		log.Fatalf("can't generate text %s\n", err)
	}

	// cache the sdf for an evaluation speedup
	s0 = sdf.Cache2D(s0)

	render.ToDXF(s0, "shape.dxf", render.NewMarchingSquaresQuadtree(600))
	render.ToSVG(s0, "shape.svg", render.NewMarchingSquaresQuadtree(600))

	s1, err := sdf.ExtrudeRounded3D(s0, 1.0, 0.2)
	if err != nil {
		log.Fatal(err)
	}
	render.ToSTL(s1, "shape.stl", render.NewMarchingCubesOctree(600))
}

//-----------------------------------------------------------------------------
