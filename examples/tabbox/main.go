//-----------------------------------------------------------------------------
/*

Demonstrate tabs connecting a box and lid.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// material shrinkage

const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const wallThickness = 3.0

func tab() (obj.Tab, error) {
	tabSize := v3.Vec{3.0 * wallThickness, 0.5 * wallThickness, wallThickness}
	const tabClearance = 0.1
	return obj.NewStraightTab(tabSize, tabClearance)
}

func tabbox(upper bool) (sdf.SDF3, error) {

	round := 0.5 * wallThickness
	oSize := v3.Vec{40, 40, 30}
	iSize := oSize.SubScalar(2.0 * wallThickness)

	outer, err := sdf.Box3D(oSize, round)
	if err != nil {
		return nil, err
	}
	inner, err := sdf.Box3D(iSize, round)
	if err != nil {
		return nil, err
	}

	box := sdf.Difference3D(outer, inner)
	lidHeight := oSize.Z * 0.25

	tab, err := tab()
	if err != nil {
		return nil, err
	}

	xOfs := 0.5 * (iSize.X + wallThickness)
	yOfs := 0.5 * (iSize.Y + wallThickness)

	mSet := []sdf.M44{
		sdf.Translate3d(v3.Vec{xOfs, 0, lidHeight}).Mul(sdf.RotateZ(sdf.DtoR(90))),
		sdf.Translate3d(v3.Vec{-xOfs, 0, lidHeight}).Mul(sdf.RotateZ(sdf.DtoR(90))),
		sdf.Translate3d(v3.Vec{0, yOfs, lidHeight}),
		sdf.Translate3d(v3.Vec{0, -yOfs, lidHeight}),
	}

	var s sdf.SDF3

	if upper == true {
		s = sdf.Cut3D(box, v3.Vec{0, 0, lidHeight}, v3.Vec{0, 0, 1})
	} else {
		s = sdf.Cut3D(box, v3.Vec{0, 0, lidHeight}, v3.Vec{0, 0, -1})
	}

	return obj.AddTabs(s, tab, upper, mSet), nil
}

//-----------------------------------------------------------------------------

func main() {

	upper, err := tabbox(true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(upper, shrink), "upper.stl", render.NewMarchingCubesOctree(300))

	lower, err := tabbox(false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(lower, shrink), "lower.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
