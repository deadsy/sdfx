//-----------------------------------------------------------------------------
/*

 */
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
)

//-----------------------------------------------------------------------------

type SDF3 interface {
	Evaluate(p V3) float64
	BoundingBox() Box3
}

type SDF2 interface {
	Evaluate(p V2) float64
	BoundingBox() Box2
}

//-----------------------------------------------------------------------------
// Basic SDF Functions

func sdf_box3d(p, s V3) float64 {
	d := p.Abs().Sub(s)
	return d.Max(V3{0, 0, 0}).Length() + math.Min(d.MaxComponent(), 0)
}

func sdf_box2d(p, s V2) float64 {
	d := p.Abs().Sub(s)
	return d.Max(V2{0, 0}).Length() + math.Min(d.MaxComponent(), 0)
}

/* alternate function - probably faster
func sdf_box2d(p, s V2) float64 {
	p = p.Abs()
	d := p.Sub(s)
	k := s[1] - s[0]
	if d[0] > 0 && d[1] > 0 {
		return d.Length()
	}
	if p[1]-p[0] > k {
		return d[1]
	}
	return d[0]
}
*/

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
// 3D Solid of Revolution

type SorSDF3 struct {
	Sdf *SDF2
}

func NewSorSDF3(sdf *SDF2) SDF3 {
	return &SorSDF3{sdf}
}

func (s *SorSDF3) Evaluate(p V3) float64 {
	return 0
}

func (s *SorSDF3) BoundingBox() Box3 {
	return Box3{}
}

//-----------------------------------------------------------------------------
// 2D Circle

type CircleSDF2 struct {
	Radius float64
}

func NewCircleSDF2(radius float64) SDF2 {
	return &CircleSDF2{radius}
}

func (s *CircleSDF2) Evaluate(p V2) float64 {
	return p.Length() - s.Radius
}

func (s *CircleSDF2) BoundingBox() Box2 {
	return Box2{}
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
	return Box2{}
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
	return sdf_box2d(s.Size, p) - s.Radius
}

func (s *RoundedBoxSDF2) BoundingBox() Box2 {
	return Box2{}
}

//-----------------------------------------------------------------------------
