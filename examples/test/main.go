//-----------------------------------------------------------------------------
/*

Short Tests for primitives and objects.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	. "github.com/deadsy/sdfx/sdf"

	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v2i"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

func test1() error {
	s0 := Box2D(v2.Vec{0.8, 1.2}, 0.05)
	s1, err := RevolveTheta3D(s0, DtoR(225))
	if err != nil {
		return err
	}
	render.ToSTL(s1, "test1.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test2() error {
	s0 := Box2D(v2.Vec{0.8, 1.2}, 0.1)
	s1 := Extrude3D(s0, 0.3)
	render.ToSTL(s1, "test2.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test3() error {
	s0, err := Circle2D(0.1)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(v2.Vec{1, 0}))
	s1, err := Revolve3D(s0)
	if err != nil {
		return err
	}
	render.ToSTL(s1, "test3.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test4() error {
	s0 := Box2D(v2.Vec{0.2, 0.4}, 0.05)
	s0 = Transform2D(s0, Translate2d(v2.Vec{1, 0}))
	s1, err := RevolveTheta3D(s0, DtoR(270))
	if err != nil {
		return err
	}
	render.ToSTL(s1, "test4.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test5() error {
	s0 := Box2D(v2.Vec{0.2, 0.4}, 0.05)
	s0 = Transform2D(s0, Rotate2d(DtoR(45)))
	s0 = Transform2D(s0, Translate2d(v2.Vec{1, 0}))
	s1, err := RevolveTheta3D(s0, DtoR(315))
	if err != nil {
		return err
	}
	render.ToSTL(s1, "test5.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test6() error {
	s0, err := Sphere3D(0.5)
	if err != nil {
		return err
	}
	d := 0.4
	s1 := Transform3D(s0, Translate3d(v3.Vec{0, d, 0}))
	s2 := Transform3D(s0, Translate3d(v3.Vec{0, -d, 0}))
	s3 := Union3D(s1, s2)
	s3.(*UnionSDF3).SetMin(PolyMin(0.1))
	render.ToSTL(s3, "test6.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test7() error {
	s0, err := Box3D(v3.Vec{0.8, 0.8, 0.05}, 0)
	if err != nil {
		return err
	}
	s1 := Transform3D(s0, Rotate3d(v3.Vec{1, 0, 0}, DtoR(60)))
	s2 := Union3D(s0, s1)
	s2.(*UnionSDF3).SetMin(PolyMin(0.1))
	s3 := Transform3D(s2, Rotate3d(v3.Vec{0, 0, 1}, DtoR(-30)))
	render.ToSTL(s3, "test7.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test9() error {
	s, err := Sphere3D(10.0)
	if err != nil {
		return err
	}
	render.ToSTL(s, "test9.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test10() error {
	s0, err := Box3D(v3.Vec{0.8, 0.8, 0.05}, 0)
	if err != nil {
		return err
	}
	s1 := Transform3D(s0, Rotate3d(v3.Vec{1, 0, 0}, DtoR(60)))
	s := Union3D(s0, s1)
	s.(*UnionSDF3).SetMin(PolyMin(0.1))
	render.ToSTL(s, "test10.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test11() error {
	s, err := Capsule3D(3.0, 1.4)
	if err != nil {
		return err
	}
	render.ToSTL(s, "test11.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test12() error {
	k := 0.1
	points := []v2.Vec{
		{0, -k},
		{k, k},
		{-k, k},
	}
	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(v2.Vec{0.8, 0}))
	s1, err := RevolveTheta3D(s0, DtoR(360))
	if err != nil {
		return err
	}
	render.ToSTL(s1, "test12.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test13() error {
	k := 0.4
	s0, err := Polygon2D([]v2.Vec{{k, -k}, {k, k}, {-k, k}, {-k, -k}})
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(v2.Vec{0.8, 0}))
	s1, err := RevolveTheta3D(s0, DtoR(270))
	if err != nil {
		return err
	}
	render.ToSTL(s1, "test13.stl", render.NewMarchingCubesOctree(200))
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

	points := []v2.Vec{
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
	render.ToSTL(s1, "test14.stl", render.NewMarchingCubesOctree(200))
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

	points := []v2.Vec{
		{0, -b},
		{a, b},
		{-a, b},
	}

	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Rotate2d(DtoR(theta)))
	s0 = Transform2D(s0, Translate2d(v2.Vec{j, k}))

	s1, err := RevolveTheta3D(s0, DtoR(300))
	if err != nil {
		return err
	}

	s1 = Transform3D(s1, Rotate3d(v3.Vec{0, 0, 1}, DtoR(30)))

	render.ToSTL(s1, "test15.stl", render.NewMarchingCubesOctree(200))
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

	points := []v2.Vec{
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
	s0 = Transform2D(s0, Translate2d(v2.Vec{j, k}))

	s1, err := RevolveTheta3D(s0, DtoR(300))
	if err != nil {
		return err
	}

	s1 = Transform3D(s1, Rotate3d(v3.Vec{0, 0, 1}, DtoR(30)))

	render.ToSTL(s1, "test16.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test17() error {
	// size
	a := 1.3
	b := 0.4
	// translate
	j := 3.0
	k := 0.0

	points := []v2.Vec{
		{a, 0},
		{-a, b},
		{-a, -b},
	}

	s0, err := Polygon2D(points)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(v2.Vec{j, k}))

	s1, err := RevolveTheta3D(s0, DtoR(300))
	if err != nil {
		return err
	}

	s1 = Transform3D(s1, Rotate3d(v3.Vec{0, 0, 1}, DtoR(30)))

	render.ToSTL(s1, "test17.stl", render.NewMarchingCubesOctree(200))
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

	points := []v2.Vec{
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

	s1 = Transform3D(s1, Rotate3d(v3.Vec{0, 0, 1}, DtoR(30)))

	render.ToSTL(s1, "test18.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test19() error {
	r := 2.0
	k := 1.9
	s0, err := Circle2D(r)
	if err != nil {
		return err
	}
	s1 := Array2D(s0, v2i.Vec{3, 7}, v2.Vec{k * r, k * r})
	s1.(*ArraySDF2).SetMin(PolyMin(0.8))
	s2 := Extrude3D(s1, 1.0)
	render.ToSTL(s2, "test19.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test20() error {
	r := 4.0
	d := 20.0
	s0, err := Circle2D(r)
	if err != nil {
		return err
	}
	s0 = Transform2D(s0, Translate2d(v2.Vec{d, 0}))
	s0 = RotateUnion2D(s0, 5, Rotate2d(DtoR(20)))
	s0.(*RotateUnionSDF2).SetMin(PolyMin(1.2))
	s1 := Extrude3D(s0, 10.0)
	render.ToSTL(s1, "test20.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test21() error {
	r := 2.0
	k := 1.9
	s0, err := Sphere3D(r)
	if err != nil {
		return err
	}
	s1 := Array3D(s0, v3i.Vec{3, 7, 5}, v3.Vec{k * r, k * r, k * r})
	s1.(*ArraySDF3).SetMin(PolyMin(0.8))
	render.ToSTL(s1, "test21.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test22() error {
	r := 4.0
	d := 20.0
	s0, err := Sphere3D(r)
	if err != nil {
		return err
	}
	s0 = Transform3D(s0, Translate3d(v3.Vec{d, 0, 0}))
	s0 = RotateUnion3D(s0, 5, Rotate3d(v3.Vec{0, 0, 1}, DtoR(20)))
	s0.(*RotateUnionSDF3).SetMin(PolyMin(1.2))
	render.ToSTL(s0, "test22.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test26() error {
	s, err := Cylinder3D(5, 2, 1)
	if err != nil {
		return err
	}
	render.ToSTL(s, "test26.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test27() error {
	r := 5.0
	posn := v3.VecSet{{2 * r, 2 * r, 0}, {-r, r, 0}, {r, -r, 0}, {-r, -r, 0}, {0, 0, 0}}
	cylinder, err := Cylinder3D(3, 1, 0)
	if err != nil {
		return err
	}
	s := Multi3D(cylinder, posn)
	render.ToSTL(s, "test27.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test28() error {
	s, err := Cone3D(20, 12, 8, 2)
	if err != nil {
		return err
	}
	render.ToSTL(s, "test28.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test29() error {
	s0 := Line2D(10, 3)
	s1 := Extrude3D(s0, 4)
	render.ToSTL(s1, "test29.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test30() error {
	s0 := Line2D(10, 3)
	s0 = Cut2D(s0, v2.Vec{4, 0}, v2.Vec{1, 1})
	s1 := Extrude3D(s0, 4)
	render.ToSTL(s1, "test30.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test31() error {
	s, err := obj.CounterSunkHole3D(30, 2)
	if err != nil {
		return err
	}
	render.ToSTL(s, "test31.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test32() error {
	s0, err := MakeFlatFlankCam(0.094, DtoR(2.0*57.5), 0.625)
	if err != nil {
		return err
	}
	s1 := Extrude3D(s0, 0.1)
	render.ToSTL(s1, "cam0.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test33() error {
	s0, err := ThreeArcCam2D(30, 20, 5, 50000)
	if err != nil {
		return err
	}
	s1 := Extrude3D(s0, 4)
	render.ToSTL(s1, "cam1.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test34() error {
	s0, err := MakeThreeArcCam(0.1, DtoR(2.0*80), 0.7, 1.1)
	if err != nil {
		return err
	}
	s1 := Extrude3D(s0, 0.1)
	render.ToSTL(s1, "cam2.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test35() error {
	r := 7.0
	d := 20.0
	s0 := Line2D(r, 1.0)
	s0 = Transform2D(s0, Translate2d(v2.Vec{d, 0}))
	s0 = RotateCopy2D(s0, 15)
	s1 := Extrude3D(s0, 10.0)
	render.ToSTL(s1, "rotate_copy.stl", render.NewMarchingCubesOctree(200))
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
	render.ToSTL(Extrude3D(sDriver, 10), "driver.stl", render.NewMarchingCubesOctree(200))
	render.ToSTL(Extrude3D(sDriven, 10), "driven.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test37() error {
	r := 5.0
	p := 2.0
	isoThread, err := ISOThread(r, p, true)
	if err != nil {
		return err
	}
	s, err := Screw3D(isoThread, 50, 0, p, 1)
	if err != nil {
		return err
	}
	render.ToSTL(s, "screw.stl", render.NewMarchingCubesOctree(400))
	return nil
}

func test39() error {
	s0 := NewFlange1(30, 20, 10)
	s1 := Extrude3D(s0, 5)
	render.ToSTL(s1, "flange.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test40() error {
	d := 30.0
	wall := 5.0
	s0, err := Box3D(v3.Vec{d, d, d}, wall/2)
	if err != nil {
		return err
	}
	s1, err := Box3D(v3.Vec{d - wall, d - wall, d}, wall/2)
	if err != nil {
		return err
	}
	s1 = Transform3D(s1, Translate3d(v3.Vec{0, 0, wall / 2}))
	s := Difference3D(s0, s1)
	s.(*DifferenceSDF3).SetMax(PolyMax(2))
	render.ToSTL(s, "rounded_box.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test41() error {
	s0, err := Cylinder3D(20.0, 5.0, 0)
	if err != nil {
		return err
	}
	s1 := Slice2D(s0, v3.Vec{0, 0, 0}, v3.Vec{0, 1, 1})
	s2, err := Revolve3D(s1)
	if err != nil {
		return err
	}
	render.ToSTL(s2, "ellipsoid_egg.stl", render.NewMarchingCubesOctree(200))
	return nil
}

func test43() error {
	s0 := Line2D(10, 3)
	s0 = Cut2D(s0, v2.Vec{4, 0}, v2.Vec{1, 1})
	s1, err := ExtrudeRounded3D(s0, 4, 1)
	if err != nil {
		return err
	}
	render.ToSTL(s1, "cut2d.stl", render.NewMarchingCubesOctree(300))
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
	render.ToSTL(s2, "loft.stl", render.NewMarchingCubesOctree(300))
	return err
}

func test49() error {
	s0, err := Circle2D(0.8)
	if err != nil {
		return err
	}
	render.RenderDXF(s0, 50, "circle_2d.dxf")
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
	render.ToSTL(s, "washer.stl", render.NewMarchingCubesOctree(300))
	return nil
}

func test51() error {
	s, err := obj.StdPipe3D("sch40:1", "mm", 100)
	if err != nil {
		return err
	}
	render.ToSTL(s, "standard_pipe.stl", render.NewMarchingCubesOctree(300))
	return nil
}

//-----------------------------------------------------------------------------

type testFunc func() error

var testFuncs = []testFunc{
	test1,
	test2,
	test3,
	test4,
	test5,
	test6,
	test7,
	//test8,
	test9,
	test10,
	test11,
	test12,
	test13,
	test14,
	test15,
	test16,
	test17,
	test18,
	test19,
	test20,
	test21,
	test22,
	//test23,
	//test24,
	//test25,
	test26,
	test27,
	test28,
	test29,
	test30,
	test31,
	test32,
	test33,
	test34,
	test35,
	test36,
	test37,
	//test38,
	test39,
	test40,
	test41,
	//test42,
	test43,
	test44,
	//test45,
	//test46,
	//test47,
	//test48,
	test49,
	test50,
	test51,
}

func main() {
	for i, test := range testFuncs {
		err := test()
		if err != nil {
			log.Fatalf("error with testFuncs[%d]: %s\n", i, err)
		}
	}
}

//-----------------------------------------------------------------------------
