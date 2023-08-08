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
	"strconv"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/mesh"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type Benchmark int

const (
	Square Benchmark = iota + 1
	Circle
	Pipe
	I
)

type ElementConfig int

const (
	Tet4       ElementConfig = iota + 1 // tet4 i.e. 4-node tetrahedron
	Tet10                               // tet10 i.e. 10-node tetrahedron
	Hex8                                // hex8 i.e. 8-node hexahedron
	Hex20                               // hex20 i.e. 20-node hexahedron
	Hex8Tet4                            // hex8 and tet4
	Hex20Tet10                          // hex20 and tet10
)

// Render SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
func main() {
	benchmark := Square
	elementconfig := Hex20Tet10

	// Optional argument from 1 to 4 to specify the benchmark to run.
	if len(os.Args) > 1 {
		bmint, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		benchmark = Benchmark(bmint)
	}

	// Optional argument from 1 to 6 to specify the elements to generate.
	if len(os.Args) > 2 {
		ecint, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		elementconfig = ElementConfig(ecint)
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
		benchmarkRun("../../files/benchmark-square.stl", 50, 0, 3, restraints, loads, elementconfig)
	case Circle:
		restraints = benchmarkCircleRestraint()
		benchmarkRun("../../files/benchmark-circle.stl", 50, 0, 3, restraints, loads, elementconfig)
	case Pipe:
		restraints = benchmarkPipeRestraint()
		// When resolution is `50`, the pipe benchmark misses some necessary finite elements.
		// By incrementing it to `51`, no element is missed.
		// There is a bug: https://github.com/deadsy/sdfx/issues/72
		benchmarkRun("../../files/benchmark-pipe.stl", 50, 0, 3, restraints, loads, elementconfig)
	case I:
		restraints = benchmarkIRestraint()
		benchmarkRun("../../files/benchmark-I.stl", 50, 0, 3, restraints, loads, elementconfig)
	default:
		// If no valid benchmark is picked, run for the teapot.
		benchmarkRun("../../files/teapot.stl", 61, 0, 18, restraints, loads, elementconfig)
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
	elementconfig ElementConfig,
) {
	// create the SDF from the STL mesh
	inSdf, err := obj.ImportSTL(stl, 20, 3, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	switch elementconfig {
	case Tet4:
		{

			err = fe(inSdf, resolution, render.Linear, render.Tetrahedral, "tet4.inp", restraints, loads)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			err = feLayers(
				inSdf, resolution, render.Linear, render.Tetrahedral,
				fmt.Sprintf("tet4--layers-%v-to-%v.inp", layerStart, layerEnd),
				restraints, loads, layerStart, layerEnd,
			)
			if err != nil {
				log.Fatalf("error: %s", err)
			}
		}
	case Tet10:
		{
			err = fe(inSdf, resolution, render.Quadratic, render.Tetrahedral, "tet10.inp", restraints, loads)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			err = feLayers(
				inSdf, resolution, render.Quadratic, render.Tetrahedral,
				fmt.Sprintf("tet10--layers-%v-to-%v.inp", layerStart, layerEnd),
				restraints, loads, layerStart, layerEnd,
			)
			if err != nil {
				log.Fatalf("error: %s", err)
			}
		}
	case Hex8:
		{
			err = fe(inSdf, resolution, render.Linear, render.Hexahedral, "hex8.inp", restraints, loads)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			err = feLayers(inSdf, resolution, render.Linear, render.Hexahedral,
				fmt.Sprintf("hex8--layers-%v-to-%v.inp", layerStart, layerEnd),
				restraints, loads, layerStart, layerEnd,
			)
			if err != nil {
				log.Fatalf("error: %s", err)
			}
		}
	case Hex20:
		{
			err = fe(inSdf, resolution, render.Quadratic, render.Hexahedral, "hex20.inp", restraints, loads)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			err = feLayers(inSdf, resolution, render.Quadratic, render.Hexahedral,
				fmt.Sprintf("hex20--layers-%v-to-%v.inp", layerStart, layerEnd),
				restraints, loads, layerStart, layerEnd,
			)
			if err != nil {
				log.Fatalf("error: %s", err)
			}
		}
	case Hex8Tet4:
		{
			err = fe(inSdf, resolution, render.Linear, render.Both, "hex8tet4.inp", restraints, loads)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			err = feLayers(inSdf, resolution, render.Linear, render.Both,
				fmt.Sprintf("hex8tet4--layers-%v-to-%v.inp", layerStart, layerEnd),
				restraints, loads, layerStart, layerEnd,
			)
			if err != nil {
				log.Fatalf("error: %s", err)
			}
		}
	case Hex20Tet10:
		{
			err = fe(inSdf, resolution, render.Quadratic, render.Both, "hex20tet10.inp", restraints, loads)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			err = feLayers(inSdf, resolution, render.Quadratic, render.Both,
				fmt.Sprintf("hex20tet10--layers-%v-to-%v.inp", layerStart, layerEnd),
				restraints, loads, layerStart, layerEnd,
			)
			if err != nil {
				log.Fatalf("error: %s", err)
			}
		}
	}
}

// Generate finite elements.
func fe(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string,
	restraints []*mesh.Restraint,
	loads []*mesh.Load,
) error {
	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))
	components := m.Components()
	fmt.Printf("components count: %v\n", len(components))
	for i, component := range components {
		fmt.Printf("component %v voxel count: %v\n", i, component.VoxelCount())
	}
	m.CleanDisconnections(components)
	components = m.Components()
	fmt.Printf("components count after clean up: %v\n", len(components))

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
	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFEUniform(resolution, order, shape))
	components := m.Components()
	fmt.Printf("components count: %v\n", len(components))
	for i, component := range components {
		fmt.Printf("component %v voxel count: %v\n", i, component.VoxelCount())
	}
	m.CleanDisconnections(components)
	components = m.Components()
	fmt.Printf("components count after clean up: %v\n", len(components))

	// Write just some layers of mesh to file.
	return m.WriteInpLayers(pth, layerStart, layerEnd, 7.85e-9, 210000, 0.3, restraints, loads, v3.Vec{X: 0, Y: 0, Z: -1}, 9810)
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
		{X: 0, Y: 2.0313, Z: 0.213498},
		{X: 0, Y: 3.97382, Z: 0.844661},
	}

	restraints = append(restraints, mesh.NewRestraint(locationPinned, true, true, true))

	// The roller end of 3D beam.
	locationRoller := []v3.Vec{
		{X: 200, Y: 0, Z: 0},
		{X: 200, Y: -2.0313, Z: 0.213498},
		{X: 200, Y: -3.97382, Z: 0.844661},
		{X: 200, Y: 2.0313, Z: 0.213498},
		{X: 200, Y: 3.97382, Z: 0.844661},
	}

	restraints = append(restraints, mesh.NewRestraint(locationRoller, false, true, true))

	return restraints
}

