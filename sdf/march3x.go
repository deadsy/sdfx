//-----------------------------------------------------------------------------
/*

Marching Cubes

Convert an SDF3 to a triangle mesh.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
	"sync"
)

//-----------------------------------------------------------------------------
// Evaluate the SDF3 via a distance cache to avoid repeated evaluations.

type dcache struct {
	origin     V3              // origin of bounding cube
	resolution float64         // size of smallest octree cube
	s          SDF3            // the SDF3 to be rendered
	cache      map[V3i]float64 // cache of distances
	lock       sync.RWMutex    // lock the the cache during reads/writes
}

func new_dcache(s SDF3) *dcache {
	return &dcache{
		cache: make(map[V3i]float64),
		s:     s,
		lock:  sync.RWMutex{},
	}
}

// read from the cache
func (dc *dcache) read(vi V3i) (float64, bool) {
	dc.lock.RLock()
	dist, found := dc.cache[vi]
	dc.lock.RUnlock()
	return dist, found
}

// write to the cache
func (dc *dcache) write(vi V3i, dist float64) {
	dc.lock.Lock()
	dc.cache[vi] = dist
	dc.lock.Unlock()
}

func (dc *dcache) evaluate(vi V3i) float64 {
	// do we have it in the cache?
	dist, found := dc.read(vi)
	if found {
		return dist
	}
	// evaluate the SDF3
	v := dc.origin.Add(vi.ToV3().MulScalar(dc.resolution))
	dist = dc.s.Evaluate(v)
	// write it to the cache
	dc.write(vi, dist)
	return dist
}

//-----------------------------------------------------------------------------

type cube struct {
	v V3i // origin of cube as integers
	s int // size of cube side as an integer
}

//-----------------------------------------------------------------------------

// MarchingCubes generates a triangle mesh for an SDF3
func MarchingCubesX(sdf SDF3, resolution float64, output chan<- *Triangle3) {

}

//-----------------------------------------------------------------------------

// is the cube empty?
func (c *cube) empty(dc *dcache) bool {
	// evaluate the SDF3 at the center of the cube
	s := c.s >> 1
	dist := dc.evaluate(c.v.Add(V3i{s, s, s}))
	// compare to the center/corner distance
	x := float64(s) * dc.resolution
	return Abs(dist) >= math.Sqrt(3.0*x*x)
}

// Process a cube. Generate triangles, or more cubes.
func (c *cube) process(dc *dcache, wq chan<- *cube, output chan<- *Triangle3) {
	if c.s == 1 {
		// this cube is at the required resolution
		values := [8]float64{
			dc.evaluate(c.v.Add(V3i{0, 0, 0})),
			dc.evaluate(c.v.Add(V3i{1, 0, 0})),
			dc.evaluate(c.v.Add(V3i{1, 1, 0})),
			dc.evaluate(c.v.Add(V3i{0, 1, 0})),
			dc.evaluate(c.v.Add(V3i{0, 0, 1})),
			dc.evaluate(c.v.Add(V3i{1, 0, 1})),
			dc.evaluate(c.v.Add(V3i{1, 1, 1})),
			dc.evaluate(c.v.Add(V3i{0, 1, 1})),
		}
		// TODO - run marching cubes

		_ = values

		output <- &Triangle3{}
	} else {
		if !c.empty(dc) {
			// the cube is not empty- generate the sub cubes
			s := c.s >> 1
			wq <- &cube{c.v.Add(V3i{0, 0, 0}), s}
			wq <- &cube{c.v.Add(V3i{s, 0, 0}), s}
			wq <- &cube{c.v.Add(V3i{s, s, 0}), s}
			wq <- &cube{c.v.Add(V3i{0, s, 0}), s}
			wq <- &cube{c.v.Add(V3i{0, 0, s}), s}
			wq <- &cube{c.v.Add(V3i{s, 0, s}), s}
			wq <- &cube{c.v.Add(V3i{s, s, s}), s}
			wq <- &cube{c.v.Add(V3i{0, s, s}), s}
		}
	}

}

//-----------------------------------------------------------------------------
