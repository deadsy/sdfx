//-----------------------------------------------------------------------------
/*

Cams

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------
// Cam Type 1: The profile is made from 2 circles and straight flanks.

type Cam1 struct {
	distance     float64 // center to center circle distance
	major_radius float64 // radius of major circle
	minor_radius float64 // radius of minor circle
	round        float64 // distance offset
	a            V2      // lower point on flank line
	u            V2      // normalised line vector for flank
	l            float64 // length of flank line
	bb           Box2    // bounding box
}

// Create a 2D cam profile.
// The profile is made from 2 circles and straight flanks.
// The major circle is centered on the origin.
// The minor circle is located on the positive y axis.
// distance = circle to circle center distance
// major_radius = radius of major circle
// minor_radius = radius of minor circle
// round = distance offset
func NewCam1(distance, major_radius, minor_radius, round float64) SDF2 {
	s := Cam1{}

	s.distance = distance
	s.major_radius = major_radius
	s.minor_radius = minor_radius
	s.round = round

	// work out the flank line
	sin := (major_radius - minor_radius) / distance
	cos := math.Sqrt(1 - sin*sin)

	// first point on line
	s.a = V2{cos, sin}.MulScalar(major_radius)
	// second point on line
	b := V2{cos, sin}.MulScalar(minor_radius).Add(V2{0, distance})

	// line information
	u := b.Sub(s.a)
	s.u = u.Normalize()
	s.l = u.Length()

	// work out the bounding box
	s.bb = Box2{V2{-major_radius, -major_radius}, V2{major_radius, distance + minor_radius}}
	return &s

}

// Return the minimum distance to the cam.
func (s *Cam1) Evaluate(p V2) float64 {
	return 0
}

// Return the bounding box for the cam.
func (s *Cam1) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
