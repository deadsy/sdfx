//-----------------------------------------------------------------------------
/*

Output a 3D triangle mesh to a 3MF file.

*/
//-----------------------------------------------------------------------------

package render

import (
	"fmt"
	"sync"

	"github.com/deadsy/sdfx/vec/conv"
	"github.com/hpinc/go3mf"
)

//-----------------------------------------------------------------------------

// Write3MF writes a stream of triangles to a 3MF file.
func Write3MF(wg *sync.WaitGroup, path string) (chan<- []*Triangle3, error) {

	f, err := go3mf.CreateWriter(path)
	if err != nil {
		return nil, err
	}

	// External code writes triangles to this channel.
	// This goroutine reads the channel and writes triangles to the file.
	c := make(chan []*Triangle3)

	var model go3mf.Model
	var mesh go3mf.Mesh

	// add the mesh to the model
	obj := &go3mf.Object{Mesh: &mesh}
	obj.ID = model.Resources.UnusedID()
	model.Resources.Objects = append(model.Resources.Objects, obj)
	model.Build.Items = append(model.Build.Items, &go3mf.Item{ObjectID: obj.ID})

	// use the mesh builder to de-dup the vertices
	mb := go3mf.NewMeshBuilder(&mesh)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer f.Close()
		// read triangles from the channel and add them to the model
		for ts := range c {
			for _, t := range ts {
				v1 := mb.AddVertex(conv.V3ToPoint3D(t.V[0]))
				v2 := mb.AddVertex(conv.V3ToPoint3D(t.V[1]))
				v3 := mb.AddVertex(conv.V3ToPoint3D(t.V[2]))
				mesh.Triangles.Triangle = append(mesh.Triangles.Triangle, go3mf.Triangle{V1: v1, V2: v2, V3: v3})
			}
		}
		// encode and write out the file
		if err := f.Encode(&model); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	return c, nil
}

//-----------------------------------------------------------------------------
