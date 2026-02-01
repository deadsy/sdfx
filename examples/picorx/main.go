//-----------------------------------------------------------------------------
/*

Pico RX Housing

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

// display0 : 320x240 TJCTM24028-SPI
func display0(thickness float64, negative bool) (sdf.SDF3, error) {

	k := obj.DisplayParms{
		Window:          v2.Vec{60, 45},
		Rounding:        2.0,
		Supports:        v2.Vec{76.08, 44.0},
		SupportHeight:   4.0,
		SupportDiameter: 5.0,
		HoleDiameter:    3.0, // 2.5M screw
		Offset:          v2.Vec{2.5, 0},
		Thickness:       thickness,
		Countersunk:     true,
	}

	return obj.Display(&k, negative)
}

//-----------------------------------------------------------------------------

// display1 : GME12864-11 (128x64 SSD1306)
func display1(thickness float64, negative bool) (sdf.SDF3, error) {

	k := obj.DisplayParms{
		Window:          v2.Vec{26, 14},
		Rounding:        1,
		Supports:        v2.Vec{23.5, 23.8},
		SupportHeight:   2.1,
		SupportDiameter: 4.0,
		HoleDiameter:    2.5, // 2M screw
		Offset:          v2.Vec{0, -2.0},
		Thickness:       thickness,
		Countersunk:     true,
	}

	return obj.Display(&k, negative)
}

//-----------------------------------------------------------------------------

func speakerGrille(thickness float64, negative bool) (sdf.SDF3, error) {

	const grilleRadius = 77.5 * 0.5

	if negative {
		// grille holes
		kGrille := obj.CircleGrilleParms{
			HoleDiameter:      4.0,
			GrilleDiameter:    2.0 * grilleRadius,
			RadialSpacing:     0.5,
			TangentialSpacing: 0.5,
			Thickness:         thickness,
		}
		return obj.CircleGrille3D(&kGrille)
	}

	// speaker wall
	kWall := obj.WasherParms{
		Thickness:   thickness,
		InnerRadius: grilleRadius,
		OuterRadius: grilleRadius + thickness,
	}
	wall, err := obj.Washer3D(&kWall)
	if err != nil {
		return nil, err
	}
	return sdf.Transform3D(wall, sdf.Translate3d(v3.Vec{0, 0, thickness})), nil
}

//-----------------------------------------------------------------------------

func picoRxBezel(thickness float64) (sdf.SDF3, error) {

	var xOfs, yOfs float64

	kPanel := obj.PanelParms{
		Size:         v2.Vec{175, 75},
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{4, 4, 4, 4},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    thickness,
		Ridge:        true,
	}
	panel, err := obj.Panel3D(&kPanel)
	if err != nil {
		return nil, err
	}

	// rotary encoder
	kRotaryEncoder := obj.KeyedHoleParms{
		Diameter:  9.4, // 9.6 == loose
		KeySize:   0.9,
		NumKeys:   2,
		Thickness: thickness,
	}
	re, err := obj.KeyedHole3D(&kRotaryEncoder)
	if err != nil {
		return nil, err
	}

	// push buttons
	pb, err := sdf.Box3D(v3.Vec{13.2, 10.8, thickness}, 0)
	if err != nil {
		return nil, err
	}
	xOfs = 22.0
	pb0 := sdf.Transform3D(pb, sdf.Translate3d(v3.Vec{xOfs, 0, 0}))
	pb1 := sdf.Transform3D(pb, sdf.Translate3d(v3.Vec{-xOfs, 0, 0}))

	// 128x64 display
	d1n, err := display1(thickness, true)
	if err != nil {
		return nil, err
	}
	d1p, err := display1(thickness, false)
	if err != nil {
		return nil, err
	}
	yOfs = 13.0
	xOfs = 47.0
	d1n = sdf.Transform3D(d1n, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))
	d1p = sdf.Transform3D(d1p, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// group and move the inputs
	yOfs = -17.0
	input := sdf.Union3D(re, pb0, pb1)
	input = sdf.Transform3D(input, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// 320x240 display
	d0n, err := display0(thickness, true)
	if err != nil {
		return nil, err
	}
	d0p, err := display0(thickness, false)
	if err != nil {
		return nil, err
	}
	yOfs = 0.0
	xOfs = -35.0
	d0n = sdf.Transform3D(d0n, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))
	d0p = sdf.Transform3D(d0p, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	return sdf.Difference3D(sdf.Union3D(panel, d0p, d1p), sdf.Union3D(input, d0n, d1n)), nil
}

//-----------------------------------------------------------------------------

func sideMount(thickness float64, lhs bool) (sdf.SDF3, error) {

	const mWidth = 15.0
	const mHeight = 95.0
	const mLength = 125.0
	const mSlope = 75.0
	const mRound0 = 2.0 // internal rounding

	mRound2 := mRound0 + thickness // external rounding

	// build an internal box that we will offset to build the inside and outside envelope.
	d := mSlope / math.Sqrt(2)
	bh := mHeight - 2.0*mRound2
	bl := mLength - 2.0*mRound2
	bw := 2.0 * (mWidth - mRound2)

	p := sdf.NewPolygon()
	p.Add(-bh*0.5, bh*0.5)
	p.Add(0, -bh).Rel()
	p.Add(bl, 0).Rel()
	p.Add(0, bh-d).Rel()
	p.Add(-d, d).Rel()
	b0, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	box := sdf.Extrude3D(b0, bw)

	inner := sdf.Offset3D(box, mRound0)
	outer := sdf.Offset3D(box, mRound2)
	s := sdf.Difference3D(outer, inner)

	s = sdf.Cut3D(s, v3.Vec{0, 0, 0}, v3.Vec{0, 0, -1})
	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, mWidth - 0.5*thickness}))

	if lhs {
		s = sdf.Transform3D(s, sdf.MirrorXZ())
	}

	return s, nil
}

//-----------------------------------------------------------------------------

func rhsMount(thickness float64) (sdf.SDF3, error) {

	rhs, err := sideMount(thickness, false)
	if err != nil {
		return nil, err
	}

	sn, err := speakerGrille(thickness, true)
	if err != nil {
		return nil, err
	}

	sp, err := speakerGrille(thickness, false)
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(sdf.Union3D(rhs, sp), sn), nil
}

//-----------------------------------------------------------------------------

func lhsMount(thickness float64) (sdf.SDF3, error) {

	lhs, err := sideMount(thickness, true)
	if err != nil {
		return nil, err
	}

	return lhs, nil
}

//-----------------------------------------------------------------------------

func main() {

	const panelThickness = 3.0

	s, err := picoRxBezel(panelThickness)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "picorx_bezel.stl", render.NewMarchingCubesOctree(500))

	s, err = rhsMount(panelThickness)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "rhs.stl", render.NewMarchingCubesOctree(500))

	s, err = lhsMount(panelThickness)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "lhs.stl", render.NewMarchingCubesOctree(500))

}

//-----------------------------------------------------------------------------
