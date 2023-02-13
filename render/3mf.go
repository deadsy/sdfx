//-----------------------------------------------------------------------------
/*

Output a 3D triangle mesh to a 3MF file.

https://3mf.io/specification/

Notes:

3D manufacturing files (3mf) generally contain meta data about the 3d object.
The files produced by this code are very basic. They are the equivalent of an
STL in 3MF format. That is: just the triangle mesh with a default 1mm unit.

File sizes for 3MF are around 7x smaller than an STL with the same mesh.

3MF files are not identical from run to run. 3MF files are a zipped archive.
The contents of the archive *are* the same but the containing zip file differs.

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
				v1 := mb.AddVertex(V3ToPoint3D(t.V[0]))
				v2 := mb.AddVertex(V3ToPoint3D(t.V[1]))
				v3 := mb.AddVertex(V3ToPoint3D(t.V[2]))
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

// V3ToPoint3D converts a 3D float vector to a go3mf 3D vector.
func V3ToPoint3D(a v3.Vec) go3mf.Point3D {
	return go3mf.Point3D{float32(a.X), float32(a.Y), float32(a.Z)}
}

//-----------------------------------------------------------------------------
