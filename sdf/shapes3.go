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

// CounterBored_Hole3D returns the SDF3 for a counterbored hole.
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

// Chamfered_Hole3D returns the SDF3 for a chamfered hole (45 degrees).
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

// CounterSunk_Hole3D returns the SDF3 for a countersunk hole (45 degrees).
func CounterSunk_Hole3D(
	l float64, // total length
	r float64, // hole radius
) SDF3 {
	return Chamfered_Hole3D(l, r, r)
}

//-----------------------------------------------------------------------------

// HexHead3D returns the rounded hex head for a nut or bolt.
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

// KnurledHead3D returns a cylindrical knurled head.
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

// KnurlProfile returns a 2D knurl profile.
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

// Knurl3D returns a knurled cylinder.
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

// Washer3D returns a washer.
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

// StandoffParms defines the parameters for a board standoff pillar.
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

// Standoff3D returns a single board standoff.
func Standoff3D(k *StandoffParms) SDF3 {
	s0 := Difference3D(Union3D(pillar(k), pillar_webs(k)), pillar_hole(k))
	if k.NumberWebs != 0 {
		// Cut off any part of the webs that protrude from the top of the pillar
		s1 := Cylinder3D(k.PillarHeight, k.WebDiameter, 0)
		return Intersect3D(s0, s1)
	}
	return s0
}

// Standoffs3D returns multiple board standoffs at various positions.
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

type box_tab_parms struct {
	Wall        float64 // wall thickness
	Length      float64 // tab length
	Hole        float64 // hole diameter >= 0 gives a larger tab with a screw hole
	HoleOffset  float64 // hole offset
	Orientation string  // orientation of tab
	Clearance   float64 // fit clearance (typically 0.05)
}

// box_tab_3d returns an oriented tab for the box side.
func box_tab_3d(k *box_tab_parms) SDF3 {

	w := k.Wall
	l := (1.0 - 2.0*k.Clearance) * k.Length

	var h float64
	if k.Hole > 0 {
		h = 6.0 * k.Wall
	} else {
		h = 4.0 * k.Wall
	}

	tab := Extrude3D(Box2D(V2{l, h}, 0.25*h), w)
	// add a slope where the tab attaches to the box, avoiding overhangs.
	tab = Cut3D(tab, V3{0, 0.5 * h, 0.5 * w}, V3{0, -1, 1})

	// add a cutout to give some tab/body clearance
	w1 := 2.0 * k.Clearance * w
	cutout := Box3D(V3{l, h - 2.0*k.Wall, w1}, 0)
	cutout = Transform3D(cutout, Translate3d(V3{0, -w, 0.5 * (w - w1)}))
	tab = Difference3D(tab, cutout)

	if k.Hole > 0 {
		// adjust the tab, 4 * k.Wall above, 2 * k.Wall below
		tab = Transform3D(tab, Translate3d(V3{0, -w, 0}))
		// put a hole in the tab
		hole := Cylinder3D(w, 0.5*k.Hole, 0)
		hole = Transform3D(hole, Translate3d(V3{0, -k.HoleOffset, 0}))
		tab = Difference3D(tab, hole)
	}

	m := Identity3d()
	switch k.Orientation {
	case "bl": // bottom, left
		m = m.Mul(Translate3d(V3{(0.5 - k.Clearance) * w, 0, -0.5 * k.Length}))
		m = m.Mul(RotateY(DtoR(90)))
		m = m.Mul(RotateX(PI))
	case "tl": // top, left
		m = m.Mul(Translate3d(V3{(0.5 - k.Clearance) * w, 0, -0.5 * k.Length}))
		m = m.Mul(RotateY(DtoR(-90)))
	case "br": // bottom, right
		m = m.Mul(Translate3d(V3{(-0.5 + k.Clearance) * w, 0, -0.5 * k.Length}))
		m = m.Mul(RotateY(DtoR(-90)))
		m = m.Mul(RotateX(PI))
	case "tr": // top, right
		m = m.Mul(Translate3d(V3{(-0.5 + k.Clearance) * w, 0, -0.5 * k.Length}))
		m = m.Mul(RotateY(DtoR(90)))
	default:
		panic("invalid tab orientation")
	}
	return Transform3D(tab, m)
}

//-----------------------------------------------------------------------------

type box_hole_parms struct {
	Length      float64 // total hole length
	Hole        float64 // hole diameter
	ZOffset     float64 // hole offset in z-direction (along body length)
	YOffset     float64 // hole offset in y-direction (along body height)
	Orientation string  // orientation of tab
}

