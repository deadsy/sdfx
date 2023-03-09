package render

import (
	"fmt"

	"github.com/deadsy/sdfx/render/buffer"
)

// To write different types of finite elements as ABAQUS or CalculiX `inp` file.
type Inp struct {
	// Output `inp` file path.
	Path string
	// Output `inp` file would include start layer.
	LayerStart int
	// Output `inp` file would exclude end layer.
	LayerEnd int
	// To write only required nodes to `inp` file.
	TempVBuff *buffer.VB
	// Mechanical properties of 3D print resin.
	MassDensity  float32
	YoungModulus float32
	PoissonRatio float32
}

func NewInp(path string, layerStart, layerEnd int, massDensity float32, youngModulus float32, poissonRatio float32) *Inp {
	return &Inp{
		Path:         path,
		LayerStart:   layerStart,
		LayerEnd:     layerEnd,
		TempVBuff:    buffer.NewVB(),
		MassDensity:  massDensity,
		YoungModulus: youngModulus,
		PoissonRatio: poissonRatio,
	}
}

func (inp *Inp) WriteTet4(m *MeshTet4) error {
	inp.writeHeader()
	return nil
}

func (inp *Inp) WriteHex8(m *MeshHex8) error {
	inp.writeHeader()
	return nil
}

func (inp *Inp) WriteHex20() error {
	inp.writeHeader()
	return nil
}

func (inp *Inp) writeHeader() error {
	if 0 <= inp.LayerStart && inp.LayerStart < inp.LayerEnd && inp.LayerEnd <= m.layerCount() {
		// Good.
	} else {
		return fmt.Errorf("start or end layer is beyond range")
	}
}
