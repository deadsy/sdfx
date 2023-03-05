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
	"github.com/deadsy/sdfx/sdf"
)

// 4-node tetrahedral elements.
//
// Render SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
func tet4FiniteElements(s sdf.SDF3, resolution int, pth string) error {
	// Create a mesh out of finite elements.
	m, _ := render.NewMeshTet4(s, render.NewMarchingTet4Uniform(resolution))

	// Write mesh to file.
	err := m.WriteInp(pth)
	if err != nil {
		return err
	}

	// Write just some layers of mesh to a file.
	err = m.WriteInpLayers("some-layers-of-"+pth, 0, 32)
	if err != nil {
		return err
	}
	return nil
}

// 8-node hexahedral elements.
//
// Render SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
func hex8FiniteElements(s sdf.SDF3, resolution int, pth string) error {
	// Create a mesh out of finite elements.
	m, _ := render.NewMeshHex8(s, render.NewMarchingHex8Uniform(resolution))

	// Write mesh to file.
	err := m.WriteInp(pth)
	if err != nil {
		return err
	}

	// Write just some layers of mesh to a file.
	err = m.WriteInpLayers("some-layers-of-"+pth, 0, 32)
	if err != nil {
		return err
	}
	return nil
}

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

	err = tet4FiniteElements(teapotSdf, 200, "teapot-tet4.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = hex8FiniteElements(teapotSdf, 200, "teapot-hex8.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
