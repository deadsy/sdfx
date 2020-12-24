//-----------------------------------------------------------------------------
/*

3D Signed Distance Functions

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"math"
)

//-----------------------------------------------------------------------------

// SDF3 is the interface to a 3d signed distance function object.
type SDF3 interface {
	Evaluate(p V3) float64
	BoundingBox() Box3
}

//-----------------------------------------------------------------------------
// Basic SDF Functions

/*
func sdfBox3d(p, s V3) float64 {
	d := p.Abs().Sub(s)
	return d.Max(V3{0, 0, 0}).Length() + Min(d.MaxComponent(), 0)
}
*/

func sdfBox3d(p, s V3) float64 {
	d := p.Abs().Sub(s)
	if d.X > 0 && d.Y > 0 && d.Z > 0 {
		return d.Length()
	}
	if d.X > 0 && d.Y > 0 {
		return V2{d.X, d.Y}.Length()
	}
	if d.X > 0 && d.Z > 0 {
		return V2{d.X, d.Z}.Length()
	}
	if d.Y > 0 && d.Z > 0 {
		return V2{d.Y, d.Z}.Length()
	}
	if d.X > 0 {
		return d.X
	}
	if d.Y > 0 {
		return d.Y
	}
	if d.Z > 0 {
		return d.Z
	}
	return d.MaxComponent()
}

//-----------------------------------------------------------------------------

// SorSDF3 solid of revolution, SDF2 to SDF3.
type SorSDF3 struct {
	sdf   SDF2
	theta float64 // angle for partial revolutions
	norm  V2      // pre-calculated normal to theta line
	bb    Box3
}

// RevolveTheta3D returns an SDF3 for a solid of revolution.
func RevolveTheta3D(sdf SDF2, theta float64) (SDF3, error) {
	if sdf == nil {
		return nil, nil
	}
	if theta < 0 {
		return nil, ErrMsg("theta < 0")
	}
	s := SorSDF3{}
	s.sdf = sdf
	// normalize theta
	s.theta = math.Mod(math.Abs(theta), Tau)
	sin := math.Sin(s.theta)
	cos := math.Cos(s.theta)
	// pre-calculate the normal to the theta line
	s.norm = V2{-sin, cos}
	// work out the bounding box
	var vset V2Set
	if s.theta == 0 {
		vset = []V2{{1, 1}, {-1, -1}}
	} else {
		vset = []V2{{0, 0}, {1, 0}, {cos, sin}}
		if s.theta > 0.5*Pi {
			vset = append(vset, V2{0, 1})
		}
		if s.theta > Pi {
			vset = append(vset, V2{-1, 0})
		}
		if s.theta > 1.5*Pi {
			vset = append(vset, V2{0, -1})
		}
	}
	bb := s.sdf.BoundingBox()
	l := math.Max(math.Abs(bb.Min.X), math.Abs(bb.Max.X))
	vmin := vset.Min().MulScalar(l)
	vmax := vset.Max().MulScalar(l)
	s.bb = Box3{V3{vmin.X, vmin.Y, bb.Min.Y}, V3{vmax.X, vmax.Y, bb.Max.Y}}
	return &s, nil
}

// Revolve3D returns an SDF3 for a solid of revolution.
func Revolve3D(sdf SDF2) (SDF3, error) {
	return RevolveTheta3D(sdf, 0)
}

// Evaluate returns the minimum distance to a solid of revolution.
func (s *SorSDF3) Evaluate(p V3) float64 {
	x := math.Sqrt(p.X*p.X + p.Y*p.Y)
	a := s.sdf.Evaluate(V2{x, p.Z})
	b := a
	if s.theta != 0 {
		// combine two vertical planes to give an intersection wedge
		d := s.norm.Dot(V2{p.X, p.Y})
		if s.theta < Pi {
			b = math.Max(-p.Y, d) // intersect
		} else {
			b = math.Min(-p.Y, d) // union
		}
	}
	// return the intersection
	return math.Max(a, b)
}

