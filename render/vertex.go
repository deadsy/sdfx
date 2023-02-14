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

		// Triangle corners.
		// Pre-initialize to make the next loop faster.
		var a v3.Vec
		var b v3.Vec
		var c v3.Vec

		// Read triangles from the channel and write them to vertices
		for ts := range writer {
			for _, t := range ts {
				a = t.V[0]
				b = t.V[1]
				c = t.V[2]
				*vertexes = append(*vertexes, a, b, c)
			}
		}
	}()

	return writer
}
