//-----------------------------------------------------------------------------
/*

Axoloti Board Mounting Kit

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var front_panel_thickness = 3.0
var front_panel_length = 170.0
var front_panel_height = 50.0
var front_panel_y_offset = 15.0

var base_width = 50.0
var base_length = 170.0
var base_thickness = 3.0

var base_foot_width = 10.0
var base_foot_corner_radius = 3.0

var pcb_width = 50.0
var pcb_length = 160.0

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

	z_ofs := 0.5 * (pillar_height + base_thickness)

	// from the board mechanicals
	positions := V3Set{
		{3.5, 10.0, z_ofs},   // H1
		{3.5, 40.0, z_ofs},   // H2
		{54.0, 40.0, z_ofs},  // H3
		{156.5, 10.0, z_ofs}, // H4
		//{54.0, 10.0, z_ofs},  // H5
		{156.5, 40.0, z_ofs}, // H6
		{44.0, 10.0, z_ofs},  // H7
		{116.0, 10.0, z_ofs}, // H8
	}

	return Standoffs3D(k, positions)
}

//-----------------------------------------------------------------------------

func base() SDF3 {
	// base
	pp := &PanelParms{
		Size:         V2{base_length, base_width},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{7.0, 20.0, 7.0, 20.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	s0 := Panel2D(pp)

	// cutout
	l := base_length - (2.0 * base_foot_width)
	w := 18.0
	s1 := Box2D(V2{l, w}, base_foot_corner_radius)
	y_ofs := 0.5 * (base_width - pcb_width)
	s1 = Transform2D(s1, Translate2d(V2{0, y_ofs}))

	s2 := Extrude3D(Difference2D(s0, s1), base_thickness)
	x_ofs := 0.5 * pcb_length
	y_ofs = pcb_width - (0.5 * base_width)
	s2 = Transform3D(s2, Translate3d(V3{x_ofs, y_ofs, 0}))

	// standoffs
	s3 := standoffs()

	s4 := Union3D(s2, s3)
	s4.(*UnionSDF3).SetMin(PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------
// front panel cutouts

type PanelHole struct {
	center V2   // center of hole
	hole   SDF2 // 2d hole
}

// button positions
var pb_x float64 = 53.0
var pb0 = V2{pb_x, 0.8}
var pb1 = V2{pb_x + 5.334, 0.8}

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

	holes := []PanelHole{
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

func front_panel() SDF3 {

	// overall panel
	pp := &PanelParms{
		Size:         V2{front_panel_length, front_panel_height},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	panel := Panel2D(pp)

	x_ofs := 0.5 * pcb_length
	y_ofs := (0.5 * front_panel_height) - front_panel_y_offset
	panel = Transform2D(panel, Translate2d(V2{x_ofs, y_ofs}))

	// extrude to 3d
	fp := Extrude3D(Difference2D(panel, fp_cutouts()), front_panel_thickness)

	// Add buttons to the finger button
	b_height := 4.0
	b := Cylinder3D(b_height, 1.4, 0)
	b0 := Transform3D(b, Translate3d(pb0.ToV3(-0.5*b_height)))
	b1 := Transform3D(b, Translate3d(pb1.ToV3(-0.5*b_height)))

	return Union3D(fp, b0, b1)
}

//-----------------------------------------------------------------------------

// Create the STLs for the axoloti mount kit
func mount_kit() {
	// front panel
	s0 := front_panel()
	RenderSTL(Transform3D(s0, RotateY(DtoR(180.0))), 400, "fp.stl")

	// base
	s1 := base()
	RenderSTL(s1, 400, "base.stl")

	// both together
	//s0 = Transform3D(s0, Translate3d(V3{0, 80, 0}))
	//RenderSTL(Union3D(s0, s1), 400, "fp_and_base.stl")
}

//-----------------------------------------------------------------------------

// Create the STLs for the axoloti enclosure
func enclosure() {

	box_wall := 2.5
	box_width := pcb_length + (4.0 * box_wall) + 5.0

	bp := PanelBoxParms{
		Size:       V3{box_width, 50.0, 70.0}, // width, height, length
		Wall:       box_wall,                  // wall thickness
		Panel:      box_wall,                  // panel thickness
		Rounding:   5.0,                       // outer corner rounding
		FrontInset: 3.0,                       // inset for front panel
		BackInset:  3.0,                       // inset for pack panel
		Hole:       2.0,                       // ? screw
		SideTabs:   "TbtbT",                   // tab pattern
	}

	box := PanelBox3D(&bp)

	RenderSTL(box[0], 300, "panel.stl")
	RenderSTL(box[1], 300, "top.stl")
	RenderSTL(box[2], 300, "bottom.stl")
}

//-----------------------------------------------------------------------------

func main() {
	mount_kit()
	//enclosure()
}

//-----------------------------------------------------------------------------
