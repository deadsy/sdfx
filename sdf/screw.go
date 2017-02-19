//-----------------------------------------------------------------------------
/*

Screws

Screws are made by taking a 2d thread profile, rotating it about the z-axis and
spiralling it upwards as we move along z.

The 2d thread profiles are a polygon of a single thread centered on the y-axis with
the x-axis as the screw axis. Most thread profiles are symmetric about the y-axis
but a few aren't (E.g. buttress threads) so in general we build the profile of
an entire pitch period.

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------
// Thread Profiles

// Return a 2d profile for an acme thread.
// radius = radius of thread
// pitch = thread to thread distance
func AcmeThread(radius, pitch float64) SDF2 {

	h := radius - 0.5*pitch
	theta := DtoR(29.0 / 2.0)
	delta := 0.25 * pitch * math.Tan(theta)
	x_ofs1 := 0.25*pitch + delta
	x_ofs0 := 0.25*pitch - delta

	acme := V2Set{
		V2{radius, 0},
		V2{radius, h},
		V2{x_ofs1, h},
		V2{x_ofs0, radius},
		V2{-x_ofs0, radius},
		V2{-x_ofs1, h},
		V2{-radius, h},
		V2{-radius, 0},
	}
	//RenderDXF(acme, "acme.dxf")
	return NewPolySDF2(acme)
}

//-----------------------------------------------------------------------------

type ScrewSDF3 struct {
	thread SDF2    // 2D thread profile
	pitch  float64 // thread to thread distance
	lead   float64 // distance per turn (starts * pitch)
	length float64 // total length of screw
	starts int     // number of thread starts
	bb     Box3    // bounding box
}

// Return a screw SDF3
// thread = 2D thread profile
// length = length of screw
// pitch = thread to thread distance
// starts = number of thread starts
func NewScrewSDF3(thread SDF2, length, pitch float64, starts int) SDF3 {
	s := ScrewSDF3{}
	s.thread = thread
	s.pitch = pitch
	s.length = length / 2
	s.lead = float64(starts) * pitch
	// Work out the bounding box.
	// The max-y axis of the sdf2 bounding box is the radius of the thread.
	bb := s.thread.BoundingBox()
	r := bb.Max.Y
	s.bb = Box3{V3{-r, -r, -s.length}, V3{r, r, s.length}}
	return &s
}

func (s *ScrewSDF3) Evaluate(p V3) float64 {
	// map the 3d point back to the xy space of the profile
	p0 := V2{}
	// the distance from the 3d z-axis maps to the 2d y-axis
	p0.Y = math.Sqrt(p.X*p.X + p.Y*p.Y)
	// the x/y angle and the z-height map to the 2d x-axis
	// ie: the position along thread pitch
	theta := math.Atan2(p.Y, p.X)
	z := p.Z + s.lead*theta/TAU
	p0.X = SawTooth(z, s.pitch)
	// get the thread profile distance
	d0 := s.thread.Evaluate(p0)
	// create a region for the screw length
	d1 := Abs(p.Z) - s.length
	// return the intersection
	return Max(d0, d1)
}

func (s *ScrewSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
