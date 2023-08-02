package buffer

type Graph struct {
	adjacencyList map[*Voxel][]*Voxel
}

func (g *Graph) AddEdge(v1, v2 *Voxel) {
	// Add an edge from v1 to v2. A new element is created if v1 is not already in the map
	g.adjacencyList[v1] = append(g.adjacencyList[v1], v2)
	// Add an edge from v2 to v1. A new element is created if v2 is not already in the map
	g.adjacencyList[v2] = append(g.adjacencyList[v2], v1)
}

func NewGraph(vg *VoxelGrid) *Graph {
	g := &Graph{
		adjacencyList: make(map[*Voxel][]*Voxel),
	}

	// For each pair of voxels, check if they share a node and add an edge if they do
	for i, v1 := range vg.Voxels {
		for _, v2 := range vg.Voxels[i+1:] {
			for _, e1 := range v1.data {
				for _, e2 := range v2.data {
					if sharesNode(e1, e2) {
						g.AddEdge(v1, v2)
						break
					}
				}
			}
		}
	}

	return g
}

// Check if two elements share a node
func sharesNode(e1, e2 *Element) bool {
	for _, n1 := range e1.Nodes {
		for _, n2 := range e2.Nodes {
			if n1 == n2 {
				return true
			}
		}
	}

	return false
}

// Perform a DFS on the voxel graph starting from voxel v
func DFS(g *Graph, v *Voxel, visited *map[*Voxel]bool) {
	(*visited)[v] = true

	for _, neighbor := range g.adjacencyList[v] {
		if !(*visited)[neighbor] {
			DFS(g, neighbor, visited)
		}
	}
}

// Returns the number of connected components in the voxel grid
func ConnectedComponents(vg *VoxelGrid) int {
	g := NewGraph(vg)
	visited := make(map[*Voxel]bool)
	components := 0

	for _, v := range vg.Voxels {
		if !visited[v] {
			components++
			DFS(g, v, &visited)
		}
	}

	return components
}
