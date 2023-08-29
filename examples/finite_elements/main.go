//-----------------------------------------------------------------------------
/*

Finite elements from triangle mesh.
Output `inp` file is consumable by ABAQUS or CalculiX.

*/
//-----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/sdf/finiteelements/mesh"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type Specs struct {
	MassDensity            float64
	YoungModulus           float64
	PoissonRatio           float64
	GravityConsidered      bool
	GravityDirectionX      float64
	GravityDirectionY      float64
	GravityDirectionZ      float64
	GravityMagnitude       float64
	Resolution             int
	LayersAllConsidered    bool
	LayerStart             int
	LayerEnd               int
	NonlinearConsidered    bool
	ExactSurfaceConsidered bool
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

// Render STL to SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
func main() {
	if len(os.Args) != 6 {
		log.Fatalf("usage: wrong argument count")
	}

	pthStl := os.Args[1]
	pthSpecs := os.Args[2]
	pthLoads := os.Args[3]
	pthRestraints := os.Args[4]
	pthResult := os.Args[5]

	jsonData, err := os.ReadFile(pthSpecs)
	if err != nil {
		log.Fatalf("%s", err)
	}

	var specs Specs
	err = json.Unmarshal(jsonData, &specs)
	if err != nil {
		log.Fatalf("%s", err)
	}

	jsonData, err = os.ReadFile(pthLoads)
	if err != nil {
		log.Fatalf("%s", err)
	}

	var loads []Load
	err = json.Unmarshal(jsonData, &loads)
	if err != nil {
		log.Fatalf("%s", err)
	}

	jsonData, err = os.ReadFile(pthRestraints)
	if err != nil {
		log.Fatalf("%s", err)
	}

	var restraints []Restraint
	err = json.Unmarshal(jsonData, &restraints)
	if err != nil {
		log.Fatalf("%s", err)
	}

	Run(pthStl, specs, restraints, loads, pthResult)
}

func Run(
	pthStl string,
	specs Specs,
	restraints []Restraint,
	loads []Load,
	pthResult string,
) {
	// create the SDF from the STL mesh
	_, err := obj.ImportSTL(pthStl, 20, 3, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// Generate finite elements.
func fe(s sdf.SDF3, resolution int, order render.Order, shape render.Shape, pth string,
	restraints []*mesh.Restraint,
	loads []*mesh.Load,
) error {
	// Create a mesh out of finite elements.
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFeUniform(resolution, order, shape))
	components := m.Components()
	fmt.Printf("components count: %v\n", len(components))
	for i, component := range components {
		fmt.Printf("component %v voxel count: %v\n", i, component.VoxelCount())
	}

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
	m, _ := mesh.NewFem(s, render.NewMarchingCubesFeUniform(resolution, order, shape))
	components := m.Components()
	fmt.Printf("components count: %v\n", len(components))
	for i, component := range components {
		fmt.Printf("component %v voxel count: %v\n", i, component.VoxelCount())
	}

	// Write just some layers of mesh to file.
	return m.WriteInpLayers(pth, layerStart, layerEnd, 7.85e-9, 210000, 0.3, restraints, loads, v3.Vec{X: 0, Y: 0, Z: -1}, 9810)
}
