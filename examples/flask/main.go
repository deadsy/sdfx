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
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
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
const cornerLength = 30.0    // corner mounting lug length (mm)
const keyDepth = 4.0         // sand key depth (mm)
const keyDraft = 60.0        // sand key draft angle (degrees)
const keyRatio = 0.85        // sand key height / height of flask
const sideDraft = 3.0        // pattern side draft angle (degrees)
const lugBaseThickness = 3.0 // pin lugs base thickness (mm)
const lugBaseDraft = 15.0    // pin lugs base draft (defgrees)
const lugHeight = 28.0       // pin lug height (mm)
const lugThickness = 14.0    // pin lug thickness (mm)
const lugDraft = 5.0         // pin lug draft angle (degrees)
const lugOffset = 1.5        // pin lug base to pin offset (mm)
const holeRadius = 1.5       // alignment/pull hole radius (mm)

// derived
const lugBaseWidth = padWidth * 0.95

//-----------------------------------------------------------------------------

// alignmentHoles returns an SDF3 for the alignment holes between the flask and pin lugs pattern.
func alignmentHoles() (sdf.SDF3, error) {
	w := lugBaseWidth
	h := (lugBaseThickness + padThickness + wallThickness + cornerLength) * 2.0
	xofs := w * 0.8 * 0.5
	cylinder, err := sdf.Cylinder3D(h, holeRadius, 0)
	if err != nil {
		return nil, err
	}
	return sdf.Multi3D(cylinder, v3.VecSet{{xofs, 0, 0}, {-xofs, 0, 0}}), nil
}

// pinLug returns a single pin lug.
func pinLug(w float64) (sdf.SDF3, error) {
	// pin
	k := obj.TruncRectPyramidParms{
		Size:        v3.Vec{w, lugThickness, lugHeight},
		BaseAngle:   sdf.DtoR(90 - lugDraft),
		BaseRadius:  lugThickness * 0.5,
		RoundRadius: lugThickness * 0.1,
	}
	return obj.TruncRectPyramid3D(&k)
}

// pinLugs returns an SDF3 for the pin lugs pattern.
func pinLugs() (sdf.SDF3, error) {
	// build the base
	w := lugBaseWidth
	r := lugThickness*0.5 + lugOffset
	k := obj.TruncRectPyramidParms{
		Size:        v3.Vec{w, w, lugBaseThickness},
		BaseAngle:   sdf.DtoR(90 - lugBaseDraft),
		BaseRadius:  r,
		RoundRadius: lugBaseThickness * 0.25,
	}
	base, err := obj.TruncRectPyramid3D(&k)
	if err != nil {
		return nil, err
	}

	// build the pin lugs
	pinWidth := w - 2.0*lugOffset
	pin, err := pinLug(pinWidth)
	if err != nil {
		return nil, err
	}
	yofs := 0.5 * (pinWidth - lugThickness)
	pin0 := sdf.Transform3D(pin, sdf.Translate3d(v3.Vec{0, yofs, lugBaseThickness}))
	pin1 := sdf.Transform3D(pin, sdf.Translate3d(v3.Vec{0, -yofs, lugBaseThickness}))

	// combine the base and pins
	s := sdf.Union3D(base, pin0, pin1)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(lugBaseThickness * 0.75))

	// add alignment holes
	holes, err := alignmentHoles()
	if err != nil {
		return nil, err
	}
	return sdf.Difference3D(s, holes), nil
}

//-----------------------------------------------------------------------------

// sandKey returns an SDF3 for the internal sand key.
func sandKey(size v3.Vec) (sdf.SDF3, error) {
	theta := sdf.DtoR(90 - keyDraft)
	r := keyDepth / math.Tan(theta)
	k := obj.TruncRectPyramidParms{
		Size:        size,
		BaseAngle:   theta,
		BaseRadius:  r,
		RoundRadius: size.X * 0.5,
	}
	return obj.TruncRectPyramid3D(&k)
}

//-----------------------------------------------------------------------------

// oddSide returns an SDF3 for the odd sides at either end of the flask pattern.
func oddSide(height float64) (sdf.SDF3, error) {

	theta45 := sdf.DtoR(45)

	d := cornerLength * math.Cos(theta45)
	sx := 2.0*d + cornerThickness
	sy := height*1.1 + 2.0*d
	sz := d

	k := obj.TruncRectPyramidParms{
		Size:        v3.Vec{sx, sy, sz},
		BaseAngle:   theta45,
		BaseRadius:  0.5 * sx,
		RoundRadius: 0,
	}
	base, _ := obj.TruncRectPyramid3D(&k)

	// mounting/pull holes
	h := 3.0 * d
	yofs := (height*1.1 - cornerThickness) * 0.5
	hole, _ := sdf.Cylinder3D(h, holeRadius, 0)
	holes := sdf.Multi3D(hole, v3.VecSet{{0, yofs, 0}, {0, -yofs, 0}})

	// hook into internal sand key
	sx = 0.8 * sx
	sy = height * keyRatio * 0.99
	sz = keyDepth
	key, _ := sandKey(v3.Vec{sx, sy, sz})
	key = sdf.Transform3D(key, sdf.Translate3d(v3.Vec{0.5 * sx, 0, 0}))

	return sdf.Difference3D(sdf.Union3D(base, key), holes), nil
}

