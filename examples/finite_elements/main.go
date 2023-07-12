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
	"math"
	"os"
	"os/exec"
	"strconv"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/mesh"
	"github.com/deadsy/sdfx/sdf"
)

// Declare the enum using iota and const
type Benchmark int

const (
	Square Benchmark = iota + 1
	Circle
	Pipe
	I
	Unknown
)

// Render SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
//
// OpenSCAD must be installed and be available on PATH as `openscad`
func main() {
	benchmark := Square

	// Optional argument from 1 to 4 to specify the benchmark to run.
	if len(os.Args) > 1 {
		bmint, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		benchmark = Benchmark(bmint)
	}

	switch benchmark {
	case Square:
		benchmarkRun("../../files/benchmark-square.scad", 50, 0, 3, restraintSquare, loadSquare)
	case Circle:
	case Pipe:
	case I:
	default:
	}
}

// Benchmark reference:
// https://github.com/calculix/CalculiX-Examples/tree/master/NonLinear/Sections
func benchmarkRun(
	cad string,
	resolution int,
	layerStart, layerEnd int,
	restraint func(x, y, z float64) (bool, bool, bool),
	load func(x, y, z float64) (float64, float64, float64),
) {
	prg := "openscad"
	stl := "benchmark.stl"
	cmd := exec.Command(prg, "-o", stl, cad)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	fmt.Println(string(stdout))

	// read the stl file.
	file, err := os.OpenFile(stl, os.O_RDONLY, 0400)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// create the SDF from the STL mesh
	inSdf, err := obj.ImportSTL(file, 20, 3, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet4 i.e. 4-node tetrahedron
	err = fe(inSdf, resolution, render.Linear, render.Tetrahedral, "tet4.inp", restraint, load)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet4 i.e. 4-node tetrahedron
	err = fePartial(
		inSdf, resolution, render.Linear, render.Tetrahedral,
		fmt.Sprintf("tet4--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraint, load, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fe(inSdf, resolution, render.Quadratic, render.Tetrahedral, "tet10.inp", restraint, load)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fePartial(
		inSdf, resolution, render.Quadratic, render.Tetrahedral,
		fmt.Sprintf("tet10--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraint, load, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fe(inSdf, resolution, render.Linear, render.Hexahedral, "hex8.inp", restraint, load)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fePartial(inSdf, resolution, render.Linear, render.Hexahedral,
		fmt.Sprintf("hex8--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraint, load, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fe(inSdf, resolution, render.Quadratic, render.Hexahedral, "hex20.inp", restraint, load)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fePartial(inSdf, resolution, render.Quadratic, render.Hexahedral,
		fmt.Sprintf("hex20--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraint, load, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fe(inSdf, resolution, render.Linear, render.Both, "hex8tet4.inp", restraint, load)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fePartial(inSdf, resolution, render.Linear, render.Both,
		fmt.Sprintf("hex8tet4--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraint, load, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fe(inSdf, resolution, render.Quadratic, render.Both, "hex20tet10.inp", restraint, load)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fePartial(inSdf, resolution, render.Quadratic, render.Both,
		fmt.Sprintf("hex20tet10--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraint, load, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// Generate finite elements.
func fe(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string,
	restraint func(x, y, z float64) (bool, bool, bool),
	load func(x, y, z float64) (float64, float64, float64),
) error {
	s = dilationErosion(s)

	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))

	// Write all layers of mesh to file.
	return m.WriteInp(pth, 1.25e-9, 900, 0.3, restraint, load)
}

// Generate finite elements.
// Only from a start layer to an end layer along the Z axis.
// Applicable to 3D print analysis that is done layer-by-layer.
func fePartial(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string,
	restraint func(x, y, z float64) (bool, bool, bool),
	load func(x, y, z float64) (float64, float64, float64),
	layerStart, layerEnd int,
) error {
	s = dilationErosion(s)

	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))

	// Write just some layers of mesh to file.
	return m.WriteInpLayers(pth, layerStart, layerEnd, 1.25e-9, 900, 0.3, restraint, load)
}

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
