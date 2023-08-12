package sdf

import (
	"sync"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Fe is a finite element.
type Fe struct {
	// Coordinates of nodes or vertices.
	V []v3.Vec
	// Coordinates of the voxel to which the element belongs.
	X int
	Y int
	Z int
}

//-----------------------------------------------------------------------------

// WriteFes writes a stream of finite elements to a slice.
func WriteFes(wg *sync.WaitGroup, elements *[]Fe) chan<- []*Fe {
	// External code writes to this channel.
	// This goroutine reads the channel and stores finite elements.
	c := make(chan []*Fe)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// read finite elements from the channel and handle them
		for fes := range c {
			for _, fe := range fes {
				*elements = append(*elements, *fe)
			}
		}
	}()

	return c
}

//-----------------------------------------------------------------------------
// Finite element Buffering

// We write finite elements to a channel to decouple the rendering routines from the
// routine that writes file output. We have a lot of finite elements and channels
// are not very fast, so it's best to bundle many finite elements into a single channel
// write. The renderer doesn't naturally do that, so we buffer finite elements before
// writing them to the channel.

// FeWriter is the interface of a finite element writer/closer object.
type FeWriter interface {
	Write(in []*Fe) error
	Close() error
}

// size the buffer to avoid re-allocations when appending.
const feBufferSize = 256

// marching cubes produces 0 or 1 finite element type of hex.
// marching cubes produces 0 to less than 20 finite element of type tet.
// TODO: can value be further calibrated?
const feBufferMargin = 20

// FeBuffer buffers finite elements before writing them to a channel.
type FeBuffer struct {
	buf  []*Fe        // finite element buffer
	out  chan<- []*Fe // output channel
	lock sync.Mutex   // lock the the buffer during access
}

// NewFeBuffer returns a FeBuffer.
func NewFeBuffer(out chan<- []*Fe) FeWriter {
	return &FeBuffer{
		buf: make([]*Fe, 0, feBufferSize+feBufferMargin),
		out: out,
	}
}

func (a *FeBuffer) Write(in []*Fe) error {
	a.lock.Lock()
	a.buf = append(a.buf, in...)
	if len(a.buf) >= tBufferSize {
		a.out <- a.buf
		a.buf = make([]*Fe, 0, feBufferSize+feBufferMargin)
	}
	a.lock.Unlock()
	return nil
}

// Close flushes out any remaining finite elements in the buffer.
func (a *FeBuffer) Close() error {
	a.lock.Lock()
	if len(a.buf) != 0 {
		a.out <- a.buf
		a.buf = nil
	}
	a.lock.Unlock()
	return nil
}

//-----------------------------------------------------------------------------
