//-----------------------------------------------------------------------------
/*

Finite elements from triangle mesh.
Output `inp` file is consumable by ABAQUS or CalculiX.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/mesh"
	"github.com/deadsy/sdfx/sdf"
)

// Generate finite elements.
func fe(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string, layerStart, layerEnd int) error {
	min := s.BoundingBox().Min
	max := s.BoundingBox().Max

	dimX := (max.X - min.X)
	dimY := (max.Y - min.Y)
	dimZ := (max.Z - min.Z)
	factor := math.Min(dimX, math.Min(dimY, dimZ)) * float64(0.05)

	// By dilating SDF a little bit we may actually get rid of bad elements like disconnected or improperly connected elements.
	dilation := sdf.Offset3D(s, factor)

	// Erode so that SDF returns to its original size, well almost.
	erosion := sdf.Offset3D(dilation, -factor)

	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(erosion, render.NewMarchingCubesFEUniform(resolution, order, shape))

	// Write just some layers of mesh to a file.
	//err := m.WriteInpLayers(pth, layerStart, layerEnd, []int{0, 1, 2}, 1.25e-9, 900, 0.3)
	err := m.WriteInp(pth, []int{0, 1, 2}, 1.25e-9, 900, 0.3)
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

	// tet4 i.e. 4-node tetrahedron
	err = fe(teapotSdf, 80, render.Linear, render.Tetrahedral, "teapot-tet4.inp", 0, 20)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fe(teapotSdf, 80, render.Quadratic, render.Tetrahedral, "teapot-tet10.inp", 0, 20)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fe(teapotSdf, 80, render.Linear, render.Hexahedral, "teapot-hex8.inp", 0, 20)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fe(teapotSdf, 80, render.Quadratic, render.Hexahedral, "teapot-hex20.inp", 0, 20)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fe(teapotSdf, 80, render.Linear, render.Both, "teapot-hex8tet4.inp", 0, 20)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fe(teapotSdf, 80, render.Quadratic, render.Both, "teapot-hex20tet10.inp", 0, 20)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
