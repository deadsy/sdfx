package buffer

import (
	"runtime"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Vertex buffer for finite elements.
// Vertex buffer avoids repeating vertices by a hash table.
// The same vertex buffer is used for 4-node tetrahedra, 8-node hexahedra, and others.
type VB struct {
	// To store index of vertices. Repeated vertices would have the same index.
	hashTable map[[3]int32]uint32
	// To store coordinates of vertices.
	V []v3.Vec
}

func NewVB() *VB {
	b := VB{
		hashTable: map[[3]int32]uint32{},
		V:         []v3.Vec{},
	}

	return &b
}

// Add vertex to buffer and get vertex ID.
// If vertex is already available on the buffer, its ID is just returned.
// So, all vertices will be unique. Not repeated.
func (b *VB) Id(v v3.Vec) uint32 {
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

func (b *VB) VertexCount() int {
	return len(b.V)
}

// To be called after adding all vertices to the vertex buffer.
// Call if you are sure that no new vertex will be added to the vertex buffer.
func (b *VB) DestroyHashTable() {
	// Clear memory.
	b.hashTable = nil
	runtime.GC()
}

func (b *VB) Vertex(i uint32) v3.Vec {
	return b.V[i]
}
