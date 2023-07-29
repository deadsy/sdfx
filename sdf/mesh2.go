//-----------------------------------------------------------------------------
/*

2D Mesh, 2d line segments connected to create closed polygons.

Uses R-trees to efficiently search for nearest features in the line set.
https://en.wikipedia.org/wiki/R-tree

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"

	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/dhconnelly/rtreego"
)

//-----------------------------------------------------------------------------

// rtree child parameters (tunables)
const minChildren = 3
const maxChildren = 5

// MeshSDF2 is SDF2 made from a set of line segments.
type MeshSDF2 struct {
	rtree *rtreego.Rtree // r-tree root
	bb    Box2           // bounding box
}

// Mesh2D returns an SDF2 made from a set of line segments.
func Mesh2D(mesh []*Line2) (SDF2, error) {
	n := len(mesh)
	if n == 0 {
		return nil, errors.New("no 2d line segments")
	}

	// r-tree bulk loading
	load := make([]rtreego.Spatial, n)

	// pre-process the line set
	bb := mesh[0].BoundingBox()
	for i, l := range mesh {
		load[i] = l
		bb = bb.Extend(l.BoundingBox())
	}

	return &MeshSDF2{
		rtree: rtreego.NewTree(2, minChildren, maxChildren, load...),
		bb:    bb,
	}, nil

}

// Evaluate returns the minimum distance for a 2d mesh.
func (s *MeshSDF2) Evaluate(p v2.Vec) float64 {
	return 0
}

// BoundingBox returns the bounding box of a 2d mesh.
func (s *MeshSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
