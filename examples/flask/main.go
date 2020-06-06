//-----------------------------------------------------------------------------
/*

Modular Casting Flask

Design credit to:
Olfoundryman: https://youtu.be/cX2u6S5qV3Q
smallcnclathes: http://www.benchtopcnc.com.au/downloads/

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const wallThickness = 5.0    // flask wall thickness (mm)
const padThickness = 5.0     // ping lugs pad thickness (mm)
const padWidth = 60.0        // pin lugs pad width (mm)
const padDraft = 30.0        // pad draft angle (degrees)
const cornerThickness = 7.0  // corner mounting lug thickness (mm)
const cornerLength = 24.0    // corner mounting lug length (mm)
const keyDepth = 4.0         // sand key depth (mm)
const keyDraft = 60.0        // sand key draft angle (degrees)
const keyRatio = 0.85        // sand key height / height of flask
const sideDraft = 3.0        // pattern side draft angle (degrees)
const lugBaseThickness = 2.0 // pin lugs base thickness (mm)
const lugBaseDraft = 15.0    // pin lugs base draft (defgrees)
const lugHeight = 26.0       // pin lug height (mm)
const lugThickness = 13.0    // pin lug thickness (mm)
const lugDraft = 5.0         // pin lug draft angle (degrees)
const lugOffset = 1.5        // pin lug base to pin offset (mm)
const alignRadius = 0.75     // alignment hole radius (mm)

// derived
const lugBaseWidth = padWidth * 0.95

//-----------------------------------------------------------------------------

func alignmentHoles() sdf.SDF3 {
	w := lugBaseWidth
	h := (lugBaseThickness + padThickness + wallThickness + cornerLength) * 2.0
	r := alignRadius
	xofs := w * 0.8 * 0.5
	return sdf.MultiCylinder3D(h, r, sdf.V2Set{{xofs, 0}, {-xofs, 0}})
}

// pinLug returns a single pin lug.
func pinLug(w float64) sdf.SDF3 {
	k := sdf.TruncRectPyramidParms{
		Size:        sdf.V3{w, lugThickness, lugHeight},
		Draft:       sdf.DtoR(lugDraft),
		BaseRadius:  lugThickness * 0.5,
		RoundRadius: lugThickness * 0.1,
	}
	return sdf.TruncRectPyramid3D(&k)
}

func pinLugs() sdf.SDF3 {
	// build the base
	w := lugBaseWidth
	r := lugThickness*0.5 + lugOffset
	k := sdf.TruncRectPyramidParms{
		Size:        sdf.V3{w, w, lugBaseThickness},
		Draft:       sdf.DtoR(lugBaseDraft),
		BaseRadius:  r,
		RoundRadius: lugBaseThickness * 0.25,
	}
	base := sdf.TruncRectPyramid3D(&k)

	// build the pin lugs
	pinWidth := w - 2.0*lugOffset
	pin := pinLug(pinWidth)
	yofs := 0.5 * (pinWidth - lugThickness)
	pin0 := sdf.Transform3D(pin, sdf.Translate3d(sdf.V3{0, yofs, lugBaseThickness}))
	pin1 := sdf.Transform3D(pin, sdf.Translate3d(sdf.V3{0, -yofs, lugBaseThickness}))

	s := sdf.Union3D(base, pin0, pin1)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(lugBaseThickness * 0.75))

	return sdf.Difference3D(s, alignmentHoles())
}

//-----------------------------------------------------------------------------

// keyProfile returns the 2d profile for the sand key.
func keyProfile(height float64) sdf.SDF2 {

	w0 := (height * 0.5) * keyRatio
	w1 := w0 - keyDepth*math.Tan(sdf.DtoR(keyDraft))

	r0 := keyDepth * 0.5

	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(w0, 0)
	p.Add(w1, keyDepth).Smooth(r0, 4)
	p.Add(0, keyDepth)

	s0 := sdf.Polygon2D(p.Vertices())
	s1 := sdf.Transform2D(s0, sdf.MirrorY())
	return sdf.Union2D(s0, s1)
}

//-----------------------------------------------------------------------------

// sideDraftProfile returns the 2d profile for the side draft of the flask pattern.
func sideDraftProfile(height float64) sdf.SDF2 {

	h0 := keyDepth + wallThickness + cornerLength
	w0 := height * 0.5
	w1 := w0 + w0
	w2 := w0 - h0*math.Tan(sdf.DtoR(sideDraft))

	p := sdf.NewPolygon()
	p.Add(w0, 0)
	p.Add(w1, 0)
	p.Add(w1, h0)
	p.Add(w2, h0)

	s0 := sdf.Polygon2D(p.Vertices())
	s1 := sdf.Transform2D(s0, sdf.MirrorY())
	return sdf.Union2D(s0, s1)
}

//-----------------------------------------------------------------------------

// flaskSideProfile returns a half 2D extrusion profile for the flask.
func flaskSideProfile(width float64) sdf.SDF2 {

	theta45 := sdf.DtoR(45)
	theta135 := sdf.DtoR(135)
	theta225 := sdf.DtoR(225)

	w0 := width * 0.5
	w1 := padWidth * 0.5
	w2 := w1 + padThickness*math.Tan(sdf.DtoR(padDraft))

	l0 := cornerLength + keyDepth/math.Sin(theta45)
	l1 := cornerLength - wallThickness/math.Sin(theta45) + cornerThickness*math.Tan(theta45)

	h0 := keyDepth + wallThickness
	h1 := keyDepth + wallThickness + padThickness

	r0 := cornerThickness * 0.25
	r1 := cornerThickness
	r2 := padThickness * 0.4

	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(w0, 0)
	p.Add(l0, theta45).Polar().Rel().Smooth(r0, 4)
	p.Add(cornerThickness, theta135).Polar().Rel().Smooth(r0, 4)
	p.Add(l1, theta225).Polar().Rel().Smooth(r1, 4)
	p.Add(w2, h0).Smooth(r2, 3)
	p.Add(w1, h1).Smooth(r2, 3)
	p.Add(0, h1)

	return sdf.Polygon2D(p.Vertices())
}

func flaskSide(width, height float64) sdf.SDF3 {

	// create the flask
	side0 := sdf.Extrude3D(flaskSideProfile(width), height)
	side1 := sdf.Transform3D(side0, sdf.MirrorYZ())
	flaskBody := sdf.Union3D(side0, side1)

	w := width + 2.0*cornerLength

	// internal sand key
	sandKey := sdf.Extrude3D(keyProfile(height), w)
	sandKey = sdf.Transform3D(sandKey, sdf.RotateY(sdf.DtoR(90)))
	// side draft
	sideDraft := sdf.Extrude3D(sideDraftProfile(height), w)
	sideDraft = sdf.Transform3D(sideDraft, sdf.RotateY(sdf.DtoR(90)))
	// alignment holes
	holes := alignmentHoles()
	holes = sdf.Transform3D(holes, sdf.RotateX(sdf.DtoR(90)))

	return sdf.Difference3D(flaskBody, sdf.Union3D(sandKey, sideDraft, holes))
}

//-----------------------------------------------------------------------------

func main() {
	widths := []float64{150, 200, 250, 300}
	height := 95.0
	for _, w := range widths {
		name := fmt.Sprintf("flask_%d.stl", int(w))
		sdf.RenderSTL(sdf.ScaleUniform3D(flaskSide(w, height), shrink), 300, name)
	}
	sdf.RenderSTL(sdf.ScaleUniform3D(pinLugs(), shrink), 120, "pins.stl")
}

//-----------------------------------------------------------------------------
