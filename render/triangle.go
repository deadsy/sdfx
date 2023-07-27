//-----------------------------------------------------------------------------
/*

Triangles

*/
//-----------------------------------------------------------------------------

package render

import (
	"sync"

	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// Triangle2 is a 2D triangle
type Triangle2 [3]v2.Vec

//-----------------------------------------------------------------------------

// Triangle3 is a 3D triangle
type Triangle3 [3]v3.Vec

// Normal returns the normal vector to the plane defined by the 3D triangle.
func (t *Triangle3) Normal() v3.Vec {
	e1 := t[1].Sub(t[0])
	e2 := t[2].Sub(t[0])
	return e1.Cross(e2).Normalize()
}

// Degenerate returns true if the triangle is degenerate.
func (t *Triangle3) Degenerate(tolerance float64) bool {
	// check for identical vertices
	if t[0].Equals(t[1], tolerance) {
		return true
	}
	if t[1].Equals(t[2], tolerance) {
		return true
	}
	if t[2].Equals(t[0], tolerance) {
		return true
	}
	// TODO more tests needed
	return false
}

//-----------------------------------------------------------------------------

// writeTriangles writes a stream of triangles to a slice.
func writeTriangles(wg *sync.WaitGroup, triangles *[]Triangle3) chan<- []*Triangle3 {
	// External code writes triangles to this channel.
	// This goroutine reads the channel and appends the triangles to a slice.
	c := make(chan []*Triangle3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// read triangles from the channel and append them to the slice
		for ts := range c {
			for _, t := range ts {
				*triangles = append(*triangles, *t)
			}
		}
	}()

	return c
}

//-----------------------------------------------------------------------------
