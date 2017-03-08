// distance minimisation for cubic splines

package main

import (
	. "github.com/deadsy/sdfx/sdf"
)

func main() {

	data := []V2{
		V2{0, 10},
		V2{5, 5},
		V2{10, 10},
	}

	s := NewCubicSpline(data)
	s.Min1(V2{5, 7})
}
