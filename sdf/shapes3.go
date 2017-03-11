//-----------------------------------------------------------------------------
/*

Common 3D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

// Counter Bored Hole
func CounterBored_Hole3D(
	l float64, // total length
	r float64, // hole radius
	cb_r float64, // counter bore radius
	cb_d float64, // counter bore depth
) SDF3 {
	s0 := Cylinder3D(l, r, 0)
	s1 := Cylinder3D(cb_d, cb_r, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - cb_d) / 2}))
	return Union3D(s0, s1)
}

// Chamfered Hole (45 degrees)
func Chamfered_Hole3D(
	l float64, // total length
	r float64, // hole radius
	ch_r float64, // chamfer radius
) SDF3 {
	s0 := Cylinder3D(l, r, 0)
	s1 := Cone3D(ch_r, r, r+ch_r, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - ch_r) / 2}))
	return Union3D(s0, s1)
}

// Countersunk Hole (45 degrees)
func CounterSunk_Hole3D(
	l float64, // total length
	r float64, // hole radius
) SDF3 {
	return Chamfered_Hole3D(l, r, r)
}

//-----------------------------------------------------------------------------

// Return a rounded hex head for a nut or bolt.
func HexHead3D(
	r float64, // radius
	h float64, // height
	round string, // (t)top, (b)bottom, (tb)top/bottom
) SDF3 {
	// basic hex body
	corner_round := r * 0.08
	hex_2d := Polygon2D(Nagon(6, r-corner_round))
	hex_2d = Offset2D(hex_2d, corner_round)
	hex_3d := Extrude3D(hex_2d, h)
	// round out the top and/or bottom as required
	if round != "" {
		top_round := r * 1.6
		d := r * math.Cos(DtoR(30))
		sphere_3d := Sphere3D(top_round)
		z_ofs := math.Sqrt(top_round*top_round-d*d) - h/2
		if round == "t" || round == "tb" {
			hex_3d = Intersect3D(hex_3d, Transform3D(sphere_3d, Translate3d(V3{0, 0, -z_ofs})))
		}
		if round == "b" || round == "tb" {
			hex_3d = Intersect3D(hex_3d, Transform3D(sphere_3d, Translate3d(V3{0, 0, z_ofs})))
		}
	}
	return hex_3d
}

// Return a cylindrical knurled head.
func KnurledHead3D(
	r float64, // radius
	h float64, // height
	pitch float64, // knurl pitch
) SDF3 {
	theta := DtoR(45)
	cylinder_round := r * 0.05
	// TODO: knurl_h is not correct
	pitch_h := pitch * math.Tan(theta)
	knurl_h := pitch_h * math.Floor((h-2*cylinder_round)/pitch_h)
	knurl_3d := Knurl3D(knurl_h, r, pitch, pitch*0.3, theta)
	return Union3D(Cylinder3D(h, r, cylinder_round), knurl_3d)
}

//-----------------------------------------------------------------------------

// Return a 2D knurl profile.
func KnurlProfile(
	radius float64, // radius of knurled cylinder
	pitch float64, // pitch of the knurl
	height float64, // height of the knurl
) SDF2 {
	knurl := NewPolygon()
	knurl.Add(pitch/2, 0)
	knurl.Add(pitch/2, radius)
	knurl.Add(0, radius+height)
	knurl.Add(-pitch/2, radius)
	knurl.Add(-pitch/2, 0)
	//knurl.Render("knurl.dxf")
	return Polygon2D(knurl.Vertices())
}

// Return a knurled cylinder.
func Knurl3D(
	length float64, // length of cylinder
	radius float64, // radius of cylinder
	pitch float64, // knurl pitch
	height float64, // knurl height
	theta float64, // knurl helix angle
) SDF3 {
	// A knurl is the the intersection of left and right hand
	// multistart "threads". Work out the number of starts using
	// the desired helix angle.
	n := int(TAU * radius * math.Tan(theta) / pitch)
	// build the knurl profile.
	knurl_2d := KnurlProfile(radius, pitch, height)
	// create the left/right hand spirals
	knurl0_3d := Screw3D(knurl_2d, length, pitch, n)
	knurl1_3d := Screw3D(knurl_2d, length, pitch, -n)
	return Intersect3D(knurl0_3d, knurl1_3d)
}

//-----------------------------------------------------------------------------
