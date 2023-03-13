//-----------------------------------------------------------------------------
/*

hole patterns

Experiments in distributing holes on a circular disk.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

func holes() (sdf.SDF2, error) {

	l := 1.0
	circleRadius := l * 0.2
	n := 6
	steps := 11

	// Base circle
	c, err := sdf.Circle2D(circleRadius)
	if err != nil {
		return nil, err
	}
	s := c

	dBase := 0.0

	for i := 1; i <= steps; i++ {
		k := i * n
		r := float64(i) * l
		dTheta := sdf.Tau / float64(k)
		c0 := sdf.Transform2D(c, sdf.Translate2d(v2.Vec{0, r}))

		for j := 0; j < k; j++ {
			c1 := sdf.Transform2D(c0, sdf.Rotate2d(dBase+float64(j)*dTheta))
			s = sdf.Union2D(s, c1)
		}

		dBase += 0.5 * dTheta
	}

	return s, nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := holes()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToDXF(s, "holes.dxf", render.NewMarchingSquaresQuadtree(300))
}

//-----------------------------------------------------------------------------
