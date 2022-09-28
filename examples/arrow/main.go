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
	v3 "github.com/deadsy/sdfx/vec/v3"
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
	return obj.Axes3D(v3.Vec{-10, -10, -10}, v3.Vec{10, 20, 20})
}

func axes2() (sdf.SDF3, error) {
	return obj.Axes3D(v3.Vec{-10, -20, -30}, v3.Vec{0, 0, 0})
}

func axes3() (sdf.SDF3, error) {
	return obj.Axes3D(v3.Vec{0, 0, 0}, v3.Vec{500, 500, 1000})
}

//-----------------------------------------------------------------------------

func main() {
	arrow1, err := arrow1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(arrow1, "arrow1.stl", render.NewMarchingCubesOctree(300))

	axes1, err := axes1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(axes1, "axes1.stl", render.NewMarchingCubesOctree(300))

	axes2, err := axes2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(axes2, "axes2.stl", render.NewMarchingCubesOctree(300))

	axes3, err := axes3()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(axes3, "axes3.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
