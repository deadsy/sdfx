//-----------------------------------------------------------------------------
/*

Finite elements from triangle mesh.
Output `inp` file is consumable by ABAQUS or CalculiX.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
)

func main() {
	stl := "../../files/teapot.stl"

	// read the stl file.
	file, err := os.OpenFile(stl, os.O_RDONLY, 0400)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// create the SDF from the STL mesh
	teapotSdf, err := obj.ImportSTL(file, 20, 3, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// Render SDF3 to finite elements.
	r := render.NewMarchingTet4Uniform(200)
	tet4s := render.ToTet4(teapotSdf, r)

	// Create a mesh out of finite elements.
	_, _, layerCountZ := r.LayerCounts(teapotSdf)
	m := render.NewMeshTet4(layerCountZ, tet4s)

	// Write mesh to a file.
	// Written file can be used by ABAQUS or CalculiX.
	err = m.WriteInp("teapot.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// Write just some layers of mesh to a file.
	err = m.WriteInpLayers("teapot-some-layers.inp", 10, 21)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
