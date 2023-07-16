package buffer

import (
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

// Index buffer for a mesh of finite elements.
type IB struct {
	Grid *VoxelGrid
}

func NewIB(voxelsX, voxelsY, voxelsZ int, mins, maxs []v3.Vec) *IB {
	ib := IB{
		Grid: NewVoxelGrid(voxelsX, voxelsY, voxelsZ, mins, maxs),
	}

	return &ib
}

// Add a finite element to buffer.
// Voxel coordinate and nodes are input.
// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (ib *IB) AddFE(x, y, z int, nodes []uint32) {
	ib.Grid.Append(x, y, z, NewElement(nodes))
}

// To iterate over all voxels and get elements inside each voxel and do stuff with them.
func (ib *IB) Iterate(f func(int, int, int, []*Element)) {
	ib.Grid.Iterate(f)
}

func (ib *IB) Size() (int, int, int) {
	return ib.Grid.Size()
}

// The closest node is identified.
// Also, the containing voxel is figured out.
func (ib *IB) Locate(location v3.Vec) (int, v3i.Vec) {
	return ib.Grid.Locate(location)
}
