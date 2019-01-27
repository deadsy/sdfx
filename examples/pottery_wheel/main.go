//-----------------------------------------------------------------------------
/*

Pottery Wheel

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------
// overall build controls

const SCALE = 1.0 / 0.98 // 2% Al shrinkage
const core_print = false // add the core print to the wheel
const pie_print = false  // create a 1/n pie segment (n = number of webs)

//-----------------------------------------------------------------------------

// dimension scaling
func dim(x float64) float64 {
	return SCALE * x
}

//-----------------------------------------------------------------------------

// draft angles
var draft_angle = DtoR(4.0)       // standard overall draft
var core_draft_angle = DtoR(10.0) // draft angle for the core print

// nominal size values (mm)
var wheel_diameter = dim(MillimetresPerInch * 8.0) // total wheel diameter
var hub_diameter = dim(40.0)                       // base diameter of central shaft hub
var hub_height = dim(53.0)                         // height of cental shaft hub
var shaft_diameter = dim(21.0)                     // 1" target size - reduced for machining allowance
var shaft_length = dim(45.0)                       // length of shaft bore
var wall_height = dim(35.0)                        // height of wheel side walls
var wall_thickness = dim(4.0)                      // base thickness of outer wheel walls
var plate_thickness = dim(7.0)                     // thickness of wheel top plate
var web_width = dim(2.0)                           // thickness of reinforcing webs
var web_height = dim(25.0)                         // height of reinforcing webs
var core_height = dim(15.0)                        // height of core print
var number_of_webs = 6                             // number of reinforcing webs

// derived values
var wheel_radius = wheel_diameter / 2
var hub_radius = hub_diameter / 2
var shaft_radius = shaft_diameter / 2
var web_length = wheel_radius - wall_thickness - hub_radius

//-----------------------------------------------------------------------------

// build 2d wheel profile
func wheel_profile() SDF2 {

	draft0 := (hub_height - plate_thickness) * math.Tan(draft_angle)
	draft1 := (wall_height - plate_thickness) * math.Tan(draft_angle)
	draft2 := wall_height * math.Tan(draft_angle)
	draft3 := core_height * math.Tan(core_draft_angle)

	s := NewPolygon()

	if core_print {
		s.Add(0, 0)
		s.Add(wheel_radius+draft2, 0)
		s.Add(wheel_radius, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness-draft1, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius+draft0, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius, hub_height).Smooth(2.0, 5)
		s.Add(shaft_radius, hub_height)
		s.Add(shaft_radius-draft3, hub_height+core_height)
		s.Add(0, hub_height+core_height)
	} else {
		s.Add(0, 0)
		s.Add(wheel_radius+draft2, 0)
		s.Add(wheel_radius, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness-draft1, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius+draft0, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius, hub_height).Smooth(2.0, 5)
		s.Add(shaft_radius, hub_height)
		s.Add(shaft_radius, hub_height-shaft_length)
		s.Add(0, hub_height-shaft_length)
	}

	//s.Render("wheel.dxf")
	return Polygon2D(s.Vertices())
}

// build 2d web profile
func web_profile() SDF2 {

	draft := web_height * math.Tan(draft_angle)
	x0 := web_width + draft
	x1 := web_width

	s := NewPolygon()
	s.Add(-x0, 0)
	s.Add(-x1, web_height).Smooth(1.0, 3)
	s.Add(x1, web_height).Smooth(1.0, 3)
	s.Add(x0, 0)

	//s.Render("web.dxf")
	return Polygon2D(s.Vertices())
}

// build the wheel pattern
func wheel_pattern() {

	// build a reinforcing webs
	web_2d := web_profile()
	web_3d := Extrude3D(web_2d, web_length)
	m := Translate3d(V3{0, plate_thickness, hub_radius + web_length/2})
	m = RotateX(DtoR(90)).Mul(m)

	// build the wheel profile
	wheel_2d := wheel_profile()
	var wheel_3d SDF3

	if pie_print {
		m = RotateZ(DtoR(120)).Mul(m)
		web_3d = Transform3D(web_3d, m)
		wheel_3d = RevolveTheta3D(wheel_2d, DtoR(60))
	} else {
		m = RotateZ(DtoR(90)).Mul(m)
		web_3d = Transform3D(web_3d, m)
		web_3d = RotateCopy3D(web_3d, 6)
		wheel_3d = Revolve3D(wheel_2d)
	}

	// add the webs to the wheel with some blending
	wheel := Union3D(wheel_3d, web_3d)
	wheel.(*UnionSDF3).SetMin(PolyMin(wall_thickness))

	RenderSTL(wheel, 200, "wheel.stl")
	RenderDXF(Slice2D(wheel, V3{0, 0, 15.0}, V3{0, 0, 1}), 200, "wheel.dxf")
}

//-----------------------------------------------------------------------------

// build 2d core profile
func core_profile() SDF2 {

	draft := core_height * math.Tan(core_draft_angle)

	s := NewPolygon()
	s.Add(0, 0)
	s.Add(shaft_radius-draft, 0)
	s.Add(shaft_radius, core_height)
	s.Add(shaft_radius, core_height+shaft_length).Smooth(2.0, 3)
	s.Add(0, core_height+shaft_length)

	//s.Render("core.dxf")
	return Polygon2D(s.Vertices())
}

// build the core box
func core_box() {

	// build the box
	w := 4.2 * shaft_radius
	d := 1.2 * shaft_radius
	h := (core_height + shaft_length) * 1.1
	box_3d := Box3D(V3{h, w, d}, 0)

	// holes in the box
	dy := w * 0.37
	dx := h * 0.4
	hole_radius := ((3.0 / 16.0) * MillimetresPerInch) / 2.0
	positions := []V2{
		{dx, dy},
		{-dx, dy},
		{dx, -dy},
		{-dx, -dy}}
	holes_3d := MultiCylinder3D(d, hole_radius, positions)

	// Drill the holes
	box_3d = Difference3D(box_3d, holes_3d)

	// build the core
	core_2d := core_profile()
	core_3d := Revolve3D(core_2d)
	m := Translate3d(V3{h / 2, 0, d / 2}).Mul(RotateY(DtoR(-90)))
	core_3d = Transform3D(core_3d, m)

	// remove the core from the box
	core_box := Difference3D(box_3d, core_3d)

	RenderSTL(core_box, 200, "core_box.stl")
}

//-----------------------------------------------------------------------------

func main() {
	wheel_pattern()
	core_box()
}

//-----------------------------------------------------------------------------
