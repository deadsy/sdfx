package main

import (
	"fmt"
	"github.com/deadsy/sdfx/sdf"
)

func main() {

	s := sdf.V2{2, 3}
	a := s.MulScalar(0.8)
	b := a.Negate()

	x := sdf.NewRectangleSDF(s)

	for i := 0; i < 20000; i++ {
		p := sdf.RandomV2(a, b)

		e0 := x.Evaluate(p)
		e1 := x.Evaluate2(p)

		if e0 != e1 {
			fmt.Printf("%+v %f %f\n", p, e0, e1)
		}

	}
}
