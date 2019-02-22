//-----------------------------------------------------------------------------
/*

2D Spirals

https://math.stackexchange.com/questions/175106/distance-between-point-and-a-spiral

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
)

//-----------------------------------------------------------------------------

// polarDist2 returns the distance squared between two polar points.
func polarDist2(p0, p1 P2) float64 {
	return (p0.R * p0.R) + (p1.R * p1.R) - 2.0*p0.R*p1.R*math.Cos(p0.Theta-p1.Theta)
}

//-----------------------------------------------------------------------------

// arcSpiral is an archimedean spiral.
type arcSpiral struct {
	a, n, k float64 // r = a * pow(theta, 1/n) + k
}

// radius returns the radius for a given theta.
func (s *arcSpiral) radius(theta float64) float64 {
	var r float64
	if s.a == 0 {
		r = s.k
	} else {
		if s.n == 1.0 {
			r = s.a*theta + s.k
		} else {
			r = math.Pow(theta, 1.0/s.n) + s.k
		}
	}
	return r
}

// theta returns the theta(s) for a given radius.
func (s *arcSpiral) theta(radius float64) (int, []float64) {
	if s.a == 0 {
		if s.k == radius {
			// infinite solutions
			return -1, nil
		}
		// no solutions
		return 0, nil
	}
	if s.n == 1.0 {
		return 1, []float64{(radius - s.k) / s.a}
	}
	return 1, []float64{math.Exp(s.n * math.Log((radius-s.k)/s.a))}
}

//-----------------------------------------------------------------------------

// ArcSpiralSDF2 is a 2d Archimedean spiral.
type ArcSpiralSDF2 struct {
	m, b       float64 // r = m*theta + b
	d          float64 // offset distance
	start, end P2      // start/end positions
	bb         Box2
}

// radius returns the spiral radius for a given theta value.
func (s *ArcSpiralSDF2) radius(theta float64) float64 {
	return (s.m * theta) + s.b
}

// theta returns the spiral theta for a given radius value.
func (s *ArcSpiralSDF2) theta(radius float64) float64 {
	return (radius - s.b) / s.m
}

// ArcSpiral2D returns a 2d Archimedean spiral (r = m*theta + b).
func ArcSpiral2D(
	m, b float64, // r = m*theta + b
	start, end float64, // start/end angle (radians)
	d float64, // offset distance
) SDF2 {

	// sanity checking
	if start == end {
		panic("start == end")
	}
	if m == 0 {
		panic("m == 0")
	}

	s := ArcSpiralSDF2{
		m: m,
		b: b,
		d: d,
	}

	// start and end points
	if start <= end {
		s.start = P2{s.radius(start), start}
		s.end = P2{s.radius(end), end}
	} else {
		s.end = P2{s.radius(start), start}
		s.start = P2{s.radius(end), end}
	}

	// bounding box
	rMax := Max(Abs(s.radius(start)), Abs(s.radius(end))) + d
	s.bb = Box2{V2{-rMax, -rMax}, V2{rMax, rMax}}
	return &s
}

// Evaluate returns the minimum distance to a 2d Archimedean spiral.
func (s *ArcSpiralSDF2) Evaluate(p V2) float64 {
	pp := p.CartesianToPolar()

	// end points
	d2 := Min(polarDist2(pp, s.start), polarDist2(pp, s.end))

	// positive radius
	sTheta := s.theta(pp.R)
	n := math.Round((pp.Theta - sTheta) / Tau)
	sTheta = pp.Theta - (Tau * n)
	if sTheta > s.start.Theta && sTheta < s.end.Theta {
		d2 = Min(d2, polarDist2(pp, P2{s.radius(sTheta), sTheta}))
	}

	// negative radius
	sTheta = s.theta(-pp.R)
	n = math.Round((pp.Theta - sTheta) / Tau)
	sTheta = pp.Theta - (Tau * n)
	if sTheta > s.start.Theta && sTheta < s.end.Theta {
		d2 = Min(d2, polarDist2(pp, P2{s.radius(sTheta), sTheta}))
	}

	return math.Sqrt(d2) - s.d
}

// BoundingBox returns the bounding box of a 2d Archimedean spiral.
func (s *ArcSpiralSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
