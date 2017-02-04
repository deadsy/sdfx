//-----------------------------------------------------------------------------
/*

3D Signed Distance Functions

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"

	"github.com/deadsy/pt/pt"
)

//-----------------------------------------------------------------------------

type SDF3 interface {
	Evaluate(p V3) float64
	BoundingBox() Box3
}

//-----------------------------------------------------------------------------
// Basic SDF Functions

func sdf_box3d(p, s V3) float64 {
	d := p.Abs().Sub(s)
	return d.Max(V3{0, 0, 0}).Length() + Min(d.MaxComponent(), 0)
}

//-----------------------------------------------------------------------------
// Create a pt.SDF from an SDF3

type PtSDF struct {
	Sdf SDF3
}

func NewPtSDF(sdf SDF3) pt.SDF {
	return &PtSDF{sdf}
}

func (s *PtSDF) Evaluate(p pt.Vector) float64 {
	return s.Sdf.Evaluate(V3{p.X, p.Y, p.Z})
}

func (s *PtSDF) BoundingBox() pt.Box {
	b := s.Sdf.BoundingBox()
	j := b.Min
	k := b.Max
	return pt.Box{Min: pt.Vector{X: j.X, Y: j.Y, Z: j.Z}, Max: pt.Vector{X: k.X, Y: k.Y, Z: k.Z}}
}

//-----------------------------------------------------------------------------
// Solid of Revolution, SDF2 -> SDF3

type SorSDF3 struct {
	Sdf   SDF2
	Theta float64 // angle for partial revolutions
	Norm  V2      // pre-calculated normal to theta line
}

func NewSorSDF3(sdf SDF2) SDF3 {
	return &SorSDF3{sdf, 0, V2{}}
}

func NewSorThetaSDF3(sdf SDF2, theta float64) SDF3 {
	// normalize theta
	theta = math.Mod(Abs(theta), TAU)
	// pre-calculate the normal to the theta line
	norm := V2{math.Sin(theta), -math.Cos(theta)}
	return &SorSDF3{sdf, theta, norm}
}

func (s *SorSDF3) Evaluate(p V3) float64 {
	x := math.Sqrt(p.X*p.X + p.Y*p.Y)
	a := s.Sdf.Evaluate(V2{x, p.Z})
	b := a
	if s.Theta != 0 {
		// combine two vertical planes to give an intersection wedge
		d := s.Norm.Dot(V2{p.X, p.Y})
		if s.Theta < PI {
			b = Max(p.Y, d) // intersect
		} else {
			b = Min(p.Y, d) // union
		}
	}
	// return the intersection
	return Max(a, b)
}

func (s *SorSDF3) BoundingBox() Box3 {
	// TODO - reduce the BB for theta != 0
	b := s.Sdf.BoundingBox()
	j := b.Min
	k := b.Max
	l := Max(Abs(j.X), Abs(k.X))
	return Box3{V3{-l, -l, j.Y}, V3{l, l, k.Y}}
}

//-----------------------------------------------------------------------------
// Extrude, SDF2 -> SDF3

type ExtrudeSDF3 struct {
	Sdf    SDF2
	Height float64
}

func NewExtrudeSDF3(sdf SDF2, height float64) SDF3 {
	return &ExtrudeSDF3{sdf, height}
}

func (s *ExtrudeSDF3) Evaluate(p V3) float64 {
	// sdf for the projected 2d surface
	a := s.Sdf.Evaluate(V2{p.X, p.Y})
	// sdf for the extrusion region: z = [0, height]
	b := Max(-p.Z, p.Z-s.Height)
	// return the intersection
	return Max(a, b)
}

func (s *ExtrudeSDF3) BoundingBox() Box3 {
	b := s.Sdf.BoundingBox()
	j := b.Min
	k := b.Max
	return Box3{V3{j.X, j.Y, 0}, V3{k.X, k.Y, s.Height}}
}

//-----------------------------------------------------------------------------

type CapsuleSDF3 struct {
	A, B   V3
	Radius float64
}

func NewCapsuleSDF3(a, b V3, radius float64) SDF3 {
	return &CapsuleSDF3{a, b, radius}
}

func (s *CapsuleSDF3) Evaluate(p V3) float64 {
	pa := p.Sub(s.A)
	ba := s.B.Sub(s.A)
	t := Clamp(pa.Dot(ba)/ba.Dot(ba), 0, 1)
	return pa.Sub(ba.MulScalar(t)).Length() - s.Radius
}

func (s *CapsuleSDF3) BoundingBox() Box3 {
	a := s.A.Min(s.B).SubScalar(s.Radius)
	b := s.A.Max(s.B).AddScalar(s.Radius)
	return Box3{a, b}
}

//-----------------------------------------------------------------------------
// 3D Box (rounded corners with radius > 0)

type BoxSDF3 struct {
	size   V3
	radius float64
	bb     Box3
}

func NewBoxSDF3(size V3, radius float64) SDF3 {
	size = size.MulScalar(0.5)
	s := BoxSDF3{}
	s.size = size.SubScalar(radius)
	s.radius = radius
	s.bb = Box3{size.Negate(), size}
	return &s
}

func (s *BoxSDF3) Evaluate(p V3) float64 {
	return sdf_box3d(p, s.size) - s.radius
}

func (s *BoxSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// 3D Sphere

type SphereSDF3 struct {
	Radius float64
}

func NewSphereSDF3(radius float64) SDF3 {
	return &SphereSDF3{radius}
}

func (s *SphereSDF3) Evaluate(p V3) float64 {
	return p.Length() - s.Radius
}

func (s *SphereSDF3) BoundingBox() Box3 {
	d := V3{s.Radius, s.Radius, s.Radius}
	return Box3{d.Negate(), d}
}

//-----------------------------------------------------------------------------
// Transform SDF3

type TransformSDF3 struct {
	Sdf     SDF3
	Matrix  M44
	Inverse M44
}

func NewTransformSDF3(sdf SDF3, matrix M44) SDF3 {
	return &TransformSDF3{sdf, matrix, matrix.Inverse()}
}

func (s *TransformSDF3) Evaluate(p V3) float64 {
	q := s.Inverse.MulPosition(p)
	return s.Sdf.Evaluate(q)
}

func (s *TransformSDF3) BoundingBox() Box3 {
	return s.Matrix.MulBox(s.Sdf.BoundingBox())
}

//-----------------------------------------------------------------------------
// Union of SDF3s

type UnionSDF3 struct {
	s0  SDF3
	s1  SDF3
	min MinFunc
	k   float64
	bb  Box3
}

// Return the union of two SDF3 objects.
func NewUnionSDF3(s0, s1 SDF3) SDF3 {
	s := UnionSDF3{}
	s.s0 = s0
	s.s1 = s1
	s.min = NormalMin
	s.bb = s0.BoundingBox().Extend(s1.BoundingBox())
	return &s
}

// Return the minimum distance to the object.
func (s *UnionSDF3) Evaluate(p V3) float64 {
	return s.min(s.s0.Evaluate(p), s.s1.Evaluate(p), s.k)
}

// Set the minimum function to control blending.
func (s *UnionSDF3) SetMin(min MinFunc, k float64) {
	s.min = min
	s.k = k
}

// Return the bounding box.
func (s *UnionSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Difference of SDF3s

type DifferenceSDF3 struct {
	s0  SDF3
	s1  SDF3
	max MaxFunc
	k   float64
	bb  Box3
}

// Return the difference of two SDF3 objects, s0 - s1.
func NewDifferenceSDF3(s0, s1 SDF3) SDF3 {
	s := DifferenceSDF3{}
	s.s0 = s0
	s.s1 = s1
	s.max = NormalMax
	s.bb = s0.BoundingBox()
	return &s
}

// Return the minimum distance to the object.
func (s *DifferenceSDF3) Evaluate(p V3) float64 {
	return s.max(s.s0.Evaluate(p), -s.s1.Evaluate(p), s.k)
}

// Set the maximum function to control blending.
func (s *DifferenceSDF3) SetMax(max MaxFunc, k float64) {
	s.max = max
	s.k = k
}

// Return the bounding box.
func (s *DifferenceSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// ArraySDF3: Create an X by Y by Z array of a given SDF3
// num = the array size
// size = the step size

type ArraySDF3 struct {
	sdf  SDF3
	num  V3i
	step V3
	min  MinFunc
	k    float64
	bb   Box3
}

func NewArraySDF3(sdf SDF3, num V3i, step V3) SDF3 {
	// check the number of steps
	if num[0] <= 0 || num[1] <= 0 || num[2] <= 0 {
		return nil
	}
	s := ArraySDF3{}
	s.sdf = sdf
	s.num = num
	s.step = step
	s.min = NormalMin
	// work out the bounding box
	bb0 := sdf.BoundingBox()
	bb1 := bb0.Translate(step.Mul(num.SubScalar(1).ToV3()))
	s.bb = bb0.Extend(bb1)
	return &s
}

// set the minimum function to control blending
func (s *ArraySDF3) SetMin(min MinFunc, k float64) {
	s.min = min
	s.k = k
}

func (s *ArraySDF3) Evaluate(p V3) float64 {
	d := math.MaxFloat64
	for j := 0; j < s.num[0]; j++ {
		for k := 0; k < s.num[1]; k++ {
			for l := 0; l < s.num[2]; l++ {
				x := p.Sub(V3{float64(j) * s.step.X, float64(k) * s.step.Y, float64(l) * s.step.Z})
				d = s.min(d, s.sdf.Evaluate(x), s.k)
			}
		}
	}
	return d
}

func (s *ArraySDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

type RotateSDF3 struct {
	sdf  SDF3
	num  int
	step M44
	min  MinFunc
	k    float64
	bb   Box3
}

func NewRotateSDF3(sdf SDF3, num int, step M44) SDF3 {
	// check the number of steps
	if num <= 0 {
		return nil
	}
	s := RotateSDF3{}
	s.sdf = sdf
	s.num = num
	s.step = step.Inverse()
	s.min = NormalMin
	// work out the bounding box
	v := sdf.BoundingBox().Vertices()
	bb_min := v[0]
	bb_max := v[0]
	for i := 0; i < s.num; i++ {
		bb_min = bb_min.Min(v.Min())
		bb_max = bb_max.Max(v.Max())
		v.MulVertices(step)
	}
	s.bb = Box3{bb_min, bb_max}
	return &s
}

// Return the minimum distance to the object.
func (s *RotateSDF3) Evaluate(p V3) float64 {
	d := math.MaxFloat64
	rot := Identity3d()
	for i := 0; i < s.num; i++ {
		x := rot.MulPosition(p)
		d = s.min(d, s.sdf.Evaluate(x), s.k)
		rot = rot.Mul(s.step)
	}
	return d
}

// Set the minimum function to control blending.
func (s *RotateSDF3) SetMin(min MinFunc, k float64) {
	s.min = min
	s.k = k
}

// Return the bounding box.
func (s *RotateSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
