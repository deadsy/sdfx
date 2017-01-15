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
	BoundingBox() Box
}

type SDF2 interface {
	Evaluate(p V2) float64
	Evaluate2(p V2) float64
}

//-----------------------------------------------------------------------------
// 3D Normal Box

type BoxSDF struct {
	Size V3
}

func NewBoxSDF(size V3) SDF3 {
	// note: store a modified size
	return &BoxSDF{size.MulScalar(0.5)}
}

func (s *BoxSDF) Evaluate(p V3) float64 {
	d := p.Abs().Sub(s.Size)
	return d.Max(V3{0, 0, 0}).Length() + math.Min(d.MaxComponent(), 0)
}

func (s *BoxSDF) BoundingBox() Box {
	return Box{s.Size.Negate(), s.Size}
}

//-----------------------------------------------------------------------------
// 3D Rounded Box

type RoundedBoxSDF struct {
	Size   V3
	Radius float64
}

func NewRoundedBoxSDF(size V3, radius float64) SDF3 {
	// note: store a modified size
	return &RoundedBoxSDF{size.MulScalar(0.5).SubScalar(radius), radius}
}

func (s *RoundedBoxSDF) Evaluate(p V3) float64 {
	d := p.Abs().Sub(s.Size)
	return d.Max(V3{0, 0, 0}).Length() + math.Min(d.MaxComponent(), 0) - s.Radius
}

func (s *RoundedBoxSDF) BoundingBox() Box {
	d := s.Size.AddScalar(s.Radius)
	return Box{d.Negate(), d}
}

//-----------------------------------------------------------------------------
// 2D Rectangle

type RectangleSDF struct {
	Size V2
	K    float64
}

func NewRectangleSDF(size V2) SDF2 {
	// note: store a modified size
	s := size.MulScalar(0.5)
	k := s[1] - s[0]
	return &RectangleSDF{s, k}
}

func (s *RectangleSDF) Evaluate(p V2) float64 {
	d := p.Abs().Sub(s.Size)
	return d.Max(V2{0, 0}).Length() + math.Min(d.MaxComponent(), 0)
}

func (s *RectangleSDF) Evaluate2(p V2) float64 {
	p = p.Abs()
	d := p.Sub(s.Size)
	if d[0] > 0 && d[1] > 0 {
		return d.Length()
	}
	if p[1]-p[0] > s.K {
		return d[1]
	}
	return d[0]
}

//-----------------------------------------------------------------------------
