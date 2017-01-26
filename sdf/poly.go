package sdf

// Smoothable 2d Polygon Vertex
type SmoothV2 struct {
	Point  V2
	Facets int
	Radius float64
}

type Smoothable []SmoothV2

func (s Smoothable) Smooth(closed bool) []V2 {
	// Just copying for the time being
	points := make([]V2, len(s))
	for i, _ := range s {
		points[i] = s[i].Point
	}
	return points
}
