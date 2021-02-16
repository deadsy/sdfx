//-----------------------------------------------------------------------------
/*

Create Eurorack Module Panels

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func standoff(h float64) (sdf.SDF3, error) {
	// standoff with screw hole
	k := &obj.StandoffParms{
		PillarHeight:   h,
		PillarDiameter: 8,
		HoleDepth:      10,
		HoleDiameter:   2.4, // #4 screw
	}
	return obj.Standoff3D(k)
}

// halfBreadBoardStandoffs returns the standoffs for an adafruit 1/2 breadboard.
func halfBreadBoardStandoffs(h float64) (sdf.SDF3, error) {
	s, err := standoff(h)
	if err != nil {
		return nil, err
	}
	positions := sdf.V3Set{
		{0, -1450 * sdf.Mil, 0},
		{0, 1450 * sdf.Mil, 0},
	}
	return sdf.Multi3D(s, positions), nil
}

//-----------------------------------------------------------------------------

// arPanel returns the panel for an attack/release module.
func arPanel() (sdf.SDF3, error) {

	const panelThickness = 2.5 // mm

	// 3u x 12hp panel
	k := obj.EuroRackParms{
		U:            3,
		HP:           12,
		CornerRadius: 3,
		HoleDiameter: 0,
		Thickness:    panelThickness,
		Ridge:        true,
	}
	s, err := obj.EuroRackPanel3D(&k)
	if err != nil {
		return nil, err
	}

	// breadboard standoffs
	const standoffHeight = 25
	so, err := halfBreadBoardStandoffs(standoffHeight)
	if err != nil {
		return nil, err
	}
	so = sdf.Transform3D(so, sdf.Translate3d(sdf.V3{0, 3, (panelThickness + standoffHeight) * 0.5}))
	s = sdf.Union3D(s, so)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(2))

	// push button
	pb := sdf.Box2D(sdf.V2{13.2, 10.8}, 0)
	pb = sdf.Transform2D(pb, sdf.Translate2d(sdf.V2{0, 0}))

	// cv input/output
	cv, _ := sdf.Circle2D(3.1)
	cv0 := sdf.Transform2D(cv, sdf.Translate2d(sdf.V2{-20, -45}))
	cv1 := sdf.Transform2D(cv, sdf.Translate2d(sdf.V2{20, -45}))

	// LED
	led, _ := sdf.Circle2D(3.5)
	led = sdf.Transform2D(led, sdf.Translate2d(sdf.V2{0, -45}))

	// attack/release pots
	pot, _ := sdf.Circle2D(4.7)
	pot0 := sdf.Transform2D(pot, sdf.Translate2d(sdf.V2{-15, 25}))
	pot1 := sdf.Transform2D(pot, sdf.Translate2d(sdf.V2{15, 25}))

	// spdt switch
	spdt, _ := sdf.Circle2D(3.1)
	spdt = sdf.Transform2D(spdt, sdf.Translate2d(sdf.V2{0, -22}))

	cutouts := sdf.Extrude3D(sdf.Union2D(pb, cv0, cv1, led, pot0, pot1, spdt), panelThickness)

	return sdf.Difference3D(s, cutouts), nil
}

//-----------------------------------------------------------------------------

func main() {
	p0, err := arPanel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(p0, shrink), 300, "ar_panel.stl")
}

//-----------------------------------------------------------------------------
