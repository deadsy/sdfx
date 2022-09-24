//-----------------------------------------------------------------------------
/*

DXF Rendering Code

*/
//-----------------------------------------------------------------------------

package render

import (
	"errors"
	"fmt"
	"sync"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/color"
	"github.com/yofu/dxf/drawing"
	"github.com/yofu/dxf/table"
)

//-----------------------------------------------------------------------------

// DXF is a dxf drawing object.
type DXF struct {
	name    string
	drawing *drawing.Drawing
}

// NewDXF returns an empty dxf drawing object.
func NewDXF(name string) *DXF {
	d := dxf.NewDrawing()
	d.AddLayer("Lines", dxf.DefaultColor, dxf.DefaultLineType, true)
	d.AddLayer("Points", color.Red, table.LT_CONTINUOUS, true)
	return &DXF{
		name:    name,
		drawing: d,
	}
}

// Line adds a line to a dxf drawing object.
func (d *DXF) Line(p0, p1 v2.Vec) {
	d.drawing.ChangeLayer("Lines")
	d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
}

// Lines adds a set of lines to a dxf drawing object.
func (d *DXF) Lines(s v2.VecSet) {
	d.drawing.ChangeLayer("Lines")
	p1 := s[0]
	for i := 0; i < len(s)-1; i++ {
		p0 := p1
		p1 = s[i+1]
		d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
	}
}

// Points adds a set of points to a dxf drawing object.
func (d *DXF) Points(s v2.VecSet, r float64) {
	d.drawing.ChangeLayer("Points")
	for _, p := range s {
		d.drawing.Circle(p.X, p.Y, 0, r)
	}
}

// Triangle adds a triangle to a dxf drawing object.
func (d *DXF) Triangle(t Triangle2) {
	d.Lines([]v2.Vec{t[0], t[1], t[2], t[0]})
}

// Save writes a dxf drawing object to a file.
func (d *DXF) Save() error {
	err := d.drawing.SaveAs(d.name)
	if err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------

// SaveDXF writes line segments to a DXF file.
func SaveDXF(path string, mesh []*Line) error {
	d := NewDXF(path)
	d.drawing.ChangeLayer("Lines")
	for i := range mesh {
		p0 := mesh[i][0]
		p1 := mesh[i][1]
		d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
	}
	err := d.Save()
	if err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------

// WriteDXF writes a stream of line segments to a DXF file.
func WriteDXF(wg *sync.WaitGroup, path string) (chan<- *Line, error) {

	d := NewDXF(path)
	d.drawing.ChangeLayer("Lines")

	// External code writes line segments to this channel.
	// This goroutine reads the channel and writes line segments to the file.
	c := make(chan *Line)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for l := range c {
			p0 := l[0]
			p1 := l[1]
			d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
		}
		err := d.Save()
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	return c, nil
}

//-----------------------------------------------------------------------------

// RenderDXF renders an SDF2 as a DXF file. (uses quadtree sampling)
func RenderDXF(
	s sdf.SDF2, //sdf2 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {

	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := conv.V2ToV2i(bbSize.MulScalar(1 / resolution))

	fmt.Printf("rendering %s (%dx%d, resolution %.2f)\n", path, cells.X, cells.Y, resolution)

	// write the line segments to a DXF file
	var wg sync.WaitGroup
	output, err := WriteDXF(&wg, path)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	// run marching squares to generate the line segments
	marchingSquaresQuadtree(s, resolution, output)

	// stop the DXF writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

// RenderDXFSlow renders an SDF2 as a DXF file. (uses uniform grid sampling)
func RenderDXFSlow(
	s sdf.SDF2, //sdf2 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.MulScalar(1 / meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := conv.V2ToV2i(bb1Size)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox2(bb0.Center(), bb1Size)

	fmt.Printf("rendering %s (%dx%d)\n", path, cells.X, cells.Y)

	// run marching squares to generate the line segments
	m := marchingSquares(s, bb, meshInc)
	err := SaveDXF(path, m)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------

// Poly outputs a polygon as a 2D DXF file.
func Poly(p *sdf.Polygon, path string) error {

	vlist := p.Vertices()
	if vlist == nil {
		return errors.New("no vertices")
	}

	fmt.Printf("rendering %s\n", path)
	d := NewDXF(path)

	for i := 0; i < len(vlist)-1; i++ {
		p0 := vlist[i]
		p1 := vlist[i+1]
		d.Line(p0, p1)
	}

	if p.Closed() {
		p0 := vlist[len(vlist)-1]
		p1 := vlist[0]
		if !p0.Equals(p1, tolerance) {
			d.Line(p0, p1)
		}
	}

	return d.Save()
}

//-----------------------------------------------------------------------------
