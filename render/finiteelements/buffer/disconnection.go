package buffer

// Count separate components consisting of disconnected finite elements.
// They cause FEA solver to throw error.
func (vg *VoxelGrid) Components() int {
	visited := make(map[*Element]bool)
	count := 0
	process := func(x, y, z int, els []*Element) {
		for _, el := range els {
			if !visited[el] {
				count++
				vg.bfs(visited, el, [3]int{x, y, z})
			}
		}
	}
	vg.Iterate(process)
	return count
}

// Algorithm: breadth-first search (bfs).
func (vg *VoxelGrid) bfs(visited map[*Element]bool, start *Element, startV [3]int) {
	queue := []*Element{start}
	quVox := [][3]int{startV} // To store the voxel of each element.
	visited[start] = true

	for len(queue) > 0 {
		e := queue[0]
		v := quVox[0]
		queue = queue[1:]
		quVox = quVox[1:]

		neighbors, neighVoxs := vg.neighbors(e, v)

		for i := range neighbors {
			n := neighbors[i]
			nv := neighVoxs[i]
			if !visited[n] {
				visited[n] = true
				queue = append(queue, n)
				quVox = append(quVox, nv)
			}
		}
	}
}

// It returns a list of neighbors.
func (vg *VoxelGrid) neighbors(e *Element, v [3]int) ([]*Element, [][3]int) {
	var neighbors []*Element
	var neighVoxs [][3]int

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
						neighVoxs = append(neighVoxs, [3]int{x, y, z})
					}
				}
			}
		}
	}

	return neighbors, neighVoxs
}

func sharesNode(e1, e2 *Element) bool {
	// The node count doesn't need to be equal for the two elements.
	// Since, the two elements could be of different types.

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

// Is voxel index inside a valid range?
func (vg *VoxelGrid) isValid(x, y, z int) bool {
	return x >= 0 && y >= 0 && z >= 0 && x < vg.Len.X && y < vg.Len.Y && z < vg.Len.Z
}

//-----------------------------------------------------------------------------
