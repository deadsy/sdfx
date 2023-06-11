package buffer

import "fmt"

// Index buffer for a mesh of finite elements.
type IB struct {
	// Every NodesPerElement indices would correspond to a finite element.
	// It's kept low-level for performance.
	// Finite elements are stored by their layer on Z axis.
	I               [][]uint32
	NodesPerElement int
}

func NewIB(layerCount, nodesPerElement int) *IB {
	if nodesPerElement < 1 {
		panic("nodes per finite element must be positive")
	}

	ib := IB{
		I:               [][]uint32{},
		NodesPerElement: nodesPerElement,
	}

	// Initialize.
	ib.I = make([][]uint32, layerCount)
	for l := 0; l < layerCount; l++ {
		ib.I[l] = make([]uint32, 0)
	}

	return &ib
}

func hasRepeatedValues(slice []uint32) bool {
	valueMap := make(map[uint32]bool)

	for _, value := range slice {
		if valueMap[value] {
			return true
		}
		valueMap[value] = true
	}

	return false
}

// Add a finite element to buffer.
// Layer number and nodes are input.
// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (ib *IB) AddFE(l int, nodes []uint32) {
	if len(nodes) != ib.NodesPerElement {
		// Don't return error.
		// Since this function is going to be called by heavy loops.
		// More efficient this way. Right?
		panic("bad sizes: nodes of finite element")
	}
	ib.I[l] = append(ib.I[l], nodes...)
	if hasRepeatedValues(nodes) {
		fmt.Println("Bad element?")
	}
}

// Number of layers along the Z axis.
func (ib *IB) LayerCount() int {
	return len(ib.I)
}

// Number of finite elements on a layer.
func (ib *IB) FECountOnLayer(l int) int {
	return len(ib.I[l]) / ib.NodesPerElement
}

// Number of finite elements for all layers.
func (ib *IB) FECount() int {
	var count int
	for _, l := range ib.I {
		count += len(l) / ib.NodesPerElement
	}
	return count
}

// Layer number is input.
// FE index on layer is input.
// FE index could be from 0 to number of FE on layer.
// Don't return error to increase performance.
func (ib *IB) FEIndicies(l, i int) []uint32 {
	indices := make([]uint32, ib.NodesPerElement)
	for n := 0; n < ib.NodesPerElement; n++ {
		indices[n] = ib.I[l][i*ib.NodesPerElement+n]
	}
	return indices
}
