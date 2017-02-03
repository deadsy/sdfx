package main

import . "github.com/deadsy/sdfx/sdf"

func main() {
	// use a 0 radius circle as a point
	s0 := NewCircleSDF2(0.0)
	// create a set of points at random locations
	var s1 SDF2
	for i := 0; i < 50; i++ {
		s1 = NewUnionSDF2(s1, NewTransformSDF2(s0, Translate2d(RandomV2(-10, 10))))
	}
	// render the distance field
	SDF2_RenderPNG(s1, "voronoi.png")
}
