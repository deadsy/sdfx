package main

import (
	"fmt"
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

func test1() {
	s0 := Box2D(V2{0.8, 1.2}, 0.05)
	s1 := RevolveTheta3D(s0, DtoR(225))
	RenderSTL(s1, 200, "test.stl")
}

func test2() {
	s0 := Box2D(V2{0.8, 1.2}, 0.1)
	s1 := Extrude3D(s0, 0.3)
	RenderSTL(s1, 200, "test.stl")
}

func test3() {
	s0 := Circle2D(0.1)
	s0 = Transform2D(s0, Translate2d(V2{1, 0}))
	s1 := Revolve3D(s0)
	RenderSTL(s1, 200, "test.stl")
}

func test4() {
	s0 := Box2D(V2{0.2, 0.4}, 0.05)
	s0 = Transform2D(s0, Translate2d(V2{1, 0}))
	s1 := RevolveTheta3D(s0, DtoR(270))
	RenderSTL(s1, 200, "test.stl")
}

func test5() {
	s0 := Box2D(V2{0.2, 0.4}, 0.05)
	s0 = Transform2D(s0, Rotate2d(DtoR(45)))
	s0 = Transform2D(s0, Translate2d(V2{1, 0}))
	s1 := RevolveTheta3D(s0, DtoR(315))

	RenderSTL(s1, 200, "test.stl")
}

func test6() {
	s0 := Sphere3D(0.5)
	d := 0.4
	s1 := Transform3D(s0, Translate3d(V3{0, d, 0}))
	s2 := Transform3D(s0, Translate3d(V3{0, -d, 0}))
	s3 := Union3D(s1, s2)
	s3.(*UnionSDF3).SetMin(PolyMin(0.1))
	RenderSTL(s3, 200, "test.stl")
}

func test7() {
	s0 := Box3D(V3{0.8, 0.8, 0.05}, 0)
	s1 := Transform3D(s0, Rotate3d(V3{1, 0, 0}, DtoR(60)))
	s2 := Union3D(s0, s1)
	s2.(*UnionSDF3).SetMin(PolyMin(0.1))
	s3 := Transform3D(s2, Rotate3d(V3{0, 0, 1}, DtoR(-30)))
	RenderSTL(s3, 200, "test.stl")
}

func test9() {
	s := Sphere3D(10.0)
	RenderSTL(s, 200, "test.stl")
}

func test10() {
	s0 := Box3D(V3{0.8, 0.8, 0.05}, 0)
	s1 := Transform3D(s0, Rotate3d(V3{1, 0, 0}, DtoR(60)))
	s := Union3D(s0, s1)
	s.(*UnionSDF3).SetMin(PolyMin(0.1))
	RenderSTL(s, 200, "test.stl")
}

func test11() {
	s := Capsule3D(0.3, 1.4)
	RenderSTL(s, 200, "test.stl")
}

func test12() {
	k := 0.1
	points := []V2{
		{0, -k},
		{k, k},
		{-k, k},
	}
	s0 := Polygon2D(points)
	s0 = Transform2D(s0, Translate2d(V2{0.8, 0}))
	s1 := RevolveTheta3D(s0, DtoR(360))
	RenderSTL(s1, 200, "test.stl")
}

func test13() {
	k := 0.4
	s0 := Polygon2D([]V2{{k, -k}, {k, k}, {-k, k}, {-k, -k}})
	s0 = Transform2D(s0, Translate2d(V2{0.8, 0}))
	s1 := RevolveTheta3D(s0, DtoR(270))
	RenderSTL(s1, 200, "test.stl")
}

func test14() {

	// size
	a := 0.3
	b := 0.7
	// rotation
	theta := 30.0
	c := math.Cos(DtoR(theta))
	s := math.Sin(DtoR(theta))
	// translate
	j := 10.0
	k := 2.0

	points := []V2{
		{j + c*a - s*b, k + s*a + c*b},
		{j - c*a - s*b, k - s*a + c*b},
		{j - c*a + s*b, k - s*a - c*b},
		{j + c*a + s*b, k + s*a - c*b},
	}

	s0 := Polygon2D(points)
	s1 := RevolveTheta3D(s0, DtoR(300))

	RenderSTL(s1, 200, "test.stl")
}

func test15() {
	// size
	a := 1.0
	b := 1.0
	// rotation
	theta := 0.0
	// translate
	j := 3.0
	k := 0.0

	points := []V2{
		{0, -b},
		{a, b},
		{-a, b},
	}

	s0 := Polygon2D(points)
	s0 = Transform2D(s0, Rotate2d(DtoR(theta)))
	s0 = Transform2D(s0, Translate2d(V2{j, k}))

	s1 := RevolveTheta3D(s0, DtoR(300))
	s1 = Transform3D(s1, Rotate3d(V3{0, 0, 1}, DtoR(30)))

	RenderSTL(s1, 200, "test.stl")
}

func test16() {
	// size
	a0 := 1.3
	b0 := 0.4
	a1 := 1.3
	b1 := 1.3
	c := 0.8
	// rotation
	theta := 20.0
	// translate
	j := 4.0
	k := 0.0

	points := []V2{
		{b0, -c},
		{a0, c},
		{-a1, c},
		{-b1, -c},
	}

	s0 := Polygon2D(points)
	s0 = Transform2D(s0, Rotate2d(DtoR(theta)))
	s0 = Transform2D(s0, Translate2d(V2{j, k}))

	s1 := RevolveTheta3D(s0, DtoR(300))
	s1 = Transform3D(s1, Rotate3d(V3{0, 0, 1}, DtoR(30)))

	RenderSTL(s1, 200, "test.stl")
}

func test17() {
	// size
	a := 1.3
	b := 0.4
	// translate
	j := 3.0
	k := 0.0

	points := []V2{
		{a, 0},
		{-a, b},
		{-a, -b},
	}

	s0 := Polygon2D(points)
	s0 = Transform2D(s0, Translate2d(V2{j, k}))

	s1 := RevolveTheta3D(s0, DtoR(300))
	s1 = Transform3D(s1, Rotate3d(V3{0, 0, 1}, DtoR(30)))

	RenderSTL(s1, 200, "test.stl")
}

func test18() {

	r0 := 10.0
	r1 := 8.0
	r2 := 7.5
	r3 := 9.0

	h0 := 4.0
	h1 := 6.0
	h2 := 5.5
	h3 := 3.5
	h4 := 1.0

	points := []V2{
		{0, 0},
		{r0, 0},
		{r0, h0},
		{r1, h1},
		{r2, h2},
		{r3, h3},
		{r3, h4},
		{0, h4},
	}

	s0 := Polygon2D(points)
	s1 := RevolveTheta3D(s0, DtoR(300))
	s1 = Transform3D(s1, Rotate3d(V3{0, 0, 1}, DtoR(30)))

	RenderSTL(s1, 200, "test.stl")
}

func test19() {
	r := 2.0
	k := 1.9
	s0 := Circle2D(r)
	s1 := Array2D(s0, V2i{3, 7}, V2{k * r, k * r})
	s1.(*ArraySDF2).SetMin(PolyMin(0.8))
	s2 := Extrude3D(s1, 1.0)
	RenderSTL(s2, 200, "test.stl")
}

func test20() {
	r := 4.0
	d := 20.0
	s0 := Circle2D(r)
	s0 = Transform2D(s0, Translate2d(V2{d, 0}))
	s0 = RotateUnion2D(s0, 5, Rotate2d(DtoR(20)))
	s0.(*RotateUnionSDF2).SetMin(PolyMin(1.2))
	s1 := Extrude3D(s0, 10.0)
	RenderSTL(s1, 200, "test.stl")
}

func test21() {
	r := 2.0
	k := 1.9
	s0 := Sphere3D(r)
	s1 := Array3D(s0, V3i{3, 7, 5}, V3{k * r, k * r, k * r})
	s1.(*ArraySDF3).SetMin(PolyMin(0.8))
	RenderSTL(s1, 200, "test.stl")
}

func test22() {
	r := 4.0
	d := 20.0
	s0 := Sphere3D(r)
	s0 = Transform3D(s0, Translate3d(V3{d, 0, 0}))
	s0 = RotateUnion3D(s0, 5, Rotate3d(V3{0, 0, 1}, DtoR(20)))
	s0.(*RotateUnionSDF3).SetMin(PolyMin(1.2))
	RenderSTL(s0, 200, "test.stl")
}

func test26() {
	s := Cylinder3D(5, 2, 1)
	RenderSTL(s, 200, "test.stl")
}

func test27() {
	r := 5.0
	posn := []V2{{2 * r, 2 * r}, {-r, r}, {r, -r}, {-r, -r}, {0, 0}}
	s := MultiCylinder3D(3, 1, posn)
	RenderSTL(s, 200, "test.stl")
}

func test28() {
	s := Cone3D(20, 12, 8, 2)
	RenderSTL(s, 200, "test.stl")
}

func test29() {
	s0 := Line2D(10, 3)
	s1 := Extrude3D(s0, 4)
	RenderSTL(s1, 200, "test.stl")
}

func test30() {
	s0 := Line2D(10, 3)
	s0 = Cut2D(s0, V2{4, 0}, V2{1, 1})
	s1 := Extrude3D(s0, 4)
	RenderSTL(s1, 200, "test.stl")
}

func test31() {
	s := CounterSunkHole3D(30, 2)
	RenderSTL(s, 200, "test.stl")
}

func test32() {
	s0, err := MakeFlatFlankCam(0.094, DtoR(2.0*57.5), 0.625)
	if err != nil {
		panic(err)
	}
	s1 := Extrude3D(s0, 0.1)
	RenderSTL(s1, 200, "test.stl")
}

func test33() {
	s0 := ThreeArcCam2D(30, 20, 5, 50000)
	fmt.Printf("%+v\n", s0)
	s1 := Extrude3D(s0, 4)
	RenderSTL(s1, 200, "test.stl")
}

func test34() {
	s0, err := MakeThreeArcCam(0.1, DtoR(2.0*80), 0.7, 1.1)
	if err != nil {
		panic(err)
	}
	s1 := Extrude3D(s0, 0.1)
	RenderSTL(s1, 200, "test.stl")
}

func test35() {
	r := 7.0
	d := 20.0
	s0 := Line2D(r, 1.0)
	s0 = Transform2D(s0, Translate2d(V2{d, 0}))
	s0 = RotateCopy2D(s0, 15)
	s1 := Extrude3D(s0, 10.0)
	RenderSTL(s1, 200, "test.stl")
}

func test36() {
	s_driver, s_driven, err := MakeGenevaCam(6, 100, 40, 80, 5, 0.5)
	if err != nil {
		panic(err)
	}
	RenderSTL(Extrude3D(s_driver, 10), 200, "driver.stl")
	RenderSTL(Extrude3D(s_driven, 10), 200, "driven.stl")
}

func test37() {
	r := 5.0
	p := 2.0
	s := Screw3D(ISOThread(r, p, "external"), 50, p, 1)
	RenderSTL(s, 400, "screw.stl")
}

func test38() {
	fmt.Printf("%+v\n", ThreadLookup("unc_1/4"))
	fmt.Printf("%+v\n", ThreadLookup("blah"))
}

func test39() {
	s0 := NewFlange1(30, 20, 10)
	fmt.Printf("%+v\n", s0)
	s1 := Extrude3D(s0, 5)
	RenderSTL(s1, 200, "test.stl")
}

func test40() {
	d := 30.0
	wall := 5.0
	s0 := Box3D(V3{d, d, d}, wall/2)
	s1 := Box3D(V3{d - wall, d - wall, d}, wall/2)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, wall / 2}))
	s := Difference3D(s0, s1)
	s.(*DifferenceSDF3).SetMax(PolyMax(2))
	RenderSTL(s, 200, "test.stl")
}

