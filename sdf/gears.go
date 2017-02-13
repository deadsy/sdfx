//-----------------------------------------------------------------------------
/*

Involute Gears

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

// return the involute coordinate
// r = base radius
// theta = involute angle
func involute(r, theta float64) V2 {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return V2{
		r * (c + theta*s),
		r * (s - theta*c),
	}
}

// return the involute angle
// r = base radius
// d = involute radial distance
func involute_angle(r, d float64) float64 {
	x := d / r
	return math.Sqrt(x*x - 1)
}

//-----------------------------------------------------------------------------

// Generate an SDF2 polygon for a single involute tooth
// number_teeth = number of gear teeth
// gear_module = pitch circle diameter / number of gear teeth
// root_radius = radius at tooth root
// base_radius = radius at the base of the involute
// outer_radius = radius at the outside of the tooth
// backlash = backlash expressed as units of pitch circumference
// facets = number of facets for involute flank
func InvoluteGearTooth(
	number_teeth int,
	gear_module float64,
	root_radius float64,
	base_radius float64,
	outer_radius float64,
	backlash float64,
	facets int,
) SDF2 {

	pitch_radius := float64(number_teeth) * gear_module / 2.0

	// work out the angular extent of the tooth on the base radius
	pitch_point := involute(base_radius, involute_angle(base_radius, pitch_radius))
	face_angle := math.Atan2(pitch_point.Y, pitch_point.X)
	backlash_angle := backlash / (2.0 * pitch_radius)
	center_angle := PI/(2.0*float64(number_teeth)) + face_angle - backlash_angle

	// work out the angles over which the involute will be used
	start_angle := involute_angle(base_radius, Max(base_radius, root_radius))
	stop_angle := involute_angle(base_radius, outer_radius)
	dtheta := (stop_angle - start_angle) / float64(facets)

	v := make([]V2, 2*(facets+1)+1)

	// lower tooth face
	m := Rotate(-center_angle)
	angle := start_angle
	for i := 0; i <= facets; i++ {
		v[i] = m.MulPosition(involute(base_radius, angle))
		angle += dtheta
	}

	// upper tooth face (mirror the lower point)
	for i := 0; i <= facets; i++ {
		p := v[facets-i]
		v[facets+1+i] = V2{p.X, -p.Y}
	}

	// add the origin to make the polygon a tooth wedge
	v[2*(facets+1)] = V2{0, 0}

	return NewPolySDF2(v)
}

//-----------------------------------------------------------------------------

// Generate an SDF2 polygon for an involute gear
// number_teeth = number of gear teeth
// gear_module = pitch circle diameter / number of gear teeth
// pressure_angle = gear pressure angle (radians)
// backlash = backlash expressed as units of pitch circumference
// clearance = additional root clearance
// ring_width = width of ring wall (from root circle)
// facets = number of facets for involute flank
func InvoluteGear(
	number_teeth int,
	gear_module float64,
	pressure_angle float64,
	backlash float64,
	clearance float64,
	ring_width float64,
	facets int,
) SDF2 {

	// pitch radius
	pitch_radius := float64(number_teeth) * gear_module / 2.0

	// base circle radius
	base_radius := pitch_radius * math.Cos(pressure_angle)

	// addendum: radial distance from pitch circle to outside circle
	addendum := gear_module * 1.0
	// dedendum: radial distance from pitch circle to root circle
	dedendum := addendum + clearance

	outer_radius := pitch_radius + addendum
	root_radius := pitch_radius - dedendum
	ring_radius := root_radius - ring_width

	tooth := InvoluteGearTooth(
		number_teeth,
		gear_module,
		root_radius,
		base_radius,
		outer_radius,
		backlash,
		facets,
	)

	gear := NewRotateSDF2(tooth, number_teeth, Rotate2d(TAU/float64(number_teeth)))
	root := NewCircleSDF2(root_radius)
	ring := NewCircleSDF2(ring_radius)

	return NewDifferenceSDF2(NewUnionSDF2(gear, root), ring)
}

//-----------------------------------------------------------------------------
