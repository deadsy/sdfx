// Package mesh provides convenient types & functions for meshes consisting of finite elements.
// Like 4-node tetrahedra, 8-node and 20-node hexahedra.
package mesh

import v3 "github.com/deadsy/sdfx/vec/v3"

// FE is a dynamic type for meshes of finite elements like Tet4, Hex8, Hex20, ...
type FE interface {
	NodesPerElement() int
	layerCount() int
	feCountOnLayer(l int) int
	feVertices(l, i int) []v3.Vec
}

//-----------------------------------------------------------------------------
