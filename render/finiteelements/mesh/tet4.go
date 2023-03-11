package mesh

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/buffer"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// A mesh of 4-node tetrahedra.
// A sophisticated data structure for mesh is required.
// The repeated nodes would be removed.
// The element connectivity would be created with unique nodes.
type MeshTet4 struct {
	// Index buffer.
	IBuff *buffer.Tet4IB
	// Vertex buffer.
	VBuff *buffer.VB
}

// To get a new mesh and number of its layers along Z-axis.
func NewMeshTet4(s sdf.SDF3, r render.RenderTet4) (*MeshTet4, int) {
	fes := render.ToTet4(s, r)

	_, _, layerCountZ := r.LayerCounts(s)

	m := newMeshTet4(layerCountZ)

	// Fill out the mesh with finite elements.
	for _, fe := range fes {
		m.addFE(fe.Layer, fe.V[0], fe.V[1], fe.V[2], fe.V[3])
	}

	defer m.VBuff.DestroyHashTable()

	return m, layerCountZ
}

func newMeshTet4(layerCount int) *MeshTet4 {
	return &MeshTet4{
		IBuff: buffer.NewTet4IB(layerCount),
		VBuff: buffer.NewVB(),
	}
}

func (m *MeshTet4) NodesPerElement() int {
	return 4
}

// Add a finite element.
// Layer number and nodes are input.
// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (m *MeshTet4) addFE(l int, a, b, c, d v3.Vec) {
	m.IBuff.AddTet4(l, m.addVertex(a), m.addVertex(b), m.addVertex(c), m.addVertex(d))
}

func (m *MeshTet4) addVertex(vert v3.Vec) uint32 {
	return m.VBuff.Id(vert)
}

func (m *MeshTet4) vertexCount() int {
	return m.VBuff.VertexCount()
}

func (m *MeshTet4) vertex(i uint32) v3.Vec {
	return m.VBuff.Vertex(i)
}

// Number of layers along the Z axis.
func (m *MeshTet4) layerCount() int {
	return m.IBuff.LayerCount()
}

// Number of tetrahedra on a layer.
func (m *MeshTet4) feCountOnLayer(l int) int {
	return m.IBuff.FECountOnLayer(l)
}

// Number of tetrahedra for all layers.
func (m *MeshTet4) feCount() int {
	return m.IBuff.FECount()
}

// Get a finite element.
// Layer number is input.
// Tetrahedron index on layer is input.
// Tetrahedron index could be from 0 to number of tetrahedra on layer.
// Don't return error to increase performance.
func (m *MeshTet4) feIndicies(l, i int) (uint32, uint32, uint32, uint32) {
	return m.IBuff.FEIndicies(l, i)
}

// Get a finite element.
// Layer number is input.
// Tetrahedron index on layer is input.
// Tetrahedron index could be from 0 to number of tetrahedra on layer.
// Don't return error to increase performance.
func (m *MeshTet4) feVertices(l, i int) []v3.Vec {
	idx0, idx1, idx2, idx3 := m.IBuff.FEIndicies(l, i)
	vertices := make([]v3.Vec, 4)

	vertices[0] = m.VBuff.Vertex(idx0)
	vertices[1] = m.VBuff.Vertex(idx1)
	vertices[2] = m.VBuff.Vertex(idx2)
	vertices[3] = m.VBuff.Vertex(idx3)

	return vertices
}

// Write mesh to ABAQUS or CalculiX `inp` file.
func (m *MeshTet4) WriteInp(
	path string,
	layersFixed []int,
	massDensity float32,
	youngModulus float32,
	poissonRatio float32,
) error {
	return m.WriteInpLayers(path, 0, m.layerCount(), layersFixed, massDensity, youngModulus, poissonRatio)
}

// Write specific layers of mesh to ABAQUS or CalculiX `inp` file.
// Result would include start layer.
// Result would exclude end layer.
func (m *MeshTet4) WriteInpLayers(
	path string,
	layerStart, layerEnd int,
	layersFixed []int,
	massDensity float32,
	youngModulus float32,
	poissonRatio float32,
) error {
	inp := NewInp(m, path, layerStart, layerEnd, layersFixed, massDensity, youngModulus, poissonRatio)
	return inp.Write()
}

//-----------------------------------------------------------------------------
