//-----------------------------------------------------------------------------
/*

Pool Model

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/dc"
	"github.com/deadsy/sdfx/sdf"
	"log"
)

//-----------------------------------------------------------------------------

const cubicInchesPerGallon = 231.0

// pool dimensions are in inches
const poolWidth = 234.0
const poolLength = 477.0

var poolDepth = []sdf.V2{
	{0.0, 43.0},
	{101.0, 46.0},
	{202.0, 58.0},
	{298.0, 83.0},
	{394.0, 96.0},
	{477.0, 96.0},
}

const vol = (7738.3005 * 1000.0) / cubicInchesPerGallon // gallons

//-----------------------------------------------------------------------------

func pool() (sdf.SDF3, error) {
	log.Printf("pool volume %f gallons\n", vol)
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.AddV2Set(poolDepth)
	p.Add(poolLength, 0)
	profile, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(profile, poolWidth), nil
}

//-----------------------------------------------------------------------------

func main() {
	pool, err := pool()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(pool, 300, "pool1.stl", &render.MarchingCubesOctree{})
	render.ToSTL(pool, 15, "pool2.stl", dc.NewDualContouringDefault())

	// Test ImportSTL
	triChan := make(chan *render.Triangle3)
	go func() {
		dc.NewDualContouringDefault().Render(pool, 15, triChan)
		close(triChan)
	}()
	poolRebuilt := obj.ImportTriMesh(triChan, 0.2, 1)
	render.ToSTL(poolRebuilt, 15, "pool3.stl", dc.NewDualContouringDefault())

	//// Test ImportSTL 2
	//triChan2 := make(chan *render.Triangle3)
	//go func() {
	//	(&render.MarchingCubesOctree{}).Render(pool, 64, triChan2)
	//	close(triChan2)
	//}()
	//poolRebuilt2 := obj.ImportTriMesh(triChan2, 0.2, 1)
	//render.ToSTL(poolRebuilt2, 64, "pool4.stl", &render.MarchingCubesOctree{})
}

//-----------------------------------------------------------------------------
