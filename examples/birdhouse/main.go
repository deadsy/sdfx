//-----------------------------------------------------------------------------
/*

An A-Frame Birdhouse.

git@github.com:deadsy/sdfx.git
https://github.com/deadsy/sdfx/tree/master/examples/birdhouse

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

const width = 120.0
const height = 85.0
const thickness = 2.0
const hookHeight = 10.0
const holeFactor = 0.9 // control the hole size 0..1

//-----------------------------------------------------------------------------

// holeRadius returns the radius for a circle inscribed within the frame triangle.
func holeRadius() float64 {
	b := width * 0.5
	h := height
	a := math.Sqrt((b * b) + (h * h))
	return b * math.Sqrt((a-b)/(a+b))
}

// hook returns a hook used to suspend the birdhouse.
func hook() (sdf.SDF3, error) {
	k := obj.WasherParms{
		Thickness:   thickness,
		InnerRadius: hookHeight * 0.5,
		OuterRadius: hookHeight,
		Remove:      0.5,
	}
	s, err := obj.Washer3D(&k)
	if err != nil {
		return nil, err
	}
	m := sdf.RotateY(sdf.DtoR(90))
	m = sdf.Translate3d(v3.Vec{0, 0, height + thickness}).Mul(m)
	return sdf.Transform3D(s, m), nil
}

// frame returns a birdhouse A-frame.
func frame() (sdf.SDF3, error) {
	p := sdf.NewPolygon()
	p.Add(width/2, 0)
	p.Add(0, height)
	p.Add(-width/2, 0)
	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	sOuter := sdf.Offset2D(s, 2*thickness)
	sInner := sdf.Offset2D(s, thickness)
	f2d := sdf.Difference2D(sOuter, sInner)
	f3d := sdf.Extrude3D(f2d, width*1.1)
	return sdf.Transform3D(f3d, sdf.RotateX(sdf.DtoR(90))), nil
}

func hole() (sdf.SDF3, error) {
	r := holeRadius()
	s, err := sdf.Cylinder3D(2*width, r*holeFactor, 0)
	if err != nil {
		return nil, err
	}
	m := sdf.RotateX(sdf.DtoR(90))
	m = sdf.Translate3d(v3.Vec{0, 0, r}).Mul(m)
	return sdf.Transform3D(s, m), nil
}

// cross returns the union of s and a copy rotated 90 degrees about the z axis.
func cross(s sdf.SDF3) sdf.SDF3 {
	s1 := sdf.Transform3D(s, sdf.RotateZ(sdf.DtoR(90)))
	return sdf.Union3D(s, s1)
}

func birdhouse() (sdf.SDF3, error) {
	frame, err := frame()
	if err != nil {
		return nil, err
	}
	hole, err := hole()
	if err != nil {
		return nil, err
	}
	hook, err := hook()
	if err != nil {
		return nil, err
	}
	s := sdf.Difference3D(cross(frame), cross(hole))
	s = sdf.Union3D(s, hook)
	return s, nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := birdhouse()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "birdhouse.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
