//-----------------------------------------------------------------------------
/*

Common 3D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

// Counter Bored Hole
// l = total length
// r = hole radius
// cb_r = counter bore radius
// cb_d = counter bore depth
func CounterBored_Hole3d(l, r, cb_r, cb_d float64) SDF3 {
	s0 := NewCylinderSDF3(l, r, 0)
	s1 := NewCylinderSDF3(cb_d, cb_r, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - cb_d) / 2}))
	return Union3D(s0, s1)
}

// Chamfered Hole (45 degrees)
// l = total length
// r = hole radius
// ch_r = chamfer radius
func Chamfered_Hole3d(l, r, ch_r float64) SDF3 {
	s0 := NewCylinderSDF3(l, r, 0)
	s1 := NewConeSDF3(ch_r, r, r+ch_r, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - ch_r) / 2}))
	return Union3D(s0, s1)
}

// Countersunk Hole (45 degrees)
// l = total length
// r = hole radius
func CounterSunk_Hole3d(l, r float64) SDF3 {
	return Chamfered_Hole3d(l, r, r)
}

//-----------------------------------------------------------------------------
