// -*- compile-command: "go build && ./utron && fstl utron.stl"; -*-

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
	"github.com/gmlewis/sdfx/examples/utron/enclosure"
	"github.com/gmlewis/sdfx/examples/utron/half-magnet"
	"github.com/gmlewis/sdfx/examples/utron/half-utron"
)

// All dimensions in mm
const (
	utronEdge    = 50.0
	magnetHeight = 25.4
	innerGap     = 70.0
	magnetDiam   = 50.8
	metalMargin  = 0.5
	magnetMargin = 10.0
)

func main() {
	utronRadius := 0.5 * math.Sqrt(2*utronEdge*utronEdge)

	base := enclosure.Base(utronEdge)
	ch := 4 * magnetHeight
	baseCutout := Cylinder3D(ch, 0.5*magnetDiam+metalMargin, 1)
	ssHeight := 0.5*(4*magnetHeight-utronEdge) - magnetMargin
	m := Translate3d(V3{0, 0, -0.5*ch - 2*magnetHeight + ssHeight + metalMargin})
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	baseCutout = Transform3D(baseCutout, m)
	base = Difference3D(base, baseCutout)

	halfUtron := half_utron.HalfUtron(utronEdge)
	utronLower := Transform3D(halfUtron, RotateX(math.Pi))
	utronLower = Transform3D(utronLower, Translate3d(V3{0, 0, utronRadius}))
	utronUpper := Transform3D(halfUtron, Translate3d(V3{0, 0, utronRadius}))

	halfMagnet := half_magnet.HalfMagnet(utronEdge, innerGap, magnetDiam, magnetHeight, magnetMargin)
	m = RotateX(0.5 * math.Pi)
	m = Translate3d(V3{-0.5 * (innerGap + magnetDiam), 0, -2 * magnetHeight}).Mul(m)
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	halfMagnetLower := Transform3D(halfMagnet, m)
	m = RotateX(-0.5 * math.Pi)
	m = Translate3d(V3{-0.5 * (innerGap + magnetDiam), 0, 2 * magnetHeight}).Mul(m)
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	halfMagnetUpper := Transform3D(halfMagnet, m)

	magnet1 := Cylinder3D(magnetHeight, 0.5*magnetDiam, 1)
	magnet1 = Transform3D(magnet1, Translate3d(V3{0, 0, -1.5 * magnetHeight}))
	magnet2 := Cylinder3D(magnetHeight, 0.5*magnetDiam, 1)
	magnet2 = Transform3D(magnet2, Translate3d(V3{0, 0, -0.5 * magnetHeight}))
	magnet3 := Cylinder3D(magnetHeight, 0.5*magnetDiam, 1)
	magnet3 = Transform3D(magnet3, Translate3d(V3{0, 0, 0.5 * magnetHeight}))
	magnet4 := Cylinder3D(magnetHeight, 0.5*magnetDiam, 1)
	magnet4 = Transform3D(magnet4, Translate3d(V3{0, 0, 1.5 * magnetHeight}))
	magnets := Union3D(magnet1, magnet2, magnet3, magnet4)
	m = Translate3d(V3{-innerGap - magnetDiam, 0, 0})
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	magnets = Transform3D(magnets, m)

	s := Union3D(base, utronLower, utronUpper, halfMagnetLower, halfMagnetUpper, magnets)
	RenderSTL(s, 400, "utron.stl")
}
