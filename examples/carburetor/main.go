//-----------------------------------------------------------------------------
/*

Carburetor/Manifold Block-Off Plates

Block off plates to stop foreign objects getting into intake manifolds and
carburetors.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// Derived from internet sources and measurement of an Edelbrock intake manifold.
const dX = 0.5 * 5.625 * sdf.MillimetresPerInch
const dY0 = 0.5 * 4.25 * sdf.MillimetresPerInch // spreadbore
const dY1 = 0.5 * 5.16 * sdf.MillimetresPerInch // holley
const holeRadius = 0.5 * (5.0 / 16.0) * sdf.MillimetresPerInch
const holeClearance = 1.05

const plateX = (2.0 * dX) + 20.0
const plateY = (2.0 * dY1) + 20.0
const plateZ = 4.0

//-----------------------------------------------------------------------------

func blockOffPlate() (sdf.SDF3, error) {

	// plate
	plate := sdf.Box2D(v2.Vec{plateX, plateY}, 1.0*plateZ)

	// holes
	hole, err := sdf.Circle2D(holeClearance * holeRadius)
	if err != nil {
		return nil, err
	}

	posn := []v2.Vec{
		v2.Vec{dX, dY0},
		v2.Vec{-dX, -dY0},
		v2.Vec{dX, -dY0},
		v2.Vec{-dX, dY0},
		v2.Vec{dX, dY1},
		v2.Vec{-dX, -dY1},
		v2.Vec{dX, -dY1},
		v2.Vec{-dX, dY1},
	}
	holes := sdf.Multi2D(hole, posn)

	return sdf.Extrude3D(sdf.Difference2D(plate, holes), plateZ), nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := blockOffPlate()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "plate.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
