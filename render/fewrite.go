package render

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------

// Define the ABAQUS or CalculiX inp file sections.

type InpComments struct {
	Text string // General comments.
}

type InpHeading struct {
	Title  string //
	Break0 string //
	Model  string //
	Tab    string //
	Date   string //
	Break1 string //
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

	// write general comments
	cmnts := InpComments{
		Text: `**
		** Structure: tetrahedral elements of a 3D model.
		**\n`,
	}
	err = binary.Write(buf, binary.LittleEndian, &cmnts)
	if err != nil {
		return nil, err
	}

	hdng := InpHeading{
		Title:  "*HEADING",
		Break0: "\n",
		Model:  "Model: 3D model",
		Tab:    "\t",
		Date:   "Date: " + time.Now().UTC().Format("2006-January-02 MST"),
		Break1: "\n",
	}
	err = binary.Write(buf, binary.LittleEndian, &hdng)
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
