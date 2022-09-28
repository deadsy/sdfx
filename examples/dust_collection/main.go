//-----------------------------------------------------------------------------
/*

Dust Collection Adapters

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// dust deputy tapered pipe
const ddOuterDiameter = 51.0
const ddLength = 39.0

var ddTaper = sdf.DtoR(2.0)

// vacuum hose 2.5" male fitting
const vhOuterDiameter = 58.0
const vhClearance = 0.6

var vhTaper = sdf.DtoR(0.4)

const wallThickness = 4.0

//-----------------------------------------------------------------------------

// dust deputy (female), 2.5" vacuum (female)
func dustDeputyToVacuumFF() (sdf.SDF3, error) {

	const t = wallThickness
	const transitionLength = 15.0
	const vhLength = 30.0

	r0 := ddOuterDiameter * 0.5
	r1 := r0 - ddLength*math.Tan(ddTaper)
	r3 := (vhOuterDiameter + vhClearance) * 0.5
	r2 := r3 - (vhLength * math.Tan(vhTaper))

	h0 := 0.0
	h1 := h0 + ddLength
	h2 := h1 + transitionLength
	h3 := h2 + vhLength

	p := sdf.NewPolygon()
	p.Add(r0+t, h0)
	p.Add(r1+t, h1).Smooth(t, 4)
	p.Add(r2+t, h2).Smooth(t, 4)
	p.Add(r3+t, h3)
	p.Add(r3, h3)
	p.Add(r2, h2).Smooth(t, 4)
	p.Add(r1, h1).Smooth(t, 4)
	p.Add(r0, h0)

	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}

	return sdf.Revolve3D(s)
}

//-----------------------------------------------------------------------------

// 2.5" vacuum (male) to pipe (male)
func vacuumToPipeMM(name string) (sdf.SDF3, error) {

	k, err := obj.PipeLookup(name, "mm")
	if err != nil {
		return nil, err
	}

	t := wallThickness
	transitionLength := 15.0

	r0 := k.Outer
	r1 := vhOuterDiameter * 0.5

	h0 := 0.0
	h1 := h0 + 35.0
	h2 := h1 + transitionLength
	h3 := h2 + 20.0

	p := sdf.NewPolygon()
	p.Add(r0, h0)
	p.Add(r0, h1).Smooth(t, 4)
	p.Add(r1, h2).Smooth(t, 4)
	p.Add(r1, h3)
	p.Add(r1-t, h3)
	p.Add(r1-t, h2).Smooth(t, 4)
	p.Add(r0-t, h1).Smooth(t, 4)
	p.Add(r0-t, h0)

	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}

	return sdf.Revolve3D(s)
}

//-----------------------------------------------------------------------------

// dust deputy (female) to pipe (male)
func dustDeputyToPipeFM(name string) (sdf.SDF3, error) {

	k, err := obj.PipeLookup(name, "mm")
	if err != nil {
		return nil, err
	}

	t := wallThickness
	transitionLength := 15.0

	r0 := k.Outer
	r2 := (ddOuterDiameter * 0.5) + t
	r1 := r2 - ddLength*math.Tan(ddTaper)

	h0 := 0.0
	h1 := h0 + 35.0
	h2 := h1 + transitionLength
	h3 := h2 + ddLength

	p := sdf.NewPolygon()
	p.Add(r0, h0)
	p.Add(r0, h1).Smooth(t, 4)
	p.Add(r1, h2).Smooth(t, 4)
	p.Add(r2, h3)
	p.Add(r2-t, h3)
	p.Add(r1-t, h2).Smooth(t, 4)
	p.Add(r0-t, h1).Smooth(t, 4)
	p.Add(r0-t, h0)

	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}

	return sdf.Revolve3D(s)
}

//-----------------------------------------------------------------------------

func main() {
	s, err := dustDeputyToVacuumFF()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "fdd_fvh25.stl", render.NewMarchingCubesOctree(150))

	s, err = vacuumToPipeMM("sch40:2")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "mvh25_mpvc.stl", render.NewMarchingCubesOctree(150))

	s, err = dustDeputyToPipeFM("sch40:2")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "fdd_mpvc.stl", render.NewMarchingCubesOctree(150))
}

//-----------------------------------------------------------------------------
