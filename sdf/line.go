//-----------------------------------------------------------------------------
/*

2D lines

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
	"sync"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// Interval is a closed interval on real numbers.
type Interval [2]float64

// Sort sorts the interval endpoints lowest to highest.
func (a Interval) Sort() Interval {
	if a[0] <= a[1] {
		return a
	}
	return Interval{a[1], a[0]}
}

// Equals returns true if a == b within the tolerance limit.
func (a Interval) Equals(b Interval, tolerance float64) bool {
	return math.Abs(a[0]-b[0]) <= tolerance && math.Abs(a[1]-b[1]) <= tolerance
}

// Overlap returns true if two intervals overlap.
func (a Interval) Overlap(b Interval) bool {
	return b[0] <= a[1] && a[0] <= b[1]
}

// Intersect returns the intersection of two intervals.
func (a Interval) Intersect(b Interval) *Interval {
	if a.Overlap(b) {
		return &Interval{math.Max(a[0], b[0]), math.Min(a[1], b[1])}
	}
	return nil
}

//-----------------------------------------------------------------------------

// Line2 is a 2d line defined with end-points.
type Line2 [2]v2.Vec

// BoundingBox returns a bounding box for the line.
func (a *Line2) BoundingBox() Box2 {
	return Box2{Min: a[0], Max: a[0]}.Include(a[1])
}

// Reverse the direction of a line segment.
func (a *Line2) Reverse() *Line2 {
	return &Line2{a[1], a[0]}
}

// Equals returns true if the lines are the same (within tolerance).
func (a *Line2) Equals(b *Line2, tolerance float64) bool {
	return a[0].Equals(b[0], tolerance) && a[1].Equals(b[1], tolerance)
}

// Degenerate returns true if the line is degenerate.
func (a Line2) Degenerate(tolerance float64) bool {
	// check for identical vertices
	return a[0].Equals(a[1], tolerance)
}

// IntersectLine intersects 2 line segments.
// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
func (a *Line2) IntersectLine(b *Line2) []v2.Vec {

	p := a[0]
	r := a[1].Sub(a[0])
	q := b[0]
	s := b[1].Sub(b[0])

	k0 := r.Cross(s)        // r x s
	k1 := q.Sub(p).Cross(r) // (q - p) x r

	if k0 == 0 {
		if k1 != 0 {
			// parallel, non-intersecting
			return nil
		}

		// collinear lines
		k2 := 1.0 / r.Dot(r)
		t0 := q.Sub(p).Dot(r) * k2
		t1 := t0 + s.Dot(r)*k2

		t := Interval{t0, t1}.Sort()
		x := t.Intersect(Interval{0, 1})
		if x != nil {
			// collinear, intersecting
			p0 := p.Add(r.MulScalar(x[0]))
			if x[0] == x[1] {
				return []v2.Vec{p0}
			}
			p1 := p.Add(r.MulScalar(x[1]))
			return []v2.Vec{p0, p1}
		}

		// collinear, non-intersecting
		return nil
	}
	// non-parallel
	u := k1 / k0
	t := q.Sub(p).Cross(s) / k0
	if u >= 0 && u <= 1 && t >= 0 && t <= 1 {
		p0 := p.Add(r.MulScalar(t))
		return []v2.Vec{p0}
	}
	// non-parallel, non-intersecting
	return nil
}

//-----------------------------------------------------------------------------
// Line2 Buffering

// We write lines to a channel to decouple the rendering routines from the
// routine that writes file output. We have a lot of lines and channels
// are not very fast, so it's best to bundle many lines into a single channel
// write. The renderer doesn't naturally do that, so we buffer lines before
// writing them to the channel.

// Line2Writer is the interface of a line writer/closer object.
type Line2Writer interface {
	Write(in []*Line2) error
	Close() error
}

// size the buffer to avoid re-allocations when appending.
const lBufferSize = 128
const lBufferMargin = 4 // marching squares produces 0 to 2 lines

// Line2Buffer buffers lines before writing them to a channel.
type Line2Buffer struct {
	buf  []*Line2        // line buffer
	out  chan<- []*Line2 // output channel
	lock sync.Mutex      // lock the the buffer during access
}

// NewLine2Buffer returns a Line2Buffer.
func NewLine2Buffer(out chan<- []*Line2) Line2Writer {
	return &Line2Buffer{
		buf: make([]*Line2, 0, lBufferSize+lBufferMargin),
		out: out,
	}
}

func (a *Line2Buffer) Write(in []*Line2) error {
	a.lock.Lock()
	a.buf = append(a.buf, in...)
	if len(a.buf) >= lBufferSize {
		a.out <- a.buf
		a.buf = make([]*Line2, 0, lBufferSize+lBufferMargin)
	}
	a.lock.Unlock()
	return nil
}

// Close flushes out any remaining lines in the buffer.
func (a *Line2Buffer) Close() error {
	a.lock.Lock()
	if len(a.buf) != 0 {
		a.out <- a.buf
		a.buf = nil
	}
	a.lock.Unlock()
	return nil
}

//-----------------------------------------------------------------------------

// geometryLine is a 2d line defined as either point/point or point/vector.
type geometryLine struct {
	segment bool    // is this a line segment?
	length  float64 // segment length
	a       v2.Vec  // line start point
	b       v2.Vec  // line end point point (if segment)
	v       v2.Vec  // normalized line vector
}

// NewLinePV returns a 2d line defined by a point and vector.
func newLinePV(p, v v2.Vec) geometryLine {
	l := geometryLine{}
	l.segment = false
	l.length = 0.0
	l.a = p
	l.v = v.Normalize()
	return l
}

// NewLinePP returns a 2d line segment defined by 2 points.
func newLinePP(a, b v2.Vec) geometryLine {
	l := geometryLine{}
	v := b.Sub(a)
	l.segment = true
	l.length = v.Length()
	l.a = a
	l.b = b
	l.v = v.Normalize()
	return l
}

// Position returns the position on the line given the t value.
func (l geometryLine) Position(t float64) v2.Vec {
	return l.a.Add(l.v.MulScalar(t))
}

// Intersect returns the t parameters for the intersection between lines l and lx
func (l geometryLine) Intersect(lx geometryLine) (float64, float64, error) {
	m := M22{l.v.X, -lx.v.X, l.v.Y, -lx.v.Y}
	if m.Determinant() == 0 {
		return 0, 0, fmt.Errorf("zero/many")
	}
	p := lx.a.Sub(l.a)
	t := m.Inverse().MulPosition(p)
	return t.X, t.Y, nil
}

// Distance returns the distance to the line.
// Greater than 0 implies to the right of the line vector.
func (l geometryLine) Distance(p v2.Vec) float64 {

	n := v2.Vec{l.v.Y, -l.v.X} // normal to line
	ap := p.Sub(l.a)           // line from a to p
	dn := ap.Dot(n)            // normal distance to line

	var d float64
	if l.segment {
		// this is a line segment - consider endpoints
		t := ap.Dot(l.v) // t-parameter of projection onto line
		if t < 0 {
			d = ap.Length()
		} else if t > l.length {
			bp := p.Sub(l.b) // line from b to p
			d = bp.Length()
		} else {
			// return the normal distance
			return dn
		}
	} else {
		// not a line segment - just return the normal distance
		return dn
	}

	if dn < 0 {
		d = -d
	}
	return d
}

//-----------------------------------------------------------------------------
