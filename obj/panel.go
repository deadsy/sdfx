//-----------------------------------------------------------------------------
/*

Create 2d/3d panels.

*/
//-----------------------------------------------------------------------------

package obj

import "github.com/deadsy/sdfx/sdf"

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
const erHoleDiameter = 3.2

// gaps between adjacent panels (doepfer 3U module spec)
const erUGap = ((3 * erU) - 128.5) * 0.5
const erHPGap = ((3 * erHP) - 15) * 0.5

// EuroRackParms defines the parameters for a eurorack panel.
type EuroRackParms struct {
	U            float64 // U-size (vertical)
	HP           float64 // HP-size (horizontal)
	CornerRadius float64 // radius of panel corners
	HoleDiameter float64 // panel holes (0 for default)
	Thickness    float64 // panel thickness (3d only)
	Ridge        bool    // add side ridges for reinforcing (3d only)
}

func erUSize(u float64) float64 {
	return (u * erU) - (2 * erUGap)
}

func erHPSize(hp float64) float64 {
	return (hp * erHP) - (2 * erHPGap)
}

// EuroRackPanel2D returns a 2d eurorack synthesizer module panel (in mm).
func EuroRackPanel2D(k *EuroRackParms) (sdf.SDF2, error) {

	if k.U < 1 {
		return nil, sdf.ErrMsg("k.U < 1")
	}
	if k.HP <= 1 {
		return nil, sdf.ErrMsg("k.HP <= 1")
	}
	if k.CornerRadius < 0 {
		return nil, sdf.ErrMsg("k.CornerRadius < 0")
	}
	if k.HoleDiameter <= 0 {
		k.HoleDiameter = erHoleDiameter
	}

	// edge to mount hole margins
	const vMargin = 3.0
	const hMargin = (3 * erHP * 0.5) - erHPGap

	x := erHPSize(k.HP)
	y := erUSize(k.U)

	pk := PanelParms{
		Size:         sdf.V2{x, y},
		CornerRadius: k.CornerRadius,
		HoleDiameter: k.HoleDiameter,
		HoleMargin:   [4]float64{vMargin, hMargin, vMargin, hMargin},
	}

	if k.HP < 8 {
		// two holes
		pk.HolePattern = [4]string{"x", "", "", "x"}
	} else {
		// four holes
		pk.HolePattern = [4]string{"x", "x", "x", "x"}
	}

	return Panel2D(&pk)
}

// EuroRackPanel3D returns a 3d eurorack synthesizer module panel (in mm).
func EuroRackPanel3D(k *EuroRackParms) (sdf.SDF3, error) {
	if k.Thickness <= 0 {
		return nil, sdf.ErrMsg("k.Thickness <= 0")
	}
	panel2d, err := EuroRackPanel2D(k)
	if err != nil {
		return nil, err
	}
	s := sdf.Extrude3D(panel2d, k.Thickness)
	if !k.Ridge {
		return s, nil
	}
	// create a reinforcing ridge
	xSize := k.Thickness
	ySize := erUSize(k.U) - 15.0
	zSize := k.Thickness * 1.5
	r, err := sdf.Box3D(sdf.V3{xSize, ySize, zSize}, 0)
	if err != nil {
		return nil, err
	}
	// add the ridges to the sides
	zOfs := 0.5 * (k.Thickness + zSize)
	xOfs := 0.5 * (erHPSize(k.HP) - xSize)
	r = sdf.Transform3D(r, sdf.Translate3d(sdf.V3{0, 0, zOfs}))
	r0 := sdf.Transform3D(r, sdf.Translate3d(sdf.V3{xOfs, 0, 0}))
	r1 := sdf.Transform3D(r, sdf.Translate3d(sdf.V3{-xOfs, 0, 0}))

	return sdf.Union3D(s, r0, r1), nil
}

//-----------------------------------------------------------------------------
