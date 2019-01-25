//-----------------------------------------------------------------------------
/*

DXF Rendering Code

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"sync"

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
func (d *DXF) Line(p0, p1 V2) {
	d.drawing.ChangeLayer("Lines")
	d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
}

// Lines adds a set of lines to a dxf drawing object.
func (d *DXF) Lines(s V2Set) {
	d.drawing.ChangeLayer("Lines")
	p1 := s[0]
	for i := 0; i < len(s)-1; i++ {
		p0 := p1
		p1 = s[i+1]
		d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
	}
}

// Points adds a set of points to a dxf drawing object.
func (d *DXF) Points(s V2Set, r float64) {
	d.drawing.ChangeLayer("Points")
	for _, p := range s {
		d.drawing.Circle(p.X, p.Y, 0, r)
	}
}

// Triangle adds a triangle to a dxf drawing object.
func (d *DXF) Triangle(t Triangle2) {
	d.Lines([]V2{t[0], t[1], t[2], t[0]})
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
func SaveDXF(path string, mesh []*Line2_PP) error {
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
func WriteDXF(wg *sync.WaitGroup, path string) (chan<- *Line2_PP, error) {

	d := NewDXF(path)
	d.drawing.ChangeLayer("Lines")

	// External code writes line segments to this channel.
	// This goroutine reads the channel and writes line segments to the file.
	c := make(chan *Line2_PP)

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
