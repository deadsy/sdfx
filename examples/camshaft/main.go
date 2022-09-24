//-----------------------------------------------------------------------------
/*

Wallaby Camshaft

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func camshaft() (sdf.SDF3, error) {

	// build the shaft from an SoR
	const l0 = 13.0 / 16.0
	const r0 = (5.0 / 16.0) / 2.0
	const l1 = (3.0/32.0)*2.0 + (5.0/16.0)*2.0 + (11.0 / 16.0) + (3.0/16.0)*4.0
	const r1 = (13.0 / 32.0) / 2.0
	const l2 = 1.0 / 2.0
	const r2 = (5.0 / 16.0) / 2.0
	const l3 = 3.0 / 8.0
	r3 := r2 - l3*math.Tan(sdf.DtoR(10.0))
	const l4 = 1.0 / 4.0

	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(r0, 0).Rel()
	p.Add(0, l0).Rel()
	p.Add(r1-r0, 0).Rel()
	p.Add(0, l1).Rel()
	p.Add(r2-r1, 0).Rel()
	p.Add(0, l2).Rel()
	p.Add(r3-r2, l3).Rel()
	p.Add(0, l4).Rel()
	p.Add(-r3, 0).Rel()

	shaft2d, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}

	shaft3d, err := sdf.Revolve3D(shaft2d)
	if err != nil {
		return nil, err
	}
	// make the cams
	const valveDiameter = 0.25
	const rockerRatio = 1.0
	const lift = valveDiameter * rockerRatio * 0.25
	const camDiameter = 5.0 / 8.0
	const camWidth = 3.0 / 16.0
	const k = 1.05
	inletTheta := sdf.DtoR(-110)

	inlet2d, _ := sdf.MakeThreeArcCam(lift, sdf.DtoR(115), camDiameter, k)
	inlet3d := sdf.Extrude3D(inlet2d, camWidth)
	exhaust2d, _ := sdf.MakeThreeArcCam(lift, sdf.DtoR(125), camDiameter, k)
	exhaust3d := sdf.Extrude3D(exhaust2d, camWidth)

	zOfs := (13.0 / 16.0) + (3.0 / 32.0) + (camWidth / 2.0)
	m := sdf.Translate3d(v3.Vec{0, 0, zOfs})
	m = sdf.RotateZ(0).Mul(m)
	ex4 := sdf.Transform3D(exhaust3d, m)

	zOfs += (5.0 / 16.0) + camWidth
	m = sdf.Translate3d(v3.Vec{0, 0, zOfs})
	m = sdf.RotateZ(inletTheta).Mul(m)
	in3 := sdf.Transform3D(inlet3d, m)

	zOfs += (11.0 / 16.0) + camWidth
	m = sdf.Translate3d(v3.Vec{0, 0, zOfs})
	m = sdf.RotateZ(inletTheta + sdf.Pi).Mul(m)
	in2 := sdf.Transform3D(inlet3d, m)

	zOfs += (5.0 / 16.0) + camWidth
	m = sdf.Translate3d(v3.Vec{0, 0, zOfs})
	m = sdf.RotateZ(sdf.Pi).Mul(m)
	ex1 := sdf.Transform3D(exhaust3d, m)

	return sdf.Union3D(shaft3d, ex1, in2, in3, ex4), nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := camshaft()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(s, 400, "camshaft.stl")
}

//-----------------------------------------------------------------------------
