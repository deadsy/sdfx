package main

import (
	"fmt"

	"github.com/deadsy/sdfx/render"
	. "github.com/deadsy/sdfx/sdf"
)

const CUBIC_INCHES_PER_GALLON = 231.0

// pool dimensions are in inches
var pool_w = 234.0
var pool_l = 477.0

var pool_depth = []V2{
	{0.0, 43.0},
	{101.0, 46.0},
	{202.0, 58.0},
	{298.0, 83.0},
	{394.0, 96.0},
	{477.0, 96.0},
}

var vol = (7738.3005 * 1000.0) / CUBIC_INCHES_PER_GALLON // gallons

func main() {
	fmt.Printf("pool volume %f gallons\n", vol)

	p := NewPolygon()
	p.Add(0, 0)
	p.AddV2Set(pool_depth)
	p.Add(pool_l, 0)

	profile := Polygon2D(p.Vertices())
	pool := Extrude3D(profile, pool_w)
	render.RenderSTL(pool, 300, "pool.stl")
}
