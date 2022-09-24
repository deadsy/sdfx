//-----------------------------------------------------------------------------
/*

Phone Holder

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// phone body
var phone_w = 78.0  // width
var phone_h = 146.5 // height
var phone_t = 13.0  // thickness
var phone_r = 12.0  // corner radius

// camera hole
var camera_w = 23.5 // width
var camera_h = 33.0 // height
var camera_r = 3.0  // corner radius
var camera_xofs = 0.0
var camera_yofs = 48.0

// speaker hole
var speaker_w = 12.5 // width
var speaker_h = 10.0 // height
var speaker_r = 3.0  // corner radius
var speaker_xofs = 23.0
var speaker_yofs = -46.0

// wall thickness
var wall_t = 3.0

//-----------------------------------------------------------------------------

func phone_body() sdf.SDF3 {
	s2d := sdf.Box2D(v2.Vec{phone_w, phone_h}, phone_r)
	s3d := sdf.Extrude3D(s2d, phone_t)
	m := sdf.Translate3d(v3.Vec{0, 0, wall_t / 2.0})
	return sdf.Transform3D(s3d, m)
}

func camera_hole() sdf.SDF3 {
	s2d := sdf.Box2D(v2.Vec{camera_w, camera_h}, camera_r)
	s3d := sdf.Extrude3D(s2d, wall_t+phone_t)
	m := sdf.Translate3d(v3.Vec{camera_xofs, camera_yofs, 0})
	return sdf.Transform3D(s3d, m)
}

func speaker_hole() sdf.SDF3 {
	s2d := sdf.Box2D(v2.Vec{speaker_w, speaker_h}, speaker_r)
	s3d := sdf.Extrude3D(s2d, wall_t+phone_t)
	m := sdf.Translate3d(v3.Vec{speaker_xofs, speaker_yofs, 0})
	return sdf.Transform3D(s3d, m)
}

//-----------------------------------------------------------------------------
// holes for buttons, jacks, etc.

var hole_r = 2.0 // corner radius

func hole_left(length, yofs, zofs float64) sdf.SDF3 {
	w := phone_t * 2.0
	xofs := -(phone_w + wall_t) / 2.0
	yofs = (phone_h-length)/2.0 - yofs
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	s2d := sdf.Box2D(v2.Vec{w, length}, hole_r)
	s3d := sdf.Extrude3D(s2d, wall_t)
	m := sdf.Translate3d(v3.Vec{xofs, yofs, zofs}).Mul(sdf.RotateY(sdf.DtoR(90)))
	return sdf.Transform3D(s3d, m)
}

func hole_right(length, yofs, zofs float64) sdf.SDF3 {
	w := phone_t * 2.0
	xofs := (phone_w + wall_t) / 2.0
	yofs = (phone_h-length)/2.0 - yofs
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	s2d := sdf.Box2D(v2.Vec{w, length}, hole_r)
	s3d := sdf.Extrude3D(s2d, wall_t)
	m := sdf.Translate3d(v3.Vec{xofs, yofs, zofs}).Mul(sdf.RotateY(sdf.DtoR(90)))
	return sdf.Transform3D(s3d, m)
}

func hole_top(length, xofs, zofs float64) sdf.SDF3 {
	w := phone_t * 2.0
	xofs = -(phone_w-length)/2.0 + xofs
	yofs := (phone_h + wall_t) / 2.0
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	s2d := sdf.Box2D(v2.Vec{length, w}, hole_r)
	s3d := sdf.Extrude3D(s2d, wall_t)
	m := sdf.Translate3d(v3.Vec{xofs, yofs, zofs}).Mul(sdf.RotateX(sdf.DtoR(90)))
	return sdf.Transform3D(s3d, m)
}

func hole_bottom(length, xofs, zofs float64) sdf.SDF3 {
	w := phone_t * 2.0
	xofs = -(phone_w-length)/2.0 + xofs
	yofs := -(phone_h + wall_t) / 2.0
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	s2d := sdf.Box2D(v2.Vec{length, w}, hole_r)
	s3d := sdf.Extrude3D(s2d, wall_t)
	m := sdf.Translate3d(v3.Vec{xofs, yofs, zofs}).Mul(sdf.RotateX(sdf.DtoR(90)))
	return sdf.Transform3D(s3d, m)
}

//-----------------------------------------------------------------------------

func outer_shell() sdf.SDF3 {
	w := phone_w + (2.0 * wall_t)
	h := phone_h + (2.0 * wall_t)
	r := phone_r + wall_t
	t := phone_t + wall_t
	s2d := sdf.Box2D(v2.Vec{w, h}, r)
	return sdf.Extrude3D(s2d, t)
}

//-----------------------------------------------------------------------------

func clip() sdf.SDF3 {
	theta := 35.0
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(12.0, 0).Rel()
	p.Add(0, 2.0).Rel()
	p.Add(-10.0, 0).Rel()
	p.Add(0, 4.5).Rel()
	p.Add(-19.5411, 0).Rel()
	p.Add(14.8717, sdf.DtoR(270.0-theta)).Polar().Rel()
	p.Add(0, -7.8612).Rel()
	p.Add(4.3306, sdf.DtoR(270.0+theta)).Polar().Rel()
	p.Add(2.0, sdf.DtoR(theta)).Polar().Rel()
	p.Add(3.7, sdf.DtoR(90.0+theta)).Polar().Rel()
	p.Add(0, 6.6).Rel()
	p.Add(13.2, sdf.DtoR(90.0-theta)).Polar().Rel()
	p.Add(16.5, 0).Rel()
	// back to the the start with a closed polygon
	p.Close()
	//p.Render("clip.dxf")
	s, _ := sdf.Polygon2D(p.Vertices())
	return sdf.Extrude3D(s, 8.0)
}

//-----------------------------------------------------------------------------

func additive() sdf.SDF3 {
	return sdf.Union3D(
		outer_shell(),
	)
}

//-----------------------------------------------------------------------------

func subtractive() sdf.SDF3 {
	return sdf.Union3D(
		phone_body(),
		camera_hole(),
		speaker_hole(),
		hole_left(31.0, 19.5, 8.0),
		hole_right(20.0, 34.0, 8.0),
		hole_top(13.0, 16.0, 8.0),
		hole_top(13.0, 49.5, 9.0),
		hole_bottom(35.0, 20.5, 9.0),
	)
}

//-----------------------------------------------------------------------------

func main() {
	render.RenderSTL(clip(), 300, "clip.stl")
	s := sdf.Difference3D(additive(), subtractive())
	render.RenderSTL(s, 300, "holder.stl")
}

//-----------------------------------------------------------------------------
