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
	PathStl                        string // Input STL file.
	PathLoadPoints                 string // File containing point loads.
	PathRestraintPoints            string // File containing point restraints.
	PathResult                     string // Result file, consumable by ABAQUS or CalculiX.
	PathResultInfo                 string // Result details and info.
	MassDensity                    float64
	YoungModulus                   float64
	PoissonRatio                   float64
	GravityDirectionX              float64
	GravityDirectionY              float64
	GravityDirectionZ              float64
	GravityMagnitude               float64
	Resolution                     int    // Number of voxels on the longest axis of 3D model AABB.
	LayerByLayerfor3dPrintAnalysis bool   // If true, multiple results will be created layer-by-layer to simulated 3D print process.
	LayerByLayerPathResult         string // Only relevant if layer-by-layer is true. Must include "#" character as placeholder for layer number.
	NonlinearConsidered            bool   // If true, nonlinear finite elements are generated.
	ExactSurfaceConsidered         bool   // If true, surface is approximated by tetrahedral finite elements.
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

	jsonData, err = os.ReadFile(specs.PathLoadPoints)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var loads []Load
	err = json.Unmarshal(jsonData, &loads)
	if err != nil {
		log.Fatalf(err.Error())
	}

	jsonData, err = os.ReadFile(specs.PathRestraintPoints)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var restraints []Restraint
	err = json.Unmarshal(jsonData, &restraints)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = Run(specs, restraints, loads)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func Run(
	specs Specs,
	restraints []Restraint,
	loads []Load,
) error {
	// create the SDF from the STL mesh
	inSdf, err := obj.ImportSTL(specs.PathStl, 20, 3, 5)
	if err != nil {
		return err
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
	ri := ResultInfo{
		VoxelsX:        voxelsX,
		VoxelsY:        voxelsY,
		VoxelsZ:        voxelsZ,
		ComponentCount: len(components),
		Components:     make([]Component, len(components)),
	}
	for i, component := range components {
		ri.Components[i] = Component{VoxelCount: component.VoxelCount()}
	}

	jsonData, err := json.MarshalIndent(ri, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(specs.PathResultInfo, jsonData, 0644)
	if err != nil {
		return err
	}

	if !specs.LayerByLayerfor3dPrintAnalysis {
		// Write all layers of mesh to file.
		return m.WriteInp(specs.PathResult,
			float32(specs.MassDensity), float32(specs.YoungModulus), float32(specs.PoissonRatio),
			restraintsConvert(restraints),
			loadsConvert(loads),
			v3.Vec{X: specs.GravityDirectionX, Y: specs.GravityDirectionY, Z: specs.GravityDirectionZ}, specs.GravityMagnitude,
		)
	} else {
		// Write layer by layer.
		// From first layer to the last layer along the Z axis.
		// Applicable to 3D print analysis that is done layer-by-layer.

		if voxelsZ < 8 {
			return fmt.Errorf("not enough voxel layers along the Z axis %d", voxelsZ)
		}

		// The first few layers are ignored.
		for z := 3; z < voxelsZ; z++ {
			err = m.WriteInpLayers(
				strings.Replace(specs.LayerByLayerPathResult, "#", fmt.Sprintf("%d", z), 1),
				0, z,
				float32(specs.MassDensity), float32(specs.YoungModulus), float32(specs.PoissonRatio),
				restraintsConvert(restraints),
				loadsConvert(loads),
				v3.Vec{X: specs.GravityDirectionX, Y: specs.GravityDirectionY, Z: specs.GravityDirectionZ}, specs.GravityMagnitude,
			)
			if err != nil {
				return err
			}
			fmt.Printf("Finite elements are generated from layer 0 to layer %v out of %v total.\n", z, voxelsZ-1)
		}
	}

	return nil
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

// TODO: For 3D print analysis, all the voxels at the same Z layer are considered as restraint.
// Since, the 3D print floor is at the same Z level.
