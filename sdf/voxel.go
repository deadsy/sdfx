//-----------------------------------------------------------------------------
/*

Voxel-based cache/smoothing to remove deep SDF2/SDF3 hierarchies and speed up evaluation

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

// VoxelSDF3 is the SDF that represents a pre-computed voxel-based SDF3.
//It can be used as a cache, or for smoothing.
//
// CACHE:
// It can be used to speed up all evaluations required by the surface mesher at the cost of scene setup time and accuracy.
//
// SMOOTHING (meshCells <<< renderer's meshCells):
// It performs trilinear mapping for inner values and may be used as a cache for any other SDF, losing some accuracy.
//
// WARNING: It may lose sharp features, even if meshCells is high.
type VoxelSDF3 struct {
	// voxelCorners are the values of this SDF in each voxel corner
	voxelCorners map[v3i.Vec]float64 // TODO: Octree + k-d tree to simplify/reduce memory consumption + speed-up access?
	// bb is the bounding box.
	bb Box3
	// Number of voxelCorners to consider
	numVoxels v3i.Vec
}

// NewVoxelSDF3 returns a VoxelSDF3.
// This populates the whole cache from the given SDF.
// The progress listener may be nil.
func NewVoxelSDF3(s SDF3, meshCells int, progress chan float64) SDF3 {
	bb := s.BoundingBox() // TODO: Use default code to avoid duplication
	bbSize := bb.Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := conv.V3ToV3i(bbSize.DivScalar(resolution))

	voxelCorners := map[v3i.Vec]float64{}
	voxelCornerIndex := v3i.Vec{}
	for voxelCornerIndex.X = 0; voxelCornerIndex.X <= cells.X; voxelCornerIndex.X++ {
		for voxelCornerIndex.Y = 0; voxelCornerIndex.Y <= cells.Y; voxelCornerIndex.Y++ {
			for voxelCornerIndex.Z = 0; voxelCornerIndex.Z <= cells.Z; voxelCornerIndex.Z++ {
				voxelCorner := bb.Min.Add(bbSize.Mul(conv.V3iToV3(voxelCornerIndex)).Div(conv.V3iToV3(cells)))
				voxelCorners[voxelCornerIndex] = s.Evaluate(voxelCorner)
			}
		}
		if progress != nil {
			progress <- float64(voxelCornerIndex.X) / float64(cells.X)
		}
	}

	return &VoxelSDF3{
		voxelCorners: voxelCorners,
		bb:           bb,
		numVoxels:    cells,
	}
}

// Evaluate returns the minimum distance to a VoxelSDF3.
func (m *VoxelSDF3) Evaluate(p v3.Vec) float64 {
	// Find the voxel's {0,0,0} corner quickly and compute p's displacement
	voxelSize := m.bb.Size().Div(conv.V3iToV3(m.numVoxels))
	voxelStartIndex := conv.V3ToV3i(p.Sub(m.bb.Min).Div(voxelSize))
	voxelStart := m.bb.Min.Add(voxelSize.Mul(conv.V3iToV3(voxelStartIndex)))
	d := p.Sub(voxelStart).Div(voxelSize) // [0, 1) for each dimension
	// Get the values at the voxel's corners
	c000 := m.voxelCorners[voxelStartIndex]
	c001 := m.voxelCorners[voxelStartIndex.Add(v3i.Vec{0, 0, 1})]
	c010 := m.voxelCorners[voxelStartIndex.Add(v3i.Vec{0, 1, 0})]
	c011 := m.voxelCorners[voxelStartIndex.Add(v3i.Vec{0, 1, 1})]
	c100 := m.voxelCorners[voxelStartIndex.Add(v3i.Vec{1, 0, 0})]
	c101 := m.voxelCorners[voxelStartIndex.Add(v3i.Vec{1, 0, 1})]
	c110 := m.voxelCorners[voxelStartIndex.Add(v3i.Vec{1, 1, 0})]
	c111 := m.voxelCorners[voxelStartIndex.Add(v3i.Vec{1, 1, 1})]
	// Perform trilinear interpolation over the voxel's corners
	// - 4 linear interpolations
	c00 := c000*(1-d.X) + c100*d.X
	c01 := c001*(1-d.X) + c101*d.X
	c10 := c010*(1-d.X) + c110*d.X
	c11 := c011*(1-d.X) + c111*d.X
	// - 2 bilinear interpolations
	c0 := c00*(1-d.Y) + c10*d.Y
	c1 := c01*(1-d.Y) + c11*d.Y
	// - 1 trilinear interpolation
	c := c0*(1-d.Z) + c1*d.Z
	return c
}

// BoundingBox returns the bounding box for a VoxelSDF3.
func (m *VoxelSDF3) BoundingBox() Box3 {
	return m.bb
}

//-----------------------------------------------------------------------------
