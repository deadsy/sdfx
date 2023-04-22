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
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// material shrinkage

const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const wallThickness = 3.0
const round = 0.5 * wallThickness
const clearance = 0.3

//-----------------------------------------------------------------------------

func box1(upper bool) (sdf.SDF3, error) {

	oSize := v3.Vec{40, 40, 20}
	iSize := oSize.SubScalar(2.0 * wallThickness)

	// build the box
	outer := sdf.Extrude3D(sdf.Box2D(v2.Vec{oSize.X, oSize.Y}, round), oSize.Z)
	inner := sdf.Extrude3D(sdf.Box2D(v2.Vec{iSize.X, iSize.Y}, round), iSize.Z)
	box := sdf.Difference3D(outer, inner)

	// add some internals wall
	yOfs := oSize.Y * 0.2
	wall, _ := sdf.Box3D(v3.Vec{oSize.X, wallThickness, oSize.Z}, 0)
	wall0 := sdf.Transform3D(wall, sdf.Translate3d(v3.Vec{0, yOfs, 0}))
	wall1 := sdf.Transform3D(wall, sdf.Translate3d(v3.Vec{0, -yOfs, 0}))
	box = sdf.Union3D(box, wall0, wall1)

	lidHeight := 0.5*oSize.Z - wallThickness

	if upper == true {
		box = sdf.Cut3D(box, v3.Vec{0, 0, lidHeight}, v3.Vec{0, 0, 1})
	} else {
		box = sdf.Cut3D(box, v3.Vec{0, 0, lidHeight}, v3.Vec{0, 0, -1})
	}

	// angled tabs
	tabSize := v3.Vec{2.5 * wallThickness, wallThickness, wallThickness}
	tab, err := obj.NewAngleTab(tabSize, clearance)
	if err != nil {
		return nil, err
	}
	xOfs := oSize.X * 0.25
	mSet := []sdf.M44{
		sdf.Translate3d(v3.Vec{xOfs, yOfs, lidHeight}),
		sdf.Translate3d(v3.Vec{xOfs, -yOfs, lidHeight}),
		sdf.Translate3d(v3.Vec{-xOfs, yOfs, lidHeight}),
		sdf.Translate3d(v3.Vec{-xOfs, -yOfs, lidHeight}),
	}
	box = obj.AddTabs(box, tab, upper, mSet)

	// screw tabs

	l := oSize.Z * 0.35
	k := obj.ScrewTab{
		Length:     l,                   // length of pillar
		Radius:     0.8 * wallThickness, // radius of pillar
		Round:      true,                // round the bottom of the pillar
		HoleUpper:  wallThickness,       // length of upper hole
		HoleLower:  0.8 * l,             // length of lower hole
		HoleRadius: 1,                   // radius of hole
	}
	tab, err = obj.NewScrewTab(&k)
	if err != nil {
		return nil, err
	}
	xOfs = 0.5*oSize.X - wallThickness
	yOfs = 0.5*oSize.Y - wallThickness
	mSet = []sdf.M44{
		sdf.Translate3d(v3.Vec{xOfs, yOfs, lidHeight}),
		sdf.Translate3d(v3.Vec{-xOfs, yOfs, lidHeight}),
		sdf.Translate3d(v3.Vec{xOfs, -yOfs, lidHeight}),
		sdf.Translate3d(v3.Vec{-xOfs, -yOfs, lidHeight}),
	}
	box = obj.AddTabs(box, tab, upper, mSet)

	return box, nil
}

//-----------------------------------------------------------------------------

func box0(upper bool) (sdf.SDF3, error) {

	oSize := v3.Vec{40, 40, 20}
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

	if upper == true {
		box = sdf.Cut3D(box, v3.Vec{0, 0, lidHeight}, v3.Vec{0, 0, 1})
	} else {
		box = sdf.Cut3D(box, v3.Vec{0, 0, lidHeight}, v3.Vec{0, 0, -1})
	}

	tabSize := v3.Vec{3.0 * wallThickness, 0.5 * wallThickness, wallThickness}
	tab, err := obj.NewStraightTab(tabSize, clearance)
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

	return obj.AddTabs(box, tab, upper, mSet), nil
}

//-----------------------------------------------------------------------------

func main() {

	s, err := box0(true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "box0_upper.stl", render.NewMarchingCubesOctree(300))

	s, err = box0(false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "box0_lower.stl", render.NewMarchingCubesOctree(300))

	s, err = box1(true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "box1_upper.stl", render.NewMarchingCubesOctree(300))

	s, err = box1(false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "box1_lower.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
