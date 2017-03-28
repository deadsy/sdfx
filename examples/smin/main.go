//-----------------------------------------------------------------------------
/*

 */
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func test2() {

	b := NewBezier()

	d1 := 5.0
	d2 := 7.0
	k1 := 2.5
	k2 := 1.0

	b.Add(0, -d1).Handle(DtoR(0), k1, k1)
	b.Add(d1, 0).Handle(DtoR(90), k1, k1)
	b.Add(d2, d2).Handle(DtoR(135), k2, k2)
	b.Add(0, d1).Handle(DtoR(180), k1, k1)
	b.Add(-d1, 0).Handle(DtoR(270), k1, k1)
	b.Close()

	p := b.Polygon()
	p.Render("spline.dxf")

	s0 := Polygon2D(p.Vertices())
	s1 := Extrude3D(s0, 1)
	RenderSTL(s1, 300, "curve.stl")
}

//-----------------------------------------------------------------------------

func test3() {

	b := NewBezier()

	b.Add(0, 0).HandleFwd(DtoR(135), 3)
	b.Add(2, 0).HandleRev(DtoR(45), 3)
	b.Close()

	p := b.Polygon()
	p.Render("spline.dxf")

}

//-----------------------------------------------------------------------------

func test5() {
	b := NewBezier()
	b.Add(0, 0)
	b.Add(3, 1).Mid()
	b.Add(0, 2)
	p := b.Polygon()
	p.Render("spline.dxf")
}

//-----------------------------------------------------------------------------

func test4() {

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

func main() {
	test4()
}

//-----------------------------------------------------------------------------
