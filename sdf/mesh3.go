//-----------------------------------------------------------------------------
/*

3D Mesh, 3d triangles connected to create manifold objects.

*/
//-----------------------------------------------------------------------------

package sdf

import v3 "github.com/deadsy/sdfx/vec/v3"

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
