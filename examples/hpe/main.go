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

func ap723hSupport() (sdf.SDF3, error) {

	const w0 = 20
	const l0 = 60
	const h0 = 6

	b0, _ := sdf.Box3D(v3.Vec{w0, l0, h0}, 0)

	const h1 = 3.7
	const l1 = 47

	b1, _ := sdf.Box3D(v3.Vec{w0, l1, h1}, 0)
	zOfs := 0.5 * (h1 - h0)
	yOfs := 0.5 * (l1 - l0)
	b1 = sdf.Transform3D(b1, sdf.Translate3d(v3.Vec{0, yOfs, zOfs}))

	hole, _ := sdf.Cylinder3D(h0-h1, 1.2, 0)
	xOfs := 0.5*w0 - 3.0
	yOfs = l1 - 0.5*l0
	zOfs = 0.5 * h1
	hole = sdf.Transform3D(hole, sdf.Translate3d(v3.Vec{xOfs, yOfs, zOfs}))

	s := sdf.Difference3D(b0, sdf.Union3D(b1, hole))

	return s, nil
}

func ap723hStandoffs() (sdf.SDF3, error) {

	// standoffs with screw holes
	k0 := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 10.0,
		HoleDepth:      10.0,
		HoleDiameter:   4.0,
	}

	k1 := &obj.StandoffParms{
		PillarHeight:   pillarHeight + 2.0,
		PillarDiameter: 5.5,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}

	s0, _ := obj.Standoff3D(k0)
	s1, _ := obj.Standoff3D(k1)
	s := sdf.Union3D(s0, s1)

	zOfs := 0.5 * (pillarHeight + baseThickness)

	positions0 := v3.VecSet{
		{0, 0, zOfs},
		{103.0, 0, zOfs},
		{103.0, 152.0, zOfs},
		{0, 152.0, zOfs},
	}

	return sdf.Multi3D(s, positions0), nil
}

func ap723hMount() (sdf.SDF3, error) {

	pcbX := 102.5
	pcbY := 152.0

	baseX := pcbX + 20.0
	baseY := pcbY + 20.0

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	// cutouts
	c1 := sdf.Box2D(v2.Vec{baseX - 35, baseY - 35}, 5.0)

	// extrude the base
	s2 := sdf.Extrude3D(sdf.Difference2D(s0, c1), baseThickness)
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

	r2, _ := sdf.Box3D(v3.Vec{0.8 * pcbX, 3.0, ribHeight}, 0)
	yOfs = pcbY
	xOfs = 0.5 * pcbX
	r2 = sdf.Transform3D(r2, sdf.Translate3d(v3.Vec{xOfs, 0, zOfs}))
	r3 := sdf.Transform3D(r2, sdf.Translate3d(v3.Vec{0, yOfs, 0}))

	s2 = sdf.Union3D(s2, r0, r1, r2, r3)

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

	support, err := ap723hSupport()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(support, "ap723h_support.stl", render.NewMarchingCubesOctree(500))

}

//-----------------------------------------------------------------------------
