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

const upperArmRadius0 = 10.0
const upperArmRadius1 = 5.0
const upperArmRadius2 = 2.5
const upperArmLength = 100.0
const upperArmThickness = 5.0
const upperArmWidth = 50.0

const servoCenterHoleRadius = 2.0
const servoMountHoleRadius = 1.0
const servoMountRadius = 6.0
const servoMountHoles = 8

func upperArm() (sdf.SDF3, error) {

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
	dx := upperArmWidth * 0.4
	dy := upperArmLength * 0.6
	g := sdf.NewPolygon()
	g.Add(-dx, dy)
	g.Add(dx, dy)
	g.Add(0, 0)
	g2d, err := sdf.Polygon2D(g.Vertices())
	if err != nil {
		return nil, err
	}
	gusset := sdf.Extrude3D(g2d, upperArmThickness*0.5)
	gusset = sdf.Transform3D(gusset, sdf.RotateY(sdf.DtoR(90)))
	yOfs := upperArmLength - dy
	gusset = sdf.Transform3D(gusset, sdf.Translate3d(sdf.V3{0, yOfs, 0}))

	// servo mounting
	h0, err := obj.BoltCircle3D(upperArmThickness, servoMountHoleRadius, servoMountRadius, servoMountHoles)
	if err != nil {
		return nil, err
	}
	h1, err := sdf.Cylinder3D(upperArmWidth, servoCenterHoleRadius, 0)
	if err != nil {
		return nil, err
	}

	// body + cylinder
	s := sdf.Union3D(body, c0)
	// add the gusset with fillets
	s = sdf.Union3D(s, gusset)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(upperArmThickness * 0.5))
	// remove the holes
	s = sdf.Difference3D(s, sdf.Union3D(c1, h0, h1))

	return s, nil
}

//-----------------------------------------------------------------------------

func main() {

	s, err := upperArm()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, 300, "arm.stl", &render.MarchingCubesOctree{})
}

//-----------------------------------------------------------------------------
