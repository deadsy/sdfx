package main

import (
	"fmt"
	"github.com/deadsy/sdfx/sdf"
)

func main() {

	s := sdf.V2{2, 3}
	a := s.MulScalar(0.8)
	b := a.Negate()

	x := sdf.NewBoxSDF2(s)

	for i := 0; i < 200; i++ {
		p := sdf.RandomV2(a, b)
		e0 := x.Evaluate(p)
		fmt.Printf("%+v %f\n", p, e0)
	}
}
