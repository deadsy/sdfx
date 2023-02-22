package render

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------

// Define the ABAQUS or CalculiX inp file sections.

type InpComments struct {
	Text [56]byte // General comments.
}

type InpHeading struct {
	Title  [8]byte  //
	Break0 [1]byte  // Line break.
	Model  [15]byte //
	Tab    [1]byte  // Tab.
	Date   [26]byte //
	Break1 [1]byte  // Line break.
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
	cmnts := InpComments{}
	copy(cmnts.Text[:], []byte("**\n** Structure: tetrahedral elements of a 3D model.\n**\n"))

	fmt.Println(len(cmnts.Text))

	err = binary.Write(buf, binary.LittleEndian, &cmnts)
	if err != nil {
		return nil, err
	}

	hdng := InpHeading{}
	copy(hdng.Title[:], []byte("*HEADING"))
	copy(hdng.Break0[:], []byte("\n"))
	copy(hdng.Model[:], []byte("Model: 3D model"))
	copy(hdng.Tab[:], []byte("\t"))
	copy(hdng.Date[:], []byte("Date: "+time.Now().UTC().Format("2006-Jan-02 MST")))
	copy(hdng.Break1[:], []byte("\n"))

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
