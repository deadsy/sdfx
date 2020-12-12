//-----------------------------------------------------------------------------
/*

Arrow

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func arrow1() (sdf.SDF3, error) {
	k := obj.ArrowParms{
		Axis:  [2]float64{50, 1},
		Head:  [2]float64{5, 2},
		Tail:  [2]float64{5, 2},
		Style: "cb",
	}
	return obj.Arrow3D(&k)
}

//-----------------------------------------------------------------------------

func main() {
	arrow1, err := arrow1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(arrow1, 300, "arrow1.stl")
}

//-----------------------------------------------------------------------------
