package mesh

import (
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//  1. Define each restraint by two things:
//     1.1. A collection of points, *not* by a single point.
//     1.2. Degrees of freedom that are fixed/free for all the points of the collection.
//  2. Assume those points are connected by straight lines.
//  3. Any voxel that intersects with those straight lines is considered rigid.
//  4. Create `*RIGID BODY` by all the elements or nodes inside those voxels.
//  5. Degree of freedom would be fixed/free for the `REF NODE` of the `*RIGID BODY`
//
// According to CCX manual, under the hood, a `*RIGID BODY` is actually a nonlinear multiple-point constraint (MPC).
//
// The objective: the stress concentration at the restraint may be alleviated by distributing it.
type Restraint struct {
	Location []v3.Vec // Exact coordinates inside rigid body.
	IsFixedX bool     // Is X degree of freedom fixed?
	IsFixedY bool     // Is Y degree of freedom fixed?
	IsFixedZ bool     // Is Z degree of freedom fixed?
	voxel    v3i.Vec  // Containing voxel: to be computed by logic.
	nodeID   uint32   // Eventual node to which the restraint is applied. To be computed.
}

// Point loads are applied to the nodes of the mesh.
type Load struct {
	Location  []v3.Vec // Exact coordinates inside rigid body.
	Magnitude v3.Vec   // X, Y, Z magnitude.
	voxel     v3i.Vec  // Containing voxel: to be computed by logic.
	nodeID    uint32   // Eventual node to which the load is applied. To be computed.
}

func NewRestraint(location []v3.Vec, isFixedX, isFixedY, isFixedZ bool) *Restraint {
	return &Restraint{
		Location: location,
		IsFixedX: isFixedX,
		IsFixedY: isFixedY,
		IsFixedZ: isFixedZ,
		voxel:    v3i.Vec{},
		nodeID:   0,
	}
}

func NewLoad(location []v3.Vec, magnitude v3.Vec) *Load {
	return &Load{
		Location:  location,
		Magnitude: magnitude,
		voxel:     v3i.Vec{},
		nodeID:    0,
	}
}
