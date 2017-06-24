package main

import . "github.com/deadsy/sdfx/sdf"

func main() {
	// define the bounding box for the point set
	b := NewBox2(V2{0, 0}, V2{20, 20})
	// use a 0 radius circle as a point
	s0 := Circle2D(0.0)
	// create a set of points at random locations
	var s1 SDF2
	for i := 0; i < 50; i++ {
		s1 = Union2D(s1, Transform2D(s0, Translate2d(b.Random())))
	}
	// render the distance field
	SDF2_RenderPNG(s1, "voronoi.png")
}
