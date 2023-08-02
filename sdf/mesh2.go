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

type vertexInfo struct {
	vertex v2.Vec   // coordinates of this vertex
	edge   []*Line2 // edges for this vertex
}

type qtNode struct {
	level int         // quadtree level
	box   Box2        // bounding box for the node
	child [4]*qtNode  // child nodes (sw, se, nw, ne)
	vInfo *vertexInfo // vertex information (non-nil for a leaf node)
}

// vertexFilter returns the set of vertices contained within the box.
func vertexFilter(vSet []int, box Box2, vInfo []vertexInfo) []int {
	var result []int
	for _, i := range vSet {
		if box.Contains(vInfo[i].vertex) {
			result = append(result, i)
		}
	}
	return result
}

func qtBuild(level int, box Box2, vInfo []vertexInfo, vSet []int) *qtNode {

	if len(vSet) == 0 {
		// empty node
		return nil
	}

	if len(vSet) == 1 {
		// leaf node
		return &qtNode{
			level: level,
			box:   box,
			vInfo: &vInfo[vSet[0]],
		}
	}

	// non-leaf node
	node := &qtNode{
		level: level,
		box:   box,
	}
	box0 := box.Quad0()
	box1 := box.Quad1()
	box2 := box.Quad2()
	box3 := box.Quad3()
	node.child[0] = qtBuild(level+1, box0, vInfo, vertexFilter(vSet, box0, vInfo))
	node.child[1] = qtBuild(level+1, box1, vInfo, vertexFilter(vSet, box1, vInfo))
	node.child[2] = qtBuild(level+1, box2, vInfo, vertexFilter(vSet, box2, vInfo))
	node.child[3] = qtBuild(level+1, box3, vInfo, vertexFilter(vSet, box3, vInfo))
	return node
}

// searchOrder returns the child search order for this node.
// Order by minimum distance to child box.
func (node *qtNode) searchOrder(p v2.Vec) [4]int {
	// translate the point so the node box center is at the origin
	p = p.Sub(node.box.Center())
	// gradient of box diagonal
	size := node.box.Size()
	m := size.Y / size.X
	if p.X >= 0 {
		if p.Y >= 0 {
			// quad3
			if p.Y >= m*p.X {
				return [4]int{3, 2, 1, 0}
			} else {
				return [4]int{3, 1, 2, 0}
			}
		} else {
			// quad1
			if p.Y <= -m*p.X {
				return [4]int{1, 0, 3, 2}
			} else {
				return [4]int{1, 3, 0, 2}
			}
		}
	} else {
		if p.Y >= 0 {
			// quad2
			if p.Y >= -m*p.X {
				return [4]int{2, 3, 0, 1}
			} else {
				return [4]int{2, 0, 3, 1}
			}
		} else {
			// quad0
			if p.Y <= m*p.X {
				return [4]int{0, 1, 2, 3}
			} else {
				return [4]int{0, 2, 1, 3}
			}
		}
	}
}

func (node *qtNode) minBoxDist2(p v2.Vec) float64 {

	return 0
}

func (node *qtNode) minDist2(p v2.Vec, dist2 float64) float64 {
	if node == nil || node.minBoxDist2(p) >= dist2 {
		return dist2
	}
	if node.vInfo != nil {
		// TODO actual calculation
		x := 100.0
		return math.Min(dist2, x)
	}
	// search the child nodes
	order := node.searchOrder(p)
	for _, i := range order {
		dist2 = node.child[i].minDist2(p, dist2)
	}
	return dist2
}

func (node *qtNode) countLeaves() int {
	if node == nil {
		return 0
	}
	if node.vInfo != nil {
		return 1
	}
	n := node.child[0].countLeaves()
	n += node.child[1].countLeaves()
	n += node.child[2].countLeaves()
	n += node.child[3].countLeaves()
	return n
}

func (node *qtNode) maxLevel() int {
	if node == nil {
		return 0
	}
	if node.vInfo != nil {
		return node.level
	}
	l0 := node.child[0].maxLevel()
	l1 := node.child[1].maxLevel()
	l2 := node.child[2].maxLevel()
	l3 := node.child[3].maxLevel()
	return maxInt(maxInt(l0, l1), maxInt(l2, l3))
}

//-----------------------------------------------------------------------------

// MeshSDF2 is SDF2 made from a set of line segments.
type MeshSDF2 struct {
	mesh  []*Line2     // polygon edges
	vInfo []vertexInfo // vertex information
	qt    *qtNode      // quadtree root
	bb    Box2         // bounding box
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
	// scale the bounding box slightly to contain vertices on the max edge
	bb = bb.ScaleAboutCenter(1.01)

	// create the vertex information
	vIndex := make(map[v2.Vec]int)
	var vInfo []vertexInfo
	for _, edge := range mesh {
		for _, vertex := range edge {
			if i, ok := vIndex[vertex]; ok {
				// existing vertex - add the edge
				vInfo[i].edge = append(vInfo[i].edge, edge)
			} else {
				// new vertex
				vInfo = append(vInfo, vertexInfo{vertex: vertex, edge: []*Line2{edge}})
				vIndex[vertex] = len(vInfo) - 1
			}
		}
	}

	// build the quadtree
	vSet := make([]int, len(vInfo))
	for i := range vSet {
		vSet[i] = i
	}
	qt := qtBuild(0, bb, vInfo, vSet)

	return &MeshSDF2{
		mesh:  mesh,
		vInfo: vInfo,
		qt:    qt,
		bb:    bb,
	}, nil

}

// Evaluate returns the minimum distance for a 2d mesh.
func (s *MeshSDF2) Evaluate(p v2.Vec) float64 {

	fmt.Printf("bb: %v\n", s.bb)
	fmt.Printf("edges: %d\n", len(s.mesh))
	fmt.Printf("vertices: %d\n", len(s.vInfo))
	fmt.Printf("leaves: %d\n", s.qt.countLeaves())
	fmt.Printf("max level: %d\n", s.qt.maxLevel())

	return 0
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
		n += 1
	}
	// create the mesh line segments
	mesh := make([]*Line2, n-1)
	for i := range mesh {
		mesh[i] = &Line2{vertex[i], vertex[i+1]}
	}
	return mesh, nil
}

//-----------------------------------------------------------------------------
