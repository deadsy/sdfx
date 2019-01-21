//-----------------------------------------------------------------------------
/*

Common 2D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

// PanelParms defines the parameters for a 2D panel.
type PanelParms struct {
	Size         V2
	CornerRadius float64
	HoleDiameter float64
	HoleMargin   [4]float64 // top, right, bottom, left
	HolePattern  [4]string  // top, right, bottom, left
}

// Panel2D returns a 2d panel with holes on the edges.
func Panel2D(k *PanelParms) SDF2 {
	// panel
	s0 := Box2D(k.Size, k.CornerRadius)
	if k.HoleDiameter <= 0.0 {
		// no holes
		return s0
	}

	// corners
	tl := V2{-0.5*k.Size.X + k.HoleMargin[3], 0.5*k.Size.Y - k.HoleMargin[0]}
	tr := V2{0.5*k.Size.X - k.HoleMargin[1], 0.5*k.Size.Y - k.HoleMargin[0]}
	br := V2{0.5*k.Size.X - k.HoleMargin[1], -0.5*k.Size.Y + k.HoleMargin[2]}
	bl := V2{-0.5*k.Size.X + k.HoleMargin[3], -0.5*k.Size.Y + k.HoleMargin[2]}

	// holes
	hole := Circle2D(0.5 * k.HoleDiameter)
	var holes []SDF2
	// clockwise: top, right, bottom, left
	holes = append(holes, LineOf2D(hole, tl, tr, k.HolePattern[0]))
	holes = append(holes, LineOf2D(hole, tr, br, k.HolePattern[1]))
	holes = append(holes, LineOf2D(hole, br, bl, k.HolePattern[2]))
	holes = append(holes, LineOf2D(hole, bl, tl, k.HolePattern[3]))

	return Difference2D(s0, Union2D(holes...))
}

//-----------------------------------------------------------------------------
// finger button

// FingerButtonParms defines the parameters for a 2D finger button.
type FingerButtonParms struct {
	Width  float64 // finger width
	Gap    float64 // gap between finger and body
	Length float64 // length of the finger
}

// FingerButton2D returns a 2D cutout for a finger button.
func FingerButton2D(k *FingerButtonParms) SDF2 {
	r0 := 0.5 * k.Width
	r1 := r0 - k.Gap
	l := 2.0 * k.Length
	s := Difference2D(Line2D(l, r0), Line2D(l, r1))
	s = Cut2D(s, V2{0, 0}, V2{0, 1})
	return Transform2D(s, Translate2d(V2{-k.Length, 0}))
}

//-----------------------------------------------------------------------------
