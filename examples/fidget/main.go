//-----------------------------------------------------------------------------
/*

Fidget Spinners

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// 608 bearing
const bearingOuterOD = 22.0  // outer diameter of outer race
const bearingOuterID = 19.2  // inner diameter of outer race
const bearingInnerOD = 12.1  // outer diameter of inner race
const bearingInnerID = 8.0   // inner diameter of inner race
const bearingThickness = 7.0 // bearing thickness

// Adjust clearance to give good interference fits for the bearings and spin caps.
const clearance = 0.0

//-----------------------------------------------------------------------------

// ball bearing counterweights
const bbLargeD = (1.0 / 2.0) * sdf.MillimetresPerInch
const bbSmallD = (5.0 / 16.0) * sdf.MillimetresPerInch

//-----------------------------------------------------------------------------

// Return an N petal bezier flower.
func flower(n int, r0, r1, r2 float64) (sdf.SDF2, error) {

	theta := sdf.Tau / float64(n)
	b := sdf.NewBezier()

	p0 := sdf.V2{r1, 0}.Add(sdf.PolarToXY(r0, sdf.DtoR(-135)))
	p1 := sdf.V2{r1, 0}.Add(sdf.PolarToXY(r0, sdf.DtoR(-45)))
	p2 := sdf.V2{r1, 0}.Add(sdf.PolarToXY(r0, sdf.DtoR(45)))
	p3 := sdf.V2{r1, 0}.Add(sdf.PolarToXY(r0, sdf.DtoR(135)))
	p4 := sdf.PolarToXY(r2, theta/2)

	m := sdf.Rotate(theta)

	for i := 0; i < n; i++ {
		ofs := float64(i) * theta

		b.AddV2(p0).Handle(ofs+sdf.DtoR(-45), r0/2, r0/2)
		b.AddV2(p1).Handle(ofs+sdf.DtoR(45), r0/2, r0/2)
		b.AddV2(p2).Handle(ofs+sdf.DtoR(135), r0/2, r0/2)
		b.AddV2(p3).Handle(ofs+sdf.DtoR(225), r0/2, r0/2)
		b.AddV2(p4).Handle(ofs+theta/2+sdf.DtoR(90), r2/1.5, r2/1.5)

		p0 = m.MulPosition(p0)
		p1 = m.MulPosition(p1)
		p2 = m.MulPosition(p2)
		p3 = m.MulPosition(p3)
		p4 = m.MulPosition(p4)
	}

	b.Close()
	p, err := b.Polygon()
	if err != nil {
		return nil, err
	}

	return sdf.Polygon2D(p.Vertices())
}

func body1() (sdf.SDF3, error) {

	n := 3
	t := bearingThickness
	r := bearingOuterOD / 2

	r0 := r + 4.0
	r1 := 45.0 - r0
	r2 := r + 4.0

	// body
	f, err := flower(n, r0, r1, r2)
	if err != nil {
		log.Fatal(err)
	}
	s1, err := sdf.ExtrudeRounded3D(f, t, t/4.0)
	if err != nil {
		return nil, err
	}

	// periphery holes
	s2, err := obj.BoltCircle3D(t, r+clearance, r1, n)
	if err != nil {
		return nil, err
	}
	// center hole
	s3, err := sdf.Cylinder3D(t, r+clearance, 0)
	if err != nil {
		return nil, err
	}
	return sdf.Difference3D(s1, sdf.Union3D(s2, s3)), nil
}

//-----------------------------------------------------------------------------

func body2() (sdf.SDF3, error) {
	t := bearingThickness
	r := bearingOuterOD / 2
	r0 := r + 4.0

	// build the arm
	p := sdf.NewPolygon()
	p.Add(r, -t/2)
	p.Add(r0, -t/2)
	p.Add(r0, t/2)
	p.Add(r, t/2)
	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	theta := sdf.DtoR(270)
	arm, err := sdf.RevolveTheta3D(s, theta)
	if err != nil {
		return nil, err
	}
	arm = sdf.Transform3D(arm, sdf.Translate3d(sdf.V3{-1.5 * r0, 0, 0}))

	// create 6 arms
	arms := sdf.RotateUnion3D(arm, 6, sdf.RotateZ(sdf.DtoR(60)))

	// add the center
	body, err := sdf.Cylinder3D(t, r0, 0)
	if err != nil {
		return nil, err
	}
	body = sdf.Union3D(body, arms)

	// remove the center hole
	hole, err := sdf.Cylinder3D(t, r, 0)
	if err != nil {
		return nil, err
	}
	return sdf.Difference3D(body, hole), nil
}

//-----------------------------------------------------------------------------

// Basic spin cap with variable pin size.
func spincap(
	pinR float64, // pin radius
	pinL float64, // pin length
) (sdf.SDF3, error) {

	t := 3.0  // thickness of the spin cap
	st := 1.0 // spacer thickness

	r0 := bearingOuterOD / 2
	r1 := bearingInnerOD / 2

	p := sdf.NewPolygon()
	p.Add(0, -t-st)
	p.Add(r0, -t-st).Smooth(t/1.5, 6)
	p.Add(r0, -st)
	p.Add(r1, -st)
	p.Add(r1, 0)
	p.Add(pinR, 0)
	p.Add(pinR, pinL)
	p.Add(0, pinL)

	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}

	return sdf.Revolve3D(s)
}

//-----------------------------------------------------------------------------

// Push to fit spincap for single spinner.
func spincapSingle() (sdf.SDF3, error) {
	gap := 1.0
	r := (bearingInnerID / 2) - clearance
	l := (bearingThickness - gap) / 2
	return spincap(r, l)
}

//-----------------------------------------------------------------------------

// Threaded spincap for double spinners.
func spincapDouble(male bool) (sdf.SDF3, error) {
	r := (bearingInnerID / 2) - clearance
	threadR := r * 0.8
	threadPitch := 1.0
	threadTolerance := 0.25
	l := bearingThickness

	if male {
		// Add an external screw thread.
		t, err := sdf.ISOThread(threadR-threadTolerance, threadPitch, true)
		if err != nil {
			return nil, err
		}
		screw, err := sdf.Screw3D(t, bearingThickness, threadPitch, 1, 0)
		if err != nil {
			return nil, err
		}
		screw, err = obj.ChamferedCylinder(screw, 0, 0.5)
		if err != nil {
			return nil, err
		}
		screw = sdf.Transform3D(screw, sdf.Translate3d(sdf.V3{0, 0, 1.5 * l}))
		sc, err := spincap(r, l+0.5)
		if err != nil {
			return nil, err
		}
		return sdf.Union3D(sc, screw), nil
	}
	// Add an internal screw thread.
	t, err := sdf.ISOThread(threadR, threadPitch, false)
	if err != nil {
		return nil, err
	}
	screw, err := sdf.Screw3D(t, bearingThickness, threadPitch, 1, 0)
	if err != nil {
		return nil, err
	}
	screw = sdf.Transform3D(screw, sdf.Translate3d(sdf.V3{0, 0, l * 0.5}))
	sc, err := spincap(r, l-0.5)
	if err != nil {
		return nil, err
	}
	return sdf.Difference3D(sc, screw), nil
}

// Inner washer for double spinner.
func spincapWasher() (sdf.SDF3, error) {
	k := obj.WasherParms{
		Thickness:   1.0,
		InnerRadius: (bearingInnerID / 2) * 1.05,
		OuterRadius: (bearingOuterOD + bearingInnerID) / 4,
	}
	s, err := obj.Washer3D(&k)
	if err != nil {
		return nil, err
	}
	return s, nil
}

//-----------------------------------------------------------------------------

func main() {
	body1, err := body1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(body1, 300, "body1.stl")

	body2, err := body2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(body2, 300, "body2.stl")

	scs, err := spincapSingle()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(scs, 150, "cap_single.stl")

	scdm, err := spincapDouble(true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(scdm, 150, "cap_double_male.stl")

	scdf, err := spincapDouble(false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(scdf, 150, "cap_double_female.stl")

	scw, err := spincapWasher()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(scw, 150, "washer.stl")
}

//-----------------------------------------------------------------------------
