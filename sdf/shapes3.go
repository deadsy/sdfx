//-----------------------------------------------------------------------------
/*

Common 3D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

// Counter-Bored Hole
// l = total length
// r = hole radius
// cb_r = counter bore radius
// cb_d = counter bore depth
func CounterBore3d(l, r, cb_r, cb_d float64) SDF3 {
	s0 := NewCylinderSDF3(l, r, 0)
	s1 := NewCylinderSDF3(cb_d, cb_r, 0)
	s1 = NewTransformSDF3(s1, Translate3d(V3{0, 0, (l - cb_d) / 2}))
	return NewUnionSDF3(s0, s1)
}

//-----------------------------------------------------------------------------
