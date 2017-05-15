//-----------------------------------------------------------------------------
/*

Finials for pillar decorations.

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func finial1() {
	base := 50.0
	base_height := 20.0
	column_radius := 15.0
	column_height := 70.0
	ball_radius := 40.0
	column_ofs := (column_height + base_height) / 2
	ball_ofs := (base_height / 2) + column_height + ball_radius*0.8
	round := ball_radius / 5

	s0 := Polygon2D(Nagon(4, base))
	s1 := Circle2D(column_radius)

	column_3d := Loft3D(s0, s1, column_height, 0)
	column_3d = Transform3D(column_3d, Translate3d(V3{0, 0, column_ofs}))

	ball_3d := Sphere3D(ball_radius)
	ball_3d = Transform3D(ball_3d, Translate3d(V3{0, 0, ball_ofs}))

	base_3d := Extrude3D(s0, base_height)

	bc_3d := Union3D(column_3d, ball_3d)
	bc_3d.(*UnionSDF3).SetMin(PolyMin(round))

	RenderSTL(Union3D(bc_3d, base_3d), 300, "f1.stl")
}

//-----------------------------------------------------------------------------

func main() {
	finial1()
}

//-----------------------------------------------------------------------------
