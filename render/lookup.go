package render

import (
	"runtime"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// To avoid repeating vertices on vertex buffer.
type Lookup struct {
	hashTable map[[3]int32]uint32 //
	V         *[]v3.Vec           // Pointer: to be able to append to it.
}

func NewLookup(V *[]v3.Vec) *Lookup {
	l := Lookup{
		hashTable: map[[3]int32]uint32{},
		V:         V,
	}
	l.hashTable = make(map[[3]int32]uint32, 0)
	return &l
}

// Add vertex to buffer and get vertex ID.
// All vertices would be unique. Not repeated.
func (l Lookup) Id(v v3.Vec) uint32 {
	// Deduplicate by removing small details and use of epsilon
	epsilon := float64(0.0001)
	key := [3]int32{int32((v.X + epsilon) * 1000), int32((v.Y + epsilon) * 1000), int32((v.Z + epsilon) * 1000)}
	if vID, ok := l.hashTable[key]; ok {
		// Vertex already exists. It's repeated.
		return vID
	}

	// Vertex is new, so append it.
	*l.V = append(*l.V, v)

	// Store index of the appended vertex.
	l.hashTable[key] = uint32(l.vertexCount() - 1)

	// Return index of the appended vertex.
	return uint32(l.vertexCount() - 1)
}

func (l *Lookup) vertexCount() int {
	return len(*l.V)
}

// To be called after adding all vertices to the vertex buffer.
func (l *Lookup) Destroy() {
	// Clear memory.
	l.hashTable = nil
	runtime.GC()
}
