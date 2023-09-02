//-----------------------------------------------------------------------------
/*

Finite elements from triangle mesh.
The result `inp` file is consumable by ABAQUS or CalculiX.

*/
//-----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf/finiteelements/mesh"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type Specs struct {
	MassDensity            float64
	YoungModulus           float64
	PoissonRatio           float64
	GravityDirectionX      float64
	GravityDirectionY      float64
	GravityDirectionZ      float64
	GravityMagnitude       float64
	Resolution             int  // Number of voxels on the longest axis of 3D model AABB.
	LayersAllConsidered    bool // If true, layer start and layer end are ignored.
	LayerStart             int
	LayerEnd               int
	NonlinearConsidered    bool // If true, nonlinear finite elements are generated.
	ExactSurfaceConsidered bool // If true, surface is approximated by tetrahedral finite elements.
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

type ResultInfo struct {
	ComponentCount int
	Components     []struct{ VoxelCount int }
}

// Render STL to SDF3 to finite elements.
// Write finite elements to an `inp` file.
// Written file can be used by ABAQUS or CalculiX.
func main() {
	if len(os.Args) != 7 {
		log.Fatalf("usage: wrong argument count")
	}

	pthStl := os.Args[1]
	pthSpecs := os.Args[2]
	pthLoads := os.Args[3]
	pthRestraints := os.Args[4]
	pthResult := os.Args[5]
	pthResultInfo := os.Args[6]

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

	err = Run(pthStl, specs, restraints, loads, pthResult, pthResultInfo)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func Run(
	pthStl string,
	specs Specs,
	restraints []Restraint,
	loads []Load,
	pthResult string,
	pthResultInfo string,
) error {
	// create the SDF from the STL mesh
	inSdf, err := obj.ImportSTL(pthStl, 20, 3, 5)
	if err != nil {
		log.Fatalf("%s", err)
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
	m, _ := mesh.NewFem(inSdf, render.NewMarchingCubesFeUniform(specs.Resolution, order, shape))

	components := m.Components()
	ri := ResultInfo{
		ComponentCount: len(components),
		Components:     make([]struct{ VoxelCount int }, len(components)),
	}
	for i, component := range components {
		ri.Components[i] = struct{ VoxelCount int }{VoxelCount: component.VoxelCount()}
	}

	jsonData, err := json.MarshalIndent(components, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = os.WriteFile(pthResultInfo, jsonData, 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if specs.LayersAllConsidered {
		// Write all layers of mesh to file.
		return m.WriteInp(pthResult,
			float32(specs.MassDensity), float32(specs.YoungModulus), float32(specs.PoissonRatio),
			restraintsConvert(restraints),
			loadsConvert(loads),
			v3.Vec{X: specs.GravityDirectionX, Y: specs.GravityDirectionY, Z: specs.GravityDirectionZ}, specs.GravityMagnitude,
		)
	} else {
		// Write just some layers of mesh to file.
		// Generate finite elements.
		// Only from a start layer to an end layer along the Z axis.
		// Applicable to 3D print analysis that is done layer-by-layer.
		return m.WriteInpLayers(
			pthResult,
			specs.LayerStart, specs.LayerEnd,
			float32(specs.MassDensity), float32(specs.YoungModulus), float32(specs.PoissonRatio),
			restraintsConvert(restraints),
			loadsConvert(loads),
			v3.Vec{X: specs.GravityDirectionX, Y: specs.GravityDirectionY, Z: specs.GravityDirectionZ}, specs.GravityMagnitude,
		)
	}
}

func restraintsConvert(rs []Restraint) []*mesh.Restraint {
	restraints := make([]*mesh.Restraint, len(rs))
	for i, r := range rs {
		restraint := mesh.NewRestraint([]v3.Vec{{X: r.LocX, Y: r.LocY, Z: r.LocZ}}, r.IsFixedX, r.IsFixedY, r.IsFixedZ)
		restraints[i] = restraint
	}
	return restraints
}

func loadsConvert(ls []Load) []*mesh.Load {
	loads := make([]*mesh.Load, len(ls))
	for i, l := range ls {
		load := mesh.NewLoad([]v3.Vec{{X: l.LocX, Y: l.LocY, Z: l.LocZ}}, v3.Vec{X: l.MagX, Y: l.MagY, Z: l.MagZ})
		loads[i] = load
	}
	return loads
}
