//-----------------------------------------------------------------------------
/*

Involute Gears

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

// return the involute coordinate for a given angle
func involuteXY(
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
func involuteTheta(
	r float64, // base radius
	d float64, // involute radial distance
) float64 {
	x := d / r
	return math.Sqrt(x*x - 1)
}

//-----------------------------------------------------------------------------

// InvoluteGearTooth returns a 2D profile for a single involute tooth.
func InvoluteGearTooth(
	numberTeeth int, // number of gear teeth
	gearModule float64, // pitch circle diameter / number of gear teeth
	rootRadius float64, // radius at tooth root
	baseRadius float64, // radius at the base of the involute
	outerRadius float64, // radius at the outside of the tooth
	backlash float64, // backlash expressed as units of pitch circumference
	facets int, // number of facets for involute flank
) SDF2 {

	pitchRadius := float64(numberTeeth) * gearModule / 2.0

	// work out the angular extent of the tooth on the base radius
	pitchPoint := involuteXY(baseRadius, involuteTheta(baseRadius, pitchRadius))
	faceAngle := math.Atan2(pitchPoint.Y, pitchPoint.X)
	backlashAngle := backlash / (2.0 * pitchRadius)
	centerAngle := PI/(2.0*float64(numberTeeth)) + faceAngle - backlashAngle

	// work out the angles over which the involute will be used
	startAngle := involuteTheta(baseRadius, Max(baseRadius, rootRadius))
	stopAngle := involuteTheta(baseRadius, outerRadius)
	dtheta := (stopAngle - startAngle) / float64(facets)

	v := make([]V2, 2*(facets+1)+1)

	// lower tooth face
	m := Rotate(-centerAngle)
	angle := startAngle
	for i := 0; i <= facets; i++ {
		v[i] = m.MulPosition(involuteXY(baseRadius, angle))
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
	numberTeeth int, // number of gear teeth
	gearModule float64, // pitch circle diameter / number of gear teeth
	pressureAngle float64, // gear pressure angle (radians)
	backlash float64, // backlash expressed as per-tooth distance at pitch circumference
	clearance float64, // additional root clearance
	ringWidth float64, // width of ring wall (from root circle)
	facets int, // number of facets for involute flank
) SDF2 {

	// pitch radius
	pitchRadius := float64(numberTeeth) * gearModule / 2.0

	// base circle radius
	baseRadius := pitchRadius * math.Cos(pressureAngle)

	// addendum: radial distance from pitch circle to outside circle
	addendum := gearModule * 1.0
	// dedendum: radial distance from pitch circle to root circle
	dedendum := addendum + clearance

	outerRadius := pitchRadius + addendum
	rootRadius := pitchRadius - dedendum
	ringRadius := rootRadius - ringWidth

	tooth := InvoluteGearTooth(
		numberTeeth,
		gearModule,
		rootRadius,
		baseRadius,
		outerRadius,
		backlash,
		facets,
	)

	gear := RotateCopy2D(tooth, numberTeeth)
	root := Circle2D(rootRadius)
	ring := Circle2D(ringRadius)

	return Difference2D(Union2D(gear, root), ring)
}

//-----------------------------------------------------------------------------
// 2D Gear Rack

// GearRackSDF2 is a 2d linear gear rack.
type GearRackSDF2 struct {
	tooth  SDF2    // polygon for rack tooth
	pitch  float64 // tooth to tooth pitch
	length float64 // half the total rack length
	bb     Box2    // bounding box
}

// GearRack2D returns the 2D profile for a gear rack.
func GearRack2D(
	numberTeeth float64, // number of rack teeth
	gearModule float64, // pitch circle diameter / number of gear teeth
	pressureAngle float64, // gear pressure angle (radians)
	backlash float64, // backlash expressed as units of pitch circumference
	baseHeight float64, // height of rack base
) SDF2 {
	s := GearRackSDF2{}

	// addendum: distance from pitch line to top of tooth
	addendum := gearModule * 1.0
	// dedendum: distance from pitch line to root of tooth
	dedendum := gearModule * 1.25
	// total tooth height
	toothHeight := baseHeight + addendum + dedendum
	// tooth_pitch: tooth to tooth distance along pitch line
	pitch := gearModule * PI

	// x size of tooth flank
	dx := (addendum + dedendum) * math.Tan(pressureAngle)
	// 1/2 x size of tooth top
	dxt := ((pitch / 2.0) - dx) / 2.0
	// x size of backlash
	bl := backlash / 2.0

	// create a half tooth profile centered on the y-axis
	tooth := []V2{
		{pitch, 0},
		{pitch, baseHeight},
		{dx + dxt - bl, baseHeight},
		{dxt - bl, toothHeight},
		{-pitch, toothHeight},
		{-pitch, 0},
	}

	s.tooth = Polygon2D(tooth)
	s.pitch = pitch
	s.length = pitch * numberTeeth / 2.0
	s.bb = Box2{V2{-s.length, 0}, V2{s.length, toothHeight}}
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
