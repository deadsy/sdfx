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
// Flat Flank Cams.

type FlatFlankCam struct {
	distance    float64 // center to center circle distance
	base_radius float64 // radius of base circle
	nose_radius float64 // radius of nose circle
	a           V2      // lower point on flank line
	u           V2      // normalised line vector for flank
	l           float64 // length of flank line
	bb          Box2    // bounding box
}

// Create a 2D cam profile.
// The profile is made from 2 circles and straight line flanks.
// The base circle is centered on the origin.
// The nose circle is located on the positive y axis.
// distance = circle to circle center distance
// base_radius = radius of base circle
// nose_radius = radius of nose circle
func FlatFlankCam2D(distance, base_radius, nose_radius float64) SDF2 {
	s := FlatFlankCam{}
	s.distance = distance
	s.base_radius = base_radius
	s.nose_radius = nose_radius
	// work out the flank line
	sin := (base_radius - nose_radius) / distance
	cos := math.Sqrt(1 - sin*sin)
	// first point on line
	s.a = V2{cos, sin}.MulScalar(base_radius)
	// second point on line
	b := V2{cos, sin}.MulScalar(nose_radius).Add(V2{0, distance})
	// line information
	u := b.Sub(s.a)
	s.u = u.Normalize()
	s.l = u.Length()
	// work out the bounding box
	s.bb = Box2{V2{-base_radius, -base_radius}, V2{base_radius, distance + nose_radius}}
	return &s
}

// Return the minimum distance to the cam.
func (s *FlatFlankCam) Evaluate(p V2) float64 {
	// we have symmetry about the y-axis
	p = V2{Abs(p.X), p.Y}
	// vector to first point of flank line
	v := p.Sub(s.a)
	// work out the t-parameter of the projection onto the flank line
	t := v.Dot(s.u)
	var d float64
	if t < 0 {
		// the nearest point is on the major circle
		d = p.Length() - s.base_radius
	} else if t <= s.l {
		// the nearest point is on the flank line
		d = v.Dot(V2{s.u.Y, -s.u.X})
	} else {
		// the nearest point is on the minor circle
		d = p.Sub(V2{0, s.distance}).Length() - s.nose_radius
	}
	return d
}

// Return the bounding box for the cam.
func (s *FlatFlankCam) BoundingBox() Box2 {
	return s.bb
}

// Create a flat flank cam profile from design parameters.
// lift = follower lift distance from base circle
// duration = angle over which the follower lifts from the base circle
// max_diameter = maximum diameter of cam rotation
func MakeFlatFlankCam(lift, duration, max_diameter float64) (SDF2, error) {

	if max_diameter <= 0 {
		return nil, fmt.Errorf("max_diameter <= 0")
	}
	if lift <= 0 {
		return nil, fmt.Errorf("lift <= 0")
	}
	if duration <= 0 || duration >= PI {
		return nil, fmt.Errorf("invalid duration")
	}

	base_radius := (max_diameter / 2.0) - lift
	if base_radius <= 0 {
		return nil, fmt.Errorf("base_radius <= 0")
	}

	delta := duration / 2.0
	c := math.Cos(delta)
	nose_radius := base_radius - (lift*c)/(1-c)
	if nose_radius <= 0 {
		return nil, fmt.Errorf("nose_radius <= 0")
	}
	distance := base_radius + lift - nose_radius
	return FlatFlankCam2D(distance, base_radius, nose_radius), nil
}

//-----------------------------------------------------------------------------
// Three Arc Cams.

type ThreeArcCam struct {
	distance     float64 // center to center circle distance
	base_radius  float64 // radius of base circle
	nose_radius  float64 // radius of nose circle
	flank_radius float64 // radius of flank circle
	flank_center V2      // center of flank circle (+ve x-axis flank arc)
	theta_base   float64 // base/flank intersection angle wrt flank center
	theta_nose   float64 // nose/flank intersection angle wrt flank center
	bb           Box2    // bounding box
}

// Create a 2D cam profile.
// The profile is made from 2 circles and circular flank arcs.
// The base circle is centered on the origin.
// The nose circle is located on the positive y axis.
// The flank arcs are tangential to the base and nose circles.
// distance = circle to circle center distance
// base_radius = radius of major circle
// nose_radius = radius of minor circle
// flank_radius = radius of flank arc
func ThreeArcCam2D(distance, base_radius, nose_radius, flank_radius float64) SDF2 {
	// check for the minimum size flank radius
	if flank_radius < (base_radius+distance+nose_radius)/2.0 {
		panic("flank_radius too small")
	}
	s := ThreeArcCam{}
	s.distance = distance
	s.base_radius = base_radius
	s.nose_radius = nose_radius
	s.flank_radius = flank_radius
	// work out the center for the flank radius
	// the flank arc center must lie on the intersection
	// of two circles about the base/nose circles
	r0 := flank_radius - base_radius
	r1 := flank_radius - nose_radius
	y := ((r0 * r0) - (r1 * r1) + (distance * distance)) / (2.0 * distance)
	x := -math.Sqrt((r0 * r0) - (y * y)) // < 0 result, +ve x-axis flank arc
	s.flank_center = V2{x, y}
	// work out theta for the intersection of flank arc and base radius
	p := V2{0, 0}.Sub(s.flank_center)
	s.theta_base = math.Atan2(p.Y, p.X)
	// work out theta for the intersection of flank arc and nose radius
	p = V2{0, distance}.Sub(s.flank_center)
	s.theta_nose = math.Atan2(p.Y, p.X)
	// work out the bounding box
	// TODO fix this - it's wrong if the flank radius is small
	s.bb = Box2{V2{-base_radius, -base_radius}, V2{base_radius, distance + nose_radius}}
	return &s
}

