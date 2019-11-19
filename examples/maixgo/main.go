//-----------------------------------------------------------------------------
/*

MAix Go Kit

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

var baseThickness = 3.0
var pillarHeight = 14.0

//-----------------------------------------------------------------------------

func bezelStandoffs() SDF3 {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 5.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}
	positions := V3Set{
		{0, 0, zOfs},
		{82, 0, zOfs},
		{0, 54, zOfs},
		{82, 54, zOfs},
	}

	return Standoffs3D(k, positions)
}

func bezel() SDF3 {

  // bezel
	b0 := Box2D(V2{90.976, 60.0033}, 2.0)

	// cutouts
	l0 := Box2D(V2{59.902, 46.433}, 2.0)

  xOfs := 0.0
  yOfs := 0.5 * (60.0033 - 46.433)

	l0 = Transform2D(l0, Translate2d(V2{xOfs, yOfs}))

	// extrude the base
	s0 := Extrude3D(Difference2D(b0, Union2D(l0,)), baseThickness)
	xOfs = 0.0
	yOfs = 0.0
	s0 = Transform3D(s0, Translate3d(V3{xOfs, yOfs, 0}))

	// add the standoffs
	s1 := bezelStandoffs()
	s2 := Union3D(s0, s1)

	return s2
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(ScaleUniform3D(bezel(), shrink), 300, "bezel.stl")

}

//-----------------------------------------------------------------------------
