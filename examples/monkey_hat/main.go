//-----------------------------------------------------------------------------
/*

Imported monkey model, with modifications

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"os"
	"time"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func monkeyWithHat() sdf.SDF3 {
	// MONKEY
	// - Open the STL file
	file, err := os.OpenFile("monkey.stl", os.O_RDONLY, 0400)
	if err != nil {
		file, err = os.OpenFile("examples/monkey_hat/monkey.stl", os.O_RDONLY, 0400)
		if err != nil {
			panic(err)
		}
	}
	// - Create the SDF from the mesh (a modified Suzanne from Blender with 366 faces)
	monkeyImported, err := obj.ImportSTL(file, 20, 3, 5)
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

	// - Cache the mesh full SDF3 hierarchy for faster evaluation (at the cost of initialization time and memory).
	//   It also smooths the mesh a little using trilinear interpolation.
	//   It is actually slower for this mesh (unless meshCells <<< renderer's meshCells), but should be faster for
	//   more complex meshes (with more triangles) or SDF3 hierarchies that take longer to evaluate.
	monkeyHat = sdf.NewVoxelSDF3(monkeyHat, 64, nil) // Use 32 for harder smoothing demo

	return monkeyHat
}

func main() {
	startTime := time.Now()
	monkeyHat := monkeyWithHat()

	render.ToSTL(monkeyHat, 128, "monkey-out.stl", &render.MarchingCubesUniform{})

	// Dual Contouring is very sensitive to noise (produced when close to shared triangle vertices)
	//render.ToSTL(monkeyHat, 64, "monkey-out.stl", dc.NewDualContouringDefault())

	log.Println("Monkey + hat rendered in", time.Since(startTime))
}
