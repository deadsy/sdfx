//-----------------------------------------------------------------------------
/*

STL Load/Save

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"encoding/binary"
	"os"
)

//-----------------------------------------------------------------------------

type STLHeader struct {
	_     [80]uint8 // Header
	Count uint32    // Number of triangles
}

type STLTriangle struct {
	Normal, Vertex1, Vertex2, Vertex3 [3]float32
	_                                 uint16 // Attribute byte count
}

//-----------------------------------------------------------------------------

func SaveSTL(path string, mesh *Mesh) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	header := STLHeader{}
	header.Count = uint32(len(mesh.Triangles))
	if err := binary.Write(file, binary.LittleEndian, &header); err != nil {
		return err
	}

	for _, triangle := range mesh.Triangles {
		n := triangle.Normal()
		d := STLTriangle{}
		d.Normal[0] = float32(n.X)
		d.Normal[1] = float32(n.Y)
		d.Normal[2] = float32(n.Z)
		d.Vertex1[0] = float32(triangle.V[0].X)
		d.Vertex1[1] = float32(triangle.V[0].Y)
		d.Vertex1[2] = float32(triangle.V[0].Z)
		d.Vertex2[0] = float32(triangle.V[1].X)
		d.Vertex2[1] = float32(triangle.V[1].Y)
		d.Vertex2[2] = float32(triangle.V[1].Z)
		d.Vertex3[0] = float32(triangle.V[2].X)
		d.Vertex3[1] = float32(triangle.V[2].Y)
		d.Vertex3[2] = float32(triangle.V[2].Z)
		if err := binary.Write(file, binary.LittleEndian, &d); err != nil {
			return err
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
