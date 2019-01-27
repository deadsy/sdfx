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

func dome(r, h, w float64) SDF3 {

	fillet := w

	// step heights
	k := 0.8
	stepH0 := h
	stepH1 := stepH0 * k
	stepH2 := stepH1 * k
	stepH3 := stepH2 * k

	height := stepH0 + stepH1 + stepH2 + stepH3
	fmt.Printf("height %f inches\n", height/MillimetresPerInch)

	// step ledges
	stepX := (r / 4.0) * 0.75
	stepX0 := stepX * 0.20
	stepX1 := stepX - stepX0

	// outer shell
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
	outer := Revolve3D(Polygon2D(p.Vertices()))

	// inner shell

	b := NewBezier()

	x := 0.0
	y := 0.0
	b.Add(x, y)

	x += r - w
	b.Add(x, y)

	x -= stepX
	y += stepH0 - w
	b.Add(x, y)

	x -= stepX
	y += stepH1
	b.Add(x, y)

	x -= stepX
	y += stepH2
	b.Add(x, y)

	y += stepH3
	b.Add(0, y)

	b.Close()

	inner := Revolve3D(Polygon2D(b.Polygon().Vertices()))

	return Difference3D(outer, inner)
}

//-----------------------------------------------------------------------------

func main() {
	radius := (9.5 * MillimetresPerInch) / 2.0
	h0 := 2.05 * MillimetresPerInch
	wall := 4.0

	s := dome(radius, h0, wall)
	//s = Cut3D(s, V3{0, 0, 0}, V3{0, 1, 0})
	RenderSTL(s, 150, "energy_dome.stl")
}

//-----------------------------------------------------------------------------
