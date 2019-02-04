//-----------------------------------------------------------------------------
/*

Fidget Spinners

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// 608 bearing
var bearing_outer_od = 22.0 // outer diameter of outer race
var bearing_outer_id = 19.2 // inner diameter of outer race
var bearing_inner_id = 8.0  // inner diameter of inner race
var bearing_inner_od = 12.1 // outer diameter of inner race
var bearing_thickness = 7.0 // bearing thickness

// Adjust clearance to give good interference fits for the bearings and spin caps.
var clearance = 0.0

//-----------------------------------------------------------------------------

// ball bearing counterweights
var bb_large_d = (1.0 / 2.0) * MillimetresPerInch
var bb_small_d = (5.0 / 16.0) * MillimetresPerInch

//-----------------------------------------------------------------------------

// Return an N petal bezier flower.
func flower(n int, r0, r1, r2 float64) SDF2 {

	theta := Tau / float64(n)
	b := NewBezier()

	p0 := V2{r1, 0}.Add(PolarToXY(r0, DtoR(-135)))
	p1 := V2{r1, 0}.Add(PolarToXY(r0, DtoR(-45)))
	p2 := V2{r1, 0}.Add(PolarToXY(r0, DtoR(45)))
	p3 := V2{r1, 0}.Add(PolarToXY(r0, DtoR(135)))
	p4 := PolarToXY(r2, theta/2)

	m := Rotate(theta)

	for i := 0; i < n; i++ {
		ofs := float64(i) * theta

		b.AddV2(p0).Handle(ofs+DtoR(-45), r0/2, r0/2)
		b.AddV2(p1).Handle(ofs+DtoR(45), r0/2, r0/2)
		b.AddV2(p2).Handle(ofs+DtoR(135), r0/2, r0/2)
		b.AddV2(p3).Handle(ofs+DtoR(225), r0/2, r0/2)
		b.AddV2(p4).Handle(ofs+theta/2+DtoR(90), r2/1.5, r2/1.5)

		p0 = m.MulPosition(p0)
		p1 = m.MulPosition(p1)
		p2 = m.MulPosition(p2)
		p3 = m.MulPosition(p3)
		p4 = m.MulPosition(p4)
	}

	b.Close()
	return Polygon2D(b.Polygon().Vertices())
}

func body1() SDF3 {

	n := 3
	t := bearing_thickness
	r := bearing_outer_od / 2

	r0 := r + 4.0
	r1 := 45.0 - r0
	r2 := r + 4.0

	// body
	s1 := ExtrudeRounded3D(flower(n, r0, r1, r2), t, t/4.0)
	// periphery holes
	s2 := MakeBoltCircle3D(t, r+clearance, r1, n)
	// center hole
	s3 := Cylinder3D(t, r+clearance, 0)

	return Difference3D(s1, Union3D(s2, s3))
}

//-----------------------------------------------------------------------------

func body2() SDF3 {
	t := bearing_thickness
	r := bearing_outer_od / 2
	r0 := r + 4.0

	// build the arm
	p := NewPolygon()
	p.Add(r, -t/2)
	p.Add(r0, -t/2)
	p.Add(r0, t/2)
	p.Add(r, t/2)
	theta := DtoR(270)
	arm := RevolveTheta3D(Polygon2D(p.Vertices()), theta)
	arm = Transform3D(arm, Translate3d(V3{-1.5 * r0, 0, 0}))

	// create 6 arms
	arms := RotateUnion3D(arm, 6, RotateZ(DtoR(60)))

	// add the center
	body := Union3D(Cylinder3D(t, r0, 0), arms)

	// remove the center hole
	return Difference3D(body, Cylinder3D(t, r, 0))
}

//-----------------------------------------------------------------------------

// Basic spin cap with variable pin size.
func spincap(
	pin_r float64, // pin radius
	pin_l float64, // pin length
) SDF3 {

	t := 3.0  // thickness of the spin cap
	st := 1.0 // spacer thickness

	r0 := bearing_outer_od / 2
	r1 := bearing_inner_od / 2

	p := NewPolygon()
	p.Add(0, -t-st)
	p.Add(r0, -t-st).Smooth(t/1.5, 6)
	p.Add(r0, -st)
	p.Add(r1, -st)
	p.Add(r1, 0)
	p.Add(pin_r, 0)
	p.Add(pin_r, pin_l)
	p.Add(0, pin_l)

	return Revolve3D(Polygon2D(p.Vertices()))
}

//-----------------------------------------------------------------------------

// Push to fit spincap for single spinner.
func spincap_single() SDF3 {
	gap := 1.0
	r := (bearing_inner_id / 2) - clearance
	l := (bearing_thickness - gap) / 2
	return spincap(r, l)
}

//-----------------------------------------------------------------------------

// Threaded spincap for double spinners.
func spincap_double(mode string) SDF3 {
	r := (bearing_inner_id / 2) - clearance
	thread_r := r * 0.8
	thread_pitch := 1.0
	thread_tolerance := 0.25
	l := bearing_thickness

	if mode == "male" {
		// Add an external screw thread.
		t := ISOThread(thread_r-thread_tolerance, thread_pitch, "external")
		screw := Screw3D(t, bearing_thickness, thread_pitch, 1)
		screw = ChamferedCylinder(screw, 0, 0.5)
		screw = Transform3D(screw, Translate3d(V3{0, 0, 1.5 * l}))
		return Union3D(spincap(r, l+0.5), screw)

	} else if mode == "female" {
		// Add an internal screw thread.
		t := ISOThread(thread_r, thread_pitch, "internal")
		screw := Screw3D(t, bearing_thickness, thread_pitch, 1)
		screw = Transform3D(screw, Translate3d(V3{0, 0, l * 0.5}))
		return Difference3D(spincap(r, l-0.5), screw)
	}

	panic("bad mode")
}

// Inner washer for double spinner.
func spincap_washer() SDF3 {
	t := 1.0
	r_inner := (bearing_inner_id / 2) * 1.05
	r_outer := (bearing_outer_od + bearing_inner_id) / 4
	return Washer3D(t, r_inner, r_outer)
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(body1(), 300, "body1.stl")
	RenderSTL(body2(), 300, "body2.stl")
	RenderSTL(spincap_single(), 150, "cap_single.stl")
	RenderSTL(spincap_double("male"), 150, "cap_double_male.stl")
	RenderSTL(spincap_double("female"), 150, "cap_double_female.stl")
	RenderSTL(spincap_washer(), 150, "washer.stl")
}

//-----------------------------------------------------------------------------
