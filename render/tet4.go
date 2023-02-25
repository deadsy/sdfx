package render

import (
	"runtime"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Tet4 is a 3D tetrahedron consisting of 4 nodes.
// It's a kind of finite element, FE.
// https://en.wikipedia.org/wiki/Tetrahedron
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
	Lookup map[[3]float64]uint32 // Used to avoid repeating vertices when adding a new tetrahedron.
}

func NewMeshTet4() *MeshTet4 {
	return &MeshTet4{
		T:      []uint32{},
		V:      []v3.Vec{},
		Lookup: map[[3]float64]uint32{},
	}
}

func (m *MeshTet4) AddTet(a, b, c, d v3.Vec) {
	// Index of current tetrahedron being added.
	i := len(m.T) / 4

	// Make room for a new tetrahedron.
	m.T = append(m.T, 0, 0, 0, 0)

	// The node numbering should follow the convention of CalculiX.
	// http://www.dhondt.de/ccx_2.20.pdf
	m.T[i*4], m.T[i*4+1], m.T[i*4+2], m.T[i*4+3] = m.AddVertex(a), m.AddVertex(b), m.AddVertex(c), m.AddVertex(d)
}

func (m *MeshTet4) AddVertex(vert v3.Vec) uint32 {
	// TODO: Binary insertion sort and search to eliminate extra allocation
	// TODO: Consider epsilon in comparison and use int (*100) for searching
	if vertID, ok := m.Lookup[[3]float64{vert.X, vert.Y, vert.Z}]; ok {
		return vertID
	}

	m.V = append(m.V, vert)

	m.Lookup[[3]float64{vert.X, vert.Y, vert.Z}] = uint32(m.VertexCount())

	return uint32(m.VertexCount())
}

func (m *MeshTet4) VertexCount() int {
	return len(m.V)
}

func (t *MeshTet4) Finalize() {
	// Clear memory.
	t.Lookup = nil
	runtime.GC()
}
