//-----------------------------------------------------------------------------
/*

https://shop.sb-components.co.uk/products/raspberry-pi-pico-hat-expansion

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
// material shrinkage

const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const baseThickness = 3
const pillarHeight = 8

const pcbX = 116
const pcbY = 61

const baseX = pcbX + 30
const baseY = pcbY + 20

//-----------------------------------------------------------------------------

func standoffs() (sdf.SDF3, error) {

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      pillarHeight,
		HoleDiameter:   2.4, // #4 screw
	}

	s, err := obj.Standoff3D(k)
	if err != nil {
		return nil, err
	}

	positions0 := v3.VecSet{
		{0, 0, 0},
		{pcbX, 0, 0},
		{pcbX, pcbY, 0},
		{0, pcbY, 0},
	}
	s = sdf.Multi3D(s, positions0)

	xOfs := -0.5 * pcbX
	yOfs := -0.5 * pcbY
	zOfs := 0.5 * (pillarHeight + baseThickness)
	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{xOfs, yOfs, zOfs}))

	return s, nil
}

func mainBoard() (sdf.SDF3, error) {

	// base
	base := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s0, err := obj.Panel2D(base)
	if err != nil {
		return nil, err
	}

	// cutout
	cutout := &obj.PanelParms{
		Size:         v2.Vec{baseX - 40, baseY - 40},
		CornerRadius: 5.0,
	}
	s1, err := obj.Panel2D(cutout)
	if err != nil {
		return nil, err
	}

	s2 := sdf.Difference2D(s0, s1)

	// extrude the base
	s3 := sdf.Extrude3D(s2, baseThickness)

	// add the standoffs
	s4, err := standoffs()
	if err != nil {
		return nil, err
	}

	s5 := sdf.Union3D(s3, s4)
	s5.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s5, nil
}

//-----------------------------------------------------------------------------

func main() {

	s, err := mainBoard()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "main_board.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
