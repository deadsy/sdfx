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

// Solid of Revolution, SDF2 to SDF3
type SorSDF3 struct {
	sdf   SDF2
	theta float64 // angle for partial revolutions
	norm  V2      // pre-calculated normal to theta line
	bb    Box3
}

// Return an SDF3 for a solid of revolution.
func NewSorThetaSDF3(sdf SDF2, theta float64) SDF3 {
	s := SorSDF3{}
	s.sdf = sdf
	// normalize theta
	s.theta = math.Mod(Abs(theta), TAU)
	sin := math.Sin(s.theta)
	cos := math.Cos(s.theta)
	// pre-calculate the normal to the theta line
	s.norm = V2{-sin, cos}
	// work out the bounding box
	var vset V2Set
	if s.theta == 0 {
		vset = []V2{V2{1, 1}, V2{-1, -1}}
	} else {
		vset = []V2{V2{0, 0}, V2{1, 0}, V2{cos, sin}}
		if s.theta > 0.5*PI {
			vset = append(vset, V2{0, 1})
		}
		if s.theta > PI {
			vset = append(vset, V2{-1, 0})
		}
		if s.theta > 1.5*PI {
			vset = append(vset, V2{0, -1})
		}
	}
	bb := s.sdf.BoundingBox()
	l := Max(Abs(bb.Min.X), Abs(bb.Max.X))
	vmin := vset.Min().MulScalar(l)
	vmax := vset.Max().MulScalar(l)
	s.bb = Box3{V3{vmin.X, vmin.Y, bb.Min.Y}, V3{vmax.X, vmax.Y, bb.Max.Y}}
	return &s
}

// Return an SDF3 for a solid of revolution.
func NewSorSDF3(sdf SDF2) SDF3 {
	return NewSorThetaSDF3(sdf, 0)
}

// Return the minimum distance to a solid of revolution.
func (s *SorSDF3) Evaluate(p V3) float64 {
	x := math.Sqrt(p.X*p.X + p.Y*p.Y)
	a := s.sdf.Evaluate(V2{x, p.Z})
	b := a
	if s.theta != 0 {
		// combine two vertical planes to give an intersection wedge
		d := s.norm.Dot(V2{p.X, p.Y})
		if s.theta < PI {
			b = Max(-p.Y, d) // intersect
		} else {
			b = Min(-p.Y, d) // union
		}
	}
	// return the intersection
	return Max(a, b)
}

