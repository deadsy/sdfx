//-----------------------------------------------------------------------------
/*

Finials for pillar decorations.

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func square1(l float64) SDF2 {
	return Polygon2D(Nagon(4, l*math.Sqrt(0.5)))
}

func square2(l float64) SDF2 {

	h := l * 0.5
	r := l * 0.1
	n := 5

	s := NewPolygon()
	s.Add(h, -h).Smooth(r, n)
	s.Add(h, h).Smooth(r, n)
	s.Add(-h, h).Smooth(r, n)
	s.Add(-h, -h).Smooth(r, n)
	s.Close()

	return Polygon2D(s.Vertices())
}

//-----------------------------------------------------------------------------

func finial2() {
	base := 100.0
	base_height := 20.0
	column_radius := 15.0
	column_height := 60.0
	ball_radius := 45.0
	column_ofs := (column_height + base_height) / 2
	ball_ofs := (base_height / 2) + column_height + ball_radius*0.8
	round := ball_radius / 5

	//s0 := Offset2D(square2(base), base*0.1)
	s0 := square2(base)
	s1 := Circle2D(column_radius)

	column_3d := Loft3D(s0, s1, column_height, 0)
	column_3d = Transform3D(column_3d, Translate3d(V3{0, 0, column_ofs}))

	ball_3d := Sphere3D(ball_radius)
	ball_3d = Transform3D(ball_3d, Translate3d(V3{0, 0, ball_ofs}))

	base_3d := Extrude3D(s0, base_height)

	bc_3d := Union3D(column_3d, ball_3d)
	bc_3d.(*UnionSDF3).SetMin(PolyMin(round))

	RenderSTLSlow(Union3D(bc_3d, base_3d), 300, "f2.stl")
}

//-----------------------------------------------------------------------------

func finial1() {
	base := 100.0
	base_height := 20.0
	column_radius := 15.0
	column_height := 60.0
	ball_radius := 45.0
	column_ofs := (column_height + base_height) / 2
	ball_ofs := (base_height / 2) + column_height + ball_radius*0.8
	round := ball_radius / 5

	s0 := Polygon2D(Nagon(4, base*math.Sqrt(0.5)))
	s1 := Circle2D(column_radius)

	column_3d := Loft3D(s0, s1, column_height, 0)
	column_3d = Transform3D(column_3d, Translate3d(V3{0, 0, column_ofs}))

	ball_3d := Sphere3D(ball_radius)
	ball_3d = Transform3D(ball_3d, Translate3d(V3{0, 0, ball_ofs}))

	base_3d := Extrude3D(s0, base_height)

	bc_3d := Union3D(column_3d, ball_3d)
	bc_3d.(*UnionSDF3).SetMin(PolyMin(round))

	RenderSTLSlow(Union3D(bc_3d, base_3d), 300, "f1.stl")
}

//-----------------------------------------------------------------------------

func main() {
	finial1()
	finial2()
}

//-----------------------------------------------------------------------------
