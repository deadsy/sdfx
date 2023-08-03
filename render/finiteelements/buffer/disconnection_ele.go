package buffer

// Count separate components consisting of disconnected finite elements.
// They cause FEA solver to throw error.
func (vg *VoxelGrid) Components() int {
	// Map key is (x, y, z) index of voxel.
	visited := make(map[*Element]bool)
	count := 0
	for z := 0; z < vg.Len.Z; z++ {
		for y := 0; y < vg.Len.Y; y++ {
			for x := 0; x < vg.Len.X; x++ {
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

	// The 3D directions to iterate over.
	//
	// Two voxels could be considered neighbors if:
	// 1) They share a face.
	// 2) They share an edge.
	// 3) They share a corner.
	//
	// You'll need to adjust:
	directions := [][3]int{
		// Share a face:
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
		// Share an edge:
		{1, 1, 0},
		{1, -1, 0},
		{-1, 1, 0},
		{-1, -1, 0},
		{0, 1, 1},
		{0, 1, -1},
		{0, -1, 1},
		{0, -1, -1},
		// Share a corner:
		{1, 1, 1},
		{1, -1, -1},
		{-1, 1, -1},
		{-1, -1, 1},
	}

	for _, direction := range directions {
		x := v[0] + direction[0]
		y := v[1] + direction[1]
		z := v[2] + direction[2]

		if !vg.isValid(x, y, z) {
			continue
		}

		for _, el := range vg.Get(x, y, z) {
			if sharesNode(e, el) {
				neighbors = append(neighbors, el)
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
