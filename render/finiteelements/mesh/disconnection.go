package mesh

import "github.com/deadsy/sdfx/render/finiteelements/buffer"

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
