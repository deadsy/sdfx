package mesh

import "github.com/deadsy/sdfx/render/finiteelements/buffer"

// Fix any disconnected finite elements.
// They cause FEA solver to throw error.
func (m *Fem) Components() int {
	return buffer.ConnectedComponents(m.IBuff.Grid)
}
