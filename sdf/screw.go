//-----------------------------------------------------------------------------
/*

Screws

Screws are made by taking a 2D thread profile, rotating it about the z-axis and
spiralling it upwards as we move along z.

The 2D thread profiles are a polygon of a single thread centered on the y-axis with
the x-axis as the screw axis. Most thread profiles are symmetric about the y-axis
but a few aren't (E.g. buttress threads) so in general we build the profile of
an entire pitch period.

This code doesn't deal with thread tolerancing. If you want threads to fit properly
the radius of the thread will need to be tweaked (+/-) to give internal/external thread
clearance.

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------
// Thread Database - lookup standard screw threads by name

// ThreadParameters stores the values that define a thread.
type ThreadParameters struct {
	Name         string  // name of screw thread
	Radius       float64 // nominal major radius of screw
	Pitch        float64 // thread to thread distance of screw
	HexFlat2Flat float64 // hex head flat to flat distance
	Units        string  // "inch" or "mm"
}

type threadDatabase map[string]*ThreadParameters

var threadDB = initThreadLookup()

// UTSAdd adds a Unified Thread Standard to the thread database.
func (m threadDatabase) UTSAdd(
	name string, // thread name
	diameter float64, // screw major diameter
	tpi float64, // threads per inch
	ftof float64, // hex head flat to flat distance
) {
	t := ThreadParameters{}
	t.Name = name
	t.Radius = diameter / 2.0
	t.Pitch = 1.0 / tpi
	t.HexFlat2Flat = ftof
	t.Units = "inch"
	m[name] = &t
}

// ISOAdd adds an ISO Thread Standard to the thread database.
func (m threadDatabase) ISOAdd(
	name string, // thread name
	diameter float64, // screw major diamater
	pitch float64, // thread pitch
	ftof float64, // hex head flat to flat distance
) {
	t := ThreadParameters{}
	t.Name = name
	t.Radius = diameter / 2.0
	t.Pitch = pitch
	t.HexFlat2Flat = ftof
	t.Units = "mm"
	m[name] = &t
}

// initThreadLookup adds a collection of standard threads to the thread database.
func initThreadLookup() threadDatabase {
	m := make(threadDatabase)
	// UTS Coarse
	m.UTSAdd("unc_1/4", 1.0/4.0, 20, 7.0/16.0)
	m.UTSAdd("unc_5/16", 5.0/16.0, 18, 1.0/2.0)
	m.UTSAdd("unc_3/8", 3.0/8.0, 16, 9.0/16.0)
	m.UTSAdd("unc_7/16", 7.0/16.0, 14, 5.0/8.0)
	m.UTSAdd("unc_1/2", 1.0/2.0, 13, 3.0/4.0)
	m.UTSAdd("unc_9/16", 9.0/16.0, 12, 13.0/16.0)
	m.UTSAdd("unc_5/8", 5.0/8.0, 11, 15.0/16.0)
	m.UTSAdd("unc_3/4", 3.0/4.0, 10, 9.0/8.0)
	m.UTSAdd("unc_7/8", 7.0/8.0, 9, 21.0/16.0)
	m.UTSAdd("unc_1", 1.0, 8, 3.0/2.0)
	// UTS Fine
	m.UTSAdd("unf_1/4", 1.0/4.0, 28, 7.0/16.0)
	m.UTSAdd("unf_5/16", 5.0/16.0, 24, 1.0/2.0)
	m.UTSAdd("unf_3/8", 3.0/8.0, 24, 9.0/16.0)
	m.UTSAdd("unf_7/16", 7.0/16.0, 20, 5.0/8.0)
	m.UTSAdd("unf_1/2", 1.0/2.0, 20, 3.0/4.0)
	m.UTSAdd("unf_9/16", 9.0/16.0, 18, 13.0/16.0)
	m.UTSAdd("unf_5/8", 5.0/8.0, 18, 15.0/16.0)
	m.UTSAdd("unf_3/4", 3.0/4.0, 16, 9.0/8.0)
	m.UTSAdd("unf_7/8", 7.0/8.0, 14, 21.0/16.0)
	m.UTSAdd("unf_1", 1.0, 12, 3.0/2.0)
	// ISO Coarse
	m.ISOAdd("M1x0.25", 1, 0.25, -1)
	m.ISOAdd("M1.2x0.25", 1.2, 0.25, -1)
	m.ISOAdd("M1.6x0.35", 1.6, 0.35, 3.2)
	m.ISOAdd("M2x0.4", 2, 0.4, 4)
	m.ISOAdd("M2.5x0.45", 2.5, 0.45, 5)
	m.ISOAdd("M3x0.5", 3, 0.5, 6)
	m.ISOAdd("M4x0.7", 4, 0.7, 7)
	m.ISOAdd("M5x0.8", 5, 0.8, 8)
	m.ISOAdd("M6x1", 6, 1, 10)
	m.ISOAdd("M8x1.25", 8, 1.25, 13)
	m.ISOAdd("M10x1.5", 10, 1.5, 17)
	m.ISOAdd("M12x1.75", 12, 1.75, 19)
	m.ISOAdd("M16x2", 16, 2, 24)
	m.ISOAdd("M20x2.5", 20, 2.5, 30)
	m.ISOAdd("M24x3", 24, 3, 36)
	m.ISOAdd("M30x3.5", 30, 3.5, 46)
	m.ISOAdd("M36x4", 36, 4, 55)
	m.ISOAdd("M42x4.5", 42, 4.5, 65)
	m.ISOAdd("M48x5", 48, 5, 75)
	m.ISOAdd("M56x5.5", 56, 5.5, 85)
	m.ISOAdd("M64x6", 64, 6, 95)
	// ISO Fine
	m.ISOAdd("M1x0.2", 1, 0.2, -1)
	m.ISOAdd("M1.2x0.2", 1.2, 0.2, -1)
	m.ISOAdd("M1.6x0.2", 1.6, 0.2, 3.2)
	m.ISOAdd("M2x0.25", 2, 0.25, 4)
	m.ISOAdd("M2.5x0.35", 2.5, 0.35, 5)
	m.ISOAdd("M3x0.35", 3, 0.35, 6)
	m.ISOAdd("M4x0.5", 4, 0.5, 7)
	m.ISOAdd("M5x0.5", 5, 0.5, 8)
	m.ISOAdd("M6x0.75", 6, 0.75, 10)
	m.ISOAdd("M8x1", 8, 1, 13)
	m.ISOAdd("M10x1.25", 10, 1.25, 17)
	m.ISOAdd("M12x1.5", 12, 1.5, 19)
	m.ISOAdd("M16x1.5", 16, 1.5, 24)
	m.ISOAdd("M20x2", 20, 2, 30)
	m.ISOAdd("M24x2", 24, 2, 36)
	m.ISOAdd("M30x2", 30, 2, 46)
	m.ISOAdd("M36x3", 36, 3, 55)
	m.ISOAdd("M42x3", 42, 3, 65)
	m.ISOAdd("M48x3", 48, 3, 75)
	m.ISOAdd("M56x4", 56, 4, 85)
	m.ISOAdd("M64x4", 64, 4, 95)
	return m
}

// ThreadLookup lookups the parameters for a thread by name.
func ThreadLookup(name string) *ThreadParameters {
	t, ok := threadDB[name]
	if !ok {
		panic("thread name not found")
	}
	return t
}

// HexRadius returns the hex head radius.
func (t *ThreadParameters) HexRadius() float64 {
	if t.HexFlat2Flat < 0 {
		panic("no hex head flat to flat distance defined for this thread")
	}
	return t.HexFlat2Flat / (2.0 * math.Cos(DtoR(30)))
}

// HexHeight returns the hex head height (empirical).
func (t *ThreadParameters) HexHeight() float64 {
	return 2.0 * t.HexRadius() * (5.0 / 12.0)
}

//-----------------------------------------------------------------------------
// Thread Profiles

// AcmeThread returns the 2d profile for an acme thread.
func AcmeThread(
	radius float64, // radius of thread
	pitch float64, // thread to thread distance
) SDF2 {

	h := radius - 0.5*pitch
	theta := DtoR(29.0 / 2.0)
	delta := 0.25 * pitch * math.Tan(theta)
	xOfs0 := 0.25*pitch - delta
	xOfs1 := 0.25*pitch + delta

	acme := NewPolygon()
	acme.Add(radius, 0)
	acme.Add(radius, h)
	acme.Add(xOfs1, h)
	acme.Add(xOfs0, radius)
	acme.Add(-xOfs0, radius)
	acme.Add(-xOfs1, h)
	acme.Add(-radius, h)
	acme.Add(-radius, 0)

	//acme.Render("acme.dxf")
	return Polygon2D(acme.Vertices())
}

// ISOThread returns the 2d profile for an ISO/UTS thread.
// https://en.wikipedia.org/wiki/ISO_metric_screw_thread
// https://en.wikipedia.org/wiki/Unified_Thread_Standard
func ISOThread(
	radius float64, // radius of thread
	pitch float64, // thread to thread distance
	mode string, // internal/external thread
) SDF2 {

	theta := DtoR(30.0)
	h := pitch / (2.0 * math.Tan(theta))
	rMajor := radius
	r0 := rMajor - (7.0/8.0)*h

	iso := NewPolygon()
	if mode == "external" {
		rRoot := (pitch / 8.0) / math.Cos(theta)
		xOfs := (1.0 / 16.0) * pitch
		iso.Add(pitch, 0)
		iso.Add(pitch, r0+h)
		iso.Add(pitch/2.0, r0).Smooth(rRoot, 5)
		iso.Add(xOfs, rMajor)
		iso.Add(-xOfs, rMajor)
		iso.Add(-pitch/2.0, r0).Smooth(rRoot, 5)
		iso.Add(-pitch, r0+h)
		iso.Add(-pitch, 0)
	} else if mode == "internal" {
		rMinor := r0 + (1.0/4.0)*h
		rCrest := (pitch / 16.0) / math.Cos(theta)
		xOfs := (1.0 / 8.0) * pitch
		iso.Add(pitch, 0)
		iso.Add(pitch, rMinor)
		iso.Add(pitch/2-xOfs, rMinor)
		iso.Add(0, r0+h).Smooth(rCrest, 5)
		iso.Add(-pitch/2+xOfs, rMinor)
		iso.Add(-pitch, rMinor)
		iso.Add(-pitch, 0)
	} else {
		panic("bad mode")
	}
	//iso.Render("iso.dxf")
	return Polygon2D(iso.Vertices())
}

// ANSIButtressThread returns the 2d profile for an ANSI 45/7 buttress thread.
// https://en.wikipedia.org/wiki/Buttress_thread
// AMSE B1.9-1973
func ANSIButtressThread(
	radius float64, // radius of thread
	pitch float64, // thread to thread distance
) SDF2 {
	t0 := math.Tan(DtoR(45.0))
	t1 := math.Tan(DtoR(7.0))
	b := 0.6 // thread engagement

	h0 := pitch / (t0 + t1)
	h1 := ((b / 2.0) * pitch) + (0.5 * h0)
	hp := pitch / 2.0

	tp := NewPolygon()
	tp.Add(pitch, 0)
	tp.Add(pitch, radius)
	tp.Add(hp-((h0-h1)*t1), radius)
	tp.Add(t0*h0-hp, radius-h1).Smooth(0.0714*pitch, 5)
	tp.Add((h0-h1)*t0-hp, radius)
	tp.Add(-pitch, radius)
	tp.Add(-pitch, 0)

	//tp.Render("buttress.dxf")
	return Polygon2D(tp.Vertices())
}

// PlasticButtressThread returns the 2d profile for a screw top style plastic buttress thread.
// Similar to ANSI 45/7 - but with more corner rounding
func PlasticButtressThread(
	radius float64, // radius of thread
	pitch float64, // thread to thread distance
) SDF2 {
	t0 := math.Tan(DtoR(45.0))
	t1 := math.Tan(DtoR(7.0))
	b := 0.6 // thread engagement

	h0 := pitch / (t0 + t1)
	h1 := ((b / 2.0) * pitch) + (0.5 * h0)
	hp := pitch / 2.0

	tp := NewPolygon()
	tp.Add(pitch, 0)
	tp.Add(pitch, radius)
	tp.Add(hp-((h0-h1)*t1), radius).Smooth(0.05*pitch, 5)
	tp.Add(t0*h0-hp, radius-h1).Smooth(0.15*pitch, 5)
	tp.Add((h0-h1)*t0-hp, radius).Smooth(0.15*pitch, 5)
	tp.Add(-pitch, radius)
	tp.Add(-pitch, 0)

	//tp.Render("buttress.dxf")
	return Polygon2D(tp.Vertices())
}

//-----------------------------------------------------------------------------

// ScrewSDF3 is a 3d screw form.
type ScrewSDF3 struct {
	thread SDF2    // 2D thread profile
	pitch  float64 // thread to thread distance
	lead   float64 // distance per turn (starts * pitch)
	length float64 // total length of screw
	starts int     // number of thread starts
	bb     Box3    // bounding box
}

// Screw3D returns a screw SDF3.
func Screw3D(
	thread SDF2, // 2D thread profile
	length float64, // length of screw
	pitch float64, // thread to thread distance
	starts int, // number of thread starts (< 0 for left hand threads)
) SDF3 {
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

// Evaluate returns the minimum distance to a 3d screw form.
func (s *ScrewSDF3) Evaluate(p V3) float64 {
	// map the 3d point back to the xy space of the profile
	p0 := V2{}
	// the distance from the 3d z-axis maps to the 2d y-axis
	p0.Y = math.Sqrt(p.X*p.X + p.Y*p.Y)
	// the x/y angle and the z-height map to the 2d x-axis
	// ie: the position along thread pitch
	theta := math.Atan2(p.Y, p.X)
	z := p.Z + s.lead*theta/Tau
	p0.X = SawTooth(z, s.pitch)
	// get the thread profile distance
	d0 := s.thread.Evaluate(p0)
	// create a region for the screw length
	d1 := Abs(p.Z) - s.length
	// return the intersection
	return Max(d0, d1)
}

// BoundingBox returns the bounding box for a 3d screw form.
func (s *ScrewSDF3) BoundingBox() Box3 {
	return s.bb
}

//-----------------------------------------------------------------------------
