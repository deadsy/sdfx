package main

import . "github.com/deadsy/sdfx/sdf"

func main() {

	t := ThreadLookup("unc_1/4")
	length := 1.0 // overall length
	shank := 0.25 // unthreaded shank length

	// hex
	z_ofs := 0.5 * (length + shank + t.Hex_Height)
	hex_2d := NewPolySDF2(Nagon(6, t.Hex_Radius))
	hex_3d := NewExtrudeSDF3(hex_2d, t.Hex_Height)
	hex_3d = NewTransformSDF3(hex_3d, Translate3d(V3{0, 0, z_ofs}))

	// shank
	z_ofs = 0.5 * length
	shank_3d := NewCylinderSDF3(shank, t.Radius, 0)
	shank_3d = NewTransformSDF3(shank_3d, Translate3d(V3{0, 0, z_ofs}))

	// screw
	screw_3d := NewScrewSDF3(ISOThread(t.Radius, t.Pitch), length-shank, t.Pitch, 1)

	s := NewUnionSDF3(hex_3d, screw_3d)
	s = NewUnionSDF3(s, shank_3d)

	RenderSTL(s, 400, "screw.stl")
}
