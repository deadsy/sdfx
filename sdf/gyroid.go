//-----------------------------------------------------------------------------
/*

Gyroids

https://en.wikipedia.org/wiki/Gyroid

*/
//-----------------------------------------------------------------------------

package sdf

import v3 "github.com/deadsy/sdfx/vec/v3"

//-----------------------------------------------------------------------------

// GyroidSDF3 is a 3d gyroid.
type GyroidSDF3 struct {
	k v3.Vec // scaling factor
}

// Gyroid3D returns a 3d gyroid.
func Gyroid3D(scale v3.Vec) (SDF3, error) {
	return &GyroidSDF3{
		k: v3.Vec{Tau / scale.X, Tau / scale.Y, Tau / scale.Z},
	}, nil
}

// Evaluate returns the minimum distance to a 3d gyroid.
func (s *GyroidSDF3) Evaluate(p v3.Vec) float64 {
	p = p.Mul(s.k)
	return p.Sin().Dot(v3.Vec{p.Y, p.Z, p.X}.Cos())
}

// BoundingBox returns the bounding box for a 3d gyroid.
func (s *GyroidSDF3) BoundingBox() Box3 {
	// The surface is defined for all xyz, so the bounding box is a point at the origin.
	// To use the surface it needs to be intersected an external bounding volume.
	return Box3{}
}

//-----------------------------------------------------------------------------
