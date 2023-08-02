package mesh

import (
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/buffer"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

// Fem is a mesh of finite elements.
// A sophisticated data structure for mesh is required.
// The repeated nodes would be removed.
// The element connectivity would be created with unique nodes.
type Fem struct {
	// Index buffer.
	IBuff *buffer.IB
	// Vertex buffer.
	VBuff *buffer.VB
}

// NewFem returns a new mesh and number of its layers along Z-axis.
func NewFem(s sdf.SDF3, r render.RenderFE) (*Fem, int) {
	fes := render.ToFem(s, r)

	voxelLen, voxelDim, mins, maxs := r.Voxels(s)

	m := newFem(voxelLen, voxelDim, mins, maxs)

	// Fill out the mesh with finite elements.
	for _, fe := range fes {
		m.addFE(fe.X, fe.Y, fe.Z, fe.V)
	}

	defer m.VBuff.DestroyHashTable()

	return m, voxelLen.Z
}

func newFem(voxelLen v3i.Vec, voxelDim v3.Vec, mins, maxs []v3.Vec) *Fem {
	return &Fem{
		IBuff: buffer.NewIB(voxelLen, voxelDim, mins, maxs),
		VBuff: buffer.NewVB(),
	}
}

func (m *Fem) Size() (int, int, int) {
	return m.IBuff.Size()
}

// Add a finite element.
// Voxel coordinate and nodes are input.
// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (m *Fem) addFE(x, y, z int, nodes []v3.Vec) {
	indices := make([]uint32, len(nodes))
	for n := 0; n < len(nodes); n++ {
		indices[n] = m.addVertex(nodes[n])
	}
	m.IBuff.AddFE(x, y, z, indices)
}

func (m *Fem) addVertex(vert v3.Vec) uint32 {
	return m.VBuff.Id(vert)
}

func (m *Fem) vertexCount() int {
	return m.VBuff.VertexCount()
}

func (m *Fem) vertex(i uint32) v3.Vec {
	return m.VBuff.Vertex(i)
}

// To iterate over all voxels and get elements inside each voxel and do stuff with them.
func (m *Fem) iterate(f func(int, int, int, []*buffer.Element)) {
	m.IBuff.Iterate(f)
}

// The closest node is identified.
// Also, the containing voxel is figured out.
//
// This logic has to be here, since we need access to any node vertex.
func (m *Fem) Locate(location v3.Vec) (uint32, v3i.Vec) {
	// Calculating voxel indices.
	idxX := int(math.Floor((location.X - m.IBuff.Grid.Voxels[0].Min.X) / (m.IBuff.Grid.Dim.X)))
	idxY := int(math.Floor((location.Y - m.IBuff.Grid.Voxels[0].Min.Y) / (m.IBuff.Grid.Dim.Y)))
	idxZ := int(math.Floor((location.Z - m.IBuff.Grid.Voxels[0].Min.Z) / (m.IBuff.Grid.Dim.Z)))

	// Ensure indices are within bounds
	if idxX >= m.IBuff.Grid.Len.X {
		idxX = m.IBuff.Grid.Len.X - 1
	}
	if idxY >= m.IBuff.Grid.Len.Y {
		idxY = m.IBuff.Grid.Len.Y - 1
	}
	if idxZ >= m.IBuff.Grid.Len.Z {
		idxZ = m.IBuff.Grid.Len.Z - 1
	}

	// Get elements in the voxel
	elements := m.IBuff.Grid.Get(idxX, idxY, idxZ)

	// Find the closest node
	var closestNode uint32
	minDistance := math.Inf(1)

	for _, element := range elements {
		for _, node := range element.Nodes {
			// A function that gives you the position of a node.
			nodePos := m.vertex(node)

			distance := location.Sub(nodePos).Length()
			if distance < minDistance {
				minDistance = distance
				closestNode = node
			}
		}
	}

	return closestNode, v3i.Vec{X: idxX, Y: idxY, Z: idxZ}
}

// Compute the bounding box of all the input points.
// Return all the voxels that are intersecting with that bounding box.
func (m *Fem) VoxelsIntersecting(points []v3.Vec) ([]v3i.Vec, v3.Vec, v3.Vec) {
	return m.IBuff.Grid.VoxelsIntersecting(points)
}

//-----------------------------------------------------------------------------

// Count separate components consisting of disconnected finite elements.
// They cause FEA solver to throw error.
func (m *Fem) CountComponents() int {
	// Map key is (x, y, z) index of voxel.
	visited := make(map[[3]int]bool, m.IBuff.Grid.Len.X*m.IBuff.Grid.Len.Y*m.IBuff.Grid.Len.Z)
	count := 0
	process := func(x, y, z int, els []*buffer.Element) {
		if !visited[[3]int{x, y, z}] {
			count++
			m.bfs(visited, [3]int{x, y, z})
		}
	}
	m.iterate(process)
	return count
}

func (m *Fem) bfs(visited map[[3]int]bool, start [3]int) {
	queue := [][3]int{start}
	visited[start] = true

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		neighbors := m.getNeighbors(v)

		for _, n := range neighbors {
			if !visited[n] {
				visited[n] = true
				queue = append(queue, n)
			}
		}
	}
}

// It returns a list of neighbor voxels that are full, i.e. not empty.
func (m *Fem) getNeighbors(v [3]int) [][3]int {
	var neighbors [][3]int
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			for k := -1; k <= 1; k++ {
				if i == 0 && j == 0 && k == 0 {
					continue
				}

				x := v[0] + i
				y := v[1] + j
				z := v[2] + k

				if x >= 0 && x < m.IBuff.Grid.Len.X &&
					y >= 0 && y < m.IBuff.Grid.Len.Y &&
					z >= 0 && z < m.IBuff.Grid.Len.Z {
					// Index is valid.
				} else {
					continue
				}

				// Is neighbor voxel empty or not?
				if len(m.IBuff.Grid.Get(x, y, z)) > 0 {
					neighbors = append(neighbors, [3]int{x, y, z})
				}
			}
		}
	}
	return neighbors
}

//-----------------------------------------------------------------------------

// WriteInp writes mesh to ABAQUS or CalculiX `inp` file.
func (m *Fem) WriteInp(
	path string,
	massDensity float32,
	youngModulus float32,
	poissonRatio float32,
	restraints []*Restraint,
	loads []*Load,
	gravityDirection v3.Vec,
	gravityMagnitude float64,
) error {
	_, _, layersZ := m.IBuff.Size()
	return m.WriteInpLayers(path, 0, layersZ, massDensity, youngModulus, poissonRatio, restraints, loads, gravityDirection, gravityMagnitude)
}

// WriteInpLayers writes specific layers of mesh to ABAQUS or CalculiX `inp` file.
// Result would include start layer.
// Result would exclude end layer.
func (m *Fem) WriteInpLayers(
	path string,
	layerStart, layerEnd int,
	massDensity float32,
	youngModulus float32,
	poissonRatio float32,
	restraints []*Restraint,
	loads []*Load,
	gravityDirection v3.Vec,
	gravityMagnitude float64,
) error {
	inp := NewInp(m, path, layerStart, layerEnd, massDensity, youngModulus, poissonRatio, restraints, loads, gravityDirection, gravityMagnitude)
	return inp.Write()
}

//-----------------------------------------------------------------------------
