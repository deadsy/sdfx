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

func main() {
	test1()
}
