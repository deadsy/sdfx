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
	"github.com/deadsy/sdfx/render/finiteelements/mesh"
	"github.com/deadsy/sdfx/sdf"
)

// Generate finite elements.
func fe(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string, layerStart, layerEnd int) error {
	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))

	// Write just some layers of mesh to a file.
	err := m.WriteInpLayers(pth, layerStart, layerEnd, []int{0, 1, 2}, 1.25e-9, 900, 0.3)
	if err != nil {
		return err
	}

	return nil
}

// Render SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
func main() {
	stl := "../../files/hinge.stl"

	// read the stl file.
	file, err := os.OpenFile(stl, os.O_RDONLY, 0400)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// create the SDF from the STL mesh
	hingeSdf, err := obj.ImportSTL(file, 20, 3, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet4 i.e. 4-node tetrahedron
	err = fe(hingeSdf, 80, render.Linear, render.Tetrahedral, "hinge-tet4.inp", 0, 10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fe(hingeSdf, 80, render.Quadratic, render.Tetrahedral, "hinge-tet10.inp", 0, 10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fe(hingeSdf, 80, render.Linear, render.Hexahedral, "hinge-hex8.inp", 0, 10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fe(hingeSdf, 80, render.Quadratic, render.Hexahedral, "hinge-hex20.inp", 0, 10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fe(hingeSdf, 80, render.Linear, render.Both, "hinge-hex8tet4.inp", 0, 10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fe(hingeSdf, 80, render.Quadratic, render.Both, "hinge-hex20tet10.inp", 0, 10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
