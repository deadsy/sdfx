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
		PillarDiameter: 4.5,
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
	bezel := V2{90.976, 60.0033}
	b0 := Box2D(bezel, 2.0)
	b0 = Transform2D(b0, Translate2d(bezel.MulScalar(0.5)))

	// lcd cutout
	lcd := V2{59.902, 46.433}
	l0 := Box2D(lcd, 2.0)
	l0 = Transform2D(l0, Translate2d(lcd.MulScalar(0.5)))
	l0 = Transform2D(l0, Translate2d(V2{9.1289, 6.8267}))

	// camera cutout
	c0 := Circle2D(7.0)
	c0 = Transform2D(c0, Translate2d(V2{81.2903, 29.8240}))

	// led hole cutout
	c1 := Circle2D(1.9221)
	c1 = Transform2D(c1, Translate2d(V2{83.1539, 9.6240}))

	// extrude the base
	s0 := Extrude3D(Difference2D(b0, Union2D(l0, c0, c1)), baseThickness)

	// standoffs
	s1 := bezelStandoffs()
	s1 = Transform3D(s1, Translate3d(V3{5.1, 3.0277, 0}))

	s2 := Union3D(s0, s1)

	return s2
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(ScaleUniform3D(bezel(), shrink), 300, "bezel.stl")
}

//-----------------------------------------------------------------------------
