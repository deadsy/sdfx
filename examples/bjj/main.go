//-----------------------------------------------------------------------------
/*

Bushing for the Box Joint Jig

https://woodgears.ca/box_joint/jig.html

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

// center hole
const ch_d = 0.755 * sdf.MillimetresPerInch
const ch_r = ch_d / 2.0

//-----------------------------------------------------------------------------

func bushing() (sdf.SDF3, error) {

	// R6-2RS 3/8 x 7/8 x 9/32 bearing
	bearing_outer_od := (7.0 / 8.0) * sdf.MillimetresPerInch // outer diameter of outer race
	//bearing_outer_id := 19.0                        // inner diameter of outer race
	bearing_inner_id := (3.0 / 8.0) * sdf.MillimetresPerInch   // inner diameter of inner race
	bearing_inner_od := 12.0                                   // outer diameter of inner race
	bearing_thickness := (9.0 / 32.0) * sdf.MillimetresPerInch // bearing thickness

	// Adjust clearance to give good interference fits for the bearing
	clearance := 0.0

	r0 := 2.3 // radius of central screw
	r1 := (bearing_outer_od + bearing_inner_od) / 4.0
	r2 := (bearing_inner_id / 2.0) - clearance

	h0 := 3.0 // height of cap
	h1 := h0 + bearing_thickness + 1.0

	p := sdf.NewPolygon()
	p.Add(r0, 0)
	p.Add(r1, 0)
	p.Add(r1, h0)
	p.Add(r2, h0)
	p.Add(r2, h1)
	p.Add(r0, h1)

	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	return sdf.Revolve3D(s)
}

//-----------------------------------------------------------------------------

// plateHoles2D returns 4 holes to attach the plate to the gear stack.
func plateHoles2D() (sdf.SDF2, error) {
	d := 17.0
	h, err := sdf.Circle2D(1.2)
	if err != nil {
		return nil, err
	}
	s0 := sdf.Transform2D(h, sdf.Translate2d(v2.Vec{d, d}))
	s1 := sdf.Transform2D(h, sdf.Translate2d(v2.Vec{-d, -d}))
	s2 := sdf.Transform2D(h, sdf.Translate2d(v2.Vec{-d, d}))
	s3 := sdf.Transform2D(h, sdf.Translate2d(v2.Vec{d, -d}))
	return sdf.Union2D(s0, s1, s2, s3), nil
}

const rod_r = (1.0 / 16.0) * sdf.MillimetresPerInch * 1.10

func lockingRod() (sdf.SDF3, error) {
	l := 62.0
	s0, err := sdf.Circle2D(rod_r)
	if err != nil {
		return nil, err
	}
	s1 := sdf.Box2D(v2.Vec{2 * rod_r, rod_r}, 0)
	s1 = sdf.Transform2D(s1, sdf.Translate2d(v2.Vec{0, -0.5 * rod_r}))
	s2 := sdf.Union2D(s0, s1)
	return sdf.Extrude3D(s2, l), nil
}

func plate() (sdf.SDF3, error) {
	r := (16.0 * gear_module / 2.0) * 0.83
	h := 5.0

	// plate
	s0, err := sdf.Cylinder3D(h, r, 0)
	if err != nil {
		return nil, err
	}

	// holes for attachment screws
	ph, err := plateHoles2D()
	if err != nil {
		return nil, err
	}
	s1 := sdf.Extrude3D(ph, h)

	// center hole
	s2, err := sdf.Cylinder3D(h, ch_r, 0)
	if err != nil {
		return nil, err
	}

	// indent for locking rod
	lr, err := lockingRod()
	if err != nil {
		return nil, err
	}
	m := sdf.Translate3d(v3.Vec{0, 0, h/2 - rod_r})
	m = m.Mul(sdf.RotateX(sdf.DtoR(-90.0)))
	s3 := sdf.Transform3D(lr, m)

	return sdf.Difference3D(s0, sdf.Union3D(s1, s2, s3)), nil
}

//-----------------------------------------------------------------------------

var gear_module = 80.0 / 16.0
var pressure_angle = sdf.DtoR(20)
var involute_facets = 10

func gears() (sdf.SDF3, error) {

	g_height := 10.0

	// 12 tooth spur gear
	k := obj.InvoluteGearParms{
		NumberTeeth:   12,
		Module:        gear_module,
		PressureAngle: pressure_angle,
		Facets:        involute_facets,
	}
	g0_2d, err := obj.InvoluteGear(&k)
	if err != nil {
		return nil, err
	}
	g0 := sdf.Extrude3D(g0_2d, g_height)

	// 16 tooth spur gear
	k = obj.InvoluteGearParms{
		NumberTeeth:   16,
		Module:        gear_module,
		PressureAngle: pressure_angle,
		Facets:        involute_facets,
	}
	g1_2d, err := obj.InvoluteGear(&k)
	if err != nil {
		return nil, err
	}
	g1 := sdf.Extrude3D(g1_2d, g_height)

	s0 := sdf.Transform3D(g0, sdf.Translate3d(v3.Vec{0, 0, g_height / 2.0}))
	s1 := sdf.Transform3D(g1, sdf.Translate3d(v3.Vec{0, 0, -g_height / 2.0}))

	// center hole
	s2, err := sdf.Cylinder3D(2.0*g_height, ch_r, 0)
	if err != nil {
		return nil, err
	}

	// holes for attachment screws
	ph, err := plateHoles2D()
	if err != nil {
		return nil, err
	}
	screw_depth := 10.0
	s3 := sdf.Extrude3D(ph, screw_depth)
	s3 = sdf.Transform3D(s3, sdf.Translate3d(v3.Vec{0, 0, screw_depth/2.0 - g_height}))

	return sdf.Difference3D(sdf.Union3D(s0, s1), sdf.Union3D(s2, s3)), nil
}

//-----------------------------------------------------------------------------

func main() {
	bushing, err := bushing()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(bushing, 100, "bushing.stl")

	gears, err := gears()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(gears, 300, "gear.stl")

	plate, err := plate()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(plate, 300, "plate.stl")
}

//-----------------------------------------------------------------------------
