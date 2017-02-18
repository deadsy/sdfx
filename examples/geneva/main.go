// 3d printable geneva cam mechanism

package main

import . "github.com/deadsy/sdfx/sdf"

func main() {

	wheel_h := 5.0 // height of wheels
	hole_r := 3.25 // radius of center hole
	hub_r := 10.0  // hub radius for driven wheel
	base_r := 30.0 // radius of base for driver wheel

	s_driver, s_driven, err := MakeGenevaCam(6, 50, 20, 40, 2.5, 0.1)
	if err != nil {
		panic(err)
	}

	// extrude the driver wheel
	driver_3d := NewExtrudeSDF3(s_driver, wheel_h)
	driver_3d = NewTransformSDF3(driver_3d, Translate3d(V3{0, 0, wheel_h / 2}))
	// add a base
	base_3d := NewCylinderSDF3(wheel_h, base_r, 0)
	base_3d = NewTransformSDF3(base_3d, Translate3d(V3{0, 0, -wheel_h / 2}))
	driver_3d = NewUnionSDF3(driver_3d, base_3d)
	// remove a center hole
	hole_3d := NewCylinderSDF3(2*wheel_h, hole_r, 0)
	driver_3d = NewDifferenceSDF3(driver_3d, hole_3d)

	// extrude the driven wheel
	driven_3d := NewExtrudeSDF3(s_driven, wheel_h)
	driven_3d = NewTransformSDF3(driven_3d, Translate3d(V3{0, 0, -wheel_h / 2}))
	// add a hub
	hub_3d := NewCylinderSDF3(wheel_h, hub_r, 0)
	hub_3d = NewTransformSDF3(hub_3d, Translate3d(V3{0, 0, wheel_h / 2}))
	driven_3d = NewUnionSDF3(driven_3d, hub_3d)
	// remove a center hole
	driven_3d = NewDifferenceSDF3(driven_3d, hole_3d)

	RenderSTL(driver_3d, "driver.stl")
	RenderSTL(driven_3d, "driven.stl")
}
