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
// Keep it low-level to increase performance. Indicated by experiments.
type MeshTet4 struct {
	T      []uint32              // Index buffer. Every 4 indices would correspond to a tetrahedron of (i, j, k, l).
	V      []float32             // Vertex buffer. Every 3 floats would correspond to a vertex of (x, y, z).
	Lookup map[[3]float32]uint32 // Used to avoid repeating vertices when adding a new tetrahedron.
}

func NewMeshTet4() *MeshTet4 {
	return &MeshTet4{
		T:      []uint32{},
		V:      []float32{},
		Lookup: map[[3]float32]uint32{},
	}
}

// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (m *MeshTet4) AddTet4(a, b, c, d [3]float32) {
	m.T = append(m.T, m.addVertex(a), m.addVertex(b), m.addVertex(c), m.addVertex(d))
}

func (m *MeshTet4) addVertex(vert [3]float32) uint32 {
	// TODO: Binary insertion sort and search to eliminate extra allocation
	// TODO: Consider epsilon in comparison and use int (*100) for searching
	if vertID, ok := m.Lookup[vert]; ok {
		return vertID
	}

	m.V = append(m.V, vert[0], vert[1], vert[2])

	m.Lookup[vert] = uint32(m.countVertex())

	return uint32(m.countVertex())
}

func (m *MeshTet4) countVertex() int {
	return len(m.V)
}

// To be called after adding all tetrahedra to the mesh.
func (t *MeshTet4) Finalize() {
	// Clear memory.
	t.Lookup = nil
	runtime.GC()
}

// Number of tetrahedra on mesh.
func (m *MeshTet4) CountTet4() int {
	return len(m.T) / 4
}

// Input from 0 to number of tetrahedra on mesh.
// Don't return error to increase performance.
func (m *MeshTet4) Tet4Indicies(i int) (uint32, uint32, uint32, uint32) {
	return m.T[i*4], m.T[i*4+1], m.T[i*4+2], m.T[i*4+3]
}

// Input from 0 to number of tetrahedra on mesh.
// Don't return error to increase performance.
func (m *MeshTet4) Tet4Vertices(i int) ([3]float32, [3]float32, [3]float32, [3]float32) {
	return [3]float32{m.V[m.T[i*4]*3], m.V[m.T[i*4]*3+1], m.V[m.T[i*4]*3+2]},
		[3]float32{m.V[m.T[i*4+1]*3], m.V[m.T[i*4+1]*3+1], m.V[m.T[i*4+1]*3+2]},
		[3]float32{m.V[m.T[i*4+2]*3], m.V[m.T[i*4+2]*3+1], m.V[m.T[i*4+2]*3+2]},
		[3]float32{m.V[m.T[i*4+3]*3], m.V[m.T[i*4+3]*3+1], m.V[m.T[i*4+3]*3+2]}
}
