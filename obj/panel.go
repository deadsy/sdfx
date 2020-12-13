//-----------------------------------------------------------------------------
/*

2D Panel with rounded corners and edge holes.

Note: The hole pattern is used to layout multiple holes along an edge.

Examples:

"x" - single hole on edge
"xx" - two holes on edge
"x.x" = two holes on edge with spacing
"xx.x.xx" = five holes on edge with spacing
etc.

*/
//-----------------------------------------------------------------------------

package obj

import "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// PanelParms defines the parameters for a 2D panel.
type PanelParms struct {
	Size         sdf.V2     // size of the panel
	CornerRadius float64    // radius of rounded corners
	HoleDiameter float64    // radius of panel holes
	HoleMargin   [4]float64 // hole margins for top, right, bottom, left
	HolePattern  [4]string  // hole pattern for top, right, bottom, left
}

// Panel2D returns a 2d panel with holes on the edges.
func Panel2D(k *PanelParms) (sdf.SDF2, error) {
	// panel
	s0 := sdf.Box2D(k.Size, k.CornerRadius)
	if k.HoleDiameter <= 0.0 {
		// no holes
		return s0, nil
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

	return sdf.Difference2D(s0, sdf.Union2D(holes...)), nil
}

//-----------------------------------------------------------------------------
