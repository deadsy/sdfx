//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// Create a Hex Head Screw/Bolt
// name = thread name
// total_length = threaded length + shank length
// shank length = non threaded length
func Hex_Screw(name string, total_length, shank_length float64) SDF3 {
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
	z_ofs := 0.5 * (total_length + shank_length + hex_h)
	round := hex_r * 0.08
	hex_2d := Polygon2D(Nagon(6, hex_r-round))
	hex_2d = Offset2D(hex_2d, round)
	hex_3d := Extrude3D(hex_2d, hex_h)
	// round off the edges
	sphere_3d := Sphere3D(hex_r * 1.55)
	sphere_3d = Transform3D(sphere_3d, Translate3d(V3{0, 0, -hex_r * 0.9}))
	hex_3d = Intersection3D(hex_3d, sphere_3d)
	// add a rounded cylinder
	hex_3d = Union3D(hex_3d, Cylinder3D(hex_h*1.05, hex_r*0.8, round))
	hex_3d = Transform3D(hex_3d, Translate3d(V3{0, 0, z_ofs}))

	// shank
	z_ofs = 0.5 * total_length
	shank_3d := Cylinder3D(shank_length, t.Radius, 0)
	shank_3d = Transform3D(shank_3d, Translate3d(V3{0, 0, z_ofs}))

	// thread
	screw_3d := Screw3D(ISOThread(t.Radius, t.Pitch, "external"), thread_length, t.Pitch, 1)

	return Union3D(hex_3d, screw_3d, shank_3d)
}

//-----------------------------------------------------------------------------

func main() {

	x_ofs := 1.5

	s0 := Hex_Screw("unc_1/4", 2.0, 0.5)
	s0 = Transform3D(s0, Translate3d(V3{-0.6 * x_ofs, 0, 0}))

	s1 := Hex_Screw("unc_1/2", 2.0, 0.5)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, 0}))

	s2 := Hex_Screw("unc_1", 2.0, 0.5)
	s2 = Transform3D(s2, Translate3d(V3{x_ofs, 0, 0}))

	s := Union3D(s0, s1)
	s = Union3D(s, s2)

	RenderSTL(s, 400, "screw.stl")
}

//-----------------------------------------------------------------------------
