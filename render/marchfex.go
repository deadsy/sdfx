//-----------------------------------------------------------------------------
/*

Marching Cubes Octree

Convert an SDF3 to finite elements.
Uses octree space subdivision.

*/
//-----------------------------------------------------------------------------

package render

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

// Process a cube. Generate finite elements, or more cubes.
func (dc *dcache3) processCubeHex8(c *cube, output chan<- []*Triangle3) {
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

// marchingCubesFEOctree generates a triangle mesh for an SDF3 using octree subdivision.
func marchingCubesFEOctree(s sdf.SDF3, resolution float64, output chan<- []*Triangle3) {
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
	dc.processCubeHex8(&cube{v3i.Vec{0, 0, 0}, levels - 1}, output)
}

//-----------------------------------------------------------------------------

// MarchingCubesFEOctree renders using marching cubes with octree space sampling.
type MarchingCubesFEOctree struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewMarchingCubesFEOctree returns a Render3 object.
func NewMarchingCubesFEOctree(meshCells int) *MarchingCubesFEOctree {
	return &MarchingCubesFEOctree{
		meshCells: meshCells,
	}
}

// Info returns a string describing the rendered volume.
func (r *MarchingCubesFEOctree) Info(s sdf.SDF3) string {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	cells := conv.V3ToV3i(bbSize.MulScalar(1 / resolution))
	return fmt.Sprintf("%dx%dx%d, resolution %.2f", cells.X, cells.Y, cells.Z, resolution)
}

// Render produces a 3d triangle mesh over the bounding volume of an sdf3.
func (r *MarchingCubesFEOctree) Render(s sdf.SDF3, output chan<- []*Triangle3) {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	marchingCubesFEOctree(s, resolution, output)
}

//-----------------------------------------------------------------------------
