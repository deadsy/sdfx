package buffer

// Count separate components consisting of disconnected finite elements.
// They cause FEA solver to throw error.
func (vg *VoxelGrid) Components() int {
	// Map key is (x, y, z) index of voxel.
	visited := make(map[*Element]bool)
	count := 0
	for x := 0; x < vg.Len.X; x++ {
		for y := 0; y < vg.Len.Y; y++ {
			for z := 0; z < vg.Len.Z; z++ {
				els := vg.Get(x, y, z)
				for _, el := range els {
					if !visited[el] {
						count++
						vg.BFS(visited, el, [3]int{x, y, z})
					}
				}
			}
		}
	}
	return count
}

// Algorithm: breadth-first search (BFS).
func (vg *VoxelGrid) BFS(visited map[*Element]bool, start *Element, v [3]int) {
	queue := []*Element{start}
	visited[start] = true

	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]

		neighbors := vg.neighbors(e, v)

		for _, n := range neighbors {
			if !visited[n] {
				visited[n] = true
				queue = append(queue, n)
			}
		}
	}
}

// It returns a list of neighbors.
func (vg *VoxelGrid) neighbors(e *Element, v [3]int) []*Element {
	var neighbors []*Element

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			for k := -1; k <= 1; k++ {
				x := v[0] + i
				y := v[1] + j
				z := v[2] + k

				if !vg.isValid(x, y, z) {
					continue
				}

				for _, el := range vg.Get(x, y, z) {
					if i == 0 && j == 0 && k == 0 {
						// The same voxel: skip the same element.
						if el == e {
							continue
						}
					}
					if sharesNode(e, el) {
						neighbors = append(neighbors, el)
					}
				}
			}
		}
	}

	return neighbors
}

func sharesNode(e1, e2 *Element) bool {
	if len(e1.Nodes) != len(e2.Nodes) {
		return false
	}
	for _, n1 := range e1.Nodes {
		if contains(e2.Nodes, n1) {
			return true
		}
	}
	return false
}

func contains(arr []uint32, i uint32) bool {
	for _, n := range arr {
		if n == i {
			return true
		}
	}
	return false
}

//-----------------------------------------------------------------------------
