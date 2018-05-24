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

func (dc *dcache) evaluate(vi V3i) (V3, float64) {
	v := dc.origin.Add(vi.ToV3().MulScalar(dc.resolution))
	// do we have it in the cache?
	dist, found := dc.read(vi)
	if found {
		return v, dist
	}
	// evaluate the SDF3
	dist = dc.s.Evaluate(v)
	// write it to the cache
	dc.write(vi, dist)
	return v, dist
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
	_, d := dc.evaluate(c.v.Add(V3i{s, s, s}))
	// compare to the center/corner distance
	x := float64(s) * dc.resolution
	// TODO - optimise the sqrt? Use a LUT?
	return Abs(d) >= math.Sqrt(3.0*x*x)
}

// Process a cube. Generate triangles, or more cubes.
func (c *cube) process(dc *dcache, wq chan<- *cube, output chan<- *Triangle3) {
	if c.s == 1 {
		// this cube is at the required resolution
		c0, d0 := dc.evaluate(c.v.Add(V3i{0, 0, 0}))
		c1, d1 := dc.evaluate(c.v.Add(V3i{1, 0, 0}))
		c2, d2 := dc.evaluate(c.v.Add(V3i{1, 1, 0}))
		c3, d3 := dc.evaluate(c.v.Add(V3i{0, 1, 0}))
		c4, d4 := dc.evaluate(c.v.Add(V3i{0, 0, 1}))
		c5, d5 := dc.evaluate(c.v.Add(V3i{1, 0, 1}))
		c6, d6 := dc.evaluate(c.v.Add(V3i{1, 1, 1}))
		c7, d7 := dc.evaluate(c.v.Add(V3i{0, 1, 1}))
		corners := [8]V3{c0, c1, c2, c3, c4, c5, c6, c7}
		values := [8]float64{d0, d1, d2, d3, d4, d5, d6, d7}
		// output the triangle(s) for this cube
		for _, t := range mc_ToTriangles(corners, values, 0) {
			output <- t
		}
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
