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

const panelHoleDiameter = 4.0
const baseHoleDiameter = 4.0

const panelThickness = 3.0
const panelWidth = 190
const panelHeight = 75

const mountWidth = 16.0

const holeMargin = 0.5 * (mountWidth + panelThickness)

//-----------------------------------------------------------------------------

// pam8302: 2.5W Class D Audio Amplifier (https://www.adafruit.com/product/2130)
func pam8302(thickness float64) (sdf.SDF3, error) {

	const pillarHeight = 4.5

	// standoff with screw holes
	k := obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 4.5,
		HoleDepth:      pillarHeight,
		HoleDiameter:   2.0, // #2 screw
	}
	s, err := obj.Standoff3D(&k)
	if err != nil {
		return nil, err
	}

	xOfs := 0.4 * sdf.MillimetresPerInch * 0.5
	zOfs := 0.5 * (thickness + pillarHeight)
	positions := v3.VecSet{
		{xOfs, 0, zOfs},
		{-xOfs, 0, zOfs},
	}
	return sdf.Multi3D(s, positions), nil
}

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
		SupportDiameter: 4.5,
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

// pcbMount0 mounts the adafruit half breadboard with the rpi-pico.
func pcbMount0() (sdf.SDF3, error) {

	const width = 60.0
	const length = 90.0
	const margin = 5.0
	const height = 10.0
	const thickness = 3.0

	pp := obj.PanelParms{
		Size:         v2.Vec{width, length},
		CornerRadius: margin,
		HoleDiameter: 3,
		HoleMargin:   [4]float64{margin, margin, margin, margin},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    thickness,
		Ridge:        v2.Vec{width - 3.0*margin, length - 3.0*margin},
	}
	panel, err := obj.Panel3D(&pp)
	if err != nil {
		return nil, err
	}

	// standoff with screw hole
	sp := obj.StandoffParms{
		PillarHeight:   height,
		PillarDiameter: 8,
		HoleDepth:      height,
		HoleDiameter:   2.4, // #4 screw
	}
	standoff, err := obj.Standoff3D(&sp)
	if err != nil {
		return nil, err
	}

	// two standoffs, 2.9" apart
	zOfs := 0.5 * (height + thickness)
	yOfs := 0.5 * 2.9 * sdf.MillimetresPerInch
	positions := v3.VecSet{
		{0, -yOfs, zOfs},
		{0, yOfs, zOfs},
	}
	standoffs := sdf.Multi3D(standoff, positions)

	return sdf.Union3D(panel, standoffs), nil
}

// pcbMount1 mounts the sdr front end.
func pcbMount1() (sdf.SDF3, error) {

	return nil, nil
}

//-----------------------------------------------------------------------------

func picoRxBezel(thickness float64) (sdf.SDF3, error) {

	var xOfs, yOfs float64

	const ridgeWidth = panelWidth - (2.0 * mountWidth) - 1.0

	kPanel := obj.PanelParms{
		Size:         v2.Vec{panelWidth, panelHeight},
		CornerRadius: 5.0,
		HoleDiameter: panelHoleDiameter,
		HoleMargin:   [4]float64{holeMargin, holeMargin, holeMargin, holeMargin},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    thickness,
		Ridge:        v2.Vec{ridgeWidth, 0},
	}
	panel, err := obj.Panel3D(&kPanel)
	if err != nil {
		return nil, err
	}

	// rotary encoder
	kRotaryEncoder := obj.KeyedHoleParms{
		Diameter:  9.2,
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

func twoHoles(thickness, diameter, distance float64) (sdf.SDF3, error) {
	h, err := sdf.Cylinder3D(thickness, 0.5*diameter, 0)
	if err != nil {
		return nil, err
	}
	xOfs := 0.5 * distance
	h0 := sdf.Transform3D(h, sdf.Translate3d(v3.Vec{xOfs, 0, 0}))
	h1 := sdf.Transform3D(h, sdf.Translate3d(v3.Vec{-xOfs, 0, 0}))
	return sdf.Union3D(h0, h1), nil
}

func sideMount(thickness float64, lhs bool) (sdf.SDF3, error) {

	const mHeight = 95.0
	const mLength = 125.0
	const mSlope = 75.0
	const mRound0 = 2.0 // internal rounding

	mRound2 := mRound0 + thickness // external rounding

	// build an internal box that we will offset to build the inside and outside envelope.
	d := mSlope / math.Sqrt(2)
	bh := mHeight - 2.0*mRound2
	bl := mLength - 2.0*mRound2
	bw := 2.0 * (mountWidth - mRound2)

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
	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, mountWidth - 0.5*thickness}))

	// base holes
	baseHoles, err := twoHoles(thickness, baseHoleDiameter, 0.7*mLength)
	if err != nil {
		return nil, err
	}
	baseHoles = sdf.Transform3D(baseHoles, sdf.RotateX(sdf.DtoR(-90)))
	xOfs := 0.5 * (mLength - mHeight)
	yOfs := 0.5 * (mHeight - thickness)
	zOfs := 0.5 * mountWidth
	baseHoles = sdf.Transform3D(baseHoles, sdf.Translate3d(v3.Vec{xOfs, -yOfs, zOfs}))

	// panel holes
	const panelHoleDistance = panelHeight - 2.0*holeMargin
	panelHoles, err := twoHoles(thickness, panelHoleDiameter, panelHoleDistance)
	if err != nil {
		return nil, err
	}
	panelHoles = sdf.Transform3D(panelHoles, sdf.RotateX(sdf.DtoR(-90)))
	panelHoles = sdf.Transform3D(panelHoles, sdf.RotateZ(sdf.DtoR(-45)))
	delta := (mRound0 + 0.5*thickness) / math.Sqrt(2)
	xOfs = bl - 0.5*(bh+d) + delta
	yOfs = 0.5*(bh-d) + delta
	zOfs = 0.5 * mountWidth
	panelHoles = sdf.Transform3D(panelHoles, sdf.Translate3d(v3.Vec{xOfs, yOfs, zOfs}))

	s = sdf.Difference3D(s, sdf.Union3D(baseHoles, panelHoles))

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

	amp, err := pam8302(thickness)
	if err != nil {
		return nil, err
	}
	amp = sdf.Transform3D(amp, sdf.Translate3d(v3.Vec{55, -10, 0}))

	return sdf.Difference3D(sdf.Union3D(rhs, sp, amp), sn), nil
}

//-----------------------------------------------------------------------------

func lhsMount(thickness float64) (sdf.SDF3, error) {
	return sideMount(thickness, true)
}

//-----------------------------------------------------------------------------

func main() {

	s, err := pcbMount0()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "pcb_mount0.stl", render.NewMarchingCubesOctree(500))

	s, err = picoRxBezel(panelThickness)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "bezel.stl", render.NewMarchingCubesOctree(500))

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
