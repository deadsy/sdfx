package render

import "github.com/deadsy/sdfx/render/buffer"

// To write different types of finite elements as ABAQUS or CalculiX `inp` file.
type Inp struct {
	// To write only required nodes to `inp` file.
	TempVBuff *buffer.VB
	// Output `inp` file path.
	Path string
}

func NewInp(path string) *Inp {
	return &Inp{
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
