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

const wallThickness = 5.0   // flask wall thickness (mm)
const padThickness = 5.0    // thickness of pad for pin lugs (mm)
const padWidth = 60.0       // width of pad for pin lugs (mm)
const padDraft = 30.0       // pad draft angle (degrees)
const cornerThickness = 7.0 // thickness of corner mount (mm)
const cornerLength = 24.0   // length of corner mount (mm)
const keyDepth = 4.0        // depth of internal sand key (mm)
const keyDraft = 60.0       // key draft angle (degrees)
const keyRatio = 0.85       // key height / height
const sideDraft = 3.0       // pattern side draft angle (degrees)

//-----------------------------------------------------------------------------

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
	s0 := sdf.Extrude3D(flaskSideProfile(width), height)
	s1 := sdf.Transform3D(s0, sdf.MirrorYZ())
	s2 := sdf.Union3D(s0, s1)

	w := width + 2.0*cornerLength

	// create the sand key
	s3 := sdf.Extrude3D(keyProfile(height), w)
	s3 = sdf.Transform3D(s3, sdf.RotateY(sdf.DtoR(90)))

	// create the side draft
	s4 := sdf.Extrude3D(sideDraftProfile(height), w)
	s4 = sdf.Transform3D(s4, sdf.RotateY(sdf.DtoR(90)))

	s5 := sdf.Union3D(s3, s4)

	return sdf.Difference3D(s2, s5)

}

//-----------------------------------------------------------------------------

func main() {
	widths := []float64{150, 200, 250, 300}
	height := 95.0
	for _, w := range widths {
		name := fmt.Sprintf("flask_%d.stl", int(w))
		sdf.RenderSTL(sdf.ScaleUniform3D(flaskSide(w, height), shrink), 300, name)
	}
}

//-----------------------------------------------------------------------------
