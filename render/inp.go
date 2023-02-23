package render

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

//-----------------------------------------------------------------------------

// Define the ABAQUS or CalculiX inp file sections.

// Don't modify the text. Its size matters.
const inpComment = "**\n** Structure: finite elements of a 3D model.\n** Generated by: https://github.com/deadsy/sdfx\n**\n"

type InpComment struct {
	Text [99]byte // Exact size of text.
}

func NewInpComment() InpComment {
	cmnts := InpComment{}
	copy(cmnts.Text[:], []byte(inpComment))
	return cmnts
}

// Don't modify the text. Its size matters.
const inpHeading = "*HEADING\nModel: 3D model Date: N/A\n"

type InpHeading struct {
	Text [35]byte // Exact size of text.
}

func NewInpHeading() InpHeading {
	hdng := InpHeading{}
	copy(hdng.Text[:], []byte(inpHeading))
	return hdng
}

// Don't modify the text. Its size matters.
const inpNode = "*NODE\n"

type InpNode struct {
	Text [6]byte // Exact size of text.
}

func NewInpNode() InpNode {
	nd := InpNode{}
	copy(nd.Text[:], []byte(inpNode))
	return nd
}

//-----------------------------------------------------------------------------

// writeFE writes a stream of finite elements in the shape of tetrahedra to an ABAQUS or CalculiX file.
func writeFE(wg *sync.WaitGroup, path string) (chan<- []*Tetrahedron, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	// Use buffered IO for optimal IO writes.
	// The default buffer size doesn't appear to limit performance.
	buf := bufio.NewWriter(f)

	err = binary.Write(buf, binary.LittleEndian, NewInpComment())
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, NewInpHeading())
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, NewInpNode())
	if err != nil {
		return nil, err
	}

	// External code writes tetrahedra to this channel.
	// This goroutine reads the channel and writes tetrahedra to the file.
	c := make(chan []*Tetrahedron)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer f.Close()

		// read tetrahedra from the channel and write them to the file
		for ts := range c {
			for _, t := range ts {
				_ = t
				// TODO.
			}
		}

		// flush the tetrahedra
		buf.Flush()
	}()

	return c, nil
}
