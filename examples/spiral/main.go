//-----------------------------------------------------------------------------
/*

Spirals

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/render/dev"
	"github.com/hajimehoshi/ebiten"
	"log"
	"os"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {
	s, err := sdf.ArcSpiral2D(1.0, 20.0, 0.25*sdf.Pi, 8*sdf.Tau, 1.0)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	c, _ := sdf.Circle2D(22.)
	s = sdf.Union2D(s, c)

	c2, _ := sdf.Circle2D(20.)
	c2 = sdf.Transform2D(c2, sdf.Translate2d(sdf.V2{X: 0}))
	s = sdf.Difference2D(s, c2)

	if os.Getenv("SDFX_TEST_DEV_RENDERER_2D") != "" {
		ebiten.SetWindowTitle("SDFX spiral 2D demo")
		ebiten.SetRunnableOnUnfocused(true)
		ebiten.SetWindowResizable(true)
		//ebiten.SetWindowSize(1920, 1040)
		err = dev.NewDevRenderer(s).Run()
		if err != nil {
			panic(err)
		}
	} else {
		render.RenderDXF(s, 400, "spiral.dxf")
	}
}

//-----------------------------------------------------------------------------
