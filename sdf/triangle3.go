//-----------------------------------------------------------------------------
/*

3D Triangles

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"sync"

	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/dhconnelly/rtreego"
)

//-----------------------------------------------------------------------------

// Triangle3 is a 3D triangle
type Triangle3 [3]v3.Vec

// Normal returns the normal vector to the plane defined by the 3D triangle.
func (t *Triangle3) Normal() v3.Vec {
	e1 := t[1].Sub(t[0])
	e2 := t[2].Sub(t[0])
	return e1.Cross(e2).Normalize()
}

// Degenerate returns true if the triangle is degenerate.
func (t *Triangle3) Degenerate(tolerance float64) bool {
	// check for identical vertices
	if t[0].Equals(t[1], tolerance) {
		return true
	}
	if t[1].Equals(t[2], tolerance) {
		return true
	}
	if t[2].Equals(t[0], tolerance) {
		return true
	}
	// TODO more tests needed
	return false
}

// BoundingBox returns a bounding box for the triangle.
func (t *Triangle3) BoundingBox() Box3 {
	return Box3{Min: t[0], Max: t[0]}.Include(t[1]).Include(t[2])
}

// Equals tests if two triangles are equal within tolerance.
func (t *Triangle3) Equals(a *Triangle3, tolerance float64) bool {
	return t[0].Equals(a[0], tolerance) &&
		t[1].Equals(a[1], tolerance) &&
		t[2].Equals(a[2], tolerance)
}

// rotateVertex rotates the vertices of a triangle
func (t *Triangle3) rotateVertex() Triangle3 {
	return Triangle3{t[2], t[0], t[1]}
}

//-----------------------------------------------------------------------------

func v3ToPoint(v v3.Vec) rtreego.Point {
	return rtreego.Point{v.X, v.Y, v.Z}
}

// Bounds returns a r-tree bounding rectangle for the triangle.
func (t *Triangle3) Bounds() *rtreego.Rect {
	b := t.BoundingBox()
	r, _ := rtreego.NewRectFromPoints(v3ToPoint(b.Min), v3ToPoint(b.Max))
	return r
}

//-----------------------------------------------------------------------------

// rotateToXY returns the transformation matrix that maps
// t[0] to the origin, t[1] to the x axis (x > 0), t[2] to the xy plane (y > 0)
func (t *Triangle3) rotateToXY() M44 {

	a := t[0]        // maps to the origin
	b := t[1].Sub(a) // maps to the x-axis
	c := t[2].Sub(a) // maps to the xy-plane

	// u maps to the x-axis
	u := b.Normalize()
	// w maps to the z-axis (normal to the triangle plane)
	w := c.Cross(u).Normalize()
	// v maps to the xy-plane (in the plane of the triangle)
	v := u.Cross(w)

	// translate to the origin
	m := Translate3d(a.Neg())

	return M44{
		u.X, u.Y, u.Z, 0,
		v.X, v.Y, v.Z, 0,
		w.X, w.Y, w.Z, 0,
		0, 0, 0, 1,
	}.Mul(m)
}

//-----------------------------------------------------------------------------

// WriteTriangles writes a stream of triangles to a slice.
func WriteTriangles(wg *sync.WaitGroup, triangles *[]Triangle3) chan<- []*Triangle3 {
	// External code writes triangles to this channel.
	// This goroutine reads the channel and appends the triangles to a slice.
	c := make(chan []*Triangle3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// read triangles from the channel and append them to the slice
		for ts := range c {
			for _, t := range ts {
				*triangles = append(*triangles, *t)
			}
		}
	}()

	return c
}

//-----------------------------------------------------------------------------
// Triangle3 Buffering

// We write triangles to a channel to decouple the rendering routines from the
// routine that writes file output. We have a lot of triangles and channels
// are not very fast, so it's best to bundle many triangles into a single channel
// write. The renderer doesn't naturally do that, so we buffer triangles before
// writing them to the channel.

// Triangle3Writer is the interface of a triangle writer/closer object.
type Triangle3Writer interface {
	Write(in []*Triangle3) error
	Close() error
}

// size the buffer to avoid re-allocations when appending.
const tBufferSize = 256
const tBufferMargin = 8 // marching cubes produces 0 to 5 triangles

// Triangle3Buffer buffers triangles before writing them to a channel.
type Triangle3Buffer struct {
	buf  []*Triangle3        // triangle buffer
	out  chan<- []*Triangle3 // output channel
	lock sync.Mutex          // lock the the buffer during access
}

// NewTriangle3Buffer returns a Triangle3Buffer.
func NewTriangle3Buffer(out chan<- []*Triangle3) Triangle3Writer {
	return &Triangle3Buffer{
		buf: make([]*Triangle3, 0, tBufferSize+tBufferMargin),
		out: out,
	}
}

func (a *Triangle3Buffer) Write(in []*Triangle3) error {
	a.lock.Lock()
	a.buf = append(a.buf, in...)
	if len(a.buf) >= tBufferSize {
		a.out <- a.buf
		a.buf = make([]*Triangle3, 0, tBufferSize+tBufferMargin)
	}
	a.lock.Unlock()
	return nil
}

// Close flushes out any remaining triangles in the buffer.
func (a *Triangle3Buffer) Close() error {
	a.lock.Lock()
	if len(a.buf) != 0 {
		a.out <- a.buf
		a.buf = nil
	}
	a.lock.Unlock()
	return nil
}

//-----------------------------------------------------------------------------
