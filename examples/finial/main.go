//-----------------------------------------------------------------------------
/*

Finials for pillar decorations.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func square1(l float64) (sdf.SDF2, error) {
	return sdf.Polygon2D(sdf.Nagon(4, l*math.Sqrt(0.5)))
}

func square2(l float64) (sdf.SDF2, error) {

	h := l * 0.5
	r := l * 0.1
	n := 5

	s := sdf.NewPolygon()
	s.Add(h, -h).Smooth(r, n)
	s.Add(h, h).Smooth(r, n)
	s.Add(-h, h).Smooth(r, n)
	s.Add(-h, -h).Smooth(r, n)
	s.Close()

	return sdf.Polygon2D(s.Vertices())
}

//-----------------------------------------------------------------------------

func finial2() (sdf.SDF3, error) {
	base := 100.0
	base_height := 20.0
	column_radius := 15.0
	column_height := 60.0
	ball_radius := 45.0
	column_ofs := (column_height + base_height) / 2
	ball_ofs := (base_height / 2) + column_height + ball_radius*0.8
	round := ball_radius / 5

	//s0 := Offset2D(square2(base), base*0.1)
	s0, err := square2(base)
	if err != nil {
		return nil, err
	}
	s1, err := sdf.Circle2D(column_radius)
	if err != nil {
		return nil, err
	}

	column_3d, err := sdf.Loft3D(s0, s1, column_height, 0)
	if err != nil {
		return nil, err
	}
	column_3d = sdf.Transform3D(column_3d, sdf.Translate3d(sdf.V3{0, 0, column_ofs}))

	ball_3d, err := sdf.Sphere3D(ball_radius)
	if err != nil {
		return nil, err
	}

	ball_3d = sdf.Transform3D(ball_3d, sdf.Translate3d(sdf.V3{0, 0, ball_ofs}))

	base_3d := sdf.Extrude3D(s0, base_height)

	bc_3d := sdf.Union3D(column_3d, ball_3d)
	bc_3d.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(round))

	return sdf.Union3D(bc_3d, base_3d), nil
}

//-----------------------------------------------------------------------------

func finial1() (sdf.SDF3, error) {
	base := 100.0
	base_height := 20.0
	column_radius := 15.0
	column_height := 60.0
	ball_radius := 45.0
	column_ofs := (column_height + base_height) / 2
	ball_ofs := (base_height / 2) + column_height + ball_radius*0.8
	round := ball_radius / 5

	s0, err := sdf.Polygon2D(sdf.Nagon(4, base*math.Sqrt(0.5)))
	if err != nil {
		return nil, err
	}
	s1, err := sdf.Circle2D(column_radius)
	if err != nil {
		return nil, err
	}

	column_3d, err := sdf.Loft3D(s0, s1, column_height, 0)
	if err != nil {
		return nil, err
	}
	column_3d = sdf.Transform3D(column_3d, sdf.Translate3d(sdf.V3{0, 0, column_ofs}))

	ball_3d, err := sdf.Sphere3D(ball_radius)
	if err != nil {
		return nil, err
	}

	ball_3d = sdf.Transform3D(ball_3d, sdf.Translate3d(sdf.V3{0, 0, ball_ofs}))

	base_3d := sdf.Extrude3D(s0, base_height)

	bc_3d := sdf.Union3D(column_3d, ball_3d)
	bc_3d.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(round))

	return sdf.Union3D(bc_3d, base_3d), nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := finial1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTLSlow(s, 300, "f1.stl")

	s, err = finial2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTLSlow(s, 300, "f2.stl")
}

//-----------------------------------------------------------------------------
