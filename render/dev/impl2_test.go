package dev

import (
	"context"
	"github.com/deadsy/sdfx/sdf"
	"image"
	"sync"
	"testing"
)

func BenchmarkDevRenderer2_Render(b *testing.B) {
	s, _ := sdf.ArcSpiral2D(1.0, 20.0, 0.25*sdf.Pi, 8*sdf.Tau, 1.0)
	impl := newDevRenderer2(s)
	b.ReportAllocs()
	fullRender := image.NewRGBA(image.Rect(0, 0, 1920, 1080))
	state := RendererState{
		ResInv: 1,
		Bb:     s.BoundingBox(),
	}
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
