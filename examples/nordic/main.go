//-----------------------------------------------------------------------------
/*

Nordic NRF52DK Board Mounting Kit

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var baseX = 120.0
var baseY = 64.0
var baseThickness = 3.0

var pcbX = 102.0
var pcbY = 63.5

var pillarHeight = 15.0

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

// standoffs1 (all with screw holes)
func standoffs1() SDF3 {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	k := &StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}

	// from the board gerbers
	positions := V3Set{
		{550.0 * Mil, 300.0 * Mil, zOfs},
		{600.0 * Mil, 2200.0 * Mil, zOfs},
		{2600.0 * Mil, 1600.0 * Mil, zOfs},
		{2600.0 * Mil, 500.0 * Mil, zOfs},
		{3800.0 * Mil, 300.0 * Mil, zOfs},
	}

	return Standoffs3D(k, positions)
}

// standoffs2 (one with a support stub)
func standoffs2() SDF3 {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}
	positions0 := V3Set{
		{550.0 * Mil, 300.0 * Mil, zOfs},
		{2600.0 * Mil, 1600.0 * Mil, zOfs},
		{2600.0 * Mil, 500.0 * Mil, zOfs},
		{3800.0 * Mil, 300.0 * Mil, zOfs},
	}
	s0 := Standoffs3D(k, positions0)

	// standoffs with support stubs
	k.HoleDepth = -2.0
	positions1 := V3Set{
		{600.0 * Mil, 2200.0 * Mil, zOfs},
	}
	s1 := Standoffs3D(k, positions1)

	return Union3D(s0, s1)
}

//-----------------------------------------------------------------------------

func base() SDF3 {
	// base
	pp := &PanelParms{
		Size:         V2{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s0 := Panel2D(pp)

	// cutouts
	c1 := Box2D(V2{53.0, 35.0}, 3.0)
	c1 = Transform2D(c1, Translate2d(V2{-22.0, 1.00}))
	c2 := Box2D(V2{20.0, 40.0}, 3.0)
	c2 = Transform2D(c2, Translate2d(V2{37.0, 3.0}))

	// extrude the base
	s2 := Extrude3D(Difference2D(s0, Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := pcbY - (0.5 * baseY)
	s2 = Transform3D(s2, Translate3d(V3{xOfs, yOfs, 0}))

	// add the standoffs
	//s3 := standoffs1() // all pillars have screw holes
	s3 := standoffs2() // one pillar has a support stub
	s4 := Union3D(s2, s3)
	s4.(*UnionSDF3).SetMin(PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(ScaleUniform3D(base(), shrink), 300, "nrf52dk.stl")
}

//-----------------------------------------------------------------------------
