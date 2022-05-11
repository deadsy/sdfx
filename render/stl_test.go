package render_test

import (
	"sync"
	"testing"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

const (
	tol     = .1
	quality = 200
)

var (
	renderer = &render.MarchingCubesOctree{}
	object   sdf.SDF3
)

func init() {
	boltParms := obj.BoltParms{
		Thread:      "M16x2",
		Style:       "hex",
		Tolerance:   tol,
		TotalLength: 50.0,
		ShankLength: 10.0,
	}
	bolt, err := obj.Bolt(&boltParms)
	if err != nil {
		panic(err)
	}
	object = bolt
}

func BenchmarkSaveSTL(b *testing.B) {
	const path = "bolt_save.stl"
	// defer os.Remove(path)
	for i := 0; i < b.N; i++ {
		output := renderer.RenderSlice(object, quality)
		err := render.SaveSTL(path, output)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStreamSTL(b *testing.B) {
	const path = "bolt_stream.stl"
	// defer os.Remove(path)
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		ch, err := render.StreamSTL(&wg, path)
		if err != nil {
			b.Fatal(err)
		}
		renderer.Render(object, quality, ch)
	}
}
