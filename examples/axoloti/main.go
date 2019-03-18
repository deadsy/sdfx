//-----------------------------------------------------------------------------
/*

Axoloti Board Mounting Kit

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var frontPanelThickness = 3.0
var frontPanelLength = 170.0
var frontPanelHeight = 50.0
var frontPanelYOffset = 15.0

var baseWidth = 50.0
var baseLength = 170.0
var baseThickness = 3.0

var baseFootWidth = 10.0
var baseFootCornerRadius = 3.0

var pcbWidth = 50.0
var pcbLength = 160.0

var pillar_height = 16.8

//-----------------------------------------------------------------------------

// multiple standoffs
func standoffs() SDF3 {

	k := &StandoffParms{
		PillarHeight:   pillar_height,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}

	zOfs := 0.5 * (pillar_height + baseThickness)

	// from the board mechanicals
	positions := V3Set{
		{3.5, 10.0, zOfs},   // H1
		{3.5, 40.0, zOfs},   // H2
		{54.0, 40.0, zOfs},  // H3
		{156.5, 10.0, zOfs}, // H4
		//{54.0, 10.0, zOfs},  // H5
		{156.5, 40.0, zOfs}, // H6
		{44.0, 10.0, zOfs},  // H7
		{116.0, 10.0, zOfs}, // H8
	}

	return Standoffs3D(k, positions)
}

//-----------------------------------------------------------------------------

// base returns the base mount.
func base() SDF3 {
	// base
	pp := &PanelParms{
		Size:         V2{baseLength, baseWidth},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{7.0, 20.0, 7.0, 20.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	s0 := Panel2D(pp)

	// cutout
	l := baseLength - (2.0 * baseFootWidth)
	w := 18.0
	s1 := Box2D(V2{l, w}, baseFootCornerRadius)
	yOfs := 0.5 * (baseWidth - pcbWidth)
	s1 = Transform2D(s1, Translate2d(V2{0, yOfs}))

	s2 := Extrude3D(Difference2D(s0, s1), baseThickness)
	xOfs := 0.5 * pcbLength
	yOfs = pcbWidth - (0.5 * baseWidth)
	s2 = Transform3D(s2, Translate3d(V3{xOfs, yOfs, 0}))

	// standoffs
	s3 := standoffs()

	s4 := Union3D(s2, s3)
	s4.(*UnionSDF3).SetMin(PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------
// front panel cutouts

type panelHole struct {
	center V2   // center of hole
	hole   SDF2 // 2d hole
}

// button positions
var pbX = 53.0
var pb0 = V2{pbX, 0.8}
var pb1 = V2{pbX + 5.334, 0.8}

// fp_cutouts returns the 2D front panel cutouts
func fp_cutouts() SDF2 {

	s_midi := Circle2D(0.5 * 17.0)
	s_jack := Circle2D(0.5 * 11.5)
	s_led := Box2D(V2{1.6, 1.6}, 0)

	fb := &FingerButtonParms{
		Width:  4.0,
		Gap:    0.6,
		Length: 20.0,
	}
	s_button := Transform2D(FingerButton2D(fb), Rotate2d(DtoR(-90)))

	jack_x := 123.0
	midi_x := 18.8
	led_x := 62.9

	holes := []panelHole{
		{V2{midi_x, 10.2}, s_midi},               // MIDI DIN Jack
		{V2{midi_x + 20.32, 10.2}, s_midi},       // MIDI DIN Jack
		{V2{jack_x, 8.14}, s_jack},               // 1/4" Stereo Jack
		{V2{jack_x + 19.5, 8.14}, s_jack},        // 1/4" Stereo Jack
		{V2{107.6, 2.3}, Circle2D(0.5 * 5.5)},    // 3.5 mm Headphone Jack
		{V2{led_x, 0.5}, s_led},                  // LED
		{V2{led_x + 3.635, 0.5}, s_led},          // LED
		{pb0, s_button},                          // Push Button
		{pb1, s_button},                          // Push Button
		{V2{84.1, 1.0}, Box2D(V2{16.0, 7.5}, 0)}, // micro SD card
		{V2{96.7, 1.0}, Box2D(V2{11.0, 7.5}, 0)}, // micro USB connector
		{V2{73.1, 7.1}, Box2D(V2{7.5, 15.0}, 0)}, // fullsize USB connector
	}

	s := make([]SDF2, len(holes))
	for i, k := range holes {
		s[i] = Transform2D(k.hole, Translate2d(k.center))
	}

	return Union2D(s...)
}

//-----------------------------------------------------------------------------

// frontPanel returns the front panel mount.
func frontPanel() SDF3 {

	// overall panel
	pp := &PanelParms{
		Size:         V2{frontPanelLength, frontPanelHeight},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	panel := Panel2D(pp)

	xOfs := 0.5 * pcbLength
	yOfs := (0.5 * frontPanelHeight) - frontPanelYOffset
	panel = Transform2D(panel, Translate2d(V2{xOfs, yOfs}))

	// extrude to 3d
	fp := Extrude3D(Difference2D(panel, fp_cutouts()), frontPanelThickness)

	// Add buttons to the finger button
	b_height := 4.0
	b := Cylinder3D(b_height, 1.4, 0)
	b0 := Transform3D(b, Translate3d(pb0.ToV3(-0.5*b_height)))
	b1 := Transform3D(b, Translate3d(pb1.ToV3(-0.5*b_height)))

	return Union3D(fp, b0, b1)
}

//-----------------------------------------------------------------------------

// mountingKit creates the STLs for the axoloti mount kit
func mountingKit() {

	s0 := frontPanel()
	RenderSTL(Transform3D(s0, RotateY(DtoR(180.0))), 400, "panel.stl")

	s1 := base()
	RenderSTL(s1, 400, "base.stl")

	// both together
	//s0 = Transform3D(s0, Translate3d(V3{0, 80, 0}))
	//RenderSTL(Union3D(s0, s1), 400, "panel_and_base.stl")
}

//-----------------------------------------------------------------------------

func main() {
	mountingKit()
}

//-----------------------------------------------------------------------------
