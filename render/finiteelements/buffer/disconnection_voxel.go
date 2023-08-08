package buffer

// Count separate components consisting of disconnected voxels.
func (vg *VoxelGrid) ComponentsByVoxel() int {
	// Map key is (x, y, z) index of voxel.
	visited := make(map[[3]int]bool)
	count := 0
	process := func(x, y, z int, els []*Element) {
		// If voxel is not empty and if it's not already visited.
		if len(els) > 0 && !visited[[3]int{x, y, z}] {
			count++
			vg.bfsByVoxel(visited, [3]int{x, y, z})
		}
	}
	vg.Iterate(process)
	return count
}

// Algorithm: breadth-first search (BFS).
func (vg *VoxelGrid) bfsByVoxel(visited map[[3]int]bool, start [3]int) {
	queue := [][3]int{start}
	visited[start] = true

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		neighbors := vg.nonemptyNeighbors(v)

		for _, n := range neighbors {
			if !visited[n] {
				visited[n] = true
				queue = append(queue, n)
			}
		}
	}
}

// It returns a list of neighbor voxels that are non-empty.
func (vg *VoxelGrid) nonemptyNeighbors(v [3]int) [][3]int {
	var neighbors [][3]int

	// The 3D directions to iterate over.
	// Two voxels are considered neighbors if they share a face,
	// i.e., they are adjacent along x, y, or z direction.
	directions := [][3]int{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}

	for _, direction := range directions {
		x := v[0] + direction[0]
		y := v[1] + direction[1]
		z := v[2] + direction[2]

		if !vg.isValid(x, y, z) {
			continue
		}

		// Is neighbor voxel non-empty?
		if len(vg.Get(x, y, z)) > 0 {
			neighbors = append(neighbors, [3]int{x, y, z})
		}
	}

	return neighbors
}

//-----------------------------------------------------------------------------
