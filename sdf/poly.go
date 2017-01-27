package sdf

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

func NewSmoother(closed bool) *Smoother {
	return &Smoother{nil, closed}
}

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
	return false
}

// Add a regular (non-smoothable) point to the list
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
	for i, _ := range s.Points {
		p[i] = s.Points[i].Point
	}
	return p
}
