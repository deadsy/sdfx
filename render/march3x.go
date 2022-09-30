//-----------------------------------------------------------------------------
/*

Marching Cubes Octree

Convert an SDF3 to a triangle mesh.
Uses octree space subdivision.

*/
//-----------------------------------------------------------------------------

package render

import (
	"fmt"
	"math"
	"sync"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

type cube struct {
	v v3i.Vec // origin of cube as integers
	n uint    // level of cube, size = 1 << n
}

//-----------------------------------------------------------------------------
// Evaluate the SDF3 via a distance cache to avoid repeated evaluations.
// Experimentally about 2/3 of lookups get a hit, and the overall speedup
// is about 2x a non-cached evaluation.

type dcache3 struct {
	origin     v3.Vec              // origin of the overall bounding cube
	resolution float64             // size of smallest octree cube
	hdiag      []float64           // lookup table of cube half diagonals
	s          sdf.SDF3            // the SDF3 to be rendered
	cache      map[v3i.Vec]float64 // cache of distances
	lock       sync.RWMutex        // lock the the cache during reads/writes
}

func newDcache3(s sdf.SDF3, origin v3.Vec, resolution float64, n uint) *dcache3 {
	// TODO heuristic for initial cache size. Maybe k * (1 << n)^3
	// Avoiding any resizing of the map seems to be worth 2-5% of speedup.
	dc := dcache3{
		origin:     origin,
		resolution: resolution,
		hdiag:      make([]float64, n),
		s:          s,
		cache:      make(map[v3i.Vec]float64),
	}
	// build a lut for cube half diagonal lengths
	for i := range dc.hdiag {
		si := 1 << uint(i)
		s := float64(si) * dc.resolution
		dc.hdiag[i] = 0.5 * math.Sqrt(3.0*s*s)
	}
	return &dc
}

// read from the cache
func (dc *dcache3) read(vi v3i.Vec) (float64, bool) {
	dc.lock.RLock()
	dist, found := dc.cache[vi]
	dc.lock.RUnlock()
	return dist, found
}

// write to the cache
func (dc *dcache3) write(vi v3i.Vec, dist float64) {
	dc.lock.Lock()
	dc.cache[vi] = dist
	dc.lock.Unlock()
}

func (dc *dcache3) evaluate(vi v3i.Vec) (v3.Vec, float64) {
	v := dc.origin.Add(conv.V3iToV3(vi).MulScalar(dc.resolution))
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

// isEmpty returns true if the cube contains no SDF surface
func (dc *dcache3) isEmpty(c *cube) bool {
	// evaluate the SDF3 at the center of the cube
	s := 1 << (c.n - 1) // half side
	_, d := dc.evaluate(c.v.AddScalar(s))
	// compare to the center/corner distance
	return math.Abs(d) >= dc.hdiag[c.n]
}

// Process a cube. Generate triangles, or more cubes.
func (dc *dcache3) processCube(c *cube, output chan<- []*Triangle3) {
	if !dc.isEmpty(c) {
		if c.n == 1 {
			// this cube is at the required resolution
			c0, d0 := dc.evaluate(c.v.Add(v3i.Vec{0, 0, 0}))
			c1, d1 := dc.evaluate(c.v.Add(v3i.Vec{2, 0, 0}))
			c2, d2 := dc.evaluate(c.v.Add(v3i.Vec{2, 2, 0}))
			c3, d3 := dc.evaluate(c.v.Add(v3i.Vec{0, 2, 0}))
			c4, d4 := dc.evaluate(c.v.Add(v3i.Vec{0, 0, 2}))
			c5, d5 := dc.evaluate(c.v.Add(v3i.Vec{2, 0, 2}))
			c6, d6 := dc.evaluate(c.v.Add(v3i.Vec{2, 2, 2}))
			c7, d7 := dc.evaluate(c.v.Add(v3i.Vec{0, 2, 2}))
			corners := [8]v3.Vec{c0, c1, c2, c3, c4, c5, c6, c7}
			values := [8]float64{d0, d1, d2, d3, d4, d5, d6, d7}
			// output the triangle(s) for this cube
			output <- mcToTriangles(corners, values, 0)
		} else {
			// process the sub cubes
			n := c.n - 1
			s := 1 << n
			// TODO - turn these into throttled go-routines
			dc.processCube(&cube{c.v.Add(v3i.Vec{0, 0, 0}), n}, output)
			dc.processCube(&cube{c.v.Add(v3i.Vec{s, 0, 0}), n}, output)
			dc.processCube(&cube{c.v.Add(v3i.Vec{s, s, 0}), n}, output)
			dc.processCube(&cube{c.v.Add(v3i.Vec{0, s, 0}), n}, output)
			dc.processCube(&cube{c.v.Add(v3i.Vec{0, 0, s}), n}, output)
			dc.processCube(&cube{c.v.Add(v3i.Vec{s, 0, s}), n}, output)
			dc.processCube(&cube{c.v.Add(v3i.Vec{s, s, s}), n}, output)
			dc.processCube(&cube{c.v.Add(v3i.Vec{0, s, s}), n}, output)
		}
	}
}

//-----------------------------------------------------------------------------

// marchingCubesOctree generates a triangle mesh for an SDF3 using octree subdivision.
func marchingCubesOctree(s sdf.SDF3, resolution float64, output chan<- []*Triangle3) {
	// Scale the bounding box about the center to make sure the boundaries
	// aren't on the object surface.
	bb := s.BoundingBox()
	bb = bb.ScaleAboutCenter(1.01)
	longAxis := bb.Size().MaxComponent()
	// We want to test the smallest cube (side == resolution) for emptiness
	// so the level = 0 cube is at half resolution.
	resolution = 0.5 * resolution
	// how many cube levels for the octree?
	levels := uint(math.Ceil(math.Log2(longAxis/resolution))) + 1
	// create the distance cache
	dc := newDcache3(s, bb.Min, resolution, levels)
	// process the octree, start at the top level
	dc.processCube(&cube{v3i.Vec{0, 0, 0}, levels - 1}, output)
}

//-----------------------------------------------------------------------------

// MarchingCubesOctree renders using marching cubes with octree space sampling.
type MarchingCubesOctree struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewMarchingCubesOctree returns a Render3 object.
func NewMarchingCubesOctree(meshCells int) *MarchingCubesOctree {
	return &MarchingCubesOctree{
		meshCells: meshCells,
	}
}

// Info returns a string describing the rendered volume.
func (r *MarchingCubesOctree) Info(s sdf.SDF3) string {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	cells := conv.V3ToV3i(bbSize.MulScalar(1 / resolution))
	return fmt.Sprintf("%dx%dx%d, resolution %.2f", cells.X, cells.Y, cells.Z, resolution)
}

// Render produces a 3d triangle mesh over the bounding volume of an sdf3.
func (r *MarchingCubesOctree) Render(s sdf.SDF3, output chan<- []*Triangle3) {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	marchingCubesOctree(s, resolution, output)
}

//-----------------------------------------------------------------------------
