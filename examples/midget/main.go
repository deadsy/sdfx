//-----------------------------------------------------------------------------
/*

Midget Motor Casting Patterns
Popular Mechanics 1936

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func cylinderCoreBox() sdf.SDF3 {
	return nil
}

func cylinderFlange() sdf.SDF3 {

	const draft = 3.0

	const x = 2.0
	const y = 0.25 * 2.0 // base flange thickness * 2
	const z = x * 0.5

	const round = 0.125

	k := sdf.TruncRectPyramidParms{
		Size:        sdf.V3{x, y, z},
		BaseAngle:   sdf.DtoR(90 - draft),
		BaseRadius:  round,
		RoundRadius: round * 1.5,
	}
	base0 := sdf.TruncRectPyramid3D(&k)
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

	const flangeOfs = 3.0 / 16.0

	const l0 = flangeOfs + 0.25 + 0.75 // bottom to finned area
	const l1 = 1.75                    // finned area length
	const l2 = 4.0 + (7.0 / 16.0)      // core length

	const r0 = 1.0 * 0.5 // inner
	const r1 = 1.5 * 0.5 // non finned area
	const r2 = 2.0 * 0.5 // finned area

	// cylinder body
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(r1, draft*r1).Rel().Smooth(smooth1, smoothN)
	p.Add(0, l0).Rel().Smooth(smooth0, smoothN)
	p.Add(r2-r1, draft*(r2-r1)).Rel().Smooth(smooth0, smoothN)
	p.Add(0, l1).Rel().Smooth(smooth1, smoothN)
	p.Add(-r2, draft*r2).Rel()
	body := sdf.Revolve3D(sdf.Polygon2D(p.Vertices()))

	// cylinder base flange
	flange := cylinderFlange()
	flange = sdf.Transform3D(flange, sdf.Translate3d(sdf.V3{0, 0, flangeOfs}))

	// add it to the body pattern
	body = sdf.Union3D(body, flange)

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

func main() {
	scale := shrink * sdf.MillimetresPerInch
	sdf.RenderSTL(sdf.ScaleUniform3D(cylinderPattern(true, true), scale), 330, "cylinder_pattern.stl")
	//sdf.RenderSTL(sdf.ScaleUniform3D(cylinderCoreBox(), shrink), 330, "cylinder_corebox.stl")
}

//-----------------------------------------------------------------------------
