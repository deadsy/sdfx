package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func test1() {
	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.8, 1.2}, 0.05)
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(225))
	sdf.Render(s1, true)
}

func test2() {
	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.8, 1.2}, 0.1)
	s1 := sdf.NewExtrudeSDF3(s0, 0.3)
	sdf.Render(s1, true)
}

func test3() {
	s0 := sdf.NewCircleSDF2(0.1)
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{1, 0}))
	s1 := sdf.NewSorSDF3(s0)
	sdf.Render(s1, true)
}

func test4() {
	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.2, 0.4}, 0.05)
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{1, 0}))
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(270))
	sdf.Render(s1, true)
}

func test5() {
	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.2, 0.4}, 0.05)
	s0 = sdf.NewTransformSDF2(s0, sdf.Rotate2d(sdf.DtoR(45)))
	s0 = sdf.NewTransformSDF2(s0, sdf.Translate2d(sdf.V2{1, 0}))
	s1 := sdf.NewSorThetaSDF3(s0, sdf.DtoR(315))
	//s1 := sdf.NewSorSDF3(s0)
	sdf.Render(s1, true)
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
	sdf.Render(s3, true)
}

func main() {
	test6()
}
