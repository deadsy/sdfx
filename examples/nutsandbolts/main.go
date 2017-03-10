//-----------------------------------------------------------------------------
/*

Nuts and Bolts

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// Return a hex body for a nut or bolt head.
func hex_body(
	r float64, // radius
	h float64, // height
	rounded int, // number of sides to round 0,1,2
) SDF3 {
	// basic hex body
	corner_round := r * 0.08
	hex_2d := Polygon2D(Nagon(6, r-corner_round))
	hex_2d = Offset2D(hex_2d, corner_round)
	hex_3d := Extrude3D(hex_2d, h)
	// round out the top and/or bottom as required
	if rounded != 0 {
		top_round := r * 1.6
		d := r * math.Cos(DtoR(30))
		sphere_3d := Sphere3D(top_round)
		z_ofs := h/2 - math.Sqrt(top_round*top_round-d*d)
		if rounded >= 1 {
			hex_3d = Intersect3D(hex_3d, Transform3D(sphere_3d, Translate3d(V3{0, 0, -z_ofs})))
		}
		if rounded == 2 {
			hex_3d = Intersect3D(hex_3d, Transform3D(sphere_3d, Translate3d(V3{0, 0, z_ofs})))
		}
	}
	return hex_3d
}

//-----------------------------------------------------------------------------

// Return a Hex Head Bolt
func Hex_Bolt(
	name string, // name of thread
	tolerance float64, // subtract from external thread radius
	total_length float64, // threaded length + shank length
	shank_length float64, //  non threaded length
) SDF3 {

	t := ThreadLookup(name)
	if t == nil {
		return nil
	}
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

	// hex head
	hex_r := t.Hex_Radius()
	hex_h := t.Hex_Height()
	hex_3d := hex_body(hex_r, hex_h, 1)

	// add a rounded cylinder
	hex_3d = Union3D(hex_3d, Cylinder3D(hex_h*1.05, hex_r*0.8, hex_r*0.08))

	// shank
	shank_length += hex_h / 2
	shank_ofs := shank_length / 2
	shank_3d := Cylinder3D(shank_length, t.Radius, hex_r*0.08)
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

	return Union3D(hex_3d, screw_3d, shank_3d)
}

//-----------------------------------------------------------------------------

// Return a Hex Nut
func Hex_Nut(
	name string, // name of thread
	tolerance float64, // add to internal thread radius
	height float64, // height of nut
) SDF3 {

	t := ThreadLookup(name)
	if t == nil {
		return nil
	}
	if height < 0 {
		return nil
	}

	// hex nut body
	hex_3d := hex_body(t.Hex_Radius(), height, 2)

	// internal thread
	thread_3d := Screw3D(ISOThread(t.Radius+tolerance, t.Pitch, "internal"), height, t.Pitch, 1)

	return Difference3D(hex_3d, thread_3d)
}

//-----------------------------------------------------------------------------

func Nut_And_Bolt(
	name string, // name of thread
	tolerance float64, // thread tolerance
	total_length float64, // threaded length + shank length
	shank_length float64, //  non threaded length
) SDF3 {
	t := ThreadLookup(name)
	bolt_3d := Hex_Bolt(name, tolerance, total_length, shank_length)
	nut_3d := Hex_Nut(name, tolerance, t.Hex_Height()/1.5)
	z_ofs := total_length + t.Hex_Height() + 0.25
	nut_3d = Transform3D(nut_3d, Translate3d(V3{0, 0, z_ofs}))
	return Union3D(nut_3d, bolt_3d)
}

//-----------------------------------------------------------------------------

func main() {

	x_ofs := 1.5

	s0 := Nut_And_Bolt("unc_1/4", 0, 2, 0.5)
	s0 = Transform3D(s0, Translate3d(V3{-0.6 * x_ofs, 0, 0}))

	s1 := Nut_And_Bolt("unc_1/2", 0, 2.0, 0.5)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, 0}))

	s2 := Nut_And_Bolt("unc_1", 0, 2.0, 0.5)
	s2 = Transform3D(s2, Translate3d(V3{x_ofs, 0, 0}))

	RenderSTL(Union3D(s0, s1, s2), 400, "nutandbolt.stl")
}

//-----------------------------------------------------------------------------
