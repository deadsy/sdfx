//-----------------------------------------------------------------------------
/*

Triangles

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

// Triangle2 is a 2D triangle
type Triangle2 [3]V2

//-----------------------------------------------------------------------------

// Triangle3 is a 3D triangle
type Triangle3 struct {
	V [3]V3
}

// NewTriangle3 returns a new 3D triangle.
func NewTriangle3(a, b, c V3) *Triangle3 {
	t := Triangle3{}
	t.V[0] = a
	t.V[1] = b
	t.V[2] = c
	return &t
}

// Normal returns the normal vector to the plane defined by the 3D triangle.
func (t *Triangle3) Normal() V3 {
	e1 := t.V[1].Sub(t.V[0])
	e2 := t.V[2].Sub(t.V[0])
	return e1.Cross(e2).Normalize()
}

// Degenerate returns true if the triangle is degenerate.
func (t *Triangle3) Degenerate(tolerance float64) bool {
	// check for identical vertices
	if t.V[0].Equals(t.V[1], tolerance) {
		return true
	}
	if t.V[1].Equals(t.V[2], tolerance) {
		return true
	}
	if t.V[2].Equals(t.V[0], tolerance) {
		return true
	}
	// TODO more tests needed
	return false
}

//-----------------------------------------------------------------------------
