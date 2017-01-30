//-----------------------------------------------------------------------------
/*

2D Signed Distance Functions

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

type SDF2 interface {
	Evaluate(p V2) float64
	BoundingBox() Box2
}

//-----------------------------------------------------------------------------
// Basic SDF Functions

func sdf_box2d(p, s V2) float64 {
	d := p.Abs().Sub(s)
	return d.Max(V2{0, 0}).Length() + Min(d.MaxComponent(), 0)
}

/* alternate function - probably faster
func sdf_box2d(p, s V2) float64 {
	p = p.Abs()
	d := p.Sub(s)
	k := s.Y - s.X
	if d.X > 0 && d.Y > 0 {
		return d.Length()
	}
	if p.Y-p.X > k {
		return d.Y
	}
	return d.X
}
*/

//-----------------------------------------------------------------------------
// 2D Circle

type CircleSDF2 struct {
	radius float64
}

func NewCircleSDF2(radius float64) SDF2 {
	return &CircleSDF2{radius}
}

func (s *CircleSDF2) Evaluate(p V2) float64 {
	return p.Length() - s.radius
}

func (s *CircleSDF2) BoundingBox() Box2 {
	d := V2{s.radius, s.radius}
	return Box2{d.Negate(), d}
}

//-----------------------------------------------------------------------------
// 2D Normal Box

type BoxSDF2 struct {
	Size V2
}

func NewBoxSDF2(size V2) SDF2 {
	// note: store a modified size
	return &BoxSDF2{size.MulScalar(0.5)}
}

func (s *BoxSDF2) Evaluate(p V2) float64 {
	return sdf_box2d(p, s.Size)
}

func (s *BoxSDF2) BoundingBox() Box2 {
	return Box2{s.Size.Negate(), s.Size}
}

//-----------------------------------------------------------------------------
// 2D Rounded Box

type RoundedBoxSDF2 struct {
	Size   V2
	Radius float64
}

func NewRoundedBoxSDF2(size V2, radius float64) SDF2 {
	// note: store a modified size
	return &RoundedBoxSDF2{size.MulScalar(0.5).SubScalar(radius), radius}
}

func (s *RoundedBoxSDF2) Evaluate(p V2) float64 {
	return sdf_box2d(p, s.Size) - s.Radius
}

func (s *RoundedBoxSDF2) BoundingBox() Box2 {
	d := s.Size.AddScalar(s.Radius)
	return Box2{d.Negate(), d}
}

//-----------------------------------------------------------------------------
// 2D Polygon

type PolySDF2 struct {
	line []Line2 // line segments
	bb   Box2    // bounding box
}

func NewPolySDF2(points []V2) SDF2 {
	n := len(points)
	if n < 3 {
		return nil
	}
	// build the line segments
	s := PolySDF2{}
	for i := 0; i < n-1; i++ {
		s.line = append(s.line, NewLine2(points[i], points[i+1]))
	}
	// close the loop
	s.line = append(s.line, NewLine2(points[n-1], points[0]))
	// build the bounding box
	l := s.line[0]
	a := l.A.Min(l.B)
	b := l.A.Max(l.B)
	for _, l := range s.line {
		a = a.Min(l.A.Min(l.B))
		b = b.Max(l.A.Max(l.B))
	}
	s.bb = Box2{a, b}
	return &s
}

func (s *PolySDF2) Evaluate(p V2) float64 {
	dd := math.MaxFloat64 // d^2 to polygon (>0)
	wn := 0               // winding number (inside/outside)
	// iterate over the line segments
	for _, l := range s.line {
		// record the minimum distance squared
		x := l.Distance2(p)
		if Abs(x) < dd {
			dd = Abs(x)
		}

		// Is the point in the polygon?
		// See: http://geomalgorithms.com/a03-_inclusion.html
		if l.A.Y <= p.Y {
			if l.B.Y > p.Y { // upward crossing
				if x < 0 { // p is to the left of the line segment
					wn += 1 // up intersect
				}
			}
		} else {
			if l.B.Y <= p.Y { // downward crossing
				if x > 0 { // p is to the right of the line segment
					wn -= 1 // down intersect
				}
			}
		}
	}
	// normalise d*d to d
	d := math.Sqrt(dd)
	if wn != 0 {
		// p is inside the polygon
		return -d
	}
	// p is outside the polygon
	return d
}

