//-----------------------------------------------------------------------------
/*

Drain Covers

This code implements a parametric drain cover. Draft angles are implemented so a
3D printed model can be used as a sand casting pattern.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// DrainCoverParms defines a grated drain pipe cover.
type DrainCoverParms struct {
	WallDiameter   float64 // outer diameter of wall
	WallHeight     float64 // height of wall
	WallThickness  float64 // thickness of wall
	WallDraft      float64 // draft angle of walls
	OuterWidth     float64 // extra width beyond the wall
	InnerWidth     float64 // width between inner wall and grate field
	CoverThickness float64 // thickness of the drain cover
	GrateNumber    int     // number of grate slots
	GrateWidth     float64 // multiple of inter-slot gap
	GrateDraft     float64 // draft angle of grate slots (radians)
	CrossBarWidth  float64 // multiple of InnerWidth
	CrossBarWeb    bool    // add a reinforcing crossbar web
}

//-----------------------------------------------------------------------------

func dcBody(k *DrainCoverParms) (sdf.SDF3, error) {

	// x drafts
	dx0 := 0.5 * k.CoverThickness * math.Tan(k.WallDraft)
	dx1 := 0.5 * k.WallHeight * math.Tan(k.WallDraft)

	// x radii
	r0 := 0.5 * k.WallDiameter
	r1 := r0 + k.OuterWidth
	r2 := r0 - k.WallThickness

	// y thicknesses
	t0 := k.CoverThickness
	t1 := t0 + k.WallHeight

	// build the 2d profile
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(r1+dx0, 0)
	p.Add(r1-dx0, t0).Smooth(0.25*k.CoverThickness, 4)
	p.Add(r0+dx1, t0).Smooth(0.25*k.WallThickness, 4)
	p.Add(r0-dx1, t1).Smooth(0.25*k.WallThickness, 4)
	p.Add(r2+dx1, t1).Smooth(0.25*k.WallThickness, 4)
	p.Add(r2-dx1, t0).Smooth(0.25*k.WallThickness, 4)
	p.Add(0, t0)

	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}

	// return the revolved profile
	return sdf.Revolve3D(s)
}

// dcGrate returns a grate (no crossbar)
func dcGrate(k *DrainCoverParms) (sdf.SDF3, error) {

	r := (0.5 * k.WallDiameter) - k.InnerWidth
	n := float64(k.GrateNumber)
	g := (2.0 * r) / (n + (n * k.GrateWidth) + 1.0)
	w := k.GrateWidth * g

	x := g + (0.5 * w) - r
	dx := g + w

	slots := make([]sdf.SDF3, k.GrateNumber)
	for i := 0; i < k.GrateNumber; i++ {
		l := 2.0 * math.Sqrt((r*r)-(x*x))
		k1 := TruncRectPyramidParms{
			Size:       v3.Vec{w, l, k.CoverThickness},
			BaseAngle:  0.5*sdf.Pi - k.GrateDraft,
			BaseRadius: 0.5 * w,
		}
		slot, err := TruncRectPyramid3D(&k1)
		if err != nil {
			return nil, err
		}
		slots[i] = sdf.Transform3D(slot, sdf.Translate3d(v3.Vec{x, 0, -k.CoverThickness}))
		x += dx
	}

	return sdf.Transform3D(sdf.Union3D(slots...), sdf.MirrorXY()), nil
}

// dcGrate returns a grate with a crossbar
func dcGrateCrossBar(k *DrainCoverParms) (sdf.SDF3, error) {

	r := (0.5 * k.WallDiameter) - k.InnerWidth
	n := float64(k.GrateNumber)
	g := (2.0 * r) / (n + (n * k.GrateWidth) + 1.0)
	w := k.GrateWidth * g

	x := g + (0.5 * w) - r
	dx := g + w
	dy := 0.5 * k.InnerWidth * k.CrossBarWidth

	slots := make([]sdf.SDF3, k.GrateNumber)
	for i := 0; i < k.GrateNumber; i++ {
		l := math.Sqrt((r*r)-(x*x)) - dy
		y := dy + (0.5 * l)

		k1 := TruncRectPyramidParms{
			Size:       v3.Vec{w, l, k.CoverThickness},
			BaseAngle:  0.5*sdf.Pi - k.GrateDraft,
			BaseRadius: 0.5 * w,
		}
		slot, err := TruncRectPyramid3D(&k1)
		if err != nil {
			return nil, err
		}
		slots[i] = sdf.Transform3D(slot, sdf.Translate3d(v3.Vec{x, y, -k.CoverThickness}))
		x += dx
	}

	g0 := sdf.Transform3D(sdf.Union3D(slots...), sdf.MirrorXY())
	g1 := sdf.Transform3D(g0, sdf.MirrorXZ())

	return sdf.Union3D(g0, g1), nil
}

func dcCrossWeb(k *DrainCoverParms) (sdf.SDF3, error) {

	l := k.WallDiameter - (2.0 * k.WallThickness)
	x := k.WallHeight * 0.6
	y := k.WallThickness * 0.5
	dy := 0.5 * x * math.Tan(k.WallDraft)

	// build the 2d profile
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(0, y+dy)
	p.Add(x, y-dy).Smooth(0.25*k.WallThickness, 4)
	p.Add(x, -y+dy).Smooth(0.25*k.WallThickness, 4)
	p.Add(0, -y-dy)

	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}

	web := sdf.Extrude3D(s, l)
	web = sdf.Transform3D(web, sdf.RotateY(sdf.DtoR(-90)))
	web = sdf.Transform3D(web, sdf.Translate3d(v3.Vec{0, 0, k.CoverThickness}))
	return web, nil
}

// DrainCover returns a grated drain pipe cover.
func DrainCover(k *DrainCoverParms) (sdf.SDF3, error) {

	body, err := dcBody(k)
	if err != nil {
		return nil, err
	}

	if k.CrossBarWeb && k.CrossBarWidth != 0 {
		web, err := dcCrossWeb(k)
		if err != nil {
			return nil, err
		}
		body = sdf.Union3D(body, web)
		body.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(k.WallThickness))
	}

	var grate sdf.SDF3
	if k.CrossBarWidth == 0 {
		grate, err = dcGrate(k)
	} else {
		grate, err = dcGrateCrossBar(k)
	}
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(body, grate), nil
}

//-----------------------------------------------------------------------------
