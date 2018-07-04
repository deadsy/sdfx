// -*- compile-command: "go build && ./base && fstl base.stl"; -*-

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

// All dimensions in mm
const (
	utronEdge   = 50.0
	utronMargin = 5.0

	magnetMargin = 10.0
	gapWidth     = 50.0
	innerGap     = 70.0
	magnetHeight = 101.6
	magnetDiam   = 50.8

	baseHeight    = 11.0
	wallThickness = 16.0

	bearingHeight   = 5.0
	bearingDiam     = 14.0
	bearingMargin   = 0.75
	bearingOverhang = 2.0
)

func main() {
	utronDiam := math.Sqrt(2 * utronEdge * utronEdge)

	// center of lower bearing is the origin.
	inside := utronDiam + 2*utronMargin
	outside := inside + 2*wallThickness
	inbox := Box3D(V3{inside, inside, 2 * outside}, 0)
	inbox = Transform3D(inbox, Translate3d(V3{0, 0, outside}))
	boxHeight := wallThickness - 1.5*bearingHeight + utronDiam
	box := Box3D(V3{outside, outside, boxHeight}, 0)
	box = Transform3D(box, Translate3d(V3{0, 0, 0.5*boxHeight - wallThickness}))
	box = Difference3D(box, inbox)
	box = Transform3D(box, Translate3d(V3{0, 0, 0.5 * bearingHeight}))
	// left cutout
	cutBox := Box3D(V3{outside, outside, outside}, 0)
	cutBox = Transform3D(cutBox, Translate3d(V3{-0.5 * outside, 0, 0.5*outside + 0.5*utronDiam - baseHeight}))
	box = Difference3D(box, cutBox)
	// lower magnet brace
	dx := math.Sqrt(2 * utronMargin * utronMargin)
	ts := 0.5*outside - wallThickness
	triangle := Polygon2D([]V2{{dx, 0}, {ts + dx, 0}, {ts + dx, ts}})
	prism := Extrude3D(triangle, outside)
	prism = Transform3D(prism, RotateX(0.5*math.Pi))
	// prism = Transform3D(prism,
	box = Union3D(box, prism)

	bearing := Cylinder3D(bearingHeight, 0.5*(bearingDiam+bearingMargin), 0)
	access := Cylinder3D(wallThickness, 0.5*(bearingDiam-bearingOverhang), 0)
	access = Transform3D(access, Translate3d(V3{0, 0, -0.5 * wallThickness}))
	bearingCutout := Union3D(bearing, access)

	s := Difference3D(box, bearingCutout)
	RenderSTL(s, 200, "base.stl")
}
