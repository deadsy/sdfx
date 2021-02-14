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
	HoleDiameter float64    // diameter of panel holes
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
	hole, err := sdf.Circle2D(0.5 * k.HoleDiameter)
	if err != nil {
		return nil, err
	}
	var holes []sdf.SDF2
	// clockwise: top, right, bottom, left
	holes = append(holes, sdf.LineOf2D(hole, tl, tr, k.HolePattern[0]))
	holes = append(holes, sdf.LineOf2D(hole, tr, br, k.HolePattern[1]))
	holes = append(holes, sdf.LineOf2D(hole, br, bl, k.HolePattern[2]))
	holes = append(holes, sdf.LineOf2D(hole, bl, tl, k.HolePattern[3]))

	return sdf.Difference2D(s0, sdf.Union2D(holes...)), nil
}

//-----------------------------------------------------------------------------
// EuroRack Module Panels: http://www.doepfer.de/a100_man/a100m_e.htm

const erU = 1.75 * sdf.MillimetresPerInch
const erHP = 0.2 * sdf.MillimetresPerInch

// EuroRackPanel returns a 2d eurorack synthesizer module panel (in mm).
func EuroRackPanel(u, hp, round float64) (sdf.SDF2, error) {

	if u < 1 {
		return nil, sdf.ErrMsg("u < 1")
	}
	if hp <= 1 {
		return nil, sdf.ErrMsg("hp <= 1")
	}
	if round < 0 {
		return nil, sdf.ErrMsg("round < 0")
	}

	// gaps between adjacent panels (doepfer 3U module spec)
	const vGap = ((3 * erU) - 128.5) * 0.5
	const hGap = ((3 * erHP) - 15) * 0.5
	// edge to mount hole margins
	const vMargin = 3.0
	const hMargin = (3 * erHP * 0.5) - hGap
	const holeDiameter = 3.2

	x := (hp * erHP) - (2 * hGap)
	y := (u * erU) - (2 * vGap)

	k := PanelParms{
		Size:         sdf.V2{x, y},
		CornerRadius: round,
		HoleDiameter: holeDiameter,
		HoleMargin:   [4]float64{vMargin, hMargin, vMargin, hMargin},
	}

	if hp < 8 {
		// two holes
		k.HolePattern = [4]string{"x", "", "", "x"}
	} else {
		// four holes
		k.HolePattern = [4]string{"x", "x", "x", "x"}
	}

	return Panel2D(&k)
}

//-----------------------------------------------------------------------------
