package render

import (
	"sync"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Fe is a finite element.
type Fe struct {
	// Coordinates of nodes or vertices.
	V []v3.Vec
	// Coordinates of the voxel to which the element belongs.
	X int
	Y int
	Z int
}

//-----------------------------------------------------------------------------

// writeFe writes a stream of finite elements to an array.
func writeFe(wg *sync.WaitGroup, elements *[]Fe) chan<- []*Fe {
	// External code writes to this channel.
	// This goroutine reads the channel and stores finite elements.
	c := make(chan []*Fe)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// read finite elements from the channel and handle them
		for fes := range c {
			for _, fe := range fes {
				*elements = append(*elements, *fe)
			}
		}
	}()

	return c
}

//-----------------------------------------------------------------------------
