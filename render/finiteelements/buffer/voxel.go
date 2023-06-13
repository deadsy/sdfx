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
	data             [][]*Element // Each voxel stores multiple elements.
	lenX, lenY, lenZ int          // Voxels count in 3 directions.
}

func NewVoxelGrid(x, y, z int) *VoxelGrid {
	return &VoxelGrid{
		data: make([][]*Element, x*y*z),
		lenX: x,
		lenY: y,
		lenZ: z,
	}
}

// To get all the elements inside a voxel.
func (vg *VoxelGrid) Get(x, y, z int) []*Element {
	return vg.data[x*vg.lenY*vg.lenZ+y*vg.lenZ+z]
}

// To set all the elements inside a voxel at once.
func (vg *VoxelGrid) Set(x, y, z int, value []*Element) {
	vg.data[x*vg.lenY*vg.lenZ+y*vg.lenZ+z] = value
}

// To append a single element to the elements inside a voxel.
func (vg *VoxelGrid) Append(x, y, z int, value *Element) {
	vg.data[x*vg.lenY*vg.lenZ+y*vg.lenZ+z] = append(vg.data[x*vg.lenY*vg.lenZ+y*vg.lenZ+z], value)
}

// To iterate over all voxels and get elements inside each voxel and do stuff with them.
func (t *VoxelGrid) Iterate(f func(x, y, z int, value []*Element)) {
	for z := 0; z < t.lenZ; z++ {
		for y := 0; y < t.lenY; y++ {
			for x := 0; x < t.lenX; x++ {
				value := t.Get(x, y, z)
				f(x, y, z, value)
			}
		}
	}
}
