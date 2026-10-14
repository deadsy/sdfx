//-----------------------------------------------------------------------------
/*

Gyroids

https://en.wikipedia.org/wiki/Gyroid

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// GyroidSDF3 is a 3d gyroid.
type GyroidSDF3 struct {
	k          v3.Vec  // scaling factor
	invStretch float64 // 1/Lipschitz; makes g a Lipschitz-1 scalar field
}

// Gyroid3D returns a 3d gyroid.
func Gyroid3D(scale v3.Vec) (SDF3, error) {
	k := v3.Vec{Tau / scale.X, Tau / scale.Y, Tau / scale.Z}
	// Gyroid is not a true SDF. The implicit function
	//   g(p) = sin(kx·x)·cos(ky·y) + sin(ky·y)·cos(kz·z) + sin(kz·z)·cos(kx·x)
	// has gradient bounded in magnitude by √3·max(k_i) (tight for isotropic
	// scales). Dividing by that constant gives a Lipschitz-1 scalar field
	// with the same zero level set — safe for the octree's isEmpty skip.
	kMax := math.Max(math.Max(math.Abs(k.X), math.Abs(k.Y)), math.Abs(k.Z))
	return &GyroidSDF3{
		k:          k,
		invStretch: 1 / (math.Sqrt(3) * kMax),
	}, nil
}

// Evaluate returns the gyroid implicit function rescaled to Lipschitz-1.
func (s *GyroidSDF3) Evaluate(p v3.Vec) float64 {
	p = p.Mul(s.k)
	return p.Sin().Dot(v3.Vec{p.Y, p.Z, p.X}.Cos()) * s.invStretch
}

// BoundingBox returns the bounding box for a 3d gyroid.
func (s *GyroidSDF3) BoundingBox() Box3 {
	// The surface is defined for all xyz, so the bounding box is a point at the origin.
	// To use the surface it needs to be intersected an external bounding volume.
	return Box3{}
}

//-----------------------------------------------------------------------------
