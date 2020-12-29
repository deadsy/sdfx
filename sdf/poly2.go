//-----------------------------------------------------------------------------
/*

SDF for 2D polygons.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
)

//-----------------------------------------------------------------------------

// PolySDF2 is an SDF2 made from a closed set of line segments.
type PolySDF2 struct {
	vertex []V2      // vertices
	vector []V2      // unit line vectors
	length []float64 // line lengths
	bb     Box2      // bounding box
}

// Polygon2D returns an SDF2 made from a closed set of line segments.
func Polygon2D(vertex []V2) (SDF2, error) {
	s := PolySDF2{}

	n := len(vertex)
	if n < 3 {
		return nil, ErrMsg("number of vertices < 3")
	}

	// Close the loop (if necessary)
	s.vertex = vertex
	if !vertex[0].Equals(vertex[n-1], tolerance) {
		s.vertex = append(s.vertex, vertex[0])
	}

	// allocate pre-calculated line segment info
	nsegs := len(s.vertex) - 1
	s.vector = make([]V2, nsegs)
	s.length = make([]float64, nsegs)

	vmin := s.vertex[0]
	vmax := s.vertex[0]

	for i := 0; i < nsegs; i++ {
		l := s.vertex[i+1].Sub(s.vertex[i])
		s.length[i] = l.Length()
		s.vector[i] = l.Normalize()
		vmin = vmin.Min(s.vertex[i])
		vmax = vmax.Max(s.vertex[i])
	}

	s.bb = Box2{vmin, vmax}
	return &s, nil
}

// Evaluate returns the minimum distance for a 2d polygon.
func (s *PolySDF2) Evaluate(p V2) float64 {
	dd := math.MaxFloat64 // d^2 to polygon (>0)
	wn := 0               // winding number (inside/outside)

	// iterate over the line segments
	nsegs := len(s.vertex) - 1
	pb := p.Sub(s.vertex[0])

	for i := 0; i < nsegs; i++ {
		a := s.vertex[i]
		b := s.vertex[i+1]

		pa := pb
		pb = p.Sub(b)

		t := pa.Dot(s.vector[i])                        // t-parameter of projection onto line
		dn := pa.Dot(V2{s.vector[i].Y, -s.vector[i].X}) // normal distance from p to line

		// Distance to line segment
		if t < 0 {
			dd = math.Min(dd, pa.Length2()) // distance to vertex[0] of line
		} else if t > s.length[i] {
			dd = math.Min(dd, pb.Length2()) // distance to vertex[1] of line
		} else {
			dd = math.Min(dd, dn*dn) // normal distance to line
		}

		// Is the point in the polygon?
		// See: http://geomalgorithms.com/a03-_inclusion.html
		if a.Y <= p.Y {
			if b.Y > p.Y { // upward crossing
				if dn < 0 { // p is to the left of the line segment
					wn++ // up intersect
				}
			}
		} else {
			if b.Y <= p.Y { // downward crossing
				if dn > 0 { // p is to the right of the line segment
					wn-- // down intersect
				}
			}
		}
	}

	// normalise d*d to d
	d := math.Sqrt(dd)
	if wn != 0 {
		// p is inside the polygon
		return -d
	}
	return d
}

// BoundingBox returns the bounding box of a 2d polygon.
func (s *PolySDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
