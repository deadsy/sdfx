//-----------------------------------------------------------------------------
/*

Involute Gears

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

// return the involute coordinate for a given angle
func involute_xy(
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
func involute_theta(
	r float64, // base radius
	d float64, // involute radial distance
) float64 {
	x := d / r
	return math.Sqrt(x*x - 1)
}

// return the closest distance between a polar point and an involute
// > 0 on convex side of involute
func involute_distance(
	b float64, // base radius
	r float64, // radius for point (>b)
	theta float64, // theta for point
) (d, d_theta float64) {
	d_theta = math.Acos(b/r) + theta
	d = math.Sqrt((r*r)-(b*b)) - (b * theta)
	return
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
	pitch_point := involute_xy(base_radius, involute_theta(base_radius, pitch_radius))
	face_angle := math.Atan2(pitch_point.Y, pitch_point.X)
	backlash_angle := backlash / (2.0 * pitch_radius)
	center_angle := PI/(2.0*float64(number_teeth)) + face_angle - backlash_angle

	// work out the angles over which the involute will be used
	start_angle := involute_theta(base_radius, Max(base_radius, root_radius))
	stop_angle := involute_theta(base_radius, outer_radius)
	dtheta := (stop_angle - start_angle) / float64(facets)

	v := make([]V2, 2*(facets+1)+1)

	// lower tooth face
	m := Rotate(-center_angle)
	angle := start_angle
	for i := 0; i <= facets; i++ {
		v[i] = m.MulPosition(involute_xy(base_radius, angle))
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
	backlash float64, // backlash expressed as per-tooth distance at pitch circumference
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

	gear := RotateCopy2D(tooth, number_teeth)
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
	base_angle   float64 // involute base angle (at base radius)
	start_angle  float64 // involute start angle (at root radius)
	stop_angle   float64 // involute stop angle (at outer radius)
	start_xy     V2      // involute start coordinate
	stop_xy      V2      // involute stop ccordinate
	bb           Box2    // bounding box
}

// InvoluteGear2D returns the 2D profile for an involute gear.
func InvoluteGear2D(
	number_teeth int, // number of gear teeth
	gear_module float64, // pitch circle diameter / number of gear teeth
	pressure_angle float64, // gear pressure angle (radians)
	backlash float64, // backlash expressed as per-tooth distance at pitch circumference
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
	s.root_radius = Max(s.base_radius, pitch_radius-dedendum)
	// radius of inner gear ring
	s.ring_radius = s.root_radius - ring_width

	// involute angles at various radii
	outer_angle := involute_theta(s.base_radius, s.outer_radius)
	pitch_angle := involute_theta(s.base_radius, pitch_radius)
	root_angle := involute_theta(s.base_radius, s.root_radius)

	// work out the half angle subtended by the top land of the tooth at the outer radius
	backlash_angle := backlash / (2.0 * pitch_radius)
	top_angle := s.tooth_angle/4.0 - (outer_angle - pitch_angle) - backlash_angle

	// store the base, start and stop angles for the involute portion of the tooth
	s.base_angle = top_angle + outer_angle
	s.start_angle = top_angle + outer_angle - root_angle
	s.stop_angle = top_angle

	// store the xy positions for the start and stop involute points
	s.start_xy = PolarToXY(s.root_radius, s.start_angle)
	s.stop_xy = PolarToXY(s.outer_radius, s.stop_angle)

	return &s
}

func (s *InvoluteGearSDF2) involute_distance(p V2, p_r, p_theta float64) float64 {
	if p_r < s.base_radius {
		return p.Sub(s.start_xy).Length()
	}
	d, d_theta := involute_distance(s.base_radius, p_r, s.base_angle-p_theta)
	d_theta = s.base_angle - d_theta
	if d_theta > s.start_angle {
		return p.Sub(s.start_xy).Length()
	}
	if d_theta < s.stop_angle {
		return p.Sub(s.stop_xy).Length()
	}
	return d
}

// Evaluate returns the minimum distance to the involute gear.
func (s *InvoluteGearSDF2) Evaluate(p V2) float64 {
	// work out the polar coordinates of p
	p_theta := math.Atan2(p.Y, p.X)
	p_r := p.Length()

	d_ring := Abs(p_r - s.ring_radius)
	d_root := Abs(p_r - s.root_radius)
	d_outer := Abs(p_r - s.outer_radius)

	// check the ring radius
	if p_r < s.ring_radius {
		// within the ring radius
		return d_ring
	}

	// map the angle back to the 0th tooth (about the x-axis)
	p_theta = SawTooth(p_theta, s.tooth_angle)
	// the tooth is symmetrical about the x-axis, only consider the 1st quadrant (+x,+y)
	p_theta = Abs(p_theta)

	if p_theta < s.stop_angle {
		if p_r > s.outer_radius {
			return d_outer
		} else {
			d_involute := Abs(s.involute_distance(p, p_r, p_theta))
			return -Min(d_involute, Min(d_outer, d_ring))
		}
	} else if p_theta < s.start_angle {

		// TODO
		return 0

	} else {
		if p_r < s.root_radius {
			return -Min(d_ring, d_root)
		} else {
			d_involute := Abs(s.involute_distance(p, p_r, p_theta))
			return Min(d_involute, d_root)
		}
	}

	panic("")
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
