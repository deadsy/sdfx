//-----------------------------------------------------------------------------
/*

Multithreaded 3D renderer

*/
//-----------------------------------------------------------------------------

package render

import (
	"github.com/barkimedes/go-deepcopy"
	"github.com/deadsy/sdfx/sdf"
	"gonum.org/v1/gonum/spatial/kdtree"
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
	// MergeVerticesEpsilon is the minimum distance between vertices to merge them. This is needed because vertices in
	// the chunk seams can be slightly different (due to running independent algorithms and operation order might change),
	// causing software to think they are different vertices when they shouldn't be.
	MergeVerticesEpsilon float64
}

var _ Render3 = &MTRenderer3{}

// NewMtRenderer3 see MTRenderer3
func NewMtRenderer3(impl Render3, overlappingCells float64) *MTRenderer3 {
	return &MTRenderer3{impl: impl, OverlappingCells: overlappingCells, NumSplits: sdf.V3i{0, 0, 0}, MergeVerticesEpsilon: 1e-8}
}

// AutoSplitsMinimum auto-partitions the space to have at least `routines` partitions.
//
// It is recommended to use the number of logical CPUs of the system or more. More partitions may result in better
// performance (as long as synchronization overhead is not too much and meshCells is big enough),
// as chunks are shared among goroutine workers and there is less wait for the last chunk to finish.
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

	// Intercept output to merge vertices generated on seams
	wg := &sync.WaitGroup{}
	var wgRet *sync.WaitGroup
	if m.MergeVerticesEpsilon > 0 {
		output, wgRet = m.mergeGeneratedVertices(output, m.MergeVerticesEpsilon*cellSize.MaxComponent(),
			fullBb.Min.AddScalar(-m.OverlappingCells/2), partitionSizeBase, cellSize)
	}

	// Generate each partition in a goroutine
	cellIndex := sdf.V3i{0, 0, 0}
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
	if m.MergeVerticesEpsilon > 0 {
		close(output) // Actually the input for internal goroutine
		wgRet.Wait()  // Wait for post-processing to complete
	}
}

//-----------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------

func mtToKdPoint(v3 sdf.V3) kdtree.Point {
	return kdtree.Point{v3.X, v3.Y, v3.Z}
}

func mtFromKdPoint(v3 kdtree.Point) sdf.V3 {
	return sdf.V3{X: v3[0], Y: v3[1], Z: v3[2]}
}

