//-----------------------------------------------------------------------------
/*

Axoloti Board Mounting Kit

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const frontPanelThickness = 3.0
const frontPanelLength = 170.0
const frontPanelHeight = 50.0
const frontPanelYOffset = 15.0

const baseWidth = 50.0
const baseLength = 170.0
const baseThickness = 3.0

const baseFootWidth = 10.0
const baseFootCornerRadius = 3.0

const pcbWidth = 50.0
const pcbLength = 160.0

const pillarHeight = 16.8

//-----------------------------------------------------------------------------

// multiple standoffs
func standoffs() (sdf.SDF3, error) {

	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// from the board mechanicals
	positions := v3.VecSet{
		{3.5, 10.0, zOfs},   // H1
		{3.5, 40.0, zOfs},   // H2
		{54.0, 40.0, zOfs},  // H3
		{156.5, 10.0, zOfs}, // H4
		//{54.0, 10.0, zOfs},  // H5
		{156.5, 40.0, zOfs}, // H6
		{44.0, 10.0, zOfs},  // H7
		{116.0, 10.0, zOfs}, // H8
	}

	s, err := obj.Standoff3D(k)
	if err != nil {
		return nil, err
	}
	return sdf.Multi3D(s, positions), nil
}

//-----------------------------------------------------------------------------

// base returns the base mount.
func base() (sdf.SDF3, error) {
	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseLength, baseWidth},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{7.0, 20.0, 7.0, 20.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	// cutout
	l := baseLength - (2.0 * baseFootWidth)
	w := 18.0
	s1 := sdf.Box2D(v2.Vec{l, w}, baseFootCornerRadius)
	yOfs := 0.5 * (baseWidth - pcbWidth)
	s1 = sdf.Transform2D(s1, sdf.Translate2d(v2.Vec{0, yOfs}))

	s2 := sdf.Extrude3D(sdf.Difference2D(s0, s1), baseThickness)
	xOfs := 0.5 * pcbLength
	yOfs = pcbWidth - (0.5 * baseWidth)
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// standoffs
	s3, err := standoffs()
	if err != nil {
		return nil, err
	}

	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
}

//-----------------------------------------------------------------------------
// front panel cutouts

type panelHole struct {
	center v2.Vec   // center of hole
	hole   sdf.SDF2 // 2d hole
}

// button positions
const pbX = 53.0

var pb0 = v2.Vec{pbX, 0.8}
var pb1 = v2.Vec{pbX + 5.334, 0.8}

// panelCutouts returns the 2D front panel cutouts
func panelCutouts() (sdf.SDF2, error) {

	sMidi, err := sdf.Circle2D(0.5 * 17.0)
	if err != nil {
		return nil, err
	}
	sJack0, err := sdf.Circle2D(0.5 * 11.5)
	if err != nil {
		return nil, err
	}
	sJack1, err := sdf.Circle2D(0.5 * 5.5)
	if err != nil {
		return nil, err
	}

	sLed := sdf.Box2D(v2.Vec{1.6, 1.6}, 0)

	k := obj.FingerButtonParms{
		Width:  4.0,
		Gap:    0.6,
		Length: 20.0,
	}
	fb, err := obj.FingerButton2D(&k)
	if err != nil {
		return nil, err
	}
	sButton := sdf.Transform2D(fb, sdf.Rotate2d(sdf.DtoR(-90)))

	jackX := 123.0
	midiX := 18.8
	ledX := 62.9

	holes := []panelHole{
		{v2.Vec{midiX, 10.2}, sMidi},                         // MIDI DIN Jack
		{v2.Vec{midiX + 20.32, 10.2}, sMidi},                 // MIDI DIN Jack
		{v2.Vec{jackX, 8.14}, sJack0},                        // 1/4" Stereo Jack
		{v2.Vec{jackX + 19.5, 8.14}, sJack0},                 // 1/4" Stereo Jack
		{v2.Vec{107.6, 2.3}, sJack1},                         // 3.5 mm Headphone Jack
		{v2.Vec{ledX, 0.5}, sLed},                            // LED
		{v2.Vec{ledX + 3.635, 0.5}, sLed},                    // LED
		{pb0, sButton},                                       // Push Button
		{pb1, sButton},                                       // Push Button
		{v2.Vec{84.1, 1.0}, sdf.Box2D(v2.Vec{16.0, 7.5}, 0)}, // micro SD card
		{v2.Vec{96.7, 1.0}, sdf.Box2D(v2.Vec{11.0, 7.5}, 0)}, // micro USB connector
		{v2.Vec{73.1, 7.1}, sdf.Box2D(v2.Vec{7.5, 15.0}, 0)}, // fullsize USB connector
	}

	s := make([]sdf.SDF2, len(holes))
	for i, k := range holes {
		s[i] = sdf.Transform2D(k.hole, sdf.Translate2d(k.center))
	}

	return sdf.Union2D(s...), nil
}

//-----------------------------------------------------------------------------

// frontPanel returns the front panel mount.
func frontPanel() (sdf.SDF3, error) {

	// overall panel
	pp := &obj.PanelParms{
		Size:         v2.Vec{frontPanelLength, frontPanelHeight},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	panel, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	xOfs := 0.5 * pcbLength
	yOfs := (0.5 * frontPanelHeight) - frontPanelYOffset
	panel = sdf.Transform2D(panel, sdf.Translate2d(v2.Vec{xOfs, yOfs}))

	// extrude to 3d
	panelCutouts, err := panelCutouts()
	if err != nil {
		return nil, err
	}
	fp := sdf.Extrude3D(sdf.Difference2D(panel, panelCutouts), frontPanelThickness)

	// Add buttons to the finger button
	bHeight := 4.0
	b, _ := sdf.Cylinder3D(bHeight, 1.4, 0)
	b0 := sdf.Transform3D(b, sdf.Translate3d(conv.V2ToV3(pb0, -0.5*bHeight)))
	b1 := sdf.Transform3D(b, sdf.Translate3d(conv.V2ToV3(pb1, -0.5*bHeight)))

	return sdf.Union3D(fp, b0, b1), nil
}

//-----------------------------------------------------------------------------

func main() {

	// front panel
	s0, err := frontPanel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	sx := sdf.Transform3D(s0, sdf.RotateY(sdf.DtoR(180.0)))
	render.RenderSTL(sdf.ScaleUniform3D(sx, shrink), 400, "panel.stl")

	// base
	s1, err := base()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(s1, shrink), 400, "base.stl")

	// both together
	s0 = sdf.Transform3D(s0, sdf.Translate3d(v3.Vec{0, 80, 0}))
	s3 := sdf.Union3D(s0, s1)
	render.RenderSTL(sdf.ScaleUniform3D(s3, shrink), 400, "panel_and_base.stl")
}

//-----------------------------------------------------------------------------
