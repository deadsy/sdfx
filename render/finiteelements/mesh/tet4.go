package mesh

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/buffer"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Tet4 is a mesh of 4-node tetrahedra.
// A sophisticated data structure for mesh is required.
// The repeated nodes would be removed.
// The element connectivity would be created with unique nodes.
type Tet4 struct {
	// Index buffer.
	IBuff *buffer.IB
	// Vertex buffer.
	VBuff *buffer.VB
}

// NewTet4 returns a new mesh and number of its layers along Z-axis.
func NewTet4(s sdf.SDF3, r render.RenderFE) (*Tet4, int) {
	fes := render.ToTet4(s, r)

	_, _, layerCountZ := r.LayerCounts(s)

	m := newTet4(layerCountZ)

	// Fill out the mesh with finite elements.
	for _, fe := range fes {
		m.addFE(fe.Layer, [4]v3.Vec{fe.V[0], fe.V[1], fe.V[2], fe.V[3]})
	}

	defer m.VBuff.DestroyHashTable()

	return m, layerCountZ
}

func newTet4(layerCount int) *Tet4 {
	return &Tet4{
		IBuff: buffer.NewIB(layerCount, 4),
		VBuff: buffer.NewVB(),
	}
}

// NodesPerElement returns nodes per element.
func (m *Tet4) NodesPerElement() int {
	return 4
}

// Add a finite element.
// Layer number and nodes are input.
// The node numbering should follow the convention of CalculiX.
// http://www.dhondt.de/ccx_2.20.pdf
func (m *Tet4) addFE(l int, nodes [4]v3.Vec) {
	indices := [4]uint32{}
	for n := 0; n < 4; n++ {
		indices[n] = m.addVertex(nodes[n])
	}
	m.IBuff.AddFE(l, indices[:])
}

func (m *Tet4) addVertex(vert v3.Vec) uint32 {
	return m.VBuff.Id(vert)
}

func (m *Tet4) vertexCount() int {
	return m.VBuff.VertexCount()
}

func (m *Tet4) vertex(i uint32) v3.Vec {
	return m.VBuff.Vertex(i)
}

// Number of layers along the Z axis.
func (m *Tet4) layerCount() int {
	return m.IBuff.LayerCount()
}

// Number of tetrahedra on a layer.
func (m *Tet4) feCountOnLayer(l int) int {
	return m.IBuff.FECountOnLayer(l)
}

// Number of tetrahedra for all layers.
func (m *Tet4) feCount() int {
	return m.IBuff.FECount()
}

// Get a finite element.
// Layer number is input.
// Tetrahedron index on layer is input.
// Tetrahedron index could be from 0 to number of tetrahedra on layer.
// Don't return error to increase performance.
func (m *Tet4) feIndicies(l, i int) []uint32 {
	return m.IBuff.FEIndicies(l, i)
}

// Get a finite element.
// Layer number is input.
// Tetrahedron index on layer is input.
// Tetrahedron index could be from 0 to number of tetrahedra on layer.
// Don't return error to increase performance.
func (m *Tet4) feVertices(l, i int) []v3.Vec {
	indices := m.IBuff.FEIndicies(l, i)
	vertices := make([]v3.Vec, 4)
	for n := 0; n < 4; n++ {
		vertices[n] = m.VBuff.Vertex(indices[n])
	}
	return vertices
}

// WriteInp writes mesh to ABAQUS or CalculiX `inp` file.
func (m *Tet4) WriteInp(
	path string,
	layersFixed []int,
	massDensity float32,
	youngModulus float32,
	poissonRatio float32,
) error {
	return m.WriteInpLayers(path, 0, m.layerCount(), layersFixed, massDensity, youngModulus, poissonRatio)
}

// WriteInpLayers writes specific layers of mesh to ABAQUS or CalculiX `inp` file.
// Result would include start layer.
// Result would exclude end layer.
func (m *Tet4) WriteInpLayers(
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
