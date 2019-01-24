//-----------------------------------------------------------------------------
/*

Delaunay Triangulation

See:
http://www.mathopenref.com/trianglecircumcircle.html
http://paulbourke.net/papers/triangulate/
Computational Geometry, Joseph O'Rourke, 2nd edition, Code 5.1

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"sort"
)

//-----------------------------------------------------------------------------

// TriangleI is a 2d/3d triangle referencing a list of vertices.
type TriangleI [3]int

// ToTriangle2 given vertex indices and the vertex array, return the triangle with real vertices.
func (t TriangleI) ToTriangle2(p []V2) Triangle2 {
	return Triangle2{p[t[0]], p[t[1]], p[t[2]]}
}

// TriangleIByIndex sorts triangles by index.
type TriangleIByIndex []TriangleI

func (a TriangleIByIndex) Len() int {
	return len(a)
}
func (a TriangleIByIndex) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a TriangleIByIndex) Less(i, j int) bool {
	if a[i][0] < a[j][0] {
		return true
	}
	if a[i][0] == a[j][0] && a[i][1] < a[j][1] {
		return true
	}
	if a[i][1] == a[j][1] && a[i][2] < a[j][2] {
		return true
	}
	return false
}

// Canonical converts a triangle to it's lowest index first form.
// Preserve the winding order.
func (t *TriangleI) Canonical() {
	if t[0] < t[1] && t[0] < t[2] {
		// ok
		return
	}
	if t[1] < t[0] && t[1] < t[2] {
		// t[1] is the smallest
		tmp := t[0]
		t[0] = t[1]
		t[1] = t[2]
		t[2] = tmp
		return
	}
	// t[2] is the smallest
	tmp := t[2]
	t[2] = t[1]
	t[1] = t[0]
	t[0] = tmp
}

// TriangleISet is a set of triangles defined by vertice indices.
type TriangleISet []TriangleI

// Canonical converts a triangle set to it's canonical form.
// This common form is used to facilitate comparison
// between the results of different implementations.
func (ts TriangleISet) Canonical() []TriangleI {
	// convert each triangle to it's lowest index first form
	for i := range ts {
		ts[i].Canonical()
	}
	// sort the triangles by index
	sort.Sort(TriangleIByIndex(ts))
	return ts
}

// Equals tests two triangle sets for equality.
func (ts TriangleISet) Equals(s TriangleISet) bool {
	if len(ts) != len(s) {
		return false
	}
	ts = ts.Canonical()
	s = s.Canonical()
	for i := range ts {
		if (ts[i][0] != s[i][0]) ||
			(ts[i][1] != s[i][1]) ||
			(ts[i][2] != s[i][2]) {
			return false
		}
	}
	return true
}

//-----------------------------------------------------------------------------

// EdgeI is a 2d/3d edge referencing a list of vertices.
type EdgeI [2]int

//-----------------------------------------------------------------------------

// SuperTriangle return the super triangle of a point set, ie: 3 vertices enclosing all points.
func (vs V2Set) SuperTriangle() (Triangle2, error) {

	if len(vs) == 0 {
		return Triangle2{}, errors.New("no vertices")
	}

	var p V2
	var k float64

	if len(vs) == 1 {
		// a single point
		p := vs[0]
		k = p.MaxComponent() * 0.125
		if k == 0 {
			k = 1
		}
	} else {
		b := Box2{vs.Min(), vs.Max()}
		p = b.Center()
		k = b.Size().MaxComponent() * 2.0
	}

	// Note: super triangles should be large enough to avoid having the circumcenter of
	// any triangle lie outside of the super triangle. This is kludgey. For thin triangles
	// on the hull the circumcenter is going to be arbitrarily far away.
	k *= 4096.0

	p0 := p.Add(V2{-k, -k})
	p1 := p.Add(V2{0, k})
	p2 := p.Add(V2{k, -k})
	return Triangle2{p0, p1, p2}, nil
}

//-----------------------------------------------------------------------------

// Circumcenter returns the circumcenter of a triangle.
func (t Triangle2) Circumcenter() (V2, error) {

	var m1, m2, mx1, mx2, my1, my2 float64
	var xc, yc float64

	x1 := t[0].X
	x2 := t[1].X
	x3 := t[2].X

	y1 := t[0].Y
	y2 := t[1].Y
	y3 := t[2].Y

	fabsy1y2 := Abs(y1 - y2)
	fabsy2y3 := Abs(y2 - y3)

	// Check for coincident points
	if fabsy1y2 < EPSILON && fabsy2y3 < EPSILON {
		return V2{}, errors.New("coincident points")
	}

	if fabsy1y2 < EPSILON {
		m2 = -(x3 - x2) / (y3 - y2)
		mx2 = (x2 + x3) / 2.0
		my2 = (y2 + y3) / 2.0
		xc = (x2 + x1) / 2.0
		yc = m2*(xc-mx2) + my2
	} else if fabsy2y3 < EPSILON {
		m1 = -(x2 - x1) / (y2 - y1)
		mx1 = (x1 + x2) / 2.0
		my1 = (y1 + y2) / 2.0
		xc = (x3 + x2) / 2.0
		yc = m1*(xc-mx1) + my1
	} else {
		m1 = -(x2 - x1) / (y2 - y1)
		m2 = -(x3 - x2) / (y3 - y2)
		mx1 = (x1 + x2) / 2.0
		mx2 = (x2 + x3) / 2.0
		my1 = (y1 + y2) / 2.0
		my2 = (y2 + y3) / 2.0
		xc = (m1*mx1 - m2*mx2 + my2 - my1) / (m1 - m2)
		if fabsy1y2 > fabsy2y3 {
			yc = m1*(xc-mx1) + my1
		} else {
			yc = m2*(xc-mx2) + my2
		}
	}

	return V2{xc, yc}, nil
}

// InCircumcircle return inside == true if the point is inside the circumcircle of the triangle.
// Returns done == true if the vertex and the subsequent x-ordered vertices are outside the circumcircle.
func (t Triangle2) InCircumcircle(p V2) (inside, done bool) {
	c, err := t.Circumcenter()
	if err != nil {
		inside = false
		done = true
		return
	}

	// radius squared of circumcircle
	dx := t[0].X - c.X
	dy := t[0].Y - c.Y
	r2 := dx*dx + dy*dy

	// distance squared from circumcenter to point
	dx = p.X - c.X
	dy = p.Y - c.Y
	d2 := dx*dx + dy*dy

	// is the point within the circumcircle?
	inside = d2-r2 <= EPSILON

	// If this vertex has an x-value beyond the circumcenter and the distance based on the x-delta
	// is greater than the circumradius, then this triangle is done for this and all subsequent vertices
	// since the vertex list has been sorted by x-value.
	done = (dx > 0) && (dx*dx > r2)

	return
}

//-----------------------------------------------------------------------------

// Delaunay2d returns the delaunay triangulation of a 2d point set.
func (vs V2Set) Delaunay2d() (TriangleISet, error) {

	// number of vertices
	n := len(vs)

	// sort the vertices by x value
	sort.Sort(V2SetByX(vs))

	// work out the super triangle
	t, err := vs.SuperTriangle()
	if err != nil {
		return nil, err
	}

	// add the super triangle to the vertex set
	vs = append(vs, t[:]...)

	// allocate the triangles
	k := (2 * n) + 1
	ts := make([]TriangleI, 0, k)
	done := make([]bool, 0, k)

	// set the super triangle as the 0th triangle
	ts = append(ts, TriangleI{n, n + 1, n + 2})
	done = append(done, false)

	// Add the vertices one at a time into the mesh
	// Note: we don't iterate over the super triangle vertices
	for i := 0; i < n; i++ {
		v := vs[i]

		// Create the edge buffer.
		// If the vertex lies inside the circumcircle of the triangle
		// then the three edges of that triangle are added to the edge
		// buffer and that triangle is removed.
		es := make([]EdgeI, 0, 64)
		nt := len(ts)
		for j := 0; j < nt; j++ {
			if done[j] {
				continue
			}

			t := ts[j].ToTriangle2(vs)
			inside, complete := t.InCircumcircle(v)
			done[j] = complete

			if inside {
				// add the triangle edges to the edge set
				es = append(es, EdgeI{ts[j][0], ts[j][1]})
				es = append(es, EdgeI{ts[j][1], ts[j][2]})
				es = append(es, EdgeI{ts[j][2], ts[j][0]})
				// remove the triangle (copy in the tail)
				ts[j] = ts[nt-1]
				done[j] = done[nt-1]
				nt--
				j--
			}
		}

		// re-size the triangle/done sets
		ts = ts[:nt]
		done = done[:nt]

		// Tag duplicate edges for removal.
		for j := 0; j < len(es)-1; j++ {
			for k := j + 1; k < len(es); k++ {
				if (es[j][0] == es[k][1] && es[j][1] == es[k][0]) ||
					(es[j][1] == es[k][1] && es[j][0] == es[k][0]) {
					es[j] = EdgeI{-1, -1}
					es[k] = EdgeI{-1, -1}
				}
			}
		}

		// Form new triangles for the current point skipping over any duplicate edges.
		for _, e := range es {
			if e[0] < 0 || e[1] < 0 {
				continue
			}
			ts = append(ts, TriangleI{e[0], e[1], i})
			done = append(done, false)
		}
	}

	// remove any triangles with vertices from the super triangle
	nt := len(ts)
	for j := 0; j < nt; j++ {
		t := ts[j]
		if t[0] >= n || t[1] >= n || t[2] >= n {
			// remove the triangle (copy in the tail)
			ts[j] = ts[nt-1]
			nt--
			j--
		}
	}

	// re-size the triangle set
	ts = ts[:nt]

	// done
	return ts, nil
}

//-----------------------------------------------------------------------------

// Delaunay2dSlow returns the delaunay triangulation of a 2d point set.
// This is a slow reference implementation for testing faster algorithms.
// See: Computational Geometry, Joseph O'Rourke, 2nd edition, Code 5.1
func (vs V2Set) Delaunay2dSlow() (TriangleISet, error) {

	// number of vertices
	n := len(vs)
	if n < 3 {
		return nil, errors.New("number of vertices < 3")
	}

	// map the 2d points onto a 3d parabola
	z := make([]float64, n)
	for i, v := range vs {
		z[i] = v.Length2()
	}

	// make the set of triangles
	ts := make([]TriangleI, 0, (2*n)+1)

	// iterate through all the possible triangles
	c := []int{0, 1, 2}

	for {

		t := TriangleI{c[0], c[1], c[2]}

		p0 := vs[t[0]].ToV3(z[t[0]])
		p1 := vs[t[1]].ToV3(z[t[1]])
		p2 := vs[t[2]].ToV3(z[t[2]])

		norm := p1.Sub(p0).Cross(p2.Sub(p1))

		// we want to consider triangles whose normal faces in the -ve z direction
		if norm.Z > 0 {
			// swap the triangle handed-ness to flip the normal
			t[1], t[2] = t[2], t[1]
			norm = norm.MulScalar(-1.0)
		}

		// Are there any vertices below this plane?
		hull := true
		for i, v := range vs {
			if i == t[0] || i == t[1] || i == t[2] {
				// on the plane
				continue
			}
			pi := v.ToV3(z[i])
			if pi.Sub(p0).Dot(norm) > 0 {
				// below the plane
				hull = false
				break
			}
		}

		if hull {
			// there are no vertices below this triangles plane
			// so it is part of the lower convex hull.
			ts = append(ts, t)
		}

		// get the next triangle
		if NextCombination(n, c) == false {
			break
		}
	}

	// done
	return ts, nil
}

//-----------------------------------------------------------------------------
