//-----------------------------------------------------------------------------
/*

Phone Holder

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var phone_l = 147.0
var phone_w = 78.0
var phone_h = 11.6
var phone_r = 11.0

func phone(z_ofs float64) SDF3 {
	s2d := Box2D(V2{phone_l, phone_w}, phone_r)
	s3d := Extrude3D(s2d, phone_h)
	m := Translate3d(V3{0, 0, z_ofs})
	return Transform3D(s3d, m)
}

//-----------------------------------------------------------------------------

var wall_thickness = 3.0

func outer_shell() SDF3 {
	l := phone_l + (2.0 * wall_thickness)
	w := phone_w + (2.0 * wall_thickness)
	r := phone_r + wall_thickness
	h := phone_h + wall_thickness
	s2d := Box2D(V2{l, w}, r)
	return Extrude3D(s2d, h)
}

//-----------------------------------------------------------------------------

func side_cutout() SDF3 {
	x := (phone_l / 2.0) + wall_thickness
	l := phone_l * 0.3
	y0 := phone_w * 0.2
	y1 := phone_w * 0.1

	b := NewBezier()
	b.Add(x, y0)
	b.Add(x-l, y1).Mid()
	b.Add(x-l, -y1).Mid()
	b.Add(x, -y0)
	b.Close()

	s2d := Polygon2D(b.Polygon().Vertices())
	return Extrude3D(s2d, phone_h+wall_thickness)
}

func side_cutouts() SDF3 {
	s0 := side_cutout()
	m := RotateZ(DtoR(180))
	s1 := Transform3D(s0, m)
	return Union3D(s0, s1)
}

//-----------------------------------------------------------------------------

func top_cutout(x_ofs float64) SDF3 {

	y := (phone_w / 2.0) + wall_thickness
	h := phone_w * 0.3
	x0 := phone_w * 0.2
	x1 := x0 * 0.8

	b := NewBezier()
	b.Add(x0, y)
	b.Add(-x0, y)
	b.Add(-x1, y-h).Mid()
	b.Add(x1, y-h).Mid()
	b.Close()

	s2d := Polygon2D(b.Polygon().Vertices())
	s2d = Transform2D(s2d, Translate2d(V2{x_ofs, 0}))
	return Extrude3D(s2d, phone_h+wall_thickness)
}

func top_cutouts() SDF3 {
	x_ofs := phone_l * 0.3
	return Union3D(top_cutout(x_ofs), top_cutout(-x_ofs))
}

//-----------------------------------------------------------------------------

func additive() SDF3 {
	return Union3D(outer_shell())
}

//-----------------------------------------------------------------------------

func subtractive() SDF3 {
	return Union3D(
		phone(wall_thickness/2.0),
		side_cutouts(),
		top_cutouts(),
	)
}

//-----------------------------------------------------------------------------

func main() {
	s := Difference3D(additive(), subtractive())
	RenderSTL(s, 300, "holder.stl")
}

//-----------------------------------------------------------------------------
