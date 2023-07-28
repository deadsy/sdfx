//-----------------------------------------------------------------------------
/*

Delaunay Triangulation

See:
http://www.mathopenref.com/trianglecircumcircle.html
http://paulbourke.net/papers/triangulate/
Computational Geometry, Joseph O'Rourke, 2nd edition, Code 5.1

*/
//-----------------------------------------------------------------------------

package render

import (
	"errors"
	"sort"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// TriangleI is a 2d/3d triangle referencing a list of vertices.
type TriangleI [3]int

// ToTriangle2 given vertex indices and the vertex array, return the triangle with real vertices.
func (t TriangleI) ToTriangle2(p []v2.Vec) sdf.Triangle2 {
	return sdf.Triangle2{p[t[0]], p[t[1]], p[t[2]]}
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

// superTriangle return the super triangle of a point set, ie: 3 vertices enclosing all points.
func superTriangle(vs v2.VecSet) (sdf.Triangle2, error) {

	if len(vs) == 0 {
		return sdf.Triangle2{}, errors.New("no vertices")
	}

	var p v2.Vec
	var k float64

	if len(vs) == 1 {
		// a single point
		p := vs[0]
		k = p.MaxComponent() * 0.125
		if k == 0 {
			k = 1
		}
	} else {
		b := sdf.Box2{vs.Min(), vs.Max()}
		p = b.Center()
		k = b.Size().MaxComponent() * 2.0
	}

	// Note: super triangles should be large enough to avoid having the circumcenter of
	// any triangle lie outside of the super triangle. This is kludgey. For thin triangles
	// on the hull the circumcenter is going to be arbitrarily far away.
	k *= 4096.0

	p0 := p.Add(v2.Vec{-k, -k})
	p1 := p.Add(v2.Vec{0, k})
	p2 := p.Add(v2.Vec{k, -k})
	return sdf.Triangle2{p0, p1, p2}, nil
}

//-----------------------------------------------------------------------------

// Delaunay2d returns the delaunay triangulation of a 2d point set.
func Delaunay2d(vs v2.VecSet) (TriangleISet, error) {

	// number of vertices
	n := len(vs)

	// sort the vertices by x value
	sort.Sort(v2.VecSetByX(vs))

	// work out the super triangle
	t, err := superTriangle(vs)
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
func Delaunay2dSlow(vs v2.VecSet) (TriangleISet, error) {

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

		p0 := conv.V2ToV3(vs[t[0]], z[t[0]])
		p1 := conv.V2ToV3(vs[t[1]], z[t[1]])
		p2 := conv.V2ToV3(vs[t[2]], z[t[2]])

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
			pi := conv.V2ToV3(v, z[i])
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
		if nextCombination(n, c) == false {
			break
		}
	}

	// done
	return ts, nil
}

//-----------------------------------------------------------------------------
