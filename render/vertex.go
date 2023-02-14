package render

import (
	"sync"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Writes a stream of triangles, keeping only the raw vertexes.
//
// Pass slice by pointer. Because the function adds new elements to the slice,
// that requires changing the slice header, which the caller will not see.
func writeVertexes(wg *sync.WaitGroup, vertexes *[]v3.Vec) chan<- []*Triangle3 {
	// External code writes triangles to this channel.
	// This goroutine reads the channel and writes triangles to vertices.
	writer := make(chan []*Triangle3)

	// Write by a goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Read triangles from the channel and write their vertices
		for ts := range writer {
			for _, t := range ts {
				*vertexes = append(*vertexes, t.V[0], t.V[1], t.V[2])
			}
		}
	}()

	return writer
}
