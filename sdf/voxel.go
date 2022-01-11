//-----------------------------------------------------------------------------
/*

Voxel-based cache/smoothing to remove deep SDF2/SDF3 hierarchies and speed up evaluation

*/
//-----------------------------------------------------------------------------

package sdf

import "sync"

//-----------------------------------------------------------------------------

// VoxelSDF3 is the SDF that represents a pre-computed voxel-based SDF3.
// It can be used as a cache and/or for smoothing.
//
// CACHE:
// It can be used to speed up all evaluations required by the surface mesher at the cost of scene setup time and accuracy.
//
// SMOOTHING (meshCells < renderer's meshCells):
// It performs trilinear interpolation for inner values and may be used as a cache for any other SDF, losing some accuracy.
//
// WARNING: It may lose sharp features, even if meshCells is high.
type VoxelSDF3 struct {
	// voxelCorners are the values of this SDF in each voxel corner (populated lazily by default)
	voxelCorners map[V3i]float64
	// s is the SDF
	s SDF3
	// Number of voxelCorners to consider
	numVoxels V3i
	// mu is the mutex for allowing concurrent access (set to nil if not necessary)
	mu *sync.RWMutex
}

// NewVoxelSDF3 returns a VoxelSDF3.
// This populates the whole cache from the given SDF.
// synchronize is required for concurrent access (multithread renderers).
func NewVoxelSDF3(s SDF3, meshCells int, synchronize bool) SDF3 {
	bb := s.BoundingBox() // TODO: Use default code to avoid duplication
	bbSize := bb.Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := bbSize.DivScalar(resolution).ToV3i()
	var mu *sync.RWMutex
	if synchronize {
		mu = &sync.RWMutex{}
	}
	return &VoxelSDF3{
		voxelCorners: map[V3i]float64{},
		s:            s,
		numVoxels:    cells,
		mu:           mu,
	}
}

// getOrCompute retrieves the distance for a specific voxel index, computing it if not cached.
// NOTE: This will also work for values outside the bounding box (within `int` limits).
func (m *VoxelSDF3) getOrCompute(voxelStartIndex V3i) float64 {
	if m.mu != nil {
		m.mu.RLock()
	}
	cached, ok := m.voxelCorners[voxelStartIndex]
	if m.mu != nil {
		m.mu.RUnlock()
	}
	// This may cause double writes, but those are not a problem (same value written) and avoids locking for writes if not needed
	if !ok {
		bb := m.BoundingBox()
		bbSize := bb.Size()
		voxelCorner := bb.Min.Add(bbSize.Mul(voxelStartIndex.ToV3()).Div(m.numVoxels.ToV3()))
		cached = m.s.Evaluate(voxelCorner)
		// Only acquire write access if absolutely necessary, as reads can be concurrent
		if m.mu != nil {
			m.mu.Lock()
		}
		m.voxelCorners[voxelStartIndex] = cached
		if m.mu != nil {
			m.mu.Unlock()
		}
	}
	return cached
}

// Populate forces the population of the full VoxelSDF (inside the bounding box), optionally publishing the progress.
// Populate may increase performance by avoiding all locking (singlethread), reducing to a minimum the number of
// synchronizations needed while calling Evaluate
func (m *VoxelSDF3) Populate(progress chan float64) map[V3i]float64 {
	cells := m.numVoxels
	voxelCorners := map[V3i]float64{}
	voxelCornerIndex := V3i{}
	prevMu := m.mu
	m.mu = nil
	for voxelCornerIndex[0] = 0; voxelCornerIndex[0] <= cells[0]; voxelCornerIndex[0]++ {
		for voxelCornerIndex[1] = 0; voxelCornerIndex[1] <= cells[1]; voxelCornerIndex[1]++ {
			for voxelCornerIndex[2] = 0; voxelCornerIndex[2] <= cells[2]; voxelCornerIndex[2]++ {
				m.getOrCompute(voxelCornerIndex)
			}
		}
		if progress != nil {
			progress <- float64(voxelCornerIndex[0]) / float64(cells[0])
		}
	}
	m.mu = prevMu
	return voxelCorners
}

// Evaluate returns the minimum distance to a VoxelSDF3.
func (m *VoxelSDF3) Evaluate(p V3) float64 {
	// Find the voxel's {0,0,0} corner quickly and compute p's displacement
	bb := m.BoundingBox()
	voxelSize := bb.Size().Div(m.numVoxels.ToV3())
	voxelStartIndex := p.Sub(bb.Min).Div(voxelSize).ToV3i()
	voxelStart := bb.Min.Add(voxelSize.Mul(voxelStartIndex.ToV3()))
	d := p.Sub(voxelStart).Div(voxelSize) // [0, 1) for each dimension
	// Get the values at the voxel's corners
	c000 := m.getOrCompute(voxelStartIndex)
	c001 := m.getOrCompute(voxelStartIndex.Add(V3i{0, 0, 1}))
	c010 := m.getOrCompute(voxelStartIndex.Add(V3i{0, 1, 0}))
	c011 := m.getOrCompute(voxelStartIndex.Add(V3i{0, 1, 1}))
	c100 := m.getOrCompute(voxelStartIndex.Add(V3i{1, 0, 0}))
	c101 := m.getOrCompute(voxelStartIndex.Add(V3i{1, 0, 1}))
	c110 := m.getOrCompute(voxelStartIndex.Add(V3i{1, 1, 0}))
	c111 := m.getOrCompute(voxelStartIndex.Add(V3i{1, 1, 1}))
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
	return m.s.BoundingBox()
}

