package main

import . "github.com/deadsy/sdfx/sdf"

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
