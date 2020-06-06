//-----------------------------------------------------------------------------
/*

4 Part Panel Box

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"strings"
)

//-----------------------------------------------------------------------------

type boxTabParms struct {
	Wall        float64 // wall thickness
	Length      float64 // tab length
	Hole        float64 // hole diameter >= 0 gives a larger tab with a screw hole
	HoleOffset  float64 // hole offset
	Orientation string  // orientation of tab
	Clearance   float64 // fit clearance (typically 0.05)
}

// boxTab3d returns an oriented tab for the box side.
func boxTab3d(k *boxTabParms) SDF3 {

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
		m = m.Mul(RotateX(Pi))
	case "tl": // top, left
		m = m.Mul(Translate3d(V3{(0.5 - k.Clearance) * w, 0, -0.5 * k.Length}))
		m = m.Mul(RotateY(DtoR(-90)))
	case "br": // bottom, right
		m = m.Mul(Translate3d(V3{(-0.5 + k.Clearance) * w, 0, -0.5 * k.Length}))
		m = m.Mul(RotateY(DtoR(-90)))
		m = m.Mul(RotateX(Pi))
	case "tr": // top, right
		m = m.Mul(Translate3d(V3{(-0.5 + k.Clearance) * w, 0, -0.5 * k.Length}))
		m = m.Mul(RotateY(DtoR(90)))
	default:
		panic("invalid tab orientation")
	}
	return Transform3D(tab, m)
}

//-----------------------------------------------------------------------------

type boxHoleParms struct {
	Length      float64 // total hole length
	Hole        float64 // hole diameter
	ZOffset     float64 // hole offset in z-direction (along body length)
	YOffset     float64 // hole offset in y-direction (along body height)
	Orientation string  // orientation of tab
}

// boxHole3d returns an oriented countersunk hole for the box side.
func boxHole3d(k *boxHoleParms) SDF3 {
	hole := CounterSunkHole3D(k.Length, 0.5*k.Hole)
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
func filterTabs(pattern string, tab rune) string {
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
	panelGap := (1.0 + (4.0 * k.Clearance)) * k.Panel

	midZ := k.Size.Z - k.FrontInset - k.BackInset - 2.0*(panelGap+2.0*k.Wall)
	if midZ <= 0.0 {
		panic("the front and back panel depths exceed the total box length")
	}

	outerSize := V2{k.Size.X, k.Size.Y}
	innerSize := outerSize.SubScalar(2.0 * k.Wall)
	ridgeSize := innerSize.SubScalar(2.0 * k.Wall)

	innerPlusSize := innerSize.AddScalar(2.0 * k.Clearance * k.Wall)
	innerMinusSize := innerSize.SubScalar(4.0 * k.Clearance * k.Wall)
	innerRounding := Max(0.0, k.Rounding-k.Wall)

	outer := Box2D(outerSize, k.Rounding)
	inner := Box2D(innerSize, innerRounding)
	innerPlus := Box2D(innerPlusSize, innerRounding)
	innerMinus := Box2D(innerMinusSize, innerRounding)
	ridge := Box2D(ridgeSize, Max(0.0, k.Rounding-2.0*k.Wall))

	// front/pack panels
	panel := Extrude3D(innerMinus, k.Panel)

	// box
	box := Extrude3D(Difference2D(outer, inner), k.Size.Z)

	// add the panel holding ridges
	pr := Extrude3D(Difference2D(innerPlus, ridge), k.Wall)
	z0 := 0.5*(k.Size.Z-k.Wall) - k.FrontInset
	z1 := z0 - k.Wall - panelGap
	z2 := 0.5*(k.Wall-k.Size.Z) + k.BackInset
	z3 := z2 + k.Wall + panelGap
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

		tabLength := midZ / float64(len(k.SideTabs))
		z0 := 0.5*k.Size.Z - k.FrontInset - 2.0*k.Wall - k.Panel
		z1 := -0.5*k.Size.Z + k.BackInset + 2.0*k.Wall + k.Panel
		x := 0.5*k.Size.X - k.Wall

		tPattern := filterTabs(k.SideTabs, 't')
		bPattern := filterTabs(k.SideTabs, 'b')

		tp := &boxTabParms{
			Wall:      k.Wall,
			Length:    tabLength,
			Clearance: k.Clearance,
		}

		// top panel left side
		tp.Orientation = "tl"
		tlTabs := LineOf3D(boxTab3d(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, tPattern)
		// top panel right side
		tp.Orientation = "tr"
		trTabs := LineOf3D(boxTab3d(tp), V3{x, 0, z0}, V3{x, 0, z1}, tPattern)
		// add tabs to the top panel
		top = Union3D(top, tlTabs, trTabs)

		// bottom panel left side
		tp.Orientation = "bl"
		blTabs := LineOf3D(boxTab3d(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, bPattern)
		// bottom panel right side
		tp.Orientation = "br"
		brTabs := LineOf3D(boxTab3d(tp), V3{x, 0, z0}, V3{x, 0, z1}, bPattern)
		// add tabs to the bottom panel
		bottom = Union3D(bottom, blTabs, brTabs)

		if k.Hole > 0 {
			// tabs with holes
			tPattern := filterTabs(k.SideTabs, 'T')
			bPattern := filterTabs(k.SideTabs, 'B')

			holeOffset := 2.0 * k.Wall

			// tabs
			tp := &boxTabParms{
				Wall:       k.Wall,
				Length:     tabLength,
				Hole:       0.85 * k.Hole,
				HoleOffset: holeOffset,
				Clearance:  k.Clearance,
			}

			// top panel left side
			tp.Orientation = "tl"
			tlTabs := LineOf3D(boxTab3d(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, tPattern)
			// top panel right side
			tp.Orientation = "tr"
			trTabs := LineOf3D(boxTab3d(tp), V3{x, 0, z0}, V3{x, 0, z1}, tPattern)
			// add tabs to the top panel
			top = Union3D(top, tlTabs, trTabs)

			// bottom panel left side
			tp.Orientation = "bl"
			blTabs := LineOf3D(boxTab3d(tp), V3{-x, 0, z0}, V3{-x, 0, z1}, bPattern)
			// bottom panel right side
			tp.Orientation = "br"
			brTabs := LineOf3D(boxTab3d(tp), V3{x, 0, z0}, V3{x, 0, z1}, bPattern)
			// add tabs to the bottom panel
			bottom = Union3D(bottom, blTabs, brTabs)

			// holes
			hp := &boxHoleParms{
				Length:  k.Wall,
				Hole:    k.Hole,
				ZOffset: 0.5 * tabLength,
				YOffset: holeOffset,
			}

			// top panel left side
			hp.Orientation = "tl"
			tlHoles := LineOf3D(boxHole3d(hp), V3{-x, 0, z0}, V3{-x, 0, z1}, bPattern)
			// top panel right side
			hp.Orientation = "tr"
			trHoles := LineOf3D(boxHole3d(hp), V3{x, 0, z0}, V3{x, 0, z1}, bPattern)
			// add holes to the top panel
			tHoles := Union3D(tlHoles, trHoles)
			top = Difference3D(top, tHoles)

			// bottom panel left side
			hp.Orientation = "bl"
			blHoles := LineOf3D(boxHole3d(hp), V3{-x, 0, z0}, V3{-x, 0, z1}, tPattern)
			// bottom panel right side
			hp.Orientation = "br"
			brHoles := LineOf3D(boxHole3d(hp), V3{x, 0, z0}, V3{x, 0, z1}, tPattern)
			// add holes to the bottom panel
			bHoles := Union3D(blHoles, brHoles)
			bottom = Difference3D(bottom, bHoles)
		}
	}

	return []SDF3{panel, top, bottom}
}

//-----------------------------------------------------------------------------
