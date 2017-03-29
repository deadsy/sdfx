//-----------------------------------------------------------------------------
/*

Solids built with Bezier Curves

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func bowling_pin() {

	b := NewBezier()
	b.Add(0, 0)
	b.Add(2.031/2.0, 0).HandleFwd(DtoR(45), 2)
	b.Add(4.766/2.0, 4.5).Handle(DtoR(90), 2, 2)
	b.Add(1.797/2.0, 10).Handle(DtoR(90), 3, 3)
	b.Add(2.547/2.0, 13.5).Handle(DtoR(90), 1, 1)
	b.Add(0, 15).HandleRev(DtoR(0), 1)
	b.Close()

	p := b.Polygon()
	p.Render("bowlingpin.dxf")

	s0 := Polygon2D(p.Vertices())
	s1 := Revolve3D(s0)
	RenderSTL(s1, 300, "bowlingpin.stl")
}

//-----------------------------------------------------------------------------

func bowl() {
	b := NewBezier()
	_ = b
}

//-----------------------------------------------------------------------------

func main() {
	bowling_pin()
}

//-----------------------------------------------------------------------------
