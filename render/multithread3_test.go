package render

import (
	"github.com/deadsy/sdfx/sdf"
	"log"
	"sync"
	"testing"
)

func TestMergeVertices(t *testing.T) {
	// Config
	mergedOut := make(chan *Triangle3)
	vertDist := 0.1
	trianglesIn, wg := NewMtRenderer3(nil, 1).mergeGeneratedVertices(mergedOut, vertDist,
		sdf.V3{X: 0, Y: -1, Z: -0.001}, sdf.V3{X: 1, Y: 1, Z: 0.002}, sdf.V3{X: 1, Y: 1, Z: 0.002})

	// Input data
	for _, triIn := range []*Triangle3{ // 4 triangles forming a rectangle with 2 split by the middle
		{V: [3]sdf.V3{{0, vertDist, 0}, {1, vertDist, 0}, {1, 1, 0}}},
		{V: [3]sdf.V3{{0, vertDist, 0}, {1, 1, 0}, {0, 1, 0}}},
		{V: [3]sdf.V3{{0, 0, 0}, {1, 0, 0}, {1, -1, 0}}},
		{V: [3]sdf.V3{{0, 0, 0}, {1, -1, 0}, {0, -1, 0}}},
	} {
		trianglesIn <- triIn
	}
	close(trianglesIn)

	// Output processing
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for tri := range mergedOut {
			log.Println(tri) // Should print 4 triangles with no hole in the middle
		}
	}()
	wg.Wait()
	close(mergedOut)
	wg2.Wait()
}
