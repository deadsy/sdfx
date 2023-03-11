package render

import (
	"sync"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Hex20 is a 3D hexahedron consisting of 20 nodes.
// It's a kind of finite element, FE.
// https://en.wikipedia.org/wiki/Hexahedron
type Hex20 struct {
	// Coordinates of 20 corner nodes or vertices.
	V [20]v3.Vec
	// The layer to which tetrahedron belongs. Layers are along Z axis.
	// For finite element analysis - FEA - of 3D printed objects, it's more efficient to store layer along Z axis.
	// The 3D print is done along the Z axis. Likewise, FEA is done along the Z axis.
	// Sampling/marching algorithm is expected to return the layer to which a finite element belongs.
	Layer int
}

//-----------------------------------------------------------------------------

// writeHex20 writes a stream of finite elements to an array.
func writeHex20(wg *sync.WaitGroup, hex20s *[]Hex20) chan<- []*Hex20 {
	// External code writes to this channel.
	// This goroutine reads the channel and stores finite elements.
	c := make(chan []*Hex20)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// read finite elements from the channel and handle them
		for fes := range c {
			for _, fe := range fes {
				*hex20s = append(*hex20s, *fe)
			}
		}
	}()

	return c
}

//-----------------------------------------------------------------------------
