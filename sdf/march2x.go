//-----------------------------------------------------------------------------
/*

Marching Squares Quadtree

Convert an SDF2 boundary to a set of line segments.
Uses quadtree space subdivision.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
	"sync"
)

//-----------------------------------------------------------------------------

type square struct {
	v V2i  // origin of square as integers
	n uint // level of square, size = 1 << n
}

//-----------------------------------------------------------------------------
// Evaluate the SDF2 via a distance cache to avoid repeated evaluations.

type dcache2 struct {
	origin     V2              // origin of the overall bounding square
	resolution float64         // size of smallest quadtree square
	hdiag      []float64       // lookup table of square half diagonals
	s          SDF2            // the SDF2 to be rendered
	cache      map[V2i]float64 // cache of distances
	lock       sync.RWMutex    // lock the the cache during reads/writes
}

func newDcache2(s SDF2, origin V2, resolution float64, n uint) *dcache2 {
	dc := dcache2{
		origin:     origin,
		resolution: resolution,
		hdiag:      make([]float64, n),
		s:          s,
		cache:      make(map[V2i]float64),
	}
	// build a lut for cube half diagonal lengths
	for i := range dc.hdiag {
		si := 1 << uint(i)
		s := float64(si) * dc.resolution
		dc.hdiag[i] = 0.5 * math.Sqrt(2.0*s*s)
	}
	return &dc
}

// read from the cache
func (dc *dcache2) read(vi V2i) (float64, bool) {
	dc.lock.RLock()
	dist, found := dc.cache[vi]
	dc.lock.RUnlock()
	return dist, found
}

// write to the cache
func (dc *dcache2) write(vi V2i, dist float64) {
	dc.lock.Lock()
	dc.cache[vi] = dist
	dc.lock.Unlock()
}

func (dc *dcache2) evaluate(vi V2i) (V2, float64) {
	v := dc.origin.Add(vi.ToV2().MulScalar(dc.resolution))
	// do we have it in the cache?
	dist, found := dc.read(vi)
	if found {
		return v, dist
	}
	// evaluate the SDF2
	dist = dc.s.Evaluate(v)
	// write it to the cache
	dc.write(vi, dist)
	return v, dist
}

// isEmpty returns true if the square contains no SDF surface
func (dc *dcache2) isEmpty(c *square) bool {
	// evaluate the SDF2 at the center of the square
	s := 1 << (c.n - 1) // half side
	_, d := dc.evaluate(c.v.AddScalar(s))
	// compare to the center/corner distance
	return Abs(d) >= dc.hdiag[c.n]
}

// Process a square. Generate line segments, or more squares.
func (dc *dcache2) processSquare(c *square, output chan<- *Line) {
	if !dc.isEmpty(c) {
		if c.n == 1 {
			// this square is at the required resolution
			c0, d0 := dc.evaluate(c.v.Add(V2i{0, 0}))
			c1, d1 := dc.evaluate(c.v.Add(V2i{2, 0}))
			c2, d2 := dc.evaluate(c.v.Add(V2i{2, 2}))
			c3, d3 := dc.evaluate(c.v.Add(V2i{0, 2}))
			corners := [4]V2{c0, c1, c2, c3}
			values := [4]float64{d0, d1, d2, d3}
			// output the line(s) for this square
			for _, l := range msToLines(corners, values, 0) {
				output <- l
			}
		} else {
			// process the sub squares
			n := c.n - 1
			s := 1 << n
			// TODO - turn these into throttled go-routines
			dc.processSquare(&square{c.v.Add(V2i{0, 0}), n}, output)
			dc.processSquare(&square{c.v.Add(V2i{s, 0}), n}, output)
			dc.processSquare(&square{c.v.Add(V2i{s, s}), n}, output)
			dc.processSquare(&square{c.v.Add(V2i{0, s}), n}, output)
		}
	}
}

//-----------------------------------------------------------------------------

// marchingSquaresQuadtree generates line segments for an SDF2 using quadtree subdivision.
func marchingSquaresQuadtree(s SDF2, resolution float64, output chan<- *Line) {
	// Scale the bounding box about the center to make sure the boundaries
	// aren't on the object surface.
	bb := s.BoundingBox()
	bb = bb.ScaleAboutCenter(1.01)
	longAxis := bb.Size().MaxComponent()
	// We want to test the smallest squares (side == resolution) for emptiness
	// so the level = 0 cube is at half resolution.
	resolution = 0.5 * resolution
	// how many cube levels for the quadtree?
	levels := uint(math.Ceil(math.Log2(longAxis/resolution))) + 1
	// create the distance cache
	dc := newDcache2(s, bb.Min, resolution, levels)
	// process the quadtree, start at the top level
	dc.processSquare(&square{V2i{0, 0}, levels - 1}, output)
}

//-----------------------------------------------------------------------------
