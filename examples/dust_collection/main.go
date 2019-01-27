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
var dd_od = 51.0
var dd_taper = DtoR(2.0)
var dd_length = 39.0

// vacuum hose 2.5" male fitting
var vh_od = 58.0
var vh_clearance = 0.6
var vh_taper = DtoR(0.4)

// pvc pipe outside diameters
var pvc3_od = 3.26 * MillimetresPerInch
var pvc2_od = 2.375 * MillimetresPerInch

var wall_thickness = 4.0

//-----------------------------------------------------------------------------

// adapter: female dust deputy, female 2.5" vacuum
func fdd_to_fvh25() {

	t := wall_thickness
	transition_length := 15.0
	vh_length := 30.0

	r0 := dd_od / 2
	r1 := r0 - dd_length*math.Tan(dd_taper)
	r3 := (vh_od + vh_clearance) / 2
	r2 := r3 - (vh_length * math.Tan(vh_taper))

	h0 := 0.0
	h1 := h0 + dd_length
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

	s := Revolve3D(Polygon2D(p.Vertices()))
	RenderSTL(s, 150, "fdd_fvh25.stl")
}

//-----------------------------------------------------------------------------

// adapter: male 2.5" vacuum, male 3" pvc
func mvh25_to_mpvc(pvc_od float64) {

	t := wall_thickness
	transition_length := 15.0

	r0 := pvc_od / 2
	r1 := vh_od / 2

	h0 := 0.0
	h1 := h0 + 35.0
	h2 := h1 + transition_length
	h3 := h2 + 20.0

	p := NewPolygon()
	p.Add(r0, h0)
	p.Add(r0, h1).Smooth(t, 4)
	p.Add(r1, h2).Smooth(t, 4)
	p.Add(r1, h3)
	p.Add(r1-t, h3)
	p.Add(r1-t, h2).Smooth(t, 4)
	p.Add(r0-t, h1).Smooth(t, 4)
	p.Add(r0-t, h0)

	s := Revolve3D(Polygon2D(p.Vertices()))
	RenderSTL(s, 150, "mvh25_mpvc.stl")
}

//-----------------------------------------------------------------------------

// adapter: female dust deputy, male 3" pvc
func fdd_to_mpvc(pvc_od float64) {

	t := wall_thickness
	transition_length := 15.0

	r0 := pvc_od / 2
	r2 := (dd_od / 2) + t
	r1 := r2 - dd_length*math.Tan(dd_taper)

	h0 := 0.0
	h1 := h0 + 35.0
	h2 := h1 + transition_length
	h3 := h2 + dd_length

	p := NewPolygon()
	p.Add(r0, h0)
	p.Add(r0, h1).Smooth(t, 4)
	p.Add(r1, h2).Smooth(t, 4)
	p.Add(r2, h3)
	p.Add(r2-t, h3)
	p.Add(r1-t, h2).Smooth(t, 4)
	p.Add(r0-t, h1).Smooth(t, 4)
	p.Add(r0-t, h0)

	s := Revolve3D(Polygon2D(p.Vertices()))
	RenderSTL(s, 150, "fdd_mpvc.stl")
}

//-----------------------------------------------------------------------------

func main() {
	fdd_to_fvh25()
	mvh25_to_mpvc(pvc2_od)
	fdd_to_mpvc(pvc2_od)
}

//-----------------------------------------------------------------------------
