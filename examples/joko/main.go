//-----------------------------------------------------------------------------
/*

Joko Engineering Part

https://www.youtube.com/c/JokoEngineeringhelp
https://grabcad.com/library/freecad-practice-part-1

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// small end
const radiusOuterSmall = 1.0
const radiusInnerSmall = 0.55
const smallThickness = 1.0

// big end
const radiusOuterBig = 1.89
const radiusInnerBig = 2.90 * 0.5

const armWidth0 = 0.4
const armWidth1 = 0.5

const smallLength = 3.0
const overallLength = 9.75
const overallHeight = 4.0

var theta0 = sdf.DtoR(65.0) * 0.5
var theta1 = 0.5*sdf.Pi - theta0

const filletRadius0 = 0.25
const filletRadius1 = 0.50

const shaftRadius = 0.55
const keyRadius = 0.77
const keyWidth = 0.35

// derived
const centerToCenter = overallLength - radiusOuterBig - radiusOuterSmall

//-----------------------------------------------------------------------------

func planView() (sdf.SDF2, error) {
	sOuter, err := sdf.FlatFlankCam2D(centerToCenter, radiusOuterBig, radiusOuterSmall)
	if err != nil {
		return nil, err
	}
	sInner := sdf.Offset2D(sOuter, -armWidth0)
	s0 := sdf.Difference2D(sOuter, sInner)

	s1, err := sdf.Circle2D(radiusOuterSmall)
	if err != nil {
		return nil, err
	}
	s1 = sdf.Transform2D(s1, sdf.Translate2d(v2.Vec{0, centerToCenter}))

	k := obj.WasherParms{
		InnerRadius: radiusInnerBig,
		OuterRadius: radiusOuterBig,
	}
	s2, err := obj.Washer2D(&k)
	if err != nil {
		return nil, err
	}

	s3 := sdf.Union2D(s0, s1, s2)
	s3.(*sdf.UnionSDF2).SetMin(sdf.PolyMin(0.3))

	return sdf.Intersect2D(sOuter, s3), nil
}

//-----------------------------------------------------------------------------

const smoothSteps = 5

func sideView() (sdf.SDF2, error) {
	dx0 := smallThickness * 0.5
	dy1 := smallLength
	dx2 := (overallHeight - smallThickness) * 0.5
	dy2 := dx2 * math.Tan(theta1)
	dy3 := overallLength - smallLength - dy2
	dx4 := -armWidth1
	dy5 := -dy3 + (armWidth1 / math.Cos(theta1)) - armWidth1*math.Tan(theta1)
	dx6 := armWidth1 - overallHeight*0.5
	dy6 := dx6 / math.Tan(theta0)

	p := sdf.NewPolygon()
	p.Add(dx0, 0)
	p.Add(0, dy1).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(dx2, dy2).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(0, dy3).Rel()
	p.Add(dx4, 0).Rel()
	p.Add(0, dy5).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(dx6, dy6).Rel().Smooth(filletRadius0, smoothSteps)
	// mirror
	p.Add(dx6, -dy6).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(0, -dy5).Rel()
	p.Add(dx4, 0).Rel()
	p.Add(0, -dy3).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(dx2, -dy2).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(0, -dy1).Rel()
	return sdf.Polygon2D(p.Vertices())
}

//-----------------------------------------------------------------------------

func shaft() (sdf.SDF3, error) {

	k := obj.KeywayParameters{
		ShaftRadius: shaftRadius,
		KeyRadius:   keyRadius,
		KeyWidth:    keyWidth,
		ShaftLength: overallHeight,
	}

	s, err := obj.Keyway3D(&k)
	if err != nil {
		return nil, err
	}

	m := sdf.RotateY(sdf.DtoR(-90))
	m = sdf.RotateX(sdf.DtoR(-30)).Mul(m)
	m = sdf.Translate3d(v3.Vec{0, radiusOuterSmall, 0}).Mul(m)
	s = sdf.Transform3D(s, m)

	return s, nil
}

//-----------------------------------------------------------------------------

func part() (sdf.SDF3, error) {

	side2d, err := sideView()
	if err != nil {
		return nil, err
	}
	side3d := sdf.Extrude3D(side2d, radiusOuterBig*2.0)

	plan2d, err := planView()
	if err != nil {
		return nil, err
	}
	plan3d := sdf.Extrude3D(plan2d, overallHeight)
	m := sdf.RotateZ(sdf.DtoR(180))
	m = sdf.Translate3d(v3.Vec{0, centerToCenter + radiusOuterSmall, 0}).Mul(m)
	m = sdf.RotateY(sdf.DtoR(90)).Mul(m)
	plan3d = sdf.Transform3D(plan3d, m)

	part := sdf.Intersect3D(plan3d, side3d)

	shaft, err := shaft()
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(part, shaft), nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := part()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "part.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
