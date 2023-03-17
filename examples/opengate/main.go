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

var baseThickness = 3.0
var pillarHeight = 15.0

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func standoffs() (sdf.SDF3, error) {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}
	positions0 := v3.VecSet{
		{3.5, 3.5, zOfs},
		{3.5 + 116.0, 3.5, zOfs},
		{3.5 + 116.0, 3.5 + 61.0, zOfs},
		{3.5, 3.5 + 61.0, zOfs},
	}
	s, err := obj.Standoff3D(k)
	if err != nil {
		return nil, err
	}

	return sdf.Multi3D(s, positions0), nil
}

func mainBoard() (sdf.SDF3, error) {

	baseX := 180.0
	baseY := 80.0
	pcbX := 123.0
	pcbY := 68.0

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

	// extrude the base
	s2 := sdf.Extrude3D(s0, baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := pcbY - (0.5 * baseY)
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// add the standoffs
	s3, err := standoffs()
	if err != nil {
		return nil, err
	}

	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
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
