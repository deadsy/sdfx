package render

import "github.com/deadsy/sdfx/render/buffer"

type FEType int

const (
	FETet4 FEType = iota + 1
	FEHex8
	FEHex20
)

// To write different types of finite elements as ABAQUS or CalculiX `inp` file.
type Inp struct {
	// To be able to write to `inp` for all types of FE.
	ElType FEType
	// To write only required nodes to `inp` file.
	TempVBuff *buffer.VB
	// Output `inp` file path.
	Path string
}

func NewInp(elType FEType, path string) *Inp {
	return &Inp{
		ElType:    elType,
		TempVBuff: buffer.NewVB(),
		Path:      path,
	}
}

func (inp *Inp) WriteTet4(m *MeshTet4) error {
	return nil
}

func (inp *Inp) WriteHex8(m *MeshHex8) error {
	return nil
}