//-----------------------------------------------------------------------------

// VoxelSDF2 is the SDF that represents a pre-computed voxel-based SDF2.
// It can be used as a cache and/or for smoothing.
//
// CACHE:
// It can be used to speed up all evaluations required by the surface mesher at the cost of scene setup time and accuracy.
//
// SMOOTHING (meshCells < renderer's meshCells):
// It performs bilinear interpolation for inner values and may be used as a cache for any other SDF, losing some accuracy.
//
// WARNING: It may lose sharp features, even if meshCells is high.
type VoxelSDF2 struct {
	// voxelCorners are the values of this SDF in each voxel corner (populated lazily by default)
	voxelCorners map[V2i]float64
	// s is the SDF
	s SDF2
	// Number of voxelCorners to consider
	numVoxels V2i
	// mu is the mutex for allowing concurrent access (set to nil if not necessary)
	mu *sync.RWMutex
}

// NewVoxelSDF2 returns a VoxelSDF2.
// This populates the whole cache from the given SDF.
// synchronize is required for concurrent access (multithread renderers).
func NewVoxelSDF2(s SDF2, meshCells int, synchronize bool) SDF2 {
	bb := s.BoundingBox() // TODO: Use default code to avoid duplication
	bbSize := bb.Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := bbSize.DivScalar(resolution).ToV2i()
	var mu *sync.RWMutex
	if synchronize {
		mu = &sync.RWMutex{}
	}
	return &VoxelSDF2{
		voxelCorners: map[V2i]float64{},
		s:            s,
		numVoxels:    cells,
		mu:           mu,
	}
}

// getOrCompute retrieves the distance for a specific voxel index, computing it if not cached.
// NOTE: This will also work for values outside the bounding box (within `int` limits).
func (m *VoxelSDF2) getOrCompute(voxelStartIndex V2i) float64 {
	if m.mu != nil {
		m.mu.RLock()
	}
	cached, ok := m.voxelCorners[voxelStartIndex]
	if m.mu != nil {
		m.mu.RUnlock()
	}
	// This may cause double writes, but those are not a problem (same value written) and avoids locking for writes if not needed
	if !ok {
		bb := m.BoundingBox()
		bbSize := bb.Size()
		voxelCorner := bb.Min.Add(bbSize.Mul(voxelStartIndex.ToV2()).Div(m.numVoxels.ToV2()))
		cached = m.s.Evaluate(voxelCorner)
		// Only acquire write access if absolutely necessary, as reads can be concurrent
		if m.mu != nil {
			m.mu.Lock()
		}
		m.voxelCorners[voxelStartIndex] = cached
		if m.mu != nil {
			m.mu.Unlock()
		}
	}
	return cached
}

// Populate forces the population of the full VoxelSDF (inside the bounding box), optionally publishing the progress
func (m *VoxelSDF2) Populate(progress chan float64) map[V2i]float64 {
	cells := m.numVoxels
	voxelCorners := map[V2i]float64{}
	voxelCornerIndex := V2i{}
	prevMu := m.mu
	m.mu = nil
	for voxelCornerIndex[0] = 0; voxelCornerIndex[0] <= cells[0]; voxelCornerIndex[0]++ {
		for voxelCornerIndex[1] = 0; voxelCornerIndex[1] <= cells[1]; voxelCornerIndex[1]++ {
			m.getOrCompute(voxelCornerIndex)
		}
		if progress != nil {
			progress <- float64(voxelCornerIndex[0]) / float64(cells[0])
		}
	}
	m.mu = prevMu
	return voxelCorners
}

// Evaluate returns the minimum distance to a VoxelSDF2.
func (m *VoxelSDF2) Evaluate(p V2) float64 {
	// Find the voxel's {0,0,0} corner quickly and compute p's displacement
	bb := m.BoundingBox()
	voxelSize := bb.Size().Div(m.numVoxels.ToV2())
	voxelStartIndex := p.Sub(bb.Min).Div(voxelSize).ToV2i()
	voxelStart := bb.Min.Add(voxelSize.Mul(voxelStartIndex.ToV2()))
	d := p.Sub(voxelStart).Div(voxelSize) // [0, 1) for each dimension
	// Get the values at the voxel's corners
	c00 := m.getOrCompute(voxelStartIndex)
	c01 := m.getOrCompute(voxelStartIndex.Add(V2i{0, 1}))
	c10 := m.getOrCompute(voxelStartIndex.Add(V2i{1, 0}))
	c11 := m.getOrCompute(voxelStartIndex.Add(V2i{1, 1}))
	// Perform bilinear interpolation over the voxel's corners
	// - 2 linear interpolations
	c0 := c00*(1-d.X) + c10*d.X
	c1 := c01*(1-d.X) + c11*d.X
	// - 1 bilinear interpolation
	c := c0*(1-d.Y) + c1*d.Y
	return c
}

// BoundingBox returns the bounding box for a VoxelSDF2.
func (m *VoxelSDF2) BoundingBox() Box2 {
	return m.s.BoundingBox()
}
