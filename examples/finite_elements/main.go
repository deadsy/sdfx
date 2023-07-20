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
	"strconv"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/mesh"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
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

	// To be set for each benchmark.
	var restraints []*mesh.Restraint

	// Gravity load is used for benchmarks. So, there is no point load.
	location := []v3.Vec{{X: 0, Y: 0, Z: 0}}
	loads := []*mesh.Load{
		mesh.NewLoad(location, v3.Vec{X: 0, Y: 0, Z: 0}),
	}

	switch benchmark {
	case Square:
		restraints = benchmarkSquareRestraint()
		benchmarkRun("../../files/benchmark-square.stl", 50, 0, 3, restraints, loads)
	case Circle:
		restraints = benchmarkCircleRestraint()
		benchmarkRun("../../files/benchmark-circle.stl", 50, 0, 3, restraints, loads)
	case Pipe:
		benchmarkRun("../../files/benchmark-pipe.stl", 50, 0, 3, restraints, loads)
	case I:
		benchmarkRun("../../files/benchmark-I.stl", 50, 0, 3, restraints, loads)
	default:
	}
}

// Benchmark reference:
// https://github.com/calculix/CalculiX-Examples/tree/master/NonLinear/Sections
func benchmarkRun(
	stl string,
	resolution int,
	layerStart, layerEnd int,
	restraints []*mesh.Restraint,
	loads []*mesh.Load,
) {
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
	err = fe(inSdf, resolution, render.Linear, render.Tetrahedral, "tet4.inp", restraints, loads)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet4 i.e. 4-node tetrahedron
	err = feLayers(
		inSdf, resolution, render.Linear, render.Tetrahedral,
		fmt.Sprintf("tet4--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraints, loads, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fe(inSdf, resolution, render.Quadratic, render.Tetrahedral, "tet10.inp", restraints, loads)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = feLayers(
		inSdf, resolution, render.Quadratic, render.Tetrahedral,
		fmt.Sprintf("tet10--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraints, loads, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fe(inSdf, resolution, render.Linear, render.Hexahedral, "hex8.inp", restraints, loads)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = feLayers(inSdf, resolution, render.Linear, render.Hexahedral,
		fmt.Sprintf("hex8--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraints, loads, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fe(inSdf, resolution, render.Quadratic, render.Hexahedral, "hex20.inp", restraints, loads)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = feLayers(inSdf, resolution, render.Quadratic, render.Hexahedral,
		fmt.Sprintf("hex20--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraints, loads, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fe(inSdf, resolution, render.Linear, render.Both, "hex8tet4.inp", restraints, loads)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = feLayers(inSdf, resolution, render.Linear, render.Both,
		fmt.Sprintf("hex8tet4--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraints, loads, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fe(inSdf, resolution, render.Quadratic, render.Both, "hex20tet10.inp", restraints, loads)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = feLayers(inSdf, resolution, render.Quadratic, render.Both,
		fmt.Sprintf("hex20tet10--layers-%v-to-%v.inp", layerStart, layerEnd),
		restraints, loads, layerStart, layerEnd,
	)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// Generate finite elements.
func fe(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string,
	restraints []*mesh.Restraint,
	loads []*mesh.Load,
) error {
	s = dilationErosion(s)

	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))

	// Write all layers of mesh to file.
	return m.WriteInp(pth, 7.85e-9, 210000, 0.3, restraints, loads, v3.Vec{X: 0, Y: 0, Z: -1}, 9810)
}

// Generate finite elements.
// Only from a start layer to an end layer along the Z axis.
// Applicable to 3D print analysis that is done layer-by-layer.
func feLayers(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string,
	restraints []*mesh.Restraint,
	loads []*mesh.Load,
	layerStart, layerEnd int,
) error {
	s = dilationErosion(s)

	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))

	// Write just some layers of mesh to file.
	return m.WriteInpLayers(pth, layerStart, layerEnd, 7.85e-9, 210000, 0.3, restraints, loads, v3.Vec{X: 0, Y: 0, Z: -1}, 9810)
}

// By dilating SDF a little bit we may actually get rid of
// bad elements like disconnected or improperly connected elements.
// Erode so that SDF returns to its original size, well almost.
//
// TODO: run benchmarks without this to test.
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

func benchmarkSquareRestraint() []*mesh.Restraint {
	restraints := []*mesh.Restraint{}

	locationPinned := []v3.Vec{}
	gap := 1.0
	var y float64
	for y <= 17.32 {
		locationPinned = append(locationPinned, v3.Vec{X: 0, Y: y, Z: 0})
		y += gap
	}

	restraints = append(restraints, mesh.NewRestraint(locationPinned, true, true, true))

	locationRoller := []v3.Vec{}
	y = 0
	for y <= 17.32 {
		locationRoller = append(locationRoller, v3.Vec{X: 200, Y: y, Z: 0})
		y += gap
	}

	restraints = append(restraints, mesh.NewRestraint(locationRoller, false, true, true))

	return restraints
}

func benchmarkCircleRestraint() []*mesh.Restraint {
	restraints := make([]*mesh.Restraint, 0)

	// The pinned end of 3D beam.
	locationPinned := []v3.Vec{
		{X: 0, Y: 0, Z: 0},
		{X: 0, Y: -2.0313, Z: 0.213498},
		{X: 0, Y: -3.97382, Z: 0.844661},
		{X: 0, Y: -5.74266, Z: 1.8659},
		{X: 0, Y: -7.26052, Z: 3.23259},
		{X: 0, Y: -8.46107, Z: 4.885},
		{X: 0, Y: -9.29182, Z: 6.7509},
		{X: 0, Y: 2.0313, Z: 0.213498},
		{X: 0, Y: 3.97382, Z: 0.844661},
		{X: 0, Y: 5.74266, Z: 1.8659},
		{X: 0, Y: 7.26052, Z: 3.23259},
		{X: 0, Y: 8.46107, Z: 4.885},
		{X: 0, Y: 9.29182, Z: 6.7509},
	}

	restraints = append(restraints, mesh.NewRestraint(locationPinned, true, true, true))

	// The roller end of 3D beam.
	locationRoller := []v3.Vec{
		{X: 200, Y: 0, Z: 0},
		{X: 200, Y: -2.0313, Z: 0.213498},
		{X: 200, Y: -3.97382, Z: 0.844661},
		{X: 200, Y: -5.74266, Z: 1.8659},
		{X: 200, Y: -7.26052, Z: 3.23259},
		{X: 200, Y: -8.46107, Z: 4.885},
		{X: 200, Y: -9.29182, Z: 6.7509},
		{X: 200, Y: 2.0313, Z: 0.213498},
		{X: 200, Y: 3.97382, Z: 0.844661},
		{X: 200, Y: 5.74266, Z: 1.8659},
		{X: 200, Y: 7.26052, Z: 3.23259},
		{X: 200, Y: 8.46107, Z: 4.885},
		{X: 200, Y: 9.29182, Z: 6.7509},
	}

	restraints = append(restraints, mesh.NewRestraint(locationRoller, false, true, true))

	return restraints
}
