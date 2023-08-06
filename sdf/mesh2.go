//-----------------------------------------------------------------------------
/*

2D Mesh, 2d line segments connected to create closed polygons.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

const qtMaxLevel = 15

type leafInfo struct {
	line []*Line2 // lines stored at this leaf
}

type qtNode struct {
	level    int        // quadtree level
	box      Box2       // bounding box for the node
	center   v2.Vec     // pre-calculated from box
	halfSide float64    // pre-calculated from box
	child    [4]*qtNode // child nodes (sw, se, nw, ne)
	leaf     *leafInfo  // leaf information (non-nil for a leaf node)
}

func qtBuild(level int, box Box2, lSet []*Line2) *qtNode {

	if len(lSet) == 0 {
		// empty node
		return nil
	}

	halfSide := 0.5 * (box.Max.X - box.Min.X)
	center := box.Center()

	if len(lSet) == 1 || level == qtMaxLevel {
		// leaf node
		return &qtNode{
			level:    level,
			box:      box,
			halfSide: halfSide,
			center:   center,
			leaf:     &leafInfo{line: lSet},
		}
	}

	// non-leaf node
	node := &qtNode{
		level:    level,
		box:      box,
		halfSide: halfSide,
		center:   center,
	}
	box0 := box.Quad0()
	box1 := box.Quad1()
	box2 := box.Quad2()
	box3 := box.Quad3()
	node.child[0] = qtBuild(level+1, box0, box0.lineFilter(lSet))
	node.child[1] = qtBuild(level+1, box1, box1.lineFilter(lSet))
	node.child[2] = qtBuild(level+1, box2, box2.lineFilter(lSet))
	node.child[3] = qtBuild(level+1, box3, box3.lineFilter(lSet))
	return node
}

// searchOrder returns the child search order for this node.
// Order by minimum distance to the child boxes.
func (node *qtNode) searchOrder(p v2.Vec) [4]int {
	// translate the point so the node box center is at the origin
	p = p.Sub(node.center)
	if p.X >= 0 {
		if p.Y >= 0 {
			// quad3
			if p.Y >= p.X {
				return [4]int{3, 2, 1, 0}
			}
			return [4]int{3, 1, 2, 0}
		}
		// quad1
		if p.Y <= -p.X {
			return [4]int{1, 0, 3, 2}
		}
		return [4]int{1, 3, 0, 2}
	}
	if p.Y >= 0 {
		// quad2
		if p.Y >= -p.X {
			return [4]int{2, 3, 0, 1}
		}
		return [4]int{2, 0, 3, 1}
	}
	// quad0
	if p.Y <= p.X {
		return [4]int{0, 1, 2, 3}
	}
	return [4]int{0, 2, 1, 3}
}

// minBoxDist2 returns the minimum distance squared from a point to the node box.
// Inside the box is a zero distance.
func (node *qtNode) minBoxDist2(p v2.Vec) float64 {
	// translate the point so the node box center is at the origin
	// work in a single quadrant
	p = p.Sub(node.center).Abs()
	dx := p.X - node.halfSide
	dy := p.Y - node.halfSide
	// inside the box
	if dx < 0 && dy < 0 {
		return 0
	}
	if dy < 0 {
		return dx * dx
	}
	if dx < 0 {
		return dy * dy
	}
	return (dx * dx) + (dy * dy)
}

var leafCount int

// minFeatureDist2 returns the minimum distance squared from a point to the leaf feature.
func (node *qtNode) minLeafDist2(p v2.Vec) float64 {
	leafCount++
	dd := math.MaxFloat64
	for _, l := range node.leaf.line {
		dd = math.Min(dd, l.minDistance2(p))
	}
	return dd
}

func (node *qtNode) minDist2(p v2.Vec, dd float64) float64 {
	if node == nil || node.minBoxDist2(p) >= dd {
		// no new minimums here
		return dd
	}
	if node.leaf != nil {
		// measure the leaf
		return math.Min(dd, node.minLeafDist2(p))
	}
	// search the child nodes
	for _, i := range node.searchOrder(p) {
		dd = node.child[i].minDist2(p, dd)
	}
	return dd
}

//-----------------------------------------------------------------------------

// MeshSDF2 is SDF2 made from a set of line segments.
type MeshSDF2 struct {
	mesh []*Line2
	qt   *qtNode // quadtree root
	bb   Box2    // bounding box
}

// Mesh2D returns an SDF2 made from a set of line segments.
func Mesh2D(mesh []*Line2) (SDF2, error) {
	n := len(mesh)
	if n == 0 {
		return nil, ErrMsg("no 2d line segments")
	}

	// work out the bounding box
	bb := mesh[0].BoundingBox()
	for _, edge := range mesh {
		bb = bb.Include(edge[0]).Include(edge[1])
	}
	// square up the bounding box
	// scale it slightly to contain vertices on the max edge
	bb = bb.Square().ScaleAboutCenter(1.01)

	// build the quadtree
	qt := qtBuild(0, bb, mesh)

	return &MeshSDF2{
		mesh: mesh,
		qt:   qt,
		bb:   bb,
	}, nil

}

// Evaluate returns the minimum distance for a 2d mesh.
func (s *MeshSDF2) Evaluate(p v2.Vec) float64 {
	leafCount = 0
	dd := s.qt.minDist2(p, math.MaxFloat64)
	fmt.Printf("fast evals %d\n", leafCount)
	return math.Sqrt(dd)
}

// EvaluateSlow returns the minimum distance for a 2d mesh (slowly).
func (s *MeshSDF2) EvaluateSlow(p v2.Vec) float64 {
	dd := math.MaxFloat64
	for _, l := range s.mesh {
		dd = math.Min(dd, l.minDistance2(p))
	}
	fmt.Printf("slow evals %d\n", len(s.mesh))
	return math.Sqrt(dd)
}

// BoundingBox returns the bounding box of a 2d mesh.
func (s *MeshSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// PolygonToMesh converts a polygon into a mesh (line segment) representation.
func PolygonToMesh(p *Polygon) ([]*Line2, error) {
	vertex := p.Vertices()
	n := len(vertex)
	if n < 3 {
		return nil, ErrMsg("number of vertices < 3")
	}
	// Close the loop (if necessary)
	if !vertex[0].Equals(vertex[n-1], tolerance) {
		vertex = append(vertex, vertex[0])
		n++
	}
	// create the mesh line segments
	mesh := make([]*Line2, n-1)
	for i := range mesh {
		mesh[i] = &Line2{vertex[i], vertex[i+1]}
	}
	return mesh, nil
}

//-----------------------------------------------------------------------------
