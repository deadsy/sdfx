package main

import (
	"fmt"
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	. "github.com/deadsy/sdfx/sdf"
)

func test1() error {
	s0 := Box2D(V2{0.8, 1.2}, 0.05)
	s1, err := RevolveTheta3D(s0, DtoR(225))
	if err != nil {
		return err
	}
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test2() error {
	s0 := Box2D(V2{0.8, 1.2}, 0.1)
	s1 := Extrude3D(s0, 0.3)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test3() error {
	s0, err := Circle2D(0.1)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(V2{1, 0}))
	s1, err := Revolve3D(s0)
	if err != nil {
		return err
	}
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test4() error {
	s0 := Box2D(V2{0.2, 0.4}, 0.05)
	s0 = Transform2D(s0, Translate2d(V2{1, 0}))
	s1, err := RevolveTheta3D(s0, DtoR(270))
	if err != nil {
		return err
	}
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test5() error {
	s0 := Box2D(V2{0.2, 0.4}, 0.05)
	s0 = Transform2D(s0, Rotate2d(DtoR(45)))
	s0 = Transform2D(s0, Translate2d(V2{1, 0}))
	s1, err := RevolveTheta3D(s0, DtoR(315))
	if err != nil {
		return err
	}
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test6() error {
	s0, err := Sphere3D(0.5)
	if err != nil {
		return err
	}
	d := 0.4
	s1 := Transform3D(s0, Translate3d(V3{0, d, 0}))
	s2 := Transform3D(s0, Translate3d(V3{0, -d, 0}))
	s3 := Union3D(s1, s2)
	s3.(*UnionSDF3).SetMin(PolyMin(0.1))
	render.RenderSTL(s3, 200, "test.stl")
	return nil
}

func test7() error {
	s0, err := Box3D(V3{0.8, 0.8, 0.05}, 0)
	if err != nil {
		return err
	}
	s1 := Transform3D(s0, Rotate3d(V3{1, 0, 0}, DtoR(60)))
	s2 := Union3D(s0, s1)
	s2.(*UnionSDF3).SetMin(PolyMin(0.1))
	s3 := Transform3D(s2, Rotate3d(V3{0, 0, 1}, DtoR(-30)))
	render.RenderSTL(s3, 200, "test.stl")
	return nil
}

func test9() error {
	s, err := Sphere3D(10.0)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 200, "test.stl")
	return nil
}

func test10() error {
	s0, err := Box3D(V3{0.8, 0.8, 0.05}, 0)
	if err != nil {
		return err
	}
	s1 := Transform3D(s0, Rotate3d(V3{1, 0, 0}, DtoR(60)))
	s := Union3D(s0, s1)
	s.(*UnionSDF3).SetMin(PolyMin(0.1))
	render.RenderSTL(s, 200, "test.stl")
	return nil
}

func test11() error {
	s, err := Capsule3D(0.3, 1.4)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 200, "test.stl")
	return nil
}

func test12() error {
	k := 0.1
	points := []V2{
		{0, -k},
		{k, k},
		{-k, k},
	}
	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(V2{0.8, 0}))
	s1, err := RevolveTheta3D(s0, DtoR(360))
	if err != nil {
		return err
	}
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test13() error {
	k := 0.4
	s0, err := Polygon2D([]V2{{k, -k}, {k, k}, {-k, k}, {-k, -k}})
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(V2{0.8, 0}))
	s1, err := RevolveTheta3D(s0, DtoR(270))
	if err != nil {
		return err
	}
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test14() error {

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

	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s1, err := RevolveTheta3D(s0, DtoR(300))
	if err != nil {
		return err
	}
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test15() error {
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

	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Rotate2d(DtoR(theta)))
	s0 = Transform2D(s0, Translate2d(V2{j, k}))

	s1, err := RevolveTheta3D(s0, DtoR(300))
	if err != nil {
		return err
	}

	s1 = Transform3D(s1, Rotate3d(V3{0, 0, 1}, DtoR(30)))

	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test16() error {
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

	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Rotate2d(DtoR(theta)))
	s0 = Transform2D(s0, Translate2d(V2{j, k}))

	s1, err := RevolveTheta3D(s0, DtoR(300))
	if err != nil {
		return err
	}

	s1 = Transform3D(s1, Rotate3d(V3{0, 0, 1}, DtoR(30)))

	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test17() error {
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

	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(V2{j, k}))

	s1, err := RevolveTheta3D(s0, DtoR(300))
	if err != nil {
		return err
	}

	s1 = Transform3D(s1, Rotate3d(V3{0, 0, 1}, DtoR(30)))

	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test18() error {

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

	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s1, err := RevolveTheta3D(s0, DtoR(300))
	if err != nil {
		return err
	}

	s1 = Transform3D(s1, Rotate3d(V3{0, 0, 1}, DtoR(30)))

	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test19() error {
	r := 2.0
	k := 1.9
	s0, err := Circle2D(r)
	if err != nil {
		return err
	}
	s1 := Array2D(s0, V2i{3, 7}, V2{k * r, k * r})
	s1.(*ArraySDF2).SetMin(PolyMin(0.8))
	s2 := Extrude3D(s1, 1.0)
	render.RenderSTL(s2, 200, "test.stl")
	return nil
}

func test20() error {
	r := 4.0
	d := 20.0
	s0, err := Circle2D(r)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(V2{d, 0}))
	s0 = RotateUnion2D(s0, 5, Rotate2d(DtoR(20)))
	s0.(*RotateUnionSDF2).SetMin(PolyMin(1.2))
	s1 := Extrude3D(s0, 10.0)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test21() error {
	r := 2.0
	k := 1.9
	s0, err := Sphere3D(r)
	if err != nil {
		return err
	}
	s1 := Array3D(s0, V3i{3, 7, 5}, V3{k * r, k * r, k * r})
	s1.(*ArraySDF3).SetMin(PolyMin(0.8))
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test22() error {
	r := 4.0
	d := 20.0
	s0, err := Sphere3D(r)
	if err != nil {
		return err
	}
	s0 = Transform3D(s0, Translate3d(V3{d, 0, 0}))
	s0 = RotateUnion3D(s0, 5, Rotate3d(V3{0, 0, 1}, DtoR(20)))
	s0.(*RotateUnionSDF3).SetMin(PolyMin(1.2))
	render.RenderSTL(s0, 200, "test.stl")
	return nil
}

func test26() error {
	s, err := Cylinder3D(5, 2, 1)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 200, "test.stl")
	return nil
}

func test27() error {
	r := 5.0
	posn := V3Set{{2 * r, 2 * r, 0}, {-r, r, 0}, {r, -r, 0}, {-r, -r, 0}, {0, 0, 0}}
	cylinder, err := Cylinder3D(3, 1, 0)
	if err != nil {
		return err
	}
	s := Multi3D(cylinder, posn)
	render.RenderSTL(s, 200, "test.stl")
	return nil
}

func test28() error {
	s, err := Cone3D(20, 12, 8, 2)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 200, "test.stl")
	return nil
}

func test29() error {
	s0 := Line2D(10, 3)
	s1 := Extrude3D(s0, 4)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test30() error {
	s0 := Line2D(10, 3)
	s0 = Cut2D(s0, V2{4, 0}, V2{1, 1})
	s1 := Extrude3D(s0, 4)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test31() error {
	s, err := obj.CounterSunkHole3D(30, 2)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 200, "test.stl")
	return nil
}

func test32() error {
	s0, err := MakeFlatFlankCam(0.094, DtoR(2.0*57.5), 0.625)
	if err != nil {
		return err
	}
	s1 := Extrude3D(s0, 0.1)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test33() error {
	s0, err := ThreeArcCam2D(30, 20, 5, 50000)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", s0)
	s1 := Extrude3D(s0, 4)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test34() error {
	s0, err := MakeThreeArcCam(0.1, DtoR(2.0*80), 0.7, 1.1)
	if err != nil {
		return err
	}
	s1 := Extrude3D(s0, 0.1)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test35() error {
	r := 7.0
	d := 20.0
	s0 := Line2D(r, 1.0)
	s0 = Transform2D(s0, Translate2d(V2{d, 0}))
	s0 = RotateCopy2D(s0, 15)
	s1 := Extrude3D(s0, 10.0)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test36() error {
	k := obj.GenevaParms{
		NumSectors:     6,
		CenterDistance: 100,
		DriverRadius:   40,
		DrivenRadius:   80,
		PinRadius:      5,
		Clearance:      0.5,
	}
	sDriver, sDriven, err := obj.Geneva2D(&k)
	if err != nil {
		return err
	}
	render.RenderSTL(Extrude3D(sDriver, 10), 200, "driver.stl")
	render.RenderSTL(Extrude3D(sDriven, 10), 200, "driven.stl")
	return nil
}

func test37() error {
	r := 5.0
	p := 2.0
	isoThread, err := ISOThread(r, p, true)
	if err != nil {
		return err
	}
	s, err := Screw3D(isoThread, 50, p, 1)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 400, "screw.stl")
	return nil
}

func test39() error {
	s0 := NewFlange1(30, 20, 10)
	fmt.Printf("%+v\n", s0)
	s1 := Extrude3D(s0, 5)
	render.RenderSTL(s1, 200, "test.stl")
	return nil
}

func test40() error {
	d := 30.0
	wall := 5.0
	s0, err := Box3D(V3{d, d, d}, wall/2)
	if err != nil {
		return err
	}
	s1, err := Box3D(V3{d - wall, d - wall, d}, wall/2)
	if err != nil {
		return err
	}
	s1 = Transform3D(s1, Translate3d(V3{0, 0, wall / 2}))
	s := Difference3D(s0, s1)
	s.(*DifferenceSDF3).SetMax(PolyMax(2))
	render.RenderSTL(s, 200, "test.stl")
	return nil
}

func test41() error {
	s0, err := Cylinder3D(20.0, 5.0, 0)
	if err != nil {
		return err
	}
	s1 := Slice2D(s0, V3{0, 0, 0}, V3{0, 1, 1})
	s2, err := Revolve3D(s1)
	if err != nil {
		return err
	}
	render.RenderSTL(s2, 200, "test.stl")
	return nil
}

func test43() error {
	s0 := Line2D(10, 3)
	s0 = Cut2D(s0, V2{4, 0}, V2{1, 1})
	s1, err := ExtrudeRounded3D(s0, 4, 1)
	if err != nil {
		return err
	}
	render.RenderSTL(s1, 300, "test.stl")
	return nil
}

func test44() error {
	r := 100.0
	s0, err := Polygon2D(Nagon(5, r))
	if err != nil {
		return err
	}
	s1, err := Circle2D(r / 2)
	if err != nil {
		return err
	}
	s2, err := Loft3D(s1, s0, 200.0, 20.0)
	if err != nil {
		return err
	}
	render.RenderSTL(s2, 300, "test.stl")
	return err
}

func test49() error {
	s0, err := Circle2D(0.8)
	if err != nil {
		return err
	}
	render.RenderDXF(s0, 50, "test.dxf")
	return nil
}

func test50() error {
	k := obj.WasherParms{
		Thickness:   10,
		InnerRadius: 40,
		OuterRadius: 50,
		Remove:      0.3,
	}
	s, err := obj.Washer3D(&k)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 300, "test.stl")
	return nil
}

func test51() error {
	s, err := obj.StdPipe3D("sch40:1", "mm", 100)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 300, "test.stl")
	return nil
}

func test52() error {
	s, err := obj.StdPipeElbow3D("sch40:1", "mm", 30, 40)
	if err != nil {
		return err
	}
	render.RenderSTL(s, 300, "test.stl")
	return nil
}

func main() {
	err := test52()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
