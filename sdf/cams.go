//-----------------------------------------------------------------------------
/*

Cams

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
)

//-----------------------------------------------------------------------------
// Flat Flank Cams

// FlatFlankCamSDF2 is 2d cam profile.
// The profile is made from a base circle, a smaller nose circle and flat, tangential flanks.
type FlatFlankCamSDF2 struct {
	distance   float64 // center to center circle distance
	baseRadius float64 // radius of base circle
	noseRadius float64 // radius of nose circle
	a          V2      // lower point on flank line
	u          V2      // normalised line vector for flank
	l          float64 // length of flank line
	bb         Box2    // bounding box
}

// FlatFlankCam2D creates a 2D cam profile.
// The profile is made from a base circle, a smaller nose circle and flat, tangential flanks.
// The base circle is centered on the origin.
// The nose circle is located on the positive y axis.
func FlatFlankCam2D(
	distance float64, // circle to circle center distance
	baseRadius float64, // radius of base circle
	noseRadius float64, // radius of nose circle
) SDF2 {
	s := FlatFlankCamSDF2{}
	s.distance = distance
	s.baseRadius = baseRadius
	s.noseRadius = noseRadius
	// work out the flank line
	sin := (baseRadius - noseRadius) / distance
	cos := math.Sqrt(1 - sin*sin)
	// first point on line
	s.a = V2{cos, sin}.MulScalar(baseRadius)
	// second point on line
	b := V2{cos, sin}.MulScalar(noseRadius).Add(V2{0, distance})
	// line information
	u := b.Sub(s.a)
	s.u = u.Normalize()
	s.l = u.Length()
	// work out the bounding box
	s.bb = Box2{V2{-baseRadius, -baseRadius}, V2{baseRadius, distance + noseRadius}}
	return &s
}

// Evaluate returns the minimum distance to the cam.
func (s *FlatFlankCamSDF2) Evaluate(p V2) float64 {
	// we have symmetry about the y-axis
	p = V2{Abs(p.X), p.Y}
	// vector to first point of flank line
	v := p.Sub(s.a)
	// work out the t-parameter of the projection onto the flank line
	t := v.Dot(s.u)
	var d float64
	if t < 0 {
		// the nearest point is on the major circle
		d = p.Length() - s.baseRadius
	} else if t <= s.l {
		// the nearest point is on the flank line
		d = v.Dot(V2{s.u.Y, -s.u.X})
	} else {
		// the nearest point is on the minor circle
		d = p.Sub(V2{0, s.distance}).Length() - s.noseRadius
	}
	return d
}

// BoundingBox returns the bounding box for the cam.
func (s *FlatFlankCamSDF2) BoundingBox() Box2 {
	return s.bb
}

// MakeFlatFlankCam makes a flat flank cam profile from design parameters.
func MakeFlatFlankCam(
	lift float64, // follower lift distance from base circle
	duration float64, // angle over which the follower lifts from the base circle
	maxDiameter float64, // maximum diameter of cam rotation
) (SDF2, error) {

	if maxDiameter <= 0 {
		return nil, fmt.Errorf("maxDiameter <= 0")
	}
	if lift <= 0 {
		return nil, fmt.Errorf("lift <= 0")
	}
	if duration <= 0 || duration >= PI {
		return nil, fmt.Errorf("invalid duration")
	}

	baseRadius := (maxDiameter / 2.0) - lift
	if baseRadius <= 0 {
		return nil, fmt.Errorf("baseRadius <= 0")
	}

	delta := duration / 2.0
	c := math.Cos(delta)
	noseRadius := baseRadius - (lift*c)/(1-c)
	if noseRadius <= 0 {
		return nil, fmt.Errorf("noseRadius <= 0")
	}
	distance := baseRadius + lift - noseRadius
	return FlatFlankCam2D(distance, baseRadius, noseRadius), nil
}

//-----------------------------------------------------------------------------
// Three Arc Cams

// ThreeArcCamSDF2 is 2d cam profile.
// The profile is made from a base circle, a smaller nose circle and circular flank arcs.
type ThreeArcCamSDF2 struct {
	distance    float64 // center to center circle distance
	baseRadius  float64 // radius of base circle
	noseRadius  float64 // radius of nose circle
	flankRadius float64 // radius of flank circle
	flankCenter V2      // center of flank circle (+ve x-axis flank arc)
	thetaBase   float64 // base/flank intersection angle wrt flank center
	thetaNose   float64 // nose/flank intersection angle wrt flank center
	bb          Box2    // bounding box
}

// ThreeArcCam2D creates a 2D cam profile.
// The profile is made from a base circle, a smaller nose circle and circular flank arcs.
// The base circle is centered on the origin.
// The nose circle is located on the positive y axis.
// The flank arcs are tangential to the base and nose circles.
func ThreeArcCam2D(
	distance float64, // circle to circle center distance
	baseRadius float64, // radius of base circle
	noseRadius float64, // radius of nose circle
	flankRadius float64, // radius of flank arc
) SDF2 {
	// check for the minimum size flank radius
	if flankRadius < (baseRadius+distance+noseRadius)/2.0 {
		panic("flankRadius too small")
	}
	s := ThreeArcCamSDF2{}
	s.distance = distance
	s.baseRadius = baseRadius
	s.noseRadius = noseRadius
	s.flankRadius = flankRadius
	// work out the center for the flank radius
	// the flank arc center must lie on the intersection
	// of two circles about the base/nose circles
	r0 := flankRadius - baseRadius
	r1 := flankRadius - noseRadius
	y := ((r0 * r0) - (r1 * r1) + (distance * distance)) / (2.0 * distance)
	x := -math.Sqrt((r0 * r0) - (y * y)) // < 0 result, +ve x-axis flank arc
	s.flankCenter = V2{x, y}
	// work out theta for the intersection of flank arc and base radius
	p := V2{0, 0}.Sub(s.flankCenter)
	s.thetaBase = math.Atan2(p.Y, p.X)
	// work out theta for the intersection of flank arc and nose radius
	p = V2{0, distance}.Sub(s.flankCenter)
	s.thetaNose = math.Atan2(p.Y, p.X)
	// work out the bounding box
	// TODO fix this - it's wrong if the flank radius is small
	s.bb = Box2{V2{-baseRadius, -baseRadius}, V2{baseRadius, distance + noseRadius}}
	return &s
}

// Evaluate returns the minimum distance to the cam.
func (s *ThreeArcCamSDF2) Evaluate(p V2) float64 {
	// we have symmetry about the y-axis
	p0 := V2{Abs(p.X), p.Y}
	// work out the theta angle wrt the flank center
	v := p0.Sub(s.flankCenter)
	t := math.Atan2(v.Y, v.X)
	// work out the minimum distance
	var d float64
	if t < s.thetaBase {
		// the closest point is on the base radius
		d = p0.Length() - s.baseRadius
	} else if t > s.thetaNose {
		// the closest point is on the nose radius
		d = p0.Sub(V2{0, s.distance}).Length() - s.noseRadius
	} else {
		// the closest point is on the flank radius
		d = v.Length() - s.flankRadius
	}
	return d
}

// BoundingBox returns the bounding box for the cam.
func (s *ThreeArcCamSDF2) BoundingBox() Box2 {
	return s.bb
}

// MakeThreeArcCam makes a three arc cam profile from design parameters.
func MakeThreeArcCam(
	lift float64, // follower lift distance from base circle
	duration float64, // angle over which the follower lifts from the base circle
	maxDiameter float64, // maximum diameter of cam rotation
	k float64, // tunable, bigger k = rounder nose, E.g. 1.05
) (SDF2, error) {

	if maxDiameter <= 0 {
		return nil, fmt.Errorf("maxDiameter <= 0")
	}
	if lift <= 0 {
		return nil, fmt.Errorf("lift <= 0")
	}
	if duration <= 0 {
		return nil, fmt.Errorf("invalid duration")
	}
	if k <= 1.0 {
		return nil, fmt.Errorf("invalid k")
	}

	baseRadius := (maxDiameter / 2.0) - lift
	if baseRadius <= 0 {
		return nil, fmt.Errorf("baseRadius <= 0")
	}

	// Given the duration we know where the flank arc intersects the base circle.
	theta := (PI - duration) / 2.0
	p0 := V2{math.Cos(theta), math.Sin(theta)}.MulScalar(baseRadius)
	// This gives us a line back to the flank arc center
	l0 := NewLine2_PV(p0, p0.Negate())

	//The flank arc intersects the y axis above the lift height.
	p1 := V2{0, k * (baseRadius + lift)}

	// The perpendicular bisector of p0 and p1 passes through the flank arc center.
	pMid := p1.Add(p0).MulScalar(0.5)
	u := p1.Sub(p0)
	l1 := NewLine2_PV(pMid, V2{u.Y, -u.X})

	// Intersect to find the flank arc center.
	flankRadius, _, err := l0.Intersect(l1)
	if err != nil {
		return nil, err
	}
	flankCenter := l0.Position(flankRadius)

	// The nose circle is tangential to the flank arcs and the lift line.
	j := baseRadius + lift
	f := flankRadius
	cx := flankCenter.X
	cy := flankCenter.Y
	noseRadius := ((cx * cx) + (cy * cy) - (f * f) + (j * j) - (2 * cy * j)) / (2 * (j - f - cy))

	// distance between base and nose circles
	distance := baseRadius + lift - noseRadius
	return ThreeArcCam2D(distance, baseRadius, noseRadius, flankRadius), nil
}

//-----------------------------------------------------------------------------

// MakeGenevaCam makes 2d profiles for the driver/driven wheels of a geneva cam.
func MakeGenevaCam(
	numSectors int, // number of sectors in the driven wheel
	centerDistance float64, // center to center distance of driver/driven wheels
	driverRadius float64, // radius of lock portion of driver wheel
	drivenRadius float64, // radius of driven wheel
	pinRadius float64, // radius of driver pin
	clearance float64, // pin/slot and wheel/wheel clearance
) (SDF2, SDF2, error) {

	if numSectors < 2 {
		return nil, nil, fmt.Errorf("invalid number of sectors, must be > 2")
	}
	if centerDistance <= 0 ||
		drivenRadius <= 0 ||
		driverRadius <= 0 ||
		pinRadius <= 0 {
		return nil, nil, fmt.Errorf("invalid dimensions, must be > 0")
	}
	if clearance < 0 {
		return nil, nil, fmt.Errorf("invalid clearance, must be >= 0")
	}
	if centerDistance > drivenRadius+driverRadius {
		return nil, nil, fmt.Errorf("center distance is too large")
	}

	// work out the pin offset from the center of the driver wheel
	theta := TAU / (2.0 * float64(numSectors))
	d := centerDistance
	r := drivenRadius
	pinOffset := math.Sqrt((d * d) + (r * r) - (2 * d * r * math.Cos(theta)))

	// driven wheel
	sDriven := Circle2D(drivenRadius - clearance)
	// cutouts for the driver wheel
	s := Circle2D(driverRadius + clearance)
	s = Transform2D(s, Translate2d(V2{centerDistance, 0}))
	s = RotateCopy2D(s, numSectors)
	sDriven = Difference2D(sDriven, s)
	// cutouts for the pin slots
	slotLength := pinOffset + drivenRadius - centerDistance
	s = Line2D(2*slotLength, pinRadius+clearance)
	s = Transform2D(s, Translate2d(V2{drivenRadius, 0}))
	s = RotateCopy2D(s, numSectors)
	s = Transform2D(s, Rotate2d(theta))
	sDriven = Difference2D(sDriven, s)

	// driver wheel
	sDriver := Circle2D(driverRadius - clearance)
	// cutout for the driven wheel
	s = Circle2D(drivenRadius + clearance)
	s = Transform2D(s, Translate2d(V2{centerDistance, 0}))
	sDriver = Difference2D(sDriver, s)
	// driver pin
	s = Circle2D(pinRadius)
	s = Transform2D(s, Translate2d(V2{pinOffset, 0}))
	sDriver = Union2D(sDriver, s)

	return sDriver, sDriven, nil
}

//-----------------------------------------------------------------------------
