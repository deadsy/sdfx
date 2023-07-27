//-----------------------------------------------------------------------------
/*

STL Load/Save

*/
//-----------------------------------------------------------------------------

package render

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	v3 "github.com/deadsy/sdfx/vec/v3"
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

// parseFloats converts float value strings to []float64.
func parseFloats(in []string) ([]float64, error) {
	out := make([]float64, len(in))
	for i := range in {
		val, err := strconv.ParseFloat(in[i], 64)
		if err != nil {
			return nil, err
		}
		out[i] = val
	}
	return out, nil
}

// loadSTLAscii loads an STL file created in ASCII format.
func loadSTLAscii(file *os.File) ([]*Triangle3, error) {
	var v []v3.Vec
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 4 && fields[0] == "vertex" {
			f, err := parseFloats(fields[1:])
			if err != nil {
				return nil, err
			}
			v = append(v, v3.Vec{f[0], f[1], f[2]})
		}
	}
	// make triangles out of every 3 vertices
	var mesh []*Triangle3
	for i := 0; i < len(v); i += 3 {
		mesh = append(mesh, NewTriangle3(v[i+0], v[i+1], v[i+2]))
	}
	return mesh, scanner.Err()
}

// loadSTLBinary loads an STL file created in binary format.
func loadSTLBinary(file *os.File) ([]*Triangle3, error) {
	r := bufio.NewReader(file)
	header := STLHeader{}
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	mesh := make([]*Triangle3, int(header.Count))
	for i := range mesh {
		d := STLTriangle{}
		if err := binary.Read(r, binary.LittleEndian, &d); err != nil {
			return nil, err
		}
		v1 := v3.Vec{float64(d.Vertex1[0]), float64(d.Vertex1[1]), float64(d.Vertex1[2])}
		v2 := v3.Vec{float64(d.Vertex2[0]), float64(d.Vertex2[1]), float64(d.Vertex2[2])}
		v3 := v3.Vec{float64(d.Vertex3[0]), float64(d.Vertex3[1]), float64(d.Vertex3[2])}
		mesh[i] = NewTriangle3(v1, v2, v3)
	}
	return mesh, nil
}

// LoadSTL loads an STL file (ascii or binary) and returns the triangle mesh.
func LoadSTL(path string) ([]*Triangle3, error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// get file size
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()

	// read header, get expected binary size
	header := STLHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	expectedSize := int64(header.Count)*50 + 84

	// rewind to start of file
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// parse ascii or binary stl
	if size == expectedSize {
		return loadSTLBinary(file)
	}
	return loadSTLAscii(file)
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

// writeSTL writes a stream of triangles to an STL file.
func writeSTL(wg *sync.WaitGroup, path string) (chan<- []*Triangle3, error) {

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
	c := make(chan []*Triangle3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer f.Close()

		var count uint32
		var d STLTriangle
		// read triangles from the channel and write them to the file
		for ts := range c {
			for _, t := range ts {
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
