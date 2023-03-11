package buffer

// Index buffer for 8-node hexahedra.
type Hex8IB struct {
	// Every 8 indices would correspond to a hexahedron.
	// It's kept low-level for performance.
	// Tetrahedra are stored by their layer on Z axis.
	I [][]uint32
}

func NewHex8IB(layerCount int) *Hex8IB {
	ib := Hex8IB{
		I: [][]uint32{},
	}

	// Initialize.
	ib.I = make([][]uint32, layerCount)
	for l := 0; l < layerCount; l++ {
		ib.I[l] = make([]uint32, 0)
	}

	return &ib
}

// Add a finite element to buffer.
// Layer number and nodes are input.
// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (ib *Hex8IB) AddFE(l int, nodes [8]uint32) {
	ib.I[l] = append(ib.I[l], nodes[:]...)
}

// Number of layers along the Z axis.
func (ib *Hex8IB) LayerCount() int {
	return len(ib.I)
}

// Number of finite elements on a layer.
func (ib *Hex8IB) FECountOnLayer(l int) int {
	return len(ib.I[l]) / 8
}

// Number of finite elements for all layers.
func (ib *Hex8IB) FECount() int {
	var count int
	for _, l := range ib.I {
		count += len(l) / 8
	}
	return count
}

// Layer number is input.
// FE index on layer is input.
// FE index could be from 0 to number of FE on layer.
// Don't return error to increase performance.
func (ib *Hex8IB) FEIndicies(l, i int) [8]uint32 {
	indices := [8]uint32{}
	for n := 0; n < 8; n++ {
		indices[n] = ib.I[l][i*8+n]
	}
	return indices
}
