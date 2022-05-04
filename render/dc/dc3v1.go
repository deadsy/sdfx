//-----------------------------------------------------------------------------
/*

Dual Contouring

Convert an SDF3 to a triangle mesh.
Uses octree space subdivision.
Supports sharp edges and octree-based mesh simplification.
Based on: https://github.com/nickgildea/DualContouringSample

*/
//-----------------------------------------------------------------------------

package dc

import (
	"fmt"
	"math"
	"sort"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	"gonum.org/v1/gonum/mat"
)

//-----------------------------------------------------------------------------

// DualContouringV1 renders using dual contouring (octree sampling, sharp edges!, automatic simplification)
type DualContouringV1 struct {
	// Simplify: how much to simplify (if >=0).
	// NOTE: Meshing might fail with simplification enabled (FIXME),
	// but the mesh might can still simplified later using external tools (the main benefit of dual contouring is sharp edges).
	Simplify float64
	// RCond [0, 1) is the parameter that controls the accuracy of sharp edges, with lower being more accurate
	// but it can cause instability leading to large wrong triangles. Leave the default if unsure.
	RCond float64
	// LockVertices makes sure each vertex stays in its voxel, avoiding small or bad triangles that may be generated
	// otherwise, but it also may remove some sharp edges.
	LockVertices bool
}

// NewDualContouringV1 see DualContouringV1
func NewDualContouringV1(simplify float64, RCond float64, lockVertices bool) *DualContouringV1 {
	return &DualContouringV1{Simplify: simplify, RCond: RCond, LockVertices: lockVertices}
}

// Info returns a string describing the rendered volume.
func (m *DualContouringV1) Info(s sdf.SDF3, meshCells int) string {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := conv.V3ToV3i(bbSize.DivScalar(resolution))
	return fmt.Sprintf("%dx%dx%d, resolution %.2f", cells.X, cells.Y, cells.Z, resolution)
}

// Render produces a 3d triangle mesh over the bounding volume of an sdf3.
func (m *DualContouringV1) Render(s sdf.SDF3, meshCells int, output chan<- *render.Triangle3) {
	if m.RCond == 0 {
		m.RCond = 1e-3
	}
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := conv.V3ToV3i(bbSize.DivScalar(resolution))
	// Build the octree
	dcOctreeRootNode := dcNewOctree(cells, m.RCond, m.LockVertices)
	dcOctreeRootNode.Populate(s)
	// Simplify it
	if m.Simplify >= 0 {
		dcOctreeRootNode.Simplify(s, m.Simplify)
	}
	// Generate the final mesh
	dcOctreeRootNode.GenerateMesh(output)
}

//-----------------------------------------------------------------------------

var dcChildMinOffsets = [8]sdf.V3i{
	{0, 0, 0},
	{0, 0, 1},
	{0, 1, 0},
	{0, 1, 1},
	{1, 0, 0},
	{1, 0, 1},
	{1, 1, 0},
	{1, 1, 1},
}

var dcEdgevmap = [12][2]int{
	{0, 4}, {1, 5}, {2, 6}, {3, 7}, // x-axis
	{0, 2}, {1, 3}, {4, 6}, {5, 7}, // y-axis
	{0, 1}, {2, 3}, {4, 5}, {6, 7}, // z-axis
}
var dcEdgemask = [3]int{5, 3, 6}

var dcVertMap = [8][3]int{
	{0, 0, 0},
	{0, 0, 1},
	{0, 1, 0},
	{0, 1, 1},
	{1, 0, 0},
	{1, 0, 1},
	{1, 1, 0},
	{1, 1, 1},
}

var dcFaceMap = [6][4]int{{4, 8, 5, 9}, {6, 10, 7, 11}, {0, 8, 1, 10}, {2, 9, 3, 11}, {0, 4, 2, 6}, {1, 5, 3, 7}}

var dcCellProcFaceMask = [12][3]int{{0, 4, 0}, {1, 5, 0}, {2, 6, 0}, {3, 7, 0}, {0, 2, 1}, {4, 6, 1}, {1, 3, 1}, {5, 7, 1}, {0, 1, 2}, {2, 3, 2}, {4, 5, 2}, {6, 7, 2}}

