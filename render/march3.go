//-----------------------------------------------------------------------------
/*

Marching Cubes

Convert an SDF3 to a triangle mesh.

*/
//-----------------------------------------------------------------------------

package render

import (
	"math"
	"runtime"
	"sync"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

type layerYZ struct {
	base          sdf.V3       // base coordinate of layer
	inc           sdf.V3       // dx, dy, dz for each step
	steps         sdf.V3i      // number of x,y,z steps
	evalProcessCh chan evalReq // the evaluation channel for parallelization
	val0          []float64    // SDF values for x layer
	val1          []float64    // SDF values for x + dx layer
}

func newLayerYZ(base, inc sdf.V3, steps sdf.V3i, evalProcessCh chan evalReq) *layerYZ {
	return &layerYZ{base, inc, steps, evalProcessCh, nil, nil}
}

// evalReq is used for processing evaluations in parallel.
//
// A slice of V3 is run through `fn`; the result of which
// is stored in the corresponding index of the `out` slice.
type evalReq struct {
	out []float64
	p   []sdf.V3
	fn  func(sdf.V3) float64
	wg  *sync.WaitGroup
}

// Evaluate the SDF for a given XY layer
func (l *layerYZ) Evaluate(s sdf.SDF3, x int) {

	// Swap the layers
	l.val0, l.val1 = l.val1, l.val0

	ny, nz := l.steps[1], l.steps[2]
	dx, dy, dz := l.inc.X, l.inc.Y, l.inc.Z

	// allocate storage
	if l.val1 == nil {
		l.val1 = make([]float64, (ny+1)*(nz+1))
	}

	// setup the loop variables
	idx := 0
	var p sdf.V3
	p.X = l.base.X + float64(x)*dx

	// define the base struct for requesting evaluation
	var eReq evalReq
	if l.evalProcessCh != nil {
		eReq = evalReq{
			wg:  new(sync.WaitGroup),
			fn:  s.Evaluate,
			out: l.val1,
		}
	}

	// evaluate the layer
	p.Y = l.base.Y

	// Performance doesn't seem to improve past 100.
	const batchSize = 100

	if l.evalProcessCh != nil {
		eReq.p = make([]sdf.V3, 0, batchSize)
	}
	for y := 0; y < ny+1; y++ {
		p.Z = l.base.Z
		for z := 0; z < nz+1; z++ {
			if l.evalProcessCh == nil { // Singlethread mode (just cache the evaluation)
				l.val1[idx] = s.Evaluate(p)
			} else { // Multithread mode: prepare and send slice for parallel processing using the evaluation goroutines
				eReq.p = append(eReq.p, p)
				if len(eReq.p) == batchSize {
					eReq.wg.Add(1)
					l.evalProcessCh <- eReq
					eReq.out = eReq.out[batchSize:]       // shift the output slice for processing
					eReq.p = make([]sdf.V3, 0, batchSize) // create a new slice for the next batch
				}
			}
			idx++
			p.Z += dz
		}
		p.Y += dy
	}

	if l.evalProcessCh != nil {
		// send any remaining points for processing
		if len(eReq.p) > 0 {
			eReq.wg.Add(1)
			l.evalProcessCh <- eReq
		}

		// Wait for all processing to complete before returning
		eReq.wg.Wait()
	}
}

func (l *layerYZ) Get(x, y, z int) float64 {
	idx := y*(l.steps[2]+1) + z
	if x == 0 {
		return l.val0[idx]
	}
	return l.val1[idx]
}

//-----------------------------------------------------------------------------

func marchingCubes(s sdf.SDF3, box sdf.Box3, step float64, out chan<- *Triangle3, goroutines int) {
	var evalProcessCh chan evalReq
	if goroutines > 1 {
		evalProcessCh = make(chan evalReq, 100)
		for i := 0; i < goroutines; i++ {
			go func() {
				var i int
				var p sdf.V3
				for r := range evalProcessCh {
					for i, p = range r.p {
						r.out[i] = r.fn(p)
					}
					r.wg.Done()
				}
			}()
		}
	}

	size := box.Size()
	base := box.Min
	steps := size.DivScalar(step).Ceil().ToV3i()
	inc := size.Div(steps.ToV3())

	// create the SDF layer cache
	l := newLayerYZ(base, inc, steps, evalProcessCh)
	// evaluate the SDF for x = 0
	l.Evaluate(s, 0)

	nx, ny, nz := steps[0], steps[1], steps[2]
	dx, dy, dz := inc.X, inc.Y, inc.Z

	var p sdf.V3
	p.X = base.X
	for x := 0; x < nx; x++ {
		// read the x + 1 layer
		l.Evaluate(s, x+1)
		// process all cubes in the x and x + 1 layers
		p.Y = base.Y
		for y := 0; y < ny; y++ {
			p.Z = base.Z
			for z := 0; z < nz; z++ {
				x0, y0, z0 := p.X, p.Y, p.Z
				x1, y1, z1 := x0+dx, y0+dy, z0+dz
				corners := [8]sdf.V3{
					{x0, y0, z0},
					{x1, y0, z0},
					{x1, y1, z0},
					{x0, y1, z0},
					{x0, y0, z1},
					{x1, y0, z1},
					{x1, y1, z1},
					{x0, y1, z1}}
				values := [8]float64{
					l.Get(0, y, z),
					l.Get(1, y, z),
					l.Get(1, y+1, z),
					l.Get(0, y+1, z),
					l.Get(0, y, z+1),
					l.Get(1, y, z+1),
					l.Get(1, y+1, z+1),
					l.Get(0, y+1, z+1)}
				//triangles = append(triangles, mcToTriangles(corners, values, 0)...)
				for _, tri := range mcToTriangles(corners, values, 0) {
					out <- tri
				}
				p.Z += dz
			}
			p.Y += dy
		}
		p.X += dx
	}
}

//-----------------------------------------------------------------------------

func mcToTriangles(p [8]sdf.V3, v [8]float64, x float64) []*Triangle3 {
	// which of the 0..255 patterns do we have?
	index := 0
	for i := 0; i < 8; i++ {
		if v[i] < x {
			index |= 1 << uint(i)
		}
	}
	// do we have any triangles to create?
	if mcEdgeTable[index] == 0 {
		return nil
	}
	// work out the interpolated points on the edges
	var points [12]sdf.V3
	for i := 0; i < 12; i++ {
		bit := 1 << uint(i)
		if mcEdgeTable[index]&bit != 0 {
			a := mcPairTable[i][0]
			b := mcPairTable[i][1]
			points[i] = mcInterpolate(p[a], p[b], v[a], v[b], x)
		}
	}
	// create the triangles
	table := mcTriangleTable[index]
	count := len(table) / 3
	result := make([]*Triangle3, 0, count)
	for i := 0; i < count; i++ {
		t := Triangle3{}
		t.V[2] = points[table[i*3+0]]
		t.V[1] = points[table[i*3+1]]
		t.V[0] = points[table[i*3+2]]
		if !t.Degenerate(0) {
			result = append(result, &t)
		}
	}
	return result
}

//-----------------------------------------------------------------------------

func mcInterpolate(p1, p2 sdf.V3, v1, v2, x float64) sdf.V3 {

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

	return sdf.V3{
		p1.X + t*(p2.X-p1.X),
		p1.Y + t*(p2.Y-p1.Y),
		p1.Z + t*(p2.Z-p1.Z),
	}
}

//-----------------------------------------------------------------------------

// MarchingCubesUniform renders using marching cubes with uniform space sampling.
type MarchingCubesUniform struct {
	// How many goroutines to spawn for parallel evaluation of the SDF3 (0 is runtime.NumCPU())
	// Set to 1 to avoid generating goroutines, useful for using the parallel renderer as a wrapper
	EvaluateGoroutines int
}

func (m *MarchingCubesUniform) Cells(s sdf.SDF3, meshCells int) (float64, sdf.V3i) {
	return DefaultRender3Cells(s, meshCells)
}

// Render produces a 3d triangle mesh over the bounding volume of an sdf3.
func (m *MarchingCubesUniform) Render(s sdf.SDF3, meshCells int, output chan<- *Triangle3) {
	if m.EvaluateGoroutines == 0 {
		m.EvaluateGoroutines = runtime.NumCPU() // Keep legacy behavior
	}
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size /*.Ceil().AddScalar(1) - Changed to work with multithread renderer: same behaviour as other renderers*/
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	marchingCubes(s, bb, meshInc, output, m.EvaluateGoroutines)
}

//-----------------------------------------------------------------------------

// These are the vertex pairs for the edges
var mcPairTable = [12][2]int{
	{0, 1}, // 0
	{1, 2}, // 1
	{2, 3}, // 2
	{3, 0}, // 3
	{4, 5}, // 4
	{5, 6}, // 5
	{6, 7}, // 6
	{7, 4}, // 7
	{0, 4}, // 8
	{1, 5}, // 9
	{2, 6}, // 10
	{3, 7}, // 11
}

// 8 vertices -> 256 possible inside/outside combinations
// A 1 bit in the value indicates an edge with a line end point.
// 12 edges -> 12 bit values, note the fwd/rev symmetry
var mcEdgeTable = [256]int{
	0x0000, 0x0109, 0x0203, 0x030a, 0x0406, 0x050f, 0x0605, 0x070c,
	0x080c, 0x0905, 0x0a0f, 0x0b06, 0x0c0a, 0x0d03, 0x0e09, 0x0f00,
	0x0190, 0x0099, 0x0393, 0x029a, 0x0596, 0x049f, 0x0795, 0x069c,
	0x099c, 0x0895, 0x0b9f, 0x0a96, 0x0d9a, 0x0c93, 0x0f99, 0x0e90,
	0x0230, 0x0339, 0x0033, 0x013a, 0x0636, 0x073f, 0x0435, 0x053c,
	0x0a3c, 0x0b35, 0x083f, 0x0936, 0x0e3a, 0x0f33, 0x0c39, 0x0d30,
	0x03a0, 0x02a9, 0x01a3, 0x00aa, 0x07a6, 0x06af, 0x05a5, 0x04ac,
	0x0bac, 0x0aa5, 0x09af, 0x08a6, 0x0faa, 0x0ea3, 0x0da9, 0x0ca0,
	0x0460, 0x0569, 0x0663, 0x076a, 0x0066, 0x016f, 0x0265, 0x036c,
	0x0c6c, 0x0d65, 0x0e6f, 0x0f66, 0x086a, 0x0963, 0x0a69, 0x0b60,
	0x05f0, 0x04f9, 0x07f3, 0x06fa, 0x01f6, 0x00ff, 0x03f5, 0x02fc,
	0x0dfc, 0x0cf5, 0x0fff, 0x0ef6, 0x09fa, 0x08f3, 0x0bf9, 0x0af0,
	0x0650, 0x0759, 0x0453, 0x055a, 0x0256, 0x035f, 0x0055, 0x015c,
	0x0e5c, 0x0f55, 0x0c5f, 0x0d56, 0x0a5a, 0x0b53, 0x0859, 0x0950,
	0x07c0, 0x06c9, 0x05c3, 0x04ca, 0x03c6, 0x02cf, 0x01c5, 0x00cc,
	0x0fcc, 0x0ec5, 0x0dcf, 0x0cc6, 0x0bca, 0x0ac3, 0x09c9, 0x08c0,
	0x08c0, 0x09c9, 0x0ac3, 0x0bca, 0x0cc6, 0x0dcf, 0x0ec5, 0x0fcc,
	0x00cc, 0x01c5, 0x02cf, 0x03c6, 0x04ca, 0x05c3, 0x06c9, 0x07c0,
	0x0950, 0x0859, 0x0b53, 0x0a5a, 0x0d56, 0x0c5f, 0x0f55, 0x0e5c,
	0x015c, 0x0055, 0x035f, 0x0256, 0x055a, 0x0453, 0x0759, 0x0650,
	0x0af0, 0x0bf9, 0x08f3, 0x09fa, 0x0ef6, 0x0fff, 0x0cf5, 0x0dfc,
	0x02fc, 0x03f5, 0x00ff, 0x01f6, 0x06fa, 0x07f3, 0x04f9, 0x05f0,
	0x0b60, 0x0a69, 0x0963, 0x086a, 0x0f66, 0x0e6f, 0x0d65, 0x0c6c,
	0x036c, 0x0265, 0x016f, 0x0066, 0x076a, 0x0663, 0x0569, 0x0460,
	0x0ca0, 0x0da9, 0x0ea3, 0x0faa, 0x08a6, 0x09af, 0x0aa5, 0x0bac,
	0x04ac, 0x05a5, 0x06af, 0x07a6, 0x00aa, 0x01a3, 0x02a9, 0x03a0,
	0x0d30, 0x0c39, 0x0f33, 0x0e3a, 0x0936, 0x083f, 0x0b35, 0x0a3c,
	0x053c, 0x0435, 0x073f, 0x0636, 0x013a, 0x0033, 0x0339, 0x0230,
	0x0e90, 0x0f99, 0x0c93, 0x0d9a, 0x0a96, 0x0b9f, 0x0895, 0x099c,
	0x069c, 0x0795, 0x049f, 0x0596, 0x029a, 0x0393, 0x0099, 0x0190,
	0x0f00, 0x0e09, 0x0d03, 0x0c0a, 0x0b06, 0x0a0f, 0x0905, 0x080c,
	0x070c, 0x0605, 0x050f, 0x0406, 0x030a, 0x0203, 0x0109, 0x0000,
}

// specify the edges used to create the triangle(s)
var mcTriangleTable = [256][]int{
	{},
	{0, 8, 3},
	{0, 1, 9},
	{1, 8, 3, 9, 8, 1},
	{1, 2, 10},
	{0, 8, 3, 1, 2, 10},
	{9, 2, 10, 0, 2, 9},
	{2, 8, 3, 2, 10, 8, 10, 9, 8},
	{3, 11, 2},
	{0, 11, 2, 8, 11, 0},
	{1, 9, 0, 2, 3, 11},
	{1, 11, 2, 1, 9, 11, 9, 8, 11},
	{3, 10, 1, 11, 10, 3},
	{0, 10, 1, 0, 8, 10, 8, 11, 10},
	{3, 9, 0, 3, 11, 9, 11, 10, 9},
	{9, 8, 10, 10, 8, 11},
	{4, 7, 8},
	{4, 3, 0, 7, 3, 4},
	{0, 1, 9, 8, 4, 7},
	{4, 1, 9, 4, 7, 1, 7, 3, 1},
	{1, 2, 10, 8, 4, 7},
	{3, 4, 7, 3, 0, 4, 1, 2, 10},
	{9, 2, 10, 9, 0, 2, 8, 4, 7},
	{2, 10, 9, 2, 9, 7, 2, 7, 3, 7, 9, 4},
	{8, 4, 7, 3, 11, 2},
	{11, 4, 7, 11, 2, 4, 2, 0, 4},
	{9, 0, 1, 8, 4, 7, 2, 3, 11},
	{4, 7, 11, 9, 4, 11, 9, 11, 2, 9, 2, 1},
	{3, 10, 1, 3, 11, 10, 7, 8, 4},
	{1, 11, 10, 1, 4, 11, 1, 0, 4, 7, 11, 4},
	{4, 7, 8, 9, 0, 11, 9, 11, 10, 11, 0, 3},
	{4, 7, 11, 4, 11, 9, 9, 11, 10},
	{9, 5, 4},
	{9, 5, 4, 0, 8, 3},
	{0, 5, 4, 1, 5, 0},
	{8, 5, 4, 8, 3, 5, 3, 1, 5},
	{1, 2, 10, 9, 5, 4},
	{3, 0, 8, 1, 2, 10, 4, 9, 5},
	{5, 2, 10, 5, 4, 2, 4, 0, 2},
	{2, 10, 5, 3, 2, 5, 3, 5, 4, 3, 4, 8},
	{9, 5, 4, 2, 3, 11},
	{0, 11, 2, 0, 8, 11, 4, 9, 5},
	{0, 5, 4, 0, 1, 5, 2, 3, 11},
	{2, 1, 5, 2, 5, 8, 2, 8, 11, 4, 8, 5},
	{10, 3, 11, 10, 1, 3, 9, 5, 4},
	{4, 9, 5, 0, 8, 1, 8, 10, 1, 8, 11, 10},
	{5, 4, 0, 5, 0, 11, 5, 11, 10, 11, 0, 3},
	{5, 4, 8, 5, 8, 10, 10, 8, 11},
	{9, 7, 8, 5, 7, 9},
	{9, 3, 0, 9, 5, 3, 5, 7, 3},
	{0, 7, 8, 0, 1, 7, 1, 5, 7},
	{1, 5, 3, 3, 5, 7},
	{9, 7, 8, 9, 5, 7, 10, 1, 2},
	{10, 1, 2, 9, 5, 0, 5, 3, 0, 5, 7, 3},
	{8, 0, 2, 8, 2, 5, 8, 5, 7, 10, 5, 2},
	{2, 10, 5, 2, 5, 3, 3, 5, 7},
	{7, 9, 5, 7, 8, 9, 3, 11, 2},
	{9, 5, 7, 9, 7, 2, 9, 2, 0, 2, 7, 11},
	{2, 3, 11, 0, 1, 8, 1, 7, 8, 1, 5, 7},
	{11, 2, 1, 11, 1, 7, 7, 1, 5},
	{9, 5, 8, 8, 5, 7, 10, 1, 3, 10, 3, 11},
	{5, 7, 0, 5, 0, 9, 7, 11, 0, 1, 0, 10, 11, 10, 0},
	{11, 10, 0, 11, 0, 3, 10, 5, 0, 8, 0, 7, 5, 7, 0},
	{11, 10, 5, 7, 11, 5},
	{10, 6, 5},
	{0, 8, 3, 5, 10, 6},
	{9, 0, 1, 5, 10, 6},
	{1, 8, 3, 1, 9, 8, 5, 10, 6},
	{1, 6, 5, 2, 6, 1},
	{1, 6, 5, 1, 2, 6, 3, 0, 8},
	{9, 6, 5, 9, 0, 6, 0, 2, 6},
	{5, 9, 8, 5, 8, 2, 5, 2, 6, 3, 2, 8},
	{2, 3, 11, 10, 6, 5},
	{11, 0, 8, 11, 2, 0, 10, 6, 5},
	{0, 1, 9, 2, 3, 11, 5, 10, 6},
	{5, 10, 6, 1, 9, 2, 9, 11, 2, 9, 8, 11},
	{6, 3, 11, 6, 5, 3, 5, 1, 3},
	{0, 8, 11, 0, 11, 5, 0, 5, 1, 5, 11, 6},
	{3, 11, 6, 0, 3, 6, 0, 6, 5, 0, 5, 9},
	{6, 5, 9, 6, 9, 11, 11, 9, 8},
	{5, 10, 6, 4, 7, 8},
	{4, 3, 0, 4, 7, 3, 6, 5, 10},
	{1, 9, 0, 5, 10, 6, 8, 4, 7},
	{10, 6, 5, 1, 9, 7, 1, 7, 3, 7, 9, 4},
	{6, 1, 2, 6, 5, 1, 4, 7, 8},
	{1, 2, 5, 5, 2, 6, 3, 0, 4, 3, 4, 7},
	{8, 4, 7, 9, 0, 5, 0, 6, 5, 0, 2, 6},
	{7, 3, 9, 7, 9, 4, 3, 2, 9, 5, 9, 6, 2, 6, 9},
	{3, 11, 2, 7, 8, 4, 10, 6, 5},
	{5, 10, 6, 4, 7, 2, 4, 2, 0, 2, 7, 11},
	{0, 1, 9, 4, 7, 8, 2, 3, 11, 5, 10, 6},
	{9, 2, 1, 9, 11, 2, 9, 4, 11, 7, 11, 4, 5, 10, 6},
	{8, 4, 7, 3, 11, 5, 3, 5, 1, 5, 11, 6},
	{5, 1, 11, 5, 11, 6, 1, 0, 11, 7, 11, 4, 0, 4, 11},
	{0, 5, 9, 0, 6, 5, 0, 3, 6, 11, 6, 3, 8, 4, 7},
	{6, 5, 9, 6, 9, 11, 4, 7, 9, 7, 11, 9},
	{10, 4, 9, 6, 4, 10},
	{4, 10, 6, 4, 9, 10, 0, 8, 3},
	{10, 0, 1, 10, 6, 0, 6, 4, 0},
	{8, 3, 1, 8, 1, 6, 8, 6, 4, 6, 1, 10},
	{1, 4, 9, 1, 2, 4, 2, 6, 4},
	{3, 0, 8, 1, 2, 9, 2, 4, 9, 2, 6, 4},
	{0, 2, 4, 4, 2, 6},
	{8, 3, 2, 8, 2, 4, 4, 2, 6},
	{10, 4, 9, 10, 6, 4, 11, 2, 3},
	{0, 8, 2, 2, 8, 11, 4, 9, 10, 4, 10, 6},
	{3, 11, 2, 0, 1, 6, 0, 6, 4, 6, 1, 10},
	{6, 4, 1, 6, 1, 10, 4, 8, 1, 2, 1, 11, 8, 11, 1},
	{9, 6, 4, 9, 3, 6, 9, 1, 3, 11, 6, 3},
	{8, 11, 1, 8, 1, 0, 11, 6, 1, 9, 1, 4, 6, 4, 1},
	{3, 11, 6, 3, 6, 0, 0, 6, 4},
	{6, 4, 8, 11, 6, 8},
	{7, 10, 6, 7, 8, 10, 8, 9, 10},
	{0, 7, 3, 0, 10, 7, 0, 9, 10, 6, 7, 10},
	{10, 6, 7, 1, 10, 7, 1, 7, 8, 1, 8, 0},
	{10, 6, 7, 10, 7, 1, 1, 7, 3},
	{1, 2, 6, 1, 6, 8, 1, 8, 9, 8, 6, 7},
	{2, 6, 9, 2, 9, 1, 6, 7, 9, 0, 9, 3, 7, 3, 9},
	{7, 8, 0, 7, 0, 6, 6, 0, 2},
	{7, 3, 2, 6, 7, 2},
	{2, 3, 11, 10, 6, 8, 10, 8, 9, 8, 6, 7},
	{2, 0, 7, 2, 7, 11, 0, 9, 7, 6, 7, 10, 9, 10, 7},
	{1, 8, 0, 1, 7, 8, 1, 10, 7, 6, 7, 10, 2, 3, 11},
	{11, 2, 1, 11, 1, 7, 10, 6, 1, 6, 7, 1},
	{8, 9, 6, 8, 6, 7, 9, 1, 6, 11, 6, 3, 1, 3, 6},
	{0, 9, 1, 11, 6, 7},
	{7, 8, 0, 7, 0, 6, 3, 11, 0, 11, 6, 0},
	{7, 11, 6},
	{7, 6, 11},
	{3, 0, 8, 11, 7, 6},
	{0, 1, 9, 11, 7, 6},
	{8, 1, 9, 8, 3, 1, 11, 7, 6},
	{10, 1, 2, 6, 11, 7},
	{1, 2, 10, 3, 0, 8, 6, 11, 7},
	{2, 9, 0, 2, 10, 9, 6, 11, 7},
	{6, 11, 7, 2, 10, 3, 10, 8, 3, 10, 9, 8},
	{7, 2, 3, 6, 2, 7},
	{7, 0, 8, 7, 6, 0, 6, 2, 0},
	{2, 7, 6, 2, 3, 7, 0, 1, 9},
	{1, 6, 2, 1, 8, 6, 1, 9, 8, 8, 7, 6},
	{10, 7, 6, 10, 1, 7, 1, 3, 7},
	{10, 7, 6, 1, 7, 10, 1, 8, 7, 1, 0, 8},
	{0, 3, 7, 0, 7, 10, 0, 10, 9, 6, 10, 7},
	{7, 6, 10, 7, 10, 8, 8, 10, 9},
	{6, 8, 4, 11, 8, 6},
	{3, 6, 11, 3, 0, 6, 0, 4, 6},
	{8, 6, 11, 8, 4, 6, 9, 0, 1},
	{9, 4, 6, 9, 6, 3, 9, 3, 1, 11, 3, 6},
	{6, 8, 4, 6, 11, 8, 2, 10, 1},
	{1, 2, 10, 3, 0, 11, 0, 6, 11, 0, 4, 6},
	{4, 11, 8, 4, 6, 11, 0, 2, 9, 2, 10, 9},
	{10, 9, 3, 10, 3, 2, 9, 4, 3, 11, 3, 6, 4, 6, 3},
	{8, 2, 3, 8, 4, 2, 4, 6, 2},
	{0, 4, 2, 4, 6, 2},
	{1, 9, 0, 2, 3, 4, 2, 4, 6, 4, 3, 8},
	{1, 9, 4, 1, 4, 2, 2, 4, 6},
	{8, 1, 3, 8, 6, 1, 8, 4, 6, 6, 10, 1},
	{10, 1, 0, 10, 0, 6, 6, 0, 4},
	{4, 6, 3, 4, 3, 8, 6, 10, 3, 0, 3, 9, 10, 9, 3},
	{10, 9, 4, 6, 10, 4},
	{4, 9, 5, 7, 6, 11},
	{0, 8, 3, 4, 9, 5, 11, 7, 6},
	{5, 0, 1, 5, 4, 0, 7, 6, 11},
	{11, 7, 6, 8, 3, 4, 3, 5, 4, 3, 1, 5},
	{9, 5, 4, 10, 1, 2, 7, 6, 11},
	{6, 11, 7, 1, 2, 10, 0, 8, 3, 4, 9, 5},
	{7, 6, 11, 5, 4, 10, 4, 2, 10, 4, 0, 2},
	{3, 4, 8, 3, 5, 4, 3, 2, 5, 10, 5, 2, 11, 7, 6},
	{7, 2, 3, 7, 6, 2, 5, 4, 9},
	{9, 5, 4, 0, 8, 6, 0, 6, 2, 6, 8, 7},
	{3, 6, 2, 3, 7, 6, 1, 5, 0, 5, 4, 0},
	{6, 2, 8, 6, 8, 7, 2, 1, 8, 4, 8, 5, 1, 5, 8},
	{9, 5, 4, 10, 1, 6, 1, 7, 6, 1, 3, 7},
	{1, 6, 10, 1, 7, 6, 1, 0, 7, 8, 7, 0, 9, 5, 4},
	{4, 0, 10, 4, 10, 5, 0, 3, 10, 6, 10, 7, 3, 7, 10},
	{7, 6, 10, 7, 10, 8, 5, 4, 10, 4, 8, 10},
	{6, 9, 5, 6, 11, 9, 11, 8, 9},
	{3, 6, 11, 0, 6, 3, 0, 5, 6, 0, 9, 5},
	{0, 11, 8, 0, 5, 11, 0, 1, 5, 5, 6, 11},
	{6, 11, 3, 6, 3, 5, 5, 3, 1},
	{1, 2, 10, 9, 5, 11, 9, 11, 8, 11, 5, 6},
	{0, 11, 3, 0, 6, 11, 0, 9, 6, 5, 6, 9, 1, 2, 10},
	{11, 8, 5, 11, 5, 6, 8, 0, 5, 10, 5, 2, 0, 2, 5},
	{6, 11, 3, 6, 3, 5, 2, 10, 3, 10, 5, 3},
	{5, 8, 9, 5, 2, 8, 5, 6, 2, 3, 8, 2},
	{9, 5, 6, 9, 6, 0, 0, 6, 2},
	{1, 5, 8, 1, 8, 0, 5, 6, 8, 3, 8, 2, 6, 2, 8},
	{1, 5, 6, 2, 1, 6},
	{1, 3, 6, 1, 6, 10, 3, 8, 6, 5, 6, 9, 8, 9, 6},
	{10, 1, 0, 10, 0, 6, 9, 5, 0, 5, 6, 0},
	{0, 3, 8, 5, 6, 10},
	{10, 5, 6},
	{11, 5, 10, 7, 5, 11},
	{11, 5, 10, 11, 7, 5, 8, 3, 0},
	{5, 11, 7, 5, 10, 11, 1, 9, 0},
	{10, 7, 5, 10, 11, 7, 9, 8, 1, 8, 3, 1},
	{11, 1, 2, 11, 7, 1, 7, 5, 1},
	{0, 8, 3, 1, 2, 7, 1, 7, 5, 7, 2, 11},
	{9, 7, 5, 9, 2, 7, 9, 0, 2, 2, 11, 7},
	{7, 5, 2, 7, 2, 11, 5, 9, 2, 3, 2, 8, 9, 8, 2},
	{2, 5, 10, 2, 3, 5, 3, 7, 5},
	{8, 2, 0, 8, 5, 2, 8, 7, 5, 10, 2, 5},
	{9, 0, 1, 5, 10, 3, 5, 3, 7, 3, 10, 2},
	{9, 8, 2, 9, 2, 1, 8, 7, 2, 10, 2, 5, 7, 5, 2},
	{1, 3, 5, 3, 7, 5},
	{0, 8, 7, 0, 7, 1, 1, 7, 5},
	{9, 0, 3, 9, 3, 5, 5, 3, 7},
	{9, 8, 7, 5, 9, 7},
	{5, 8, 4, 5, 10, 8, 10, 11, 8},
	{5, 0, 4, 5, 11, 0, 5, 10, 11, 11, 3, 0},
	{0, 1, 9, 8, 4, 10, 8, 10, 11, 10, 4, 5},
	{10, 11, 4, 10, 4, 5, 11, 3, 4, 9, 4, 1, 3, 1, 4},
	{2, 5, 1, 2, 8, 5, 2, 11, 8, 4, 5, 8},
	{0, 4, 11, 0, 11, 3, 4, 5, 11, 2, 11, 1, 5, 1, 11},
	{0, 2, 5, 0, 5, 9, 2, 11, 5, 4, 5, 8, 11, 8, 5},
	{9, 4, 5, 2, 11, 3},
	{2, 5, 10, 3, 5, 2, 3, 4, 5, 3, 8, 4},
	{5, 10, 2, 5, 2, 4, 4, 2, 0},
	{3, 10, 2, 3, 5, 10, 3, 8, 5, 4, 5, 8, 0, 1, 9},
	{5, 10, 2, 5, 2, 4, 1, 9, 2, 9, 4, 2},
	{8, 4, 5, 8, 5, 3, 3, 5, 1},
	{0, 4, 5, 1, 0, 5},
	{8, 4, 5, 8, 5, 3, 9, 0, 5, 0, 3, 5},
	{9, 4, 5},
	{4, 11, 7, 4, 9, 11, 9, 10, 11},
	{0, 8, 3, 4, 9, 7, 9, 11, 7, 9, 10, 11},
	{1, 10, 11, 1, 11, 4, 1, 4, 0, 7, 4, 11},
	{3, 1, 4, 3, 4, 8, 1, 10, 4, 7, 4, 11, 10, 11, 4},
	{4, 11, 7, 9, 11, 4, 9, 2, 11, 9, 1, 2},
	{9, 7, 4, 9, 11, 7, 9, 1, 11, 2, 11, 1, 0, 8, 3},
	{11, 7, 4, 11, 4, 2, 2, 4, 0},
	{11, 7, 4, 11, 4, 2, 8, 3, 4, 3, 2, 4},
	{2, 9, 10, 2, 7, 9, 2, 3, 7, 7, 4, 9},
	{9, 10, 7, 9, 7, 4, 10, 2, 7, 8, 7, 0, 2, 0, 7},
	{3, 7, 10, 3, 10, 2, 7, 4, 10, 1, 10, 0, 4, 0, 10},
	{1, 10, 2, 8, 7, 4},
	{4, 9, 1, 4, 1, 7, 7, 1, 3},
	{4, 9, 1, 4, 1, 7, 0, 8, 1, 8, 7, 1},
	{4, 0, 3, 7, 4, 3},
	{4, 8, 7},
	{9, 10, 8, 10, 11, 8},
	{3, 0, 9, 3, 9, 11, 11, 9, 10},
	{0, 1, 10, 0, 10, 8, 8, 10, 11},
	{3, 1, 10, 11, 3, 10},
	{1, 2, 11, 1, 11, 9, 9, 11, 8},
	{3, 0, 9, 3, 9, 11, 1, 2, 9, 2, 11, 9},
	{0, 2, 11, 8, 0, 11},
	{3, 2, 11},
	{2, 3, 8, 2, 8, 10, 10, 8, 9},
	{9, 10, 2, 0, 9, 2},
	{2, 3, 8, 2, 8, 10, 0, 1, 8, 1, 10, 8},
	{1, 10, 2},
	{1, 3, 8, 9, 1, 8},
	{0, 9, 1},
	{0, 3, 8},
	{},
}

//-----------------------------------------------------------------------------