//-----------------------------------------------------------------------------

// sideDraftProfile returns the 2d profile for the side draft of the flask pattern.
func sideDraftProfile(height float64) (sdf.SDF2, error) {

	h0 := keyDepth + wallThickness + cornerLength
	w0 := height * 0.5
	w1 := w0 + w0
	w2 := w0 - h0*math.Tan(sdf.DtoR(sideDraft))

	p := sdf.NewPolygon()
	p.Add(w0, 0)
	p.Add(w1, 0)
	p.Add(w1, h0)
	p.Add(w2, h0)

	s0, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	s1 := sdf.Transform2D(s0, sdf.MirrorY())
	return sdf.Union2D(s0, s1), nil
}

//-----------------------------------------------------------------------------

// flaskSideProfile returns a half 2D extrusion profile for the flask.
func flaskSideProfile(width float64) (sdf.SDF2, error) {

	theta45 := sdf.DtoR(45)
	theta135 := sdf.DtoR(135)
	theta225 := sdf.DtoR(225)

	w0 := width * 0.5
	w1 := padWidth * 0.5
	w2 := w1 + padThickness*math.Tan(sdf.DtoR(padDraft))

	h0 := keyDepth + wallThickness
	h1 := keyDepth + wallThickness + padThickness

	l0 := cornerLength + cornerThickness - (keyDepth+wallThickness)/math.Sin(theta45)

	r0 := cornerThickness * 0.25
	r1 := cornerThickness
	r2 := padThickness * 0.4

	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(w0, 0)
	p.Add(cornerLength, theta45).Polar().Rel().Smooth(r0, 4)
	p.Add(cornerThickness, theta135).Polar().Rel().Smooth(r0, 4)
	p.Add(l0, theta225).Polar().Rel().Smooth(r1, 4)
	p.Add(w2, h0).Smooth(r2, 3)
	p.Add(w1, h1).Smooth(r2, 3)
	p.Add(0, h1)

	return sdf.Polygon2D(p.Vertices())
}

// pullHoles returns an SDF3 for the flask pull holes.
func pullHoles(width float64) (sdf.SDF3, error) {
	h := (wallThickness + keyDepth) * 2.0
	xofs := width * 0.9 * 0.5
	hole, _ := sdf.Cylinder3D(h, holeRadius, 0)
	return sdf.Multi3D(hole, v3.VecSet{{xofs, 0, 0}, {-xofs, 0, 0}}), nil
}

func flaskHalf(width, height float64) (sdf.SDF3, error) {
	s, err := flaskSideProfile(width)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, height), nil
}

// flaskSide returns an SDF3 for the flask side.
func flaskSide(width, height float64) (sdf.SDF3, error) {

	// create the flask
	side0, err := flaskHalf(width, height)
	if err != nil {
		return nil, err
	}
	side1 := sdf.Transform3D(side0, sdf.MirrorYZ())
	flaskBody := sdf.Union3D(side0, side1)

	w := width + 2.0*cornerLength

	// internal sand key
	key, err := sandKey(v3.Vec{w, height * keyRatio, keyDepth})
	if err != nil {
		return nil, err
	}
	key = sdf.Transform3D(key, sdf.RotateX(sdf.DtoR(-90)))

	// side draft
	s, err := sideDraftProfile(height)
	if err != nil {
		return nil, err
	}
	sideDraft := sdf.Extrude3D(s, w)
	sideDraft = sdf.Transform3D(sideDraft, sdf.RotateY(sdf.DtoR(90)))

	// alignment holes
	aHoles, err := alignmentHoles()
	if err != nil {
		return nil, err
	}
	aHoles = sdf.Transform3D(aHoles, sdf.RotateX(sdf.DtoR(90)))

	// pull holes
	pHoles, err := pullHoles(width)
	if err != nil {
		return nil, err
	}
	pHoles = sdf.Transform3D(pHoles, sdf.RotateX(sdf.DtoR(90)))

	return sdf.Difference3D(flaskBody, sdf.Union3D(key, sideDraft, aHoles, pHoles)), nil
}

//-----------------------------------------------------------------------------

func main() {
	widths := []float64{150, 200, 250, 300}
	height := 95.0
	for _, w := range widths {
		s, err := flaskSide(w, height)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		// rotate for the preferred print orientation
		s = sdf.Transform3D(s, sdf.RotateX(-sdf.DtoR(sideDraft)))
		name := fmt.Sprintf("flask_%d.stl", int(w))
		render.RenderSTL(sdf.ScaleUniform3D(s, shrink), 300, name)
	}

	pinLugs, err := pinLugs()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(pinLugs, shrink), 120, "pins.stl")

	oddSide, err := oddSide(height)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(oddSide, shrink), 300, "odd_side.stl")
}

//-----------------------------------------------------------------------------
