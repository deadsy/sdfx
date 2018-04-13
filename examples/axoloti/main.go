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
var front_panel_radius = 5.0

var base_width = 50.0
var base_length = 170.0
var base_thickness = 3.0
var base_corner_radius = 5.0

var base_foot_width = 10.0
var base_foot_corner_radius = 3.0
var base_hole_radius = 1.0

var pcb_thickness = 1.4
var pcb_width = 50.0
var pcb_length = 160.0

var pillar_height = 10.0

//-----------------------------------------------------------------------------

// one standoff
func standoff() SDF3 {
	standoff_parms := &Standoff_Parms{
		Pillar_height: pillar_height,
		Pillar_radius: 6.0 / 2.0,
		Hole_depth:    7.0,
		Hole_radius:   1.5 / 2.0,
		Number_webs:   0,
		Web_height:    0,
		Web_radius:    0,
		Web_width:     0,
	}
	return Standoff3D(standoff_parms)
}

// multiple standoffs
func standoffs() SDF3 {
	// from the board mechanicals
	positions := V2Set{
		{3.5, 10.0},   // H1
		{3.5, 40.0},   // H2
		{54.0, 40.0},  // H3
		{156.5, 10.0}, // H4
		{54.0, 10.0},  // H5
		{156.5, 40.0}, // H6
		{44.0, 10.0},  // H7
		{116.0, 10.0}, // H8
	}
	s := make([]SDF3, len(positions))
	for i, p := range positions {
		s[i] = Transform3D(standoff(), Translate3d(V3{p.X, p.Y, 0}))
	}
	return Union3D(s...)
}

//-----------------------------------------------------------------------------

func base_holes() SDF2 {
	// from the board mechanicals
	positions := V2Set{
		{-60.0, 15.0},
		{-60.0, -15.0},
		{0, 15.0},
		{0, -15.0},
		{60.0, 15.0},
		{60.0, -15.0},
	}
	s := make([]SDF2, len(positions))
	for i, p := range positions {
		s[i] = Transform2D(Circle2D(base_hole_radius), Translate2d(V2{p.X, p.Y}))
	}
	return Union2D(s...)
}

//-----------------------------------------------------------------------------

func base() SDF3 {
	// base
	s0 := Box2D(V2{base_length, base_width}, base_corner_radius)

	// cutout
	l := base_length - (2.0 * base_foot_width)
	w := 18.0
	s1 := Box2D(V2{l, w}, base_foot_corner_radius)
	s1 = Union2D(base_holes(), s1)
	y_ofs := 0.5 * (base_width - pcb_width)
	s1 = Transform2D(s1, Translate2d(V2{0, y_ofs}))

	s2 := Extrude3D(Difference2D(s0, s1), base_thickness)
	x_ofs := 0.5 * pcb_length
	y_ofs = pcb_width - (0.5 * base_width)
	s2 = Transform3D(s2, Translate3d(V3{x_ofs, y_ofs, 0}))

	// standoffs
	z_ofs := 0.5 * (pillar_height + base_thickness)
	s3 := Transform3D(standoffs(), Translate3d(V3{0, 0, z_ofs}))

	s4 := Union3D(s2, s3)
	s4.(*UnionSDF3).SetMin(PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------

const (
	CUTOUT_CIRCLE = iota
	CUTOUT_RECT
)

type CutoutParms struct {
	cutout_type int       // circle or rectangle
	center      V2        // center of cutout
	parms       []float64 // width, height or diameter
}

func front_panel() SDF3 {

	jack_x := 123.0
	midi_x := 18.2
	led_x := 62.7
	pb_x := 52.8

	midi_d := 15.5
	jack_d := 11.5

	cp := []CutoutParms{
		{CUTOUT_CIRCLE, V2{midi_x, 9.3}, []float64{midi_d}},         // MIDI DIN Jack
		{CUTOUT_CIRCLE, V2{midi_x + 20.32, 9.3}, []float64{midi_d}}, // MIDI DIN Jack
		{CUTOUT_CIRCLE, V2{jack_x, 8.14}, []float64{jack_d}},        // 1/4" Stereo Jack
		{CUTOUT_CIRCLE, V2{jack_x + 19.5, 8.14}, []float64{jack_d}}, // 1/4" Stereo Jack
		{CUTOUT_RECT, V2{led_x, 0.5}, []float64{1.6, 1.6}},          // LED
		{CUTOUT_RECT, V2{led_x + 3.635, 0.5}, []float64{1.6, 1.6}},  // LED
		{CUTOUT_RECT, V2{pb_x, 0.8}, []float64{3.5, 1.6}},           // Push Button
		{CUTOUT_RECT, V2{pb_x + 5.334, 0.8}, []float64{3.5, 1.6}},   // Push Button

		{CUTOUT_CIRCLE, V2{107.4, 2.3}, []float64{5.5}},    // 3.5 mm Headphone Jack
		{CUTOUT_RECT, V2{84.1, 1.0}, []float64{14.3, 2.0}}, // micro SD card
		{CUTOUT_RECT, V2{96.7, 1.3}, []float64{8.0, 3.1}},  // micro USB connector
		{CUTOUT_RECT, V2{72.6, 7.6}, []float64{7.1, 14.8}}, // fullsize USB connector
	}

	s := make([]SDF2, len(cp))
	for i, p := range cp {
		var s0 SDF2
		if p.cutout_type == CUTOUT_CIRCLE {
			r := p.parms[0] / 2.0
			s0 = Circle2D(r)
		} else if p.cutout_type == CUTOUT_RECT {
			w := p.parms[0]
			h := p.parms[1]
			s0 = Box2D(V2{w, h}, 0.0)
		} else {
			panic("bad cutout type")
		}
		s[i] = Transform2D(s0, Translate2d(p.center))
	}
	cutouts := Union2D(s...)

	// overall panel
	pp := &PanelParms{
		Size:         V2{front_panel_length, front_panel_height},
		CornerRadius: front_panel_radius,
		HoleRadius:   1.0,
		HoleOffset:   5.0,
	}
	panel := Panel2D(pp)

	x_ofs := 0.5 * pcb_length
	y_ofs := (0.5 * front_panel_height) - pcb_thickness - pillar_height - base_thickness
	panel = Transform2D(panel, Translate2d(V2{x_ofs, y_ofs}))

	return Extrude3D(Difference2D(panel, cutouts), front_panel_thickness)
}

//-----------------------------------------------------------------------------

func main() {
	s0 := front_panel()
	s0 = Transform3D(s0, Translate3d(V3{0, 80, 0}))
	s1 := base()
	RenderSTL(Union3D(s0, s1), 400, "x.stl")

	//RenderSTL(front_panel(), 400, "front.stl")
	//RenderSTL(base(), 400, "base.stl")

	fb := &FingerButtonParms{
		Width:  3.5,
		Gap:    0.5,
		Length: 30.0,
	}

	RenderDXF(FingerButton2D(fb), 400, "b.dxf")

}

//-----------------------------------------------------------------------------
