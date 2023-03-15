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
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

func distanceCache(s sdf.SDF3, resolution float64) (*dcache3, uint) {
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
	return dc, levels
}

// marchingCubesTet4Octree generates finite elements for an SDF3 using octree subdivision.
func marchingCubesTet4Octree(s sdf.SDF3, resolution float64, output chan<- []*Tet4) {
	dc, levels := distanceCache(s, resolution)
	// process the octree, start at the top level
	dc.processCubeTet4(&cube{v3i.Vec{0, 0, 0}, levels - 1}, output)
}

// marchingCubesHex8Octree generates finite elements for an SDF3 using octree subdivision.
func marchingCubesHex8Octree(s sdf.SDF3, resolution float64, output chan<- []*Hex8) {
	dc, levels := distanceCache(s, resolution)
	// process the octree, start at the top level
	dc.processCubeHex8(&cube{v3i.Vec{0, 0, 0}, levels - 1}, output)
}

// marchingCubesHex20Octree generates finite elements for an SDF3 using octree subdivision.
func marchingCubesHex20Octree(s sdf.SDF3, resolution float64, output chan<- []*Hex20) {
	dc, levels := distanceCache(s, resolution)
	// process the octree, start at the top level
	dc.processCubeHex20(&cube{v3i.Vec{0, 0, 0}, levels - 1}, output)
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

// LayerCounts computes number of layes in X, Y, Z directions.
// We are specifically interested in Z direction.
func (r *MarchingCubesFEOctree) LayerCounts(s sdf.SDF3) (int, int, int) {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	cells := conv.V3ToV3i(bbSize.MulScalar(1 / resolution))
	return cells.X, cells.Y, cells.Z
}

// RenderTet4 produces finite elements over the bounding volume of an sdf3.
func (r *MarchingCubesFEOctree) RenderTet4(s sdf.SDF3, output chan<- []*Tet4) {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	marchingCubesTet4Octree(s, resolution, output)
}

// RenderHex8 produces finite elements over the bounding volume of an sdf3.
func (r *MarchingCubesFEOctree) RenderHex8(s sdf.SDF3, output chan<- []*Hex8) {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	marchingCubesHex8Octree(s, resolution, output)
}

// RenderHex20 produces finite elements over the bounding volume of an sdf3.
func (r *MarchingCubesFEOctree) RenderHex20(s sdf.SDF3, output chan<- []*Hex20) {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	marchingCubesHex20Octree(s, resolution, output)
}

//-----------------------------------------------------------------------------
