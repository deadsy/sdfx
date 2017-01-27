//-----------------------------------------------------------------------------
/*

QuadTrees

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

const (
	TL = iota // top left
	TR        // top right
	BL        // bottom left
	BR        // bottom right
)

type QTMeta struct {
	posn V2
}

type QTNode struct {
	Parent   *QTNode    // pointer to the node parent
	Children [4]*QTNode // pointers to the node children
	Meta     [4]*QTMeta // meta data associated with each node vertex/corner
	Posn     int        // position wrt parent TL, TR, BL, BR

	//topleft  V2  // the top left position for this node
	//size     V2  // size of this node
	//level    int // level in quadtree, root node == 0
}

type QTree struct {
	Root *QTNode
}

//-----------------------------------------------------------------------------

// return the meta data for this node vertex
func (n *QTNode) GetMeta(i int) *QTMeta {

	if n.Parent == nil {
		// We are at the top of the tree, return what we have
		return n.Meta[i]
	}

	// Note: Implicit in this evaluation is the idea that meta data will
	// be created for the nodes in a specific order:
	// parent, top left, top right, bottom left, bottom right

	switch n.Posn {
	case TL:
		switch i {
		case TL:
			return n.Parent.GetMeta(TL)
		}
	case TR:
		switch i {
		case TL:
			return n.Parent.Children[TL].GetMeta(TR)
		case TR:
			return n.Parent.GetMeta(TR)
		case BL:
			return n.Parent.Children[TL].GetMeta(BR)
		}
	case BL:
		switch i {
		case TL:
			return n.Parent.Children[TL].GetMeta(BL)
		case TR:
			return n.Parent.Children[TL].GetMeta(BR)
		case BL:
			return n.Parent.GetMeta(BL)
		}
	case BR:
		switch i {
		case TL:
			return n.Parent.Children[TL].GetMeta(BR)
		case TR:
			return n.Parent.Children[TR].GetMeta(BR)
		case BL:
			return n.Parent.Children[BL].GetMeta(BR)
		case BR:
			return n.Parent.GetMeta(BR)
		}
	}

	return n.Meta[i]
}

//-----------------------------------------------------------------------------
