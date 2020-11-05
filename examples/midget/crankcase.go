//-----------------------------------------------------------------------------
/*

Crankcase Pattern and Core Box

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

const crankcaseOuterRadius = 1.0 + (5.0 / 16.0)
const crankcaseInnerRadius = 1.0 + (1.0 / 8.0)
const crankcaseOuterHeight = 7.0 / 8.0
const crankcaseInnerHeight = 5.0 / 8.0
const boltLugRadius = 0.5 * (7.0 / 16.0)

const mountLength = 4.75
const mountWidth = 4.75
const mountThickness = 0.25

//-----------------------------------------------------------------------------

// mountLugs returns the lugs used to mount the motor.
func mountLugs() sdf.SDF3 {
	const draft = 3.0
	const thickness = 0.25

	k := sdf.TruncRectPyramidParms{
		Size:        sdf.V3{4.75, thickness, crankcaseOuterHeight},
		BaseAngle:   sdf.DtoR(90 - draft),
		BaseRadius:  crankcaseOuterHeight * 0.1,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}

	s := sdf.TruncRectPyramid3D(&k)
	return sdf.Transform3D(s, sdf.Translate3d(sdf.V3{0, thickness * 0.5, 0}))
}

//-----------------------------------------------------------------------------

func cylinderMount() sdf.SDF3 {
	const draft = 3.0

	k := sdf.TruncRectPyramidParms{
		Size:        sdf.V3{2.0, 5.0 / 16.0, 1 + (3.0 / 16.0)},
		BaseAngle:   sdf.DtoR(90 - draft),
		BaseRadius:  crankcaseOuterHeight * 0.1,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}

	s := sdf.TruncRectPyramid3D(&k)
	return sdf.Transform3D(s, sdf.Translate3d(sdf.V3{0, crankcaseInnerRadius, 0}))
}

//-----------------------------------------------------------------------------

// boltLugs returns lugs that hold the crankcase halves together.
func boltLugs() sdf.SDF3 {

	const draft = 3.0

	k := sdf.TruncRectPyramidParms{
		Size:        sdf.V3{0, 0, crankcaseOuterHeight},
		BaseAngle:   sdf.DtoR(90 - draft),
		BaseRadius:  boltLugRadius,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}
	lug := sdf.TruncRectPyramid3D(&k)

	// position the lugs
	r := crankcaseOuterRadius
	d := r * math.Cos(sdf.DtoR(45))
	dy0 := 0.75
	dx0 := -math.Sqrt(r*r - dy0*dy0)
	positions := sdf.V3Set{
		{dx0, dy0, 0},
		{1.0, 13.0 / 16.0, 0},
		{-d, -d, 0},
		{d, -d, 0},
	}

	return sdf.Multi3D(lug, positions)
}

//-----------------------------------------------------------------------------

func basePattern() sdf.SDF3 {

	const draft = 3.0

	k := sdf.TruncRectPyramidParms{
		Size:        sdf.V3{0, 0, crankcaseOuterHeight},
		BaseAngle:   sdf.DtoR(90 - draft),
		BaseRadius:  crankcaseOuterRadius,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}

	body := sdf.TruncRectPyramid3D(&k)

	// add the bolt/mount lugs to the body with filleting
	s := sdf.Union3D(body, boltLugs(), mountLugs())
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(0.1))

	// cleanup the top artifacts caused by the filleting
	s = sdf.Cut3D(s, sdf.V3{0, 0, crankcaseOuterHeight}, sdf.V3{0, 0, -1})

	// add the cylinder mount
	s = sdf.Union3D(s, cylinderMount())
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(0.1))

	// cleanup the bottom artifacts caused by the filleting
	s = sdf.Cut3D(s, sdf.V3{0, 0, 0}, sdf.V3{0, 0, 1})

	return s
}

//-----------------------------------------------------------------------------

func ccRearPattern() sdf.SDF3 {
	s := basePattern()
	return s
}

func ccFrontPattern() sdf.SDF3 {
	s := basePattern()
	return s
}

//-----------------------------------------------------------------------------
