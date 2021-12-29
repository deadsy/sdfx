//-----------------------------------------------------------------------------
/*

Multithreaded 3D renderer

*/
//-----------------------------------------------------------------------------

package render

import (
	"github.com/barkimedes/go-deepcopy"
	"github.com/deadsy/sdfx/sdf"
	"math"
	"sync"
)

// MTRenderer3 converts a SDF3 to a triangle mesh, parallelizing the rendering process by splitting the space.
// It supports MarchingCubesUniform and dc.DualContouringV2 (not MarchingCubesOctree for some reason),
// and should support other Render3 voxel-based implementation.
type MTRenderer3 struct {
	// impl the base implementation to copy and use for rendering partial surfaces.
	impl Render3
	// NumSplits is the number of times to split each dimension (0 splits means 1 partition).
	NumSplits sdf.V3i
	// OverlappingCells provides overlapping boundaries so that N cells are shared between contiguous renderers.
	// Set to 0 for Marching Cubes renderer and to 1 for Dual Contouring (places vertices instead of faces per voxel)
	OverlappingCells float64
}

var _ Render3 = &MTRenderer3{}

// NewMtRenderer3 see MTRenderer3
func NewMtRenderer3(impl Render3, overlappingCells float64) *MTRenderer3 {
	return &MTRenderer3{impl: impl, OverlappingCells: overlappingCells, NumSplits: sdf.V3i{0, 0, 0}}
}

// AutoSplitsMinimum auto-partitions the space to have at least `routines` partitions.
//
// Example: `MTRenderer3.AutoSplitsMinimum(runtime.NumCPU())`
func (m *MTRenderer3) AutoSplitsMinimum(minPartitions int) {
	// Try to share splits in all dimensions, forcing last dimension to get more splits if needed (may overshoot `minPartitions`)
	splits := float64(minPartitions) * math.Pow(2, -3 /* dimensions */)
	m.NumSplits = sdf.V3i{int(splits), int(splits), int(splits)}
	// If more splits are needed, give the first dimension more splits (generate at least minPartitions splits)
	// NOTE: partitions = splits + 1
	m.NumSplits[0] += (minPartitions - (m.NumSplits[0]+1)*(m.NumSplits[1]+1)*(m.NumSplits[2]+1)) /
		((m.NumSplits[1] + 1) * (m.NumSplits[2] + 1))
}

func (m *MTRenderer3) Cells(sdf3 sdf.SDF3, meshCells int) (float64, sdf.V3i) {
	return m.impl.Cells(sdf3, meshCells)
}

func (m *MTRenderer3) Render(sdf3 sdf.SDF3, meshCells int, output chan<- *Triangle3) {
	// Get cells to render on each dimension (change Render3 Info)
	originalResolution, originalCells := m.Cells(sdf3, meshCells)
	fullBb := sdf3.BoundingBox()
	fullBbSize := fullBb.Size()
	cellSize := fullBbSize.Div(originalCells.ToV3())
	numPartitions := m.NumSplits.AddScalar(1)

	// The priority is to keep the cellSize the same on all partitions (in case overlapping cells are needed)
	// Then, try to have all partitions be of the same size to avoid waiting for one goroutine to finish
	cellsPerPartitionBase := originalCells.ToV3().Div(numPartitions.ToV3())
	partitionSizeBase := cellSize.Mul(cellsPerPartitionBase)
	partitionSize := cellSize.Mul(cellsPerPartitionBase.AddScalar(m.OverlappingCells))
	subMeshCells := int( /*math.Ceil*/ partitionSize.MaxComponent() / originalResolution)

	cellIndex := sdf.V3i{0, 0, 0}
	wg := &sync.WaitGroup{}
	for cellIndex[0] = 0; cellIndex[0] <= m.NumSplits[0]; cellIndex[0]++ {
		for cellIndex[1] = 0; cellIndex[1] <= m.NumSplits[1]; cellIndex[1]++ {
			for cellIndex[2] = 0; cellIndex[2] <= m.NumSplits[2]; cellIndex[2]++ {
				// Get bounded sub-surface
				// WARNING: Will be slightly outside the original bounding-box (for simplicity, if OverlappingCells > 0),
				// but keeps partition size and cells per partition the same so that OverlappingCells works
				from := fullBb.Min.Add(partitionSizeBase.Mul(cellIndex.ToV3())).AddScalar(-m.OverlappingCells / 2)
				to := from.Add(partitionSize)
				subSurface := &customBbSdf{impl: sdf3, bb: sdf.Box3{Min: from, Max: to}}
				// Spawn the rendering goroutine
				go func(impl Render3, subSurface sdf.SDF3) {
					// Debug
					//resolution2, cells2 := impl.Cells(subSurface, subMeshCells)
					//fmt.Printf("rendering part (%dx%dx%d, resolution %.2f) from %s to %s\n",
					//	cells2[0], cells2[1], cells2[2], resolution2, fmt.Sprint(subSurface.BoundingBox().Min),
					//	fmt.Sprint(subSurface.BoundingBox().Max))
					// Just call render on the sub-surface, generating triangles to the same output
					impl.Render(subSurface, subMeshCells, output)
					wg.Done()
				}(deepcopy.MustAnything(m.impl).(Render3), subSurface)
				wg.Add(1)
			}
		}
	}
	wg.Wait() // Wait for all spawned goroutines (space partition renderers) to finish
}

var _ sdf.SDF3 = &customBbSdf{}

type customBbSdf struct {
	impl sdf.SDF3
	bb   sdf.Box3
}

func (m *customBbSdf) Evaluate(p sdf.V3) float64 {
	return m.impl.Evaluate(p)
}

func (m *customBbSdf) BoundingBox() sdf.Box3 {
	return m.bb
}
