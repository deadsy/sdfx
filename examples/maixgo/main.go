//-----------------------------------------------------------------------------
/*

MAix Go Bezel

https://www.sipeed.com
https://wiki.sipeed.com/en/maix/board/go.html
https://www.seeedstudio.com/Sipeed-MAix-GO-Suit-for-RISC-V-AI-IoT-p-2874.html

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

var baseThickness = 3.0

//-----------------------------------------------------------------------------

func boardStandoffs() sdf.SDF3 {
	pillarHeight := 14.0
	zOfs := 0.5 * (pillarHeight + baseThickness)
	// standoffs with screw holes
	k := &obj.StandoffParms{
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
	positions := sdf.V3Set{
		{x0, y0, zOfs},
		{x0 + x, y0, zOfs},
		{x0, y0 + y, zOfs},
		{x0 + x, y0 + y, zOfs},
	}
	return sdf.Multi3D(obj.Standoff3D(k), positions)
}

//-----------------------------------------------------------------------------

func bezelStandoffs() sdf.SDF3 {
	pillarHeight := 22.0
	zOfs := 0.5 * (pillarHeight + baseThickness)
	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      11.0,
		HoleDiameter:   2.4, // #4 screw
	}
	x := 140.0
	y := 55.0
	x0 := -0.5 * x
	y0 := -0.5 * y
	positions := sdf.V3Set{
		{x0, y0, zOfs},
		{x0 + x, y0, zOfs},
		{x0, y0 + y, zOfs},
		{x0 + x, y0 + y, zOfs},
	}
	return sdf.Multi3D(obj.Standoff3D(k), positions)
}

//-----------------------------------------------------------------------------

func speakerHoles(d float64, ofs sdf.V2) sdf.SDF2 {
	holeRadius := 1.7
	s0 := sdf.Circle2D(holeRadius)
	s1 := sdf.MakeBoltCircle2D(holeRadius, d*0.3, 6)
	return sdf.Transform2D(sdf.Union2D(s0, s1), sdf.Translate2d(ofs))
}

func speakerHolder(d float64, ofs sdf.V2) (sdf.SDF3, error) {
	thickness := 3.0
	zOfs := 0.5 * (thickness + baseThickness)
	k := obj.WasherParms{
		Thickness:   thickness,
		InnerRadius: 0.5 * d,
		OuterRadius: 0.5 * (d + 4.0),
		Remove:      0.3,
	}
	s, err := obj.Washer3D(&k)
	if err != nil {
		return nil, err
	}
	s = sdf.Transform3D(s, sdf.RotateZ(sdf.Pi))
	return sdf.Transform3D(s, sdf.Translate3d(sdf.V3{ofs.X, ofs.Y, zOfs})), nil
}

//-----------------------------------------------------------------------------

func bezel() (sdf.SDF3, error) {

	speakerOfs := sdf.V2{60, 14}
	speakerDiameter := 20.3

	// bezel
	bezel := sdf.V2{150, 65}
	b0 := sdf.Box2D(bezel, 2)

	// lcd cutout
	lcd := sdf.V2{60, 46}
	l0 := sdf.Box2D(lcd, 2)

	// camera cutout
	c0 := sdf.Circle2D(7.25)
	c0 = sdf.Transform2D(c0, sdf.Translate2d(sdf.V2{42, 0}))

	// led hole cutout
	c1 := sdf.Circle2D(2)
	c1 = sdf.Transform2D(c1, sdf.Translate2d(sdf.V2{44, -20}))

	// speaker holes cutout
	c2 := speakerHoles(speakerDiameter, speakerOfs)

	// extrude the bezel
	s0 := sdf.Extrude3D(sdf.Difference2D(b0, sdf.Union2D(l0, c0, c1, c2)), baseThickness)

	// add the board standoffs
	s0 = sdf.Union3D(s0, boardStandoffs())

	// add the bezel standoffs (with foot rounding)
	s1 := sdf.Union3D(s0, bezelStandoffs())
	s1.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	// speaker holder
	s3, err := speakerHolder(speakerDiameter, speakerOfs)
	if err != nil {
		return nil, err
	}

	return sdf.Union3D(s1, s3), nil
}

//-----------------------------------------------------------------------------

func main() {
	b, err := bezel()
	if err != nil {
		log.Fatal(err)
	}
	sdf.RenderSTL(sdf.ScaleUniform3D(b, shrink), 330, "bezel.stl")
}

//-----------------------------------------------------------------------------
