package buffer

import "fmt"

// Rather than connecting disconnected components, we can keep the largest one and delete the rest.
func (vg *VoxelGrid) CleanDisconnections(components []*Component) {
	// Find largest component. Consider volume criterion.
	maxComponentIndex := -1
	maxVoxelCount := 0
	for i, component := range components {
		if component.VoxelCount() > maxVoxelCount {
			maxVoxelCount = component.VoxelCount()
			maxComponentIndex = i
		}
	}

	if maxComponentIndex != -1 {
		fmt.Printf("Component %v has the largest voxel count: %v\n", maxComponentIndex, maxVoxelCount)
	} else {
		fmt.Printf("No components found.")
	}

	// Remove elements inside the voxels of smaller components.
	for i, component := range components {
		if i != maxComponentIndex {
			for v := range component.Voxels {
				vg.DelAll(v[0], v[1], v[2])
			}
		}
	}
}
