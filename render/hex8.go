package render

import (
	"sync"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Hex8 is a 3D hexahedron consisting of 8 nodes.
// It's a kind of finite element, FE.
// https://en.wikipedia.org/wiki/Hexahedron
type Hex8 struct {
	// Coordinates of 8 corner nodes or vertices.
	V [8]v3.Vec
	// The layer to which tetrahedron belongs. Layers are along Z axis.
	// For finite element analysis - FEA - of 3D printed objects, it's more efficient to store layer along Z axis.
	// The 3D print is done along the Z axis. Likewise, FEA is done along the Z axis.
	// Sampling/marching algorithm is expected to return the layer to which a finite element belongs.
	Layer int
}

//-----------------------------------------------------------------------------

// writeHex8 writes a stream of finite elements to an array.
func writeHex8(wg *sync.WaitGroup, hex8s *[]Hex8) chan<- []*Hex8 {
	// External code writes to this channel.
	// This goroutine reads the channel and stores finite elements.
	c := make(chan []*Hex8)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// read finite elements from the channel and handle them
		for fes := range c {
			for _, fe := range fes {
				*hex8s = append(*hex8s, *fe)
			}
		}
	}()

	return c
}

//-----------------------------------------------------------------------------
