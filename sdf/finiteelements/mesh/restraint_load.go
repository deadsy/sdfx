package mesh

import (
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

// A voxel is detected that intersects with the restraint location.
// Boundary is created for all the nodes inside that voxel.
// This way, the stress concentration at the restraint may be alleviated by distributing it.
type Restraint struct {
	Location v3.Vec  // Exact coordinates inside restraint.
	IsFixedX bool    // Is X degree of freedom fixed?
	IsFixedY bool    // Is Y degree of freedom fixed?
	IsFixedZ bool    // Is Z degree of freedom fixed?
	voxel    v3i.Vec // Intersecting voxel: to be computed by logic.
}

// A voxel is detected that intersects with the point load location.
// A vertex inside that voxel is selected, i.e. the vertex closest to load location.
// The vertex is assigned the load.
// This way, point load is applied to a single node/vertex of the mesh.
type Load struct {
	Location  v3.Vec  // Exact coordinates inside restraint.
	Magnitude v3.Vec  // X, Y, Z magnitude.
	voxel     v3i.Vec // Intersecting voxel: to be computed by logic.
	nodeREF   v3.Vec  // Eventual vertex/node to which the load is applied. To be computed.
}

func NewRestraint(location v3.Vec, isFixedX, isFixedY, isFixedZ bool) *Restraint {
	return &Restraint{
		Location: location,
		IsFixedX: isFixedX,
		IsFixedY: isFixedY,
		IsFixedZ: isFixedZ,
		voxel:    v3i.Vec{X: -1, Y: -1, Z: -1}, // We depend on -1 value to see if voxel is valid.
	}
}

// All the nodes inside the input voxel will be considered as restraint.
func NewRestraintByVoxel(voxel v3i.Vec, isFixedX, isFixedY, isFixedZ bool) *Restraint {
	return &Restraint{
		IsFixedX: isFixedX,
		IsFixedY: isFixedY,
		IsFixedZ: isFixedZ,
		voxel:    voxel,
	}
}

func NewLoad(location v3.Vec, magnitude v3.Vec) *Load {
	return &Load{
		Location:  location,
		Magnitude: magnitude,
		voxel:     v3i.Vec{X: -1, Y: -1, Z: -1}, // We depend on -1 value to see if voxel is valid.
	}
}