// Return the minimum distance to the cam.
func (s *ThreeArcCam) Evaluate(p V2) float64 {
	// we have symmetry about the y-axis
	p0 := V2{Abs(p.X), p.Y}
	// work out the theta angle wrt the flank center
	v := p0.Sub(s.flank_center)
	t := math.Atan2(v.Y, v.X)
	// work out the minimum distance
	var d float64
	if t < s.theta_base {
		// the closest point is on the base radius
		d = p0.Length() - s.base_radius
	} else if t > s.theta_nose {
		// the closest point is on the nose radius
		d = p0.Sub(V2{0, s.distance}).Length() - s.nose_radius
	} else {
		// the closest point is on the flank radius
		d = v.Length() - s.flank_radius
	}
	return d
}

// Return the bounding box for the cam.
func (s *ThreeArcCam) BoundingBox() Box2 {
	return s.bb
}

// Create a three arc cam profile from design parameters.
// lift = follower lift distance from base circle
// duration = angle over which the follower lifts from the base circle
// max_diameter = maximum diameter of cam rotation
// k = tunable, bigger k = rounder nose, E.g. 1.05
func MakeThreeArcCam(lift, duration, max_diameter, k float64) (SDF2, error) {

	if max_diameter <= 0 {
		return nil, fmt.Errorf("max_diameter <= 0")
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

	base_radius := (max_diameter / 2.0) - lift
	if base_radius <= 0 {
		return nil, fmt.Errorf("base_radius <= 0")
	}

	// Given the duration we know where the flank arc intersects the base circle.
	theta := (PI - duration) / 2.0
	p0 := V2{math.Cos(theta), math.Sin(theta)}.MulScalar(base_radius)
	// This gives us a line back to the flank arc center
	l0 := NewLine2_PV(p0, p0.Negate())

	//The flank arc intersects the y axis above the lift height.
	p1 := V2{0, k * (base_radius + lift)}

	// The perpendicular bisector of p0 and p1 passes through the flank arc center.
	p_mid := p1.Add(p0).MulScalar(0.5)
	u := p1.Sub(p0)
	l1 := NewLine2_PV(p_mid, V2{u.Y, -u.X})

	// Intersect to find the flank arc center.
	flank_radius, _, err := l0.Intersect(l1)
	if err != nil {
		return nil, err
	}
	flank_center := l0.Position(flank_radius)

	// The nose circle is tangential to the flank arcs and the lift line.
	j := base_radius + lift
	f := flank_radius
	cx := flank_center.X
	cy := flank_center.Y
	nose_radius := ((cx * cx) + (cy * cy) - (f * f) + (j * j) - (2 * cy * j)) / (2 * (j - f - cy))

	// distance between base and nose circles
	distance := base_radius + lift - nose_radius
	return ThreeArcCam2D(distance, base_radius, nose_radius, flank_radius), nil
}

//-----------------------------------------------------------------------------

// Make 2D profiles for the driver/driven wheels of a geneva cam.
// num_sectors = number of sectors in the driven wheel
// center_distance = center to center distance of driver/driven wheels
// driver_radius = radius of lock portion of driver wheel
// driven_radius = radius of driven wheel
// pin_radius = radius of driver pin
// clearance = pin/slot and wheel/wheel clearance
func MakeGenevaCam(
	num_sectors int,
	center_distance float64,
	driver_radius float64,
	driven_radius float64,
	pin_radius float64,
	clearance float64,
) (SDF2, SDF2, error) {

	if num_sectors < 2 {
		return nil, nil, fmt.Errorf("invalid number of sectors, must be > 2")
	}
	if center_distance <= 0 ||
		driven_radius <= 0 ||
		driver_radius <= 0 ||
		pin_radius <= 0 {
		return nil, nil, fmt.Errorf("invalid dimensions, must be > 0")
	}
	if clearance < 0 {
		return nil, nil, fmt.Errorf("invalid clearance, must be >= 0")
	}
	if center_distance > driven_radius+driver_radius {
		return nil, nil, fmt.Errorf("center distance is too large")
	}

	// work out the pin offset from the center of the driver wheel
	theta := TAU / (2.0 * float64(num_sectors))
	d := center_distance
	r := driven_radius
	pin_offset := math.Sqrt((d * d) + (r * r) - (2 * d * r * math.Cos(theta)))

	// driven wheel
	s_driven := Circle2D(driven_radius - clearance)
	// cutouts for the driver wheel
	s := Circle2D(driver_radius + clearance)
	s = Transform2D(s, Translate2d(V2{center_distance, 0}))
	s = RotateCopy2D(s, num_sectors)
	s_driven = Difference2D(s_driven, s)
	// cutouts for the pin slots
	slot_length := pin_offset + driven_radius - center_distance
	s = Line2D(2*slot_length, pin_radius+clearance)
	s = Transform2D(s, Translate2d(V2{driven_radius, 0}))
	s = RotateCopy2D(s, num_sectors)
	s = Transform2D(s, Rotate2d(theta))
	s_driven = Difference2D(s_driven, s)

	// driver wheel
	s_driver := Circle2D(driver_radius - clearance)
	// cutout for the driven wheel
	s = Circle2D(driven_radius + clearance)
	s = Transform2D(s, Translate2d(V2{center_distance, 0}))
	s_driver = Difference2D(s_driver, s)
	// driver pin
	s = Circle2D(pin_radius)
	s = Transform2D(s, Translate2d(V2{pin_offset, 0}))
	s_driver = Union2D(s_driver, s)

	return s_driver, s_driven, nil
}

//-----------------------------------------------------------------------------