func benchmarkPipeRestraint() []*mesh.Restraint {
	restraints := make([]*mesh.Restraint, 0)

	// The pinned end of 3D beam.
	locationPinned := []v3.Vec{
		{X: 0, Y: 0, Z: 0},
		{X: 0, Y: -2.0313, Z: 0.213498},
		{X: 0, Y: -3.97382, Z: 0.844661},
		{X: 0, Y: 2.0313, Z: 0.213498},
		{X: 0, Y: 3.97382, Z: 0.844661},
	}

	restraints = append(restraints, mesh.NewRestraint(locationPinned, true, true, true))

	// The roller end of 3D beam.
	locationRoller := []v3.Vec{
		{X: 200, Y: 0, Z: 0},
		{X: 200, Y: -2.0313, Z: 0.213498},
		{X: 200, Y: -3.97382, Z: 0.844661},
		{X: 200, Y: 2.0313, Z: 0.213498},
		{X: 200, Y: 3.97382, Z: 0.844661},
	}

	restraints = append(restraints, mesh.NewRestraint(locationRoller, false, true, true))

	return restraints
}

func benchmarkIRestraint() []*mesh.Restraint {
	restraints := []*mesh.Restraint{}

	locationPinned := []v3.Vec{}
	gap := 1.0
	var y float64
	for y <= 25 {
		locationPinned = append(locationPinned, v3.Vec{X: 0, Y: y, Z: 0})
		y += gap
	}

	restraints = append(restraints, mesh.NewRestraint(locationPinned, true, true, true))

	locationRoller := []v3.Vec{}
	y = 0
	for y <= 25 {
		locationRoller = append(locationRoller, v3.Vec{X: 200, Y: y, Z: 0})
		y += gap
	}

	restraints = append(restraints, mesh.NewRestraint(locationRoller, false, true, true))

	return restraints
}
