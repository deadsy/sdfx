//-----------------------------------------------------------------------------
/*

Tachometer Holding Bracket

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

const tachoRadius = 0.5 * 3.5 * sdf.MillimetresPerInch
const bracketHeight = 15.0
const bracketWidth = 10.0
const tabWidth = 20.0
const tabLength = 20.0
const slotWidth = 4.0
const screwRadius = 1.1 * 0.5 * (5.0 / 32.0) * sdf.MillimetresPerInch

//-----------------------------------------------------------------------------

func tachoBracket() (sdf.SDF3, error) {

	// outer bracket
	const outerRadius = tachoRadius + bracketWidth
	body, err := sdf.Circle2D(outerRadius)
	if err != nil {
		return nil, err
	}

	// inner hole
	hole, err := sdf.Circle2D(tachoRadius)
	if err != nil {
		return nil, err
	}

	// side tabs
	tabs := sdf.Box2D(v2.Vec{2.0 * (outerRadius + tabLength), tabWidth}, 0.07*(tabWidth+tabLength))

	// slot
	l := bracketWidth + tabLength
	slot := sdf.Box2D(v2.Vec{l, slotWidth}, 0)
	slot = sdf.Transform2D(slot, sdf.Translate2d(v2.Vec{0.5*l + tachoRadius, 0}))

	// panel hole
	panelHole, err := sdf.Circle2D(screwRadius)
	if err != nil {
		return nil, err
	}
	const xOfs = tachoRadius + bracketWidth + 0.5*tabLength
	panelHole = sdf.Transform2D(panelHole, sdf.Translate2d(v2.Vec{-xOfs, 0}))

	// outer body
	s3 := sdf.Union2D(body, tabs)
	s3.(*sdf.UnionSDF2).SetMin(sdf.PolyMin(bracketWidth))

	// remove the holes
	s4 := sdf.Difference2D(s3, sdf.Union2D(hole, slot, panelHole))
	bracket := sdf.Extrude3D(s4, bracketHeight)

	// clamp hole
	clampHole, err := sdf.Cylinder3D(1.1*tabWidth, screwRadius, 0)
	clampHole = sdf.Transform3D(clampHole, sdf.RotateX(0.5*sdf.Pi))
	clampHole = sdf.Transform3D(clampHole, sdf.Translate3d(v3.Vec{xOfs, 0, 0}))

	bracket = sdf.Difference3D(bracket, clampHole)

	return bracket, nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := tachoBracket()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "tacho.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