func (s *PolySDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Transform SDF2

type TransformSDF2 struct {
	Sdf     SDF2
	Matrix  M33
	Inverse M33
}

func NewTransformSDF2(sdf SDF2, matrix M33) SDF2 {
	return &TransformSDF2{sdf, matrix, matrix.Inverse()}
}

func (s *TransformSDF2) Evaluate(p V2) float64 {
	q := s.Inverse.MulPosition(p)
	return s.Sdf.Evaluate(q)
}

func (s *TransformSDF2) BoundingBox() Box2 {
	return s.Matrix.MulBox(s.Sdf.BoundingBox())
}

//-----------------------------------------------------------------------------
// ArraySDF2: Create an X by Y array of a given SDF2
// num = the array size
// size = the step size

type ArraySDF2 struct {
	sdf  SDF2
	num  V2i
	step V2
	min  MinFunc
	k    float64
	bb   Box2
}

func NewArraySDF2(sdf SDF2, num V2i, step V2) SDF2 {
	// check the number of x/y steps
	if num[0] <= 0 || num[1] <= 0 {
		return nil
	}
	s := ArraySDF2{}
	s.sdf = sdf
	s.num = num
	s.step = step
	s.min = NormalMin
	// work out the bounding box
	bb0 := sdf.BoundingBox()
	bb1 := bb0.Translate(step.Mul(num.SubScalar(1).ToV2()))
	s.bb = bb0.Extend(bb1)
	return &s
}

// set the minimum function to control blending
func (s *ArraySDF2) SetMin(min MinFunc, k float64) {
	s.min = min
	s.k = k
}

func (s *ArraySDF2) Evaluate(p V2) float64 {
	d := math.MaxFloat64
	for j := 0; j < s.num[0]; j++ {
		for k := 0; k < s.num[1]; k++ {
			x := p.Sub(V2{float64(j) * s.step.X, float64(k) * s.step.Y})
			d = s.min(d, s.sdf.Evaluate(x), s.k)
		}
	}
	return d
}

func (s *ArraySDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

type RotateSDF2 struct {
	sdf  SDF2
	num  int
	step M33
	min  MinFunc
	k    float64
	bb   Box2
}

func NewRotateSDF2(sdf SDF2, num int, step M33) SDF2 {
	// check the number of steps
	if num <= 0 {
		return nil
	}
	s := RotateSDF2{}
	s.sdf = sdf
	s.num = num
	s.step = step.Inverse()
	s.min = NormalMin
	// work out the bounding box
	// TODO: It could be smaller based on num * step
	bb := sdf.BoundingBox()
	size := bb.Size()
	tr := bb.Max
	bl := bb.Min
	br := bl.Add(V2{size.X, 0})
	tl := bl.Add(V2{0, size.Y})
	r := math.Sqrt(Max(Max(tl.Length2(), tr.Length2()), Max(bl.Length2(), br.Length2())))
	s.bb = Box2{V2{-r, -r}, V2{r, r}}
	return &s
}

func (s *RotateSDF2) Evaluate(p V2) float64 {
	d := math.MaxFloat64
	rot := Identity2d()
	for i := 0; i < s.num; i++ {
		x := rot.MulPosition(p)
		d = s.min(d, s.sdf.Evaluate(x), s.k)
		rot = rot.Mul(s.step)
	}
	return d
}

// set the minimum function to control blending
func (s *RotateSDF2) SetMin(min MinFunc, k float64) {
	s.min = min
	s.k = k
}

func (s *RotateSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
