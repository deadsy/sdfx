package mesh

// Fix any disconnected finite elements.
// They cause FEA solver to throw error.
func (m *Fem) Connect() {
}
