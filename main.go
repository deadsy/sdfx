package main

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

func test1() {
	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.8, 1.2}, 0.05)
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(225))
	sdf.RenderPNG(s1, true)
}

func test2() {
	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.8, 1.2}, 0.1)
	s1 := sdf.NewExtrudeSDF3(s0, 0.3)
	sdf.RenderPNG(s1, true)
}

func test3() {
	s0 := sdf.NewCircleSDF2(0.1)
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{1, 0}))
	s1 := sdf.NewSorSDF3(s0)
	sdf.RenderPNG(s1, true)
}

func test4() {
	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.2, 0.4}, 0.05)
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{1, 0}))
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(270))
	sdf.RenderPNG(s1, true)
}

func test5() {
	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.2, 0.4}, 0.05)
	s0 = sdf.NewTransformSDF2(s0, sdf.Rotate2d(sdf.DtoR(45)))
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{1, 0}))
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(315))

	m := sdf.NewSDFMesh(s1, s1.BoundingBox().Scale(1.1), 0.01)
	err := sdf.SaveSTL("test.stl", m)
	if err != nil {
		fmt.Printf("%s", err)
	}

	//sdf.RenderPNG(s1, true)
}

func test6() {
	s0 := sdf.NewSphereSDF3(0.5)
	d := 0.4
	s1 := sdf.NewTransformSDF3(s0, sdf.Translate3d(sdf.V3{0, d, 0}))
	s2 := sdf.NewTransformSDF3(s0, sdf.Translate3d(sdf.V3{0, -d, 0}))
	//s3 := sdf.NewUnionRoundSDF3(s1, s2, 0.1)
	//s3 := sdf.NewUnionExpSDF3(s1, s2, 32)
	//s3 := sdf.NewUnionPowSDF3(s1, s2, 8)
	s3 := sdf.NewUnionPolySDF3(s1, s2, 0.1)
	//s3 := sdf.NewUnionChamferSDF3(s1, s2, 0.1)
	sdf.RenderPNG(s3, true)
}

func test7() {
	s0 := sdf.NewBoxSDF3(sdf.V3{0.8, 0.8, 0.05})
	s1 := sdf.NewTransformSDF3(s0, sdf.Rotate3d(sdf.V3{1, 0, 0}, sdf.DtoR(60)))
	s2 := sdf.NewUnionPolySDF3(s0, s1, 0.1)
	s3 := sdf.NewTransformSDF3(s2, sdf.Rotate3d(sdf.V3{0, 0, 1}, sdf.DtoR(-30)))
	sdf.RenderPNG(s3, true)
}

func test8() {
	a := sdf.V3{0, 0, 0}
	b := sdf.V3{1, 0, 0}
	c := sdf.V3{0, 1, 0}
	d := sdf.V3{0, 0, 1}
	t1 := sdf.NewTriangle(a, b, d)
	t2 := sdf.NewTriangle(a, c, b)
	t3 := sdf.NewTriangle(a, d, c)
	t4 := sdf.NewTriangle(b, c, d)
	m := sdf.NewMesh([]*sdf.Triangle{t1, t2, t3, t4})
	err := sdf.SaveSTL("test.stl", m)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func test9() {
	s := sdf.NewSphereSDF3(10.0)
	sdf.RenderSTL(s)
}

func test10() {
	s0 := sdf.NewBoxSDF3(sdf.V3{0.8, 0.8, 0.05})
	s1 := sdf.NewTransformSDF3(s0, sdf.Rotate3d(sdf.V3{1, 0, 0}, sdf.DtoR(60)))
	s := sdf.NewUnionPolySDF3(s0, s1, 0.1)
	sdf.RenderSTL(s)
}

func test11() {
	s := sdf.NewCapsuleSDF3(sdf.V3{0, -0.7, 0}, sdf.V3{0, 0.7, 0}, 0.3)
	sdf.RenderSTL(s)
}

func test12() {
	k := 0.1
	points := []sdf.V2{
		sdf.V2{0, -k},
		sdf.V2{k, k},
		sdf.V2{-k, k},
	}
	s0 := sdf.NewPolySDF2(points)
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{0.8, 0}))
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(360))
	sdf.RenderSTL(s1)
	//sdf.RenderPNG(s1, true)
}

func test13() {
	k := 0.4
	s0 := sdf.NewPolySDF2([]sdf.V2{sdf.V2{k, -k}, sdf.V2{k, k}, sdf.V2{-k, k}, sdf.V2{-k, -k}})
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{0.8, 0}))
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(270))
	sdf.RenderSTL(s1)
	//sdf.RenderPNG(s1, true)
}

func test14() {

	// size
	a := 0.3
	b := 0.7
	// rotation
	theta := 30.0
	c := math.Cos(sdf.DtoR(theta))
	s := math.Sin(sdf.DtoR(theta))
	// translate
	j := 10.0
	k := 2.0

	points := []sdf.V2{
		sdf.V2{j + c*a - s*b, k + s*a + c*b},
		sdf.V2{j - c*a - s*b, k - s*a + c*b},
		sdf.V2{j - c*a + s*b, k - s*a - c*b},
		sdf.V2{j + c*a + s*b, k + s*a - c*b},
	}

	s0 := sdf.NewPolySDF2(points)
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(300))

	sdf.RenderSTL(s1)
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

	points := []sdf.V2{
		sdf.V2{0, -b},
		sdf.V2{a, b},
		sdf.V2{-a, b},
	}

	s0 := sdf.NewPolySDF2(points)
	s0 = sdf.NewTransformSDF2(s0, sdf.Rotate2d(sdf.DtoR(theta)))
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{j, k}))

	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(300))
	s1 = sdf.NewTransformSDF3(s1, sdf.Rotate3d(sdf.V3{0, 0, 1}, sdf.DtoR(30)))

	sdf.RenderSTL(s1)
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

	points := []sdf.V2{
		sdf.V2{b0, -c},
		sdf.V2{a0, c},
		sdf.V2{-a1, c},
		sdf.V2{-b1, -c},
	}

	s0 := sdf.NewPolySDF2(points)
	s0 = sdf.NewTransformSDF2(s0, sdf.Rotate2d(sdf.DtoR(theta)))
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{j, k}))

	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(300))
	s1 = sdf.NewTransformSDF3(s1, sdf.Rotate3d(sdf.V3{0, 0, 1}, sdf.DtoR(30)))

	sdf.RenderSTL(s1)
}

func test17() {
	// size
	a := 1.3
	b := 0.4
	// translate
	j := 3.0
	k := 0.0

	points := []sdf.V2{
		sdf.V2{a, 0},
		sdf.V2{-a, b},
		sdf.V2{-a, -b},
	}

	s0 := sdf.NewPolySDF2(points)
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{j, k}))

	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(300))
	s1 = sdf.NewTransformSDF3(s1, sdf.Rotate3d(sdf.V3{0, 0, 1}, sdf.DtoR(30)))

	sdf.RenderSTL(s1)
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

	points := []sdf.V2{
		sdf.V2{0, 0},
		sdf.V2{r0, 0},
		sdf.V2{r0, h0},
		sdf.V2{r1, h1},
		sdf.V2{r2, h2},
		sdf.V2{r3, h3},
		sdf.V2{r3, h4},
		sdf.V2{0, h4},
	}

	s0 := sdf.NewPolySDF2(points)
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(300))
	s1 = sdf.NewTransformSDF3(s1, sdf.Rotate3d(sdf.V3{0, 0, 1}, sdf.DtoR(30)))

	sdf.RenderSTL(s1)
}

func main() {
	wheel()
}
