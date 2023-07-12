package mesh

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/finiteelements/buffer"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
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

	layersX, layersY, layersZ := r.LayerCounts(s)

	m := newFem(layersX, layersY, layersZ)

	// Fill out the mesh with finite elements.
	for _, fe := range fes {
		m.addFE(fe.X, fe.Y, fe.Z, fe.V)
	}

	defer m.VBuff.DestroyHashTable()

	return m, layersZ
}

func newFem(layersX, layersY, layersZ int) *Fem {
	return &Fem{
		IBuff: buffer.NewIB(layersX, layersY, layersZ),
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

// WriteInp writes mesh to ABAQUS or CalculiX `inp` file.
func (m *Fem) WriteInp(
	path string,
	massDensity float32,
	youngModulus float32,
	poissonRatio float32,
	restraint func(x, y, z float64) (bool, bool, bool),
	load func(x, y, z float64) (float64, float64, float64),
	gravityDirection v3.Vec,
	gravityMagnitude float64,
) error {
	_, _, layersZ := m.IBuff.Size()
	return m.WriteInpLayers(path, 0, layersZ, massDensity, youngModulus, poissonRatio, restraint, load, gravityDirection, gravityMagnitude)
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
	restraint func(x, y, z float64) (bool, bool, bool),
	load func(x, y, z float64) (float64, float64, float64),
	gravityDirection v3.Vec,
	gravityMagnitude float64,
) error {
	inp := NewInp(m, path, layerStart, layerEnd, massDensity, youngModulus, poissonRatio, restraint, load, gravityDirection, gravityMagnitude)
	return inp.Write()
}

//-----------------------------------------------------------------------------
