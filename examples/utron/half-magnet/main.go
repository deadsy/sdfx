// -*- compile-command: "go build && ./half-magnet && fstl half-magnet.stl"; -*-

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

// All dimensions in mm
const (
	magnetMargin = 10.0
	gapWidth     = 50.0
	innerGap     = 70.0
	magnetHeight = 101.6
	magnetDiam   = 50.8
)

func main() {
	r := 0.5 * (innerGap + magnetDiam)
	torus := torus3D(0.5*magnetDiam, r)
	block := Box3D(V3{4 * r, 2 * r, 2 * r}, 0)
	block = Transform3D(block, Translate3d(V3{0, r, 0}))
	halfTorus := Difference3D(torus, block)

	// straight section
	ssHeight := 0.5*(magnetHeight-gapWidth) - magnetMargin
	ss := Cylinder3D(ssHeight, 0.5*magnetDiam, 0)
	ss = Transform3D(ss, RotateX(0.5*math.Pi))
	ss = Transform3D(ss, Translate3d(V3{r, 0.5 * ssHeight, 0}))

	s := Union3D(halfTorus, ss)
	RenderSTL(s, 200, "half-magnet.stl")
}

func torus3D(minorRadius, majorRadius float64) SDF3 {
	circle := Circle2D(minorRadius)
	circle = Transform2D(circle, Translate2d(V2{majorRadius, 0}))
	return Revolve3D(circle)
}
