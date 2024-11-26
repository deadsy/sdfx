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
// pcb mount base for rs232/ttl serial converter

func serialConverter() (sdf.SDF3, error) {

	pcb := v3.Vec{21.5, 40.0, 1.5}

	const margin0 = 0.5
	const baseY0 = 3.0
	const baseY1 = 15.0

	baseX0 := (0.5 * pcb.Y) - margin0
	const baseX1 = 5.0
	baseX := baseX0 + baseX1
	baseZ := pcb.X + 6.0

	// body profile
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(baseX, 0).Rel()
	p.Add(0, baseY0+baseY1).Rel()
	p.Add(-baseX1, 0).Rel()
	p.Add(baseX0, baseY0)
	p.Add(-baseX0, 0).Rel()
	p.Add(0, 0)
	s2d, _ := sdf.Polygon2D(p.Vertices())
	base0 := sdf.Extrude3D(s2d, baseZ)

	// indent for pcb board
	const margin1 = 0.3
	boardY := (baseY0 + baseY1) * 0.5

	p = sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(0.5*pcb.Y+margin1, 0).Rel()
	p.Add(baseX0, 2.0*pcb.Z)
	p.Add(-baseX0, 0).Rel()
	p.Add(0, 0)
	s2d, _ = sdf.Polygon2D(p.Vertices())
	pcbIndent := sdf.Extrude3D(s2d, pcb.X+margin1)
	pcbIndent = sdf.Transform3D(pcbIndent, sdf.Translate3d(v3.Vec{0, boardY, 0}))
	base0 = sdf.Difference3D(base0, pcbIndent)

	// base mounting hole
	hole, _ := sdf.Cylinder3D(3*baseY0, baseHoleDiameter*0.5, 0)
	hole = sdf.Transform3D(hole, sdf.RotateX(sdf.DtoR(90)))
	hole = sdf.Transform3D(hole, sdf.Translate3d(v3.Vec{baseX0 * 0.7, 0, 0}))
	base0 = sdf.Difference3D(base0, hole)

	base1 := sdf.Transform3D(base0, sdf.MirrorYZ())
	base := sdf.Union3D(base0, base1)

	return base, nil
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
}

//-----------------------------------------------------------------------------
