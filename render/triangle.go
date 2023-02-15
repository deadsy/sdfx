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
type Triangle3 struct {
	V [3]v3.Vec
}

// NewTriangle3 returns a new 3D triangle.
func NewTriangle3(a, b, c v3.Vec) *Triangle3 {
	t := Triangle3{}
	t.V[0] = a
	t.V[1] = b
	t.V[2] = c
	return &t
}

// Normal returns the normal vector to the plane defined by the 3D triangle.
func (t *Triangle3) Normal() v3.Vec {
	e1 := t.V[1].Sub(t.V[0])
	e2 := t.V[2].Sub(t.V[0])
	return e1.Cross(e2).Normalize()
}

// Degenerate returns true if the triangle is degenerate.
func (t *Triangle3) Degenerate(tolerance float64) bool {
	// check for identical vertices
	if t.V[0].Equals(t.V[1], tolerance) {
		return true
	}
	if t.V[1].Equals(t.V[2], tolerance) {
		return true
	}
	if t.V[2].Equals(t.V[0], tolerance) {
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