var dcCellProcEdgeMask = [6][5]int{{0, 1, 2, 3, 0}, {4, 5, 6, 7, 0}, {0, 4, 1, 5, 1}, {2, 6, 3, 7, 1}, {0, 2, 4, 6, 2}, {1, 3, 5, 7, 2}}

var dcFaceProcFaceMask = [3][4][3]int{
	{{4, 0, 0}, {5, 1, 0}, {6, 2, 0}, {7, 3, 0}},
	{{2, 0, 1}, {6, 4, 1}, {3, 1, 1}, {7, 5, 1}},
	{{1, 0, 2}, {3, 2, 2}, {5, 4, 2}, {7, 6, 2}},
}

var dcFaceProcEdgeMask = [3][4][6]int{
	{{1, 4, 0, 5, 1, 1}, {1, 6, 2, 7, 3, 1}, {0, 4, 6, 0, 2, 2}, {0, 5, 7, 1, 3, 2}},
	{{0, 2, 3, 0, 1, 0}, {0, 6, 7, 4, 5, 0}, {1, 2, 0, 6, 4, 2}, {1, 3, 1, 7, 5, 2}},
	{{1, 1, 0, 3, 2, 0}, {1, 5, 4, 7, 6, 0}, {0, 1, 5, 0, 4, 1}, {0, 3, 7, 2, 6, 1}},
}

var dcEdgeProcEdgeMask = [3][2][5]int{
	{{3, 2, 1, 0, 0}, {7, 6, 5, 4, 0}},
	{{5, 1, 4, 0, 1}, {7, 3, 6, 2, 1}},
	{{6, 4, 2, 0, 2}, {7, 5, 3, 1, 2}},
}

var dcProcessEdgeMask = [3][4]int{{3, 2, 1, 0}, {7, 5, 6, 4}, {11, 10, 9, 8}}

type dcOctreeNodeType uint

const (
	dcOctreeNodeTypeInternal   dcOctreeNodeType = iota
	dcOctreeNodeTypePseudoLeaf                  // A simplified leaf node
	dcOctreeNodeTypeLeaf
)

type dcOctree struct {
	kind           dcOctreeNodeType
	minOffset      sdf.V3i
	size, meshSize int
	cellCounts     sdf.V3i
	children       [8]*dcOctree
	drawInfo       *dcOctreeDrawInfo
	// Extra parameters
	rCond        float64
	lockVertices bool
}

type dcOctreeDrawInfo struct {
	index, corners          int
	position, averageNormal sdf.V3
	qef                     *dcQefSolver
}

