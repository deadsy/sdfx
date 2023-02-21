package render

import v3 "github.com/deadsy/sdfx/vec/v3"

// Tetrahedron is a 3D tetrahedron.
// It's a type of finite element, FE.
// https://en.wikipedia.org/wiki/Tetrahedron
type Tetrahedron struct {
	V [4]v3.Vec
}
