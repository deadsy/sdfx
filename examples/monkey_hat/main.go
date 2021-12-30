//-----------------------------------------------------------------------------
/*

Imported monkey model, with modifications

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"os"
)

func main() {

	// MONKEY
	// - Open the STL file
	file, err := os.OpenFile("monkey.stl", os.O_RDONLY, 0400)
	if err != nil {
		file, err = os.OpenFile("examples/monkey_hat/monkey.stl", os.O_RDONLY, 0400)
		if err != nil {
			panic(err)
		}
	}
	// - Create the SDF from the mesh
	monkeyImported, err := obj.ImportSTL(file, 0.2, 5)
	if err != nil {
		panic(err)
	}

	// HAT
	hatHeight := 0.5
	hat, err := sdf.Cylinder3D(hatHeight, 0.6, 0)
	if err != nil {
		panic(err)
	}
	edge, err := sdf.Cylinder3D(hatHeight*0.4, 1, 0)
	if err != nil {
		panic(err)
	}
	edge = sdf.Transform3D(edge, sdf.Translate3d(sdf.V3{Z: -hatHeight / 2}))
	fullHat := sdf.Union3D(hat, edge)

	// Union
	fullHat = sdf.Transform3D(fullHat, sdf.Translate3d(sdf.V3{Y: 0.15, Z: 1}))
	monkeyHat := sdf.Union3D(monkeyImported, fullHat)

	render.ToSTL(monkeyHat, 256, "monkey-out.stl", &render.MarchingCubesOctree{})
	//render.ToSTL(monkeyHat, 16, "monkey-out.stl", dc.NewDualContouringDefault())
}
