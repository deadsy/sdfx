//-----------------------------------------------------------------------------
/*

Voronoi Diagram and Delaunay Triangulation

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v2i"
)

//-----------------------------------------------------------------------------

func main() {
	// create a random set of vertices
	b := sdf.NewBox2(v2.Vec{0, 0}, v2.Vec{20, 20})
	s := b.RandomSet(20)
	pixels := v2i.Vec{800, 800}
	k := 1.5
	path := "voronoi.png"

	// use a 0 radius circle as a point
	point, err := sdf.Circle2D(0.0)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// build an SDF for the points
	var s0 sdf.SDF2
	for i := range s {
		s0 = sdf.Union2D(s0, sdf.Transform2D(point, sdf.Translate2d(s[i])))
	}

	// work out the region we will sample
	bb := s0.BoundingBox().ScaleAboutCenter(k)
	log.Printf("rendering %s (%dx%d)\n", path, pixels.X, pixels.Y)
	d, err := render.NewPNG(path, bb, pixels)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	d.RenderSDF2(s0)

	// create the delaunay triangulation
	ts, _ := render.Delaunay2d(s)
	// render the triangles
	for _, t := range ts {
		d.Triangle(t.ToTriangle2(s))
	}

	d.Save()
}

//-----------------------------------------------------------------------------
