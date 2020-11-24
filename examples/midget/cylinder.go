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

func cylinderCoreBox() sdf.SDF3 {
	return nil
}

func cylinderBase() sdf.SDF3 {

	const draft = 3.0

	const x = cylinderBodyRadius * 2.0
	const y = cylinderBaseThickness * 2.0
	const z = cylinderBodyRadius

	const round = 0.125

	k := obj.TruncRectPyramidParms{
		Size:        sdf.V3{x, y, z},
		BaseAngle:   sdf.DtoR(90 - draft),
		BaseRadius:  round,
		RoundRadius: round * 1.5,
	}
	base0, _ := obj.TruncRectPyramid3D(&k)
	base1 := sdf.Transform3D(base0, sdf.MirrorXY())
	base := sdf.Union3D(base0, base1)
	base = sdf.Cut3D(base, sdf.V3{0, 0, 0}, sdf.V3{0, 1, 0})
	return sdf.Transform3D(base, sdf.RotateX(sdf.DtoR(90)))
}

func cylinderPattern(core, split bool) sdf.SDF3 {

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
	body := sdf.Revolve3D(sdf.Polygon2D(p.Vertices()))

	// cylinder base
	base := cylinderBase()
	base = sdf.Transform3D(base, sdf.Translate3d(sdf.V3{0, 0, cylinderBaseOffset}))

	// add the base to the body pattern
	body = sdf.Union3D(body, base)

	// core print
	p = sdf.NewPolygon()
	p.Add(0, -0.75)
	p.Add(r0, draft*r0).Rel().Smooth(smooth1, smoothN)
	p.Add(0, l2).Rel().Smooth(smooth1, smoothN)
	p.Add(-r0, draft*r0).Rel()
	corePrint := sdf.Revolve3D(sdf.Polygon2D(p.Vertices()))

	var cylinder sdf.SDF3
	if core {
		cylinder = sdf.Union3D(body, corePrint)
	} else {
		cylinder = sdf.Difference3D(body, corePrint)
	}

	if split {
		cylinder = sdf.Cut3D(cylinder, sdf.V3{0, 0, 0}, sdf.V3{0, 1, 0})
	}

	return cylinder
}

//-----------------------------------------------------------------------------
