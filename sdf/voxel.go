//-----------------------------------------------------------------------------
/*

Voxel-based cache to remove deep SDF3 hierarchies at setup and speed up evaluation

*/
//-----------------------------------------------------------------------------

package sdf

// VoxelSdf is the SDF that represents a pre-computed voxel-based SDF3.
//It can be used as a cache, or for smoothing.
//
// CACHE:
// It can be used to speed up all evaluations required by the surface mesher at the cost of scene setup time and accuracy.
//
// SMOOTHING (meshCells <<< renderer's meshCells):
// It performs trilinear mapping for inner values and may be used as a cache for any other SDF, losing some accuracy.
//
// WARNING: It may lose sharp features, even if meshCells is high.
type VoxelSdf struct {
	// voxelCorners are the values of this SDF in each voxel corner
	voxelCorners map[V3i]float64 // TODO: Octree + k-d tree to simplify/reduce memory consumption + speed-up access?
	// bb is the bounding box.
	bb Box3
	// Number of voxelCorners to consider
	numVoxels V3i
}

// NewVoxelSDF see VoxelSdf. This populates the whole cache from the given SDF. The progress listener may be nil.
func NewVoxelSDF(s SDF3, meshCells int, progress chan float64) SDF3 {
	bb := s.BoundingBox() // TODO: Use default code to avoid duplication
	bbSize := bb.Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := bbSize.DivScalar(resolution).ToV3i()

	voxelCorners := map[V3i]float64{}
	voxelCornerIndex := V3i{}
	for voxelCornerIndex[0] = 0; voxelCornerIndex[0] <= cells[0]; voxelCornerIndex[0]++ {
		for voxelCornerIndex[1] = 0; voxelCornerIndex[1] <= cells[1]; voxelCornerIndex[1]++ {
			for voxelCornerIndex[2] = 0; voxelCornerIndex[2] <= cells[2]; voxelCornerIndex[2]++ {
				voxelCorner := bb.Min.Add(bbSize.Mul(voxelCornerIndex.ToV3()).Div(cells.ToV3()))
				voxelCorners[voxelCornerIndex] = s.Evaluate(voxelCorner)
			}
		}
		if progress != nil {
			progress <- float64(voxelCornerIndex[0]) / float64(cells[0])
		}
	}

	return &VoxelSdf{
		voxelCorners: voxelCorners,
		bb:           bb,
		numVoxels:    cells,
	}
}

func (m *VoxelSdf) Evaluate(p V3) float64 {
	// Find the voxel's {0,0,0} corner quickly and compute p's displacement
	voxelSize := m.bb.Size().Div(m.numVoxels.ToV3())
	voxelStartIndex := p.Sub(m.bb.Min).Div(voxelSize).ToV3i()
	voxelStart := m.bb.Min.Add(voxelSize.Mul(voxelStartIndex.ToV3()))
	d := p.Sub(voxelStart).Div(voxelSize) // [0, 1) for each dimension
	// Get the values at the voxel's corners
	c000 := m.voxelCorners[voxelStartIndex]
	c001 := m.voxelCorners[voxelStartIndex.Add(V3i{0, 0, 1})]
	c010 := m.voxelCorners[voxelStartIndex.Add(V3i{0, 1, 0})]
	c011 := m.voxelCorners[voxelStartIndex.Add(V3i{0, 1, 1})]
	c100 := m.voxelCorners[voxelStartIndex.Add(V3i{1, 0, 0})]
	c101 := m.voxelCorners[voxelStartIndex.Add(V3i{1, 0, 1})]
	c110 := m.voxelCorners[voxelStartIndex.Add(V3i{1, 1, 0})]
	c111 := m.voxelCorners[voxelStartIndex.Add(V3i{1, 1, 1})]
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

func (m *VoxelSdf) BoundingBox() Box3 {
	return m.bb
}
