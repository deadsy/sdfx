// distance minimisation for cubic splines

package main

import (
	. "github.com/deadsy/sdfx/sdf"
)

func test1() {

	data := []V2{
		V2{0, 10},
		V2{5, 5},
		V2{10, 10},
	}

	s := NewCubicSpline(data)
	s.Polygonize()
	s.Min1(V2{5, 7})
}

func test2() {
	data := []V2{
		V2{0, 1},
		V2{1, 2},
		V2{2, 3},
		V2{3, 4},
		V2{4, 3},
		V2{5, 2},
		V2{6, 1},
		V2{7, 0},
		V2{8, 1},
		V2{9, 2},
		V2{10, 3},
	}
	s := NewCubicSpline(data)
	s.Polygonize()
}

func main() {

	test2()

}
