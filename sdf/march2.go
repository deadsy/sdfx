//-----------------------------------------------------------------------------
/*

Marching Squares

Convert an SDF2 boundary to a set of line segments.

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

type LineCache struct {
	base  V2        // base coordinate of line
	inc   V2        // dx, dy for each step
	steps V2i       // number of x,y steps
	val0  []float64 // SDF values for x line
	val1  []float64 // SDF values for x + dx line
}

func NewLineCache(base, inc V2, steps V2i) *LineCache {
	return &LineCache{base, inc, steps, nil, nil}
}

// Evaluate the SDF for a given Y line
func (l *LineCache) Evaluate(sdf SDF2, x int) {

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
	var p V2
	p.X = l.base.X + float64(x)*dx

	// evaluate the line
	p.Y = l.base.Y
	for y := 0; y < ny+1; y++ {
		l.val1[idx] = sdf.Evaluate(p)
		idx += 1
		p.Y += dy
	}
}

func (l *LineCache) Get(x, y int) float64 {
	if x == 0 {
		return l.val0[y]
	}
	return l.val1[y]
}

//-----------------------------------------------------------------------------

func MarchingSquares(sdf SDF2, box Box2, step float64) []*Line2_PP {

	var lines []*Line2_PP
	size := box.Size()
	base := box.Min
	steps := size.DivScalar(step).Ceil().ToV2i()
	inc := size.Div(steps.ToV2())

	// create the line cache
	l := NewLineCache(base, inc, steps)
	// evaluate the SDF for x = 0
	l.Evaluate(sdf, 0)

	nx, ny := steps[0], steps[1]
	dx, dy := inc.X, inc.Y

	var p V2
	p.X = base.X
	for x := 0; x < nx; x++ {
		// read the x + 1 layer
		l.Evaluate(sdf, x+1)
		// process all squares in the x and x + 1 layers
		p.Y = base.Y
		for y := 0; y < ny; y++ {
			x0, y0 := p.X, p.Y
			x1, y1 := x0+dx, y0+dy
			corners := [4]V2{
				{x0, y0},
				{x1, y0},
				{x1, y1},
				{x0, y1},
			}
			values := [4]float64{
				l.Get(0, y),
				l.Get(1, y),
				l.Get(1, y+1),
				l.Get(0, y+1),
			}
			lines = append(lines, ms_ToLines(corners, values, 0)...)
			p.Y += dy
		}
		p.X += dx
	}

	return lines
}

//-----------------------------------------------------------------------------

// generate the line segments for a square
func ms_ToLines(p [4]V2, v [4]float64, x float64) []*Line2_PP {
	// which of the 0..15 patterns do we have?
	index := 0
	for i := 0; i < 4; i++ {
		if v[i] < x {
			index |= 1 << uint(i)
		}
	}
	// do we have any lines to create?
	if ms_edge_table[index] == 0 {
		return nil
	}
	// work out the interpolated points on the edges
	var points [4]V2
	for i := 0; i < 4; i++ {
		bit := 1 << uint(i)
		if ms_edge_table[index]&bit != 0 {
			a := ms_pair_table[i][0]
			b := ms_pair_table[i][1]
			points[i] = ms_Interpolate(p[a], p[b], v[a], v[b], x)
		}
	}
	// create the line segments
	table := ms_line_table[index]
	count := len(table) / 2
	result := make([]*Line2_PP, count)
	for i := 0; i < count; i++ {
		line := Line2_PP{}
		line[1] = points[table[i*2+0]]
		line[0] = points[table[i*2+1]]
		result[i] = &line
	}
	return result
}

//-----------------------------------------------------------------------------

func ms_Interpolate(p1, p2 V2, v1, v2, x float64) V2 {
	if Abs(x-v1) < EPS {
		return p1
	}
	if Abs(x-v2) < EPS {
		return p2
	}
	if Abs(v1-v2) < EPS {
		return p1
	}
	t := (x - v1) / (v2 - v1)
	return V2{
		p1.X + t*(p2.X-p1.X),
		p1.Y + t*(p2.Y-p1.Y),
	}
}

//-----------------------------------------------------------------------------

// These are the vertex pairs for the edges
var ms_pair_table = [4][2]int{
	{0, 1}, // 0
	{1, 2}, // 1
	{2, 3}, // 2
	{3, 0}, // 3
}

// 4 vertices -> 16 possible inside/outside combinations
// A 1 bit in the value indicates an edge with a line end point.
// 4 edges -> 4 bit values, note the fwd/rev symmetry
var ms_edge_table = [16]int{
	0x0, 0x9, 0x3, 0xa,
	0x6, 0xf, 0x5, 0xc,
	0xc, 0x5, 0xf, 0x6,
	0xa, 0x3, 0x9, 0x0,
}

// specify the edges used to create the line(s)
var ms_line_table = [16][]int{
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
