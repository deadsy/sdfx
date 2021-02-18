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

const panelThickness = 2.5 // mm

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
// panel holes and/or indents for mounted components

// pot0 return the panel hole/indent for a potentiometer
func pot0() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  9.4,
		Thickness: panelThickness,
		Indent:    sdf.V3{2, 4, 2},
		Offset:    11.0,
		//Orientation: sdf.DtoR(0),
	}
	return obj.PanelHole3D(&k)
}

// pot1 return the panel hole/indent for a potentiometer
func pot1() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  7.2,
		Thickness: panelThickness,
		Indent:    sdf.V3{2, 2, 1.5},
		Offset:    7.0,
		//Orientation: sdf.DtoR(0),
	}
	return obj.PanelHole3D(&k)
}

// spdt return the panel hole/indent for a spdt switch
func spdt() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  6.2,
		Thickness: panelThickness,
		Indent:    sdf.V3{2, 2, 1.5},
		Offset:    5.4,
		//Orientation: sdf.DtoR(0),
	}
	return obj.PanelHole3D(&k)
}

// led returns the panel hole for an led bezel
func led() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  7.0,
		Thickness: panelThickness,
	}
	return obj.PanelHole3D(&k)
}

// jack35 returns the panel hole/indent for a 3.5 mm audio jack
func jack35() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  6.4,
		Thickness: panelThickness,
		Indent:    sdf.V3{2, 2, 1.5},
		Offset:    4.9,
		//Orientation: sdf.DtoR(0),
	}
	return obj.PanelHole3D(&k)
}

//-----------------------------------------------------------------------------

// arPanel returns the panel for an attack/release module.
func arPanel() (sdf.SDF3, error) {

	// 3u x 12hp panel
	k := obj.EuroRackParms{
		U:            3,
		HP:           12,
		CornerRadius: 3,
		HoleDiameter: 3.6,
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
	pb, err := sdf.Box3D(sdf.V3{13.2, 10.8, panelThickness}, 0)
	if err != nil {
		return nil, err
	}
	pb = sdf.Transform3D(pb, sdf.Translate3d(sdf.V3{0, 0, 0}))

	// cv input/output
	cv, err := jack35()
	if err != nil {
		return nil, err
	}
	cv0 := sdf.Transform3D(cv, sdf.Translate3d(sdf.V3{-20, -45, 0}))
	cv1 := sdf.Transform3D(cv, sdf.Translate3d(sdf.V3{20, -45, 0}))

	// LED
	led, err := led()
	if err != nil {
		return nil, err
	}
	led = sdf.Transform3D(led, sdf.Translate3d(sdf.V3{0, -45, 0}))

	// attack/release pots
	pot, err := pot0()
	if err != nil {
		return nil, err
	}
	pot0 := sdf.Transform3D(pot, sdf.Translate3d(sdf.V3{-15, 25, 0}))
	pot1 := sdf.Transform3D(pot, sdf.Translate3d(sdf.V3{15, 25, 0}))

	// spdt switch
	spdt, err := spdt()
	if err != nil {
		return nil, err
	}
	spdt = sdf.Transform3D(spdt, sdf.Translate3d(sdf.V3{0, -22, 0}))

	return sdf.Difference3D(s, sdf.Union3D(pb, cv0, cv1, led, pot0, pot1, spdt)), nil
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
