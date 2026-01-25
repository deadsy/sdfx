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

	const displayX = 60.0
	const displayY = 45.0
	const cornerRounding = 2.0

	if negative {
		// return a panel hole for the display
		s0 := sdf.Box2D(v2.Vec{displayX, displayY}, cornerRounding)
		return sdf.Extrude3D(s0, thickness), nil
	}

	// return the display support standoffs
	const standOffZ = 4.0

	// standoffs with screw holes
	k0 := &obj.StandoffParms{
		PillarHeight:   standOffZ,
		PillarDiameter: 5.0,
		HoleDepth:      standOffZ,
		HoleDiameter:   2.4, // #4 screw
	}

	s0, err := obj.Standoff3D(k0)
	if err != nil {
		return nil, err
	}

	const xOfs = 0.5 * (86.0 - 3.0 - 6.92)
	const yOfs = 0.5 * (50.0 - 3.0 - 3.0)
	zOfs := 0.5 * (standOffZ + thickness)
	skew := xOfs - 0.5*displayX - 4.5

	positions := v3.VecSet{
		{xOfs + skew, yOfs, zOfs},
		{xOfs + skew, -yOfs, zOfs},
		{-xOfs + skew, yOfs, zOfs},
		{-xOfs + skew, -yOfs, zOfs},
	}
	return sdf.Multi3D(s0, positions), nil
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
		HoleDiameter:    2.0, // 2M screw
		Offset:          v2.Vec{0, -2.0},
		Thickness:       thickness,
		Countersunk:     true,
	}

	return obj.Display(&k, negative)
}

//-----------------------------------------------------------------------------

func bezel3() (sdf.SDF3, error) {

	const panelThickness = 3.0
	var xOfs, yOfs float64

	kPanel := obj.PanelParms{
		Size:         v2.Vec{100, 100},
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{4, 4, 4, 4},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    panelThickness,
	}
	panel, err := obj.Panel3D(&kPanel)
	if err != nil {
		return nil, err
	}

	// rotary encoder
	kRotaryEncoder := obj.KeyedHoleParms{
		Diameter:  9.6,
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

	yOfs = 20
	d1n = sdf.Transform3D(d1n, sdf.Translate3d(v3.Vec{0, yOfs, 0}))
	d1p = sdf.Transform3D(d1p, sdf.Translate3d(v3.Vec{0, yOfs, 0}))

	// group and move the inputs
	yOfs = -10
	input := sdf.Union3D(re, pb0, pb1)
	input = sdf.Transform3D(input, sdf.Translate3d(v3.Vec{0, yOfs, 0}))

	return sdf.Difference3D(sdf.Union3D(panel, d1p), sdf.Union3D(input, d1n)), nil
}

//-----------------------------------------------------------------------------

func main() {

	b3, err := bezel3()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(b3, "bezel3.stl", render.NewMarchingCubesOctree(500))
}

//-----------------------------------------------------------------------------
