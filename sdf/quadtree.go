//-----------------------------------------------------------------------------
/*

QuadTrees

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

const (
	ROOT = iota - 1 // root node
	TL              // top left
	TR              // top right
	BL              // bottom left
	BR              // bottom right
)

const BASE_INC = (1 << 63)

//-----------------------------------------------------------------------------

const (
	VERTEX0 = iota // distance to vertex 0 on line
	VERTEX1        // distance to vertex 1 on line
	LINE           // distance to line body
)

type QTInfo struct {
	dtype  int  // distance type
	index  int  // line index
	inside bool // true if the point is inside the polygon
}

func (a QTInfo) Equals(b QTInfo) bool {
	return (a.dtype == b.dtype) && (a.index == b.index) && (a.inside == b.inside)
}

func (t *QTree) GetInfo(n *QTNode, posn int) QTInfo {
	return QTInfo{}
}

//-----------------------------------------------------------------------------

type QTNode struct {
	child  [4]*QTNode // pointers to the node children
	leaf   bool       // true if this node is a leaf (no children)
	corner V2         // top left corner of node box
	size   V2         // size of node box in x,y directions
	xn, yn uint64     // top left corner x,y integer names
	inc    uint64     // integer increment
}

type QTree struct {
	root   *QTNode // root node
	corner V2      // top left corner of bounding box
	size   V2      // size of bounding box in x,y directions
}

//-----------------------------------------------------------------------------

func (t *QTree) NewQTNode(parent *QTNode, posn int) *QTNode {
	n := QTNode{}

	// work out the node position values
	if posn == ROOT {
		// we are the root node and have no parent
		// get the parameters from the tree
		n.corner = t.corner
		n.size = t.size
		n.xn = 0
		n.yn = 0
		n.inc = BASE_INC
	} else {
		// the node size and increment is half the parent size
		n.size = parent.size.MulScalar(0.5)
		n.inc = parent.inc / 2
		switch posn {
		case TL:
			n.corner = parent.corner
			n.xn = parent.xn
			n.yn = parent.yn
		case TR:
			n.corner = parent.corner.Add(V2{n.size.X, 0})
			n.xn = parent.xn + n.inc
			n.yn = parent.yn
		case BL:
			n.corner = parent.corner.Add(V2{0, n.size.Y})
			n.xn = parent.xn
			n.yn = parent.yn + n.inc
		case BR:
			n.corner = parent.corner.Add(n.size)
			n.xn = parent.xn + n.inc
			n.yn = parent.yn + n.inc
		}
	}

	// evaluate the corner positions
	i0 := t.GetInfo(&n, TL)
	i1 := t.GetInfo(&n, TR)
	i2 := t.GetInfo(&n, BL)
	i3 := t.GetInfo(&n, BR)

	// create children if we have to
	if i0.Equals(i1) && i2.Equals(i3) && i0.Equals(i2) {
		// they are all the same, this is a leaf node
		n.leaf = true
	} else {
		// make the children
		n.child[TL] = t.NewQTNode(&n, TL)
		n.child[TR] = t.NewQTNode(&n, TR)
		n.child[BL] = t.NewQTNode(&n, BL)
		n.child[BR] = t.NewQTNode(&n, BR)
	}

	return &n
}

//-----------------------------------------------------------------------------

func NewQTree(lines []Line2, bb Box2) *QTree {
	t := QTree{}
	t.root = t.NewQTNode(nil, ROOT)
	return &t
}

//-----------------------------------------------------------------------------
