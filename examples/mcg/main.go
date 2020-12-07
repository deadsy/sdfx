//-----------------------------------------------------------------------------
/*

Mac Cheese Grater Plate
http://saccade.com/blog/2019/06/how-to-make-apples-mac-pro-holes/

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

// colSpace returns the space between columns
func colSpace(radius float64) float64 {
	return (4.0 * radius) / math.Sqrt(3.0)
}

// rowSpace returns the space between rows
func rowSpace(radius float64) float64 {
	return 2.0 * radius
}

// xOffset returns the x-offset between adjacent rows
func xOffset(radius float64) float64 {
	return (2.0 * radius) / math.Sqrt(3.0)
}

// yOffset returns the y-offset between adjacent rows
func yOffset(radius float64) float64 {
	return (2.0 * radius) / 3.0
}

// zOffset returns the z-offset between ball grids
func zOffset(radius float64) float64 {
	return (4.0 * radius) / 3.0
}

//-----------------------------------------------------------------------------

// ballRow returns a ball row
func ballRow(ncol int, radius float64) sdf.SDF3 {

	space := colSpace(radius)
	x := sdf.V3{-0.5 * ((float64(ncol) - 1) * space), 0, 0}
	dx := sdf.V3{space, 0, 0}

	var balls []sdf.SDF3
	s := sdf.Sphere3D(radius)
	for i := 0; i < ncol; i++ {
		balls = append(balls, sdf.Transform3D(s, sdf.Translate3d(x)))
		x = x.Add(dx)
	}
	return sdf.Union3D(balls...)
}

// ballGrid returns a ball grid
func ballGrid(
	ncol int, // number of columns
	nrow int, // number of rows
	radius float64, // radius of ball
) sdf.SDF3 {

	space := rowSpace(radius)
	x := sdf.V3{0, -0.5 * ((float64(nrow) - 1) * space), 0}
	dy0 := sdf.V3{-xOffset(radius), space, 0}
	dy1 := sdf.V3{xOffset(radius), space, 0}

	var rows []sdf.SDF3
	s := ballRow(ncol, radius)
	for i := 0; i < nrow; i++ {
		rows = append(rows, sdf.Transform3D(s, sdf.Translate3d(x)))
		if i%2 == 0 {
			x = x.Add(dy0)
		} else {
			x = x.Add(dy1)
		}
	}
	return sdf.Union3D(rows...)
}

// macCheeseGrater returns a Apple Mac style cheese grater plate.
func macCheeseGrater(
	ncol int, // number of columns
	nrow int, // number of rows
	radius float64, // radius of ball
) sdf.SDF3 {

	dx := sdf.V3{xOffset(radius), yOffset(radius), zOffset(radius)}.MulScalar(0.5)
	g := ballGrid(ncol, nrow, radius)
	g0 := sdf.Transform3D(g, sdf.Translate3d(dx.Neg()))
	g1 := sdf.Transform3D(g, sdf.Translate3d(dx))
	balls := sdf.Union3D(g0, g1)

	pX := colSpace(radius) * (float64(ncol) - 1)
	pY := rowSpace(radius) * (float64(nrow) - 1)
	pZ := 0.5 * colSpace(radius)
	plate := sdf.Box3D(sdf.V3{pX, pY, pZ}, 0)

	return sdf.Difference3D(plate, balls)
}

//-----------------------------------------------------------------------------

func main() {
	s := macCheeseGrater(15, 6, 10.0)
	render.RenderSTL(sdf.ScaleUniform3D(s, shrink), 500, "mcg.stl")
}

//-----------------------------------------------------------------------------
