package render

import (
	"sync"
)

// Writes a stream of triangles.
//
// Pass slice by pointer. Because the function adds new elements to the slice,
// that requires changing the slice header, which the caller will not see.
func writeTriangles(wg *sync.WaitGroup, triangles *[]Triangle3) chan<- []*Triangle3 {
	// External code writes triangles to this channel.
	// This goroutine reads the channel and re-writes triangles.
	writer := make(chan []*Triangle3)

	// Write by a goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Read triangles from the channel and re-write them.
		for ts := range writer {
			for _, t := range ts {
				*triangles = append(*triangles, *t)
			}
		}
	}()

	return writer
}