func (m *MTRenderer3) mergeGeneratedVertices(output chan<- *Triangle3, mergeVerticesEpsilon float64, sdfStart sdf.V3, partitionSize sdf.V3, cellSize sdf.V3) (chan *Triangle3, *sync.WaitGroup) {
	input := make(chan *Triangle3)
	mergeVerticesEpsilonSq := mergeVerticesEpsilon * mergeVerticesEpsilon // Squared euclidean distance used in k-d tree
	mmod := func(a, b float64) float64 {
		r := math.Mod(a, b)
		if r > b/2 {
			r = b - r
		}
		return r
	}
	isVertCloseToSeam := func(v sdf.V3) bool { // Filter to speed-up postprocessing
		offset := v.Sub(sdfStart)
		return mmod(offset.X, partitionSize.X) < mergeVerticesEpsilon+cellSize.X ||
			mmod(offset.Y, partitionSize.Y) < mergeVerticesEpsilon+cellSize.Y ||
			mmod(offset.Z, partitionSize.Z) < mergeVerticesEpsilon+cellSize.Z
	}

	vertToPassOriginalTris := map[sdf.V3][]*Triangle3{}

	wg := &sync.WaitGroup{}
	wgRet := &sync.WaitGroup{}
	wg.Add(1)
	go func() { // Collect all initial vertices in k-d tree
		for tri := range input {
			if bypassVertTris(isVertCloseToSeam, nil, []*Triangle3{tri}) {
				output <- tri
			} else {
				for _, vert := range tri.V {
					// Do not modify the vertices with math operations, as floating point operations might make the map fail
					prevTris, _ := vertToPassOriginalTris[vert]
					prevTris = append(prevTris, tri)
					vertToPassOriginalTris[vert] = prevTris
				}
			}
		}
		wg.Done()
	}()

	wgRet.Add(1)
	go func() { // Merge all close together vertices (modifying generated triangles)
		wg.Wait()
		// Only the subset of triangles that are close to the seams (others already published)
		trianglesToPublish := map[Triangle3]*Triangle3{}
		// Build K-D tree for this pass
		allVertices := make(kdtree.Points, 0, len(vertToPassOriginalTris))
		for vert, _ := range vertToPassOriginalTris {
			allVertices = append(allVertices, mtToKdPoint(vert))
		}
		tree := kdtree.New(allVertices, false)
		//log.Println("Post-processing vertex count", len(vertToPassOriginalTris))
		movedVertices := make([]sdf.V3Set, 0) // {from, to}
		for vert, originalTriangles := range vertToPassOriginalTris {
			//log.Println("vertToPassOriginalTris", vert)
			closest := kdtree.NewNKeeper(3) // The first is always a perfect match with the current vertex
			tree.NearestSet(closest, mtToKdPoint(vert))
			for _, comparableDist := range closest.Heap[1:] {
				nthPos := mtFromKdPoint(comparableDist.Comparable.(kdtree.Point))
				nthDist := comparableDist.Dist
				modifiedTriangles := make([]*Triangle3, 0)
				for _, originalTriangle := range originalTriangles {
					modifiedTriangle, ok := trianglesToPublish[*originalTriangle]
					if !ok {
						modifiedTriangle = &*originalTriangle // Clone
					}
					modifiedTriangles = append(modifiedTriangles, modifiedTriangle)
				}
				if nthDist > mergeVerticesEpsilonSq {
					for i, originalTriangle := range originalTriangles {
						trianglesToPublish[*originalTriangle] = modifiedTriangles[i]
					}
					break
				} else if nthDist > 0 {
					//log.Println("Nth closest #", i, "tris", len(originalTriangles), "dist", nthDist, "from", vert, "to", nthPos)
					if nthPos.X > vert.X || nthPos.X == vert.X && (nthPos.Y > vert.Y || nthPos.Y == vert.Y && (nthPos.Z > vert.Z)) {
						continue // Only merge in a specific direction to avoid too many merges
					}
					// Get the modified triangle (we might have already changed another vertex)
					// Move vert to the found closest triangle
					for i, originalTriangle := range originalTriangles {
						modifiedTriangle := modifiedTriangles[i]
						if modifiedTriangle == nil {
							continue // Deleted triangle
						}
						for i, v := range modifiedTriangle.V {
							if v == vert {
								modifiedTriangle.V[i] = nthPos
								movedVertices = append(movedVertices, sdf.V3Set{vert, nthPos})
								break
							}
						}
						if modifiedTriangle.Degenerate(0) { // We created a degenerate triangle by merging vertices
							trianglesToPublish[*originalTriangle] = nil
						} else {
							trianglesToPublish[*originalTriangle] = modifiedTriangle
						}
					}
					break
				}
			}
		}
		for _, vertMoved := range movedVertices {
			prevTris := vertToPassOriginalTris[vertMoved[1]]
			addedTris := vertToPassOriginalTris[vertMoved[0]]
			prevTris = append(prevTris, addedTris...)
			vertToPassOriginalTris[vertMoved[1]] = mtRemoveDuplicatesAndNils(prevTris)
			delete(vertToPassOriginalTris, vertMoved[0])
		}
		for _, modifiedTriangle := range trianglesToPublish {
			if modifiedTriangle == nil {
				continue // Deleted
			}
			output <- modifiedTriangle
		}
		// DO NOT CLOSE: will be closed by parent: close(output)
		wgRet.Done()
	}()

	return input, wgRet
}

func bypassVertTris(isVertCloseToSeam func(v sdf.V3) bool, v *sdf.V3, tris []*Triangle3) bool {
	bypass := v == nil || !isVertCloseToSeam(*v)
	if bypass {
		for _, tri := range tris {
			bypass = bypass && !isVertCloseToSeam(tri.V[0]) && !isVertCloseToSeam(tri.V[1]) && !isVertCloseToSeam(tri.V[2])
			if !bypass {
				break
			}
		}
	}
	return bypass
}

func mtRemoveDuplicatesAndNils(strList []*Triangle3) []*Triangle3 {
	var list []*Triangle3
	for _, item := range strList {
		if item != nil && mtContains(list, item) == false {
			list = append(list, item)
		}
	}
	return list
}
func mtContains(s []*Triangle3, e *Triangle3) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
