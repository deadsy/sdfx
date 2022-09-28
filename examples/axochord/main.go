//-----------------------------------------------------------------------------
/*

OmniKeys 3 x 12 chord keys

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"
	"strings"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

var bR0 = 13.0 * 0.5        // major radius
var bR1 = 7.0 * 0.5         // minor radius
var bH0 = 6.0               // cavity height for button body
var bH1 = 1.5               // thru panel thickness
var bDeltaV = 22.0          // vertical inter-button distance
var bDeltaH = 20.0          // horizontal inter-button distance
var bTheta = sdf.DtoR(20.0) // button angle

const buttonsV = 3 // number of vertical buttons
const buttonsH = 3 //12 // number of horizontal buttons

func buttonCavity() (sdf.SDF3, error) {
	p := sdf.NewPolygon()
	p.Add(0, -(bH0 + bH1))
	p.Add(bR0, 0).Rel()
	p.Add(0, bH0).Rel()
	p.Add(bR1-bR0, 0).Rel()
	p.Add(0, bH1).Rel()
	p.Add(bR0-bR1, 0).Rel()
	p.Add(0, bH0).Rel()
	p.Add(bR1-bR0, 0).Rel()
	p.Add(0, bH1).Rel()
	p.Add(-bR1, 0).Rel()
	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	return sdf.Revolve3D(s)
}

// return the button matrix
func buttons() (sdf.SDF3, error) {
	// single key column
	d := buttonsV * bDeltaV
	p := v3.Vec{-math.Sin(bTheta) * d, math.Cos(bTheta) * d, 0}
	bc, err := buttonCavity()
	if err != nil {
		return nil, err
	}
	col := sdf.LineOf3D(bc, v3.Vec{}, p, strings.Repeat("x", buttonsV))
	// multiple key columns
	d = buttonsH * bDeltaH
	p = v3.Vec{d, 0, 0}
	matrix := sdf.LineOf3D(col, v3.Vec{}, p, strings.Repeat("x", buttonsH))
	// centered on the origin
	d = (buttonsV - 1) * bDeltaV
	dx := 0.5 * (((buttonsH - 1) * bDeltaH) - (d * math.Sin(bTheta)))
	dy := 0.5 * d * math.Cos(bTheta)
	return sdf.Transform3D(matrix, sdf.Translate3d(v3.Vec{-dx, -dy, 0})), nil
}

//-----------------------------------------------------------------------------

// https://geekhack.org/index.php?topic=47744.0
// https://cdn.sparkfun.com/datasheets/Components/Switches/MX%20Series.pdf

var cherryD0 = 0.551 * sdf.MillimetresPerInch
var cherryD1 = 0.614 * sdf.MillimetresPerInch
var cherryD2 = 0.1378 * sdf.MillimetresPerInch
var cherryD3 = 0.0386 * sdf.MillimetresPerInch

// cherryMX returns the SDF2 for a cherry MX plate cutout.
func cherryMX() (sdf.SDF2, error) {

	cherryOfs := ((cherryD0 / 2.0) - cherryD3) - (cherryD2 / 2.0)

	r0 := sdf.Box2D(v2.Vec{cherryD0, cherryD0}, 0)
	r1 := sdf.Box2D(v2.Vec{cherryD1, cherryD2}, 0)

	r2 := sdf.Transform2D(r1, sdf.Translate2d(v2.Vec{0, cherryOfs}))
	r3 := sdf.Transform2D(r1, sdf.Translate2d(v2.Vec{0, -cherryOfs}))

	r4 := sdf.Union2D(r2, r3)
	r5 := sdf.Transform2D(r4, sdf.Rotate2d(sdf.Pi*0.5))

	return sdf.Union2D(r0, r4, r5), nil
}

//-----------------------------------------------------------------------------

func panel() (sdf.SDF3, error) {
	v := (buttonsV - 1) * bDeltaV
	vx := v * math.Sin(bTheta)
	vy := v * math.Cos(bTheta)

	sx := ((buttonsH-1)*bDeltaH + vx) * 1.5
	sy := vy * 1.9

	pp := &obj.PanelParms{
		Size:         v2.Vec{sx, sy},
		CornerRadius: 5.0,
		HoleDiameter: 3.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	panel, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}
	// extrude to 3d
	return sdf.Extrude3D(panel, 2.0*(bH0+bH1)), nil
}

//-----------------------------------------------------------------------------

func main() {
	panel, err := panel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	buttons, err := buttons()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s := sdf.Difference3D(panel, buttons)
	upper := sdf.Cut3D(s, v3.Vec{}, v3.Vec{0, 0, 1})
	lower := sdf.Cut3D(s, v3.Vec{}, v3.Vec{0, 0, -1})

	render.ToSTL(upper, "upper.stl", render.NewMarchingCubesOctree(400))
	render.ToSTL(lower, "lower.stl", render.NewMarchingCubesOctree(400))

	cherryMX, err := cherryMX()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToDXF(cherryMX, "plate.dxf", render.NewMarchingSquaresQuadtree(400))
}

//-----------------------------------------------------------------------------
