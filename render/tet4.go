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

// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (m *MeshTet4) AddTet4(a, b, c, d v3.Vec) {
	m.T = append(m.T, m.addVertex(a), m.addVertex(b), m.addVertex(c), m.addVertex(d))
}

func (m *MeshTet4) addVertex(vert v3.Vec) uint32 {
	// TODO: Binary insertion sort and search to eliminate extra allocation
	// TODO: Consider epsilon in comparison and use int (*100) for searching
	if vertID, ok := m.Lookup[[3]float64{vert.X, vert.Y, vert.Z}]; ok {
		return vertID
	}

	m.V = append(m.V, vert)

	m.Lookup[[3]float64{vert.X, vert.Y, vert.Z}] = uint32(m.vertexCount())

	return uint32(m.vertexCount())
}

func (m *MeshTet4) vertexCount() int {
	return len(m.V)
}

// To be called after adding all tetrahedra to the mesh.
func (t *MeshTet4) Finalize() {
	// Clear memory.
	t.Lookup = nil
	runtime.GC()
}

// Number of tetrahedra on mesh.
func (m *MeshTet4) Tet4Count() int {
	return len(m.T) / 4
}

// Input from 0 to number of tetrahedra on mesh.
// Don't return error to increase performance.
func (m *MeshTet4) Tet4Indicies(i int) (uint32, uint32, uint32, uint32) {
	return m.T[i*4], m.T[i*4+1], m.T[i*4+2], m.T[i*4+3]
}

// Input from 0 to number of tetrahedra on mesh.
// Don't return error to increase performance.
func (m *MeshTet4) Tet4Vertices(i int) (v3.Vec, v3.Vec, v3.Vec, v3.Vec) {
	return m.V[m.T[i*4]], m.V[m.T[i*4+1]], m.V[m.T[i*4+2]], m.V[m.T[i*4+3]]
}
