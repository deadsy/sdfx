//-----------------------------------------------------------------------------
/*

SVG Rendering Code

*/
//-----------------------------------------------------------------------------

package render

import (
	"fmt"
	"os"
	"sync"

	svg "github.com/ajstarks/svgo/float"
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
)

//-----------------------------------------------------------------------------

// SVG represents an SVG renderer.
type SVG struct {
	filename  string
	lineStyle string
	p0s, p1s  []sdf.V2
	min, max  sdf.V2
}

// NewSVG returns an SVG renderer.
func NewSVG(filename, lineStyle string) *SVG {
	return &SVG{
		filename:  filename,
		lineStyle: lineStyle,
	}
}

// Line outputs a line to the SVG file.
func (s *SVG) Line(p0, p1 sdf.V2) {
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
func SaveSVG(path, lineStyle string, mesh []*Line) error {
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
func WriteSVG(wg *sync.WaitGroup, path, lineStyle string) (chan<- *Line, error) {

	s := NewSVG(path, lineStyle)

	// External code writes line segments to this channel.
	// This goroutine reads the channel and writes line segments to the file.
	c := make(chan *Line)

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

// RenderSVG renders an SDF2 as an SVG file. (uses quadtree sampling)
func RenderSVG(
	s sdf.SDF2, // sdf2 to render
	meshCells int, // number of cells on the longest axis. e.g 200
	path string, // path to filename
	lineStyle string, // SVG line style
) error {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := conv.V2ToV2i(bbSize.DivScalar(resolution))

	fmt.Printf("rendering %s (%dx%d, resolution %.2f)\n", path, cells.X, cells.Y, resolution)

	// write the line segments to an SVG file
	var wg sync.WaitGroup
	output, err := WriteSVG(&wg, path, lineStyle)
	if err != nil {
		return err
	}

	// run marching squares to generate the line segments
	marchingSquaresQuadtree(s, resolution, output)

	// stop the SVG writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
	return nil
}

// RenderSVGSlow renders an SDF2 as an SVG file. (uses uniform grid sampling)
func RenderSVGSlow(
	s sdf.SDF2, // sdf2 to render
	meshCells int, // number of cells on the longest axis. e.g 200
	path string, // path to filename
	lineStyle string, // SVG line style
) error {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := conv.V2ToV2i(bb1Size)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox2(bb0.Center(), bb1Size)

	fmt.Printf("rendering %s (%dx%d)\n", path, cells.X, cells.Y)

	// run marching squares to generate the line segments
	m := marchingSquares(s, bb, meshInc)
	return SaveSVG(path, lineStyle, m)
}

//-----------------------------------------------------------------------------
