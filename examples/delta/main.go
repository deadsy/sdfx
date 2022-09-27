//-----------------------------------------------------------------------------
/*

Delta Robot Parts

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const upperArmWidth = 30.0
const upperArmRadius0 = 15.0
const upperArmRadius1 = 5.0
const upperArmRadius2 = 3.9 * 0.5
const upperArmLength = 100.0

func upperArm() (sdf.SDF3, error) {

	const upperArmThickness = 5.0 * 2.0
	const gussetThickness = 0.5

	// body
	b, err := sdf.FlatFlankCam2D(upperArmLength, upperArmRadius0, upperArmRadius1)
	if err != nil {
		return nil, err
	}
	body := sdf.Extrude3D(b, upperArmThickness)

	// end cylinder
	c0, err := sdf.Cylinder3D(upperArmWidth*2.0, upperArmRadius1, 0)
	if err != nil {
		return nil, err
	}
	c0 = sdf.Transform3D(c0, sdf.Translate3d(v3.Vec{0, upperArmLength, 0}))

	// end cylinder hole
	c1, err := sdf.Cylinder3D(upperArmWidth*2.0, upperArmRadius2, 0)
	if err != nil {
		return nil, err
	}
	c1 = sdf.Transform3D(c1, sdf.Translate3d(v3.Vec{0, upperArmLength, 0}))

	// gusset
	const dx = upperArmWidth * 2.0 * 0.4
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
	gusset = sdf.Transform3D(gusset, sdf.Translate3d(v3.Vec{0, yOfs, 0}))

	// servo mounting
	k := obj.ServoHornParms{
		CenterRadius: 3,
		NumHoles:     4,
		CircleRadius: 14 * 0.5,
		HoleRadius:   1.9,
	}
	h0, err := obj.ServoHorn(&k)
	if err != nil {
		return nil, err
	}
	horn := sdf.Extrude3D(h0, upperArmThickness)

	const hornRadius = 10
	const hornThickness = 2.3
	hornBody, err := sdf.Cylinder3D(hornThickness, hornRadius, 0)
	if err != nil {
		return nil, err
	}
	zOfs := (upperArmThickness - hornThickness) * 0.5
	hornBody = sdf.Transform3D(hornBody, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	// body + cylinder
	s := sdf.Union3D(body, c0)
	// add the gusset with fillets
	s = sdf.Union3D(s, gusset)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(upperArmThickness * gussetThickness))
	// remove the holes
	s = sdf.Difference3D(s, sdf.Union3D(c1, horn, hornBody))

	// cut in half
	s = sdf.Cut3D(s, v3.Vec{}, v3.Vec{0, 0, 1})

	return s, nil
}

//-----------------------------------------------------------------------------

const servoMountUprightLength = 66.0
const servoMountBaseLength = 35.0
const servoMountThickness = 3.5
const servoMountWidth = 35.0
const servoMountHoleRadius = 2.4

func servoMountHoles(h float64) (sdf.SDF3, error) {
	// base holes
	hole, err := sdf.Cylinder3D(h, servoMountHoleRadius, 0)
	if err != nil {
		return nil, err
	}
	hole = sdf.Transform3D(hole, sdf.Translate3d(v3.Vec{(servoMountBaseLength + servoMountThickness) * 0.5, 0, 0}))
	dx := (servoMountBaseLength * 0.5) - servoMountThickness - 4.0
	dy := (servoMountWidth * 0.5) - servoMountThickness - 6.0
	holes := sdf.Multi3D(hole, []v3.Vec{{dx, dy, 0}, {-dx, dy, 0}, {dx, -dy, 0}, {-dx, -dy, 0}})
	return holes, nil
}

func servoMount() (sdf.SDF3, error) {

	const servoOffset = servoMountUprightLength - 20.0

	m := sdf.NewPolygon()
	m.Add(0, 0)
	m.Add(servoMountBaseLength, 0)
	m.Add(servoMountBaseLength, servoMountThickness)
	m.Add(servoMountThickness, servoMountUprightLength)
	m.Add(0, servoMountUprightLength)
	m2d, err := sdf.Polygon2D(m.Vertices())
	if err != nil {
		return nil, err
	}
	mount := sdf.Extrude3D(m2d, servoMountWidth)

	// cavity
	c := sdf.NewPolygon()
	c.Add(servoMountThickness, servoMountThickness)
	c.Add(servoMountBaseLength, servoMountThickness)
	c.Add(servoMountThickness, servoMountUprightLength)
	c2d, err := sdf.Polygon2D(c.Vertices())
	cavity := sdf.Extrude3D(c2d, servoMountWidth-2*servoMountThickness)

	mount = sdf.Difference3D(mount, cavity)
	mount = sdf.Transform3D(mount, sdf.RotateX(sdf.DtoR(90)))

	// base holes
	holes, err := servoMountHoles(servoMountThickness)
	if err != nil {
		return nil, err
	}
	holes = sdf.Transform3D(holes, sdf.Translate3d(v3.Vec{0, 0, servoMountThickness * 0.5}))
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
	servo := sdf.Extrude3D(servo2d, servoMountThickness)
	servo = sdf.Transform3D(servo, sdf.RotateY(sdf.DtoR(90)))
	servo = sdf.Transform3D(servo, sdf.Translate3d(v3.Vec{servoMountThickness * 0.5, 0, servoOffset}))

	s := sdf.Difference3D(mount, servo)

	return s, nil
}

//-----------------------------------------------------------------------------

const baseSide = 150
const baseThickness = 7
const basePillarHeight = 20
const baseHoleRadius = 7

var servoY = -baseSide * math.Tan(sdf.DtoR(30)) * 0.5
var servoX = 25.0 - upperArmWidth*0.5

func deltaBase() (sdf.SDF3, error) {

	// servo holes
	holes, err := servoMountHoles(baseThickness)
	if err != nil {
		return nil, err
	}
	holes = sdf.Transform3D(holes, sdf.Translate3d(v3.Vec{servoX, servoY, -baseThickness * 0.5}))
	holes = sdf.RotateUnion3D(holes, 3, sdf.RotateZ(sdf.DtoR(120)))

	// base
	base, err := sdf.Cylinder3D(baseThickness, baseSide*0.5*1.05, 0)
	if err != nil {
		return nil, err
	}
	base = sdf.Transform3D(base, sdf.Translate3d(v3.Vec{0, 0, -baseThickness * 0.5}))

	// pillars
	k := obj.StandoffParms{
		PillarHeight:   basePillarHeight,
		PillarDiameter: 15,
		HoleDepth:      15,
		HoleDiameter:   3,
	}
	pillars, err := obj.Standoff3D(&k)
	if err != nil {
		return nil, err
	}
	pillars = sdf.Transform3D(pillars, sdf.RotateX(sdf.DtoR(180)))
	pillars = sdf.Transform3D(pillars, sdf.Translate3d(v3.Vec{0, -baseSide * 0.4, -(0.5*basePillarHeight + baseThickness)}))
	pillars = sdf.RotateUnion3D(pillars, 3, sdf.RotateZ(sdf.DtoR(120)))

	// hole for servo wires
	baseHole, err := sdf.Cylinder3D(baseThickness, baseHoleRadius, 0)
	if err != nil {
		return nil, err
	}
	baseHole = sdf.Transform3D(baseHole, sdf.Translate3d(v3.Vec{0, 0, -baseThickness * 0.5}))
	holes = sdf.Union3D(holes, baseHole)

	// base/pillar fillet
	base = sdf.Union3D(base, pillars)
	base.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(baseThickness))

	return sdf.Difference3D(base, holes), nil
}

//-----------------------------------------------------------------------------

const rodRadius = (6.0 + 0.5) * 0.5
const holderRadius = (11.7 + 0.2) * 0.5
const holderHeight = 2.6

func rodEnd() (sdf.SDF3, error) {

	const endRadius = holderRadius * 1.5
	const endHeight = rodRadius * 2.0 * 1.5
	const round = endHeight * 0.1
	end, err := sdf.Cylinder3D(endHeight, endRadius, round)
	if err != nil {
		return nil, err
	}

	const endX = 3 * endHeight
	box, err := sdf.Box3D(v3.Vec{endX, endHeight, endHeight}, round)
	if err != nil {
		return nil, err
	}
	box = sdf.Transform3D(box, sdf.Translate3d(v3.Vec{0.5 * endX, 0, 0}))

	const rodHole = (endX - endRadius) * 0.9
	const ofsX = endX - 0.5*rodHole
	rod, err := sdf.Cylinder3D(rodHole, rodRadius, 0)
	if err != nil {
		return nil, err
	}
	rod = sdf.Transform3D(rod, sdf.RotateY(sdf.DtoR(90)))
	rod = sdf.Transform3D(rod, sdf.Translate3d(v3.Vec{ofsX, 0, 0}))

	holder, err := sdf.Cylinder3D(holderHeight, holderRadius, 0)
	if err != nil {
		return nil, err
	}
	const ofsZ = (endHeight - holderHeight) * 0.5
	holder = sdf.Transform3D(holder, sdf.Translate3d(v3.Vec{0, 0, ofsZ}))

	// end + box with fillets
	s := sdf.Union3D(end, box)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(endRadius * 0.1))

	// bump removal
	s = sdf.Cut3D(s, v3.Vec{0, 0, -endHeight * 0.5}, v3.Vec{0, 0, 1})
	s = sdf.Cut3D(s, v3.Vec{0, 0, endHeight * 0.5}, v3.Vec{0, 0, -1})

	// remove the cavities
	s = sdf.Difference3D(s, sdf.Union3D(holder, rod))

	return s, nil
}

//-----------------------------------------------------------------------------

const platformSide = 50
const platformThickness = 10.0

func platform() (sdf.SDF3, error) {

	pHalf := platformSide * 0.5
	pShort := pHalf / math.Sqrt(3)
	pLong := 2 * pShort

	c0 := v3.Vec{0, -pLong, 0}
	c1 := v3.Vec{pHalf, pShort, 0}
	c2 := v3.Vec{-pHalf, pShort, 0}

	// platform
	pp := sdf.NewPolygon()
	pp.Add(c0.X, c0.Y)
	pp.Add(c1.X, c1.Y)
	pp.Add(c2.X, c2.Y)
	p2d, err := sdf.Polygon2D(pp.Vertices())
	if err != nil {
		return nil, err
	}
	platform := sdf.Extrude3D(p2d, platformThickness)

	// connection arms

	arm0, err := obj.Pipe3D(platformThickness*0.5, upperArmRadius2, upperArmWidth)
	if err != nil {
		return nil, err
	}
	arm0 = sdf.Transform3D(arm0, sdf.RotateY(sdf.DtoR(90)))
	arm0 = sdf.Transform3D(arm0, sdf.Translate3d(c0))

	arm1 := sdf.Transform3D(arm0, sdf.RotateZ(sdf.DtoR(120)))
	arm2 := sdf.Transform3D(arm0, sdf.RotateZ(sdf.DtoR(-120)))

	s := sdf.Union3D(platform, arm0, arm1, arm2)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(platformThickness * 0.7))

	// bump removal
	s = sdf.Cut3D(s, v3.Vec{0, 0, -platformThickness * 0.5}, v3.Vec{0, 0, 1})
	s = sdf.Cut3D(s, v3.Vec{0, 0, platformThickness * 0.5}, v3.Vec{0, 0, -1})

	return s, nil
}

//-----------------------------------------------------------------------------

func baseWithServos() (sdf.SDF3, error) {

	// servos
	servos, err := servoMount()
	if err != nil {
		return nil, err
	}
	servos = sdf.Transform3D(servos, sdf.Translate3d(v3.Vec{servoX, servoY, 0}))
	servos = sdf.RotateUnion3D(servos, 3, sdf.RotateZ(sdf.DtoR(120)))

	// base
	base, err := deltaBase()
	if err != nil {
		return nil, err
	}

	return sdf.Union3D(base, servos), nil
}

//-----------------------------------------------------------------------------

func main() {

	s, err := upperArm()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, "arm.stl", render.NewMarchingCubesOctree(500))

	s, err = servoMount()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, "servomount.stl", render.NewMarchingCubesOctree(250))

	s, err = deltaBase()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, "base.stl", render.NewMarchingCubesOctree(300))

	s, err = rodEnd()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, "rodend.stl", render.NewMarchingCubesOctree(100))

	s, err = platform()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, "platform.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
