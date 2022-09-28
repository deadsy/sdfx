//-----------------------------------------------------------------------------
/*

Create a test panel for hole sizes.

Smallest diameter is 3 mm
Increase diameter in 0.2 mm increments
Largest diameter is 10.8 mm

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

// testHoles returns a panel with various holes for test fitting.
func testHoles() (sdf.SDF3, error) {

	const xInc = 15
	const yInc = 15
	const rInc = 0.1

	const nX = 5
	const nY = 8

	xOfs := 0.0
	yOfs := 0.0
	r := 1.5

	s := make([]sdf.SDF2, nX*nY)
	i := 0

	for j := 0; j < nY; j++ {
		for k := 0; k < nX; k++ {
			c, _ := sdf.Circle2D(r)
			s[i] = sdf.Transform2D(c, sdf.Translate2d(v2.Vec{xOfs, yOfs}))
			i++
			r += rInc
			xOfs += xInc
		}
		xOfs = 0.0
		yOfs += yInc
	}

	holes := sdf.Union2D(s...)
	xOfs = -float64(nX-1) * xInc * 0.5
	yOfs = -float64(nY-1) * yInc * 0.5
	holes = sdf.Transform2D(holes, sdf.Translate2d(v2.Vec{xOfs, yOfs}))

	// make a panel
	k := obj.PanelParms{
		Size:         v2.Vec{(nX + 1) * xInc, (nY + 1) * yInc},
		CornerRadius: xInc * 0.2,
	}
	panel, err := obj.Panel2D(&k)
	if err != nil {
		return nil, err
	}

	return sdf.Extrude3D(sdf.Difference2D(panel, holes), 3), nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := testHoles()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "test_holes.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
