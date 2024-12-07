//-----------------------------------------------------------------------------
/*

Create Eurorack Module Panels

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
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const panelThickness = 2.5 // mm

//-----------------------------------------------------------------------------

func standoff(h float64) (sdf.SDF3, error) {
	// standoff with screw hole
	k := &obj.StandoffParms{
		PillarHeight:   h,
		PillarDiameter: 8,
		HoleDepth:      10,
		HoleDiameter:   2.4, // #4 screw
	}
	return obj.Standoff3D(k)
}

// halfBreadBoardStandoffs returns the standoffs for an adafruit 1/2 breadboard.
func halfBreadBoardStandoffs(h float64) (sdf.SDF3, error) {
	s, err := standoff(h)
	if err != nil {
		return nil, err
	}
	positions := v3.VecSet{
		{0, -1450 * sdf.Mil, 0},
		{0, 1450 * sdf.Mil, 0},
	}
	return sdf.Multi3D(s, positions), nil
}

//-----------------------------------------------------------------------------
// panel holes and/or indents for mounted components

// pot0 return the panel hole/indent for a potentiometer
func pot0() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  9.4,
		Thickness: panelThickness,
		Indent:    v3.Vec{2, 4, 2},
		Offset:    11.0,
		//Orientation: sdf.DtoR(0),
	}
	return obj.PanelHole3D(&k)
}

// pot1 return the panel hole/indent for a potentiometer
func pot1() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  7.2,
		Thickness: panelThickness,
		Indent:    v3.Vec{2, 2, 1.5},
		Offset:    7.0,
		//Orientation: sdf.DtoR(0),
	}
	return obj.PanelHole3D(&k)
}

// spdt return the panel hole/indent for a spdt switch
func spdt() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  6.2,
		Thickness: panelThickness,
		Indent:    v3.Vec{2, 2, 1.5},
		Offset:    5.4,
		//Orientation: sdf.DtoR(0),
	}
	return obj.PanelHole3D(&k)
}

// led returns the panel hole for an led bezel
func led() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  7.0,
		Thickness: panelThickness,
	}
	return obj.PanelHole3D(&k)
}

// jack35 returns the panel hole/indent for a 3.5 mm audio jack
func jack35() (sdf.SDF3, error) {
	k := obj.PanelHoleParms{
		Diameter:  6.4,
		Thickness: panelThickness,
		Indent:    v3.Vec{2, 2, 1.5},
		Offset:    4.9,
		//Orientation: sdf.DtoR(0),
	}
	return obj.PanelHole3D(&k)
}

//-----------------------------------------------------------------------------

// powerBoardMount returns a pcb mount for a SynthRotek Noise Filtering Power Distribution Board.
func powerBoardMount() (sdf.SDF3, error) {

	const baseThickness = 3
	const standoffHeight = 10
	const xSpace = 0.9 * sdf.MillimetresPerInch
	const ySpace = 1.1 * sdf.MillimetresPerInch

	// standoffs
	s0, err := standoff(standoffHeight)
	if err != nil {
		return nil, err
	}
	// 4x2 sections
	const zOfs = (baseThickness + standoffHeight) * 0.5
	positions := v3.VecSet{
		{-1.5 * xSpace, -0.5 * ySpace, zOfs},
		{-1.5 * xSpace, 0.5 * ySpace, zOfs},
		{-0.5 * xSpace, -0.5 * ySpace, zOfs},
		{-0.5 * xSpace, 0.5 * ySpace, zOfs},
		{0.5 * xSpace, -0.5 * ySpace, zOfs},
		{0.5 * xSpace, 0.5 * ySpace, zOfs},
		{1.5 * xSpace, -0.5 * ySpace, zOfs},
		{1.5 * xSpace, 0.5 * ySpace, zOfs},
	}
	s1 := sdf.Multi3D(s0, positions)

	// base
	const baseX = (4 - 0.1) * xSpace
	const baseY = 2.0 * ySpace
	k := obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	s2, err := obj.Panel2D(&k)
	if err != nil {
		return nil, err
	}

	// cutout
	c0 := sdf.Box2D(v2.Vec{3 * xSpace, 0.5 * ySpace}, 3.0)
	s3 := sdf.Extrude3D(sdf.Difference2D(s2, c0), baseThickness)

	s4 := sdf.Union3D(s3, s1)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4, nil
}

//-----------------------------------------------------------------------------

var psuSize = v3.Vec{98, 129, 38}

// rt65b returns a model of a Meanwell RT-65B PSU
func rt65b() (sdf.SDF3, error) {

	body, err := sdf.Box3D(psuSize, 0)
	if err != nil {
		return nil, err
	}
	// move to the positive octant
	body = sdf.Transform3D(body, sdf.Translate3d(psuSize.MulScalar(0.5)))

	// m3 screw holes
	s0, err := obj.CounterBoredHole3D(12, 3.8*0.5, 10.6*0.5, 3.5)
	if err != nil {
		return nil, err
	}

	// vertical screws
	vs0 := sdf.Transform3D(s0, sdf.RotateY(sdf.DtoR(180)))
	vs0 = sdf.Transform3D(vs0, sdf.Translate3d(v3.Vec{31, 4.5 + 73.5, 0}))
	vs1 := sdf.Transform3D(vs0, sdf.Translate3d(v3.Vec{33, 0, 0}))

	// horizontal screws
	hs0 := sdf.Transform3D(s0, sdf.RotateY(sdf.DtoR(90)))
	hs0 = sdf.Transform3D(hs0, sdf.Translate3d(v3.Vec{psuSize.X, 32, 38 - 18.5}))
	hs1 := sdf.Transform3D(hs0, sdf.Translate3d(v3.Vec{0, 77, 9}))
	hs2 := sdf.Transform3D(hs0, sdf.Translate3d(v3.Vec{0, 77, -9}))

	psu := sdf.Union3D(body, vs0, vs1, hs0, hs1, hs2)

	return psu, nil
}

// psuMount returns a mount for a Meanwell RT-65B PSU
func psuMount() (sdf.SDF3, error) {

	// base
	const baseThickness = 6
	baseSize := v2.Vec{135, 145}
	k0 := obj.PanelParms{
		Size:         baseSize,
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    baseThickness,
	}
	base, err := obj.Panel3D(&k0)
	if err != nil {
		return nil, err
	}

	// cutout 0
	k2 := obj.PanelParms{
		Size:         v2.Vec{90, 55},
		CornerRadius: 4.0,
		Thickness:    baseThickness,
	}
	c0, err := obj.Panel3D(&k2)
	if err != nil {
		return nil, err
	}
	c0 = sdf.Transform3D(c0, sdf.Translate3d(v3.Vec{-10, -27.5, 0}))

	// cutout 1
	k3 := obj.PanelParms{
		Size:         v2.Vec{90, 30},
		CornerRadius: 4.0,
		Thickness:    baseThickness,
	}
	c1, err := obj.Panel3D(&k3)
	if err != nil {
		return nil, err
	}
	c1 = sdf.Transform3D(c1, sdf.Translate3d(v3.Vec{-10, 40, 0}))

	// upright mount
	uprightSize := v2.Vec{psuSize.Z + baseThickness*0.5, baseSize.Y}
	k1 := obj.PanelParms{
		Size:         uprightSize,
		CornerRadius: 3.0,
		Thickness:    baseThickness,
	}
	upright, err := obj.Panel3D(&k1)
	if err != nil {
		return nil, err
	}
	upright = sdf.Transform3D(upright, sdf.RotateY(sdf.DtoR(90)))
	uprightOffset := v3.Vec{psuSize.X + baseThickness, 0, uprightSize.X}.MulScalar(0.5)
	upright = sdf.Transform3D(upright, sdf.Translate3d(uprightOffset))

	psu, err := rt65b()
	if err != nil {
		return nil, err
	}
	psuOffset := v3.Vec{-psuSize.X, -psuSize.Y, baseThickness}.MulScalar(0.5)
	psu = sdf.Transform3D(psu, sdf.Translate3d(psuOffset))

	mount := sdf.Difference3D(sdf.Union3D(base, upright), sdf.Union3D(psu, c0, c1))

	return mount, nil
}

//-----------------------------------------------------------------------------

// powerPanel returns a mounting panel for a ac-14-f16a power connector.
func powerPanel() (sdf.SDF3, error) {

	const baseThickness = 4

	k := obj.PanelParms{
		Size:         v2.Vec{85, 95},
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}

	s, err := obj.Panel2D(&k)
	if err != nil {
		return nil, err
	}

	// panel cutout
	c0 := sdf.Box2D(v2.Vec{28, 48}, 3)

	// mounting holes
	hole, err := sdf.Circle2D(0.5 * 4.5)
	if err != nil {
		return nil, err
	}
	c1 := sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{20, 0}))
	c2 := sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{-20, 0}))

	cutouts := sdf.Union2D(c0, c1, c2)

	return sdf.Extrude3D(sdf.Difference2D(s, cutouts), baseThickness), nil
}

// powerPanelRouting returns a routing pattern for the power panel.
func powerPanelRouting() (sdf.SDF3, error) {

	const baseThickness = 4

	k := obj.PanelParms{
		Size:         v2.Vec{85, 95},
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}

	s, err := obj.Panel2D(&k)
	if err != nil {
		return nil, err
	}

	// panel cutout
	c := sdf.Box2D(v2.Vec{55, 65}, 3)

	return sdf.Extrude3D(sdf.Difference2D(s, c), baseThickness), nil
}

//-----------------------------------------------------------------------------

// arPanel returns the panel for an attack/release module.
func arPanel() (sdf.SDF3, error) {

	// 3u x 12hp panel
	k := obj.EuroRackParms{
		U:            3,
		HP:           12,
		CornerRadius: 3,
		HoleDiameter: 3.6,
		Thickness:    panelThickness,
		Ridge:        true,
	}
	s, err := obj.EuroRackPanel3D(&k)
	if err != nil {
		return nil, err
	}

	// breadboard standoffs
	const standoffHeight = 25
	so, err := halfBreadBoardStandoffs(standoffHeight)
	if err != nil {
		return nil, err
	}
	so = sdf.Transform3D(so, sdf.Translate3d(v3.Vec{0, 3, (panelThickness + standoffHeight) * 0.5}))
	s = sdf.Union3D(s, so)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(2))

	// push button
	pb, err := sdf.Box3D(v3.Vec{13.2, 10.8, panelThickness}, 0)
	if err != nil {
		return nil, err
	}
	pb = sdf.Transform3D(pb, sdf.Translate3d(v3.Vec{0, 0, 0}))

	// cv input/output
	cv, err := jack35()
	if err != nil {
		return nil, err
	}
	cv0 := sdf.Transform3D(cv, sdf.Translate3d(v3.Vec{-20, -45, 0}))
	cv1 := sdf.Transform3D(cv, sdf.Translate3d(v3.Vec{20, -45, 0}))

	// LED
	led, err := led()
	if err != nil {
		return nil, err
	}
	led = sdf.Transform3D(led, sdf.Translate3d(v3.Vec{0, -45, 0}))

	// attack/release pots
	pot, err := pot0()
	if err != nil {
		return nil, err
	}
	pot0 := sdf.Transform3D(pot, sdf.Translate3d(v3.Vec{-15, 25, 0}))
	pot1 := sdf.Transform3D(pot, sdf.Translate3d(v3.Vec{15, 25, 0}))

	// spdt switch
	spdt, err := spdt()
	if err != nil {
		return nil, err
	}
	spdt = sdf.Transform3D(spdt, sdf.Translate3d(v3.Vec{0, -22, 0}))

	return sdf.Difference3D(s, sdf.Union3D(pb, cv0, cv1, led, pot0, pot1, spdt)), nil
}

//-----------------------------------------------------------------------------

// bbPanel returns a panel for mounting a half bread board.
func bbPanel() (sdf.SDF3, error) {

	// 3u x 12hp panel
	k := obj.EuroRackParms{
		U:            3,
		HP:           12,
		CornerRadius: 3,
		HoleDiameter: 3.6,
		Thickness:    panelThickness,
		Ridge:        true,
	}
	s, err := obj.EuroRackPanel3D(&k)
	if err != nil {
		return nil, err
	}

	// breadboard standoffs
	const standoffHeight = 12
	so, err := halfBreadBoardStandoffs(standoffHeight)
	if err != nil {
		return nil, err
	}
	so = sdf.Transform3D(so, sdf.Translate3d(v3.Vec{0, 3, (panelThickness + standoffHeight) * 0.5}))
	s = sdf.Union3D(s, so)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(2))

	return s, nil
}

//-----------------------------------------------------------------------------

func main() {

	p, err := arPanel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(p, shrink), "ar_panel.stl", render.NewMarchingCubesOctree(300))

	p, err = powerBoardMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(p, shrink), "pwr_mount.stl", render.NewMarchingCubesOctree(300))

	p, err = powerPanel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(p, shrink), "pwr_panel.stl", render.NewMarchingCubesOctree(300))

	p, err = powerPanelRouting()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(p, shrink), "pwr_panel_routing.stl", render.NewMarchingCubesOctree(300))

	p, err = psuMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(p, shrink), "psu_mount.stl", render.NewMarchingCubesOctree(300))

	p, err = bbPanel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(p, shrink), "bb_panel.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
