// Provide convenient types & functions for meshes consisting of finite elements.
// Like 4-node tetrahedra, 8-node and 20-node hexahedra.
package mesh

import v3 "github.com/deadsy/sdfx/vec/v3"

// A dynamic type for meshes of finite elements like MeshTet4, MeshHex8, MeshHex20, ...
type MeshFE interface {
	NodesPerElement() int
	layerCount() int
	feCountOnLayer(l int) int
	feVertices(l, i int) []v3.Vec
}

//-----------------------------------------------------------------------------
