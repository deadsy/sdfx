//-----------------------------------------------------------------------------
/*

Pico RX Housing

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

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

func picoRxBezel() (sdf.SDF3, error) {

	const panelThickness = 3.0
	var xOfs, yOfs float64

	kPanel := obj.PanelParms{
		Size:         v2.Vec{175, 75},
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{4, 4, 4, 4},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    panelThickness,
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
		Thickness: panelThickness,
	}
	re, err := obj.KeyedHole3D(&kRotaryEncoder)
	if err != nil {
		return nil, err
	}

	// push buttons
	pb, err := sdf.Box3D(v3.Vec{13.2, 10.8, panelThickness}, 0)
	if err != nil {
		return nil, err
	}
	xOfs = 22.0
	pb0 := sdf.Transform3D(pb, sdf.Translate3d(v3.Vec{xOfs, 0, 0}))
	pb1 := sdf.Transform3D(pb, sdf.Translate3d(v3.Vec{-xOfs, 0, 0}))

	// 128x64 display
	d1n, err := display1(panelThickness, true)
	if err != nil {
		return nil, err
	}
	d1p, err := display1(panelThickness, false)
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
	d0n, err := display0(panelThickness, true)
	if err != nil {
		return nil, err
	}
	d0p, err := display0(panelThickness, false)
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

func speaker() (sdf.SDF3, error) {

	const panelThickness = 3.0

	kPanel := obj.PanelParms{
		Size:         v2.Vec{90, 90},
		CornerRadius: 5.0,
		Thickness:    panelThickness,
	}
	panel, err := obj.Panel3D(&kPanel)
	if err != nil {
		return nil, err
	}

	const grilleRadius = 77.5 * 0.5

	kGrille := obj.CircleGrilleParms{
		HoleDiameter:      4.0,
		GrilleDiameter:    2.0 * grilleRadius,
		RadialSpacing:     0.5,
		TangentialSpacing: 0.5,
		Thickness:         panelThickness,
	}
	grille, err := obj.CircleGrille3D(&kGrille)
	if err != nil {
		return nil, err
	}

	kWall := obj.WasherParms{
		Thickness:   panelThickness,
		InnerRadius: grilleRadius,
		OuterRadius: grilleRadius + panelThickness,
	}
	wall, err := obj.Washer3D(&kWall)
	if err != nil {
		return nil, err
	}
	wall = sdf.Transform3D(wall, sdf.Translate3d(v3.Vec{0, 0, panelThickness}))

	return sdf.Difference3D(sdf.Union3D(panel, wall), grille), nil
}

//-----------------------------------------------------------------------------

func main() {

	s, err := picoRxBezel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "picorx_bezel.stl", render.NewMarchingCubesOctree(500))

	s, err = speaker()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "speaker.stl", render.NewMarchingCubesOctree(500))

}

//-----------------------------------------------------------------------------
