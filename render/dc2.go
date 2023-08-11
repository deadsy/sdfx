//-----------------------------------------------------------------------------
/*

2d Dual Contouring Renderer

Resources:

https://www.mattkeeter.com/projects/contours/
https://www.graphics.rwth-aachen.de/publication/131/feature1.pdf
https://www.boristhebrave.com/2018/04/15/dual-contouring-tutorial/


1) create the dual graph
2) position the vertices within the non empty leaf nodes
3) merge leaf cells for mesh simplification

*/
//-----------------------------------------------------------------------------

package render

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v2i"
)

//-----------------------------------------------------------------------------

// norm2 returns the normal to the SDF2 at a point.
func norm2(s sdf.SDF2, p v2.Vec, epsilon float64) v2.Vec {
	return v2.Vec{
		s.Evaluate(v2.Vec{p.X + epsilon, p.Y}) - s.Evaluate(v2.Vec{p.X - epsilon, p.Y}),
		s.Evaluate(v2.Vec{p.X, p.Y + epsilon}) - s.Evaluate(v2.Vec{p.X, p.Y - epsilon}),
	}.Normalize()
}

//-----------------------------------------------------------------------------

type node2 struct {
	v     v2i.Vec // origin of square as integers
	n     uint    // level of square, size = 1 << n
	child []node2 // child nodes
}

type dc2 struct {
	origin     v2.Vec              // origin of the overall bounding square
	resolution float64             // size of smallest quadtree square
	hdiag      []float64           // lookup table of square half diagonals
	s          sdf.SDF2            // the SDF2 to be rendered
	cache      map[v2i.Vec]float64 // cache of distances
}

func newDualContouring2(s sdf.SDF2, origin v2.Vec, resolution float64, n uint) *dc2 {
	dc := dc2{
		origin:     origin,
		resolution: resolution,
		hdiag:      make([]float64, n),
		s:          s,
		cache:      make(map[v2i.Vec]float64),
	}
	// build a lut for cube half diagonal lengths
	for i := range dc.hdiag {
		si := 1 << uint(i)
		s := float64(si) * dc.resolution
		dc.hdiag[i] = 0.5 * math.Sqrt(2.0*s*s)
	}
	return &dc
}

// read from the cache
func (dc *dc2) read(vi v2i.Vec) (float64, bool) {
	dist, found := dc.cache[vi]
	return dist, found
}

// write to the cache
func (dc *dc2) write(vi v2i.Vec, dist float64) {
	dc.cache[vi] = dist
}

func (dc *dc2) evaluate(vi v2i.Vec) (v2.Vec, float64) {
	v := dc.origin.Add(conv.V2iToV2(vi).MulScalar(dc.resolution))
	// do we have it in the cache?
	dist, found := dc.read(vi)
	if found {
		return v, dist
	}
	// evaluate the SDF2
	dist = dc.s.Evaluate(v)
	// write it to the cache
	dc.write(vi, dist)
	return v, dist
}

// isEmpty returns true if the node contains no SDF surface
func (dc *dc2) isEmpty(c *node2) bool {
	// evaluate the SDF2 at the center of the square
	s := 1 << (c.n - 1) // half side
	_, d := dc.evaluate(c.v.AddScalar(s))
	// compare to the center/corner distance
	return math.Abs(d) >= dc.hdiag[c.n]
}

func (dc *dc2) processNode(node *node2) {
	if !dc.isEmpty(node) {
		if node.n == 1 {

		} else {
			// create the sub-nodes
			n := node.n - 1
			s := 1 << n
			node.child = make([]node2, 4)
			node.child[0].v = node.v.Add(v2i.Vec{0, 0})
			node.child[0].n = n
			node.child[1].v = node.v.Add(v2i.Vec{s, 0})
			node.child[1].n = n
			node.child[2].v = node.v.Add(v2i.Vec{s, s})
			node.child[2].n = n
			node.child[3].v = node.v.Add(v2i.Vec{0, s})
			node.child[3].n = n
			// process the sub-nodes
			dc.processNode(&node.child[0])
			dc.processNode(&node.child[1])
			dc.processNode(&node.child[2])
			dc.processNode(&node.child[3])
		}
	}
}

//-----------------------------------------------------------------------------

func (dc *dc2) corner(vi v2i.Vec) v2.Vec {
	return dc.origin.Add(conv.V2iToV2(vi).MulScalar(dc.resolution))
}

func (dc *dc2) drawNode(node *node2, output sdf.Line2Writer) {

	k := int(node.n) * 2

	c0 := dc.corner(node.v.Add(v2i.Vec{0, 0}))
	c1 := dc.corner(node.v.Add(v2i.Vec{k, 0}))
	c2 := dc.corner(node.v.Add(v2i.Vec{k, k}))
	c3 := dc.corner(node.v.Add(v2i.Vec{0, k}))

	l0 := sdf.Line2{c0, c1}
	l1 := sdf.Line2{c1, c2}
	l2 := sdf.Line2{c2, c3}
	l3 := sdf.Line2{c3, c0}

	output.Write([]*sdf.Line2{&l0, &l1, &l2, &l3})
}

func (dc *dc2) qtOutput(node *node2, output sdf.Line2Writer) {
	if node.child != nil {
		dc.qtOutput(&node.child[0], output)
		dc.qtOutput(&node.child[1], output)
		dc.qtOutput(&node.child[2], output)
		dc.qtOutput(&node.child[3], output)
	} else {

		if node.n == 1 {
			dc.drawNode(node, output)
		}
	}
}

// dualContouring2D generates line segments for an SDF2 using dual contouring.
func dualContouring2D(s sdf.SDF2, resolution float64, output sdf.Line2Writer) {
	// Scale the bounding box about the center to make sure the boundaries
	// aren't on the object surface.
	bb := s.BoundingBox()
	bb = bb.ScaleAboutCenter(1.01)
	longAxis := bb.Size().MaxComponent()
	// We want to test the smallest squares (side == resolution) for emptiness
	// so the level = 0 cube is at half resolution.
	resolution = 0.5 * resolution
	// how many cube levels for the quadtree?
	levels := uint(math.Ceil(math.Log2(longAxis/resolution))) + 1
	// create the dual contouring state
	dc := newDualContouring2(s, bb.Min, resolution, levels)
	// process the quadtree, start at the top level
	topNode := node2{v: v2i.Vec{0, 0}, n: levels - 1}
	dc.processNode(&topNode)
	dc.qtOutput(&topNode, output)
	output.Close()
}

//-----------------------------------------------------------------------------

// DualContouring2D renders is a 2D dual contouring renderer.
type DualContouring2D struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewDualContouring2D returns a Render2 object.
func NewDualContouring2D(meshCells int) *DualContouring2D {
	return &DualContouring2D{
		meshCells: meshCells,
	}
}

// Info returns a string describing the rendered area.
func (r *DualContouring2D) Info(s sdf.SDF2) string {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	cells := conv.V2ToV2i(bbSize.MulScalar(1 / resolution))
	return fmt.Sprintf("%dx%d, resolution %.2f", cells.X, cells.Y, resolution)
}

// Render produces a 2d line mesh over the bounding area of an sdf2.
func (r *DualContouring2D) Render(s sdf.SDF2, output sdf.Line2Writer) {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	dualContouring2D(s, resolution, output)
}

//-----------------------------------------------------------------------------
