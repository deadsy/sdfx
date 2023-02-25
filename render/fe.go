package render

import v3 "github.com/deadsy/sdfx/vec/v3"

// Tet4 is a 3D tetrahedron consisting of 4 nodes.
// It's a kind of finite element, FE.
// https://en.wikipedia.org/wiki/Tet4
type Tet4 struct {
	V [4]v3.Vec
}

// A mesh of tetrahedra with 4 nodes.
// A sophisticated data structure for mesh is required to store tetrahedra.
// The repeated nodes would be removed.
// The element connectivity would be created with unique nodes.
type MeshTet4 struct {
	T      []uint32              // Index buffer. Every 4 indices would correspond to a tetrahedron.
	V      []v3.Vec              // Vertex buffer. All unique.
	Lookup map[[4]float32]uint32 // Used to avoid repeating vertices when adding a new tetrahedron.
}
