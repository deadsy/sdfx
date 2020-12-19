//-----------------------------------------------------------------------------
/*

Linear Gear Rack

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------
// 2D Gear Rack

// GearRackParms defines the parameters for a gear rack.
type GearRackParms struct {
	NumberTeeth   int     // number of rack teeth
	Module        float64 // pitch circle diameter / number of gear teeth
	PressureAngle float64 // gear pressure angle (radians)
	Backlash      float64 // backlash expressed as units of pitch circumference
	BaseHeight    float64 // height of rack base
}

// GearRackSDF2 is a 2d linear gear rack.
type GearRackSDF2 struct {
	tooth  SDF2    // polygon for rack tooth
	pitch  float64 // tooth to tooth pitch
	length float64 // half the total rack length
	bb     Box2    // bounding box
}

// GearRack2D returns the 2D profile for a gear rack.
func GearRack2D(k *GearRackParms) (SDF2, error) {

	if k.NumberTeeth <= 0 {
		return nil, ErrMsg("NumberTeeth <= 0")
	}
	if k.Module <= 0 {
		return nil, ErrMsg("Module <= 0")
	}
	if k.PressureAngle <= 0 {
		return nil, ErrMsg("PressureAngle <= 0")
	}
	if k.Backlash < 0 {
		return nil, ErrMsg("Backlash <= 0")
	}
	if k.BaseHeight < 0 {
		return nil, ErrMsg("BaseHeight < 0")
	}

	s := GearRackSDF2{}

	// addendum: distance from pitch line to top of tooth
	addendum := k.Module * 1.0
	// dedendum: distance from pitch line to root of tooth
	dedendum := k.Module * 1.25
	// total tooth height
	toothHeight := k.BaseHeight + addendum + dedendum
	// tooth_pitch: tooth to tooth distance along pitch line
	pitch := k.Module * Pi

	// x size of tooth flank
	dx := (addendum + dedendum) * math.Tan(k.PressureAngle)
	// 1/2 x size of tooth top
	dxt := ((pitch / 2.0) - dx) / 2.0
	// x size of backlash
	bl := k.Backlash / 2.0

	// create a half tooth profile centered on the y-axis
	tooth := []V2{
		{pitch, 0},
		{pitch, k.BaseHeight},
		{dx + dxt - bl, k.BaseHeight},
		{dxt - bl, toothHeight},
		{-pitch, toothHeight},
		{-pitch, 0},
	}
	tp, err := Polygon2D(tooth)
	if err != nil {
		return nil, err
	}

	s.tooth = tp
	s.pitch = pitch
	s.length = pitch * float64(k.NumberTeeth) * 0.5
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
