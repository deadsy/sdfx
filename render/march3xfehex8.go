package render

import (
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

// Process a cube. Generate finite elements, or more cubes.
func (dc *dcache3) processCubeHex8(c *cube, output chan<- []*Hex8) {
	if !dc.isEmpty(c) {
		if c.n == 1 {
			// this cube is at the required resolution
			c0, d0 := dc.evaluate(c.v.Add(v3i.Vec{0, 0, 0}))
			c1, d1 := dc.evaluate(c.v.Add(v3i.Vec{2, 0, 0}))
			c2, d2 := dc.evaluate(c.v.Add(v3i.Vec{2, 2, 0}))
			c3, d3 := dc.evaluate(c.v.Add(v3i.Vec{0, 2, 0}))
			c4, d4 := dc.evaluate(c.v.Add(v3i.Vec{0, 0, 2}))
			c5, d5 := dc.evaluate(c.v.Add(v3i.Vec{2, 0, 2}))
			c6, d6 := dc.evaluate(c.v.Add(v3i.Vec{2, 2, 2}))
			c7, d7 := dc.evaluate(c.v.Add(v3i.Vec{0, 2, 2}))
			corners := [8]v3.Vec{c0, c1, c2, c3, c4, c5, c6, c7}
			values := [8]float64{d0, d1, d2, d3, d4, d5, d6, d7}
			// output the triangle(s) for this cube
			output <- mcToHex8(corners, values, 0, 0)
		} else {
			// process the sub cubes
			n := c.n - 1
			s := 1 << n
			// TODO - turn these into throttled go-routines
			dc.processCubeHex8(&cube{c.v.Add(v3i.Vec{0, 0, 0}), n}, output)
			dc.processCubeHex8(&cube{c.v.Add(v3i.Vec{s, 0, 0}), n}, output)
			dc.processCubeHex8(&cube{c.v.Add(v3i.Vec{s, s, 0}), n}, output)
			dc.processCubeHex8(&cube{c.v.Add(v3i.Vec{0, s, 0}), n}, output)
			dc.processCubeHex8(&cube{c.v.Add(v3i.Vec{0, 0, s}), n}, output)
			dc.processCubeHex8(&cube{c.v.Add(v3i.Vec{s, 0, s}), n}, output)
			dc.processCubeHex8(&cube{c.v.Add(v3i.Vec{s, s, s}), n}, output)
			dc.processCubeHex8(&cube{c.v.Add(v3i.Vec{0, s, s}), n}, output)
		}
	}
}
