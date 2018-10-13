//-----------------------------------------------------------------------------
/*

Devo Energy Dome

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func shell(r, h float64) SDF3 {

	k := 0.8
	stepH0 := h
	stepH1 := stepH0 * k
	stepH2 := stepH1 * k
	stepH3 := stepH2 * k

	height := stepH0 + stepH1 + stepH2 + stepH3
	fmt.Printf("height %f inches\n", height/MM_PER_INCH)

	stepX := (r / 4.0) * 0.75
	stepX0 := stepX * 0.20
	stepX1 := stepX - stepX0

	fillet := 4.0

	p := NewPolygon()

	p.Add(0, 0)
	p.Add(r, 0).Rel()

	p.Add(-stepX0, stepH0).Rel().Smooth(fillet, 4)
	p.Add(-stepX1, 0).Rel().Smooth(fillet, 4)

	p.Add(-stepX0, stepH1).Rel().Smooth(fillet, 4)
	p.Add(-stepX1, 0).Rel().Smooth(fillet, 4)

	p.Add(-stepX0, stepH2).Rel().Smooth(fillet, 4)
	p.Add(-stepX1, 0).Rel().Smooth(fillet, 4)

	p.Add(-stepX0, stepH3).Rel().Smooth(fillet, 4)
	p.Add(0, height)

	return Revolve3D(Polygon2D(p.Vertices()))
}

//-----------------------------------------------------------------------------

func main() {
	radius := (9.5 * MM_PER_INCH) / 2.0
	h0 := 2.05 * MM_PER_INCH
	wall := 4.0

	outer := shell(radius, h0)

	inner := shell(radius-wall, h0)
	inner = Transform3D(inner, Translate3d(V3{0, 0, -wall}))

	RenderSTL(Difference3D(outer, inner), 150, "energy_dome.stl")
}

//-----------------------------------------------------------------------------
