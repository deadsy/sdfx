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

const panelThickness = 3 // mm

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

// testHoles returns a panel with various holes for test fitting.
func testHoles() (sdf.SDF3, error) {

	const xInc = 15
	const yInc = 15
	const rInc = 0.1

	const nX = 5
	const nY = 8

	xOfs := 0.0
	yOfs := 0.0
	r := 1.5

	s := make([]sdf.SDF2, nX*nY)
	i := 0

	for j := 0; j < nY; j++ {
		for k := 0; k < nX; k++ {
			c, _ := sdf.Circle2D(r)
			s[i] = sdf.Transform2D(c, sdf.Translate2d(sdf.V2{xOfs, yOfs}))
			i += 1
			r += rInc
			xOfs += xInc
		}
		xOfs = 0.0
		yOfs += yInc
	}

	h := sdf.Union2D(s...)
	xOfs = -float64(nX-1) * xInc * 0.5
	yOfs = -float64(nY-1) * yInc * 0.5
	h = sdf.Transform2D(h, sdf.Translate2d(sdf.V2{xOfs, yOfs}))
	holes := sdf.Extrude3D(h, panelThickness)

	k := obj.PanelParms{
		Size:         sdf.V2{(nX + 1) * xInc, (nY + 1) * yInc},
		CornerRadius: xInc * 0.2,
	}

	p, err := obj.Panel2D(&k)
	if err != nil {
		return nil, err
	}
	panel := sdf.Extrude3D(p, panelThickness)

	return sdf.Difference3D(panel, holes), nil

}

//-----------------------------------------------------------------------------

// arPanel returns the panel for an attack/release module.
func arPanel() (sdf.SDF3, error) {

	// 3u x 12hp panel
	ep, err := obj.EuroRackPanel(3, 12, 3)
	if err != nil {
		return nil, err
	}
	s := sdf.Extrude3D(ep, panelThickness)

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
	cv, _ := sdf.Circle2D(3)
	cv0 := sdf.Transform2D(cv, sdf.Translate2d(sdf.V2{-20, -45}))
	cv1 := sdf.Transform2D(cv, sdf.Translate2d(sdf.V2{20, -45}))

	// LED
	led, _ := sdf.Circle2D(3.25)
	led = sdf.Transform2D(led, sdf.Translate2d(sdf.V2{0, -45}))

	// attack/release pots
	pot, _ := sdf.Circle2D(4)
	pot0 := sdf.Transform2D(pot, sdf.Translate2d(sdf.V2{-15, 25}))
	pot1 := sdf.Transform2D(pot, sdf.Translate2d(sdf.V2{15, 25}))

	// spdt switch
	spdt, _ := sdf.Circle2D(2.5)
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

	p1, err := testHoles()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(p1, shrink), 300, "holes.stl")

}

//-----------------------------------------------------------------------------
