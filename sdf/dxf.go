//-----------------------------------------------------------------------------
/*

DXF Rendering Code

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/color"
	"github.com/yofu/dxf/drawing"
	"github.com/yofu/dxf/table"
)

//-----------------------------------------------------------------------------

type DXF struct {
	name    string
	drawing *drawing.Drawing
}

func NewDXF(name string) *DXF {
	d := dxf.NewDrawing()
	d.AddLayer("Lines", dxf.DefaultColor, dxf.DefaultLineType, true)
	d.AddLayer("Points", color.Red, table.LT_CONTINUOUS, true)
	return &DXF{
		name:    name,
		drawing: d,
	}
}

func (d *DXF) Line(p0, p1 V2) {
	d.drawing.ChangeLayer("Lines")
	d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
}

func (d *DXF) Lines(s V2Set) {
	d.drawing.ChangeLayer("Lines")
	p1 := s[0]
	for i := 0; i < len(s)-1; i++ {
		p0 := p1
		p1 = s[i+1]
		d.drawing.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
	}
}

func (d *DXF) Points(s V2Set, r float64) {
	d.drawing.ChangeLayer("Points")
	for _, p := range s {
		d.drawing.Circle(p.X, p.Y, 0, r)
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