func test41() {
	s0 := Cylinder3D(20.0, 5.0, 0)
	s1 := Slice2D(s0, V3{0, 0, 0}, V3{0, 1, 1})
	s2 := Revolve3D(s1)
	RenderSTL(s2, 200, "test.stl")
}

func test42() {
	p := NewPolygon()
	p.Add(0, 0)
	p.Add(1, 1).Rel()
	p.Add(1, 0).Rel().Arc(-2, 4)
	p.Add(1, -1).Rel()
	p.Render("test.dxf")
}

func test43() {
	s0 := Line2D(10, 3)
	s0 = Cut2D(s0, V2{4, 0}, V2{1, 1})
	s1 := ExtrudeRounded3D(s0, 4, 1)
	RenderSTL(s1, 300, "test.stl")
}

func test44() {
	r := 100.0
	s0 := Polygon2D(Nagon(5, r))
	s1 := Circle2D(r / 2)
	s2 := Loft3D(s1, s0, 200.0, 20.0)
	RenderSTL(s2, 300, "test.stl")
}

func test45() {
	d := NewDXF("test.dxf")
	k := 1.0
	b := Box2{V2{0, 0}, V2{k, k}}
	s := b.RandomSet(1000)
	t, _ := s.SuperTriangle()
	r := k / (20.0 * math.Sqrt(float64(len(s))))
	d.Points(s, r)
	d.Triangle(t)
	d.Save()
}

