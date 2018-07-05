package half_utron

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

// All dimensions in mm
const (
	minThickness = 3.0
)

func HalfUtron(utronEdge float64) SDF3 {
	cr := math.Sqrt(0.5 * utronEdge * utronEdge)
	cone := Cone3D(cr, 0, cr, 0.5)
	cone = Transform3D(cone, Rotate3d(V3{1, 0, 0}, math.Pi))
	cone = Transform3D(cone, Translate3d(V3{0, 0, 0.5 * cr}))

	sd := utronEdge - 2.0*minThickness
	sphere := Sphere3D(0.5 * sd)

	return Difference3D(cone, sphere)
}
