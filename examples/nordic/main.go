//-----------------------------------------------------------------------------
/*

Nordic nRF52x Development Board Mounting Kits

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var baseThickness = 3.0
var pillarHeight = 15.0

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------
// nRF52DK
// https://www.nordicsemi.com/Software-and-tools/Development-Kits/nRF52-DK

func nRF52dkStandoffs() SDF3 {

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
	s0 := Multi3D(Standoff3D(k), positions0)

	// standoffs with support stubs
	k.HoleDepth = -2.0
	positions1 := V3Set{
		{600.0 * Mil, 2200.0 * Mil, zOfs},
	}
	s1 := Multi3D(Standoff3D(k), positions1)

	return Union3D(s0, s1)
}

func nRF52dk() SDF3 {

	baseX := 120.0
	baseY := 64.0
	pcbX := 102.0
	pcbY := 63.5

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
	s3 := nRF52dkStandoffs()
	s4 := Union3D(s2, s3)
	s4.(*UnionSDF3).SetMin(PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------
// nRF52833DK
// https://www.nordicsemi.com/Software-and-tools/Development-Kits/nRF52833-DK

func nRF52833dkStandoffs() SDF3 {

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
		{2600.0 * Mil, 500.0 * Mil, zOfs},
		{2600.0 * Mil, 1600.0 * Mil, zOfs},
		{5050.0 * Mil, 1825.0 * Mil, zOfs},
	}
	s0 := Multi3D(Standoff3D(k), positions0)

	// standoffs with support stubs
	k.HoleDepth = -2.0
	positions1 := V3Set{
		{600.0 * Mil, 2200.0 * Mil, zOfs},
		{3550.0 * Mil, 2200.0 * Mil, zOfs},
		{3800.0 * Mil, 300.0 * Mil, zOfs},
	}
	s1 := Multi3D(Standoff3D(k), positions1)

	return Union3D(s0, s1)
}

func nRF52833dk() SDF3 {

	baseX := 154.0
	baseY := 64.0
	pcbX := 136.53
	pcbY := 63.50

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
	c1 = Transform2D(c1, Translate2d(V2{-40.0, 0}))
	c2 := Box2D(V2{40.0, 35.0}, 3.0)
	c2 = Transform2D(c2, Translate2d(V2{32.0, 0}))

	// extrude the base
	s2 := Extrude3D(Difference2D(s0, Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := pcbY - (0.5 * baseY)
	s2 = Transform3D(s2, Translate3d(V3{xOfs, yOfs, 0}))

	// add the standoffs
	s3 := nRF52833dkStandoffs()
	s4 := Union3D(s2, s3)
	s4.(*UnionSDF3).SetMin(PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(ScaleUniform3D(nRF52dk(), shrink), 300, "nrf52dk.stl")
	RenderSTL(ScaleUniform3D(nRF52833dk(), shrink), 300, "nrf52833dk.stl")
}

//-----------------------------------------------------------------------------
