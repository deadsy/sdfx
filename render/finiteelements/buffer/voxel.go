package buffer

import (
	"math"

	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

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

func NewVoxel(min, max v3.Vec) *Voxel {
	return &Voxel{
		data: make([]*Element, 0),
		min:  min,
		max:  max,
	}
}

// Acts like a three-dimensional nested slice using
// a one-dimensional slice under the hood.
// To increase performance.
type VoxelGrid struct {
	voxels           []*Voxel //
	lenX, lenY, lenZ int      // Voxels count in 3 directions.
}

func NewVoxelGrid(x, y, z int, mins, maxs []v3.Vec) *VoxelGrid {
	vg := &VoxelGrid{
		voxels: make([]*Voxel, x*y*z),
		lenX:   x,
		lenY:   y,
		lenZ:   z,
	}

	// Assign the min corner and max corner of each voxel.
	for i := range vg.voxels {
		vg.voxels[i] = NewVoxel(mins[i], maxs[i])
	}

	return vg
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

// The closest node is identified.
// Also, the containing voxel is figured out.
func (vg *VoxelGrid) Locate(location v3.Vec) (int, v3i.Vec) {
	// Calculating voxel indices.
	// Assumes that the voxels are evenly distributed across the grid.
	idxX := int((location.X - vg.voxels[0].min.X) / (vg.voxels[0].max.X - vg.voxels[0].min.X) * float64(vg.lenX))
	idxY := int((location.Y - vg.voxels[0].min.Y) / (vg.voxels[0].max.Y - vg.voxels[0].min.Y) * float64(vg.lenY))
	idxZ := int((location.Z - vg.voxels[0].min.Z) / (vg.voxels[0].max.Z - vg.voxels[0].min.Z) * float64(vg.lenZ))

	// Ensure indices are within bounds
	if idxX >= vg.lenX {
		idxX = vg.lenX - 1
	}
	if idxY >= vg.lenY {
		idxY = vg.lenY - 1
	}
	if idxZ >= vg.lenZ {
		idxZ = vg.lenZ - 1
	}

	// Get elements in the voxel
	elements := vg.Get(idxX, idxY, idxZ)

	// Find the closest node
	closestNode := -1
	minDistance := math.Inf(1)

	for _, element := range elements {
		for _, node := range element.Nodes {
			// Assuming you have a function that gives you the position of a node
			nodePos := GetNodePosition(node)

			distance := location.Sub(nodePos).Length()
			if distance < minDistance {
				minDistance = distance
				closestNode = int(node)
			}
		}
	}

	return closestNode, v3i.Vec{X: idxX, Y: idxY, Z: idxZ}
}
