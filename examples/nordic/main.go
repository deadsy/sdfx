//-----------------------------------------------------------------------------
/*

Nordic nRF52x Development Board Mounting Kits

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

var baseThickness = 3.0
var pillarHeight = 15.0

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------
// nRF52DK
// https://www.nordicsemi.com/Software-and-tools/Development-Kits/nRF52-DK

func nRF52dkStandoffs() sdf.SDF3 {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}
	positions0 := sdf.V3Set{
		{550.0 * sdf.Mil, 300.0 * sdf.Mil, zOfs},
		{2600.0 * sdf.Mil, 1600.0 * sdf.Mil, zOfs},
		{2600.0 * sdf.Mil, 500.0 * sdf.Mil, zOfs},
		{3800.0 * sdf.Mil, 300.0 * sdf.Mil, zOfs},
	}
	s0 := sdf.Multi3D(obj.Standoff3D(k), positions0)

	// standoffs with support stubs
	k.HoleDepth = -2.0
	positions1 := sdf.V3Set{
		{600.0 * sdf.Mil, 2200.0 * sdf.Mil, zOfs},
	}
	s1 := sdf.Multi3D(obj.Standoff3D(k), positions1)

	return sdf.Union3D(s0, s1)
}

func nRF52dk() sdf.SDF3 {

	baseX := 120.0
	baseY := 64.0
	pcbX := 102.0
	pcbY := 63.5

	// base
	pp := &obj.PanelParms{
		Size:         sdf.V2{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s0 := obj.Panel2D(pp)

	// cutouts
	c1 := sdf.Box2D(sdf.V2{53.0, 35.0}, 3.0)
	c1 = sdf.Transform2D(c1, sdf.Translate2d(sdf.V2{-22.0, 1.00}))
	c2 := sdf.Box2D(sdf.V2{20.0, 40.0}, 3.0)
	c2 = sdf.Transform2D(c2, sdf.Translate2d(sdf.V2{37.0, 3.0}))

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, sdf.Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := pcbY - (0.5 * baseY)
	s2 = sdf.Transform3D(s2, sdf.Translate3d(sdf.V3{xOfs, yOfs, 0}))

	// add the standoffs
	s3 := nRF52dkStandoffs()
	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------
// nRF52833DK
// https://www.nordicsemi.com/Software-and-tools/Development-Kits/nRF52833-DK

func nRF52833dkStandoffs() sdf.SDF3 {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}
	positions0 := sdf.V3Set{
		{550.0 * sdf.Mil, 300.0 * sdf.Mil, zOfs},
		{2600.0 * sdf.Mil, 500.0 * sdf.Mil, zOfs},
		{2600.0 * sdf.Mil, 1600.0 * sdf.Mil, zOfs},
		{5050.0 * sdf.Mil, 1825.0 * sdf.Mil, zOfs},
	}
	s0 := sdf.Multi3D(obj.Standoff3D(k), positions0)

	// standoffs with support stubs
	k.HoleDepth = -2.0
	positions1 := sdf.V3Set{
		{600.0 * sdf.Mil, 2200.0 * sdf.Mil, zOfs},
		{3550.0 * sdf.Mil, 2200.0 * sdf.Mil, zOfs},
		{3800.0 * sdf.Mil, 300.0 * sdf.Mil, zOfs},
	}
	s1 := sdf.Multi3D(obj.Standoff3D(k), positions1)

	return sdf.Union3D(s0, s1)
}

func nRF52833dk() sdf.SDF3 {

	baseX := 154.0
	baseY := 64.0
	pcbX := 136.53
	pcbY := 63.50

	// base
	pp := &obj.PanelParms{
		Size:         sdf.V2{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s0 := obj.Panel2D(pp)

	// cutouts
	c1 := sdf.Box2D(sdf.V2{53.0, 35.0}, 3.0)
	c1 = sdf.Transform2D(c1, sdf.Translate2d(sdf.V2{-40.0, 0}))
	c2 := sdf.Box2D(sdf.V2{40.0, 35.0}, 3.0)
	c2 = sdf.Transform2D(c2, sdf.Translate2d(sdf.V2{32.0, 0}))

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, sdf.Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := pcbY - (0.5 * baseY)
	s2 = sdf.Transform3D(s2, sdf.Translate3d(sdf.V3{xOfs, yOfs, 0}))

	// add the standoffs
	s3 := nRF52833dkStandoffs()
	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------

func main() {
	sdf.RenderSTL(sdf.ScaleUniform3D(nRF52dk(), shrink), 300, "nrf52dk.stl")
	sdf.RenderSTL(sdf.ScaleUniform3D(nRF52833dk(), shrink), 300, "nrf52833dk.stl")
}

//-----------------------------------------------------------------------------
