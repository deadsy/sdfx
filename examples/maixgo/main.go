//-----------------------------------------------------------------------------
/*

MAix Go Bezel

https://www.sipeed.com
https://wiki.sipeed.com/en/maix/board/go.html
https://www.seeedstudio.com/Sipeed-MAix-GO-Suit-for-RISC-V-AI-IoT-p-2874.html

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

//-----------------------------------------------------------------------------

func boardStandoffs() SDF3 {
	pillarHeight := 14.0
	zOfs := 0.5 * (pillarHeight + baseThickness)
	// standoffs with screw holes
	k := &StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 4.5,
		HoleDepth:      11.0,
		HoleDiameter:   2.6, // #4 screw
		NumberWebs:     2,
		WebHeight:      10,
		WebDiameter:    12,
		WebWidth:       3.5,
	}
	x := 82.0
	y := 54.0
	x0 := -34.0
	y0 := -0.5 * y
	positions := V3Set{
		{x0, y0, zOfs},
		{x0 + x, y0, zOfs},
		{x0, y0 + y, zOfs},
		{x0 + x, y0 + y, zOfs},
	}
	return Standoffs3D(k, positions)
}

//-----------------------------------------------------------------------------

func bezelStandoffs() SDF3 {
	pillarHeight := 22.0
	zOfs := 0.5 * (pillarHeight + baseThickness)
	// standoffs with screw holes
	k := &StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      11.0,
		HoleDiameter:   2.4, // #4 screw
	}
	x := 140.0
	y := 55.0
	x0 := -0.5 * x
	y0 := -0.5 * y
	positions := V3Set{
		{x0, y0, zOfs},
		{x0 + x, y0, zOfs},
		{x0, y0 + y, zOfs},
		{x0 + x, y0 + y, zOfs},
	}
	return Standoffs3D(k, positions)
}

//-----------------------------------------------------------------------------

func speakerHoles(d float64, ofs V2) SDF2 {
	holeRadius := 1.7
	s0 := Circle2D(holeRadius)
	s1 := MakeBoltCircle2D(holeRadius, d*0.3, 6)
	return Transform2D(Union2D(s0, s1), Translate2d(ofs))
}

func speakerHolder(d float64, ofs V2) SDF3 {
	thickness := 3.0
	zOfs := 0.5 * (thickness + baseThickness)
	k := WasherParms{
		Thickness:   thickness,
		InnerRadius: 0.5 * d,
		OuterRadius: 0.5 * (d + 4.0),
		Remove:      0.3,
	}
	s := Washer3D(&k)
	s = Transform3D(s, RotateZ(Pi))
	return Transform3D(s, Translate3d(V3{ofs.X, ofs.Y, zOfs}))
}

//-----------------------------------------------------------------------------

func bezel() SDF3 {

	speakerOfs := V2{60, 14}
	speakerDiameter := 20.3

	// bezel
	bezel := V2{150, 65}
	b0 := Box2D(bezel, 2)

	// lcd cutout
	lcd := V2{60, 46}
	l0 := Box2D(lcd, 2)

	// camera cutout
	c0 := Circle2D(7.25)
	c0 = Transform2D(c0, Translate2d(V2{42, 0}))

	// led hole cutout
	c1 := Circle2D(2)
	c1 = Transform2D(c1, Translate2d(V2{44, -20}))

	// speaker holes cutout
	c2 := speakerHoles(speakerDiameter, speakerOfs)

	// extrude the bezel
	s0 := Extrude3D(Difference2D(b0, Union2D(l0, c0, c1, c2)), baseThickness)

	// add the board standoffs
	s0 = Union3D(s0, boardStandoffs())

	// add the bezel standoffs (with foot rounding)
	s1 := Union3D(s0, bezelStandoffs())
	s1.(*UnionSDF3).SetMin(PolyMin(3.0))

	// speaker holder
	s3 := speakerHolder(speakerDiameter, speakerOfs)

	return Union3D(s1, s3)
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(ScaleUniform3D(bezel(), shrink), 330, "bezel.stl")
}

//-----------------------------------------------------------------------------
