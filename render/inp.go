package render

import (
	"sync"
)

//-----------------------------------------------------------------------------

// Define the ABAQUS or CalculiX inp file requirements, if any.

//-----------------------------------------------------------------------------

// writeInpTet4 writes a stream of finite elements in the shape of tetrahedra to an ABAQUS or CalculiX `inp` file.
func writeInpTet4(wg *sync.WaitGroup, path string) (chan<- []*Tet4, error) {
	// External code writes tetrahedra to this channel.
	// This goroutine reads the channel and writes tetrahedra to the file.
	c := make(chan []*Tet4)

	wg.Add(1)
	go func() {
		defer wg.Done()

		m := NewMeshTet4()
		defer m.WriteInp(path)

		// read finite elements from the channel and handle them
		for ts := range c {
			for _, t := range ts {
				m.AddTet4(t.V[0], t.V[1], t.V[2], t.V[3])
			}
		}
	}()

	return c, nil
}
