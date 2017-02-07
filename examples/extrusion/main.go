//-----------------------------------------------------------------------------
/*

Extrusions

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func hex() SDF2 {
	return NewPolySDF2(Nagon(6, 20))
}

func extrude1() SDF3 {

	// make the extrusions
	s_linear := NewExtrudeSDF3(NewOffsetSDF2(hex(), 8), 100)
	s_fwd := NewTwistExtrudeSDF3(NewOffsetSDF2(hex(), 8), 100, TAU)
	s_rev := NewTwistExtrudeSDF3(NewOffsetSDF2(hex(), 8), 100, -TAU)
	s_combo := NewUnionSDF3(s_fwd, s_rev)

	// position them on the y-axis
	d := 60.0
	s_linear = NewTransformSDF3(s_linear, Translate3d(V3{0, -1.5 * d, 0}))
	s_fwd = NewTransformSDF3(s_fwd, Translate3d(V3{0, -0.5 * d, 0}))
	s_rev = NewTransformSDF3(s_rev, Translate3d(V3{0, 0.5 * d, 0}))
	s_combo = NewTransformSDF3(s_combo, Translate3d(V3{0, 1.5 * d, 0}))

	// return a union of them all
	return NewUnionSDF3(s_linear, NewUnionSDF3(s_fwd, NewUnionSDF3(s_rev, s_combo)))
}

func extrude2() SDF3 {
	s0 := NewScaleExtrudeSDF3(NewOffsetSDF2(hex(), 8), 80, V2{.25, .5})
	s1 := NewScaleTwistExtrudeSDF3(NewOffsetSDF2(hex(), 8), 80, PI, V2{.25, .5})

	// position them on the y-axis
	d := 30.0
	s0 = NewTransformSDF3(s0, Translate3d(V3{0, -d, 0}))
	s1 = NewTransformSDF3(s1, Translate3d(V3{0, d, 0}))

	return NewUnionSDF3(s0, s1)
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(extrude1(), "extrude1.stl")
	RenderSTL(extrude2(), "extrude2.stl")
}

//-----------------------------------------------------------------------------
