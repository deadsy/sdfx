//-----------------------------------------------------------------------------
/*

Triangles

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

// Triangle3 is a 3D triangle
type Triangle3 struct {
	V [3]V3
}

// Triangle2 is a 2D triangle
type Triangle2 [3]V2

//-----------------------------------------------------------------------------

// NewTriangle3 returns a new 3D triangle.
func NewTriangle3(a, b, c V3) *Triangle3 {
	t := Triangle3{}
	t.V[0] = a
	t.V[1] = b
	t.V[2] = c
	return &t
}

//-----------------------------------------------------------------------------

// Normal returns the normal vector to the plane defined by the 3D triangle.
func (t *Triangle3) Normal() V3 {
	e1 := t.V[1].Sub(t.V[0])
	e2 := t.V[2].Sub(t.V[0])
	return e1.Cross(e2).Normalize()
}

//-----------------------------------------------------------------------------
