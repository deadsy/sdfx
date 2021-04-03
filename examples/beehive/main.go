//-----------------------------------------------------------------------------
/*

Bee Hive Parts

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func holePattern(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = byte('x')
	}
	return string(s)
}

func entranceReducer() (sdf.SDF3, error) {

	const zSize = 4.0
	const xSize = 6.0 * sdf.MillimetresPerInch
	const ySize = 1.9 * sdf.MillimetresPerInch

	k := obj.PanelParms{
		Size:         sdf.V2{xSize, ySize},
		CornerRadius: 5.0,
	}
	s, err := obj.Panel2D(&k)
	if err != nil {
		return nil, err
	}

	const holeHeight = (3.0 / 8.0) * sdf.MillimetresPerInch * 0.5
	const holeRadius = (3.0 / 8.0) * sdf.MillimetresPerInch * 0.5
	hole := sdf.Line2D(2*holeHeight, holeRadius)
	hole = sdf.Transform2D(hole, sdf.Rotate2d(sdf.DtoR(90)))

	const entranceSize = 4.0 * sdf.MillimetresPerInch
	const n = 6
	const gap = (entranceSize - (n * holeRadius)) / (n + 1)
	const yOfs = -ySize * 0.5
	const xOfs = (n - 1) * (holeRadius + gap) * 0.5
	p0 := sdf.V2{-xOfs, yOfs}
	p1 := sdf.V2{xOfs + holeRadius + gap, yOfs}
	hole = sdf.LineOf2D(hole, p0, p1, holePattern(n))

	return sdf.Extrude3D(sdf.Difference2D(s, hole), zSize), nil
}

//-----------------------------------------------------------------------------

func main() {

	p0, err := entranceReducer()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(p0, shrink), 300, "reducer.stl")

}

//-----------------------------------------------------------------------------
