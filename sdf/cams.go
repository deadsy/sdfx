//-----------------------------------------------------------------------------
/*

Cams

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------
// Cam Type 1: The cam shape is made from 2 circles and straight line flanks.

type Cam1 struct {
	distance     float64 // center to center circle distance
	major_radius float64 // radius of major circle
	minor_radius float64 // radius of minor circle
	a            V2      // lower point on flank line
	u            V2      // normalised line vector for flank
	l            float64 // length of flank line
	bb           Box2    // bounding box
}

// Create a 2D cam profile.
// The profile is made from 2 circles and straight line flanks.
// The major circle is centered on the origin.
// The minor circle is located on the positive y axis.
// distance = circle to circle center distance
// major_radius = radius of major circle
// minor_radius = radius of minor circle
func NewCam1(distance, major_radius, minor_radius float64) SDF2 {
	s := Cam1{}
	s.distance = distance
	s.major_radius = major_radius
	s.minor_radius = minor_radius
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
	// we have symmetry about the y-axis
	p0 := V2{Abs(p.X), p.Y}
	// vector to first point of flank line
	v := p0.Sub(s.a)
	// work out the t-parameter of the projection onto the flank line
	t := v.Dot(s.u)
	var d float64
	if t < 0 {
		// the nearest point is on the major circle
		d = p0.Length() - s.major_radius
	} else if t <= s.l {
		// the nearest point is on the flank line
		d = v.Dot(V2{s.u.Y, -s.u.X})
	} else {
		// the nearest point is on the minor circle
		d = p0.Sub(V2{0, s.distance}).Length() - s.minor_radius
	}
	return d
}

// Return the bounding box for the cam.
func (s *Cam1) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
