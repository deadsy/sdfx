//-----------------------------------------------------------------------------
/*

Flanges

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

// Flange1 is a flange shape made from a center circle with two side circles.
type Flange1 struct {
	distance     float64 // distance from center to side
	centerRadius float64 // radius of center circle
	sideRadius   float64 // radius of side circle
	a            V2      // center point on flank line
	u            V2      // normalised line vector for flank
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
	s.a = V2{sin, cos}.MulScalar(centerRadius)
	// second point on line
	b := V2{sin, cos}.MulScalar(sideRadius).Add(V2{distance, 0})
	// line information
	u := b.Sub(s.a)
	s.u = u.Normalize()
	s.l = u.Length()
	// work out the bounding box
	w := distance + sideRadius
	h := centerRadius
	s.bb = Box2{V2{-w, -h}, V2{w, h}}
	return &s
}

// Evaluate returns the minimum distance to the flange.
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
		d = p.Length() - s.centerRadius
	} else if t <= s.l {
		// the nearest point is on the flank line
		d = v.Dot(V2{-s.u.Y, s.u.X})
	} else {
		// the nearest point is on the side circle
		d = p.Sub(V2{s.distance, 0}).Length() - s.sideRadius
	}
	return d
}

// BoundingBox returns the bounding box for the flange.
func (s *Flange1) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// MakeBoltCircle2D returns a 2D profile for a flange bolt circle.
func MakeBoltCircle2D(
	holeRadius float64, // radius of bolt holes
	circleRadius float64, // radius of bolt circle
	numHoles int, // number of bolts
) SDF2 {
	s := Circle2D(holeRadius)
	s = Transform2D(s, Translate2d(V2{circleRadius, 0}))
	s = RotateCopy2D(s, numHoles)
	return s
}

// MakeBoltCircle3D returns a 3D object for a flange bolt circle.
func MakeBoltCircle3D(
	holeDepth float64, // depth of bolt holes
	holeRadius float64, // radius of bolt holes
	circleRadius float64, // radius of bolt circle
	numHoles int, // number of bolts
) SDF3 {
	s := MakeBoltCircle2D(holeRadius, circleRadius, numHoles)
	return Extrude3D(s, holeDepth)
}

//-----------------------------------------------------------------------------
