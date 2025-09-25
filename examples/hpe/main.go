//-----------------------------------------------------------------------------
/*

HPE AP-745 Mounting Board

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

var baseThickness = 3.0
var pillarHeight = 15.0

//-----------------------------------------------------------------------------

func ap745Standoffs() (sdf.SDF3, error) {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 138.0 * sdf.Mil * 2.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}

	positions0 := v3.VecSet{
		{0, 6102.36 * sdf.Mil, zOfs},                 // 138 mil
		{0, 0, zOfs},                                 // 138 mil
		{6102.36 * sdf.Mil, 6102.36 * sdf.Mil, zOfs}, // 138 mil
		{6102.36 * sdf.Mil, 0, zOfs},                 // 138 mil
		{3937.01 * sdf.Mil, 7047.24 * sdf.Mil, zOfs}, // 118 mil
		{1240.16 * sdf.Mil, 5570.87 * sdf.Mil, zOfs}, // 118 mil
		{2648.46 * sdf.Mil, 3485.15 * sdf.Mil, zOfs}, // 118 mil
		{3693.46 * sdf.Mil, 610.15 * sdf.Mil, zOfs},  // 118 mil
	}

	s, _ := obj.Standoff3D(k)
	s0 := sdf.Multi3D(s, positions0)

	return s0, nil
}

//-----------------------------------------------------------------------------

func ap745mount() (sdf.SDF3, error) {

	baseX := 180.0
	baseY := 200.0

	pcbX := 6102.36 * sdf.Mil
	pcbY := 7047.24 * sdf.Mil

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "xx", "xx", "xx"},
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	c1 := sdf.Box2D(v2.Vec{140, 50.0}, 3.0)
	c1 = sdf.Transform2D(c1, sdf.Translate2d(v2.Vec{0, -37}))

	c2 := sdf.Box2D(v2.Vec{90.0, 50.0}, 3.0)
	c2 = sdf.Transform2D(c2, sdf.Translate2d(v2.Vec{15, 45}))

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, sdf.Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := 0.5 * pcbY
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// add the standoffs
	s3, err := ap745Standoffs()
	if err != nil {
		return nil, err
	}

	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
}

//-----------------------------------------------------------------------------

func main() {

	ap745mount, err := ap745mount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(ap745mount, shrink), "ap745.stl", render.NewMarchingCubesOctree(500))

}

//-----------------------------------------------------------------------------
