//-----------------------------------------------------------------------------
/*

Raspberry Pi Display Stand

*/
//-----------------------------------------------------------------------------

package stand

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

const displayAngle = 15.0 // degrees
var tanTheta = math.Tan(sdf.DtoR(displayAngle))
var invCosTheta = 1.0 / math.Cos(sdf.DtoR(displayAngle))

const filletRadius = 10.0

const baseHeight = 8.0
const baseWidth = 100.0
const baseLength = 160.0

const baseFootX = 30.0
const baseFootY = 15.0

const supportPosn = 0.25 // fraction of baseWidth
const supportHeight = 120.0
const supportThickness = 5.0
const supportLength = 20.0

const webSize = 7.0
const webLength = 5.0

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

func baseHole() (sdf.SDF3, error) {
	holeSize := v2.Vec{baseLength - 2*baseFootX, baseWidth - 2*baseFootY}
	s2d := sdf.Box2D(holeSize, filletRadius)
	s := sdf.Extrude3D(s2d, baseHeight)
	s = sdf.Transform3D(s, sdf.RotateX(sdf.DtoR(90)))
	s = sdf.Transform3D(s, sdf.RotateY(sdf.DtoR(90)))
	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0.5 * baseWidth, 0.5 * baseHeight, 0}))
	return s, nil
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

	hole, err := baseHole()
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(sdf.Union3D(base, webs, supports), sdf.Union3D(hole)), nil
}

//-----------------------------------------------------------------------------
