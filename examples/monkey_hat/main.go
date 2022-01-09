//-----------------------------------------------------------------------------
/*

Import an existing STL. Modify it. Re-render.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func monkeyWithHat() (sdf.SDF3, error) {
	// read the monkey head stl file.
	file, err := os.OpenFile("monkey.stl", os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}

	// create the SDF from the mesh (a modified Suzanne from Blender with 366 faces)
	monkeyImported, err := obj.ImportSTL(file, 20, 3, 5)
	if err != nil {
		return nil, err
	}

	// build the hat
	hatHeight := 0.5
	hat, err := sdf.Cylinder3D(hatHeight, 0.6, 0)
	if err != nil {
		return nil, err
	}
	edge, err := sdf.Cylinder3D(hatHeight*0.4, 1, 0)
	if err != nil {
		return nil, err
	}
	edge = sdf.Transform3D(edge, sdf.Translate3d(sdf.V3{Z: -hatHeight / 2}))
	fullHat := sdf.Union3D(hat, edge)

	// put the hat on the monkey
	fullHat = sdf.Transform3D(fullHat, sdf.Translate3d(sdf.V3{Y: 0.15, Z: 1}))
	monkeyHat := sdf.Union3D(monkeyImported, fullHat)

	// Cache the mesh full SDF3 hierarchy for faster evaluation (at the cost of initialization time and memory).
	// It also smooths the mesh a little using trilinear interpolation.
	// It is actually slower for this mesh (unless meshCells <<< renderer's meshCells), but should be faster for
	// more complex meshes (with more triangles) or SDF3 hierarchies that take longer to evaluate.
	monkeyHat = sdf.NewVoxelSDF3(monkeyHat, 64, nil)

	return monkeyHat, nil
}

//-----------------------------------------------------------------------------

func main() {
	monkeyHat, err := monkeyWithHat()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(monkeyHat, 128, "monkey-out.stl", &render.MarchingCubesUniform{})
	//render.ToSTL(monkeyHat, 128, "monkey-out.stl", &render.MarchingCubesOctree{})
	//render.ToSTL(monkeyHat, 64, "monkey-out.stl", dc.NewDualContouringDefault())
}

//-----------------------------------------------------------------------------
