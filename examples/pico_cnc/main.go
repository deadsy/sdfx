//-----------------------------------------------------------------------------
/*

Pico CNC Board Mounting Kit

https://github.com/phil-barrett/PicoCNC

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

const baseHoleDiameter = 3.5

//-----------------------------------------------------------------------------
// keypad panel

func keypadPanel() (sdf.SDF3, error) {

	const panelThickness = 5.5
	const panelX = 75
	const panelYa = 25
	const panelYb = 45
	const panelY = 2 * (panelYa + panelYb)

	k := &obj.PanelParms{
		Size:         v2.Vec{panelX, panelY},
		CornerRadius: 4,
		HoleDiameter: baseHoleDiameter,
		HoleMargin:   [4]float64{7, 7, 7, 7},
		HolePattern:  [4]string{"x", "xx", "x", "xx"},
		Thickness:    panelThickness,
	}

	// key hole
	const holeRadius = (22.0 + 1.5) * 0.5

	hole0, err := sdf.Cylinder3D(panelThickness, holeRadius, 0)
	if err != nil {
		return nil, err
	}
	hole1 := sdf.Transform3D(hole0, sdf.Translate3d(v3.Vec{0, panelYb, 0}))
	hole2 := sdf.Transform3D(hole0, sdf.Translate3d(v3.Vec{0, -panelYb, 0}))

	panel, err := obj.Panel3D(k)
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(panel, sdf.Union3D(hole0, hole1, hole2)), nil
}

//-----------------------------------------------------------------------------
// pcb mount base for rs232/ttl serial converter

func serialConverter() (sdf.SDF3, error) {

	// v3.Vec{0, 0.4, 0.1} // too tight

	pcb := v3.Vec{21.5, 40.0, 1.5}.Add(v3.Vec{0, 0.8, 0.4})

	wallThickness := 5.0
	innerBox := v3.Vec{pcb.X, pcb.Y - 3.0, 15}
	outerBox := innerBox.Add(v3.Vec{wallThickness, 2.0 * wallThickness, wallThickness})

	outer, _ := sdf.Box3D(outerBox, 0.5*wallThickness)
	inner, _ := sdf.Box3D(innerBox, 0)

	// body
	s := sdf.Difference3D(outer, inner)
	s = sdf.Cut3D(s, v3.Vec{0.5 * innerBox.X, 0, 0}, v3.Vec{-1, 0, 0})
	s = sdf.Cut3D(s, v3.Vec{0, 0, 0.5 * innerBox.Z}, v3.Vec{0, 0, -1})

	// base mounting hole
	hole0, _ := sdf.Cylinder3D(10*wallThickness, baseHoleDiameter*0.5, 0)
	hole0 = sdf.Transform3D(hole0, sdf.Translate3d(v3.Vec{0, 0.35 * innerBox.Y, 0}))
	hole1 := sdf.Transform3D(hole0, sdf.MirrorXZ())
	holes := sdf.Union3D(hole0, hole1)

	// pcb
	board, _ := sdf.Box3D(pcb, 0)

	s = sdf.Difference3D(s, holes)
	s = sdf.Difference3D(s, board)

	return s, nil
}

//-----------------------------------------------------------------------------
// pcb mount base for pico cnc

const baseThickness = 3.0
const pcbX = 92.0
const pcbY = 94.5
const pcbHoleMargin = 3.5

func picoCncStandoffs() (sdf.SDF3, error) {

	const pillarHeight = 15.0
	const zOfs = 0.5 * (pillarHeight + baseThickness)
	const dx = pcbX - (2.0 * pcbHoleMargin)
	const dy = pcbY - (2.0 * pcbHoleMargin)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}
	positions0 := v3.VecSet{
		{0, 0, zOfs},
		{dx, 0, zOfs},
		{0, dy, zOfs},
		{dx, dy, zOfs},
	}
	s, err := obj.Standoff3D(k)
	if err != nil {
		return nil, err
	}
	return sdf.Multi3D(s, positions0), nil
}

func picoCnc() (sdf.SDF3, error) {

	const holeMargin = 3.0
	const baseX = pcbX + (2.0 * holeMargin)
	const baseY = pcbY + (2.0 * holeMargin)
	const cutoutMargin = 12.0
	const cutoutX = baseX - (2.0 * cutoutMargin)
	const cutoutY = baseY - (2.0 * cutoutMargin)

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: baseHoleDiameter,
		HoleMargin:   [4]float64{6.0, 6.0, 6.0, 6.0},
		HolePattern:  [4]string{".x...x", ".x...x", ".x...x", ".x...x"},
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	// cutouts
	c0 := sdf.Box2D(v2.Vec{cutoutX, cutoutY}, 3.0)

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, c0), baseThickness)

	const xOfs = (0.5 * baseX) - holeMargin - pcbHoleMargin
	const yOfs = (0.5 * baseY) - holeMargin - pcbHoleMargin
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// add the standoffs
	s3, err := picoCncStandoffs()
	if err != nil {
		return nil, err
	}

	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := picoCnc()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "pico_cnc.stl", render.NewMarchingCubesOctree(300))

	s, err = serialConverter()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "serial.stl", render.NewMarchingCubesOctree(300))

	s, err = keypadPanel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "keypad_panel.stl", render.NewMarchingCubesOctree(300))

	s, err = penHolder()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "pen_holder.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