// nextPowerOfTwo is https://stackoverflow.com/questions/466204/rounding-up-to-next-power-of-2
func nextPowerOfTwo(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

// dcNewOctree builds the whole octree structure (without simplification) for the given size.
func dcNewOctree(cellCounts sdf.V3i, rCond float64, lockVertices bool) *dcOctree {
	cellCounts = sdf.V3i{ // Need powers of 2 for this algorithm (round-up for more precision)
		nextPowerOfTwo(cellCounts.X),
		nextPowerOfTwo(cellCounts.Y),
		nextPowerOfTwo(cellCounts.Z),
	}
	// Compute the complete octree with the largest component as the size and then ignoring cells outside of bounds
	cubicSize := int(conv.V3iToV3(cellCounts).MaxComponent())
	rootNode := &dcOctree{
		kind:         dcOctreeNodeTypeInternal,
		minOffset:    sdf.V3i{0, 0, 0},
		size:         cubicSize,
		meshSize:     cubicSize,
		cellCounts:   cellCounts,
		children:     [8]*dcOctree{},
		drawInfo:     nil,
		rCond:        rCond,
		lockVertices: lockVertices,
	}
	return rootNode
}

func (node *dcOctree) Populate(d sdf.SDF3) {
	minOffset := node.minOffset
	meshSize := node.meshSize
	cellCounts := node.cellCounts
	maxOffset := minOffset.AddScalar(meshSize)
	// Avoid generating any octree node outside the bounding volume (may filter before reaching leaves)
	if minOffset.X > (meshSize+cellCounts.X)/2 || maxOffset.X < (meshSize-cellCounts.X)/2 ||
		minOffset.Y > (meshSize+cellCounts.Y)/2 || maxOffset.Y < (meshSize-cellCounts.Y)/2 ||
		minOffset.Z > (meshSize+cellCounts.Z)/2 || maxOffset.Z < (meshSize-cellCounts.Z)/2 {
		return
	}
	childSize := node.size / 2
	for i := 0; i < 8; i++ {
		childMinOffset := minOffset.Add(conv.V3ToV3i(conv.V3iToV3(dcChildMinOffsets[i]).MulScalar(float64(childSize))))
		node.children[i] = &dcOctree{
			kind:         dcOctreeNodeTypeInternal,
			minOffset:    childMinOffset,
			size:         childSize,
			meshSize:     meshSize,
			cellCounts:   cellCounts,
			children:     [8]*dcOctree{},
			drawInfo:     nil,
			rCond:        node.rCond,
			lockVertices: node.lockVertices,
		}
		// Recursive children or a leaf node
		if childSize > 1 {
			node.children[i].Populate(d)
		} else {
			node.children[i].computeOctreeLeaf(d)
		}
	}
}

func (node *dcOctree) relToSDF(d sdf.SDF3, i sdf.V3i) sdf.V3 {
	bb := d.BoundingBox()
	return bb.Min.Add(bb.Size().Mul(conv.V3iToV3(i).DivScalar(float64(node.meshSize)).
		Div(conv.V3iToV3(node.cellCounts).DivScalar(float64(node.meshSize)))))
}

// computeOctreeLeaf computes the required leaf information that later will be used for meshing
func (node *dcOctree) computeOctreeLeaf(d sdf.SDF3) {
	corners := 0
	for i := 0; i < 8; i++ {
		cornerPos := node.relToSDF(d, node.minOffset.Add(dcChildMinOffsets[i]))
		isSolid := d.Evaluate(cornerPos) < 0
		if isSolid {
			corners = corners | (1 << i)
		}
	}
	if corners == 0 || corners == 255 {
		// voxel is fully inside or outside the volume: store nil for this child
		return
	}
	// otherwise, the voxel contains the surface, so find the edge intersections
	const maxCrossings = 6
	edgeCount := 0
	normalSum := sdf.V3{X: 0, Y: 0, Z: 0}
	qefSolver := new(dcQefSolver)
	for i := 0; i < 12 && edgeCount < maxCrossings; i++ {
		c1 := dcEdgevmap[i][0]
		c2 := dcEdgevmap[i][1]
		m1 := (corners >> c1) & 1
		m2 := (corners >> c2) & 1
		if (m1 == 1 && m2 == 1) || (m1 == 0 && m2 == 0) {
			// no zero crossing on this edge
			continue
		}
		p1 := node.relToSDF(d, node.minOffset.Add(dcChildMinOffsets[c1]))
		p2 := node.relToSDF(d, node.minOffset.Add(dcChildMinOffsets[c2]))
		p := dcApproximateZeroCrossingPosition(d, p1, p2)
		n := dcCalculateSurfaceNormal(d, p)
		qefSolver.Add(p, n)
		normalSum = normalSum.Add(n)
		edgeCount++
	}
	qefPosition := qefSolver.Solve(node.rCond)
	// See documentation of the next function
	if node.lockVertices {
		qefPosition = dcBoundVertexPosition(d, node, qefPosition, qefSolver)
	}
	node.drawInfo = &dcOctreeDrawInfo{
		index:         -1,
		corners:       corners,
		position:      qefPosition,
		averageNormal: normalSum.DivScalar(float64(edgeCount)).Normalize(),
		qef:           qefSolver,
	}
	node.kind = dcOctreeNodeTypeLeaf
}

// dcBoundVertexPosition binds the given vertex to their right voxel by using the mass point if out of bounds.
// NOTE: The next code avoids small triangles and even bad meshes (on noisy fields?), but reduces sharp edge accuracy
func dcBoundVertexPosition(d sdf.SDF3, leaf *dcOctree, qefPosition sdf.V3, qefSolver *dcQefSolver) sdf.V3 {
	// Avoid placing vertex outside node bounds
	min := leaf.relToSDF(d, leaf.minOffset)
	max := leaf.relToSDF(d, leaf.minOffset.Add(sdf.V3i{leaf.size, leaf.size, leaf.size}))
	if qefPosition.X < min.X || qefPosition.Y < min.Y || qefPosition.Z < min.Z ||
		qefPosition.X > max.X || qefPosition.Y > max.Y || qefPosition.Z > max.Z {
		//log.Println("Fixing vertex position", qefPosition, "-->", qefSolver.massPointSum)
		qefPosition = qefSolver.MassPoint()
	} else {
		//log.Println("NOT fixing vertex position", qefPosition)
	}
	return qefPosition
}

// Simplify optionally simplifies the octree structure merging planar faces before meshing
// should be called several times to support multi-level simplification (until false is returned)
func (node *dcOctree) Simplify(d sdf.SDF3, threshold float64) {
	if node == nil {
		return
	}
	if node.kind != dcOctreeNodeTypeInternal {
		return // can't Simplify!
	}
	isCollapsible := true
	qefSolver := new(dcQefSolver)
	signs := [8]int{-1, -1, -1, -1, -1, -1, -1, -1}
	midSign := -1
	edgeCount := 0
	for i := 0; i < 8; i++ {
		node.children[i].Simplify(d, threshold)
		if node.children[i] != nil {
			if node.children[i].kind == dcOctreeNodeTypeInternal {
				isCollapsible = false
			} else {
				qefSolver.AddSolver(node.children[i].drawInfo.qef)
				midSign = (node.children[i].drawInfo.corners >> (7 - i)) & 1
				signs[i] = (node.children[i].drawInfo.corners >> i) & 1
				edgeCount++
			}
		}
	}
	if !isCollapsible { // at least one child is an internal node, can't collapse
		return
	}
	// If no children have surface, force simplifying this node (position shouldn't matter, and qef will be left empty for no influence on parents)
	qefPosition := node.relToSDF(d, node.minOffset)
	if qefSolver.numPoints > 0 {
		// Otherwise, solve the qef of our children
		qefPosition = qefSolver.Solve(node.rCond)
		qefError := qefSolver.GetError() // Errors caused by forced simplification
		// See documentation of the next function
		if node.lockVertices {
			qefPosition = dcBoundVertexPosition(d, node, qefPosition, qefSolver)
		}
		if qefError > threshold {
			return
		}
	}
	// Build the pseudo leaf node as all checks passed
	node.kind = dcOctreeNodeTypePseudoLeaf
	node.drawInfo = new(dcOctreeDrawInfo)
	node.drawInfo.position = qefPosition
	node.drawInfo.qef = qefSolver
	for i := 0; i < 8; i++ {
		if signs[i] == -1 { // Undetermined, use centre sign instead
			node.drawInfo.corners |= midSign << i
		} else {
			node.drawInfo.corners |= signs[i] << i
		}
	}
	for i := 0; i < 8; i++ {
		child := node.children[i]
		if child != nil && (child.kind == dcOctreeNodeTypePseudoLeaf || child.kind == dcOctreeNodeTypeLeaf) {
			node.drawInfo.averageNormal = node.drawInfo.averageNormal.Add(child.drawInfo.averageNormal)
		}
	}
	node.drawInfo.averageNormal = node.drawInfo.averageNormal.Normalize()
	// Remove simplified children
	for i := 0; i < 8; i++ {
		node.children[i] = nil
	}
	return
}

func (node *dcOctree) generateVertexIndices(vertexBuffer *[]sdf.V3) {
	if node == nil { // Does not contain the surface
		return
	}
	if node.kind == dcOctreeNodeTypeInternal { // Add vertices to children
		for i := 0; i < 8; i++ {
			node.children[i].generateVertexIndices(vertexBuffer)
		}
	} else { // Leaf or pseudo-leaf node: add one vertex
		node.drawInfo.index = len(*vertexBuffer)
		*vertexBuffer = append(*vertexBuffer, node.drawInfo.position)
	}
}

func (node *dcOctree) contourCellProc(indexBuffer *[]int) {
	if node == nil { // Does not contain the surface
		return
	}
	if node.kind == dcOctreeNodeTypeInternal {
		for i := 0; i < 8; i++ {
			node.children[i].contourCellProc(indexBuffer)
		}
		for i := 0; i < 12; i++ {
			c := dcCellProcFaceMask[i][0:2]
			faceNodes0 := node.children[c[0]]
			faceNodes1 := node.children[c[1]]
			dcContourFaceProc([2]*dcOctree{faceNodes0, faceNodes1}, dcCellProcFaceMask[i][2], indexBuffer)
		}
		for i := 0; i < 6; i++ {
			edgeNodes := [4]*dcOctree{}
			c := dcCellProcEdgeMask[i][0:4]
			for j := 0; j < 4; j++ {
				edgeNodes[j] = node.children[c[j]]
			}
			dcContourEdgeProc(edgeNodes, dcCellProcEdgeMask[i][4], indexBuffer)
		}
	}
}

func dcContourFaceProc(node [2]*dcOctree, dir int, indexBuffer *[]int) {
	if node[0] == nil || node[1] == nil { // Does not contain the surface
		return
	}
	if node[0].kind == dcOctreeNodeTypeInternal || node[1].kind == dcOctreeNodeTypeInternal {
		for i := 0; i < 4; i++ {
			faceNodes := [2]*dcOctree{}
			c := dcFaceProcFaceMask[dir][i][0:2]
			for j := 0; j < 2; j++ {
				if node[j].kind != dcOctreeNodeTypeInternal {
					faceNodes[j] = node[j]
				} else {
					faceNodes[j] = node[j].children[c[j]]
				}
			}
			dcContourFaceProc(faceNodes, dcFaceProcFaceMask[dir][i][2], indexBuffer)
		}
		orders := [2][4]int{
			{0, 0, 1, 1},
			{0, 1, 0, 1},
		}
		for i := 0; i < 4; i++ {
			edgeNodes := [4]*dcOctree{}
			c := dcFaceProcEdgeMask[dir][i][1:5]
			order := orders[dcFaceProcEdgeMask[dir][i][0]]
			for j := 0; j < 4; j++ {
				if node[order[j]].kind == dcOctreeNodeTypeLeaf || node[order[j]].kind == dcOctreeNodeTypePseudoLeaf {
					edgeNodes[j] = node[order[j]]
				} else {
					edgeNodes[j] = node[order[j]].children[c[j]]
				}
			}
			dcContourEdgeProc(edgeNodes, dcFaceProcEdgeMask[dir][i][5], indexBuffer)
		}
	}
}

func dcContourEdgeProc(node [4]*dcOctree, dir int, indexBuffer *[]int) {
	if node[0] == nil || node[1] == nil || node[2] == nil || node[3] == nil { // Does not contain the surface
		return
	}
	if node[0].kind != dcOctreeNodeTypeInternal && node[1].kind != dcOctreeNodeTypeInternal &&
		node[2].kind != dcOctreeNodeTypeInternal && node[3].kind != dcOctreeNodeTypeInternal {
		dcContourProcessEdge(node, dir, indexBuffer)
	} else {
		for i := 0; i < 2; i++ {
			edgeNodes := [4]*dcOctree{}
			c := dcEdgeProcEdgeMask[dir][i][0:4]
			for j := 0; j < 4; j++ {
				if node[j].kind != dcOctreeNodeTypeInternal {
					edgeNodes[j] = node[j]
				} else {
					edgeNodes[j] = node[j].children[c[j]]
				}
			}
			dcContourEdgeProc(edgeNodes, dcEdgeProcEdgeMask[dir][i][4], indexBuffer)
		}
	}
}

func dcContourProcessEdge(node [4]*dcOctree, dir int, indexBuffer *[]int) {
	minSize := math.MaxInt
	minIndex := 0
	indices := [4]int{-1, -1, -1, -1}
	flip := false
	signChange := [4]bool{false, false, false, false}
	for i := 0; i < 4; i++ {
		edge := dcProcessEdgeMask[dir][i]
		c1 := dcEdgevmap[edge][0]
		c2 := dcEdgevmap[edge][1]
		m1 := (node[i].drawInfo.corners >> c1) & 1
		m2 := (node[i].drawInfo.corners >> c2) & 1
		if node[i].size < minSize {
			minSize = node[i].size
			minIndex = i
			flip = m1 != 0 // Make the triangles face the right way
		}
		indices[i] = node[i].drawInfo.index
		signChange[i] = m1 != m2
	}
	if signChange[minIndex] {
		if !flip {
			*indexBuffer = append(*indexBuffer, indices[0])
			*indexBuffer = append(*indexBuffer, indices[1])
			*indexBuffer = append(*indexBuffer, indices[3])

			*indexBuffer = append(*indexBuffer, indices[0])
			*indexBuffer = append(*indexBuffer, indices[3])
			*indexBuffer = append(*indexBuffer, indices[2])
		} else {
			*indexBuffer = append(*indexBuffer, indices[0])
			*indexBuffer = append(*indexBuffer, indices[3])
			*indexBuffer = append(*indexBuffer, indices[1])

			*indexBuffer = append(*indexBuffer, indices[0])
			*indexBuffer = append(*indexBuffer, indices[2])
			*indexBuffer = append(*indexBuffer, indices[3])
		}
	}
}

func (node *dcOctree) GenerateMesh(output chan<- *render.Triangle3) {
	vertexBuffer := new([]sdf.V3)
	indexBuffer := new([]int)
	// Populate buffers
	node.generateVertexIndices(vertexBuffer)
	node.contourCellProc(indexBuffer)
	// Return triangles
	for tri := 0; tri < len(*indexBuffer)/3; tri++ {
		triangle := &render.Triangle3{
			V: [3]sdf.V3{
				(*vertexBuffer)[(*indexBuffer)[tri*3]],
				(*vertexBuffer)[(*indexBuffer)[tri*3+1]],
				(*vertexBuffer)[(*indexBuffer)[tri*3+2]],
			},
		}
		//log.Println("Outputting triangle:", triangle)
		output <- triangle
	}
}

// dcQefSolver is used for vertex position estimation (sharp edges!)
type dcQefSolver struct {
	ata                  *mat.SymDense
	atb, massPointSum, x sdf.V3
	btb                  float64
	numPoints            int
	hasSolution          bool
}

func (q *dcQefSolver) Add(p, n sdf.V3) {
	n = n.Normalize()
	if q.ata == nil {
		q.ata = mat.NewSymDense(3, nil)
	}
	q.ata.SetSym(0, 0, q.ata.At(0, 0)+n.X*n.X)
	q.ata.SetSym(0, 1, q.ata.At(0, 1)+n.X*n.Y)
	q.ata.SetSym(0, 2, q.ata.At(0, 2)+n.X*n.Z)
	q.ata.SetSym(1, 1, q.ata.At(1, 1)+n.Y*n.Y)
	q.ata.SetSym(1, 2, q.ata.At(1, 2)+n.Y*n.Z)
	q.ata.SetSym(2, 2, q.ata.At(2, 2)+n.Z*n.Z)
	dot := p.Dot(n)
	q.atb = q.atb.Add(n.MulScalar(dot))
	q.btb += dot * dot
	q.massPointSum = q.massPointSum.Add(p)
	q.numPoints++
	q.hasSolution = false
}

func (q *dcQefSolver) AddSolver(q2 *dcQefSolver) {
	if q.ata == nil {
		q.ata = mat.NewSymDense(3, nil)
	}
	if q2.ata != nil {
		q.ata.AddSym(q.ata, q2.ata)
	}
	q.atb = q.atb.Add(q2.atb)
	q.btb += q2.btb
	q.massPointSum = q.massPointSum.Add(q2.massPointSum)
	q.numPoints += q2.numPoints
	q.hasSolution = false
}

func (q *dcQefSolver) MassPoint() sdf.V3 {
	return q.massPointSum.DivScalar(float64(q.numPoints))
}

func (q *dcQefSolver) Solve(rCond float64) sdf.V3 {
	// assert q.ata != nil (some points inserted)
	// VecUtils::scale(this->massPointSum, 1.0f / this->data.numPoints);
	massPointClone := q.MassPoint()
	// MatUtils::vmul_symmetric(tmpv, this->ata, this->massPoint);
	var tmpV mat.Dense
	tmpV.Mul(q.ata, toVec(massPointClone))
	// VecUtils::sub(this->atb, this->atb, tmpv); (with later reset in same function: declare new variable)
	atb := q.atb.Sub(toV3(&tmpV))

	// const float result = Svd::solveSymmetric(this->ata, this->atb, this->x, svd_tol, svd_sweeps, pinv_tol);
	var x mat.VecDense
	// SVD
	svd := new(mat.SVD)
	if !svd.Factorize(q.ata, mat.SVDThin) {
		return massPointClone // If factorization fails (for example for Box), return the mass point
	}
	_ = svd.SolveVecTo(&x, toVec(atb), svd.Rank(rCond))
	// QR (needs stabilization)
	//qr := new(mat.QR)
	//qr.Factorize(q.ata)
	//_ = qr.SolveVecTo(&x, true, toVec(atb))

	// VecUtils::addScaled(this->x, 1, this->massPoint); (previous clear in this function makes this ok)
	q.x = massPointClone.Add(toV3(&x))
	// VecUtils::addScaled(this->x, 1, this->massPointSum);
	q.hasSolution = true
	return q.x
}

func (q *dcQefSolver) GetError() float64 {
	return q.getErrorPos(&q.x)
}

func (q *dcQefSolver) getErrorPos(pos *sdf.V3) float64 {
	//MatUtils::vmul_symmetric(atax, this->ata, pos);
	var atax mat.Dense
	atax.Mul(q.ata, toVec(*pos))
	// return VecUtils::dot(pos, atax) - 2 * VecUtils::dot(pos, this->atb) + this->data.btb;
	return pos.Dot(toV3(&atax)) - 2*pos.Dot(q.atb) + q.btb
}

func toVec(massPointClone sdf.V3) *mat.VecDense {
	return mat.NewVecDense(3, []float64{massPointClone.X, massPointClone.Y, massPointClone.Z})
}

func toV3(x mat.Matrix) sdf.V3 {
	return sdf.V3{X: x.At(0, 0), Y: x.At(1, 0), Z: x.At(2, 0)}
}

//-----------------------------------------------------------------------------

func dcApproximateZeroCrossingPosition(d sdf.SDF3, p0, p1 sdf.V3) sdf.V3 {
	const steps = 8. // good enough precision? Note that errors easily explode or cause noise/instability on future operations
	// Original implementation:
	//minValue := math.MaxFloat64
	//t := 0.
	//const increment = 1. / steps // Relative to p0 <--> p1 edge length
	//for currentT := 0.; currentT <= 1.; currentT += increment {
	//	p := p0.Add(p1.Sub(p0).MulScalar(currentT))
	//	d := math.Abs(d.Evaluate(p))
	//	if d < minValue {
	//		minValue = d
	//		t = currentT
	//	}
	//}
	//return p0.Add(p1.Sub(p0).MulScalar(t))
	// Alternative: binary search. IMPORTANT: leads to better simplification!
	fakeElems := math.Pow(2, steps)
	searchSolid := !(d.Evaluate(p0) < 0)
	foundIndex := sort.Search(int(fakeElems), func(fakeElem int) bool {
		currentT := float64(fakeElem) / fakeElems
		p := p0.Add(p1.Sub(p0).MulScalar(currentT))
		foundSolid := d.Evaluate(p) < 0
		return searchSolid && foundSolid || !searchSolid && !foundSolid
	})
	t := float64(foundIndex) / fakeElems
	//log.Println("t:", t, "val:", d.Evaluate(p0.Add(p1.Sub(p0).MulScalar(t))))
	return p0.Add(p1.Sub(p0).MulScalar(t))
}

func dcCalculateSurfaceNormal(d sdf.SDF3, p sdf.V3) sdf.V3 {
	const eps = 0.001
	return sdf.V3{
		X: d.Evaluate(p.Add(sdf.V3{X: eps})) - d.Evaluate(p.Add(sdf.V3{X: -eps})),
		Y: d.Evaluate(p.Add(sdf.V3{Y: eps})) - d.Evaluate(p.Add(sdf.V3{Y: -eps})),
		Z: d.Evaluate(p.Add(sdf.V3{Z: eps})) - d.Evaluate(p.Add(sdf.V3{Z: -eps})),
	}.Normalize()
}
