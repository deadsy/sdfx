package render

import (
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

// Process a cube. Generate finite elements, or more cubes.
func (dc *dcache3) processCubeHex8(c *cube, output chan<- []*Hex8) {
	if !dc.isEmpty(c) {

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

		anyPositive := false
		for i := 0; i < 8; i++ {
			if values[i] > 0 {
				anyPositive = true
				break
			}
		}

		if !anyPositive {
			fe := Hex8{
				V:     [8]v3.Vec{},
				Layer: 0,
			}

			// Refer to CalculiX solver documentation:
			// http://www.dhondt.de/ccx_2.20.pdf

			fe.V[7] = corners[7]
			fe.V[6] = corners[6]
			fe.V[5] = corners[5]
			fe.V[4] = corners[4]
			fe.V[3] = corners[3]
			fe.V[2] = corners[2]
			fe.V[1] = corners[1]
			fe.V[0] = corners[0]

			// output the finite element(s) for this cube
			output <- []*Hex8{&fe}
		}

		if c.n == 1 {
			// this cube is at the required resolution
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
