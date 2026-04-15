//-----------------------------------------------------------------------------
/*

Mesh Analysis Utilities

Functions for collecting rendered triangles and checking mesh quality
(watertightness, boundary edges).

*/
//-----------------------------------------------------------------------------

package render

import (
	"math"
	"sync"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// CollectTriangles renders an SDF3 with the given renderer and returns all
// triangles as a flat slice.
func CollectTriangles(s sdf.SDF3, r Render3) []sdf.Triangle3 {
	ch := make(chan []*sdf.Triangle3)
	var tris []sdf.Triangle3
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for batch := range ch {
			for _, t := range batch {
				tris = append(tris, *t)
			}
		}
	}()
	r.Render(s, sdf.NewTriangle3Buffer(ch))
	close(ch)
	wg.Wait()
	return tris
}

//-----------------------------------------------------------------------------

// meshVertex is a quantized 3D position for vertex deduplication.
// Quantized to 1e-4 resolution — well below meaningful geometry scale
// but avoids overflow in int32 for typical part sizes (up to ~200mm).
type meshVertex struct {
	x, y, z int32
}

func quantizeVertex(x, y, z float64) meshVertex {
	return meshVertex{int32(x * 1e4), int32(y * 1e4), int32(z * 1e4)}
}

// meshEdge is an unordered pair of vertices.
type meshEdge struct {
	a, b meshVertex
}

// makeMeshEdge returns a canonical (sorted) edge so (a,b) == (b,a).
func makeMeshEdge(a, b meshVertex) meshEdge {
	if a.x < b.x || (a.x == b.x && a.y < b.y) || (a.x == b.x && a.y == b.y && a.z < b.z) {
		return meshEdge{a, b}
	}
	return meshEdge{b, a}
}

//-----------------------------------------------------------------------------

// CountBoundaryEdges returns the number of edges shared by only one triangle.
// A watertight mesh has zero boundary edges — every edge borders exactly
// two triangles.
func CountBoundaryEdges(tris []sdf.Triangle3) int {
	edgeCount := make(map[meshEdge]int, len(tris)*3)
	for _, t := range tris {
		v0 := quantizeVertex(t[0].X, t[0].Y, t[0].Z)
		v1 := quantizeVertex(t[1].X, t[1].Y, t[1].Z)
		v2 := quantizeVertex(t[2].X, t[2].Y, t[2].Z)
		edgeCount[makeMeshEdge(v0, v1)]++
		edgeCount[makeMeshEdge(v1, v2)]++
		edgeCount[makeMeshEdge(v2, v0)]++
	}
	boundary := 0
	for _, count := range edgeCount {
		if count == 1 {
			boundary++
		}
	}
	return boundary
}

// IsWatertight returns true if the rendered mesh has no boundary edges.
func IsWatertight(tris []sdf.Triangle3) bool {
	return CountBoundaryEdges(tris) == 0
}

//-----------------------------------------------------------------------------

// MaxZ returns the maximum Z coordinate across all triangle vertices.
func MaxZ(tris []sdf.Triangle3) float64 {
	z := -math.MaxFloat64
	for _, t := range tris {
		for _, v := range t {
			if v.Z > z {
				z = v.Z
			}
		}
	}
	return z
}

//-----------------------------------------------------------------------------
