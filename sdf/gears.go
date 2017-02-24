//-----------------------------------------------------------------------------
/*

Involute Gears

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

// return the involute coordinate for a given angle
func involute(
	r float64, // base radius
	theta float64, // involute angle
) V2 {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return V2{
		r * (c + theta*s),
		r * (s - theta*c),
	}
}

// return the involute angle for a given radial distance
func involute_angle(
	r float64, // base radius
	d float64, // involute radial distance
) float64 {
	x := d / r
	return math.Sqrt(x*x - 1)
}

//-----------------------------------------------------------------------------

// InvoluteGearTooth returns a 2D profile for a single involute tooth.
func InvoluteGearTooth(
	number_teeth int, // number of gear teeth
	gear_module float64, // pitch circle diameter / number of gear teeth
	root_radius float64, // radius at tooth root
	base_radius float64, // radius at the base of the involute
	outer_radius float64, // radius at the outside of the tooth
	backlash float64, // backlash expressed as units of pitch circumference
	facets int, // number of facets for involute flank
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

	return Polygon2D(v)
}

//-----------------------------------------------------------------------------

// InvoluteGear returns an 2D polygon for an involute gear.
func InvoluteGear(
	number_teeth int, // number of gear teeth
	gear_module float64, // pitch circle diameter / number of gear teeth
	pressure_angle float64, // gear pressure angle (radians)
	backlash float64, // backlash expressed as units of pitch circumference
	clearance float64, // additional root clearance
	ring_width float64, // width of ring wall (from root circle)
	facets int, // number of facets for involute flank
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

	gear := Rotate2D(tooth, number_teeth, Rotate2d(TAU/float64(number_teeth)))
	root := Circle2D(root_radius)
	ring := Circle2D(ring_radius)

	return Difference2D(Union2D(gear, root), ring)
}

//-----------------------------------------------------------------------------
// 2D Involute Gear

type InvoluteGearSDF2 struct {
	base_radius  float64 // base radius for the involute
	outer_radius float64 // radius for outside of gear
	root_radius  float64 // radius for root of gear tooth
	ring_radius  float64 // radius of inner gear ring
	tooth_angle  float64 // angle subtended by a single gear tooth
	flank_angle  float64 // angle subtended by one gear tooth flank (involute portion)
	top_angle    float64 // half angle subtended by top of gear tooth
	bb           Box2    // bounding box
}

// InvoluteGear2D returns the 2D profile for an involute gear.
func InvoluteGear2D(
	number_teeth int, // number of gear teeth
	gear_module float64, // pitch circle diameter / number of gear teeth
	pressure_angle float64, // gear pressure angle (radians)
	backlash float64, // backlash expressed as distance at pitch circumference
	clearance float64, // additional root clearance
	ring_width float64, // width of ring wall (from root circle)
) SDF2 {
	s := InvoluteGearSDF2{}
	// tooth angle
	s.tooth_angle = TAU / float64(number_teeth)
	// radius at gear pitch line
	pitch_radius := float64(number_teeth) * gear_module / 2.0
	// radius for base circle of involute
	s.base_radius = pitch_radius * math.Cos(pressure_angle)
	// addendum: radial distance from pitch circle to outside circle
	addendum := gear_module * 1.0
	// dedendum: radial distance from pitch circle to root circle
	dedendum := addendum + clearance
	// radius for outside of gear
	s.outer_radius = pitch_radius + addendum
	// radius for root of gear tooth
	s.root_radius = pitch_radius - dedendum
	// radius of inner gear ring
	s.ring_radius = s.root_radius - ring_width

	// work out the angle for the flank of the gear tooth (involute portion)
	theta0 := involute_angle(s.base_radius, Max(s.base_radius, s.root_radius))
	theta1 := involute_angle(s.base_radius, s.outer_radius)
	s.flank_angle = theta1 - theta0

	// work out the angle for the top of the gear tooth
	backlash_angle := PI * backlash / pitch_radius
	s.top_angle = ((s.tooth_angle / 2) - s.flank_angle - backlash_angle) / 2

	return &s
}

// Evaluate returns the minimum distance to the involute gear.
func (s *InvoluteGearSDF2) Evaluate(p V2) float64 {
	// work out the polar coordinates
	p_theta := math.Atan2(p.Y, p.X)
	p_distance := p.Length()
	// map the angle back to the 0th tooth (on the x-axis)
	// the tooth is symmetrical about the x-axis, so only consider the 1st quadrant
	p_theta = Abs(SawTooth(p_theta, s.tooth_angle))

	// work out the involute angle for the closest point on the involute
	i_theta := math.Acos(s.base_radius/p_distance) + p_theta

	_ = i_theta

	// distance between circle tangent and point

	return 0
}

// BoundingBox returns the bounding box for the involute gear.
func (s *InvoluteGearSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// 2D Gear Rack

type GearRackSDF2 struct {
	tooth  SDF2    // polygon for rack tooth
	pitch  float64 // tooth to tooth pitch
	length float64 // half the total rack length
	bb     Box2    // bounding box
}

// GearRack2D returns the 2D profile for a gear rack.
func GearRack2D(
	number_teeth float64, // number of rack teeth
	gear_module float64, // pitch circle diameter / number of gear teeth
	pressure_angle float64, // gear pressure angle (radians)
	backlash float64, // backlash expressed as units of pitch circumference
	base_height float64, // height of rack base
) SDF2 {
	s := GearRackSDF2{}

	// addendum: distance from pitch line to top of tooth
	addendum := gear_module * 1.0
	// dedendum: distance from pitch line to root of tooth
	dedendum := gear_module * 1.25
	// total tooth height
	tooth_height := base_height + addendum + dedendum
	// tooth_pitch: tooth to tooth distance along pitch line
	pitch := gear_module * PI

	// x size of tooth flank
	dx := (addendum + dedendum) * math.Tan(pressure_angle)
	// 1/2 x size of tooth top
	dxt := ((pitch / 2.0) - dx) / 2.0
	// x size of backlash
	bl := backlash / 2.0

	// create a half tooth profile centered on the y-axis
	tooth := []V2{
		V2{pitch, 0},
		V2{pitch, base_height},
		V2{dx + dxt - bl, base_height},
		V2{dxt - bl, tooth_height},
		V2{-pitch, tooth_height},
		V2{-pitch, 0},
	}

	s.tooth = Polygon2D(tooth)
	s.pitch = pitch
	s.length = pitch * number_teeth / 2.0
	s.bb = Box2{V2{-s.length, 0}, V2{s.length, tooth_height}}
	return &s
}

// Evaluate returns the minimum distance to the gear rack.
func (s *GearRackSDF2) Evaluate(p V2) float64 {
	// map p.X back to the [0,half_pitch) domain
	p0 := V2{Abs(SawTooth(p.X, s.pitch)), p.Y}
	// get the tooth profile distance
	d0 := s.tooth.Evaluate(p0)
	// create a region for the rack length
	d1 := Abs(p.X) - s.length
	// return the intersection
	return Max(d0, d1)
}

// BoundingBox returns the bounding box for the gear rack.
func (s *GearRackSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
