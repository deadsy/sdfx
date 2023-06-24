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

// By dilating SDF a little bit we may actually get rid of
// bad elements like disconnected or improperly connected elements.
// Erode so that SDF returns to its original size, well almost.
func dilationErosion(s sdf.SDF3) sdf.SDF3 {
	min := s.BoundingBox().Min
	max := s.BoundingBox().Max

	dimX := (max.X - min.X)
	dimY := (max.Y - min.Y)
	dimZ := (max.Z - min.Z)

	// What percent is preferred? Calibration is done a bit.
	factor := math.Min(dimX, math.Min(dimY, dimZ)) * float64(0.02)

	dilation := sdf.Offset3D(s, factor)
	erosion := sdf.Offset3D(dilation, -factor)

	return erosion
}

// Generate finite elements.
func fe(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string) error {
	s = dilationErosion(s)

	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))

	// Write all layers of mesh to file.
	return m.WriteInp(pth, []int{0, 1, 2}, 1.25e-9, 900, 0.3)
}

// Generate finite elements.
// Only from a start layer to an end layer along the Z axis.
// Applicable to 3D print analysis that is done layer-by-layer.
func fePartial(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string, layerStart, layerEnd int) error {
	s = dilationErosion(s)

	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))

	// Write just some layers of mesh to file.
	return m.WriteInpLayers(pth, layerStart, layerEnd, []int{0, 1, 2}, 1.25e-9, 900, 0.3)
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
	err = fe(hingeSdf, 50, render.Linear, render.Tetrahedral, "hinge-tet4.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet4 i.e. 4-node tetrahedron
	err = fePartial(hingeSdf, 50, render.Linear, render.Tetrahedral, "hinge-partial-tet4.inp", 0, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fe(hingeSdf, 50, render.Quadratic, render.Tetrahedral, "hinge-tet10.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fePartial(hingeSdf, 50, render.Quadratic, render.Tetrahedral, "hinge-partial-tet10.inp", 0, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fe(hingeSdf, 50, render.Linear, render.Hexahedral, "hinge-hex8.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fePartial(hingeSdf, 50, render.Linear, render.Hexahedral, "hinge-partial-hex8.inp", 0, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fe(hingeSdf, 50, render.Quadratic, render.Hexahedral, "hinge-hex20.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fePartial(hingeSdf, 50, render.Quadratic, render.Hexahedral, "hinge-partial-hex20.inp", 0, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fe(hingeSdf, 50, render.Linear, render.Both, "hinge-hex8tet4.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fePartial(hingeSdf, 50, render.Linear, render.Both, "hinge-partial-hex8tet4.inp", 0, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fe(hingeSdf, 50, render.Quadratic, render.Both, "hinge-hex20tet10.inp")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fePartial(hingeSdf, 50, render.Quadratic, render.Both, "hinge-partial-hex20tet10.inp", 0, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
