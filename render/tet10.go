package render

import (
	"sync"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Tet10 is a 3D tetrahedron consisting of 10 nodes.
// It's a kind of finite element, FE.
// https://en.wikipedia.org/wiki/Tetrahedron
type Tet10 struct {
	// Coordinates of corner nodes or vertices.
	V [10]v3.Vec
	// The Layer to which tetrahedron belongs. Layers are along Z axis.
	// For finite element analysis - FEA - of 3D printed objects, it's more efficient to store Layer along Z axis.
	// The 3D print is done along the Z axis. Likewise, FEA is done along the Z axis.
	// Sampling/marching algorithm is expected to return the Layer to which a finite element belongs.
	Layer int
}

//-----------------------------------------------------------------------------

// writeTet10 writes a stream of finite elements to an array.
func writeTet10(wg *sync.WaitGroup, tet10s *[]Tet10) chan<- []*Tet10 {
	// External code writes to this channel.
	// This goroutine reads the channel and stores finite elements.
	c := make(chan []*Tet10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// read finite elements from the channel and handle them
		for fes := range c {
			for _, fe := range fes {
				*tet10s = append(*tet10s, *fe)
			}
		}
	}()

	return c
}

//-----------------------------------------------------------------------------
