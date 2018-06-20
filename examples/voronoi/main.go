//-----------------------------------------------------------------------------
/*

Voronoi Diagram and Delaunay Triangulation

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {
	// create a random set of vertices
	b := NewBox2(V2{0, 0}, V2{20, 20})
	s := b.RandomSet(20)
	pixels := V2i{800, 800}
	k := 1.5
	path := "voronoi.png"

	// use a 0 radius circle as a point
	point := Circle2D(0.0)

	// build an SDF for the points
	var s0 SDF2
	for i := range s {
		s0 = Union2D(s0, Transform2D(point, Translate2d(s[i])))
	}

	// work out the region we will sample
	bb := s0.BoundingBox().ScaleAboutCenter(k)

	fmt.Printf("rendering %s (%dx%d)\n", path, pixels[0], pixels[1])
	d, err := NewPNG(path, bb, pixels)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	d.RenderSDF2(s0)

	// create the delaunay triangulation
	ts, _ := s.Delaunay2d()
	// render the triangles
	for _, t := range ts {
		d.Triangle(t.ToTriangle2(s))
	}

	d.Save()
}

//-----------------------------------------------------------------------------
