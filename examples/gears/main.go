package main

import . "github.com/deadsy/sdfx/sdf"

func main() {

	module := (5.0 / 8.0) / 20.0
	pa := DtoR(20.0)
	h := 0.15
	number_teeth := 20

	gear_2d := InvoluteGear(
		number_teeth, // number_teeth
		module,       // gear_module
		pa,           // pressure_angle
		0.0,          // backlash
		0.0,          // clearance
		0.05,         // ring_width
		7,            // facets
	)
	gear_3d := Extrude3D(gear_2d, h)
	m := Rotate3d(V3{0, 0, 1}, DtoR(180.0/float64(number_teeth)))
	m = Translate3d(V3{0, 0.39, 0}).Mul(m)
	gear_3d = Transform3D(gear_3d, m)

	rack_2d := NewGearRack(
		11,     // number_teeth
		module, // gear_module
		pa,     // pressure_angle
		0.00,   // backlash
		0.025,  // base_height
	)
	rack_3d := Extrude3D(rack_2d, h)

	s := Union3D(rack_3d, gear_3d)
	RenderSTL(s, 200, "gear.stl")
}
