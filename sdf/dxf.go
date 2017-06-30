//-----------------------------------------------------------------------------
/*

DXF Rendering Code

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/drawing"
)

//-----------------------------------------------------------------------------

type DXF struct {
	name    string
	drawing *drawing.Drawing
}

func NewDXF(name string) *DXF {
	return &DXF{
		name:    name,
		drawing: dxf.NewDrawing(),
	}
}

func (d *DXF) Line(p0, p1 V2) {
	d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
}

func (d *DXF) Lines(s V2Set) {
	p1 := s[0]
	for i := 0; i < len(s)-1; i++ {
		p0 := p1
		p1 = s[i+1]
		d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
	}
}

func (d *DXF) Points(s V2Set) {
	for _, p := range s {
		d.drawing.Point(p.X, p.Y, 0)
	}
}

func (d *DXF) Triangle(t Triangle2) {
	d.Lines([]V2{t[0], t[1], t[2], t[0]})
}

func (d *DXF) Save() error {
	err := d.drawing.SaveAs(d.name)
	if err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------
