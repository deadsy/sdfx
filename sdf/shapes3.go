//-----------------------------------------------------------------------------
/*

Common 3D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

import (
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

// KnurlProfile returns a 2D knurl profile.
func KnurlProfile(
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
	knurl2d := KnurlProfile(radius, pitch, height)
	// create the left/right hand spirals
	knurl0_3d := Screw3D(knurl2d, length, pitch, n)
	knurl1_3d := Screw3D(knurl2d, length, pitch, -n)
	return Intersect3D(knurl0_3d, knurl1_3d)
}

//-----------------------------------------------------------------------------

// WasherParms defines the parameters for a washer.
type WasherParms struct {
	Thickness   float64 // thickness
	InnerRadius float64 // inner radius
	OuterRadius float64 // outer radius
	Remove      float64 // fraction of complete washer removed
}

// Washer3D returns a washer.
// This is also used to create circular walls.
func Washer3D(k *WasherParms) SDF3 {
	if k.Thickness <= 0 {
		panic("Thickness <= 0")
	}
	if k.InnerRadius >= k.OuterRadius {
		panic("InnerRadius >= OuterRadius")
	}
	if k.Remove < 0 || k.Remove >= 1.0 {
		panic("Remove must be [0..1)")
	}

	var s SDF3
	if k.Remove == 0 {
		// difference of cylinders
		outer := Cylinder3D(k.Thickness, k.OuterRadius, 0)
		inner := Cylinder3D(k.Thickness, k.InnerRadius, 0)
		s = Difference3D(outer, inner)
	} else {
		// build a 2d profile box
		dx := k.OuterRadius - k.InnerRadius
		dy := k.Thickness
		xofs := 0.5 * (k.InnerRadius + k.OuterRadius)
		b := Box2D(V2{dx, dy}, 0)
		b = Transform2D(b, Translate2d(V2{xofs, 0}))
		// rotate about the z-axis
		theta := Tau * (1.0 - k.Remove)
		s = RevolveTheta3D(b, theta)
		// center the removed portion on the x-axis
		dtheta := 0.5 * (Tau - theta)
		s = Transform3D(s, RotateZ(dtheta))
	}
	return s
}

//-----------------------------------------------------------------------------
// Board standoffs

// StandoffParms defines the parameters for a board standoff pillar.
type StandoffParms struct {
	PillarHeight   float64
	PillarDiameter float64
	HoleDepth      float64 // > 0 is a hole, < 0 is a support stub
	HoleDiameter   float64
	NumberWebs     int // number of triangular gussets around the standoff base
	WebHeight      float64
	WebDiameter    float64
	WebWidth       float64
}

// single web
func pillarWeb(k *StandoffParms) SDF3 {
	w := NewPolygon()
	w.Add(0, 0)
	w.Add(0.5*k.WebDiameter, 0)
	w.Add(0, k.WebHeight)
	s := Extrude3D(Polygon2D(w.Vertices()), k.WebWidth)
	m := Translate3d(V3{0, 0, -0.5 * k.PillarHeight}).Mul(RotateX(DtoR(90.0)))
	return Transform3D(s, m)
}

// multiple webs
func pillarWebs(k *StandoffParms) SDF3 {
	if k.NumberWebs == 0 {
		return nil
	}
	return RotateCopy3D(pillarWeb(k), k.NumberWebs)
}

// pillar
func pillar(k *StandoffParms) SDF3 {
	return Cylinder3D(k.PillarHeight, 0.5*k.PillarDiameter, 0)
}

// pillar hole
func pillarHole(k *StandoffParms) SDF3 {
	if k.HoleDiameter == 0.0 || k.HoleDepth == 0.0 {
		return nil
	}
	s := Cylinder3D(Abs(k.HoleDepth), 0.5*k.HoleDiameter, 0)
	zOfs := 0.5 * (k.PillarHeight - k.HoleDepth)
	return Transform3D(s, Translate3d(V3{0, 0, zOfs}))
}

// Standoff3D returns a single board standoff.
func Standoff3D(k *StandoffParms) SDF3 {
	s0 := Union3D(pillar(k), pillarWebs(k))
	if k.NumberWebs != 0 {
		// Cut off any part of the webs that protrude from the top of the pillar
		s0 = Intersect3D(s0, Cylinder3D(k.PillarHeight, k.WebDiameter, 0))
	}
	// Add the pillar hole/stub
	if k.HoleDepth >= 0.0 {
		// hole
		s0 = Difference3D(s0, pillarHole(k))
	} else {
		// support stub
		s0 = Union3D(s0, pillarHole(k))
	}
	return s0
}

// Standoffs3D returns multiple board standoffs at various positions.
func Standoffs3D(k *StandoffParms, positions V3Set) SDF3 {
	if len(positions) == 0 {
		return nil
	}
	s0 := Standoff3D(k)
	if s0 == nil {
		return nil
	}
	s := make([]SDF3, len(positions))
	for i, p := range positions {
		s[i] = Transform3D(s0, Translate3d(p))
	}
	return Union3D(s...)
}

//-----------------------------------------------------------------------------
// truncated rectangular pyramid (with rounded edges)

type TruncRectPyramidParms struct {
	Size        V3      // size of truncated pyramid
	Draft       float64 // draft of pyramid (radians)
	BaseRadius  float64 // base corner radius
	RoundRadius float64 // edge rounding radius
}

// TruncRectPyramid3D returns a truncated rectangular pyramid with rounded edges.
func TruncRectPyramid3D(k *TruncRectPyramidParms) SDF3 {
	h := k.Size.Z
	dr := h * math.Tan(k.Draft)
	rb := k.BaseRadius + dr
	rt := Max(k.BaseRadius-dr, 0)
	round := Min(0.5*rt, k.RoundRadius)
	s := Cone3D(2.0*h, rb, rt, round)
	wx := Max(k.Size.X-2.0*k.BaseRadius, 0)
	wy := Max(k.Size.Y-2.0*k.BaseRadius, 0)
	s = Elongate3D(s, V3{wx, wy, 0})
	s = Cut3D(s, V3{0, 0, 0}, V3{0, 0, 1})
	return s
}

//-----------------------------------------------------------------------------
