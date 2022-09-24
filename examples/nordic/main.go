//-----------------------------------------------------------------------------
/*

Nordic nRF52x Development Board Mounting Kits

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
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

func nRF52dkStandoffs() (sdf.SDF3, error) {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}
	positions0 := v3.VecSet{
		{550.0 * sdf.Mil, 300.0 * sdf.Mil, zOfs},
		{2600.0 * sdf.Mil, 1600.0 * sdf.Mil, zOfs},
		{2600.0 * sdf.Mil, 500.0 * sdf.Mil, zOfs},
		{3800.0 * sdf.Mil, 300.0 * sdf.Mil, zOfs},
	}
	s, _ := obj.Standoff3D(k)
	s0 := sdf.Multi3D(s, positions0)

	// standoffs with support stubs
	k.HoleDepth = -2.0
	positions1 := v3.VecSet{
		{600.0 * sdf.Mil, 2200.0 * sdf.Mil, zOfs},
	}
	s, _ = obj.Standoff3D(k)
	s1 := sdf.Multi3D(s, positions1)

	return sdf.Union3D(s0, s1), nil
}

func nRF52dk() (sdf.SDF3, error) {

	baseX := 120.0
	baseY := 64.0
	pcbX := 102.0
	pcbY := 63.5

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	// cutouts
	c1 := sdf.Box2D(v2.Vec{53.0, 35.0}, 3.0)
	c1 = sdf.Transform2D(c1, sdf.Translate2d(v2.Vec{-22.0, 1.00}))
	c2 := sdf.Box2D(v2.Vec{20.0, 40.0}, 3.0)
	c2 = sdf.Transform2D(c2, sdf.Translate2d(v2.Vec{37.0, 3.0}))

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, sdf.Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := pcbY - (0.5 * baseY)
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// add the standoffs
	s3, err := nRF52dkStandoffs()
	if err != nil {
		return nil, err
	}
	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
}

//-----------------------------------------------------------------------------
// nRF52833DK
// https://www.nordicsemi.com/Software-and-tools/Development-Kits/nRF52833-DK

func nRF52833dkStandoffs() (sdf.SDF3, error) {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}
	positions0 := v3.VecSet{
		{550.0 * sdf.Mil, 300.0 * sdf.Mil, zOfs},
		{2600.0 * sdf.Mil, 500.0 * sdf.Mil, zOfs},
		{2600.0 * sdf.Mil, 1600.0 * sdf.Mil, zOfs},
		{5050.0 * sdf.Mil, 1825.0 * sdf.Mil, zOfs},
	}
	s, err := obj.Standoff3D(k)
	if err != nil {
		return nil, err
	}

	s0 := sdf.Multi3D(s, positions0)

	// standoffs with support stubs
	k.HoleDepth = -2.0
	positions1 := v3.VecSet{
		{600.0 * sdf.Mil, 2200.0 * sdf.Mil, zOfs},
		{3550.0 * sdf.Mil, 2200.0 * sdf.Mil, zOfs},
		{3800.0 * sdf.Mil, 300.0 * sdf.Mil, zOfs},
	}
	s, err = obj.Standoff3D(k)
	if err != nil {
		return nil, err
	}

	s1 := sdf.Multi3D(s, positions1)

	return sdf.Union3D(s0, s1), nil
}

func nRF52833dk() (sdf.SDF3, error) {

	baseX := 154.0
	baseY := 64.0
	pcbX := 136.53
	pcbY := 63.50

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}
	// cutouts
	c1 := sdf.Box2D(v2.Vec{53.0, 35.0}, 3.0)
	c1 = sdf.Transform2D(c1, sdf.Translate2d(v2.Vec{-40.0, 0}))
	c2 := sdf.Box2D(v2.Vec{40.0, 35.0}, 3.0)
	c2 = sdf.Transform2D(c2, sdf.Translate2d(v2.Vec{32.0, 0}))

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, sdf.Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := pcbY - (0.5 * baseY)
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// add the standoffs
	s3, err := nRF52833dkStandoffs()
	if err != nil {
		return nil, err
	}

	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
}

//-----------------------------------------------------------------------------

func main() {

	nRF52dk, err := nRF52dk()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(nRF52dk, shrink), 300, "nrf52dk.stl")

	nRF52833dk, err := nRF52833dk()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(nRF52833dk, shrink), 300, "nrf52833dk.stl")
}

//-----------------------------------------------------------------------------
