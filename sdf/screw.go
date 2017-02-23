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
// Thread Database - lookup standard screw threads by name

type ThreadParameters struct {
	Name   string  // name of screw thread
	Radius float64 // major radius of screw
	Pitch  float64 // thread to thread distance of screw
	Units  string  // "inch" or "mm"
}

type ThreadDatabase map[string]*ThreadParameters

var thread_db = Init_ThreadLookup()

// Unified Thread Standard
// name = thread name
// diameter = screw major diameter
// tpi = threads per inch
func (m ThreadDatabase) UTSAdd(name string, diameter, tpi float64) {
	t := ThreadParameters{}
	t.Name = name
	t.Radius = diameter / 2.0
	t.Pitch = 1.0 / tpi
	t.Units = "inch"
	m[name] = &t
}

// ISO Thread Standard
// name = thread name
// diameter = screw major diameter
// pitch = thread pitch
func (m ThreadDatabase) ISOAdd(name string, diameter, pitch float64) {
	t := ThreadParameters{}
	t.Name = name
	t.Radius = diameter / 2.0
	t.Pitch = pitch
	t.Units = "mm"
	m[name] = &t
}

func Init_ThreadLookup() ThreadDatabase {
	m := make(ThreadDatabase)
	m.UTSAdd("unc_1", 1.0, 8)
	m.UTSAdd("unf_1", 1.0, 12)
	m.UTSAdd("unc_1/4", 1.0/4.0, 20)
	m.UTSAdd("unf_1/4", 1.0/4.0, 28)
	m.UTSAdd("unc_1/2", 1.0/2.0, 13)
	m.UTSAdd("unf_1/2", 1.0/2.0, 20)

	m.ISOAdd("m6c", 6, 1)
	m.ISOAdd("m6f", 6, 0.75)

	return m
}

// lookup the parameters for a thread by name
func ThreadLookup(name string) *ThreadParameters {
	return thread_db[name]
}

// Hex Head Radius (empirical)
func (t *ThreadParameters) Hex_Radius() float64 {
	screw_d := t.Radius * 2.0
	hex_w := screw_d * 1.6
	hex_r := hex_w / (2.0 * math.Cos(DtoR(30)))
	return hex_r
}

// Hex Head Height (empirical)
func (t *ThreadParameters) Hex_Height() float64 {
	hex_r := t.Hex_Radius()
	hex_h := 2.0 * hex_r * (5.0 / 12.0)
	return hex_h
}

//-----------------------------------------------------------------------------

// Create a Hex Head Screw/Bolt
// name = thread name
// total_length = threaded length + shank length
// shank length = non threaded length
func Hex_Screw(name string, total_length, shank_length float64) SDF3 {
	t := ThreadLookup(name)
	if t == nil {
		return nil
	}
	if total_length < 0 {
		return nil
	}
	if shank_length < 0 {
		return nil
	}
	thread_length := total_length - shank_length
	if thread_length < 0 {
		thread_length = 0
	}

	// hex head
	hex_r := t.Hex_Radius()
	hex_h := t.Hex_Height()
	z_ofs := 0.5 * (total_length + shank_length + hex_h)
	round := hex_r * 0.08
	hex_2d := Polygon2D(Nagon(6, hex_r-round))
	hex_2d = Offset2D(hex_2d, round)
	hex_3d := Extrude3D(hex_2d, hex_h)
	// round off the edges
	sphere_3d := Sphere3D(hex_r * 1.55)
	sphere_3d = Transform3D(sphere_3d, Translate3d(V3{0, 0, -hex_r * 0.9}))
	hex_3d = Intersection3D(hex_3d, sphere_3d)
	// add a rounded cylinder
	hex_3d = Union3D(hex_3d, Cylinder3D(hex_h*1.05, hex_r*0.8, round))
	hex_3d = Transform3D(hex_3d, Translate3d(V3{0, 0, z_ofs}))

	// shank
	z_ofs = 0.5 * total_length
	shank_3d := Cylinder3D(shank_length, t.Radius, 0)
	shank_3d = Transform3D(shank_3d, Translate3d(V3{0, 0, z_ofs}))

	// thread
	screw_3d := Screw3D(ISOThread(t.Radius, t.Pitch), thread_length, t.Pitch, 1)

	s := Union3D(hex_3d, screw_3d)
	s = Union3D(s, shank_3d)
	return s
}

//-----------------------------------------------------------------------------
// Thread Profiles

// Return a 2d profile for an acme thread.
// radius = radius of thread
// pitch = thread to thread distance
func AcmeThread(radius, pitch float64) SDF2 {

	h := radius - 0.5*pitch
	theta := DtoR(29.0 / 2.0)
	delta := 0.25 * pitch * math.Tan(theta)
	x_ofs0 := 0.25*pitch - delta
	x_ofs1 := 0.25*pitch + delta

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
	return Polygon2D(acme)
}

// Return the 2d profile for an ISO/UTS thread.
// https://en.wikipedia.org/wiki/ISO_metric_screw_thread
// https://en.wikipedia.org/wiki/Unified_Thread_Standard
// radius = radius of thread
// pitch = thread to thread distance
func ISOThread(radius, pitch float64) SDF2 {

	theta := DtoR(30.0)
	h := pitch / (2.0 * math.Tan(theta))
	r_maj := radius
	r_root := (1.0 / 6.0) * h
	r_min := radius - (7.0/8.0)*h + r_root
	x_ofs0 := (1.0 / 16.0) * pitch
	x_ofs1 := (3.0 / 8.0) * pitch

	iso := NewSmoother(false)
	iso.Add(V2{radius, 0})
	iso.Add(V2{radius, r_min})
	iso.AddSmooth(V2{x_ofs1, r_min}, 3, r_root)
	iso.Add(V2{x_ofs0, r_maj})
	iso.Add(V2{-x_ofs0, r_maj})
	iso.AddSmooth(V2{-x_ofs1, r_min}, 3, r_root)
	iso.Add(V2{-radius, r_min})
	iso.Add(V2{-radius, 0})
	iso.Smooth()

	//RenderDXF(iso.Vertices(), "iso.dxf")
	return Polygon2D(iso.Vertices())
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
// starts = number of thread starts (< 0 for left hand threads)
func Screw3D(thread SDF2, length, pitch float64, starts int) SDF3 {
	s := ScrewSDF3{}
	s.thread = thread
	s.pitch = pitch
	s.length = length / 2
	s.lead = -pitch * float64(starts)
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
