package mesh

import (
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//  1. Define each restraint by two things:
//     1.1. A collection of points, *not* by a single point.
//     1.2. Degrees of freedom that are fixed/free for all the points of the collection.
//  2. Assume a bounding box around all those points.
//  3. Any voxel that intersects with the bounding box is considered inside restraint.
//  4. Create `*BOUNDARY` for all the nodes inside those voxels.
//  5. Degree of freedom would be fixed/free for all those nodes.
//
// The objective: the stress concentration at the restraint may be alleviated by distributing it.
type Restraint struct {
	Location []v3.Vec  // Exact coordinates inside rigid body.
	IsFixedX bool      // Is X degree of freedom fixed?
	IsFixedY bool      // Is Y degree of freedom fixed?
	IsFixedZ bool      // Is Z degree of freedom fixed?
	voxels   []v3i.Vec // Intersecting voxels: to be computed by logic.
}

// Point loads are applied to the nodes of the mesh.
type Load struct {
	Location  []v3.Vec  // Exact coordinates inside rigid body.
	Magnitude v3.Vec    // X, Y, Z magnitude.
	voxels    []v3i.Vec // Intersecting voxels: to be computed by logic.
	nodeREF   uint32    // Eventual node to which the load is applied. To be computed.
}

func NewRestraint(location []v3.Vec, isFixedX, isFixedY, isFixedZ bool) *Restraint {
	return &Restraint{
		Location: location,
		IsFixedX: isFixedX,
		IsFixedY: isFixedY,
		IsFixedZ: isFixedZ,
		voxels:   make([]v3i.Vec, 0),
	}
}

func NewLoad(location []v3.Vec, magnitude v3.Vec) *Load {
	return &Load{
		Location:  location,
		Magnitude: magnitude,
		voxels:    make([]v3i.Vec, 0),
		nodeREF:   0,
	}
}
