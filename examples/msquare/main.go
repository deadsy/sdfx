//-----------------------------------------------------------------------------
/*

Monster Square Casting Pattern

Inspired by: https://fireballtool.com/monster-square/

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

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

const scale = sdf.MillimetresPerInch

const draft = 2.0

//-----------------------------------------------------------------------------

type msParms struct {
	size          float64 // length of square side
	width         float64 // width of square
	wallThickness float64 // outside wall thickness
	webThickness  float64 // web thickness
	holeRadius    float64 // radius of internal web hole
	holeOffset    float64 // offset of internal web hole
	allowance     float64 // machining allowance
}

//-----------------------------------------------------------------------------

// envelope for the outside machined surfaces
func envelope(k *msParms, machined bool) (sdf.SDF3, error) {

	if machined {

		s, err := sdf.Polygon2D([]sdf.V2{{0, 0}, {k.size, 0}, {0, k.size}})
		if err != nil {
			return nil, err
		}
		return sdf.Extrude3D(s, k.width), nil

	}

	k0 := -k.allowance
	k1 := k.size + (1.0*math.Sqrt(2.0))*k.allowance
	k2 := k.width + 2.0*k.allowance

	s, err := sdf.Polygon2D([]sdf.V2{{k0, k0}, {k0 + k1, 0}, {0, k0 + k1}})
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, k2), nil

}

//-----------------------------------------------------------------------------

// wall returns a wall for the msquare
func wall(k *msParms, l float64) (sdf.SDF3, error) {
	trp := &obj.TruncRectPyramidParms{
		Size:        sdf.V3{l, k.wallThickness + k.allowance, (k.width * 0.5) + k.allowance},
		BaseAngle:   sdf.DtoR(90.0 - draft),
		BaseRadius:  (k.wallThickness + k.allowance) * 0.5,
		RoundRadius: k.wallThickness * 0.25,
	}
	s, err := obj.TruncRectPyramid3D(trp)
	if err != nil {
		return nil, err
	}
	ofs := (k.wallThickness - k.allowance) * 0.5
	s = sdf.Transform3D(s, sdf.Translate3d(sdf.V3{0, ofs, 0}))
	return s, nil
}

func walls(k *msParms) (sdf.SDF3, error) {

	ofs := 0.5 * k.size

	w, err := wall(k, k.size+2.0*k.allowance)
	if err != nil {
		return nil, err
	}

	// build the x-wall
	w0 := sdf.Transform3D(w, sdf.Translate3d(sdf.V3{ofs, 0, 0}))

	// build the y-wall
	w1 := sdf.Transform3D(w, sdf.RotateZ(sdf.DtoR(-90)))
	w1 = sdf.Transform3D(w1, sdf.Translate3d(sdf.V3{0, ofs, 0}))

	// build the 45-wall
	w, err = wall(k, math.Sqrt(2.0)*k.size+2.0*k.allowance)
	if err != nil {
		return nil, err
	}
	w2 := sdf.Transform3D(w, sdf.RotateZ(sdf.DtoR(135)))
	w2 = sdf.Transform3D(w2, sdf.Translate3d(sdf.V3{ofs, ofs, 0}))

	// build the flipped walls
	w0f := sdf.Transform3D(w0, sdf.MirrorXY())
	w1f := sdf.Transform3D(w1, sdf.MirrorXY())
	w2f := sdf.Transform3D(w2, sdf.MirrorXY())

	return sdf.Union3D(w0, w1, w2, w0f, w1f, w2f), nil
}

//-----------------------------------------------------------------------------

// webHole returns the clamp hole within the web.
func webHole(k *msParms) (sdf.SDF2, error) {
	r := k.holeRadius + 0.5*k.webThickness
	l := 2.0*k.holeOffset + k.webThickness
	s := sdf.Line2D(l, r)
	return sdf.Cut2D(s, sdf.V2{0, 0}, sdf.V2{0, 1}), nil
}

// web2d returns the 2d internal web.
func web2d(k *msParms) (sdf.SDF2, error) {
	ofs := k.wallThickness * 0.9
	l := k.size - ofs*(2.0+math.Sqrt(2.0))

	s, err := sdf.Polygon2D([]sdf.V2{{ofs, ofs}, {ofs + l, ofs}, {ofs, ofs + l}})
	if err != nil {
		return nil, err
	}

	if k.holeRadius == 0 {
		return s, nil
	}

	hole, err := webHole(k)
	if err != nil {
		return nil, err
	}

	k0 := k.size * 0.3
	k1 := k.size * 0.6
	k2 := k.size * 0.5

	h0 := sdf.Transform2D(hole, sdf.Translate2d(sdf.V2{0, k0}))
	h1 := sdf.Transform2D(hole, sdf.Translate2d(sdf.V2{0, k1}))

	hole = sdf.Transform2D(hole, sdf.Rotate2d(sdf.DtoR(90.0)))
	h2 := sdf.Transform2D(hole, sdf.Translate2d(sdf.V2{k0, 0}))
	h3 := sdf.Transform2D(hole, sdf.Translate2d(sdf.V2{k1, 0}))

	hole = sdf.Transform2D(hole, sdf.Rotate2d(sdf.DtoR(135.0)))
	h4 := sdf.Transform2D(hole, sdf.Translate2d(sdf.V2{k2, k2}))

	return sdf.Difference2D(s, sdf.Union2D(h0, h1, h2, h3, h4)), nil
}

// web returns the internal web.
func web(k *msParms) (sdf.SDF3, error) {
	s0, err := web2d(k)
	if err != nil {
		return nil, err
	}
	return sdf.ExtrudeRounded3D(s0, k.webThickness, 0.5*k.webThickness)
}

//-----------------------------------------------------------------------------

func corner90(k *msParms) (sdf.SDF3, error) {
	r := 1.5 * k.wallThickness
	trp := &obj.TruncRectPyramidParms{
		Size:        sdf.V3{2.0 * r, 2.0 * r, (k.width * 0.5) + k.allowance},
		BaseAngle:   sdf.DtoR(90.0 - 2.0*draft),
		BaseRadius:  r,
		RoundRadius: k.wallThickness * 0.25,
	}
	s, err := obj.TruncRectPyramid3D(trp)
	if err != nil {
		return nil, err
	}
	ofs := 1.2 * r
	s = sdf.Transform3D(s, sdf.Translate3d(sdf.V3{ofs, ofs, 0}))
	return s, nil
}

func corner45(k *msParms) (sdf.SDF3, error) {
	r := 1.8 * k.wallThickness
	trp := &obj.TruncRectPyramidParms{
		Size:        sdf.V3{2.0 * r, 2.0 * r, (k.width * 0.5) + k.allowance},
		BaseAngle:   sdf.DtoR(90.0 - 2.0*draft),
		BaseRadius:  r,
		RoundRadius: k.wallThickness * 0.25,
	}
	s, err := obj.TruncRectPyramid3D(trp)
	if err != nil {
		return nil, err
	}

	dy := 0.8 * r
	dx := dy * (1.0 + math.Sqrt(2.0))

	s = sdf.Transform3D(s, sdf.Translate3d(sdf.V3{k.size - dx, dy, 0}))
	return s, nil
}

func corners(k *msParms) (sdf.SDF3, error) {

	// build the 90 degreee corner
	c0, err := corner90(k)
	if err != nil {
		return nil, err
	}

	// build the 45 degree corners
	c1, err := corner45(k)
	if err != nil {
		return nil, err
	}
	c2 := sdf.Transform3D(c1, sdf.MirrorXeqY())

	// build the flipped corners
	c0f := sdf.Transform3D(c0, sdf.MirrorXY())
	c1f := sdf.Transform3D(c1, sdf.MirrorXY())
	c2f := sdf.Transform3D(c2, sdf.MirrorXY())

	return sdf.Union3D(c0, c0f, c1, c1f, c2, c2f), nil
}

//-----------------------------------------------------------------------------

func mSquare(k *msParms, machined bool) ([3]sdf.SDF3, error) {

	var bad [3]sdf.SDF3

	web, err := web(k)
	if err != nil {
		return bad, err
	}

	walls, err := walls(k)
	if err != nil {
		return bad, err
	}

	corners, err := corners(k)
	if err != nil {
		return bad, err
	}

	s := sdf.Union3D(web, walls, corners)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(k.webThickness))

	env, err := envelope(k, machined)
	if err != nil {
		return bad, err
	}
	s = sdf.Intersect3D(s, env)

	sUpper := sdf.Cut3D(s, sdf.V3{0, 0, 0}, sdf.V3{0, 0, 1})
	sLower := sdf.Cut3D(s, sdf.V3{0, 0, 0}, sdf.V3{0, 0, -1})

	return [3]sdf.SDF3{s, sUpper, sLower}, nil
}

//-----------------------------------------------------------------------------

func main() {

	k := &msParms{
		size:          8.0,
		width:         3.0,
		wallThickness: 0.375,
		webThickness:  0.25,
		holeRadius:    0.5,
		holeOffset:    1.0,
		allowance:     0.0625,
	}

	ss, err := mSquare(k, true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s := sdf.ScaleUniform3D(ss[0], shrink*scale)
	render.RenderSTLSlow(s, 300, "ms8.stl")

	s = sdf.ScaleUniform3D(ss[1], shrink*scale)
	render.RenderSTLSlow(s, 300, "ms8_upper.stl")

	s = sdf.ScaleUniform3D(ss[2], shrink*scale)
	render.RenderSTLSlow(s, 300, "ms8_lower.stl")

}

//-----------------------------------------------------------------------------
