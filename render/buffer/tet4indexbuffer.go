package buffer

// Index buffer for 4-node tetrahedra.
type Tet4IB struct {
	// Every 4 indices would correspond to a tetrahedron.
	// It's kept low-level for performance.
	// Tetrahedra are stored by their layer on Z axis.
	I [][]uint32
}

func NewTet4IB(layerCount int) *Tet4IB {
	ib := Tet4IB{
		I: [][]uint32{},
	}

	// Initialize.
	ib.I = make([][]uint32, layerCount)
	for l := 0; l < layerCount; l++ {
		ib.I[l] = make([]uint32, 0)
	}

	return &ib
}

// Layer number and 4 nodes are input.
// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (ib *Tet4IB) AddTet4(l int, a, b, c, d uint32) {
	ib.I[l] = append(ib.I[l], a, b, c, d)
}

// Number of layers along the Z axis.
func (ib *Tet4IB) LayerCount() int {
	return len(ib.I)
}

// Number of tetrahedra on a layer.
func (ib *Tet4IB) Tet4CountOnLayer(l int) int {
	return len(ib.I[l]) / 4
}

// Number of tetrahedra for all layers.
func (ib *Tet4IB) Tet4Count() int {
	var count int
	for _, l := range ib.I {
		count += len(l) / 4
	}
	return count
}

// Layer number is input.
// Tetrahedron index on layer is input.
// Tetrahedron index could be from 0 to number of tetrahedra on layer.
// Don't return error to increase performance.
func (ib *Tet4IB) Tet4Indicies(l, i int) (uint32, uint32, uint32, uint32) {
	return ib.I[l][i*4], ib.I[l][i*4+1], ib.I[l][i*4+2], ib.I[l][i*4+3]
}
