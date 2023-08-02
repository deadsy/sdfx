package mesh

// Count separate components consisting of disconnected finite elements.
// They cause FEA solver to throw error.
func (m *Fem) CountComponents() int {
	return 0
}
