//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"math"
)

//-----------------------------------------------------------------------------

// Smoothable 2d polygon vertex
type SmoothV2 struct {
	Vertex V2      // vertex coordinates
	Facets int     // number of polygon facets to create when smoothing
	Radius float64 // radius of smoothing (0 == none)
}

// Set of smoothable 2d polygon vertices
type Smoother struct {
	VList  []SmoothV2 // vertex list
	Closed bool       // is the set of points closed?
}

//-----------------------------------------------------------------------------

// Return the next vertex on the list
func (s *Smoother) next_vertex(i int) *SmoothV2 {
	if i == len(s.VList)-1 {
		if s.Closed {
			return &s.VList[0]
		} else {
			return nil
		}
	}
	return &s.VList[i+1]
}

// Return the previous vertex on list
func (s *Smoother) prev_vertex(i int) *SmoothV2 {
	if i == 0 {
		if s.Closed {
			return &s.VList[len(s.VList)-1]
		} else {
			return nil
		}
	}
	return &s.VList[i-1]
}

// Smooth the i-th vertex, return true if we smoothed it
func (s *Smoother) smooth_vertex(i int) bool {

	p := s.VList[i]
	if p.Radius == 0 {
		// fixed point
		return false
	}

	// get the next and previous points
	pn := s.next_vertex(i)
	pp := s.prev_vertex(i)
	if pp == nil || pn == nil {
		// can't smooth the endpoints of an open polygon
		return false
	}

	// work out the angle
	v0 := pp.Vertex.Sub(p.Vertex).Normalize()
	v1 := pn.Vertex.Sub(p.Vertex).Normalize()
	theta := math.Acos(v0.Dot(v1))

	// distance from vertex to circle tangent
	d1 := p.Radius / math.Tan(theta/2.0)
	if d1 > pp.Vertex.Sub(p.Vertex).Length() || d1 > pn.Vertex.Sub(p.Vertex).Length() {
		// unable to smooth - radius is too large
		return false
	}

	// tangent points
	p0 := p.Vertex.Add(v0.MulScalar(d1))

	// distance from vertex to circle center
	d2 := p.Radius / math.Sin(theta/2.0)
	// center of circle
	vc := v0.Add(v1).Normalize()
	c := p.Vertex.Add(vc.MulScalar(d2))

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
	s.VList = append(s.VList[:i], append(points, s.VList[i+1:]...)...)
	return true
}

//-----------------------------------------------------------------------------

func NewSmoother(closed bool) *Smoother {
	return &Smoother{nil, closed}
}

// Add a non-smoothable vertex to the list
func (s *Smoother) Add(p V2) {
	s.VList = append(s.VList, SmoothV2{p, 0, 0})
}

// Add a smoothable vertex to the list
func (s *Smoother) AddSmooth(p V2, facets int, radius float64) {
	s.VList = append(s.VList, SmoothV2{p, facets, radius})
}

// Smooth the vertex list
func (s *Smoother) Smooth() {
	// smooth the points
	done := false
	for done == false {
		done = true
		for i, _ := range s.VList {
			if s.smooth_vertex(i) {
				done = false
			}
		}
	}
}

// Return a list of the smoother vertices
func (s *Smoother) Vertices() []V2 {
	// return the vertex list
	vlist := make([]V2, len(s.VList))
	for i, _ := range vlist {
		vlist[i] = s.VList[i].Vertex
	}
	return vlist
}

//-----------------------------------------------------------------------------

// Return the vertices of a N sided regular polygon
func Nagon(n int, radius float64) V2Set {
	if n < 3 {
		return nil
	}
	m := Rotate(TAU / float64(n))
	v := make(V2Set, n)
	p := V2{radius, 0}
	for i := 0; i < n; i++ {
		v[i] = p
		p = m.MulPosition(p)
	}
	return v
}

//-----------------------------------------------------------------------------
