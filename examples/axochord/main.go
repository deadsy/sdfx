//-----------------------------------------------------------------------------
/*

OmniKeys 3 x 12 chord keys

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"
	"strings"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

var bR0 = 13.0 * 0.5    // major radius
var bR1 = 7.0 * 0.5     // minor radius
var bH0 = 6.0           // cavity height for button body
var bH1 = 1.5           // thru panel thickness
var bDeltaV = 22.0      // vertical inter-button distance
var bDeltaH = 20.0      // horizontal inter-button distance
var bTheta = DtoR(20.0) // button angle

const buttonsV = 3 // number of vertical buttons
const buttonsH = 3 //12 // number of horizontal buttons

func buttonCavity() SDF3 {
	p := NewPolygon()
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
	return Revolve3D(Polygon2D(p.Vertices()))
}

// return the button matrix
func buttons() SDF3 {
	// single key column
	d := buttonsV * bDeltaV
	p := V3{-math.Sin(bTheta) * d, math.Cos(bTheta) * d, 0}
	col := LineOf3D(buttonCavity(), V3{}, p, strings.Repeat("x", buttonsV))
	// multiple key columns
	d = buttonsH * bDeltaH
	p = V3{d, 0, 0}
	matrix := LineOf3D(col, V3{}, p, strings.Repeat("x", buttonsH))
	// centered on the origin
	d = (buttonsV - 1) * bDeltaV
	dx := 0.5 * (((buttonsH - 1) * bDeltaH) - (d * math.Sin(bTheta)))
	dy := 0.5 * d * math.Cos(bTheta)
	return Transform3D(matrix, Translate3d(V3{-dx, -dy, 0}))
}

//-----------------------------------------------------------------------------

// https://geekhack.org/index.php?topic=47744.0
// https://cdn.sparkfun.com/datasheets/Components/Switches/MX%20Series.pdf

var cherryD0 = 0.551 * MillimetresPerInch
var cherryD1 = 0.614 * MillimetresPerInch
var cherryD2 = 0.1378 * MillimetresPerInch
var cherryD3 = 0.0386 * MillimetresPerInch

// cherryMX returns the SDF2 for a cherry MX plate cutout.
func cherryMX() SDF2 {

	cherryOfs := ((cherryD0 / 2.0) - cherryD3) - (cherryD2 / 2.0)

	r0 := Box2D(V2{cherryD0, cherryD0}, 0)
	r1 := Box2D(V2{cherryD1, cherryD2}, 0)

	r2 := Transform2D(r1, Translate2d(V2{0, cherryOfs}))
	r3 := Transform2D(r1, Translate2d(V2{0, -cherryOfs}))

	r4 := Union2D(r2, r3)
	r5 := Transform2D(r4, Rotate2d(Pi*0.5))

	return Union2D(r0, r4, r5)
}

//-----------------------------------------------------------------------------

func panel() SDF3 {
	v := (buttonsV - 1) * bDeltaV
	vx := v * math.Sin(bTheta)
	vy := v * math.Cos(bTheta)

	sx := ((buttonsH-1)*bDeltaH + vx) * 1.5
	sy := vy * 1.9

	pp := &PanelParms{
		Size:         V2{sx, sy},
		CornerRadius: 5.0,
		HoleDiameter: 3.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	// extrude to 3d
	return Extrude3D(Panel2D(pp), 2.0*(bH0+bH1))
}

//-----------------------------------------------------------------------------

func main() {
	s := Difference3D(panel(), buttons())
	upper := Cut3D(s, V3{}, V3{0, 0, 1})
	lower := Cut3D(s, V3{}, V3{0, 0, -1})

	RenderSTL(upper, 400, "upper.stl")
	RenderSTL(lower, 400, "lower.stl")
	RenderDXF(cherryMX(), 400, "plate.dxf")
}

//-----------------------------------------------------------------------------
