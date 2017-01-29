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
// Rotate SDF3

type RotateSDF3 struct {
	sdf      SDF3
	rotation M44
	num      int
}

func NewRotateSDF3(sdf SDF3, axis V3, theta float64, num int) SDF3 {
	s := RotateSDF3{}
	s.sdf = sdf
	s.rotation = Rotate3d(axis, theta)
	s.num = num
	return &s
}

func (s *RotateSDF3) Evaluate(p V3) float64 {
	return 0
}

func (s *RotateSDF3) BoundingBox() Box3 {
	return Box3{}
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
// 3D Normal Box

type BoxSDF3 struct {
	Size V3
}

func NewBoxSDF3(size V3) SDF3 {
	// note: store a modified size
	return &BoxSDF3{size.MulScalar(0.5)}
}

func (s *BoxSDF3) Evaluate(p V3) float64 {
	return sdf_box3d(p, s.Size)
}

func (s *BoxSDF3) BoundingBox() Box3 {
	return Box3{s.Size.Negate(), s.Size}
}

//-----------------------------------------------------------------------------
// 3D Rounded Box

type RoundedBoxSDF3 struct {
	Size   V3
	Radius float64
}

func NewRoundedBoxSDF3(size V3, radius float64) SDF3 {
	// note: store a modified size
	return &RoundedBoxSDF3{size.MulScalar(0.5).SubScalar(radius), radius}
}

func (s *RoundedBoxSDF3) Evaluate(p V3) float64 {
	return sdf_box3d(p, s.Size) - s.Radius
}

func (s *RoundedBoxSDF3) BoundingBox() Box3 {
	d := s.Size.AddScalar(s.Radius)
	return Box3{d.Negate(), d}
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
// Union of SDF3

type UnionSDF3 struct {
	s0  SDF3
	s1  SDF3
	min MinFunc
	k   float64
}

func NewUnionSDF3(s0, s1 SDF3) SDF3 {
	return &UnionSDF3{s0, s1, NormalMin, 0}
}

func NewUnionRoundSDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, RoundMin, k}
}

func NewUnionExpSDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, ExpMin, k}
}

func NewUnionPowSDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, PowMin, k}
}

func NewUnionPolySDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, PolyMin, k}
}

func NewUnionChamferSDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, ChamferMin, k}
}

func (s *UnionSDF3) Evaluate(p V3) float64 {
	a := s.s0.Evaluate(p)
	b := s.s1.Evaluate(p)
	return s.min(a, b, s.k)
}

func (s *UnionSDF3) BoundingBox() Box3 {
	bb := s.s0.BoundingBox()
	return bb.Extend(s.s1.BoundingBox())
}

//-----------------------------------------------------------------------------
