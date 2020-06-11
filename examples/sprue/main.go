//-----------------------------------------------------------------------------
/*

Tapered Casting Sprue

Generate a model for a tapered pouring sprue.
Metal accelerates as it falls through the sprue but to maintain laminar flow the
vol/time at any point in the sprue must be constant. To have this the cross
sectional area gets smaller as the metal falls through the sprue.

In general:

a = sprue cross sectional area
h = sprue height

a * sqrt(h) = constant

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const steps = 20

func sprue(r, l, k float64) sdf.SDF3 {

	a0 := math.Pi * r * r
	h0 := math.Pow(k/a0, 2)
	dh := l / steps

	p := sdf.NewPolygon()
	p.Add(0, 0)
	for h := 0.0; h <= l; h += dh {
		a := k / math.Sqrt(h+h0)
		r := math.Sqrt(a / math.Pi)
		p.Add(r, h)
	}
	p.Add(0, l)

	s := sdf.Polygon2D(p.Vertices())
	return sdf.Revolve3D(s)
}

//-----------------------------------------------------------------------------

func main() {
	sdf.RenderSTL(sdf.ScaleUniform3D(sprue(20, 100, 3000), shrink), 300, "sprue.stl")
}

//-----------------------------------------------------------------------------
