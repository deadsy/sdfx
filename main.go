package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {

	//s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.8, 1.2}, 0.05)
	s0 := sdf.NewBoxSDF2(sdf.V2{0.8, 1.2})
	s1 := sdf.NewSorSDF3(&s0)
	//	s := sdf.NewRoundedBoxSDF3(sdf.V3{0.4, 0.8, 1.2}, 0.05)
	sdf.Render(&s1, true)
}
