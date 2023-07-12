package main

import v3 "github.com/deadsy/sdfx/vec/v3"

// The 3D beam is simply supported i.e. one end is pinned and the other end is roller.
// Benchmark reference:
// https://github.com/calculix/CalculiX-Examples/tree/master/NonLinear/Sections
func restraintSquare(x, y, z float64) (bool, bool, bool) {
	node := v3.Vec{X: x, Y: y, Z: z}

	// All three degrees of freedom are fixed.
	if node.Equals(v3.Vec{X: 0, Y: 0, Z: 0}, 2) {
		return true, true, true
	}
	if node.Equals(v3.Vec{X: 0, Y: 17.32, Z: 0}, 2) {
		return true, true, true
	}

	// Some degrees of freedom are fixed.
	if node.Equals(v3.Vec{X: 200, Y: 0, Z: 0}, 2) {
		return false, true, true
	}
	if node.Equals(v3.Vec{X: 200, Y: 17.32, Z: 0}, 2) {
		return false, true, true
	}
	return false, false, false
}

func loadSquare(x, y, z float64) (float64, float64, float64) {
	return 0, 0, 0
}
