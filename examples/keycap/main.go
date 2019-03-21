//-----------------------------------------------------------------------------
/*

KeyCaps for Cherry MX key switches

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

const stemDiameter = 5.6
const stemCrossLength = 4.1
const stemCrossWidth = 1.35

// stem2d returns the 2D profile of a keycap stem.
func stem2d() sdf.SDF2 {
	s0 := sdf.Circle2D(stemDiameter * 0.5)
	s1 := sdf.Box2D(sdf.V2{stemCrossLength, stemCrossWidth}, 0.1*stemCrossWidth)
	s2 := sdf.Transform2D(s1, sdf.Rotate2d(sdf.DtoR(90)))
	return sdf.Difference2D(s0, sdf.Union2D(s1, s2))
}

// stem3d returns a keycap stem of a given length.
func stem3d(length float64) sdf.SDF3 {
	return sdf.Extrude3D(stem2d(), length)
}

//-----------------------------------------------------------------------------

func main() {
	sdf.RenderSTL(stem3d(5.0), 200, "stem.stl")
}

//-----------------------------------------------------------------------------