// box_hole_3d returns an oriented countersunk hole for the box side.
func box_hole_3d(k *box_hole_parms) SDF3 {
	hole := CounterSunk_Hole3D(k.Length, 0.5*k.Hole)
	hole = Transform3D(hole, Translate3d(V3{0, 0, 0.5 * k.Length}))
	m := Identity3d()
	switch k.Orientation {
	case "bl": // bottom, left
		m = m.Mul(Translate3d(V3{0, -k.YOffset, -k.ZOffset}))
		m = m.Mul(RotateY(DtoR(-90)))
	case "tl": // top, left
		m = m.Mul(Translate3d(V3{0, k.YOffset, -k.ZOffset}))
		m = m.Mul(RotateY(DtoR(-90)))
	case "br": // bottom, right
		m = m.Mul(Translate3d(V3{0, -k.YOffset, -k.ZOffset}))
		m = m.Mul(RotateY(DtoR(90)))
	case "tr": // top, right
		m = m.Mul(Translate3d(V3{0, k.YOffset, -k.ZOffset}))
		m = m.Mul(RotateY(DtoR(90)))
	default:
		panic("invalid hole orientation")
	}
	return Transform3D(hole, m)
}

//-----------------------------------------------------------------------------
// 4 part panel box

// Convert the tab pattern to "..x.." form with the tab type of interest.
func filter_tabs(pattern string, tab rune) string {
	out := make([]byte, len(pattern))
	for i, c := range pattern {
		if c == tab {
			out[i] = byte('x')
		} else {
			out[i] = byte('.')
		}
	}
	return string(out)
}

// PanelBoxParms defines the parameters for a 4 part panel box.
type PanelBoxParms struct {
	Size       V3      // outer box dimensions (width, height, length)
	Wall       float64 // wall thickness
	Panel      float64 // front/back panel thickness
	Rounding   float64 // radius of corner rounding
	FrontInset float64 // inset depth of box front
	BackInset  float64 // inset depth of box back
	Clearance  float64 // fit clearance (typically 0.05)
	Hole       float64 // diameter of screw holes
	SideTabs   string  // tab pattern b/B (bottom) t/T (top) . (empty)
}

