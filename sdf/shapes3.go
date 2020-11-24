//-----------------------------------------------------------------------------
/*

Common 3D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"fmt"
	"math"
)

//-----------------------------------------------------------------------------

// CounterBoredHole3D returns the SDF3 for a counterbored hole.
func CounterBoredHole3D(
	l float64, // total length
	r float64, // hole radius
	cbRadius float64, // counter bore radius
	cbDepth float64, // counter bore depth
) SDF3 {
	s0 := Cylinder3D(l, r, 0)
	s1 := Cylinder3D(cbDepth, cbRadius, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - cbDepth) / 2}))
	return Union3D(s0, s1)
}

// ChamferedHole3D returns the SDF3 for a chamfered hole (45 degrees).
func ChamferedHole3D(
	l float64, // total length
	r float64, // hole radius
	chRadius float64, // chamfer radius
) SDF3 {
	s0 := Cylinder3D(l, r, 0)
	s1 := Cone3D(chRadius, r, r+chRadius, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - chRadius) / 2}))
	return Union3D(s0, s1)
}

// CounterSunkHole3D returns the SDF3 for a countersunk hole (45 degrees).
func CounterSunkHole3D(
	l float64, // total length
	r float64, // hole radius
) SDF3 {
	return ChamferedHole3D(l, r, r)
}

//-----------------------------------------------------------------------------

// HexHead3D returns the rounded hex head for a nut or bolt.
func HexHead3D(
	r float64, // radius
	h float64, // height
	round string, // (t)top, (b)bottom, (tb)top/bottom
) SDF3 {
	// basic hex body
	cornerRound := r * 0.08
	hex2d := Polygon2D(Nagon(6, r-cornerRound))
	hex2d = Offset2D(hex2d, cornerRound)
	hex3d := Extrude3D(hex2d, h)
	// round out the top and/or bottom as required
	if round != "" {
		topRound := r * 1.6
		d := r * math.Cos(DtoR(30))
		sphere3d := Sphere3D(topRound)
		zOfs := math.Sqrt(topRound*topRound-d*d) - h/2
		if round == "t" || round == "tb" {
			hex3d = Intersect3D(hex3d, Transform3D(sphere3d, Translate3d(V3{0, 0, -zOfs})))
		}
		if round == "b" || round == "tb" {
			hex3d = Intersect3D(hex3d, Transform3D(sphere3d, Translate3d(V3{0, 0, zOfs})))
		}
	}
	return hex3d
}

// KnurledHead3D returns a cylindrical knurled head.
func KnurledHead3D(
	r float64, // radius
	h float64, // height
	pitch float64, // knurl pitch
) SDF3 {
	theta := DtoR(45)
	cylinderRound := r * 0.05
	knurlH := pitch * math.Floor((h-cylinderRound)/pitch)
	knurl3d := Knurl3D(knurlH, r, pitch, pitch*0.3, theta)
	return Union3D(Cylinder3D(h, r, cylinderRound), knurl3d)
}

//-----------------------------------------------------------------------------

// knurlProfile returns a 2D knurl profile.
func knurlProfile(
	radius float64, // radius of knurled cylinder
	pitch float64, // pitch of the knurl
	height float64, // height of the knurl
) SDF2 {
	knurl := NewPolygon()
	knurl.Add(pitch/2, 0)
	knurl.Add(pitch/2, radius)
	knurl.Add(0, radius+height)
	knurl.Add(-pitch/2, radius)
	knurl.Add(-pitch/2, 0)
	//knurl.Render("knurl.dxf")
	return Polygon2D(knurl.Vertices())
}

// Knurl3D returns a knurled cylinder.
func Knurl3D(
	length float64, // length of cylinder
	radius float64, // radius of cylinder
	pitch float64, // knurl pitch
	height float64, // knurl height
	theta float64, // knurl helix angle
) SDF3 {
	// A knurl is the the intersection of left and right hand
	// multistart "threads". Work out the number of starts using
	// the desired helix angle.
	n := int(Tau * radius * math.Tan(theta) / pitch)
	// build the knurl profile.
	knurl2d := knurlProfile(radius, pitch, height)
	// create the left/right hand spirals
	knurl0_3d := Screw3D(knurl2d, length, pitch, n)
	knurl1_3d := Screw3D(knurl2d, length, pitch, -n)
	return Intersect3D(knurl0_3d, knurl1_3d)
}

//-----------------------------------------------------------------------------

// ChamferedCylinder intersects a chamfered cylinder with an SDF3.
func ChamferedCylinder(s SDF3, kb, kt float64) SDF3 {
	// get the length and radius from the bounding box
	l := s.BoundingBox().Max.Z
	r := s.BoundingBox().Max.X
	p := NewPolygon()
	p.Add(0, -l)
	p.Add(r, -l).Chamfer(r * kb)
	p.Add(r, l).Chamfer(r * kt)
	p.Add(0, l)
	return Intersect3D(s, Revolve3D(Polygon2D(p.Vertices())))
}

//-----------------------------------------------------------------------------
// Simple Bolt for 3d printing.

// BoltParms defines the parameters for a bolt.
type BoltParms struct {
	Thread      string  // name of thread
	Style       string  // head style "hex" or "knurl"
	Tolerance   float64 // subtract from external thread radius
	TotalLength float64 // threaded length + shank length
	ShankLength float64 // non threaded length
}

// Bolt returns a simple bolt suitable for 3d printing.
func Bolt(k *BoltParms) (SDF3, error) {
	// validate parameters
	t, err := ThreadLookup(k.Thread)
	if err != nil {
		return nil, err
	}
	if k.TotalLength < 0 {
		return nil, errors.New("total length < 0")
	}
	if k.ShankLength < 0 {
		return nil, errors.New("shank length < 0")
	}
	if k.Tolerance < 0 {
		return nil, errors.New("tolerance < 0")
	}

	// head
	var head SDF3

	hr := t.HexRadius()
	hh := t.HexHeight()
	switch k.Style {
	case "hex":
		head = HexHead3D(hr, hh, "b")
	case "knurl":
		head = KnurledHead3D(hr, hh, hr*0.25)
	default:
		return nil, fmt.Errorf("unknown style \"%s\"", k.Style)
	}

	// shank
	shankLength := k.ShankLength + hh/2
	shankOffset := shankLength / 2
	shank := Cylinder3D(shankLength, t.Radius, hh*0.08)
	shank = Transform3D(shank, Translate3d(V3{0, 0, shankOffset}))

	// external thread
	threadLength := k.TotalLength - k.ShankLength
	if threadLength < 0 {
		threadLength = 0
	}
	var thread SDF3
	if threadLength != 0 {
		r := t.Radius - k.Tolerance
		threadOffset := threadLength/2 + shankLength
		isoThread := ISOThread(r, t.Pitch, true)
		thread = Screw3D(isoThread, threadLength, t.Pitch, 1)
		// chamfer the thread
		thread = ChamferedCylinder(thread, 0, 0.5)
		thread = Transform3D(thread, Translate3d(V3{0, 0, threadOffset}))
	}

	return Union3D(head, shank, thread), nil
}

//-----------------------------------------------------------------------------
// Simple Nut for 3d printing.

// NutParms defines the parameters for a nut.
type NutParms struct {
	Thread    string  // name of thread
	Style     string  // head style "hex" or "knurl"
	Tolerance float64 // add to internal thread radius
}

// Nut returns a simple nut suitable for 3d printing.
func Nut(k *NutParms) (SDF3, error) {
	// validate parameters
	t, err := ThreadLookup(k.Thread)
	if err != nil {
		return nil, err
	}
	if k.Tolerance < 0 {
		return nil, errors.New("tolerance < 0")
	}

	// nut body
	var nut SDF3
	nr := t.HexRadius()
	nh := t.HexHeight()
	switch k.Style {
	case "hex":
		nut = HexHead3D(nr, nh, "tb")
	case "knurl":
		nut = KnurledHead3D(nr, nh, nr*0.25)
	default:
		return nil, fmt.Errorf("unknown style \"%s\"", k.Style)
	}

	// internal thread
	isoThread := ISOThread(t.Radius+k.Tolerance, t.Pitch, false)
	thread := Screw3D(isoThread, nh, t.Pitch, 1)

	return Difference3D(nut, thread), nil
}

//-----------------------------------------------------------------------------
