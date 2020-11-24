//-----------------------------------------------------------------------------
/*

2D Panel with Edge Holes and Rounded Corners

*/
//-----------------------------------------------------------------------------

package obj

import "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// PanelParms defines the parameters for a 2D panel.
type PanelParms struct {
	Size         sdf.V2
	CornerRadius float64
	HoleDiameter float64
	HoleMargin   [4]float64 // top, right, bottom, left
	HolePattern  [4]string  // top, right, bottom, left
}

// Panel2D returns a 2d panel with holes on the edges.
func Panel2D(k *PanelParms) sdf.SDF2 {
	// panel
	s0 := sdf.Box2D(k.Size, k.CornerRadius)
	if k.HoleDiameter <= 0.0 {
		// no holes
		return s0
	}

	// corners
	tl := sdf.V2{-0.5*k.Size.X + k.HoleMargin[3], 0.5*k.Size.Y - k.HoleMargin[0]}
	tr := sdf.V2{0.5*k.Size.X - k.HoleMargin[1], 0.5*k.Size.Y - k.HoleMargin[0]}
	br := sdf.V2{0.5*k.Size.X - k.HoleMargin[1], -0.5*k.Size.Y + k.HoleMargin[2]}
	bl := sdf.V2{-0.5*k.Size.X + k.HoleMargin[3], -0.5*k.Size.Y + k.HoleMargin[2]}

	// holes
	hole := sdf.Circle2D(0.5 * k.HoleDiameter)
	var holes []sdf.SDF2
	// clockwise: top, right, bottom, left
	holes = append(holes, sdf.LineOf2D(hole, tl, tr, k.HolePattern[0]))
	holes = append(holes, sdf.LineOf2D(hole, tr, br, k.HolePattern[1]))
	holes = append(holes, sdf.LineOf2D(hole, br, bl, k.HolePattern[2]))
	holes = append(holes, sdf.LineOf2D(hole, bl, tl, k.HolePattern[3]))

	return sdf.Difference2D(s0, sdf.Union2D(holes...))
}

//-----------------------------------------------------------------------------
