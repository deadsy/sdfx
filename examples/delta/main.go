//-----------------------------------------------------------------------------
/*

Delta Robot Parts

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
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func upperArm() (sdf.SDF3, error) {

	const upperArmRadius0 = 16.0
	const upperArmRadius1 = 5.0
	const upperArmRadius2 = 2.5
	const upperArmLength = 120.0
	const upperArmThickness = 5.0
	const upperArmWidth = 50.0
	const gussetThickness = 0.7

	// body
	b, err := sdf.FlatFlankCam2D(upperArmLength, upperArmRadius0, upperArmRadius1)
	if err != nil {
		return nil, err
	}
	body := sdf.Extrude3D(b, upperArmThickness)

	// end cylinder
	c0, err := sdf.Cylinder3D(upperArmWidth, upperArmRadius1, 0)
	if err != nil {
		return nil, err
	}
	c0 = sdf.Transform3D(c0, sdf.Translate3d(sdf.V3{0, upperArmLength, 0}))

	// end cylinder hole
	c1, err := sdf.Cylinder3D(upperArmWidth, upperArmRadius2, 0)
	if err != nil {
		return nil, err
	}
	c1 = sdf.Transform3D(c1, sdf.Translate3d(sdf.V3{0, upperArmLength, 0}))

	// gusset
	const dx = upperArmWidth * 0.4
	const dy = upperArmLength * 0.6
	g := sdf.NewPolygon()
	g.Add(-dx, dy)
	g.Add(dx, dy)
	g.Add(0, 0)
	g2d, err := sdf.Polygon2D(g.Vertices())
	if err != nil {
		return nil, err
	}
	gusset := sdf.Extrude3D(g2d, upperArmThickness*gussetThickness)
	gusset = sdf.Transform3D(gusset, sdf.RotateY(sdf.DtoR(90)))
	yOfs := upperArmLength - dy
	gusset = sdf.Transform3D(gusset, sdf.Translate3d(sdf.V3{0, yOfs, 0}))

	// servo mounting
	k := obj.ServoHornParms{
		CenterRadius: 4,
		NumHoles:     6,
		CircleRadius: 10,
		HoleRadius:   1,
	}
	h0, err := obj.ServoHorn(&k)
	if err != nil {
		return nil, err
	}
	horn := sdf.Extrude3D(h0, upperArmThickness)

	// body + cylinder
	s := sdf.Union3D(body, c0)
	// add the gusset with fillets
	s = sdf.Union3D(s, gusset)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(upperArmThickness * gussetThickness))
	// remove the holes
	s = sdf.Difference3D(s, sdf.Union3D(c1, horn))

	return s, nil
}

//-----------------------------------------------------------------------------

func servoMount() (sdf.SDF3, error) {

	const uprightLength = 66.0
	const baseLength = 35.0
	const thickness = 3.5
	const width = 35.0
	const servoOffset = uprightLength - 20.0
	const mountHoleRadius = 2.4

	m := sdf.NewPolygon()
	m.Add(0, 0)
	m.Add(baseLength, 0)
	m.Add(baseLength, thickness)
	m.Add(thickness, uprightLength)
	m.Add(0, uprightLength)
	m2d, err := sdf.Polygon2D(m.Vertices())
	if err != nil {
		return nil, err
	}
	mount := sdf.Extrude3D(m2d, width)

	// cavity
	c := sdf.NewPolygon()
	c.Add(thickness, thickness)
	c.Add(baseLength, thickness)
	c.Add(thickness, uprightLength)
	c2d, err := sdf.Polygon2D(c.Vertices())
	cavity := sdf.Extrude3D(c2d, width-2*thickness)

	mount = sdf.Difference3D(mount, cavity)
	mount = sdf.Transform3D(mount, sdf.RotateX(sdf.DtoR(90)))

	// base holes
	hole, err := sdf.Cylinder3D(thickness, mountHoleRadius, 0)
	hole = sdf.Transform3D(hole, sdf.Translate3d(sdf.V3{(baseLength + thickness) * 0.5, 0, thickness * 0.5}))
	dx := (baseLength * 0.5) - thickness - 4.0
	dy := (width * 0.5) - thickness - 6.0
	holes := sdf.Multi3D(hole, []sdf.V3{{dx, dy, 0}, {-dx, dy, 0}, {dx, -dy, 0}, {-dx, -dy, 0}})

	mount = sdf.Difference3D(mount, holes)

	// servo
	k, err := obj.ServoLookup("annimos_ds3218")
	if err != nil {
		return nil, err
	}
	servo2d, err := obj.Servo2D(k, 2.1)
	if err != nil {
		return nil, err
	}
	servo := sdf.Extrude3D(servo2d, thickness)
	servo = sdf.Transform3D(servo, sdf.RotateY(sdf.DtoR(90)))
	servo = sdf.Transform3D(servo, sdf.Translate3d(sdf.V3{thickness * 0.5, 0, servoOffset}))

	s := sdf.Difference3D(mount, servo)

	return s, nil
}

//-----------------------------------------------------------------------------

func servoControllerMount() (sdf.SDF3, error) {

	// standoff
	k0 := obj.StandoffParms{
		PillarHeight:   0.5 * sdf.MillimetresPerInch,
		PillarDiameter: 5,
		HoleDepth:      10,
		HoleDiameter:   2.4, // #4 screw
	}
	standoff, err := obj.Standoff3D(&k0)
	if err != nil {
		return nil, err
	}

	// standoffs
	h0 := sdf.V3{-0.45, -0.8, 0.25}.MulScalar(sdf.MillimetresPerInch)
	h1 := sdf.V3{0.05, 0.8, 0.25}.MulScalar(sdf.MillimetresPerInch)
	standoffs := sdf.Multi3D(standoff, []sdf.V3{h0, h1})

	// base
	k1 := obj.PanelParms{
		Size:         sdf.V2{1.1, 1.8}.MulScalar(sdf.MillimetresPerInch),
		CornerRadius: 2,
		HoleDiameter: 2.4, // #4 screw
		HoleMargin:   [4]float64{4, 4, 4, 4},
		HolePattern:  [4]string{"x", "x", ".x", ""},
		Thickness:    3,
	}
	base, err := obj.Panel3D(&k1)
	if err != nil {
		return nil, err
	}

	return sdf.Union3D(base, standoffs), nil
}

//-----------------------------------------------------------------------------

func main() {

	s, err := upperArm()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, 500, "arm.stl", &render.MarchingCubesOctree{})

	s, err = servoMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, 300, "servomount.stl", &render.MarchingCubesOctree{})

	s, err = servoControllerMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, 300, "ctrlmount.stl", &render.MarchingCubesOctree{})

}

//-----------------------------------------------------------------------------
