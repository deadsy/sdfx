//-----------------------------------------------------------------------------
/*

Monster Square Casting Pattern

Inspired by: https://fireballtool.com/monster-square/

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

const scale = sdf.MillimetresPerInch

const draft = 2.0

const pinRadius = 0.5 * 4.4 * sdf.InchesPerMillimetre

//-----------------------------------------------------------------------------

type msParms struct {
	name          string  // name of model
	size          float64 // length of square side
	width         float64 // width of square
	wallThickness float64 // outside wall thickness
	webThickness  float64 // web thickness
	holeRadius    float64 // radius of internal web hole
	holeOffset    float64 // offset of internal web hole
	allowance     float64 // machining allowance
	pinRadius     float64 // alignment pin radius
	nose          float64 // machined nose length on 45 degree corners
}

//-----------------------------------------------------------------------------

// envelope for the outside machined/cast surfaces
func envelope(k *msParms, machined bool) (sdf.SDF3, error) {
	c := k.nose
	l := k.size - c
	s0, err := sdf.Polygon2D([]v2.Vec{{0, 0}, {l, 0}, {l, c}, {c, l}, {0, l}})
	if err != nil {
		return nil, err
	}

	// machined
	if machined {
		return sdf.Extrude3D(s0, k.width), nil
	}

	// cast
	s1 := sdf.Extrude3D(s0, k.width+2.0*k.allowance)
	if err != nil {
		return nil, err
	}
	s2, err := walls(k)
	if err != nil {
		return nil, err
	}
	return sdf.Union3D(s1, s2), nil
}

//-----------------------------------------------------------------------------

// wall returns an outside wall with casting draft of length l
func wall(k *msParms, l float64) (sdf.SDF3, error) {
	trp := &obj.TruncRectPyramidParms{
		Size:        v3.Vec{l, k.wallThickness + k.allowance, (k.width * 0.5) + k.allowance},
		BaseAngle:   sdf.DtoR(90.0 - draft),
		BaseRadius:  (k.wallThickness + k.allowance) * 0.5,
		RoundRadius: k.wallThickness * 0.25,
	}
	s, err := obj.TruncRectPyramid3D(trp)
	if err != nil {
		return nil, err
	}
	ofs := (k.wallThickness - k.allowance) * 0.5
	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, ofs, 0}))
	return s, nil
}

// walls returns the 3 walls for a 45 degreee right angle triangle
func walls(k *msParms) (sdf.SDF3, error) {

	k0 := math.Sqrt(2.0)
	k1 := 1 + k0
	k2 := 2 + k0
	r := 0.5 * (k.wallThickness + k.allowance)

	// build the x-wall
	l0 := k.size + (k2 * k.allowance) - (k0 * r)
	w, err := wall(k, l0)
	if err != nil {
		return nil, err
	}
	ofs := 0.5*l0 - k.allowance
	w0 := sdf.Transform3D(w, sdf.Translate3d(v3.Vec{ofs, 0, 0}))

	// build the y-wall
	w1 := sdf.Transform3D(w0, sdf.MirrorXeqY())

	// build the 45-wall
	l1 := (k0 * k.size) + (2.0 * k1 * k.allowance) - (2.0 * k0 * r)
	w, err = wall(k, l1)
	if err != nil {
		return nil, err
	}
	ofs = 0.5 * k.size
	w2 := sdf.Transform3D(w, sdf.RotateZ(sdf.DtoR(135)))
	w2 = sdf.Transform3D(w2, sdf.Translate3d(v3.Vec{ofs, ofs, 0}))

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
	return sdf.Cut2D(s, v2.Vec{0, 0}, v2.Vec{0, 1}), nil
}

// web2d returns the 2d internal web.
func web2d(k *msParms) (sdf.SDF2, error) {
	ofs := k.wallThickness * 0.9
	l := k.size - ofs*(2.0+math.Sqrt(2.0))

	s, err := sdf.Polygon2D([]v2.Vec{{ofs, ofs}, {ofs + l, ofs}, {ofs, ofs + l}})
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

	h0 := sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{0, k0}))
	h1 := sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{0, k1}))

	hole = sdf.Transform2D(hole, sdf.Rotate2d(sdf.DtoR(90.0)))
	h2 := sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{k0, 0}))
	h3 := sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{k1, 0}))

	hole = sdf.Transform2D(hole, sdf.Rotate2d(sdf.DtoR(135.0)))
	h4 := sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{k2, k2}))

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
	r := 2.0 * k.wallThickness
	trp := &obj.TruncRectPyramidParms{
		Size:        v3.Vec{2.0 * r, 2.0 * r, (k.width * 0.5) + k.allowance},
		BaseAngle:   sdf.DtoR(90.0 - 3.0*draft),
		BaseRadius:  r,
		RoundRadius: k.wallThickness * 0.25,
	}
	s, err := obj.TruncRectPyramid3D(trp)
	if err != nil {
		return nil, err
	}
	ofs := 0.8 * r
	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{ofs, ofs, 0}))
	return s, nil
}

func corner45(k *msParms) (sdf.SDF3, error) {
	r := 2.3 * k.wallThickness
	trp := &obj.TruncRectPyramidParms{
		Size:        v3.Vec{2.0 * r, 2.0 * r, (k.width * 0.5) + k.allowance},
		BaseAngle:   sdf.DtoR(90.0 - 3.0*draft),
		BaseRadius:  r,
		RoundRadius: k.wallThickness * 0.25,
	}
	s, err := obj.TruncRectPyramid3D(trp)
	if err != nil {
		return nil, err
	}
	dy := 0.7 * r
	dx := dy * (1.0 + math.Sqrt(2.0))
	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{k.size - dx, dy, 0}))
	return s, nil
}

func corners(k *msParms) (sdf.SDF3, error) {

	// build the 90 degree corner
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

func pin(k *msParms) (sdf.SDF3, error) {
	return sdf.Cylinder3D(k.width*0.8, k.pinRadius, 0)
}

// pins returns split-casting alignment pins
func pins(k *msParms) (sdf.SDF3, error) {

	if k.pinRadius == 0 {
		return nil, nil
	}

	// build the pin at the 90 degree corner
	p0, err := pin(k)
	if err != nil {
		return nil, err
	}
	ofs := 1.5 * k.wallThickness
	p0 = sdf.Transform3D(p0, sdf.Translate3d(v3.Vec{ofs, ofs, 0}))

	// build the pins at the 45 degree corners
	p1, err := pin(k)
	if err != nil {
		return nil, err
	}
	dy := 1.5 * k.wallThickness
	dx := dy * (1.0 + math.Sqrt(2.0))
	p1 = sdf.Transform3D(p1, sdf.Translate3d(v3.Vec{k.size - dx, dy, 0}))
	p2 := sdf.Transform3D(p1, sdf.MirrorXeqY())

	return sdf.Union3D(p0, p1, p2), nil
}

//-----------------------------------------------------------------------------

func mSquare(k *msParms, machined bool) error {

	web, err := web(k)
	if err != nil {
		return err
	}

	walls, err := walls(k)
	if err != nil {
		return err
	}

	corners, err := corners(k)
	if err != nil {
		return err
	}

	s := sdf.Union3D(web, walls, corners)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(k.webThickness))

	// remove the pin cavities
	pins, err := pins(k)
	if err != nil {
		return err
	}
	s = sdf.Difference3D(s, pins)

	// cleanup with the outside envelope
	env, err := envelope(k, machined)
	if err != nil {
		return err
	}
	s = sdf.Intersect3D(s, env)

	s = sdf.ScaleUniform3D(s, shrink*scale)
	render.ToSTL(s, fmt.Sprintf("%s.stl", k.name), render.NewMarchingCubesOctree(300))

	sUpper := sdf.Cut3D(s, v3.Vec{0, 0, 0}, v3.Vec{0, 0, 1})
	render.ToSTL(sUpper, fmt.Sprintf("%s_upper.stl", k.name), render.NewMarchingCubesOctree(300))

	sLower := sdf.Cut3D(s, v3.Vec{0, 0, 0}, v3.Vec{0, 0, -1})
	render.ToSTL(sLower, fmt.Sprintf("%s_lower.stl", k.name), render.NewMarchingCubesOctree(300))

	return nil
}

//-----------------------------------------------------------------------------

func main() {

	/*

		k := &msParms{
			name:          "ms6",
			size:          6.0,
			width:         2.0,
			wallThickness: 0.25,
			webThickness:  0.25,
			holeRadius:    0.5,
			holeOffset:    0.75,
			allowance:     0.0625,
			pinRadius:     pinRadius,
			nose:          0.25,
		}

		k := &msParms{
			name:          "ms8",
			size:          8.0,
			width:         3.0,
			wallThickness: 0.375,
			webThickness:  0.25,
			holeRadius:    0.5,
			holeOffset:    1.0,
			allowance:     0.0625,
			pinRadius:     pinRadius,
			nose:          0.375,
		}

		k := &msParms{
			name:          "ms12",
			size:          12.0,
			width:         3.0,
			wallThickness: 0.375,
			webThickness:  0.25,
			holeRadius:    0.75,
			holeOffset:    1.5,
			allowance:     0.0625,
			pinRadius:     pinRadius,
			nose:          0.5,
		}

	*/

	k := &msParms{
		name:          "ms6",
		size:          6.0,
		width:         2.0,
		wallThickness: 0.25,
		webThickness:  0.25,
		holeRadius:    0.5,
		holeOffset:    0.75,
		allowance:     0.0625,
		pinRadius:     pinRadius,
		nose:          0.25,
	}

	err := mSquare(k, false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

}

//-----------------------------------------------------------------------------
