//-----------------------------------------------------------------------------
/*

3D Mesh, 3d triangles connected to create manifold objects.

See:


*/
//-----------------------------------------------------------------------------

package sdf

import (
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// triangleInfo stores pre-calculated triangle information.
type triangleInfo struct {
	m M44       // rotate to XY matrix
	t [3]v2.Vec // transformed triangle vertices
	e [3]v2.Vec // anti-clockwise (from +z) unit edge vectors
	n [3]v2.Vec // outward pointing unit normals to edge vectors
}

// newTriangleInfo pre-calculates the triangle information.
func newTriangleInfo(t *Triangle3) *triangleInfo {

	m := t.rotateToXY()

	x1 := m.MulPosition(t[1]) // maps to x axis
	x2 := m.MulPosition(t[2]) // maps to xy plane

	// triangle vertices on xy plane
	t0 := v2.Vec{0, 0}
	t1 := v2.Vec{x1.X, 0}
	t2 := v2.Vec{x2.X, x2.Y}

	// triangle edge vectors
	e0 := t1.Sub(t0).Normalize()
	e1 := t2.Sub(t1).Normalize()
	e2 := t0.Sub(t2).Normalize()

	// normals to triangle edges
	n0 := v2.Vec{e0.Y, -e0.X}
	n1 := v2.Vec{e1.Y, -e1.X}
	n2 := v2.Vec{e2.Y, -e2.X}

	return &triangleInfo{
		m: m,
		t: [3]v2.Vec{t0, t1, t2},
		e: [3]v2.Vec{e0, e1, e2},
		n: [3]v2.Vec{n0, n1, n2},
	}
}

func convertTriangles(tSet []*Triangle3) []*triangleInfo {
	ti := make([]*triangleInfo, len(tSet))
	for i := range tSet {
		ti[i] = newTriangleInfo(tSet[i])
	}
	return ti
}

// minDistance2 returns the minium distance squared between a point and the triangle.
func (a *triangleInfo) minDistance2(p v3.Vec) float64 {

	// See: https://www.researchgate.net/publication/243787422_3D_Distance_from_a_Point_to_a_Triangle
	// We use the 2d method. Rotate/translate the point and triangle so the triangle is in the XY
	// plane and then project p onto the xy plane for consideration of the plane/edge/vertex cases.

	// rotate/translate the point so the triangle is in the xy plane.
	p = a.m.MulPosition(p)

	// pXY is the closest point on the XY plane
	pXY := v2.Vec{p.X, p.Y}

	// pXY wrt the triangle vertices
	pXY0 := pXY.Sub(a.t[0])
	pXY1 := pXY.Sub(a.t[1])
	pXY2 := pXY.Sub(a.t[2])

	d2 := p.Z * p.Z

	// edge 0
	if pXY0.Cross(a.e[0]) > 0 {
		// right of edge 0
		if pXY0.Cross(a.n[0]) > 0 {
			// closest to vertex 0
			return d2 + pXY0.Length2()
		}
		if pXY1.Cross(a.n[0]) < 0 {
			// closest to vertex 1
			return d2 + pXY1.Length2()
		}
		// closest to edge 0
		dn := pXY0.Dot(a.n[0])
		return d2 + (dn * dn)
	}

	// edge 1
	if pXY1.Cross(a.e[1]) > 0 {
		// right of edge 1
		if pXY1.Cross(a.n[1]) > 0 {
			// closest to vertex 1
			return d2 + pXY1.Length2()
		}
		if pXY2.Cross(a.n[1]) < 0 {
			// closest to vertex 2
			return d2 + pXY2.Length2()
		}
		// closest to edge 1
		dn := pXY1.Dot(a.n[1])
		return d2 + (dn * dn)
	}

	// edge 2
	if pXY2.Cross(a.e[2]) > 0 {
		// right of edge 2
		if pXY2.Cross(a.n[2]) > 0 {
			// closest to vertex 2
			return d2 + pXY2.Length2()
		}
		if pXY0.Cross(a.n[2]) < 0 {
			// closest to vertex 0
			return d2 + pXY0.Length2()
		}
		// closest to edge 2
		dn := pXY2.Dot(a.n[2])
		return d2 + (dn * dn)
	}

	// left of all edges, pXY is in the triangle
	return d2
}

//-----------------------------------------------------------------------------
// Mesh3D. 3D mesh evaluation with octree speedup.

// MeshSDF3 is an SDF3 made from a set of 3d triangles.
type MeshSDF3 struct {
	mesh []*Triangle3
	bb   Box3 // bounding box
}

// Mesh3D returns an SDF3 made from a set of triangles.
func Mesh3D(mesh []*Triangle3) (SDF3, error) {
	n := len(mesh)
	if n == 0 {
		return nil, ErrMsg("no triangles")
	}

	// work out the bounding box
	bb := mesh[0].BoundingBox()
	for _, t := range mesh {
		bb = bb.Extend(t.BoundingBox())
	}

	return &MeshSDF3{
		mesh: mesh,
		bb:   bb,
	}, nil
}

// Evaluate returns the minimum distance for a 2d mesh.
func (s *MeshSDF3) Evaluate(p v3.Vec) float64 {
	// TODO
	return 0
}

// BoundingBox returns the bounding box of a 3d mesh.
func (s *MeshSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Mesh3D Slow. Provided for testing and benchmarking purposes.

// MeshSDF3Slow is an SDF3 made from a set of 3d triangles.
type MeshSDF3Slow struct {
	mesh []*Triangle3
	bb   Box3 // bounding box
}

// Mesh3DSlow returns an SDF3 made from a set of triangles.
func Mesh3DSlow(mesh []*Triangle3) (SDF3, error) {
	n := len(mesh)
	if n == 0 {
		return nil, ErrMsg("no triangles")
	}

	// work out the bounding box
	bb := mesh[0].BoundingBox()
	for _, t := range mesh {
		bb = bb.Extend(t.BoundingBox())
	}

	return &MeshSDF3Slow{
		mesh: mesh,
		bb:   bb,
	}, nil
}

// Evaluate returns the minimum distance for a 2d mesh.
func (s *MeshSDF3Slow) Evaluate(p v3.Vec) float64 {
	// TODO
	return 0
}

// BoundingBox returns the bounding box of a 3d mesh.
func (s *MeshSDF3Slow) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
