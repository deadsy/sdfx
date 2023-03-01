package buffer

import (
	"runtime"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Vertex buffer for 4-node tetrahedra.
// To avoid repeating vertices on vertex buffer.
type Tet4VB struct {
	// To store index of vertices. Repeated vertices would have the same index.
	hashTable map[[3]int32]uint32
	// To store coordinates of vertices.
	V []v3.Vec
}

func NewTet4VB() *Tet4VB {
	b := Tet4VB{
		hashTable: map[[3]int32]uint32{},
		V:         []v3.Vec{},
	}

	return &b
}

// Add vertex to buffer and get vertex ID.
// If vertex is already available on the buffer, its ID is just returned.
// All vertices would be unique. Not repeated.
func (b *Tet4VB) Id(v v3.Vec) uint32 {
	// Deduplicate by removing small details and use of epsilon
	epsilon := float64(0.0001)
	key := [3]int32{int32((v.X + epsilon) * 1000), int32((v.Y + epsilon) * 1000), int32((v.Z + epsilon) * 1000)}
	if vID, ok := b.hashTable[key]; ok {
		// Vertex already exists. It's repeated.
		return vID
	}

	// Vertex is new, so append it.
	b.V = append(b.V, v)

	// Store index of the appended vertex.
	b.hashTable[key] = uint32(b.VertexCount() - 1)

	// Return index of the appended vertex.
	return uint32(b.VertexCount() - 1)
}

func (b *Tet4VB) VertexCount() int {
	return len(b.V)
}

// To be called after adding all vertices to the vertex buffer.
// Call if you are sure that no new vertex will be added to the vertex buffer.
func (b *Tet4VB) DestroyHashTable() {
	// Clear memory.
	b.hashTable = nil
	runtime.GC()
}

func (b *Tet4VB) Vertex(i uint32) v3.Vec {
	return b.V[i]
}