// PanelBox3D returns a 4 part panel box.
func PanelBox3D(k *PanelBoxParms) []SDF3 {
	// sanity checks
	if k.Size.X <= 0 || k.Size.Y <= 0 || k.Size.Z <= 0 {
		panic("invalid box size")
	}
	if k.Wall <= 0 {
		panic("invalid wall size")
	}
	if k.Panel <= 0 {
		panic("invalid panel size")
	}
	if k.Rounding < 0 {
		panic("invalid rounding size")
	}
	if k.FrontInset < 0 || k.BackInset < 0 {
		panic("invalid front/back inset size")
	}
	if k.Clearance < 0 || k.Clearance > 1.0 {
		panic("invalid clearance")
	}
	if k.Clearance == 0 {
		// set a default
		k.Clearance = 0.05
	}
	if k.Hole < 0 {
		panic("invalid screw hole size")
	}
	if k.Hole > 0 {
		if !strings.Contains(k.SideTabs, "T") && !strings.Contains(k.SideTabs, "B") {
			panic("screw hole is non-zero, but there are no screw tabs (T/B)")
		}
	}

	// the panel gap is slightly larger than the panel thickness
	panel_gap := (1.0 + (4.0 * k.Clearance)) * k.Panel

	mid_z := k.Size.Z - k.FrontInset - k.BackInset - 2.0*(panel_gap+2.0*k.Wall)
	if mid_z <= 0.0 {
		panic("the front and back panel depths exceed the total box length")
	}

	outer_size := V2{k.Size.X, k.Size.Y}
	inner_size := outer_size.SubScalar(2.0 * k.Wall)
	ridge_size := inner_size.SubScalar(2.0 * k.Wall)

	inner_plus_size := inner_size.AddScalar(2.0 * k.Clearance * k.Wall)
	inner_minus_size := inner_size.SubScalar(4.0 * k.Clearance * k.Wall)
	inner_rounding := Max(0.0, k.Rounding-k.Wall)

	outer := Box2D(outer_size, k.Rounding)
	inner := Box2D(inner_size, inner_rounding)
	inner_plus := Box2D(inner_plus_size, inner_rounding)
	inner_minus := Box2D(inner_minus_size, inner_rounding)
	ridge := Box2D(ridge_size, Max(0.0, k.Rounding-2.0*k.Wall))

	// front/pack panels
	panel := Extrude3D(inner_minus, k.Panel)

	// box
	box := Extrude3D(Difference2D(outer, inner), k.Size.Z)

	// add the panel holding ridges
	pr := Extrude3D(Difference2D(inner_plus, ridge), k.Wall)
	z0 := 0.5*(k.Size.Z-k.Wall) - k.FrontInset
	z1 := z0 - k.Wall - panel_gap
	z2 := 0.5*(k.Wall-k.Size.Z) + k.BackInset
	z3 := z2 + k.Wall + panel_gap
	pr0 := Transform3D(pr, Translate3d(V3{0, 0, z0}))
	pr1 := Transform3D(pr, Translate3d(V3{0, 0, z1}))
	pr2 := Transform3D(pr, Translate3d(V3{0, 0, z2}))
	pr3 := Transform3D(pr, Translate3d(V3{0, 0, z3}))
	box = Union3D(box, pr0, pr1, pr2, pr3)

	// cut the top and bottom box halves
	top := Cut3D(box, V3{}, V3{0, 1, 0})
	bottom := Cut3D(box, V3{}, V3{0, -1, 0})

	if k.SideTabs != "" {
		// tabs with no holes

		tab_length := mid_z / float64(len(k.SideTabs))
		z0 := 0.5*k.Size.Z - k.FrontInset - 2.0*k.Wall - k.Panel
		z1 := -0.5*k.Size.Z + k.BackInset + 2.0*k.Wall + k.Panel
		x := 0.5*k.Size.X - k.Wall

		t_pattern := filter_tabs(k.SideTabs, 't')
		b_pattern := filter_tabs(k.SideTabs, 'b')

		tp := &box_tab_parms{
			Wall:      k.Wall,
			Length:    tab_length,
			Clearance: k.Clearance,
		}

		// top panel left side
		tp.Orientation = "tl"
		tl_tabs := LineOf3D(box_tab_3d(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, t_pattern)
		// top panel right side
		tp.Orientation = "tr"
		tr_tabs := LineOf3D(box_tab_3d(tp), V3{x, 0, z0}, V3{x, 0, z1}, t_pattern)
		// add tabs to the top panel
		top = Union3D(top, tl_tabs, tr_tabs)

		// bottom panel left side
		tp.Orientation = "bl"
		bl_tabs := LineOf3D(box_tab_3d(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, b_pattern)
		// bottom panel right side
		tp.Orientation = "br"
		br_tabs := LineOf3D(box_tab_3d(tp), V3{x, 0, z0}, V3{x, 0, z1}, b_pattern)
		// add tabs to the bottom panel
		bottom = Union3D(bottom, bl_tabs, br_tabs)

		if k.Hole > 0 {
			// tabs with holes
			t_pattern := filter_tabs(k.SideTabs, 'T')
			b_pattern := filter_tabs(k.SideTabs, 'B')

			hole_offset := 2.0 * k.Wall

			// tabs
			tp := &box_tab_parms{
				Wall:       k.Wall,
				Length:     tab_length,
				Hole:       0.85 * k.Hole,
				HoleOffset: hole_offset,
				Clearance:  k.Clearance,
			}

			// top panel left side
			tp.Orientation = "tl"
			tl_tabs := LineOf3D(box_tab_3d(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, t_pattern)
			// top panel right side
			tp.Orientation = "tr"
			tr_tabs := LineOf3D(box_tab_3d(tp), V3{x, 0, z0}, V3{x, 0, z1}, t_pattern)
			// add tabs to the top panel
			top = Union3D(top, tl_tabs, tr_tabs)

			// bottom panel left side
			tp.Orientation = "bl"
			bl_tabs := LineOf3D(box_tab_3d(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, b_pattern)
			// bottom panel right side
			tp.Orientation = "br"
			br_tabs := LineOf3D(box_tab_3d(tp), V3{x, 0, z0}, V3{x, 0, z1}, b_pattern)
			// add tabs to the bottom panel
			bottom = Union3D(bottom, bl_tabs, br_tabs)

			// holes
			hp := &box_hole_parms{
				Length:  k.Wall,
				Hole:    k.Hole,
				ZOffset: 0.5 * tab_length,
				YOffset: hole_offset,
			}

			// top panel left side
			hp.Orientation = "tl"
			tl_holes := LineOf3D(box_hole_3d(hp), V3{-x, 0, z0}, V3{-x, 0, z1}, b_pattern)
			// top panel right side
			hp.Orientation = "tr"
			tr_holes := LineOf3D(box_hole_3d(hp), V3{x, 0, z0}, V3{x, 0, z1}, b_pattern)
			// add holes to the top panel
			t_holes := Union3D(tl_holes, tr_holes)
			top = Difference3D(top, t_holes)

			// bottom panel left side
			hp.Orientation = "bl"
			bl_holes := LineOf3D(box_hole_3d(hp), V3{-x, 0, z0}, V3{-x, 0, z1}, t_pattern)
			// bottom panel right side
			hp.Orientation = "br"
			br_holes := LineOf3D(box_hole_3d(hp), V3{x, 0, z0}, V3{x, 0, z1}, t_pattern)
			// add holes to the bottom panel
			b_holes := Union3D(bl_holes, br_holes)
			bottom = Difference3D(bottom, b_holes)
		}
	}

	return []SDF3{panel, top, bottom}
}

//-----------------------------------------------------------------------------
