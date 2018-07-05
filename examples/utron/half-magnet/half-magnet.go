package half_magnet

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

// All dimensions in mm

func HalfMagnet(utronEdge, innerGap, magnetDiam, magnetHeight, magnetMargin float64) SDF3 {
	r := 0.5 * (innerGap + magnetDiam)
	torus := torus3D(0.5*magnetDiam, r)
	block := Box3D(V3{4 * r, 2 * r, 2 * r}, 0)
	block = Transform3D(block, Translate3d(V3{0, r, 0}))
	halfTorus := Difference3D(torus, block)

	// straight section
	ssHeight := 0.5*(4*magnetHeight-utronEdge) - magnetMargin
	// Add overlap to avoid chamfer at join
	overlap := 1.0
	ss := Cylinder3D(ssHeight+overlap, 0.5*magnetDiam, 0)
	ss = Transform3D(ss, RotateX(0.5*math.Pi))
	ss = Transform3D(ss, Translate3d(V3{r, 0.5*ssHeight - overlap, 0}))

	return Union3D(halfTorus, ss)
}

func torus3D(minorRadius, majorRadius float64) SDF3 {
	circle := Circle2D(minorRadius)
	circle = Transform2D(circle, Translate2d(V2{majorRadius, 0}))
	return Revolve3D(circle)
}
