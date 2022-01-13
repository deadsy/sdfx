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
)

//-----------------------------------------------------------------------------

// lineCache is a cache of SDF2 evaluations samples over a 2d line.
type lineCache struct {
	base  sdf.V2    // base coordinate of line
	inc   sdf.V2    // dx, dy for each step
	steps sdf.V2i   // number of x,y steps
	val0  []float64 // SDF values for x line
	val1  []float64 // SDF values for x + dx line
}

// newLineCache returns a line cache.
func newLineCache(base, inc sdf.V2, steps sdf.V2i) *lineCache {
	return &lineCache{base, inc, steps, nil, nil}
}

// evaluate the SDF2 over a given x line.
func (l *lineCache) evaluate(s sdf.SDF2, x int) {

	// Swap the layers
	l.val0, l.val1 = l.val1, l.val0

	ny := l.steps[1]
	dx, dy := l.inc.X, l.inc.Y

	// allocate storage
	if l.val1 == nil {
		l.val1 = make([]float64, ny+1)
	}

	// setup the loop variables
	idx := 0
	var p sdf.V2
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

// MarchingSquaresUniform renders using marching squares with uniform space sampling.
type MarchingSquaresUniform struct {
}

// Info returns a string describing the rendered shape.
func (m *MarchingSquaresUniform) Info(s sdf.SDF2, meshCells int) string {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := bb1Size.ToV2i()
	bb1Size = bb1Size.MulScalar(meshInc)
	return fmt.Sprintf("%dx%d", cells[0], cells[1])
}

// Render produces a 2D line mesh over the bounding volume of an sdf2.
func (m *MarchingSquaresUniform) Render(s sdf.SDF2, meshCells int, output chan<- *Line) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox2(bb0.Center(), bb1Size)
	for _, tri := range marchingSquares(s, bb, meshInc) {
		output <- tri
	}
}

//-----------------------------------------------------------------------------

func marchingSquares(s sdf.SDF2, box sdf.Box2, step float64) []*Line {

	var lines []*Line
	size := box.Size()
	base := box.Min
	steps := size.DivScalar(step).Ceil().ToV2i()
	inc := size.Div(steps.ToV2())

	// create the line cache
	l := newLineCache(base, inc, steps)
	// evaluate the SDF for x = 0
	l.evaluate(s, 0)

	nx, ny := steps[0], steps[1]
	dx, dy := inc.X, inc.Y

	var p sdf.V2
	p.X = base.X
	for x := 0; x < nx; x++ {
		// read the x + 1 layer
		l.evaluate(s, x+1)
		// process all squares in the x and x + 1 layers
		p.Y = base.Y
		for y := 0; y < ny; y++ {
			x0, y0 := p.X, p.Y
			x1, y1 := x0+dx, y0+dy
			corners := [4]sdf.V2{
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

// generate the line segments for a square
func msToLines(p [4]sdf.V2, v [4]float64, x float64) []*Line {
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
	var points [4]sdf.V2
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
	result := make([]*Line, 0, count)
	for i := 0; i < count; i++ {
		l := Line{}
		l[1] = points[table[i*2+0]]
		l[0] = points[table[i*2+1]]
		if !l.Degenerate(0) {
			result = append(result, &l)
		}
	}
	return result
}

//-----------------------------------------------------------------------------

func msInterpolate(p1, p2 sdf.V2, v1, v2, x float64) sdf.V2 {

	closeToV1 := math.Abs(x-v1) < epsilon
	closeToV2 := math.Abs(x-v2) < epsilon

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
		t = (x - v1) / (v2 - v1)
	}

	return sdf.V2{
		p1.X + t*(p2.X-p1.X),
		p1.Y + t*(p2.Y-p1.Y),
	}
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
