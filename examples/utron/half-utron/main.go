// -*- compile-command: "go build && ./half-utron && fstl half-utron.stl"; -*-

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

const utronEdge = 50.0   // mm
const minThickness = 3.0 // mm

func main() {
	cr := math.Sqrt(0.5 * utronEdge * utronEdge)
	cone := Cone3D(cr, 0, cr, 0.5)
	cone = Transform3D(cone, Rotate3d(V3{1, 0, 0}, math.Pi))
	cone = Transform3D(cone, Translate3d(V3{0, 0, 0.5 * cr}))

	sd := utronEdge - 2.0*minThickness
	sphere := Sphere3D(0.5 * sd)

	s := Difference3D(cone, sphere)
	RenderSTL(s, 200, "half-utron.stl")
}
