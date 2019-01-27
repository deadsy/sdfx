// wallaby camshaft

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

func main() {

	// build the shaft from an SoR
	l0 := 13.0 / 16.0
	r0 := (5.0 / 16.0) / 2.0
	l1 := (3.0/32.0)*2.0 + (5.0/16.0)*2.0 + (11.0 / 16.0) + (3.0/16.0)*4.0
	r1 := (13.0 / 32.0) / 2.0
	l2 := 1.0 / 2.0
	r2 := (5.0 / 16.0) / 2.0
	l3 := 3.0 / 8.0
	r3 := r2 - l3*math.Tan(DtoR(10.0))
	l4 := 1.0 / 4.0

	p := NewPolygon()
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

	shaft_2d := Polygon2D(p.Vertices())
	shaft_3d := Revolve3D(shaft_2d)

	// make the cams
	valve_diameter := 0.25
	rocker_ratio := 1.0
	lift := valve_diameter * rocker_ratio * 0.25
	cam_diameter := 5.0 / 8.0
	cam_width := 3.0 / 16.0
	k := 1.05
	inlet_theta := DtoR(-110)

	inlet_2d, _ := MakeThreeArcCam(lift, DtoR(115), cam_diameter, k)
	inlet_3d := Extrude3D(inlet_2d, cam_width)
	exhaust_2d, _ := MakeThreeArcCam(lift, DtoR(125), cam_diameter, k)
	exhaust_3d := Extrude3D(exhaust_2d, cam_width)

	z_ofs := (13.0 / 16.0) + (3.0 / 32.0) + (cam_width / 2.0)
	m := Translate3d(V3{0, 0, z_ofs})
	m = RotateZ(0).Mul(m)
	ex4 := Transform3D(exhaust_3d, m)

	z_ofs += (5.0 / 16.0) + cam_width
	m = Translate3d(V3{0, 0, z_ofs})
	m = RotateZ(inlet_theta).Mul(m)
	in3 := Transform3D(inlet_3d, m)

	z_ofs += (11.0 / 16.0) + cam_width
	m = Translate3d(V3{0, 0, z_ofs})
	m = RotateZ(inlet_theta + Pi).Mul(m)
	in2 := Transform3D(inlet_3d, m)

	z_ofs += (5.0 / 16.0) + cam_width
	m = Translate3d(V3{0, 0, z_ofs})
	m = RotateZ(Pi).Mul(m)
	ex1 := Transform3D(exhaust_3d, m)

	RenderSTL(Union3D(shaft_3d, ex1, in2, in3, ex4), 400, "camshaft.stl")
}
