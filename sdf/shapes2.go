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
func Panel2D(p *PanelParms) SDF2 {
	// panel
	s0 := Box2D(p.Size, p.CornerRadius)
	if p.HoleRadius <= 0.0 {
		// no holes
		return s0
	}
	// holes
	x := 0.5*p.Size.X - p.HoleOffset
	y := 0.5*p.Size.Y - p.HoleOffset
	positions := V2Set{
		{x, y},
		{x, -y},
		{-x, y},
		{-x, -y},
	}
	holes := make([]SDF2, len(positions))
	for i, posn := range positions {
		holes[i] = Transform2D(Circle2D(p.HoleRadius), Translate2d(V2{posn.X, posn.Y}))
	}
	return Difference2D(s0, Union2D(holes...))
}

//-----------------------------------------------------------------------------
