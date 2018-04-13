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
	HoleRadius   float64
	HoleOffset   float64
}

// Return a 2d panel with holes in the corners
func Panel2D(k *PanelParms) SDF2 {
	// panel
	s0 := Box2D(k.Size, k.CornerRadius)
	if k.HoleRadius <= 0.0 {
		// no holes
		return s0
	}
	// holes
	x := 0.5*k.Size.X - k.HoleOffset
	y := 0.5*k.Size.Y - k.HoleOffset
	positions := V2Set{
		{x, y},
		{x, -y},
		{-x, y},
		{-x, -y},
	}
	holes := make([]SDF2, len(positions))
	for i, posn := range positions {
		holes[i] = Transform2D(Circle2D(k.HoleRadius), Translate2d(V2{posn.X, posn.Y}))
	}
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
