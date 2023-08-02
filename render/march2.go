//-----------------------------------------------------------------------------
/*

Marching Squares

Convert an SDF2 boundary to a set of line segments.

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

// lineCache is a cache of SDF2 evaluations samples over a 2d line.
type lineCache struct {
	base  v2.Vec    // base coordinate of line
	inc   v2.Vec    // dx, dy for each step
	steps v2i.Vec   // number of x,y steps
	val0  []float64 // SDF values for x line
	val1  []float64 // SDF values for x + dx line
}

// newLineCache returns a line cache.
func newLineCache(base, inc v2.Vec, steps v2i.Vec) *lineCache {
	return &lineCache{base, inc, steps, nil, nil}
}

// evaluate the SDF2 over a given x line.
func (l *lineCache) evaluate(s sdf.SDF2, x int) {

	// Swap the layers
	l.val0, l.val1 = l.val1, l.val0

	ny := l.steps.Y
	dx, dy := l.inc.X, l.inc.Y

	// allocate storage
	if l.val1 == nil {
		l.val1 = make([]float64, ny+1)
	}

	// setup the loop variables
	idx := 0
	var p v2.Vec
	p.X = l.base.X + float64(x)*dx

	// evaluate the line
	p.Y = l.base.Y
	for y := 0; y < ny+1; y++ {
		l.val1[idx] = s.Evaluate(p)
		idx++
		p.Y += dy
	}
}

// get a value from a line cache.
func (l *lineCache) get(x, y int) float64 {
	if x == 0 {
		return l.val0[y]
	}
	return l.val1[y]
}

//-----------------------------------------------------------------------------

func marchingSquares(s sdf.SDF2, resolution float64) []*sdf.Line2 {
	// Scale the bounding box about the center to make sure the boundaries
	// aren't on the object surface.
	bb := s.BoundingBox()
	bb = bb.ScaleAboutCenter(1.01)

	size := bb.Size()
	base := bb.Min
	steps := conv.V2ToV2i(size.MulScalar(1 / resolution).Ceil())
	inc := size.Div(conv.V2iToV2(steps))

	// create the line cache
	l := newLineCache(base, inc, steps)
	// evaluate the SDF for x = 0
	l.evaluate(s, 0)

	nx, ny := steps.X, steps.Y
	dx, dy := inc.X, inc.Y

	var lines []*sdf.Line2
	var p v2.Vec
	p.X = base.X
	for x := 0; x < nx; x++ {
		// read the x + 1 layer
		l.evaluate(s, x+1)
		// process all squares in the x and x + 1 layers
		p.Y = base.Y
		for y := 0; y < ny; y++ {
			x0, y0 := p.X, p.Y
			x1, y1 := x0+dx, y0+dy
			corners := [4]v2.Vec{
				{x0, y0},
				{x1, y0},
				{x1, y1},
				{x0, y1},
			}
			values := [4]float64{
				l.get(0, y),
				l.get(1, y),
				l.get(1, y+1),
				l.get(0, y+1),
			}
			lines = append(lines, msToLines(corners, values, 0)...)
			p.Y += dy
		}
		p.X += dx
	}

	return lines
}

//-----------------------------------------------------------------------------

// MarchingSquaresUniform renders using marching squares with uniform area sampling.
type MarchingSquaresUniform struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewMarchingSquaresUniform returns a Render2 object.
func NewMarchingSquaresUniform(meshCells int) *MarchingSquaresUniform {
	return &MarchingSquaresUniform{
		meshCells: meshCells,
	}
}

// Info returns a string describing the rendered area.
func (r *MarchingSquaresUniform) Info(s sdf.SDF2) string {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	cells := conv.V2ToV2i(bbSize.MulScalar(1 / resolution))
	return fmt.Sprintf("%dx%d, resolution %.2f", cells.X, cells.Y, resolution)
}

// Render produces a 2d line mesh over the bounding area of an sdf2.
func (r *MarchingSquaresUniform) Render(s sdf.SDF2, output chan<- []*sdf.Line2) {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(r.meshCells)
	output <- marchingSquares(s, resolution)
}

//-----------------------------------------------------------------------------

// generate the line segments for a square
func msToLines(p [4]v2.Vec, v [4]float64, x float64) []*sdf.Line2 {
	// which of the 0..15 patterns do we have?
	index := 0
	for i := 0; i < 4; i++ {
		if v[i] < x {
			index |= 1 << uint(i)
		}
	}
	// do we have any lines to create?
	if msEdgeTable[index] == 0 {
		return nil
	}
	// work out the interpolated points on the edges
	var points [4]v2.Vec
	for i := 0; i < 4; i++ {
		bit := 1 << uint(i)
		if msEdgeTable[index]&bit != 0 {
			a := msPairTable[i][0]
			b := msPairTable[i][1]
			points[i] = msInterpolate(p[a], p[b], v[a], v[b], x)
		}
	}
	// create the line segments
	table := msLineTable[index]
	count := len(table) / 2
	result := make([]*sdf.Line2, 0, count)
	for i := 0; i < count; i++ {
		l := sdf.Line2{}
		l[1] = points[table[i*2+0]]
		l[0] = points[table[i*2+1]]
		if !l.Degenerate(0) {
			result = append(result, &l)
		}
	}
	return result
}

//-----------------------------------------------------------------------------

func msInterpolate(p1, p2 v2.Vec, k1, k2, x float64) v2.Vec {

	closeToV1 := math.Abs(x-k1) < epsilon
	closeToV2 := math.Abs(x-k2) < epsilon

	if closeToV1 && !closeToV2 {
		return p1
	}
	if closeToV2 && !closeToV1 {
		return p2
	}

	var t float64

	if closeToV1 && closeToV2 {
		// Pick the half way point
		t = 0.5
	} else {
		// linear interpolation
		t = (x - k1) / (k2 - k1)
	}
	return v2.Vec{p1.X + t*(p2.X-p1.X), p1.Y + t*(p2.Y-p1.Y)}
}

//-----------------------------------------------------------------------------

// These are the vertex pairs for the edges
var msPairTable = [4][2]int{
	{0, 1}, // 0
	{1, 2}, // 1
	{2, 3}, // 2
	{3, 0}, // 3
}

// 4 vertices -> 16 possible inside/outside combinations
// A 1 bit in the value indicates an edge with a line end point.
// 4 edges -> 4 bit values, note the fwd/rev symmetry
var msEdgeTable = [16]int{
	0x0, 0x9, 0x3, 0xa,
	0x6, 0xf, 0x5, 0xc,
	0xc, 0x5, 0xf, 0x6,
	0xa, 0x3, 0x9, 0x0,
}

// specify the edges used to create the line(s)
var msLineTable = [16][]int{
	{},           // 0
	{0, 3},       // 1
	{0, 1},       // 2
	{1, 3},       // 3
	{1, 2},       // 4
	{0, 1, 2, 3}, // 5
	{0, 2},       // 6
	{2, 3},       // 7
	{2, 3},       // 8
	{0, 2},       // 9
	{0, 3, 1, 2}, // 10
	{1, 2},       // 11
	{1, 3},       // 12
	{0, 1},       // 13
	{0, 3},       // 14
	{},           // 15
}

//-----------------------------------------------------------------------------
