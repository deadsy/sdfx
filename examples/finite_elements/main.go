//-----------------------------------------------------------------------------
/*

Finite elements from triangle mesh.
Output `inp` file is consumable by ABAQUS or CalculiX.

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/mesh"
	"github.com/deadsy/sdfx/sdf"
)

// 4-node tetrahedral elements.
func tet4(s sdf.SDF3, resolution int, pth string) error {
	// Create a mesh out of finite elements.
	m, _ := mesh.NewTet4(s, render.NewMarchingCubesFEUniform(resolution))

	lyrStart := 0
	lyrEnd := 20

	// Write just some layers of mesh to a file.
	err := m.WriteInpLayers(pth, lyrStart, lyrEnd, []int{0, 1, 2}, 1.25e-9, 900, 0.3)
	if err != nil {
		return err
	}

	fmt.Println("Tet4 FE is not implemented, there is TODO in function mcToTet4")

	return nil
}

// 8-node hexahedral elements.
func hex8(s sdf.SDF3, resolution int, pth string) error {
	// Create a mesh out of finite elements.
	m, _ := mesh.NewHex8(s, render.NewMarchingCubesFEUniform(resolution))

	lyrStart := 0
	lyrEnd := 20

	// Write just some layers of mesh to a file.
	//
	// Units are mm,N,sec.
	// Force per area = N/mm2 or MPa
	// Mass density = Ns2/mm4
	// Refer to the "Units" chapter of:
	// http://www.dhondt.de/ccx_2.20.pdf
	//
	// Mechanical properties are based on typical SLA resins.
	//
	// TODO: Correct resin specifications.
	err := m.WriteInpLayers(pth, lyrStart, lyrEnd, []int{0, 1, 2}, 1.25e-9, 900, 0.3)
	if err != nil {
		return err
	}
	return nil
}

// 20-node hexahedral elements.
func hex20(s sdf.SDF3, resolution int, pth string) error {
	// Create a mesh out of finite elements.
	m, _ := mesh.NewHex20(s, render.NewMarchingCubesFEUniform(resolution))

	lyrStart := 0
	lyrEnd := 20

	// Write just some layers of mesh to a file.
	//
	// Units are mm,N,sec.
	// Force per area = N/mm2 or MPa
	// Mass density = Ns2/mm4
	// Refer to the "Units" chapter of:
	// http://www.dhondt.de/ccx_2.20.pdf
	//
	// Mechanical properties are based on typical SLA resins.
	//
	// TODO: Correct resin specifications.
	err := m.WriteInpLayers(pth, lyrStart, lyrEnd, []int{0, 1, 2}, 1.25e-9, 900, 0.3)
	if err != nil {
		return err
	}
	return nil
}

// 8-node hexahedral elements.
// With adaptive mesh refinement.
func hex8adaptive(s sdf.SDF3, resolution int, pth string) error {
	// Create a mesh out of finite elements.
	m, _ := mesh.NewHex8(s, render.NewMarchingCubesFEOctree(resolution))

	err := m.WriteInp(pth, []int{0}, 1.25e-9, 900, 0.3)
	if err != nil {
		return err
	}
	return nil
}

// Render SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
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

	err = tet4(teapotSdf, 80, "teapot-tet4.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = hex8(teapotSdf, 80, "teapot-hex8.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = hex20(teapotSdf, 80, "teapot-hex20.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = hex8adaptive(teapotSdf, 80, "teapot-hex20-adaptive.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
