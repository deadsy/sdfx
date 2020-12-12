//-----------------------------------------------------------------------------
/*

Arrow/Axes Example

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

func axes1() (sdf.SDF3, error) {
	return obj.Axes3D(sdf.V3{-10, -10, -10}, sdf.V3{10, 20, 20})
}

func axes2() (sdf.SDF3, error) {
	return obj.Axes3D(sdf.V3{-10, -20, -30}, sdf.V3{0, 0, 0})
}

func axes3() (sdf.SDF3, error) {
	return obj.Axes3D(sdf.V3{0, 0, 0}, sdf.V3{500, 500, 1000})
}

//-----------------------------------------------------------------------------

func main() {
	arrow1, err := arrow1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(arrow1, 300, "arrow1.stl")

	axes1, err := axes1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(axes1, 300, "axes1.stl")

	axes2, err := axes2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(axes2, 300, "axes2.stl")

	axes3, err := axes3()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(axes3, 300, "axes3.stl")

}

//-----------------------------------------------------------------------------
