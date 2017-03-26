//-----------------------------------------------------------------------------
/*

distance minimisation for cubic splines

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func test1() {

	knot := []V2{
		V2{5, 0},
		V2{2, 2},
		V2{0, 4},
		V2{-1, 1},
		V2{-3, 0},
		V2{-2, -2},
		V2{0, -6},
		V2{2, -2},
		V2{5, 0},
	}

	s0 := CubicSpline2D(knot)
	s1 := s0.(*CubicSplineSDF2).PolySpline2D(300)
	s2 := Extrude3D(s1, 1)
	RenderSTL(s2, 300, "spline.stl")
	SDF2_RenderPNG(s1, "spline.png")
}

//-----------------------------------------------------------------------------

func test4() {

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

func main() {
	test4()
}

//-----------------------------------------------------------------------------
