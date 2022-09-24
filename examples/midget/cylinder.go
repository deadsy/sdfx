//-----------------------------------------------------------------------------
/*

Cylinder Pattern and Core Box

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

const cylinderBaseOffset = 3.0 / 16.0
const cylinderBaseThickness = 0.25
const cylinderWaistLength = 0.75
const cylinderBodyLength = 1.75
const cylinderCoreLength = 4.0 + (7.0 / 16.0)

const cylinderInnerRadius = 1.0 * 0.5
const cylinderWaistRadius = 1.5 * 0.5
const cylinderBodyRadius = 2.0 * 0.5

func cylinderCoreBox() (sdf.SDF3, error) {
	return nil, nil
}

func cylinderBase() (sdf.SDF3, error) {

	const draft = 3.0

	const x = cylinderBodyRadius * 2.0
	const y = cylinderBaseThickness * 2.0
	const z = cylinderBodyRadius

	const round = 0.125

	k := obj.TruncRectPyramidParms{
		Size:        v3.Vec{x, y, z},
		BaseAngle:   sdf.DtoR(90 - draft),
		BaseRadius:  round,
		RoundRadius: round * 1.5,
	}
	base0, err := obj.TruncRectPyramid3D(&k)
	if err != nil {
		return nil, err
	}
	base1 := sdf.Transform3D(base0, sdf.MirrorXY())
	base := sdf.Union3D(base0, base1)
	base = sdf.Cut3D(base, v3.Vec{0, 0, 0}, v3.Vec{0, 1, 0})
	return sdf.Transform3D(base, sdf.RotateX(sdf.DtoR(90))), nil
}

func cylinderPattern(core, split bool) (sdf.SDF3, error) {

	draft := math.Tan(sdf.DtoR(3.0))
	const smooth0 = 0.125
	const smooth1 = smooth0 * 0.5
	const smoothN = 5

	const l0 = cylinderBaseOffset + cylinderBaseThickness + cylinderWaistLength
	const l1 = cylinderBodyLength
	const l2 = cylinderCoreLength

	const r0 = cylinderInnerRadius
	const r1 = cylinderWaistRadius
	const r2 = cylinderBodyRadius

	// cylinder body
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(r1, draft*r1).Rel().Smooth(smooth1, smoothN)
	p.Add(0, l0).Rel().Smooth(smooth0, smoothN)
	p.Add(r2-r1, draft*(r2-r1)).Rel().Smooth(smooth0, smoothN)
	p.Add(0, l1).Rel().Smooth(smooth1, smoothN)
	p.Add(-r2, draft*r2).Rel()
	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	body, err := sdf.Revolve3D(s)
	if err != nil {
		return nil, err
	}
	// cylinder base
	base, err := cylinderBase()
	if err != nil {
		return nil, err
	}
	base = sdf.Transform3D(base, sdf.Translate3d(v3.Vec{0, 0, cylinderBaseOffset}))

	// add the base to the body pattern
	body = sdf.Union3D(body, base)

	// core print
	p = sdf.NewPolygon()
	p.Add(0, -0.75)
	p.Add(r0, draft*r0).Rel().Smooth(smooth1, smoothN)
	p.Add(0, l2).Rel().Smooth(smooth1, smoothN)
	p.Add(-r0, draft*r0).Rel()
	s, err = sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	corePrint, err := sdf.Revolve3D(s)
	if err != nil {
		return nil, err
	}

	var cylinder sdf.SDF3
	if core {
		cylinder = sdf.Union3D(body, corePrint)
	} else {
		cylinder = sdf.Difference3D(body, corePrint)
	}

	if split {
		cylinder = sdf.Cut3D(cylinder, v3.Vec{0, 0, 0}, v3.Vec{0, 1, 0})
	}

	return cylinder, nil
}

//-----------------------------------------------------------------------------
