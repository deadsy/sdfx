//-----------------------------------------------------------------------------
/*

Common 2D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

type PanelParms struct {
	Size         V2
	CornerRadius float64
	HoleDiameter float64
	HoleMargin   [4]float64 // top, right, bottom, left
	HolePattern  [4]string  // top, right, bottom, left
}

// Return a series of holes along a line.
func line_of_holes(hole SDF2, p0, p1 V2, pattern string) []SDF2 {
	var holes []SDF2
	if pattern != "" {
		x := p0
		dx := p1.Sub(p0).DivScalar(float64(len(pattern)))
		for _, c := range pattern {
			if c == 'x' {
				holes = append(holes, Transform2D(hole, Translate2d(x)))
			}
			x = x.Add(dx)
		}
	}
	return holes
}

// Return a 2d panel with holes on the edges.
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
	holes = append(holes, line_of_holes(hole, tl, tr, k.HolePattern[0])...)
	holes = append(holes, line_of_holes(hole, tr, br, k.HolePattern[1])...)
	holes = append(holes, line_of_holes(hole, br, bl, k.HolePattern[2])...)
	holes = append(holes, line_of_holes(hole, bl, tl, k.HolePattern[3])...)

	return Difference2D(s0, Union2D(holes...))
}

//-----------------------------------------------------------------------------
// finger button

type FingerButtonParms struct {
	Width  float64 // finger width
	Gap    float64 // gap between finger and body
	Length float64 // length of the finger
}

// Return the 2D cutout for a finger button.
func FingerButton2D(k *FingerButtonParms) SDF2 {
	r := 0.5 * k.Width
	l := 2.0 * k.Length
	s := Difference2D(Line2D(l, r+k.Gap), Line2D(l, r))
	s = Cut2D(s, V2{0, 0}, V2{0, 1})
	return Transform2D(s, Translate2d(V2{-k.Length, 0}))
}

//-----------------------------------------------------------------------------
