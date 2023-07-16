package buffer

import v3 "github.com/deadsy/sdfx/vec/v3"

type Element struct {
	Nodes []uint32 // Node indices
}

// Declare the enum using iota and const
type ElementType int

const (
	C3D4 ElementType = iota + 1
	C3D10
	C3D8
	C3D20R
	Unknown
)

func (e *Element) Type() ElementType {
	if len(e.Nodes) == 4 {
		return C3D4
	} else if len(e.Nodes) == 10 {
		return C3D10
	} else if len(e.Nodes) == 8 {
		return C3D8
	} else if len(e.Nodes) == 20 {
		return C3D20R
	}
	return Unknown
}

func NewElement(nodes []uint32) *Element {
	e := Element{
		Nodes: nodes,
	}
	return &e
}

type Voxel struct {
	data []*Element // Each voxel stores multiple elements.
	min  v3.Vec     // Min corner of voxel.
	max  v3.Vec     // Max corner of voxel.
}

// Acts like a three-dimensional nested slice using
// a one-dimensional slice under the hood.
// To increase performance.
type VoxelGrid struct {
	voxels           []*Voxel //
	lenX, lenY, lenZ int      // Voxels count in 3 directions.
}

func NewVoxelGrid(x, y, z int) *VoxelGrid {
	return &VoxelGrid{
		voxels: make([]*Voxel, x*y*z),
		lenX:   x,
		lenY:   y,
		lenZ:   z,
	}
}

func (vg *VoxelGrid) Size() (int, int, int) {
	return vg.lenX, vg.lenY, vg.lenZ
}

// To get all the elements inside a voxel.
func (vg *VoxelGrid) Get(x, y, z int) []*Element {
	return vg.voxels[x*vg.lenY*vg.lenZ+y*vg.lenZ+z].data
}

// To set all the elements inside a voxel at once.
func (vg *VoxelGrid) Set(x, y, z int, value []*Element) {
	vg.voxels[x*vg.lenY*vg.lenZ+y*vg.lenZ+z].data = value
}

// To append a single element to the elements inside a voxel.
func (vg *VoxelGrid) Append(x, y, z int, value *Element) {
	vg.voxels[x*vg.lenY*vg.lenZ+y*vg.lenZ+z].data = append(vg.voxels[x*vg.lenY*vg.lenZ+y*vg.lenZ+z].data, value)
}

// To iterate over all voxels and get elements inside each voxel and do stuff with them.
func (vg *VoxelGrid) Iterate(f func(int, int, int, []*Element)) {
	for z := 0; z < vg.lenZ; z++ {
		for y := 0; y < vg.lenY; y++ {
			for x := 0; x < vg.lenX; x++ {
				value := vg.Get(x, y, z)
				f(x, y, z, value)
			}
		}
	}
}
