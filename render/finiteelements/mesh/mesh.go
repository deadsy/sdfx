// Package mesh provides convenient types & functions for meshes consisting of finite elements.
// Like 4-node tetrahedra, 8-node and 20-node hexahedra.
package mesh

import v3 "github.com/deadsy/sdfx/vec/v3"

// FE is a dynamic type for a mesh of finite elements like Tet4, Hex8, Hex20, ...
type FE interface {
	// Number of nodes per element.
	Npe() int
	// Number of layers along the Z axis.
	layerCount() int
	// Number of finite elements on a layer.
	feCountOnLayer(l int) int
	// Get a finite element.
	// FE vertices are returned.
	// Layer index is input.
	// FE index on layer is input.
	feVertices(l, i int) []v3.Vec
}

//-----------------------------------------------------------------------------
