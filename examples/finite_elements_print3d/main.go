//-----------------------------------------------------------------------------
/*

Finite elements from triangle mesh.
The result `inp` file is consumable by ABAQUS or CalculiX.

*/
//-----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf/finiteelements/mesh"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type Specs struct {
	PathStl                   string // Input STL file.
	PathResultWithPlaceholder string // Result file, consumable by ABAQUS or CalculiX. Must include "#" character as placeholder for layer number.
	PathResultInfo            string // Result details and info.
	LayerToStartFea           int    // FEA will be done after this layer.
	MassDensity               float64
	YoungModulus              float64
	PoissonRatio              float64
	TensileStrength           float64
	GravityDirectionX         float64
	GravityDirectionY         float64
	GravityDirectionZ         float64
	GravityMagnitude          float64
	Resolution                int  // Number of voxels on the longest axis of 3D model AABB.
	NonlinearConsidered       bool // If true, nonlinear finite elements are generated.
	ExactSurfaceConsidered    bool // If true, surface is approximated by tetrahedral finite elements.
}

type Restraint struct {
	LocX     float64
	LocY     float64
	LocZ     float64
	IsFixedX bool
	IsFixedY bool
	IsFixedZ bool
}

type Load struct {
	LocX float64
	LocY float64
	LocZ float64
	MagX float64
	MagY float64
	MagZ float64
}

type Component struct {
	VoxelCount int
}

type ResultInfo struct {
	VoxelsX        int
	VoxelsY        int
	VoxelsZ        int
	ComponentCount int
	Components     []Component
}

// Render STL to SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: wrong argument count")
	}

	pthSpecs := os.Args[1]

	jsonData, err := os.ReadFile(pthSpecs)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var specs Specs
	err = json.Unmarshal(jsonData, &specs)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// create the SDF from the STL mesh
	inSdf, err := obj.ImportSTL(specs.PathStl, 20, 3, 5)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var order render.Order
	if specs.NonlinearConsidered {
		order = render.Quadratic
	} else {
		order = render.Linear
	}

	var shape render.Shape
	if specs.ExactSurfaceConsidered {
		shape = render.HexAndTet
	} else {
		shape = render.Hexahedral
	}

	// Create a mesh of finite elements.
	m, voxelsX, voxelsY, voxelsZ := mesh.NewFem(inSdf, render.NewMarchingCubesFeUniform(specs.Resolution, order, shape))

	components := m.Components()
	info := ResultInfo{
		VoxelsX:        voxelsX,
		VoxelsY:        voxelsY,
		VoxelsZ:        voxelsZ,
		ComponentCount: len(components),
		Components:     make([]Component, len(components)),
	}
	for i, component := range components {
		info.Components[i] = Component{VoxelCount: component.VoxelCount()}
	}

	jsonData, err = json.MarshalIndent(info, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = os.WriteFile(specs.PathResultInfo, jsonData, 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Generate finite elements layer-by-layer.
	// Applicable to 3D print analysis that is done layer-by-layer.

	if voxelsZ < specs.LayerToStartFea+3 {
		log.Fatalf("not enough voxel layers along the Z axis %d < %d", voxelsZ, specs.LayerToStartFea+3)
	}

	// The first few layers are ignored.
	for z := specs.LayerToStartFea; z <= voxelsZ; z++ {
		err := m.WriteInpLayers(
			strings.Replace(specs.PathResultWithPlaceholder, "#", fmt.Sprintf("%d", z), 1),
			0, z, // Note that the start layer is included, the end layer is excluded.
			float32(specs.MassDensity), float32(specs.YoungModulus), float32(specs.PoissonRatio),
			restraintsPrintFloor(m),
			[]*mesh.Load{}, // Load is empty since only gravity is assumed responsible for 3D print collapse.
			v3.Vec{X: specs.GravityDirectionX, Y: specs.GravityDirectionY, Z: specs.GravityDirectionZ}, specs.GravityMagnitude,
			true,
		)
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Finite elements are generated from layer 0 to layer %v out of %v total.\n", z, voxelsZ)
	}

	if err != nil {
		log.Fatalf(err.Error())
	}
}

// For 3D print analysis, all the voxels at the first layer along Z axis are considered as restraint.
// Since, the 3D print floor is at the first Z level.
func restraintsPrintFloor(m *mesh.Fem) []*mesh.Restraint {
	voxels := m.VoxelsOn1stLayerZ()
	restraints := make([]*mesh.Restraint, len(voxels))
	for i, voxel := range voxels {
		restraints[i] = mesh.NewRestraintByVoxel(voxel, true, true, true)
	}
	return restraints
}
