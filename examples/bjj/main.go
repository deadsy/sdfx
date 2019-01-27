//-----------------------------------------------------------------------------
/*

Bushing for the Box Joint Jig

https://woodgears.ca/box_joint/jig.html

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// center hole
const ch_d = 0.755 * MillimetresPerInch
const ch_r = ch_d / 2.0

//-----------------------------------------------------------------------------

func bushing() SDF3 {

	// R6-2RS 3/8 x 7/8 x 9/32 bearing
	bearing_outer_od := (7.0 / 8.0) * MillimetresPerInch // outer diameter of outer race
	//bearing_outer_id := 19.0                        // inner diameter of outer race
	bearing_inner_id := (3.0 / 8.0) * MillimetresPerInch   // inner diameter of inner race
	bearing_inner_od := 12.0                               // outer diameter of inner race
	bearing_thickness := (9.0 / 32.0) * MillimetresPerInch // bearing thickness

	// Adjust clearance to give good interference fits for the bearing
	clearance := 0.0

	r0 := 2.3 // radius of central screw
	r1 := (bearing_outer_od + bearing_inner_od) / 4.0
	r2 := (bearing_inner_id / 2.0) - clearance

	h0 := 3.0 // height of cap
	h1 := h0 + bearing_thickness + 1.0

	p := NewPolygon()
	p.Add(r0, 0)
	p.Add(r1, 0)
	p.Add(r1, h0)
	p.Add(r2, h0)
	p.Add(r2, h1)
	p.Add(r0, h1)
	return Revolve3D(Polygon2D(p.Vertices()))
}

//-----------------------------------------------------------------------------

// 4 holes to attach the plate to the gear stack.
func plate_holes_2d() SDF2 {
	d := 17.0
	h := Circle2D(1.2)
	s0 := Transform2D(h, Translate2d(V2{d, d}))
	s1 := Transform2D(h, Translate2d(V2{-d, -d}))
	s2 := Transform2D(h, Translate2d(V2{-d, d}))
	s3 := Transform2D(h, Translate2d(V2{d, -d}))
	return Union2D(s0, s1, s2, s3)
}

const rod_r = (1.0 / 16.0) * MillimetresPerInch * 1.10

func locking_rod() SDF3 {
	l := 62.0
	s0 := Circle2D(rod_r)
	s1 := Box2D(V2{2 * rod_r, rod_r}, 0)
	s1 = Transform2D(s1, Translate2d(V2{0, -0.5 * rod_r}))
	s2 := Union2D(s0, s1)
	return Extrude3D(s2, l)
}

func plate() SDF3 {
	r := (16.0 * gear_module / 2.0) * 0.83
	h := 5.0

	// plate
	s0 := Cylinder3D(h, r, 0)
	// holes for attachment screws
	s1 := Extrude3D(plate_holes_2d(), h)
	// center hole
	s2 := Cylinder3D(h, ch_r, 0)
	// indent for locking rod
	m := Translate3d(V3{0, 0, h/2 - rod_r}).Mul(RotateX(DtoR(-90.0)))
	s3 := Transform3D(locking_rod(), m)

	return Difference3D(s0, Union3D(s1, s2, s3))
}

//-----------------------------------------------------------------------------

var gear_module = 80.0 / 16.0
var pressure_angle = 20.0
var involute_facets = 10

func gears() SDF3 {

	g_height := 10.0

	// 12 tooth spur gear
	g0_teeth := 12
	g0_pd := float64(g0_teeth) * gear_module
	g0_2d := InvoluteGear(
		g0_teeth,
		gear_module,
		DtoR(pressure_angle),
		0.0,
		0.0,
		g0_pd/2.0,
		involute_facets,
	)
	g0 := Extrude3D(g0_2d, g_height)

	// 16 tooth spur gear
	g1_teeth := 16
	g1_pd := float64(g1_teeth) * gear_module
	g1_2d := InvoluteGear(
		g1_teeth,
		gear_module,
		DtoR(pressure_angle),
		0.0,
		0.0,
		g1_pd/2.0,
		involute_facets,
	)
	g1 := Extrude3D(g1_2d, g_height)

	s0 := Transform3D(g0, Translate3d(V3{0, 0, g_height / 2.0}))
	s1 := Transform3D(g1, Translate3d(V3{0, 0, -g_height / 2.0}))

	// center hole
	s2 := Cylinder3D(2.0*g_height, ch_r, 0)

	// holes for attachment screws
	screw_depth := 10.0
	s3 := Extrude3D(plate_holes_2d(), screw_depth)
	s3 = Transform3D(s3, Translate3d(V3{0, 0, screw_depth/2.0 - g_height}))

	return Difference3D(Union3D(s0, s1), Union3D(s2, s3))
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(bushing(), 100, "bushing.stl")
	RenderSTL(gears(), 300, "gear.stl")
	RenderSTL(plate(), 300, "plate.stl")
}

//-----------------------------------------------------------------------------
