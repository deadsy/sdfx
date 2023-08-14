//-----------------------------------------------------------------------------
/*

2D Evaluation Cache

In some cases (E.g. extrusion) 2d SDFs get evaluated repeatedly at the same points.
If a map lookup is cheaper than a distance evaluation it's possible to save time
by caching evaluation results. This SDF2 wraps an underlying SDF2 and caches the
evaluations.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// CacheSDF2 is an SDF2 cache.
type CacheSDF2 struct {
	sdf         SDF2
	cache       map[v2.Vec]float64
	reads, hits uint
}

// Cache2D wraps the passed SDF2 with an evaluation cache.
func Cache2D(sdf SDF2) SDF2 {
	return &CacheSDF2{
		sdf:   sdf,
		cache: make(map[v2.Vec]float64),
	}
}

func (s *CacheSDF2) String() string {
	r := float64(s.hits) / float64(s.reads)
	return fmt.Sprintf("reads %d hits %d (%.2f)", s.reads, s.hits, r)
}

// Evaluate returns the minimum distance to a cached 2d sdf.
func (s *CacheSDF2) Evaluate(p v2.Vec) float64 {
	s.reads++
	if d, ok := s.cache[p]; ok {
		s.hits++
		return d
	}
	d := s.sdf.Evaluate(p)
	s.cache[p] = d
	return d
}

// BoundingBox returns the bounding box of a cached 2d sdf.
func (s *CacheSDF2) BoundingBox() Box2 {
	return s.sdf.BoundingBox()
}

//-----------------------------------------------------------------------------
