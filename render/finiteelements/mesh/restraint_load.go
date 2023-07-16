package mesh

import (
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

// Single point constraint: one or more degrees of freedom are fixed for a given node.
type Restraint struct {
	Location v3.Vec  // Exact coordinate.
	IsFixedX bool    // Is X degree of freedom fixed?
	IsFixedY bool    // Is Y degree of freedom fixed?
	IsFixedZ bool    // Is Z degree of freedom fixed?
	voxel    v3i.Vec // Containing voxel: to be computed by logic.
	nodeID   int     // Eventual node to which the restraint is applied. To be computed.
}

// Point loads are applied to the nodes of the mesh.
type Load struct {
	Location  v3.Vec  // Exact coordinate.
	Magnitude v3.Vec  // X, Y, Z magnitude.
	voxel     v3i.Vec // Containing voxel: to be computed by logic.
	nodeID    int     // Eventual node to which the load is applied. To be computed.
}

func NewRestraint(location v3.Vec, isFixedX, isFixedY, isFixedZ bool) *Restraint {
	return &Restraint{
		Location: location,
		IsFixedX: isFixedX,
		IsFixedY: isFixedY,
		IsFixedZ: isFixedZ,
		voxel:    v3i.Vec{},
		nodeID:   0,
	}
}

func NewLoad(location, magnitude v3.Vec) *Load {
	return &Load{
		Location:  location,
		Magnitude: magnitude,
		voxel:     v3i.Vec{},
		nodeID:    0,
	}
}

// Identify the node to which restraint is applied.
// Also, the containing voxel is figured out.
func (r *Restraint) FindNode() {
	// TODO.
}
