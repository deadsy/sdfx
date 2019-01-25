//-----------------------------------------------------------------------------
/*

SVG Rendering Code

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"os"
	"sync"

	svg "github.com/ajstarks/svgo/float"
)

//-----------------------------------------------------------------------------

// SVG represents an SVG renderer.
type SVG struct {
	filename  string
	lineStyle string
	p0s, p1s  []V2
	min, max  V2
}

// NewSVG returns an SVG renderer.
func NewSVG(filename, lineStyle string) *SVG {
	return &SVG{
		filename:  filename,
		lineStyle: lineStyle,
	}
}

// Line outputs a line to the SVG file.
func (s *SVG) Line(p0, p1 V2) {
	if len(s.p0s) == 0 {
		s.min = p0.Min(p1)
		s.max = p0.Max(p1)
	} else {
		s.min = s.min.Min(p0)
		s.min = s.min.Min(p1)
		s.max = s.max.Max(p0)
		s.max = s.max.Max(p1)
	}
	s.p0s = append(s.p0s, p0)
	s.p1s = append(s.p1s, p1)
}

// Save closes the SVG file.
func (s *SVG) Save() error {
	f, err := os.Create(s.filename)
	if err != nil {
		return err
	}

	width := s.max.X - s.min.X
	height := s.max.Y - s.min.Y
	canvas := svg.New(f)
	canvas.Start(width, height)
	for i, p0 := range s.p0s {
		p1 := s.p1s[i]
		canvas.Line(p0.X-s.min.X, s.max.Y-p0.Y, p1.X-s.min.X, s.max.Y-p1.Y, s.lineStyle)
	}
	canvas.End()
	return f.Close()
}

//-----------------------------------------------------------------------------

// SaveSVG writes line segments to an SVG file.
func SaveSVG(path, lineStyle string, mesh []*Line2_PP) error {
	s := NewSVG(path, lineStyle)
	for _, v := range mesh {
		s.Line(v[0], v[1])
	}
	if err := s.Save(); err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------

// WriteSVG writes a stream of line segments to an SVG file.
func WriteSVG(wg *sync.WaitGroup, path, lineStyle string) (chan<- *Line2_PP, error) {

	s := NewSVG(path, lineStyle)

	// External code writes line segments to this channel.
	// This goroutine reads the channel and writes line segments to the file.
	c := make(chan *Line2_PP)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range c {
			s.Line(v[0], v[1])
		}
		if err := s.Save(); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	return c, nil
}

//-----------------------------------------------------------------------------
