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
	vCount uint32                // TODO: Remove?
}

func NewMeshTet4() *MeshTet4 {
	return &MeshTet4{
		T:      []uint32{},
		V:      []v3.Vec{},
		Lookup: map[[3]float64]uint32{},
	}
}

func (m *MeshTet4) Allocate(tetCount uint32) {
	m.T = make([]uint32, tetCount*4)

	// Just an optimization by experimenting.
	// Affects the speed according to experiments.
	m.V = make([]v3.Vec, tetCount/4*2)

	m.Lookup = make(map[[3]float64]uint32, tetCount/4*2)
}

func (m *MeshTet4) AddTet(i uint32, a, b, c, d v3.Vec) {
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
	if m.VertexCount() <= int(m.vCount) {
		m.V = append(m.V, vert)
	} else {
		m.V[m.vCount] = vert
	}
	m.Lookup[[3]float64{vert.X, vert.Y, vert.Z}] = m.vCount
	m.vCount++
	return m.vCount - 1
}

func (m *MeshTet4) VertexCount() int {
	return len(m.V)
}

func (t *MeshTet4) Finalize() {
	t.V = t.V[:t.vCount]
	// Clear memory.
	t.Lookup = nil
	runtime.GC()
}
