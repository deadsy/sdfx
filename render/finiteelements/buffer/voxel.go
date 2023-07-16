package buffer

import (
	v3 "github.com/deadsy/sdfx/vec/v3"
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
	Min  v3.Vec     // Min corner of voxel.
	Max  v3.Vec     // Max corner of voxel.
}

func NewVoxel(min, max v3.Vec) *Voxel {
	return &Voxel{
		data: make([]*Element, 0),
		Min:  min,
		Max:  max,
	}
}

// Acts like a three-dimensional nested slice using
// a one-dimensional slice under the hood.
// To increase performance.
type VoxelGrid struct {
	Voxels           []*Voxel //
	LenX, LenY, LenZ int      // Voxels count in 3 directions.
}

func NewVoxelGrid(x, y, z int, mins, maxs []v3.Vec) *VoxelGrid {
	vg := &VoxelGrid{
		Voxels: make([]*Voxel, x*y*z),
		LenX:   x,
		LenY:   y,
		LenZ:   z,
	}

	// Assign the min corner and max corner of each voxel.
	for i := range vg.Voxels {
		vg.Voxels[i] = NewVoxel(mins[i], maxs[i])
	}

	return vg
}

func (vg *VoxelGrid) Size() (int, int, int) {
	return vg.LenX, vg.LenY, vg.LenZ
}

// To get all the elements inside a voxel.
func (vg *VoxelGrid) Get(x, y, z int) []*Element {
	return vg.Voxels[x*vg.LenY*vg.LenZ+y*vg.LenZ+z].data
}

// To set all the elements inside a voxel at once.
func (vg *VoxelGrid) Set(x, y, z int, value []*Element) {
	vg.Voxels[x*vg.LenY*vg.LenZ+y*vg.LenZ+z].data = value
}

// To append a single element to the elements inside a voxel.
func (vg *VoxelGrid) Append(x, y, z int, value *Element) {
	vg.Voxels[x*vg.LenY*vg.LenZ+y*vg.LenZ+z].data = append(vg.Voxels[x*vg.LenY*vg.LenZ+y*vg.LenZ+z].data, value)
}

// To iterate over all voxels and get elements inside each voxel and do stuff with them.
func (vg *VoxelGrid) Iterate(f func(int, int, int, []*Element)) {
	for z := 0; z < vg.LenZ; z++ {
		for y := 0; y < vg.LenY; y++ {
			for x := 0; x < vg.LenX; x++ {
				value := vg.Get(x, y, z)
				f(x, y, z, value)
			}
		}
	}
}
