//-----------------------------------------------------------------------------
/*

3D Printable Nuts and Bolts

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// const MM_TOLERANCE = 0.4 // a bit loose
// const MM_TOLERANCE = 0.2 // very tight
const MM_TOLERANCE = 0.3
const INCH_TOLERANCE = MM_TOLERANCE / MM_PER_INCH
const QUALITY = 300 // stl mesh size

//-----------------------------------------------------------------------------

// Return a Bolt
func Bolt(
	name string, // name of thread
	style string, // head style hex,knurl
	tolerance float64, // subtract from external thread radius
	total_length float64, // threaded length + shank length
	shank_length float64, //  non threaded length
) SDF3 {

	t := ThreadLookup(name)
	if total_length < 0 {
		return nil
	}
	if shank_length < 0 {
		return nil
	}
	thread_length := total_length - shank_length
	if thread_length < 0 {
		thread_length = 0
	}

	var head_3d SDF3
	head_r := t.Hex_Radius()
	head_h := t.Hex_Height()
	if style == "hex" {
		head_3d = HexHead3D(head_r, head_h, "b")
	} else if style == "knurl" {
		head_3d = KnurledHead3D(head_r, head_h, head_r*0.25)
	} else {
		panic("unknown style")
	}

	// shank
	shank_length += head_h / 2
	shank_ofs := shank_length / 2
	shank_3d := Cylinder3D(shank_length, t.Radius, head_h*0.08)
	shank_3d = Transform3D(shank_3d, Translate3d(V3{0, 0, shank_ofs}))

	// thread
	r := t.Radius - tolerance
	l := thread_length
	screw_ofs := l/2 + shank_length
	screw_3d := Screw3D(ISOThread(r, t.Pitch, "external"), l, t.Pitch, 1)

	// chamfer the thread
	p := NewPolygon()
	p.Add(0, -l/2)
	p.Add(r, -l/2)
	p.Add(r, l/2).Chamfer(r / 2)
	p.Add(0, l/2)
	screw_3d = Intersect3D(screw_3d, Revolve3D(Polygon2D(p.Vertices())))
	screw_3d = Transform3D(screw_3d, Translate3d(V3{0, 0, screw_ofs}))

	return Union3D(head_3d, screw_3d, shank_3d)
}

//-----------------------------------------------------------------------------

// Return a Nut.
func Nut(
	name string, // name of thread
	style string, // head style hex,knurl
	tolerance float64, // add to internal thread radius
) SDF3 {

	t := ThreadLookup(name)

	var nut_3d SDF3
	nut_r := t.Hex_Radius()
	nut_h := t.Hex_Height()
	if style == "hex" {
		nut_3d = HexHead3D(nut_r, nut_h, "tb")
	} else if style == "knurl" {
		nut_3d = KnurledHead3D(nut_r, nut_h, nut_r*0.25)
	} else {
		panic("unknown style")
	}

	// internal thread
	thread_3d := Screw3D(ISOThread(t.Radius+tolerance, t.Pitch, "internal"), nut_h, t.Pitch, 1)

	return Difference3D(nut_3d, thread_3d)
}

//-----------------------------------------------------------------------------

func inch() {
	// bolt
	bolt_3d := Bolt("unc_1/2", "knurl", INCH_TOLERANCE, 2.0, 0)
	bolt_3d = Scale3D(bolt_3d, MM_PER_INCH)
	RenderSTL(bolt_3d, QUALITY, "bolt.stl")
	// nut
	nut_3d := Nut("unc_1/2", "knurl", INCH_TOLERANCE)
	nut_3d = Scale3D(nut_3d, MM_PER_INCH)
	RenderSTL(nut_3d, QUALITY, "nut.stl")
}

//-----------------------------------------------------------------------------

func metric() {
	// bolt
	bolt_3d := Bolt("M16x2", "knurl", MM_TOLERANCE, 50, 10)
	RenderSTL(bolt_3d, QUALITY, "bolt.stl")
	// nut
	nut_3d := Nut("M16x2", "knurl", MM_TOLERANCE)
	RenderSTL(nut_3d, QUALITY, "nut.stl")
}

//-----------------------------------------------------------------------------

func main() {
	inch()
	//metric()
}

//-----------------------------------------------------------------------------