// Return the bounding box for a solid of revolution.
func (s *SorSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// Extrude, SDF2 to SDF3
type ExtrudeSDF3 struct {
	sdf     SDF2
	height  float64
	extrude ExtrudeFunc
	bb      Box3
}

// Linear Extrude
func NewExtrudeSDF3(sdf SDF2, height float64) SDF3 {
	s := ExtrudeSDF3{}
	s.sdf = sdf
	s.height = height / 2
	s.extrude = NormalExtrude
	// work out the bounding box
	bb := sdf.BoundingBox()
	s.bb = Box3{V3{bb.Min.X, bb.Min.Y, -s.height}, V3{bb.Max.X, bb.Max.Y, s.height}}
	return &s
}

// Twist Extrude - rotate by twist radians over the height of the extrusion
func NewTwistExtrudeSDF3(sdf SDF2, height, twist float64) SDF3 {
	s := ExtrudeSDF3{}
	s.sdf = sdf
	s.height = height / 2
	s.extrude = TwistExtrude(height, twist)
	// work out the bounding box
	bb := sdf.BoundingBox()
	l := bb.Max.Length()
	s.bb = Box3{V3{-l, -l, -s.height}, V3{l, l, s.height}}
	return &s
}

// Scale Extrude - scale over the height of the extrusion
func NewScaleExtrudeSDF3(sdf SDF2, height float64, scale V2) SDF3 {
	s := ExtrudeSDF3{}
	s.sdf = sdf
	s.height = height / 2
	s.extrude = ScaleExtrude(height, scale)
	// work out the bounding box
	bb := sdf.BoundingBox()
	bb = bb.Extend(Box2{bb.Min.Mul(scale), bb.Max.Mul(scale)})
	s.bb = Box3{V3{bb.Min.X, bb.Min.Y, -s.height}, V3{bb.Max.X, bb.Max.Y, s.height}}
	return &s
}

// Scale + Twist Extrude - scale and then twist over the height of the extrusion
func NewScaleTwistExtrudeSDF3(sdf SDF2, height, twist float64, scale V2) SDF3 {
	s := ExtrudeSDF3{}
	s.sdf = sdf
	s.height = height / 2
	s.extrude = ScaleTwistExtrude(height, twist, scale)
	// work out the bounding box
	bb := sdf.BoundingBox()
	bb = bb.Extend(Box2{bb.Min.Mul(scale), bb.Max.Mul(scale)})
	l := bb.Max.Length()
	s.bb = Box3{V3{-l, -l, -s.height}, V3{l, l, s.height}}
	return &s
}

func (s *ExtrudeSDF3) Evaluate(p V3) float64 {
	// sdf for the projected 2d surface
	a := s.sdf.Evaluate(s.extrude(p))
	// sdf for the extrusion region: z = [-height, height]
	b := Abs(p.Z) - s.height
	// return the intersection
	return Max(a, b)
}

// Set the evaluation function to control extrusion.
func (s *ExtrudeSDF3) SetExtrude(extrude ExtrudeFunc) {
	s.extrude = extrude
}

func (s *ExtrudeSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// 3D Box
type BoxSDF3 struct {
	size  V3
	round float64
	bb    Box3
}

// Return an SDF3 for a box (rounded corners with round > 0).
func NewBoxSDF3(size V3, round float64) SDF3 {
	size = size.MulScalar(0.5)
	s := BoxSDF3{}
	s.size = size.SubScalar(round)
	s.round = round
	s.bb = Box3{size.Negate(), size}
	return &s
}

// Return the minimum distance to a box.
func (s *BoxSDF3) Evaluate(p V3) float64 {
	return sdf_box3d(p, s.size) - s.round
}

// Return the bounding box for a box.
func (s *BoxSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// 3D Sphere
type SphereSDF3 struct {
	radius float64
	bb     Box3
}

// Return an SDF3 for a sphere.
func NewSphereSDF3(radius float64) SDF3 {
	s := SphereSDF3{}
	s.radius = radius
	d := V3{radius, radius, radius}
	s.bb = Box3{d.Negate(), d}
	return &s
}

// Return the minimum distance to a sphere.
func (s *SphereSDF3) Evaluate(p V3) float64 {
	return p.Length() - s.radius
}

// Return the bounding box for a sphere.
func (s *SphereSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// 3D Cylinder
type CylinderSDF3 struct {
	height float64
	radius float64
	round  float64
	bb     Box3
}

// Return an SDF3 for a cylinder (rounded edges with round > 0).
func NewCylinderSDF3(height, radius, round float64) SDF3 {
	s := CylinderSDF3{}
	s.height = (height / 2) - round
	s.radius = radius - round
	s.round = round
	d := V3{radius, radius, height / 2}
	s.bb = Box3{d.Negate(), d}
	return &s
}

// Return an SDF3 for a capsule.
func NewCapsuleSDF3(radius, height float64) SDF3 {
	return NewCylinderSDF3(radius, height, radius)
}

// Return the minimum distance to a cylinder.
func (s *CylinderSDF3) Evaluate(p V3) float64 {
	d := sdf_box2d(V2{V2{p.X, p.Y}.Length(), p.Z}, V2{s.radius, s.height})
	return d - s.round
}

// Return the bounding box for a cylinder.
func (s *CylinderSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Cylinders of the same radius and height at various x/y positions
// (E.g. drilling patterns) are useful enough to warrant their own SDF3 function.

// Multiple Cylinders
type MultiCylinderSDF3 struct {
	height    float64
	radius    float64
	positions V2Set
	bb        Box3
}

// Return an SDF3 for multiple cylinders.
func NewMultiCylinderSDF3(height, radius float64, positions V2Set) SDF3 {
	s := MultiCylinderSDF3{}
	s.height = height / 2
	s.radius = radius
	s.positions = positions
	// work out the bounding box
	pmin := positions.Min().Sub(V2{radius, radius})
	pmax := positions.Max().Add(V2{radius, radius})
	s.bb = Box3{V3{pmin.X, pmin.Y, -height / 2}, V3{pmax.X, pmax.Y, height / 2}}
	return &s
}

// Return the minimum distance to multiple cylinders.
func (s *MultiCylinderSDF3) Evaluate(p V3) float64 {
	d := math.MaxFloat64
	for _, posn := range s.positions {
		l := V2{p.X, p.Y}.Sub(posn).Length()
		d = Min(d, sdf_box2d(V2{l, p.Z}, V2{s.radius, s.height}))
	}
	return d
}

// Return the bounding box for multiple cylinders.
func (s *MultiCylinderSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// Truncated Cone
type ConeSDF3 struct {
	r0     float64 // base radius
	r1     float64 // top radius
	height float64 // half height
	round  float64 // rounding offset
	u      V2      // normalized cone slope vector
	n      V2      // normal to cone slope (points outward)
	l      float64 // length of cone slope
	bb     Box3    // bounding box
}

// Return a new trucated cone (round > 0 gives rounded edges)
func NewConeSDF3(height, r0, r1, round float64) SDF3 {
	s := ConeSDF3{}
	s.height = (height / 2) - round
	s.round = round
	// cone slope vector, normal and length
	s.u = V2{r1 - r0, height}.Normalize()
	s.n = V2{s.u.Y, -s.u.X}
	s.l = V2{r0 - r1, 2 * s.height}.Length()
	// inset the radii because of the rounding
	ofs := round / s.n.X
	s.r0 = r0 - ofs
	s.r1 = r1 - ofs
	// work out the bounding box
	r := Max(r0, r1)
	s.bb = Box3{V3{-r, -r, -height / 2}, V3{r, r, height / 2}}
	return &s
}

// Return the minimum distance to the trucated cone.
func (s *ConeSDF3) Evaluate(p V3) float64 {
	// convert to SoR 2d coordinates
	p2 := V2{V2{p.X, p.Y}.Length(), p.Z}
	// is p2 above the cone?
	if p2.Y >= s.height && p2.X <= s.r1 {
		return p2.Y - s.height - s.round
	}
	// is p2 below the cone?
	if p2.Y <= -s.height && p2.X <= s.r0 {
		return -p2.Y - s.height - s.round
	}
	// distance to slope line
	v := p2.Sub(V2{s.r0, -s.height})
	d_slope := v.Dot(s.n)
	// is p2 inside the cone?
	if (d_slope < 0) && (Abs(p2.Y) < s.height) {
		return -Min(-d_slope, s.height-Abs(p2.Y)) - s.round
	}
	// is p2 closest to the slope line?
	t := v.Dot(s.u)
	if t >= 0 && t <= s.l {
		return d_slope - s.round
	}
	// is p2 closest to the base radius vertex?
	if t < 0 {
		return v.Length() - s.round
	}
	// p2 is closest to the top radius vertex
	return p2.Sub(V2{s.r1, s.height}).Length() - s.round
}

// Return the bounding box for the trucated cone.
func (s *ConeSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Offset an SDF3 - add a constant to the distance function
// offset > 0, enlarges and adds rounding to convex corners of the SDF
// offset < 0, skeletonizes the SDF

type OffsetSDF3 struct {
	sdf    SDF3
	offset float64
	bb     Box3
}

func NewOffsetSDF3(sdf SDF3, offset float64) SDF3 {
	s := OffsetSDF3{}
	s.sdf = sdf
	s.offset = offset
	// work out the bounding box
	bb := sdf.BoundingBox()
	s.bb = NewBox3(bb.Center(), bb.Size().AddScalar(2*offset))
	return &s
}

func (s *OffsetSDF3) Evaluate(p V3) float64 {
	return s.sdf.Evaluate(p) - s.offset
}

func (s *OffsetSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Transform SDF3

type TransformSDF3 struct {
	sdf     SDF3
	matrix  M44
	inverse M44
	bb      Box3
}

func NewTransformSDF3(sdf SDF3, matrix M44) SDF3 {
	s := TransformSDF3{}
	s.sdf = sdf
	s.matrix = matrix
	s.inverse = matrix.Inverse()
	s.bb = matrix.MulBox(sdf.BoundingBox())
	return &s
}

func (s *TransformSDF3) Evaluate(p V3) float64 {
	return s.sdf.Evaluate(s.inverse.MulPosition(p))
}

func (s *TransformSDF3) BoundingBox() Box3 {
	return s.bb
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

// Intersection of SDF3s
type IntersectionSDF3 struct {
	s0  SDF3
	s1  SDF3
	max MaxFunc
	k   float64
	bb  Box3
}

// Return the intersection of two SDF3 objects, s0 with s1.
func NewIntersectionSDF3(s0, s1 SDF3) SDF3 {
	s := IntersectionSDF3{}
	s.s0 = s0
	s.s1 = s1
	s.max = NormalMax
	s.bb = s0.BoundingBox()
	return &s
}

// Return the minimum distance to the object.
func (s *IntersectionSDF3) Evaluate(p V3) float64 {
	return s.max(s.s0.Evaluate(p), s.s1.Evaluate(p), s.k)
}

// Set the maximum function to control blending.
func (s *IntersectionSDF3) SetMax(max MaxFunc, k float64) {
	s.max = max
	s.k = k
}

// Return the bounding box.
func (s *IntersectionSDF3) BoundingBox() Box3 {
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
