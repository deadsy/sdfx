package buffer

import (
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
	Voxels []*Voxel //
	Len    v3i.Vec  // Voxel count in 3 directions.
	Dim    v3.Vec   // Voxel dimension in 3 directions.
}

func NewVoxelGrid(len v3i.Vec, dim v3.Vec, mins, maxs []v3.Vec) *VoxelGrid {
	vg := &VoxelGrid{
		Voxels: make([]*Voxel, len.X*len.Y*len.Z),
		Len:    len,
		Dim:    dim,
	}

	// Assign the min corner and max corner of each voxel.
	for i := range vg.Voxels {
		vg.Voxels[i] = NewVoxel(mins[i], maxs[i])
	}

	return vg
}

// This func must be consistent with `(r *MarchingCubesFEUniform) Voxels` func.
// This func must be consistent with `marchingCubesFE` func too.
func (vg *VoxelGrid) index1Dto3D(i int) (int, int, int) {
	z := i % vg.Len.Z
	y := (i / vg.Len.Z) % vg.Len.Y
	x := i / (vg.Len.Z * vg.Len.Y)
	return x, y, z
}

// This func must be consistent with `(r *MarchingCubesFEUniform) Voxels` func.
// This func must be consistent with `marchingCubesFE` func too.
func (vg *VoxelGrid) index3Dto1D(x, y, z int) int {
	return x*vg.Len.Y*vg.Len.Z + y*vg.Len.Z + z
}

func (vg *VoxelGrid) Size() (int, int, int) {
	return vg.Len.X, vg.Len.Y, vg.Len.Z
}

// To get all the elements inside a voxel.
func (vg *VoxelGrid) Get(x, y, z int) []*Element {
	return vg.Voxels[vg.index3Dto1D(x, y, z)].data
}

// To set all the elements inside a voxel at once.
func (vg *VoxelGrid) Set(x, y, z int, value []*Element) {
	vg.Voxels[vg.index3Dto1D(x, y, z)].data = value
}

// To append a single element to the elements inside a voxel.
func (vg *VoxelGrid) Append(x, y, z int, value *Element) {
	vg.Voxels[vg.index3Dto1D(x, y, z)].data = append(vg.Voxels[vg.index3Dto1D(x, y, z)].data, value)
}

// Compute the bounding box of all the input points.
// Return all the voxels that are intersecting with that bounding box.
func (vg *VoxelGrid) VoxelsIntersecting(points []v3.Vec) ([]v3i.Vec, v3.Vec, v3.Vec) {
	if len(points) == 0 {
		return nil, v3.Vec{}, v3.Vec{}
	}

	// compute the bounding box of all the input points
	min, max := points[0], points[0]
	for _, point := range points {
		min = min.Min(point)
		max = max.Max(point)
	}

	var intersectingVoxels []v3i.Vec

	// iterate over all the voxels
	for i, voxel := range vg.Voxels {
		// check if the voxel intersects with the bounding box
		if !DoesIntersect(voxel.Min, voxel.Max, min, max) {
			continue
		}

		// convert the 1D index to a 3D index
		x, y, z := vg.index1Dto3D(i)

		intersectingVoxels = append(intersectingVoxels, v3i.Vec{X: x, Y: y, Z: z})
	}

	return intersectingVoxels, min, max
}

// Does two voxels or b-boxes intersect with each other?
func DoesIntersect(aMin, aMax, bMin, bMax v3.Vec) bool {
	if aMin.X > bMax.X || aMin.Y > bMax.Y || aMin.Z > bMax.Z {
		return false
	}
	if bMin.X > aMax.X || bMin.Y > aMax.Y || bMin.Z > aMax.Z {
		return false
	}
	return true
}

// To iterate over all voxels and get elements inside each voxel and do stuff with them.
func (vg *VoxelGrid) Iterate(f func(int, int, int, []*Element)) {
	for z := 0; z < vg.Len.Z; z++ {
		for y := 0; y < vg.Len.Y; y++ {
			for x := 0; x < vg.Len.X; x++ {
				value := vg.Get(x, y, z)
				f(x, y, z, value)
			}
		}
	}
}
