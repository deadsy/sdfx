//-----------------------------------------------------------------------------
/*

Phone Holder

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// phone body
var phone_w = 78.0  // width
var phone_h = 146.5 // height
var phone_t = 11.6  // thickness
var phone_r = 11.2  // corner radius

// camera hole
var camera_w = 23.5 // width
var camera_h = 33.0 // height
var camera_r = 3.0  // corner radius
var camera_xofs = 0.0
var camera_yofs = ((phone_h - camera_h) / 2.0) - 9.0

// speaker hole
var speaker_w = 12.5 // width
var speaker_h = 10.0 // height
var speaker_r = 3.0  // corner radius
var speaker_xofs = 8.5 - ((phone_w - speaker_w) / 2.0)
var speaker_yofs = ((phone_h - speaker_h) / 2.0) - 114.0

// wall thickness
var wall_t = 3.0

//-----------------------------------------------------------------------------

func phone_body() SDF3 {
	s2d := Box2D(V2{phone_w, phone_h}, phone_r)
	s3d := Extrude3D(s2d, phone_t)
	m := Translate3d(V3{0, 0, wall_t / 2.0})
	return Transform3D(s3d, m)
}

func camera_hole() SDF3 {
	s2d := Box2D(V2{camera_w, camera_h}, camera_r)
	s3d := Extrude3D(s2d, wall_t+phone_t)
	m := Translate3d(V3{camera_xofs, camera_yofs, 0})
	return Transform3D(s3d, m)
}

func speaker_hole() SDF3 {
	s2d := Box2D(V2{speaker_w, speaker_h}, speaker_r)
	s3d := Extrude3D(s2d, wall_t+phone_t)
	m := Translate3d(V3{speaker_xofs, speaker_yofs, 0})
	return Transform3D(s3d, m)
}

//-----------------------------------------------------------------------------
// holes for buttons, jacks, etc.

var hole_r = 2.0 // corner radius

func hole_left(length, yofs, zofs float64) SDF3 {
	w := phone_t * 2.0
	xofs := -(phone_w + wall_t) / 2.0
	yofs = (phone_h-length)/2.0 - yofs
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	s2d := Box2D(V2{w, length}, hole_r)
	s3d := Extrude3D(s2d, wall_t)
	m := Translate3d(V3{xofs, yofs, zofs}).Mul(RotateY(DtoR(90)))
	return Transform3D(s3d, m)
}

func hole_right(length, yofs, zofs float64) SDF3 {
	w := phone_t * 2.0
	xofs := (phone_w + wall_t) / 2.0
	yofs = (phone_h-length)/2.0 - yofs
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	s2d := Box2D(V2{w, length}, hole_r)
	s3d := Extrude3D(s2d, wall_t)
	m := Translate3d(V3{xofs, yofs, zofs}).Mul(RotateY(DtoR(90)))
	return Transform3D(s3d, m)
}

func hole_top(length, xofs, zofs float64) SDF3 {
	w := phone_t * 2.0
	xofs = -(phone_w-length)/2.0 + xofs
	yofs := (phone_h + wall_t) / 2.0
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	s2d := Box2D(V2{length, w}, hole_r)
	s3d := Extrude3D(s2d, wall_t)
	m := Translate3d(V3{xofs, yofs, zofs}).Mul(RotateX(DtoR(90)))
	return Transform3D(s3d, m)
}

func hole_bottom(length, xofs, zofs float64) SDF3 {
	w := phone_t * 2.0
	xofs = -(phone_w-length)/2.0 + xofs
	yofs := -(phone_h + wall_t) / 2.0
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	s2d := Box2D(V2{length, w}, hole_r)
	s3d := Extrude3D(s2d, wall_t)
	m := Translate3d(V3{xofs, yofs, zofs}).Mul(RotateX(DtoR(90)))
	return Transform3D(s3d, m)
}

//-----------------------------------------------------------------------------

func outer_shell() SDF3 {
	w := phone_w + (2.0 * wall_t)
	h := phone_h + (2.0 * wall_t)
	r := phone_r + wall_t
	t := phone_t + wall_t
	s2d := Box2D(V2{w, h}, r)
	return Extrude3D(s2d, t)
}

//-----------------------------------------------------------------------------

func additive() SDF3 {
	return Union3D(
		outer_shell(),
	)
}

//-----------------------------------------------------------------------------

func subtractive() SDF3 {
	return Union3D(
		phone_body(),
		camera_hole(),
		speaker_hole(),
		hole_left(31.0, 16.5, 8.0),
		hole_right(20.0, 32.0, 8.0),
		hole_top(13.0, 15.0, 8.0),
		hole_top(13.0, 49.5, 9.0),
		hole_bottom(35.0, 17.0, 9.0),
	)
}

//-----------------------------------------------------------------------------

func main() {
	s := Difference3D(additive(), subtractive())
	RenderSTL(s, 300, "holder.stl")
}

//-----------------------------------------------------------------------------
