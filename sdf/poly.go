//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"math"
)

//-----------------------------------------------------------------------------

// Smoothable 2d Polygon Vertex
type SmoothV2 struct {
	Point  V2
	Facets int
	Radius float64
}

// Set of smoothable 2d polygon points
type Smoother struct {
	Points []SmoothV2
	Closed bool // is the set of points closed?
}

//-----------------------------------------------------------------------------

// Return the next point on the list
func (s *Smoother) next_point(i int) *SmoothV2 {
	if i == len(s.Points)-1 {
		if s.Closed {
			return &s.Points[0]
		} else {
			return nil
		}

	}
	return &s.Points[i+1]
}

// Return the previous point on list
func (s *Smoother) prev_point(i int) *SmoothV2 {
	if i == 0 {
		if s.Closed {
			return &s.Points[len(s.Points)-1]
		} else {
			return nil
		}

	}
	return &s.Points[i-1]
}

// Smooth the i-th point, return true if we smoothed it
func (s *Smoother) smooth_point(i int) bool {

	p := s.Points[i]
	if p.Radius == 0 {
		// fixed point
		return false
	}

	// get the next and previous points
	pn := s.next_point(i)
	pp := s.prev_point(i)
	if pp == nil || pn == nil {
		// can't smooth the endpoints of an open polygon
		return false
	}

	// work out the angle
	v0 := pp.Point.Sub(p.Point).Normalize()
	v1 := pn.Point.Sub(p.Point).Normalize()
	theta := math.Acos(v0.Dot(v1))

	// distance from vertex to circle tangent
	d1 := p.Radius / math.Tan(theta/2.0)
	if d1 > pp.Point.Sub(p.Point).Length() || d1 > pn.Point.Sub(p.Point).Length() {
		// unable to smooth - radius is too large
		return false
	}

	// tangent points
	p0 := p.Point.Add(v0.MulScalar(d1))

	// distance from vertex to circle center
	d2 := p.Radius / math.Sin(theta/2.0)
	// center of circle
	vc := v0.Add(v1).Normalize()
	c := p.Point.Add(vc.MulScalar(d2))

	// rotation angle
	dtheta := Sign(v1.Cross(v0)) * (PI - theta) / float64(p.Facets)
	// rotation matrix
	rm := Rotate(dtheta)
	// radius vector
	rv := p0.Sub(c)

	// work out the new points
	points := make([]SmoothV2, p.Facets+1)
	for j, _ := range points {
		points[j] = SmoothV2{c.Add(rv), 0, 0}
		rv = rm.MulPosition(rv)
	}

	// replace the old point with the new points
	s.Points = append(s.Points[:i], append(points, s.Points[i+1:]...)...)
	return true
}

//-----------------------------------------------------------------------------

func NewSmoother(closed bool) *Smoother {
	return &Smoother{nil, closed}
}

// Add a non-smoothable point to the list
func (s *Smoother) Add(p V2) {
	s.Points = append(s.Points, SmoothV2{p, 0, 0})
}

// Add a smoothable point to the list
func (s *Smoother) AddSmooth(p V2, facets int, radius float64) {
	s.Points = append(s.Points, SmoothV2{p, facets, radius})
}

// Smooth the point list, return the polygon vertex list
func (s *Smoother) Smooth() []V2 {
	// smooth the points
	done := false
	for done == false {
		done = true
		for i, _ := range s.Points {
			if s.smooth_point(i) {
				done = false
			}
		}
	}
	// return the point list
	p := make([]V2, len(s.Points))
	for i, _ := range p {
		p[i] = s.Points[i].Point
	}
	return p
}

//-----------------------------------------------------------------------------
