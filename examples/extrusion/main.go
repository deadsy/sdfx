//-----------------------------------------------------------------------------
/*

Extrusions

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func hex() SDF2 {
	return Polygon2D(Nagon(6, 20))
}

func extrude1() SDF3 {

	// make the extrusions
	s_linear := Extrude3D(Offset2D(hex(), 8), 100)
	s_fwd := TwistExtrude3D(Offset2D(hex(), 8), 100, Tau)
	s_rev := TwistExtrude3D(Offset2D(hex(), 8), 100, -Tau)
	s_combo := Union3D(s_fwd, s_rev)

	// position them on the y-axis
	d := 60.0
	s_linear = Transform3D(s_linear, Translate3d(V3{0, -1.5 * d, 0}))
	s_fwd = Transform3D(s_fwd, Translate3d(V3{0, -0.5 * d, 0}))
	s_rev = Transform3D(s_rev, Translate3d(V3{0, 0.5 * d, 0}))
	s_combo = Transform3D(s_combo, Translate3d(V3{0, 1.5 * d, 0}))

	// return a union of them all
	return Union3D(s_linear, s_fwd, s_rev, s_combo)
}

func extrude2() SDF3 {
	s0 := ScaleExtrude3D(Offset2D(hex(), 8), 80, V2{.25, .5})
	s1 := ScaleTwistExtrude3D(Offset2D(hex(), 8), 80, Pi, V2{.25, .5})

	// position them on the y-axis
	d := 30.0
	s0 = Transform3D(s0, Translate3d(V3{0, -d, 0}))
	s1 = Transform3D(s1, Translate3d(V3{0, d, 0}))

	return Union3D(s0, s1)
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTLSlow(extrude1(), 200, "extrude1.stl")
	RenderSTLSlow(extrude2(), 200, "extrude2.stl")
}

//-----------------------------------------------------------------------------
