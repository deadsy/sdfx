package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"sync"
	"testing"
)

func BenchmarkDevRenderer3_Render(b *testing.B) {
	s, _ := sdf.ArcSpiral2D(1.0, 20.0, 0.25*sdf.Pi, 8*sdf.Tau, 1.0)
	s3, _ := sdf.ExtrudeRounded3D(s, 4, 1)
	impl := newDevRenderer3(s3)
	b.ReportAllocs()
	state := RendererState{
		ResInv: 8,
		Bb:     s.BoundingBox(),
	}
	fullRender := image.NewRGBA(image.Rect(0, 0, 1920/state.ResInv, 1080/state.ResInv))
	lock1 := &sync.RWMutex{}
	lock2 := &sync.RWMutex{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := impl.Render(context.Background(), &state, lock1, lock2, nil, fullRender)
		if err != nil {
			b.Fatal(err)
		}
	}
}