// BoundingBox returns the bounding box for a solid of revolution.
func (s *SorSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// ExtrudeSDF3 extrudes an SDF2 to an SDF3.
type ExtrudeSDF3 struct {
	sdf     SDF2
	height  float64
	extrude ExtrudeFunc
	bb      Box3
}

// Extrude3D does a linear extrude on an SDF3.
func Extrude3D(sdf SDF2, height float64) SDF3 {
	s := ExtrudeSDF3{}
	s.sdf = sdf
	s.height = height / 2
	s.extrude = NormalExtrude
	// work out the bounding box
	bb := sdf.BoundingBox()
	s.bb = Box3{V3{bb.Min.X, bb.Min.Y, -s.height}, V3{bb.Max.X, bb.Max.Y, s.height}}
	return &s
}

// TwistExtrude3D extrudes an SDF2 while rotating by twist radians over the height of the extrusion.
func TwistExtrude3D(sdf SDF2, height, twist float64) SDF3 {
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

// ScaleExtrude3D extrudes an SDF2 and scales it over the height of the extrusion.
func ScaleExtrude3D(sdf SDF2, height float64, scale V2) SDF3 {
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

// ScaleTwistExtrude3D extrudes an SDF2 and scales and twists it over the height of the extrusion.
func ScaleTwistExtrude3D(sdf SDF2, height, twist float64, scale V2) SDF3 {
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

// Evaluate returns the minimum distance to an extrusion.
func (s *ExtrudeSDF3) Evaluate(p V3) float64 {
	// sdf for the projected 2d surface
	a := s.sdf.Evaluate(s.extrude(p))
	// sdf for the extrusion region: z = [-height, height]
	b := math.Abs(p.Z) - s.height
	// return the intersection
	return math.Max(a, b)
}

// SetExtrude sets the extrusion control function.
func (s *ExtrudeSDF3) SetExtrude(extrude ExtrudeFunc) {
	s.extrude = extrude
}

// BoundingBox returns the bounding box for an extrusion.
func (s *ExtrudeSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Linear extrude an SDF2 with rounded edges.
// Note: The height of the extrusion is adjusted for the rounding.
// The underlying SDF2 shape is not modified.

// ExtrudeRoundedSDF3 extrudes an SDF2 to an SDF3 with rounded edges.
type ExtrudeRoundedSDF3 struct {
	sdf    SDF2
	height float64
	round  float64
	bb     Box3
}

// ExtrudeRounded3D extrudes an SDF2 to an SDF3 with rounded edges.
func ExtrudeRounded3D(sdf SDF2, height, round float64) (SDF3, error) {
	if round == 0 {
		// revert to non-rounded case
		return Extrude3D(sdf, height), nil
	}
	if sdf == nil {
		return nil, errors.New("sdf == nil")
	}
	if height <= 0 {
		return nil, errors.New("height <= 0")
	}
	if round < 0 {
		return nil, errors.New("round < 0")
	}
	if height < 2*round {
		return nil, errors.New("height < 2 * round")
	}
	s := ExtrudeRoundedSDF3{
		sdf:    sdf,
		height: (height / 2) - round,
		round:  round,
	}
	// work out the bounding box
	bb := sdf.BoundingBox()
	s.bb = Box3{V3{bb.Min.X, bb.Min.Y, -s.height}.SubScalar(round), V3{bb.Max.X, bb.Max.Y, s.height}.AddScalar(round)}
	return &s, nil
}

// Evaluate returns the minimum distance to a rounded extrusion.
func (s *ExtrudeRoundedSDF3) Evaluate(p V3) float64 {
	// sdf for the projected 2d surface
	a := s.sdf.Evaluate(V2{p.X, p.Y})
	b := math.Abs(p.Z) - s.height
	var d float64
	if b > 0 {
		// outside the object Z extent
		if a < 0 {
			// inside the boundary
			d = b
		} else {
			// outside the boundary
			d = math.Sqrt((a * a) + (b * b))
		}
	} else {
		// within the object Z extent
		if a < 0 {
			// inside the boundary
			d = math.Max(a, b)
		} else {
			// outside the boundary
			d = a
		}
	}
	return d - s.round
}

// BoundingBox returns the bounding box for a rounded extrusion.
func (s *ExtrudeRoundedSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Extrude/Loft (with rounded edges)
// Blend between sdf0 and sdf1 as we move from bottom to top.

// LoftSDF3 is an extrusion between two SDF2s.
type LoftSDF3 struct {
	sdf0, sdf1 SDF2
	height     float64
	round      float64
	bb         Box3
}

// Loft3D extrudes an SDF3 that transitions between two SDF2 shapes.
func Loft3D(sdf0, sdf1 SDF2, height, round float64) (SDF3, error) {
	if sdf0 == nil {
		return nil, errors.New("sdf0 == nil")
	}
	if sdf1 == nil {
		return nil, errors.New("sdf1 == nil")
	}
	if height <= 0 {
		return nil, errors.New("height <= 0")
	}
	if round < 0 {
		return nil, errors.New("round < 0")
	}
	if height < 2*round {
		return nil, errors.New("height < 2 * round")
	}
	s := LoftSDF3{
		sdf0:   sdf0,
		sdf1:   sdf1,
		height: (height / 2) - round,
		round:  round,
	}
	// work out the bounding box
	bb0 := sdf0.BoundingBox()
	bb1 := sdf1.BoundingBox()
	bb := bb0.Extend(bb1)
	s.bb = Box3{V3{bb.Min.X, bb.Min.Y, -s.height}.SubScalar(round), V3{bb.Max.X, bb.Max.Y, s.height}.AddScalar(round)}
	return &s, nil
}

// Evaluate returns the minimum distance to a loft extrusion.
func (s *LoftSDF3) Evaluate(p V3) float64 {
	// work out the mix value as a function of height
	k := Clamp((0.5*p.Z/s.height)+0.5, 0, 1)
	// mix the 2D SDFs
	a0 := s.sdf0.Evaluate(V2{p.X, p.Y})
	a1 := s.sdf1.Evaluate(V2{p.X, p.Y})
	a := Mix(a0, a1, k)

	b := math.Abs(p.Z) - s.height
	var d float64
	if b > 0 {
		// outside the object Z extent
		if a < 0 {
			// inside the boundary
			d = b
		} else {
			// outside the boundary
			d = math.Sqrt((a * a) + (b * b))
		}
	} else {
		// within the object Z extent
		if a < 0 {
			// inside the boundary
			d = math.Max(a, b)
		} else {
			// outside the boundary
			d = a
		}
	}
	return d - s.round
}

// BoundingBox returns the bounding box for a loft extrusion.
func (s *LoftSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Box (exact distance field)

// BoxSDF3 is a 3d box.
type BoxSDF3 struct {
	size  V3
	round float64
	bb    Box3
}

// Box3D return an SDF3 for a 3d box (rounded corners with round > 0).
func Box3D(size V3, round float64) (SDF3, error) {
	if size.LTEZero() {
		return nil, ErrMsg("size <= 0")
	}
	if round < 0 {
		return nil, ErrMsg("round < 0")
	}
	size = size.MulScalar(0.5)
	s := BoxSDF3{}
	s.size = size.SubScalar(round)
	s.round = round
	s.bb = Box3{size.Neg(), size}
	return &s, nil
}

// Evaluate returns the minimum distance to a 3d box.
func (s *BoxSDF3) Evaluate(p V3) float64 {
	return sdfBox3d(p, s.size) - s.round
}

// BoundingBox returns the bounding box for a 3d box.
func (s *BoxSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Sphere (exact distance field)

// SphereSDF3 is a sphere.
type SphereSDF3 struct {
	radius float64
	bb     Box3
}

// Sphere3D return an SDF3 for a sphere.
func Sphere3D(radius float64) (SDF3, error) {
	if radius <= 0 {
		return nil, ErrMsg("radius <= 0")
	}
	s := SphereSDF3{}
	s.radius = radius
	d := V3{radius, radius, radius}
	s.bb = Box3{d.Neg(), d}
	return &s, nil
}

// Evaluate returns the minimum distance to a sphere.
func (s *SphereSDF3) Evaluate(p V3) float64 {
	return p.Length() - s.radius
}

// BoundingBox returns the bounding box for a sphere.
func (s *SphereSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Cylinder (exact distance field)

// CylinderSDF3 is a cylinder.
type CylinderSDF3 struct {
	height float64
	radius float64
	round  float64
	bb     Box3
}

// Cylinder3D return an SDF3 for a cylinder (rounded edges with round > 0).
func Cylinder3D(height, radius, round float64) (SDF3, error) {
	if radius <= 0 {
		return nil, ErrMsg("radius <= 0")
	}
	if round < 0 {
		return nil, ErrMsg("round < 0")
	}
	if round > radius {
		return nil, ErrMsg("round > radius")
	}
	if height < 2.0*round {
		return nil, ErrMsg("height < 2 * round")
	}
	s := CylinderSDF3{}
	s.height = (height / 2) - round
	s.radius = radius - round
	s.round = round
	d := V3{radius, radius, height / 2}
	s.bb = Box3{d.Neg(), d}
	return &s, nil
}

// Capsule3D return an SDF3 for a capsule.
func Capsule3D(height, radius float64) (SDF3, error) {
	return Cylinder3D(height, radius, radius)
}

// Evaluate returns the minimum distance to a cylinder.
func (s *CylinderSDF3) Evaluate(p V3) float64 {
	d := sdfBox2d(V2{V2{p.X, p.Y}.Length(), p.Z}, V2{s.radius, s.height})
	return d - s.round
}

// BoundingBox returns the bounding box for a cylinder.
func (s *CylinderSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Truncated Cone (exact distance field)

// ConeSDF3 is a truncated cone.
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

// Cone3D returns the SDF3 for a trucated cone (round > 0 gives rounded edges).
func Cone3D(height, r0, r1, round float64) (SDF3, error) {
	if height <= 0 {
		return nil, ErrMsg("height <= 0")
	}
	if round < 0 {
		return nil, ErrMsg("round < 0")
	}
	if height < 2.0*round {
		return nil, ErrMsg("height < 2 * round")
	}
	s := ConeSDF3{}
	s.height = (height / 2) - round
	s.round = round
	// cone slope vector and normal
	s.u = V2{r1, height / 2}.Sub(V2{r0, -height / 2}).Normalize()
	s.n = V2{s.u.Y, -s.u.X}
	// inset the radii for the rounding
	ofs := round / s.n.X
	s.r0 = r0 - (1+s.n.Y)*ofs
	s.r1 = r1 - (1-s.n.Y)*ofs
	// cone slope length
	s.l = V2{s.r1, s.height}.Sub(V2{s.r0, -s.height}).Length()
	// work out the bounding box
	r := math.Max(s.r0+round, s.r1+round)
	s.bb = Box3{V3{-r, -r, -height / 2}, V3{r, r, height / 2}}
	return &s, nil
}

// Evaluate returns the minimum distance to a trucated cone.
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
	dSlope := v.Dot(s.n)
	// is p2 inside the cone?
	if dSlope < 0 && math.Abs(p2.Y) < s.height {
		return -math.Min(-dSlope, s.height-math.Abs(p2.Y)) - s.round
	}
	// is p2 closest to the slope line?
	t := v.Dot(s.u)
	if t >= 0 && t <= s.l {
		return dSlope - s.round
	}
	// is p2 closest to the base radius vertex?
	if t < 0 {
		return v.Length() - s.round
	}
	// p2 is closest to the top radius vertex
	return p2.Sub(V2{s.r1, s.height}).Length() - s.round
}

// BoundingBox return the bounding box for the trucated cone..
func (s *ConeSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Transform SDF3 (rotation, translation - distance preserving)

// TransformSDF3 is an SDF3 transformed with a 4x4 transformation matrix.
type TransformSDF3 struct {
	sdf     SDF3
	matrix  M44
	inverse M44
	bb      Box3
}

// Transform3D applies a transformation matrix to an SDF3.
func Transform3D(sdf SDF3, matrix M44) SDF3 {
	s := TransformSDF3{}
	s.sdf = sdf
	s.matrix = matrix
	s.inverse = matrix.Inverse()
	s.bb = matrix.MulBox(sdf.BoundingBox())
	return &s
}

// Evaluate returns the minimum distance to a transformed SDF3.
// Distance is *not* preserved with scaling.
func (s *TransformSDF3) Evaluate(p V3) float64 {
	return s.sdf.Evaluate(s.inverse.MulPosition(p))
}

// BoundingBox returns the bounding box of a transformed SDF3.
func (s *TransformSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Uniform XYZ Scaling of SDF3s (we can work out the distance)

// ScaleUniformSDF3 is an SDF3 scaled uniformly in XYZ directions.
type ScaleUniformSDF3 struct {
	sdf     SDF3
	k, invK float64
	bb      Box3
}

// ScaleUniform3D uniformly scales an SDF3 on all axes.
func ScaleUniform3D(sdf SDF3, k float64) SDF3 {
	m := Scale3d(V3{k, k, k})
	return &ScaleUniformSDF3{
		sdf:  sdf,
		k:    k,
		invK: 1.0 / k,
		bb:   m.MulBox(sdf.BoundingBox()),
	}
}

// Evaluate returns the minimum distance to a uniformly scaled SDF3.
// The distance is correct with scaling.
func (s *ScaleUniformSDF3) Evaluate(p V3) float64 {
	q := p.MulScalar(s.invK)
	return s.sdf.Evaluate(q) * s.k
}

// BoundingBox returns the bounding box of a uniformly scaled SDF3.
func (s *ScaleUniformSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// UnionSDF3 is a union of SDF3s.
type UnionSDF3 struct {
	sdf []SDF3
	min MinFunc
	bb  Box3
}

// Union3D returns the union of multiple SDF3 objects.
func Union3D(sdf ...SDF3) SDF3 {
	if len(sdf) == 0 {
		return nil
	}
	s := UnionSDF3{}
	// strip out any nils
	s.sdf = make([]SDF3, 0, len(sdf))
	for _, x := range sdf {
		if x != nil {
			s.sdf = append(s.sdf, x)
		}
	}
	if len(s.sdf) == 0 {
		return nil
	}
	if len(s.sdf) == 1 {
		// only one sdf - not really a union
		return s.sdf[0]
	}
	// work out the bounding box
	bb := s.sdf[0].BoundingBox()
	for _, x := range s.sdf {
		bb = bb.Extend(x.BoundingBox())
	}
	s.bb = bb
	s.min = math.Min
	return &s
}

// Evaluate returns the minimum distance to an SDF3 union.
func (s *UnionSDF3) Evaluate(p V3) float64 {
	var d float64
	for i, x := range s.sdf {
		if i == 0 {
			d = x.Evaluate(p)
		} else {
			d = s.min(d, x.Evaluate(p))
		}
	}
	return d
}

// SetMin sets the minimum function to control blending.
func (s *UnionSDF3) SetMin(min MinFunc) {
	s.min = min
}

// BoundingBox returns the bounding box of an SDF3 union.
func (s *UnionSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// DifferenceSDF3 is the difference of two SDF3s, s0 - s1.
type DifferenceSDF3 struct {
	s0  SDF3
	s1  SDF3
	max MaxFunc
	bb  Box3
}

// Difference3D returns the difference of two SDF3s, s0 - s1.
func Difference3D(s0, s1 SDF3) SDF3 {
	if s1 == nil {
		return s0
	}
	if s0 == nil {
		return nil
	}
	s := DifferenceSDF3{}
	s.s0 = s0
	s.s1 = s1
	s.max = math.Max
	s.bb = s0.BoundingBox()
	return &s
}

// Evaluate returns the minimum distance to the SDF3 difference.
func (s *DifferenceSDF3) Evaluate(p V3) float64 {
	return s.max(s.s0.Evaluate(p), -s.s1.Evaluate(p))
}

// SetMax sets the maximum function to control blending.
func (s *DifferenceSDF3) SetMax(max MaxFunc) {
	s.max = max
}

// BoundingBox returns the bounding box of the SDF3 difference.
func (s *DifferenceSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// ElongateSDF3 is the elongation of an SDF3.
type ElongateSDF3 struct {
	sdf    SDF3 // the sdf being elongated
	hp, hn V3   // positive/negative elongation vector
	bb     Box3 // bounding box
}

// Elongate3D returns the elongation of an SDF3.
func Elongate3D(sdf SDF3, h V3) SDF3 {
	h = h.Abs()
	s := ElongateSDF3{
		sdf: sdf,
		hp:  h.MulScalar(0.5),
		hn:  h.MulScalar(-0.5),
	}
	// bounding box
	bb := sdf.BoundingBox()
	bb0 := bb.Translate(s.hp)
	bb1 := bb.Translate(s.hn)
	s.bb = bb0.Extend(bb1)
	return &s
}

// Evaluate returns the minimum distance to a elongated SDF2.
func (s *ElongateSDF3) Evaluate(p V3) float64 {
	q := p.Sub(p.Clamp(s.hn, s.hp))
	return s.sdf.Evaluate(q)
}

// BoundingBox returns the bounding box of an elongated SDF3.
func (s *ElongateSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// IntersectionSDF3 is the intersection of two SDF3s.
type IntersectionSDF3 struct {
	s0  SDF3
	s1  SDF3
	max MaxFunc
	bb  Box3
}

// Intersect3D returns the intersection of two SDF3s.
func Intersect3D(s0, s1 SDF3) SDF3 {
	if s0 == nil || s1 == nil {
		return nil
	}
	s := IntersectionSDF3{}
	s.s0 = s0
	s.s1 = s1
	s.max = math.Max
	// TODO fix bounding box
	s.bb = s0.BoundingBox()
	return &s
}

// Evaluate returns the minimum distance to the SDF3 intersection.
func (s *IntersectionSDF3) Evaluate(p V3) float64 {
	return s.max(s.s0.Evaluate(p), s.s1.Evaluate(p))
}

// SetMax sets the maximum function to control blending.
func (s *IntersectionSDF3) SetMax(max MaxFunc) {
	s.max = max
}

// BoundingBox returns the bounding box of an SDF3 intersection.
func (s *IntersectionSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// CutSDF3 makes a planar cut through an SDF3.
type CutSDF3 struct {
	sdf SDF3
	a   V3   // point on plane
	n   V3   // normal to plane
	bb  Box3 // bounding box
}

// Cut3D cuts an SDF3 along a plane passing through a with normal n.
// The SDF3 on the same side as the normal remains.
func Cut3D(sdf SDF3, a, n V3) SDF3 {
	s := CutSDF3{}
	s.sdf = sdf
	s.a = a
	s.n = n.Normalize().Neg()
	// TODO - cut the bounding box
	s.bb = sdf.BoundingBox()
	return &s
}

// Evaluate returns the minimum distance to the cut SDF3.
func (s *CutSDF3) Evaluate(p V3) float64 {
	return math.Max(p.Sub(s.a).Dot(s.n), s.sdf.Evaluate(p))
}

// BoundingBox returns the bounding box of the cut SDF3.
func (s *CutSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// ArraySDF3 stores an XYZ array of a given SDF3
type ArraySDF3 struct {
	sdf  SDF3
	num  V3i
	step V3
	min  MinFunc
	bb   Box3
}

// Array3D returns an XYZ array of a given SDF3
func Array3D(sdf SDF3, num V3i, step V3) SDF3 {
	// check the number of steps
	if num[0] <= 0 || num[1] <= 0 || num[2] <= 0 {
		return nil
	}
	s := ArraySDF3{}
	s.sdf = sdf
	s.num = num
	s.step = step
	s.min = math.Min
	// work out the bounding box
	bb0 := sdf.BoundingBox()
	bb1 := bb0.Translate(step.Mul(num.SubScalar(1).ToV3()))
	s.bb = bb0.Extend(bb1)
	return &s
}

// SetMin sets the minimum function to control blending.
func (s *ArraySDF3) SetMin(min MinFunc) {
	s.min = min
}

// Evaluate returns the minimum distance to an XYZ SDF3 array.
func (s *ArraySDF3) Evaluate(p V3) float64 {
	d := math.MaxFloat64
	for j := 0; j < s.num[0]; j++ {
		for k := 0; k < s.num[1]; k++ {
			for l := 0; l < s.num[2]; l++ {
				x := p.Sub(V3{float64(j) * s.step.X, float64(k) * s.step.Y, float64(l) * s.step.Z})
				d = s.min(d, s.sdf.Evaluate(x))
			}
		}
	}
	return d
}

// BoundingBox returns the bounding box of an XYZ SDF3 array.
func (s *ArraySDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// RotateUnionSDF3 creates a union of SDF3s rotated about the z-axis.
type RotateUnionSDF3 struct {
	sdf  SDF3
	num  int
	step M44
	min  MinFunc
	bb   Box3
}

// RotateUnion3D creates a union of SDF3s rotated about the z-axis.
func RotateUnion3D(sdf SDF3, num int, step M44) SDF3 {
	// check the number of steps
	if num <= 0 {
		return nil
	}
	s := RotateUnionSDF3{}
	s.sdf = sdf
	s.num = num
	s.step = step.Inverse()
	s.min = math.Min
	// work out the bounding box
	v := sdf.BoundingBox().Vertices()
	bbMin := v[0]
	bbMax := v[0]
	for i := 0; i < s.num; i++ {
		bbMin = bbMin.Min(v.Min())
		bbMax = bbMax.Max(v.Max())
		v.MulVertices(step)
	}
	s.bb = Box3{bbMin, bbMax}
	return &s
}

// Evaluate returns the minimum distance to a rotate/union object.
func (s *RotateUnionSDF3) Evaluate(p V3) float64 {
	d := math.MaxFloat64
	rot := Identity3d()
	for i := 0; i < s.num; i++ {
		x := rot.MulPosition(p)
		d = s.min(d, s.sdf.Evaluate(x))
		rot = rot.Mul(s.step)
	}
	return d
}

// SetMin sets the minimum function to control blending.
func (s *RotateUnionSDF3) SetMin(min MinFunc) {
	s.min = min
}

// BoundingBox returns the bounding box of a rotate/union object.
func (s *RotateUnionSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// RotateCopySDF3 rotates and creates N copies of an SDF3 about the z-axis.
type RotateCopySDF3 struct {
	sdf   SDF3
	theta float64
	bb    Box3
}

// RotateCopy3D rotates and creates N copies of an SDF3 about the z-axis.
func RotateCopy3D(
	sdf SDF3, // SDF3 to rotate and copy
	num int, // number of copies
) SDF3 {
	// check the number of steps
	if num <= 0 {
		return nil
	}
	s := RotateCopySDF3{}
	s.sdf = sdf
	s.theta = Tau / float64(num)
	// work out the bounding box
	bb := sdf.BoundingBox()
	zmax := bb.Max.Z
	zmin := bb.Min.Z
	rmax := 0.0
	// find the bounding box vertex with the greatest distance from the z-axis
	// TODO - revisit - should go by real vertices
	for _, v := range bb.Vertices() {
		l := V2{v.X, v.Y}.Length()
		if l > rmax {
			rmax = l
		}
	}
	s.bb = Box3{V3{-rmax, -rmax, zmin}, V3{rmax, rmax, zmax}}
	return &s
}

// Evaluate returns the minimum distance to a rotate/copy SDF3.
func (s *RotateCopySDF3) Evaluate(p V3) float64 {
	// Map p to a point in the first copy sector.
	p2 := V2{p.X, p.Y}
	p2 = PolarToXY(p2.Length(), SawTooth(math.Atan2(p2.Y, p2.X), s.theta))
	return s.sdf.Evaluate(V3{p2.X, p2.Y, p.Z})
}

// BoundingBox returns the bounding box of a rotate/copy SDF3.
func (s *RotateCopySDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

/* WIP

// Connector3 defines a 3d connection point.
type Connector3 struct {
	Name     string
	Position V3
	Vector   V3
	Angle    float64
}

// ConnectedSDF3 is an SDF3 with connection points defined.
type ConnectedSDF3 struct {
	sdf        SDF3
	connectors []Connector3
}

// AddConnector adds connection points to an SDF3.
func AddConnector(sdf SDF3, connectors ...Connector3) SDF3 {
	// is the sdf already connected?
	if s, ok := sdf.(*ConnectedSDF3); ok {
		// append connection points
		s.connectors = append(s.connectors, connectors...)
		return s
	}
	// return a new connected sdf
	return &ConnectedSDF3{
		sdf:        sdf,
		connectors: connectors,
	}
}

// Evaluate returns the minimum distance to a connected SDF3.
func (s *ConnectedSDF3) Evaluate(p V3) float64 {
	return s.sdf.Evaluate(p)
}

// BoundingBox returns the bounding box of a connected SDF3.
func (s *ConnectedSDF3) BoundingBox() Box3 {
	return s.sdf.BoundingBox()
}

*/

//-----------------------------------------------------------------------------

// OffsetSDF3 offsets the distance function of an existing SDF3.
type OffsetSDF3 struct {
	sdf    SDF3    // the underlying SDF
	offset float64 // the distance the SDF is offset by
	bb     Box3    // bounding box
}

// Offset3D returns an SDF3 that offsets the distance function of another SDF3.
func Offset3D(sdf SDF3, offset float64) SDF3 {
	s := OffsetSDF3{
		sdf:    sdf,
		offset: offset,
	}
	// bounding box
	bb := sdf.BoundingBox()
	s.bb = NewBox3(bb.Center(), bb.Size().AddScalar(2*offset))
	return &s
}

// Evaluate returns the minimum distance to an offset SDF3.
func (s *OffsetSDF3) Evaluate(p V3) float64 {
	return s.sdf.Evaluate(p) - s.offset
}

// BoundingBox returns the bounding box of an offset SDF3.
func (s *OffsetSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------

// LineOf3D returns a union of 3D objects positioned along a line from p0 to p1.
func LineOf3D(s SDF3, p0, p1 V3, pattern string) SDF3 {
	var objects []SDF3
	if pattern != "" {
		x := p0
		dx := p1.Sub(p0).DivScalar(float64(len(pattern)))
		for _, c := range pattern {
			if c == 'x' {
				objects = append(objects, Transform3D(s, Translate3d(x)))
			}
			x = x.Add(dx)
		}
	}
	return Union3D(objects...)
}

//-----------------------------------------------------------------------------

// Multi3D creates a union of an SDF3 at translated positions.
func Multi3D(s SDF3, positions V3Set) SDF3 {
	if (s == nil) || (len(positions) == 0) {
		return nil
	}
	objects := make([]SDF3, len(positions))
	for i, p := range positions {
		objects[i] = Transform3D(s, Translate3d(p))
	}
	return Union3D(objects...)
}

//-----------------------------------------------------------------------------

// Orient3D creates a union of an SDF3 at oriented directions.
func Orient3D(s SDF3, base V3, directions V3Set) SDF3 {
	if (s == nil) || (len(directions) == 0) {
		return nil
	}
	objects := make([]SDF3, len(directions))
	for i, d := range directions {
		objects[i] = Transform3D(s, base.RotateToVector(d))
	}
	return Union3D(objects...)
}

//-----------------------------------------------------------------------------
