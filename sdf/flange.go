//-----------------------------------------------------------------------------
/*

Flanges

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// Flange1 is a flange shape made from a center circle with two side circles.
type Flange1 struct {
	distance     float64 // distance from center to side
	centerRadius float64 // radius of center circle
	sideRadius   float64 // radius of side circle
	a            v2.Vec  // center point on flank line
	u            v2.Vec  // normalised line vector for flank
	l            float64 // length of flank line
	bb           Box2    // bounding box
}

// NewFlange1 returns a flange shape made from a center circle with two side circles.
func NewFlange1(
	distance float64, // distance from center to side circle
	centerRadius float64, // radius of center circle
	sideRadius float64, // radius of side circle
) SDF2 {
	s := Flange1{}
	s.distance = distance
	s.centerRadius = centerRadius
	s.sideRadius = sideRadius
	// work out the flank line
	sin := (centerRadius - sideRadius) / distance
	cos := math.Sqrt(1 - sin*sin)
	// first point on line
	s.a = v2.Vec{sin, cos}.MulScalar(centerRadius)
	// second point on line
	b := v2.Vec{sin, cos}.MulScalar(sideRadius).Add(v2.Vec{distance, 0})
	// line information
	u := b.Sub(s.a)
	s.u = u.Normalize()
	s.l = u.Length()
	// work out the bounding box
	w := distance + sideRadius
	h := centerRadius
	s.bb = Box2{v2.Vec{-w, -h}, v2.Vec{w, h}}
	return &s
}

// Evaluate returns the minimum distance to the flange.
func (s *Flange1) Evaluate(p v2.Vec) float64 {
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
		d = p.Length() - s.centerRadius
	} else if t <= s.l {
		// the nearest point is on the flank line
		d = v.Dot(v2.Vec{-s.u.Y, s.u.X})
	} else {
		// the nearest point is on the side circle
		d = p.Sub(v2.Vec{s.distance, 0}).Length() - s.sideRadius
	}
	return d
}

// BoundingBox returns the bounding box for the flange.
func (s *Flange1) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
