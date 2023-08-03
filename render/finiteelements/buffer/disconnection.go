package buffer

// Count separate components consisting of disconnected finite elements.
// They cause FEA solver to throw error.
func (vg *VoxelGrid) CountComponents() int {
	// Map key is (x, y, z) index of voxel.
	visited := make(map[[3]int]bool)
	count := 0
	for z := 0; z < vg.Len.Z; z++ {
		for y := 0; y < vg.Len.Y; y++ {
			for x := 0; x < vg.Len.X; x++ {
				// If voxel is not empty and if it's not already visited.
				if len(vg.Get(x, y, z)) > 0 && !visited[[3]int{x, y, z}] {
					count++
					vg.bfs(visited, [3]int{x, y, z})
				}
			}
		}
	}
	return count
}

// Algorithm: breadth-first search (BFS).
// This is much faster than depth first search (DFS) algorithm.
func (vg *VoxelGrid) bfs(visited map[[3]int]bool, start [3]int) {
	queue := [][3]int{start}
	visited[start] = true

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		neighbors := vg.getNeighbors(v)

		for _, n := range neighbors {
			if !visited[n] {
				visited[n] = true
				queue = append(queue, n)
			}
		}
	}
}

// It returns a list of neighbor voxels that are full, i.e. not empty.
func (vg *VoxelGrid) getNeighbors(v [3]int) [][3]int {
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

				if x >= 0 && x < vg.Len.X &&
					y >= 0 && y < vg.Len.Y &&
					z >= 0 && z < vg.Len.Z {
					// Index is valid.
				} else {
					continue
				}

				// Is neighbor voxel empty or not?
				if len(vg.Get(x, y, z)) > 0 {
					neighbors = append(neighbors, [3]int{x, y, z})
				}
			}
		}
	}
	return neighbors
}

//-----------------------------------------------------------------------------
