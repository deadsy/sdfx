package render

import (
	"runtime"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// To avoid repeating vertices on vertex buffer.
type VertexBuffer struct {
	hashTable map[[3]int32]uint32 //
	V         *[]v3.Vec           // Pointer: to be able to append to it.
}

func NewVertexBuffer(V *[]v3.Vec) *VertexBuffer {
	b := VertexBuffer{
		hashTable: map[[3]int32]uint32{},
		V:         V,
	}
	b.hashTable = make(map[[3]int32]uint32, 0)
	return &b
}

// Add vertex to buffer and get vertex ID.
// If vertex is already available on the buffer, its ID is just returned.
// All vertices would be unique. Not repeated.
func (b VertexBuffer) Id(v v3.Vec) uint32 {
	// Deduplicate by removing small details and use of epsilon
	epsilon := float64(0.0001)
	key := [3]int32{int32((v.X + epsilon) * 1000), int32((v.Y + epsilon) * 1000), int32((v.Z + epsilon) * 1000)}
	if vID, ok := b.hashTable[key]; ok {
		// Vertex already exists. It's repeated.
		return vID
	}

	// Vertex is new, so append it.
	*b.V = append(*b.V, v)

	// Store index of the appended vertex.
	b.hashTable[key] = uint32(b.vertexCount() - 1)

	// Return index of the appended vertex.
	return uint32(b.vertexCount() - 1)
}

func (b *VertexBuffer) vertexCount() int {
	return len(*b.V)
}

// To be called after adding all vertices to the vertex buffer.
func (b *VertexBuffer) Destroy() {
	// Clear memory.
	b.hashTable = nil
	runtime.GC()
}
