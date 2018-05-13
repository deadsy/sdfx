//-----------------------------------------------------------------------------
/*

Nordic NRF52DK Board Mounting Kit

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var base_x = 120.0
var base_y = 70.0
var base_thickness = 3.0

var pcb_x = 102.0
var pcb_y = 63.5

var pillar_height = 15.0

//-----------------------------------------------------------------------------

const MIL = (25.4 / 1000.0)

// multiple standoffs
func standoffs() SDF3 {

	k := &StandoffParms{
		PillarHeight:   pillar_height,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}

	z_ofs := 0.5 * (pillar_height + base_thickness)

	// from the board gerbers
	positions := V3Set{
		{550.0 * MIL, 300.0 * MIL, z_ofs},
		{600.0 * MIL, 2200.0 * MIL, z_ofs},
		{2600.0 * MIL, 1600.0 * MIL, z_ofs},
		{2600.0 * MIL, 500.0 * MIL, z_ofs},
		{3800.0 * MIL, 300.0 * MIL, z_ofs},
	}

	return Standoffs3D(k, positions)
}

//-----------------------------------------------------------------------------

func base() SDF3 {
	// base
	pp := &PanelParms{
		Size:         V2{base_x, base_y},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s0 := Panel2D(pp)

	// cutouts
	c1 := Box2D(V2{53.0, 35.0}, 3.0)
	c1 = Transform2D(c1, Translate2d(V2{-22.0, 3.0}))

	c2 := Box2D(V2{16.0, 40.0}, 3.0)
	c2 = Transform2D(c2, Translate2d(V2{35.0, 5.0}))

	s2 := Extrude3D(Difference2D(s0, Union2D(c1, c2)), base_thickness)
	x_ofs := 0.5 * pcb_x
	y_ofs := pcb_y - (0.5 * base_y)
	s2 = Transform3D(s2, Translate3d(V3{x_ofs, y_ofs, 0}))

	// standoffs
	s3 := standoffs()

	s4 := Union3D(s2, s3)
	s4.(*UnionSDF3).SetMin(PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(base(), 300, "nrf52dk.stl")
}

//-----------------------------------------------------------------------------
