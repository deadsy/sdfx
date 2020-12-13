//-----------------------------------------------------------------------------
/*

Linear Gear Rack

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

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
) (SDF2, error) {
	s := GearRackSDF2{}

	// addendum: distance from pitch line to top of tooth
	addendum := gearModule * 1.0
	// dedendum: distance from pitch line to root of tooth
	dedendum := gearModule * 1.25
	// total tooth height
	toothHeight := baseHeight + addendum + dedendum
	// tooth_pitch: tooth to tooth distance along pitch line
	pitch := gearModule * Pi

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
	return &s, nil
}

// Evaluate returns the minimum distance to the gear rack.
func (s *GearRackSDF2) Evaluate(p V2) float64 {
	// map p.X back to the [0,half_pitch) domain
	p0 := V2{math.Abs(SawTooth(p.X, s.pitch)), p.Y}
	// get the tooth profile distance
	d0 := s.tooth.Evaluate(p0)
	// create a region for the rack length
	d1 := math.Abs(p.X) - s.length
	// return the intersection
	return math.Max(d0, d1)
}

// BoundingBox returns the bounding box for the gear rack.
func (s *GearRackSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
