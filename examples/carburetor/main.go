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
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// air intake cover: derived from measurement of an Edelbrock carburetor.

const airIntakeRadius = 0.5 * 5.125 * sdf.MillimetresPerInch
const airIntakeWall = (3.0 / 16.0) * sdf.MillimetresPerInch
const airIntakeHeight = 1.375 * sdf.MillimetresPerInch
const airIntakeHole = 0.5 * (5.0 / 16.0) * sdf.MillimetresPerInch

func airIntakeCover() (sdf.SDF3, error) {

	const h0 = 2.0 * (airIntakeHeight + airIntakeWall)
	const r0 = airIntakeRadius + airIntakeWall
	body, err := sdf.Cylinder3D(h0, r0, 0.75*airIntakeWall)
	if err != nil {
		return nil, err
	}

	const h1 = 2.0 * airIntakeHeight
	const r1 = airIntakeRadius
	cavity, err := sdf.Cylinder3D(h1, r1, 0)
	if err != nil {
		return nil, err
	}

	const h2 = h0
	const r2 = airIntakeHole
	hole, err := sdf.Cylinder3D(h2, r2, 0)
	if err != nil {
		return nil, err
	}

	cover := sdf.Difference3D(body, sdf.Union3D(cavity, hole))
	return sdf.Cut3D(cover, v3.Vec{0, 0, 0}, v3.Vec{0, 0, 1}), nil
}

//-----------------------------------------------------------------------------
// manifold blockoff plate: derived from measurement of an Edelbrock intake manifold.

const dX = 0.5 * 5.625 * sdf.MillimetresPerInch
const dY0 = 0.5 * 4.25 * sdf.MillimetresPerInch // spreadbore
const dY1 = 0.5 * 5.16 * sdf.MillimetresPerInch // holley
const holeRadius = 0.5 * (5.0 / 16.0) * sdf.MillimetresPerInch
const holeClearance = 1.05

const plateX = (2.0 * dX) + 20.0
const plateY = (2.0 * dY1) + 20.0
const plateZ = 4.0

func blockOffPlate() (sdf.SDF3, error) {

	// plate
	plate := sdf.Box2D(v2.Vec{plateX, plateY}, 1.0*plateZ)

	// holes
	hole, err := sdf.Circle2D(holeClearance * holeRadius)
	if err != nil {
		return nil, err
	}

	posn := []v2.Vec{
		{dX, dY0},
		{-dX, -dY0},
		{dX, -dY0},
		{-dX, dY0},
		{dX, dY1},
		{-dX, -dY1},
		{dX, -dY1},
		{-dX, dY1},
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

	s, err = airIntakeCover()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "air.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
