//-----------------------------------------------------------------------------
/*

Axoloti Board Enclosure

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func front_panel() SDF3 {

	var cutouts SDF2

	// 1/4" Stereo Jack (x2)
	stereo_d := 11.2 // front panel cutout
	stereo_y := 8.14 // pcb to center of barrel
	stereo_x0 := 0.0
	stereo_x1 := 19.4
	stereo_r := stereo_d / 2.0
	cutouts = Union2D(cutouts, Transform2D(Circle2D(stereo_r), Translate2d(V2{stereo_x0, stereo_y})))
	cutouts = Union2D(cutouts, Transform2D(Circle2D(stereo_r), Translate2d(V2{stereo_x1, stereo_y})))

	// MIDI DIN Jack (x2)
	midi_d := 15.0 // front panel cutout
	midi_y := 10.0 // pcb to center of connector
	midi_x0 := 103.4
	midi_x1 := 124.0
	midi_r := midi_d / 2.0
	cutouts = Union2D(cutouts, Transform2D(Circle2D(midi_r), Translate2d(V2{midi_x0, midi_y})))
	cutouts = Union2D(cutouts, Transform2D(Circle2D(midi_r), Translate2d(V2{midi_x1, midi_y})))

	// 3.5 mm Headphone Jack
	headphone_d := 5.2 // front panel cutout
	headphone_y := 2.3 // pcb to center of barrel
	headphone_x := 34.9
	headphone_r := headphone_d / 2.0
	cutouts = Union2D(cutouts, Transform2D(Circle2D(headphone_r), Translate2d(V2{headphone_x, headphone_y})))

	// micro SD card
	micro_sd_w := 14.3
	micro_sd_h := 2.0
	micro_sd_x := 58.2
	micro_sd_y := 1.0
	cutouts = Union2D(cutouts, Transform2D(Box2D(V2{micro_sd_w, micro_sd_h}, 0.0), Translate2d(V2{micro_sd_x, micro_sd_y})))

	// micro USB connector
	micro_usb_w := 8.0
	micro_usb_h := 3.1
	micro_usb_x := 45.5
	micro_usb_y := 1.3
	cutouts = Union2D(cutouts, Transform2D(Box2D(V2{micro_usb_w, micro_usb_h}, 0.0), Translate2d(V2{micro_usb_x, micro_usb_y})))

	// fullsize USB connector
	fs_usb_w := 7.1
	fs_usb_h := 14.8
	fs_usb_x := 69.6
	fs_usb_y := 7.6
	cutouts = Union2D(cutouts, Transform2D(Box2D(V2{fs_usb_w, fs_usb_h}, 0.0), Translate2d(V2{fs_usb_x, fs_usb_y})))

	// LEDs (x2)
	led_w := 1.6
	led_h := 1.0
	led_y := 0.5
	led_x0 := 75.8
	led_x1 := 79.4
	cutouts = Union2D(cutouts, Transform2D(Box2D(V2{led_w, led_h}, 0.0), Translate2d(V2{led_x0, led_y})))
	cutouts = Union2D(cutouts, Transform2D(Box2D(V2{led_w, led_h}, 0.0), Translate2d(V2{led_x1, led_y})))

	// Push Buttons (x2)
	pb_w := 3.5
	pb_h := 1.6
	pb_y := 0.8
	pb_x0 := 83.9
	pb_x1 := 89.2
	cutouts = Union2D(cutouts, Transform2D(Box2D(V2{pb_w, pb_h}, 0.0), Translate2d(V2{pb_x0, pb_y})))
	cutouts = Union2D(cutouts, Transform2D(Box2D(V2{pb_w, pb_h}, 0.0), Translate2d(V2{pb_x1, pb_y})))

	// overall panel
	panel_w := 160.0
	panel_h := 35.0
	panel := Box2D(V2{panel_w, panel_h}, 0.0)
	cutouts = Transform2D(cutouts, Translate2d(V2{-60.0, -10.0}))

	return Extrude3D(Difference2D(panel, cutouts), 3.0)
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(front_panel(), 300, "front_panel.stl")
}

//-----------------------------------------------------------------------------
