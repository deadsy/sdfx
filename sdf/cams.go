//-----------------------------------------------------------------------------
/*

Cams

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------
// Flat Flank Cams

// FlatFlankCamSDF2 is 2d cam profile.
// The profile is made from a base circle, a smaller nose circle and flat, tangential flanks.
type FlatFlankCamSDF2 struct {
	distance   float64 // center to center circle distance
	baseRadius float64 // radius of base circle
	noseRadius float64 // radius of nose circle
	a          v2.Vec  // lower point on flank line
	u          v2.Vec  // normalised line vector for flank
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
) (SDF2, error) {
	s := FlatFlankCamSDF2{}
	s.distance = distance
	s.baseRadius = baseRadius
	s.noseRadius = noseRadius
	// work out the flank line
	sin := (baseRadius - noseRadius) / distance
	cos := math.Sqrt(1 - sin*sin)
	// first point on line
	s.a = v2.Vec{cos, sin}.MulScalar(baseRadius)
	// second point on line
	b := v2.Vec{cos, sin}.MulScalar(noseRadius).Add(v2.Vec{0, distance})
	// line information
	u := b.Sub(s.a)
	s.u = u.Normalize()
	s.l = u.Length()
	// work out the bounding box
	s.bb = Box2{v2.Vec{-baseRadius, -baseRadius}, v2.Vec{baseRadius, distance + noseRadius}}
	return &s, nil
}

// Evaluate returns the minimum distance to the cam.
func (s *FlatFlankCamSDF2) Evaluate(p v2.Vec) float64 {
	// we have symmetry about the y-axis
	p = v2.Vec{math.Abs(p.X), p.Y}
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
		d = v.Dot(v2.Vec{s.u.Y, -s.u.X})
	} else {
		// the nearest point is on the minor circle
		d = p.Sub(v2.Vec{0, s.distance}).Length() - s.noseRadius
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
		return nil, errors.New("maxDiameter <= 0")
	}
	if lift <= 0 {
		return nil, errors.New("lift <= 0")
	}
	if duration <= 0 || duration >= Pi {
		return nil, errors.New("invalid duration")
	}

	baseRadius := (maxDiameter / 2.0) - lift
	if baseRadius <= 0 {
		return nil, errors.New("baseRadius <= 0")
	}

	delta := duration / 2.0
	c := math.Cos(delta)
	noseRadius := baseRadius - (lift*c)/(1-c)
	if noseRadius <= 0 {
		return nil, errors.New("noseRadius <= 0")
	}
	distance := baseRadius + lift - noseRadius
	return FlatFlankCam2D(distance, baseRadius, noseRadius)
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
	flankCenter v2.Vec  // center of flank circle (+ve x-axis flank arc)
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
) (SDF2, error) {
	// check for the minimum size flank radius
	if flankRadius < (baseRadius+distance+noseRadius)/2.0 {
		return nil, errors.New("flankRadius too small")
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
	s.flankCenter = v2.Vec{x, y}
	// work out theta for the intersection of flank arc and base radius
	p := v2.Vec{0, 0}.Sub(s.flankCenter)
	s.thetaBase = math.Atan2(p.Y, p.X)
	// work out theta for the intersection of flank arc and nose radius
	p = v2.Vec{0, distance}.Sub(s.flankCenter)
	s.thetaNose = math.Atan2(p.Y, p.X)
	// work out the bounding box
	// TODO fix this - it's wrong if the flank radius is small
	s.bb = Box2{v2.Vec{-baseRadius, -baseRadius}, v2.Vec{baseRadius, distance + noseRadius}}
	return &s, nil
}

// Evaluate returns the minimum distance to the cam.
func (s *ThreeArcCamSDF2) Evaluate(p v2.Vec) float64 {
	// we have symmetry about the y-axis
	p0 := v2.Vec{math.Abs(p.X), p.Y}
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
		d = p0.Sub(v2.Vec{0, s.distance}).Length() - s.noseRadius
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
		return nil, errors.New("maxDiameter <= 0")
	}
	if lift <= 0 {
		return nil, errors.New("lift <= 0")
	}
	if duration <= 0 {
		return nil, errors.New("invalid duration")
	}
	if k <= 1.0 {
		return nil, errors.New("invalid k")
	}

	baseRadius := (maxDiameter / 2.0) - lift
	if baseRadius <= 0 {
		return nil, errors.New("baseRadius <= 0")
	}

	// Given the duration we know where the flank arc intersects the base circle.
	theta := (Pi - duration) / 2.0
	p0 := v2.Vec{math.Cos(theta), math.Sin(theta)}.MulScalar(baseRadius)
	// This gives us a line back to the flank arc center
	l0 := newLinePV(p0, p0.Neg())

	//The flank arc intersects the y axis above the lift height.
	p1 := v2.Vec{0, k * (baseRadius + lift)}

	// The perpendicular bisector of p0 and p1 passes through the flank arc center.
	pMid := p1.Add(p0).MulScalar(0.5)
	u := p1.Sub(p0)
	l1 := newLinePV(pMid, v2.Vec{u.Y, -u.X})

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
	return ThreeArcCam2D(distance, baseRadius, noseRadius, flankRadius)
}

//-----------------------------------------------------------------------------
