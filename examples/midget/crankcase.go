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

//-----------------------------------------------------------------------------

func crankCaseBoltLugs() sdf.SDF3 {

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

func crankCaseFrontPattern() sdf.SDF3 {

	const draft = 3.0

	k := sdf.TruncRectPyramidParms{
		Size:        sdf.V3{0, 0, crankcaseOuterHeight},
		BaseAngle:   sdf.DtoR(90 - draft),
		BaseRadius:  crankcaseOuterRadius,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}

	body := sdf.TruncRectPyramid3D(&k)
	lugs := crankCaseBoltLugs()

	// add the lugs to the body with filleting
	s := sdf.Union3D(body, lugs)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(0.1))

	// cleanup the top/bottom artifacts caused by the filleting
	s = sdf.Cut3D(s, sdf.V3{0, 0, 0}, sdf.V3{0, 0, 1})
	s = sdf.Cut3D(s, sdf.V3{0, 0, crankcaseOuterHeight}, sdf.V3{0, 0, -1})
	return s
}

//-----------------------------------------------------------------------------
