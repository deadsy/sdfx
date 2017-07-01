package main

import (
	"fmt"

	. "github.com/deadsy/sdfx/sdf"
)

func main() {
	// create a random set of vertices
	b := NewBox2(V2{0, 0}, V2{20, 20})
	s := b.RandomSet(20)
	pixels := V2i{800, 800}
	k := 1.5
	path := "voronoi.png"

	// use a 0 radius circle as a point
	s0 := Circle2D(0.0)
	// build an SDF for the points
	var s1 SDF2
	for _, p := range s {
		s1 = Union2D(s1, Transform2D(s0, Translate2d(p)))
	}

	// work out the region we will sample
	bb := s1.BoundingBox().ScaleAboutCenter(V2{k, k})

	fmt.Printf("rendering %s (%dx%d)\n", path, pixels[0], pixels[1])
	d, err := NewPNG(path, bb, pixels)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	d.RenderSDF2(s1)

	// create the delaunay triangulation
	ts, _ := s.Delaunay2d()
	// render the triangles
	for _, t := range ts {
		d.Triangle(t.ToTriangle2(s))
	}

	d.Save()
}
