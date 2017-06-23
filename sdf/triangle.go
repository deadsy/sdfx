//-----------------------------------------------------------------------------
/*

Triangles and Edges

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

// 3d triangle
type Triangle3 struct {
	V [3]V3
}

// 2d triangle
type Triangle2 struct {
	V [3]V2
}

// 3d edge
type Edge3 struct {
	V [2]V3
}

// 2d edge
type Edge2 struct {
	V [2]V2
}

//-----------------------------------------------------------------------------

func NewTriangle3(a, b, c V3) *Triangle3 {
	t := Triangle3{}
	t.V[0] = a
	t.V[1] = b
	t.V[2] = c
	return &t
}

func NewTriangle2(a, b, c V2) *Triangle2 {
	t := Triangle2{}
	t.V[0] = a
	t.V[1] = b
	t.V[2] = c
	return &t
}

//-----------------------------------------------------------------------------

// return the normal vector to the plane defined by the triangle
func (t *Triangle3) Normal() V3 {
	e1 := t.V[1].Sub(t.V[0])
	e2 := t.V[2].Sub(t.V[0])
	return e1.Cross(e2).Normalize()
}

//-----------------------------------------------------------------------------

// return the super triangle of the point set, ie: A triangle enclosing all the points
func (s V2Set) SuperTriangle() *Triangle2 {

	if len(s) == 0 {
		// no points
		return nil
	}

	if len(s) == 1 {
		// a single point
		p := s[0]
		return NewTriangle2(V2{p.X - 1, p.Y - 1}, V2{p.X, p.Y + 1}, V2{p.X + 1, p.Y - 1})
	}

	// TODO
	return nil
}

//-----------------------------------------------------------------------------

// Return true if the point is within the circumcircle of the triangle.
// See: http://www.mathopenref.com/trianglecircumcircle.html
// See: http://paulbourke.net/papers/triangulate/
func (t Triangle2) InCircumcircle(p V2) bool {
	// TODO
	return false
}

//-----------------------------------------------------------------------------

// Return true if two edges are the same.
func (a Edge2) Equals(b Edge2, tolerance float64) bool {
	return a.V[0].Equals(b.V[0], tolerance) &&
		a.V[1].Equals(b.V[1], tolerance)
}

// Return true if two edges are the same.
func (a Edge3) Equals(b Edge3, tolerance float64) bool {
	return a.V[0].Equals(b.V[0], tolerance) &&
		a.V[1].Equals(b.V[1], tolerance)
}

//-----------------------------------------------------------------------------
