//-----------------------------------------------------------------------------
/*

2D Spirals

https://math.stackexchange.com/questions/175106/distance-between-point-and-a-spiral

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"math"

	"github.com/deadsy/sdfx/vec/conv"
	"github.com/deadsy/sdfx/vec/p2"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// polarDist2 returns the distance squared between two polar points.
func polarDist2(p0, p1 p2.Vec) float64 {
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
func (s *arcSpiral) theta(radius float64) ([]float64, error) {
	if s.a == 0 {
		if s.k == radius {
			// infinite solutions
			return nil, errors.New("inf")
		}
		// no solutions
		return nil, nil
	}
	if s.n == 1.0 {
		return []float64{(radius - s.k) / s.a}, nil
	}
	return []float64{math.Exp(s.n * math.Log((radius-s.k)/s.a))}, nil
}

//-----------------------------------------------------------------------------

// ArcSpiralSDF2 is a 2d Archimedean spiral.
type ArcSpiralSDF2 struct {
	spiral     arcSpiral
	d          float64 // offset distance
	start, end p2.Vec  // start/end positions
	bb         Box2
}

// ArcSpiral2D returns a 2d Archimedean spiral (r = m*theta + b).
func ArcSpiral2D(
	a, k float64, // r = m*theta + b
	start, end float64, // start/end angle (radians)
	d float64, // offset distance
) (SDF2, error) {

	// sanity checking
	if start == end {
		return nil, errors.New("start == end")
	}
	if a == 0 {
		return nil, errors.New("a == 0")
	}

	s := ArcSpiralSDF2{
		spiral: arcSpiral{a, 1.0, k},
		d:      d,
	}

	// start and end points
	if start > end {
		start, end = end, start
	}
	s.start = p2.Vec{s.spiral.radius(start), start}
	s.end = p2.Vec{s.spiral.radius(end), end}

	// bounding box
	rMax := math.Max(math.Abs(s.spiral.radius(start)), math.Abs(s.spiral.radius(end))) + d
	s.bb = Box2{v2.Vec{-rMax, -rMax}, v2.Vec{rMax, rMax}}
	return &s, nil
}

// Evaluate returns the minimum distance to a 2d Archimedean spiral.
func (s *ArcSpiralSDF2) Evaluate(p v2.Vec) float64 {
	pp := conv.V2ToP2(p)

	// end points
	d2 := math.Min(polarDist2(pp, s.start), polarDist2(pp, s.end))

	thetas, err := s.spiral.theta(pp.R)
	if err == nil {
		for _, theta := range thetas {
			n := math.Round((pp.Theta - theta) / Tau)
			theta = pp.Theta - (Tau * n)

			if theta >= s.start.Theta && theta <= s.end.Theta {
				d2 = math.Min(d2, polarDist2(pp, p2.Vec{s.spiral.radius(theta), theta}))
			} else {

				if theta < s.start.Theta {
					for theta < s.start.Theta {
						theta += Tau
					}
					if theta < s.end.Theta {
						d2 = math.Min(d2, polarDist2(pp, p2.Vec{s.spiral.radius(theta), theta}))
					}
				}

				if theta > s.end.Theta {
					for theta > s.end.Theta {
						theta -= Tau
					}
					if theta > s.start.Theta {
						d2 = math.Min(d2, polarDist2(pp, p2.Vec{s.spiral.radius(theta), theta}))
					}
				}

			}
		}
	}

	return math.Sqrt(d2) - s.d
}

// BoundingBox returns the bounding box of a 2d Archimedean spiral.
func (s *ArcSpiralSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
