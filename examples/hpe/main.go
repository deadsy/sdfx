//-----------------------------------------------------------------------------
/*

HPE AP-745 and AP-725 Mounting Boards

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

var baseThickness = 3.0
var pillarHeight = 15.0

//-----------------------------------------------------------------------------

func ap723hStandoffs() (sdf.SDF3, error) {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 9.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}

	positions0 := v3.VecSet{
		{0, 0, zOfs},
		{103.0, 0, zOfs},
		{103.0, 152.0, zOfs},
		{0, 152.0, zOfs},
	}

	s, _ := obj.Standoff3D(k)
	s0 := sdf.Multi3D(s, positions0)

	return s0, nil
}

func ap723hMount() (sdf.SDF3, error) {

	baseX := 120.0
	baseY := 175.0

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	var c1, c2 sdf.SDF2

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, sdf.Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * baseX
	yOfs := 0.5 * baseY
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// add the standoffs
	s3, err := ap723hStandoffs()
	if err != nil {
		return nil, err
	}

	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
}

//-----------------------------------------------------------------------------

func ap725Standoffs() (sdf.SDF3, error) {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 138.0 * sdf.Mil * 2.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}

	positions0 := v3.VecSet{
		{0, 0, zOfs},
		{5984.255 * sdf.Mil, 0, zOfs},
		{5551.185 * sdf.Mil, 4704.72 * sdf.Mil, zOfs},
		{433.071 * sdf.Mil, 4704.72 * sdf.Mil, zOfs},
		{2700.795 * sdf.Mil, 5389.76 * sdf.Mil, zOfs},
		{3714.565 * sdf.Mil, 1708.66 * sdf.Mil, zOfs},
	}

	s, _ := obj.Standoff3D(k)
	s0 := sdf.Multi3D(s, positions0)

	return s0, nil
}

func ap725Mount() (sdf.SDF3, error) {

	baseX := 165.0
	baseY := 150.0

	pcbX := 5984.255 * sdf.Mil
	pcbY := 5389.76 * sdf.Mil

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		//HolePattern:  [4]string{"xx", "xxx", ".xx", ".xx"},
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	// cutouts
	c1 := sdf.Box2D(v2.Vec{100, 65.0}, 3.0)
	c1 = sdf.Transform2D(c1, sdf.Translate2d(v2.Vec{0, 20}))

	c2 := sdf.Box2D(v2.Vec{135, 25.0}, 3.0)
	c2 = sdf.Transform2D(c2, sdf.Translate2d(v2.Vec{0, -50}))

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, sdf.Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := 0.5 * pcbY
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// reinforcing ribs
	const ribHeight = 5.0
	r0, _ := sdf.Box3D(v3.Vec{3.0, 0.9 * pcbY, ribHeight}, 0)
	yOfs = 0.5 * pcbY
	xOfs = pcbX
	zOfs := 0.5 * (ribHeight + baseThickness)
	r0 = sdf.Transform3D(r0, sdf.Translate3d(v3.Vec{0, yOfs, zOfs}))
	r1 := sdf.Transform3D(r0, sdf.Translate3d(v3.Vec{xOfs, 0, 0}))

	r2, _ := sdf.Box3D(v3.Vec{0.9 * pcbX, 3.0, ribHeight}, 0)
	xOfs = 0.5 * pcbX
	r2 = sdf.Transform3D(r2, sdf.Translate3d(v3.Vec{xOfs, 0, zOfs}))

	s2 = sdf.Union3D(s2, r0, r1, r2)

	// add the standoffs
	s3, err := ap725Standoffs()
	if err != nil {
		return nil, err
	}

	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
}

//-----------------------------------------------------------------------------

const holeSquare = 6102.36 * sdf.Mil

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
		{0, holeSquare, zOfs},          // 138 mil
		{0, 0, zOfs},                   // 138 mil
		{holeSquare, holeSquare, zOfs}, // 138 mil
		{holeSquare, 0, zOfs},          // 138 mil

		{3937.01 * sdf.Mil, 7047.24 * sdf.Mil, zOfs}, // 118 mil
		{1240.16 * sdf.Mil, 5570.87 * sdf.Mil, zOfs}, // 118 mil
		{2648.46 * sdf.Mil, 3485.15 * sdf.Mil, zOfs}, // 118 mil
		{3693.46 * sdf.Mil, 610.15 * sdf.Mil, zOfs},  // 118 mil
	}

	s, _ := obj.Standoff3D(k)
	s0 := sdf.Multi3D(s, positions0)

	return s0, nil
}

func ap745Mount() (sdf.SDF3, error) {

	baseX := 170.0
	baseY := 190.0

	pcbX := 6102.36 * sdf.Mil
	pcbY := 7047.24 * sdf.Mil

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		//HolePattern:  [4]string{"xx", "xxx", ".xx", ".xx"},
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	// cutouts
	c1 := sdf.Box2D(v2.Vec{140, 50.0}, 3.0)
	c1 = sdf.Transform2D(c1, sdf.Translate2d(v2.Vec{0, -37}))
	c2 := sdf.Box2D(v2.Vec{90.0, 50.0}, 3.0)
	c2 = sdf.Transform2D(c2, sdf.Translate2d(v2.Vec{15, 45}))

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, sdf.Union2D(c1, c2)), baseThickness)
	xOfs := 0.5 * pcbX
	yOfs := 0.5 * pcbY
	s2 = sdf.Transform3D(s2, sdf.Translate3d(v3.Vec{xOfs, yOfs, 0}))

	// reinforcing ribs
	const ribHeight = 5.0
	r0, _ := sdf.Box3D(v3.Vec{3.0, 0.75 * pcbY, ribHeight}, 0)
	yOfs = 0.5*pcbY - 12.0
	zOfs := 0.5 * (ribHeight + baseThickness)
	r0 = sdf.Transform3D(r0, sdf.Translate3d(v3.Vec{0, yOfs, zOfs}))
	r1 := sdf.Transform3D(r0, sdf.Translate3d(v3.Vec{holeSquare, 0, 0}))
	s2 = sdf.Union3D(s2, r0, r1)

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

	ap725, err := ap725Mount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(ap725, "ap725.stl", render.NewMarchingCubesOctree(500))

	ap745, err := ap745Mount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(ap745, "ap745.stl", render.NewMarchingCubesOctree(500))

	ap723h, err := ap723hMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(ap723h, "ap723h.stl", render.NewMarchingCubesOctree(500))

}

//-----------------------------------------------------------------------------
