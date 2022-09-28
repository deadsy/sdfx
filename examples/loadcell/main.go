//-----------------------------------------------------------------------------
/*

Load Cell Holder

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func holder() (sdf.SDF3, error) {

	// dimensions taken from loadcell
	const xLoadcell = 34.0
	const yLoadCell = 34.0
	const zLoadCell = 3.0
	const rLoadCell = 8.0
	const innerMargin = 4.0

	// dimensions to outside body
	const outerMargin = 4.0
	const bodyHeight = 2.0 * 8.0

	// body
	bodySize := v2.Vec{
		xLoadcell + 2.0*outerMargin,
		yLoadCell + 2.0*outerMargin,
	}
	bodyRadius := rLoadCell + outerMargin
	body2d := sdf.Box2D(bodySize, bodyRadius)
	body3d, err := sdf.ExtrudeRounded3D(body2d, bodyHeight, 2.0)
	if err != nil {
		return nil, err
	}

	// tabs
	tabX := 15.0
	tabSize := v2.Vec{
		bodySize.X + 2.0*tabX,
		0.5 * bodySize.Y,
	}
	tabHeight := bodyHeight * 0.75
	tab2d := sdf.Box2D(tabSize, bodyRadius*0.25)
	tab3d, err := sdf.ExtrudeRounded3D(tab2d, tabHeight, 2.0)
	if err != nil {
		return nil, err
	}

	// screw holes
	screw0, err := obj.CounterSunkHole3D(tabHeight, 2.0)
	if err != nil {
		return nil, err
	}
	screwOfs := 0.5*(bodySize.X+tabX) + 1.0
	screwL := sdf.Transform3D(screw0, sdf.Translate3d(v3.Vec{-screwOfs, 0, 0}))
	screwR := sdf.Transform3D(screw0, sdf.Translate3d(v3.Vec{screwOfs, 0, 0}))
	screw3d := sdf.Union3D(screwL, screwR)

	// inner hole
	holeSize := v2.Vec{
		xLoadcell - 2.0*innerMargin,
		yLoadCell - 2.0*innerMargin,
	}
	holeRadius := rLoadCell - innerMargin
	hole2d := sdf.Box2D(holeSize, holeRadius)
	hole3d := sdf.Extrude3D(hole2d, bodyHeight)

	// recess
	recessSize := v2.Vec{
		xLoadcell,
		yLoadCell,
	}
	recess2d := sdf.Box2D(recessSize, rLoadCell)
	recess3d := sdf.Extrude3D(recess2d, zLoadCell)
	zOfs := 0.5 * (bodyHeight - zLoadCell)
	recess3d = sdf.Transform3D(recess3d, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	// wire recess
	wireSize := v3.Vec{2.0, 2.0, 3.0 * outerMargin}
	wire3d, err := sdf.Box3D(wireSize, 0)
	if err != nil {
		return nil, err
	}
	wire3d = sdf.Transform3D(wire3d, sdf.RotateX(sdf.DtoR(90)))
	zOfs = 0.5 * (bodyHeight - wireSize.X)
	yOfs := 0.5 * (yLoadCell + outerMargin)
	wire3d = sdf.Transform3D(wire3d, sdf.Translate3d(v3.Vec{0, yOfs, zOfs}))

	holder := sdf.Union3D(body3d, tab3d)
	// add some filleting
	holder.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(2.0))
	// remove the holes
	holder = sdf.Difference3D(holder, sdf.Union3D(hole3d, recess3d, screw3d, wire3d))
	// cut it along the xy plane
	holder = sdf.Cut3D(holder, v3.Vec{0, 0, 0}, v3.Vec{0, 0, 1})

	return holder, nil
}

func main() {
	s, err := holder()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, "holder.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