func test46() {
	k := 10.0
	b := Box2{V2{0, 0}, V2{k, k}}
	s := b.RandomSet(325)

	ts1, _ := s.Delaunay2d()
	fmt.Printf("ts1 %d triangles\n", len(ts1))
	d := NewDXF("test1.dxf")
	for _, t := range ts1 {
		d.Triangle(t.ToTriangle2(s))
	}
	d.Save()

	ts2, _ := s.Delaunay2dSlow()
	fmt.Printf("ts2 %d triangles\n", len(ts2))
	d = NewDXF("test2.dxf")
	for _, t := range ts2 {
		d.Triangle(t.ToTriangle2(s))
	}
	d.Save()

	if ts1.Equals(ts2) {
		fmt.Printf("same\n")
	} else {
		fmt.Printf("different\n")
	}
}

func test47() {
	xsz := 3
	ysz := 3

	m, _ := NewMap2(NewBox2(V2{0, 0}, V2{20, 20}), V2i{xsz, ysz}, false)
	s := make(V2Set, 0, xsz*ysz)
	for i := 0; i < xsz; i++ {
		for j := 0; j < ysz; j++ {
			s = append(s, m.ToV2(V2i{i, j}))
		}
	}

	ts1, _ := s.Delaunay2d()
	fmt.Printf("ts1 %d triangles\n", len(ts1))
	d := NewDXF("test1.dxf")
	for _, t := range ts1 {
		d.Triangle(t.ToTriangle2(s))
	}
	d.Save()

	ts2, _ := s.Delaunay2dSlow()
	fmt.Printf("ts2 %d triangles\n", len(ts2))
	d = NewDXF("test2.dxf")
	for _, t := range ts2 {
		d.Triangle(t.ToTriangle2(s))
	}
	d.Save()

	if ts1.Equals(ts2) {
		fmt.Printf("same\n")
	} else {
		fmt.Printf("different\n")
	}
}

func test48() {
	s1 := Circle2D(0.8)
	s0 := Box2D(V2{3.0, 4.0}, 0.5)
	s := Difference2D(s0, Transform2D(s1, Translate2d(V2{0.5, 0.5})))

	p, err := GenerateMesh2D(s, V2i{20, 40})
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	ts, err := p.Delaunay2d()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("ts %d triangles\n", len(ts))
	d := NewDXF("test.dxf")
	for _, t := range ts {
		d.Triangle(t.ToTriangle2(p))
	}
	d.Save()
}

func test49() {
	s0 := Circle2D(0.8)
	RenderDXF(s0, 50, "test.dxf")
}

func main() {
	test49()
}
