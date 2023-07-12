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
	"strconv"

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
		benchmarkSquare()
	case Circle:
	case Pipe:
	case I:
	default:
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
	return m.WriteInp(pth, []int{0, 1, 2}, 1.25e-9, 900, 0.3, restraint, load)
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
	return m.WriteInpLayers(pth, layerStart, layerEnd, []int{0, 1, 2}, 1.25e-9, 900, 0.3, restraint, load)
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
