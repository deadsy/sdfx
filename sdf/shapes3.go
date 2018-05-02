//-----------------------------------------------------------------------------
/*

Common 3D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
	"strings"
)

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
	knurl_h := pitch * math.Floor((h-cylinder_round)/pitch)
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

// Return a washer.
func Washer3D(
	t float64, // thickness
	r_inner float64, // inner radius
	r_outer float64, // outer radius
) SDF3 {
	if t <= 0 {
		panic("t <= 0")
	}
	if r_inner >= r_outer {
		panic("r_inner >= r_outer")
	}
	return Difference3D(Cylinder3D(t, r_outer, 0), Cylinder3D(t, r_inner, 0))
}

//-----------------------------------------------------------------------------
// Board standoffs

type StandoffParms struct {
	PillarHeight   float64
	PillarDiameter float64
	HoleDepth      float64
	HoleDiameter   float64
	NumberWebs     int
	WebHeight      float64
	WebDiameter    float64
	WebWidth       float64
}

// single web
func pillar_web(k *StandoffParms) SDF3 {
	w := NewPolygon()
	w.Add(0, 0)
	w.Add(0.5*k.WebDiameter, 0)
	w.Add(0, k.WebHeight)
	s := Extrude3D(Polygon2D(w.Vertices()), k.WebWidth)
	m := Translate3d(V3{0, 0, -0.5 * k.PillarHeight}).Mul(RotateX(DtoR(90.0)))
	return Transform3D(s, m)
}

// multiple webs
func pillar_webs(k *StandoffParms) SDF3 {
	if k.NumberWebs == 0 {
		return nil
	}
	return RotateCopy3D(pillar_web(k), k.NumberWebs)
}

// pillar
func pillar(k *StandoffParms) SDF3 {
	return Cylinder3D(k.PillarHeight, 0.5*k.PillarDiameter, 0)
}

// pillar hole
func pillar_hole(k *StandoffParms) SDF3 {
	if k.HoleDiameter == 0.0 || k.HoleDepth == 0.0 {
		return nil
	}
	s := Cylinder3D(k.HoleDepth, 0.5*k.HoleDiameter, 0)
	z_ofs := 0.5 * (k.PillarHeight - k.HoleDepth)
	return Transform3D(s, Translate3d(V3{0, 0, z_ofs}))
}

// Return a single board standoff.
func Standoff3D(k *StandoffParms) SDF3 {
	s0 := Difference3D(Union3D(pillar(k), pillar_webs(k)), pillar_hole(k))
	if k.NumberWebs != 0 {
		// Cut off any part of the webs that protrude from the top of the pillar
		s1 := Cylinder3D(k.PillarHeight, k.WebDiameter, 0)
		return Intersect3D(s0, s1)
	}
	return s0
}

// Multiple board standoffs at various positions
func Standoffs3D(k *StandoffParms, positions V3Set) SDF3 {
	if len(positions) == 0 {
		return nil
	}
	s0 := Standoff3D(k)
	if s0 == nil {
		return nil
	}
	s := make([]SDF3, len(positions))
	for i, p := range positions {
		s[i] = Transform3D(s0, Translate3d(p))
	}
	return Union3D(s...)
}

//-----------------------------------------------------------------------------

type BoxTabParms struct {
	Size        V3 // tab dimensions (width, height, length)
	Orientation string
	Clearance   float64 // fit clearance (typically 0.05)
}

func BoxTab3D(k *BoxTabParms) SDF3 {

	w := k.Size.X                       // width
	h := k.Size.Y                       // height
	l := (1.0 - k.Clearance) * k.Size.Z // length

	p := NewPolygon()
	p.Add(0, -0.5*h)
	p.Add(w, -0.25*h)
	p.Add(w, 0.5*h)
	p.Add(0, 0.5*h)
	tab := Extrude3D(Polygon2D(p.Vertices()), l)

	m := Translate3d(V3{0, 0, -0.5 * k.Size.Z})
	switch k.Orientation {
	case "bl": // bottom, left
	case "tl": // top, left
		m = m.Mul(RotateX(PI))
	case "br": // bottom, right
		m = m.Mul(RotateZ(PI))
		m = m.Mul(RotateX(PI))
	case "tr": // top, right
		m = m.Mul(RotateZ(PI))
	default:
		panic("invalid tab orientation")
	}

	return Transform3D(tab, m)
}

//-----------------------------------------------------------------------------
// 4 part panel box

type PanelBoxParms struct {
	Size       V3      // outer box dimensions (width, height, length)
	Wall       float64 // wall thickness
	Panel      float64 // front/back panel thickness
	Rounding   float64 // radius of corner rounding
	FrontInset float64 // inset depth of box front
	BackInset  float64 // inset depth of box back
	Clearance  float64 // fit clearance (typically 0.05)
	SideTabs   string  // side tab pattern ^ (bottom) v (top) . (empty)
}

// PanelBox3D returns a 4 part panel box
func PanelBox3D(k *PanelBoxParms) []SDF3 {

	mid_z := k.Size.Z - k.FrontInset - k.BackInset - 2.0*k.Panel - 4.0*k.Wall
	if mid_z <= 0.0 {
		panic("the front and back panel depths exceed the total box length")
	}

	outer_size := V2{k.Size.X, k.Size.Y}
	inner_size := outer_size.SubScalar(2.0 * k.Wall)
	ridge_size := inner_size.SubScalar(2.0 * k.Wall)

	outer := Box2D(outer_size, k.Rounding)
	inner := Box2D(inner_size, Max(0.0, k.Rounding-k.Wall))
	ridge := Box2D(ridge_size, Max(0.0, k.Rounding-2.0*k.Wall))

	// front/pack panels
	shrink := 1.0 - k.Clearance
	panel := Extrude3D(inner, k.Panel)
	panel = Transform3D(panel, Scale3d(V3{shrink, shrink, shrink}))

	// box
	box := Extrude3D(Difference2D(outer, inner), k.Size.Z)

	// add the panel holding ridges
	pr := Extrude3D(Difference2D(inner, ridge), k.Wall)
	z0 := 0.5*(k.Size.Z-k.Wall) - k.FrontInset
	z1 := z0 - k.Wall - k.Panel
	z2 := 0.5*(k.Wall-k.Size.Z) + k.BackInset
	z3 := z2 + k.Wall + k.Panel
	pr0 := Transform3D(pr, Translate3d(V3{0, 0, z0}))
	pr1 := Transform3D(pr, Translate3d(V3{0, 0, z1}))
	pr2 := Transform3D(pr, Translate3d(V3{0, 0, z2}))
	pr3 := Transform3D(pr, Translate3d(V3{0, 0, z3}))
	box = Union3D(box, pr0, pr1, pr2, pr3)

	// cut the top and bottom box halves
	top := Cut3D(box, V3{}, V3{0, 1, 0})
	bottom := Cut3D(box, V3{}, V3{0, -1, 0})

	if k.SideTabs != "" {

		tab_length := mid_z / float64(len(k.SideTabs))
		z0 := 0.5*k.Size.Z - k.FrontInset - 2.0*k.Wall - k.Panel
		z1 := -0.5*k.Size.Z + k.BackInset + 2.0*k.Wall + k.Panel
		x := 0.5*k.Size.X - k.Wall
		t_pattern := strings.Replace(k.SideTabs, "v", "x", -1)
		t_pattern = strings.Replace(t_pattern, "^", ".", -1)
		b_pattern := strings.Replace(k.SideTabs, "^", "x", -1)
		b_pattern = strings.Replace(b_pattern, "v", ".", -1)

		tp := &BoxTabParms{
			Size:      V3{k.Wall, k.Wall * 4.0, tab_length},
			Clearance: 0.05,
		}

		// top panel left side
		tp.Orientation = "tl"
		tl_tabs := LineOf3D(BoxTab3D(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, t_pattern)
		// top panel right side
		tp.Orientation = "tr"
		tr_tabs := LineOf3D(BoxTab3D(tp), V3{x, 0, z0}, V3{x, 0, z1}, t_pattern)
		// add tabs to the top panel
		top = Union3D(top, tl_tabs, tr_tabs)

		// bottom panel left side
		tp.Orientation = "bl"
		bl_tabs := LineOf3D(BoxTab3D(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, b_pattern)
		// bottom panel right side
		tp.Orientation = "br"
		br_tabs := LineOf3D(BoxTab3D(tp), V3{x, 0, z0}, V3{x, 0, z1}, b_pattern)
		// add tabs to the bottom panel
		bottom = Union3D(bottom, bl_tabs, br_tabs)
	}

	return []SDF3{panel, top, bottom}
}

//-----------------------------------------------------------------------------
