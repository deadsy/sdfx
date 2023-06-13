package buffer

type Element struct {
	Nodes []uint32 // Node indices
}

func NewElement(nodes []uint32) *Element {
	e := Element{
		Nodes: nodes,
	}
	return &e
}

// Acts like a three-dimensional nested slice using
// a one-dimensional slice under the hood.
// To increase performance.
type VoxelGrid struct {
	data             [][]Element // Each voxel stores multiple elements.
	xLen, yLen, zLen int         // Voxels count in 3 directions.
}

func NewVoxelGrid(x, y, z int) *VoxelGrid {
	return &VoxelGrid{
		data: make([][]Element, x*y*z),
		xLen: x,
		yLen: y,
		zLen: z,
	}
}

// To get all the elements inside a voxel.
func (vg *VoxelGrid) Get(x, y, z int) []Element {
	return vg.data[x*vg.yLen*vg.zLen+y*vg.zLen+z]
}

// To set all the elements inside a voxel at once.
func (vg *VoxelGrid) Set(x, y, z int, value []Element) {
	vg.data[x*vg.yLen*vg.zLen+y*vg.zLen+z] = value
}

// To append a single element to the elements inside a voxel.
func (vg *VoxelGrid) Append(x, y, z int, value Element) {
	vg.data[x*vg.yLen*vg.zLen+y*vg.zLen+z] = append(vg.data[x*vg.yLen*vg.zLen+y*vg.zLen+z], value)
}
