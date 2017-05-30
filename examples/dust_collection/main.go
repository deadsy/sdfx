//-----------------------------------------------------------------------------
/*

Dust collection adapters

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// dust deputy tapered pipe
var dd_outer_d = 51.0
var dd_length = 39.0
var dd_taper = DtoR(2.0)

// vaccum hose 2.5" male fitting
var vh_outer_d = 58.0
var vh_length = 30.0
var vh_clearance = 0.5
var vh_taper = DtoR(0.5)

var wall_thickness = 4.0
var transition_length = 15.0

//-----------------------------------------------------------------------------

// dust deputy to 2.5" vacuum hose
func dd_to_hose25() {

	t := wall_thickness

	r0 := dd_outer_d / 2
	r1 := r0 - dd_length*math.Tan(dd_taper)
	r3 := (vh_outer_d + vh_clearance) / 2
	r2 := r3 - (vh_length * math.Tan(vh_taper))

	h0 := 0.0
	h1 := dd_length
	h2 := h1 + transition_length
	h3 := h2 + vh_length

	p := NewPolygon()
	p.Add(r0+t, h0)
	p.Add(r1+t, h1).Smooth(t, 4)
	p.Add(r2+t, h2).Smooth(t, 4)
	p.Add(r3+t, h3)
	p.Add(r3, h3)
	p.Add(r2, h2).Smooth(t, 4)
	p.Add(r1, h1).Smooth(t, 4)
	p.Add(r0, h0)

	s_2d := Polygon2D(p.Vertices())
	s_3d := Revolve3D(s_2d)

	RenderSTL(s_3d, 300, "adapter1.stl")
}

//-----------------------------------------------------------------------------

func main() {
	dd_to_hose25()
}

//-----------------------------------------------------------------------------
