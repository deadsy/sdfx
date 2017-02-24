//-----------------------------------------------------------------------------
/*

Flanges

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

type Flange1 struct {
	distance      float64 // distance from center to side
	center_radius float64 // radius of center circle
	side_radius   float64 // radius of side circle
	a             V2      // center point on flank line
	u             V2      // normalised line vector for flank
	l             float64 // length of flank line
	bb            Box2    // bounding box
}

// Return a flange shape made from a center circle with two side circles.
// distance = distance from center to side
// center_radius = radius of center circle
// side_radius = radius of side circle
func NewFlange1(distance, center_radius, side_radius float64) SDF2 {
	s := Flange1{}
	s.distance = distance
	s.center_radius = center_radius
	s.side_radius = side_radius
	// work out the flank line
	sin := (center_radius - side_radius) / distance
	cos := math.Sqrt(1 - sin*sin)
	// first point on line
	s.a = V2{sin, cos}.MulScalar(center_radius)
	// second point on line
	b := V2{sin, cos}.MulScalar(side_radius).Add(V2{distance, 0})
	// line information
	u := b.Sub(s.a)
	s.u = u.Normalize()
	s.l = u.Length()
	// work out the bounding box
	w := distance + side_radius
	h := center_radius
	s.bb = Box2{V2{-w, -h}, V2{w, h}}
	return &s
}

// Return the minimum distance to the flange.
func (s *Flange1) Evaluate(p V2) float64 {
	// We are symmetrical about the x and y axis.
	// So- only consider the 1st quadrant.
	p = p.Abs()
	// vector to first point of flank line
	v := p.Sub(s.a)
	// work out the t-parameter of the projection onto the flank line
	t := v.Dot(s.u)
	var d float64
	if t < 0 {
		// the nearest point is on the center circle
		d = p.Length() - s.center_radius
	} else if t <= s.l {
		// the nearest point is on the flank line
		d = v.Dot(V2{-s.u.Y, s.u.X})
	} else {
		// the nearest point is on the side circle
		d = p.Sub(V2{s.distance, 0}).Length() - s.side_radius
	}
	return d
}

// Return the bounding box for the flange.
func (s *Flange1) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// MakeBoltCircle2D returns a 2D profile for a flange bolt circle.
func MakeBoltCircle2D(
	hole_radius float64, // radius of bolt holes
	circle_radius float64, // radius of bolt circle
	num_holes int, // number of bolts
) SDF2 {
	s := Circle2D(hole_radius)
	s = Transform2D(s, Translate2d(V2{circle_radius, 0}))
	s = RotateCopy2D(s, num_holes)
	return s
}

// MakeBoltCircle3D returns a 3D object for a flange bolt circle.
func MakeBoltCircle3D(
	hole_depth float64, // depth of bolt holes
	hole_radius float64, // radius of bolt holes
	circle_radius float64, // radius of bolt circle
	num_holes int, // number of bolts
) SDF3 {
	s := MakeBoltCircle2D(hole_radius, circle_radius, num_holes)
	return Extrude3D(s, hole_depth)
}

//-----------------------------------------------------------------------------
