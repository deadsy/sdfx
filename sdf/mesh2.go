//-----------------------------------------------------------------------------
/*

2D Mesh, 2d line segments connected to create closed polygons.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// lineInfo stores pre-calculated line information.
type lineInfo struct {
	line       *Line2  // line segment
	unitVector v2.Vec  // unit vector for line segment
	length     float64 // length of line segment
}

// newLineInfo pre-calculates the line segment information.
func newLineInfo(l *Line2) *lineInfo {
	v := l[1].Sub(l[0])
	return &lineInfo{
		line:       l,
		unitVector: v.Normalize(),
		length:     v.Length(),
	}
}

func convertLines(lSet []*Line2) []*lineInfo {
	li := make([]*lineInfo, len(lSet))
	for i := range lSet {
		li[i] = newLineInfo(lSet[i])
	}
	return li
}

// minDistance2 returns the minium distance squared between a point and the line.
func (a *lineInfo) minDistance2(p v2.Vec) float64 {
	var d2 float64
	pa := p.Sub(a.line[0])
	// t-parameter of projection onto line
	t := pa.Dot(a.unitVector)
	if t < 0 {
		// distance to vertex 0 of line
		d2 = a.line[0].Sub(p).Length2()
	} else if t > a.length {
		// distance to vertex 1 of line
		d2 = a.line[1].Sub(p).Length2()
	} else {
		// normal distance from p to line
		dn := pa.Dot(v2.Vec{a.unitVector.Y, -a.unitVector.X})
		d2 = dn * dn
	}
	return d2
}

//-----------------------------------------------------------------------------

const qtMaxLevel = 10

type qtNode struct {
	level    int         // quadtree level
	box      Box2        // bounding box for the node
	center   v2.Vec      // pre-calculated from box
	halfSide float64     // pre-calculated from box
	child    [4]*qtNode  // child nodes (sw, se, nw, ne)
	leaf     []*lineInfo // leaf information (non-nil for a leaf node)
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
			leaf:     convertLines(lSet),
		}
	}

	// non-leaf node
	box0 := box.Quad0()
	box1 := box.Quad1()
	box2 := box.Quad2()
	box3 := box.Quad3()
	return &qtNode{
		level:    level,
		box:      box,
		halfSide: halfSide,
		center:   center,
		child: [4]*qtNode{
			qtBuild(level+1, box0, box0.lineFilter(lSet)),
			qtBuild(level+1, box1, box1.lineFilter(lSet)),
			qtBuild(level+1, box2, box2.lineFilter(lSet)),
			qtBuild(level+1, box3, box3.lineFilter(lSet)),
		},
	}
}

// boxes returns the set of boxes used by this node.
func (node *qtNode) boxes() []*Box2 {
	if node == nil {
		return nil
	}
	if node.leaf != nil {
		return []*Box2{&node.box}
	}
	boxes := []*Box2{&node.box}
	boxes = append(boxes, node.child[0].boxes()...)
	boxes = append(boxes, node.child[1].boxes()...)
	boxes = append(boxes, node.child[2].boxes()...)
	boxes = append(boxes, node.child[3].boxes()...)
	return boxes
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

// minFeatureDist2 returns the minimum distance squared from a point to the leaf feature.
func (node *qtNode) minLeafDist2(p v2.Vec) float64 {
	dd := math.MaxFloat64
	for _, li := range node.leaf {
		dd = math.Min(dd, li.minDistance2(p))
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
// Mesh2D. 2D mesh evaluation with quadtree speedup.

// MeshSDF2 is SDF2 made from a set of line segments.
type MeshSDF2 struct {
	qt *qtNode // quadtree root
	bb Box2    // bounding box
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

	// The quadtree box is derived from the bounding box.
	// Square it up for simpler math.
	// Scale it slightly to contain line segments on the top/right edges.
	qtBox := bb.Square().ScaleAboutCenter(1.01)

	// build the quadtree
	qt := qtBuild(0, qtBox, mesh)

	return &MeshSDF2{
		qt: qt,
		bb: bb,
	}, nil
}

// Evaluate returns the minimum distance for a 2d mesh.
func (s *MeshSDF2) Evaluate(p v2.Vec) float64 {
	d2 := s.qt.minDist2(p, math.MaxFloat64)
	return math.Sqrt(d2)
}

// Boxes returns the full set of quadtree boxes.
func (s *MeshSDF2) Boxes() []*Box2 {
	return s.qt.boxes()
}

// BoundingBox returns the bounding box of a 2d mesh.
func (s *MeshSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Mesh2D Slow. Provided for testing and benchmarking purposes.

// MeshSDF2Slow is SDF2 made from a set of line segments.
type MeshSDF2Slow struct {
	mesh []*lineInfo
	bb   Box2 // bounding box
}

// Mesh2DSlow returns an SDF2 made from a set of line segments.
func Mesh2DSlow(mesh []*Line2) (SDF2, error) {
	n := len(mesh)
	if n == 0 {
		return nil, ErrMsg("no 2d line segments")
	}

	// work out the bounding box
	bb := mesh[0].BoundingBox()
	for _, edge := range mesh {
		bb = bb.Include(edge[0]).Include(edge[1])
	}

	return &MeshSDF2Slow{
		mesh: convertLines(mesh),
		bb:   bb,
	}, nil
}

// Evaluate returns the minimum distance for a 2d mesh.
func (s *MeshSDF2Slow) Evaluate(p v2.Vec) float64 {
	d2 := math.MaxFloat64
	for _, li := range s.mesh {
		d2 = math.Min(d2, li.minDistance2(p))
	}
	return math.Sqrt(d2)
}

// BoundingBox returns the bounding box of a 2d mesh.
func (s *MeshSDF2Slow) BoundingBox() Box2 {
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
