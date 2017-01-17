package main

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
)

func test1() {

	//s0 := sdf.NewRoundedBoxSDF3(sdf.V3{0.4, 0.8, 1.2}, 0.05)
	//s0 := sdf.NewBoxSDF2(sdf.V2{0.8, 1.2})

	s0 := sdf.NewRoundedBoxSDF2(sdf.V2{0.8, 1.2}, 0.05)
	s1 := sdf.NewSorSDF3(s0)

	sdf.Render(s1, true)
}

func test2() {
	a := sdf.Test33
	b := sdf.Test33_Inv
	c := a.Mul(b)

	fmt.Printf("%+v\n", a)
	fmt.Printf("%+v\n", b)
	fmt.Printf("%+v\n", c)
	fmt.Printf("%f\n", sdf.Test33_Det.Determinant())

}

func main() {
	test2()
}
