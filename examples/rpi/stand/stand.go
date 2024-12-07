//-----------------------------------------------------------------------------
/*

Raspberry Pi Display Stand

*/
//-----------------------------------------------------------------------------

package stand

import (
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

const displayAngle = 15.0 // degrees from vertical
var tanTheta = math.Tan(sdf.DtoR(displayAngle))
var invCosTheta = 1.0 / math.Cos(sdf.DtoR(displayAngle))

const filletRadius = 10.0

const baseHeight = 8.0
const baseWidth = 100.0
const baseLength = 160.0

const baseFootX = 30.0
const baseFootY = 15.0
const baseHoleRadius = 2.0

var baseHolePosn = v2.Vec{0.7, 0.8}

const supportPosn = 0.25 // fraction of baseWidth
const supportHeight = 120.0
const supportThickness = 5.0
const supportLength = 20.0

const webSize = 7.0
const webLength = 5.0

// 4 x M3 mounting holes on display
const displayW = 126.2
const displayH = 65.65
const displayHoleRadius = 0.5 * 3.9
const displayPosn = 0.7 // fraction of supportHeight

//-----------------------------------------------------------------------------

// sideProfile returns the 2d web/support profile
func sideProfile(t float64) (sdf.SDF2, error) {
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(baseWidth, 0).Rel()
	p.Add(-baseHeight*tanTheta, baseHeight).Rel()
	p.Add(-baseWidth*supportPosn, 0).Rel().Smooth(filletRadius, 5)
	p.Add(-supportHeight*tanTheta, supportHeight).Rel()
	p.Add(-invCosTheta*(supportThickness+t), 0).Rel()
	p.Add(tanTheta*(supportHeight-t), t-supportHeight).Rel().Smooth(filletRadius, 7)
	p.Add(0, baseHeight+t)
	p.Add(0, 0)
	return sdf.Polygon2D(p.Vertices())
}

func webs() (sdf.SDF3, error) {
	s2d, err := sideProfile(webSize)
	if err != nil {
		return nil, err
	}
	l := webLength
	s := sdf.Extrude3D(s2d, l)
	ofs := 0.5 * (baseLength - l)
	s0 := sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, ofs}))
	s1 := sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, -ofs}))
	return sdf.Union3D(s0, s1), nil
}

func supports() (sdf.SDF3, error) {
	s2d, err := sideProfile(0)
	if err != nil {
		return nil, err
	}
	l := supportLength + webLength
	s := sdf.Extrude3D(s2d, l)
	ofs := 0.5 * (baseLength - l)
	s0 := sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, ofs}))
	s1 := sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, -ofs}))
	return sdf.Union3D(s0, s1), nil
}

//-----------------------------------------------------------------------------

// baseProfile returns the 2d base profile
func baseProfile() (sdf.SDF2, error) {
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(baseWidth, 0).Rel()
	p.Add(-baseHeight*tanTheta, baseHeight).Rel()
	p.Add(0, baseHeight)
	p.Add(0, 0)
	return sdf.Polygon2D(p.Vertices())
}

func base() (sdf.SDF3, error) {
	s2d, err := baseProfile()
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s2d, baseLength), nil
}

//-----------------------------------------------------------------------------

func baseCutout() (sdf.SDF3, error) {
	holeSize := v2.Vec{baseLength - 2*baseFootX, baseWidth - 2*baseFootY}
	s2d := sdf.Box2D(holeSize, filletRadius)
	s := sdf.Extrude3D(s2d, baseHeight)
	s = sdf.Transform3D(s, sdf.RotateX(sdf.DtoR(90)))
	s = sdf.Transform3D(s, sdf.RotateY(sdf.DtoR(90)))
	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0.5 * baseWidth, 0.5 * baseHeight, 0}))
	return s, nil
}

func baseHole() (sdf.SDF3, error) {
	s, err := obj.CounterSunkHole3D(baseHeight, baseHoleRadius)
	if err != nil {
		return nil, err
	}
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, 0.5 * baseHeight})), nil
}

func baseHoles() (sdf.SDF3, error) {
	s, err := baseHole()
	if err != nil {
		return nil, err
	}

	dx := 0.5 * baseHolePosn.X * baseWidth
	dy := 0.5 * baseHolePosn.Y * baseLength

	holes := sdf.Multi3D(s, v3.VecSet{{dx, dy, 0}, {-dx, dy, 0}, {dx, -dy, 0}, {-dx, -dy, 0}})
	holes = sdf.Transform3D(holes, sdf.RotateX(sdf.DtoR(-90)))
	holes = sdf.Transform3D(holes, sdf.Translate3d(v3.Vec{0.5 * baseWidth, 0, 0}))
	return holes, nil
}

//-----------------------------------------------------------------------------

func displayHoles() (sdf.SDF3, error) {

	s, err := sdf.Cylinder3D(2*supportThickness, displayHoleRadius, 0)
	if err != nil {
		return nil, err
	}

	dx := 0.5 * displayW
	dy := 0.5 * displayH

	holes := sdf.Multi3D(s, v3.VecSet{{dx, dy, 0}, {-dx, dy, 0}, {dx, -dy, 0}, {-dx, -dy, 0}})
	holes = sdf.Transform3D(holes, sdf.RotateY(sdf.DtoR(90)))
	holes = sdf.Transform3D(holes, sdf.RotateZ(sdf.DtoR(15)))

	yOfs := displayPosn * supportHeight
	xOfs := (1-supportPosn)*baseWidth - (baseHeight * tanTheta) - (yOfs * tanTheta)
	holes = sdf.Transform3D(holes, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	return holes, nil
}

//-----------------------------------------------------------------------------

func DisplayStand() (sdf.SDF3, error) {

	base, err := base()
	if err != nil {
		return nil, err
	}

	supports, err := supports()
	if err != nil {
		return nil, err
	}

	webs, err := webs()
	if err != nil {
		return nil, err
	}

	cutout, err := baseCutout()
	if err != nil {
		return nil, err
	}

	baseHoles, err := baseHoles()
	if err != nil {
		return nil, err
	}

	displayHoles, err := displayHoles()
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(sdf.Union3D(base, webs, supports), sdf.Union3D(cutout, baseHoles, displayHoles)), nil
}

//-----------------------------------------------------------------------------
