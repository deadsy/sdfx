package render

import v3 "github.com/deadsy/sdfx/vec/v3"

// Tetrahedron is a 3D tetrahedron.
// It's a kind of finite element, FE.
// https://en.wikipedia.org/wiki/Tetrahedron
type Tetrahedron struct {
	V [4]v3.Vec
}

// A mesh of tetrahedra with 4 nodes.
type MeshTet4 struct {
}
