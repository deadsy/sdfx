//-----------------------------------------------------------------------------
/*

STL Load/Save

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
)

//-----------------------------------------------------------------------------

// STLHeader defines the STL file header.
type STLHeader struct {
	_     [80]uint8 // Header
	Count uint32    // Number of triangles
}

// STLTriangle defines the triangle data within an STL file.
type STLTriangle struct {
	Normal, Vertex1, Vertex2, Vertex3 [3]float32
	_                                 uint16 // Attribute byte count
}

//-----------------------------------------------------------------------------

// SaveSTL writes a triangle mesh to an STL file.
func SaveSTL(path string, mesh []*Triangle3) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	header := STLHeader{}
	header.Count = uint32(len(mesh))
	if err := binary.Write(buf, binary.LittleEndian, &header); err != nil {
		return err
	}

	var d STLTriangle
	for _, triangle := range mesh {
		n := triangle.Normal()
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
		if err := binary.Write(buf, binary.LittleEndian, &d); err != nil {
			return err
		}
	}

	return buf.Flush()
}

//-----------------------------------------------------------------------------

// WriteSTL writes a stream of triangles to an STL file.
func WriteSTL(wg *sync.WaitGroup, path string) (chan<- *Triangle3, error) {

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	// Use buffered IO for optimal IO writes.
	// The default buffer size doesn't appear to limit performance.
	buf := bufio.NewWriter(f)

	// write an empty header
	hdr := STLHeader{}
	if err := binary.Write(buf, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}

	// External code writes triangles to this channel.
	// This goroutine reads the channel and writes triangles to the file.
	c := make(chan *Triangle3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer f.Close()

		var count uint32
		var d STLTriangle
		// read triangles from the channel and write them to the file
		for t := range c {
			n := t.Normal()
			d.Normal[0] = float32(n.X)
			d.Normal[1] = float32(n.Y)
			d.Normal[2] = float32(n.Z)
			d.Vertex1[0] = float32(t.V[0].X)
			d.Vertex1[1] = float32(t.V[0].Y)
			d.Vertex1[2] = float32(t.V[0].Z)
			d.Vertex2[0] = float32(t.V[1].X)
			d.Vertex2[1] = float32(t.V[1].Y)
			d.Vertex2[2] = float32(t.V[1].Z)
			d.Vertex3[0] = float32(t.V[2].X)
			d.Vertex3[1] = float32(t.V[2].Y)
			d.Vertex3[2] = float32(t.V[2].Z)
			if err := binary.Write(buf, binary.LittleEndian, &d); err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			count++
		}
		// flush the triangles
		buf.Flush()

		// back to the start of the file
		if _, err := f.Seek(0, 0); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		// rewrite the header with the correct mesh count
		hdr.Count = count
		if err := binary.Write(f, binary.LittleEndian, &hdr); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	return c, nil
}

//-----------------------------------------------------------------------------
