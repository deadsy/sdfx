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
)

const (
	svgHeader = `<?xml version="1.0" encoding="utf-8"?>
<!-- Generator: http://github.com/gmlewis/sdfx -->
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
  width="%[1]vpx" height="%[2]vpx" viewBox="0 0 %[1]v %[2]v" enable-background="new 0 0 %[1]v %[2]v"
  xml:space="preserve">
`
)

//-----------------------------------------------------------------------------

// SVG represents an SVG renderer.
type SVG struct {
	name     string
	p0s, p1s []V2
	min, max V2
}

// NewSVG returns an SVG renderer.
func NewSVG(name string) *SVG {
	return &SVG{
		name: name,
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
	f, err := os.Create(s.name)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, svgHeader, s.max.X-s.min.X, s.max.Y-s.min.Y)
	for i, p0 := range s.p0s {
		p1 := s.p1s[i]
		fmt.Fprintf(f, `<line x1="%v" y1="%v" x2="%v" y2="%v"/>`+"\n", p0.X-s.min.X, s.max.Y-p0.Y, p1.X-s.min.X, s.max.Y-p1.Y)
	}
	fmt.Fprintln(f, "</svg>")
	return f.Close()
}

//-----------------------------------------------------------------------------

// SaveSVG writes line segments to an SVG file.
func SaveSVG(path string, mesh []*Line2_PP) error {
	s := NewSVG(path)
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
func WriteSVG(wg *sync.WaitGroup, path string) (chan<- *Line2_PP, error) {

	s := NewSVG(path)

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
