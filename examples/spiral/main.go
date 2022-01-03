//-----------------------------------------------------------------------------
/*

Spirals

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/dev"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/profile"
	"log"
	"os"
	"os/exec"
)

//-----------------------------------------------------------------------------

func spiralSdf() (sdf.SDF2, error) {
	s, err := sdf.ArcSpiral2D(1.0, 20.0, 0.25*sdf.Pi, 8*sdf.Tau, 1.0)
	if err != nil {
		return nil, err
	}

	c, err := sdf.Circle2D(22.)
	if err != nil {
		return nil, err
	}
	s = sdf.Union2D(s, c)

	c2, err := sdf.Circle2D(20.)
	if err != nil {
		return nil, err
	}
	c2 = sdf.Transform2D(c2, sdf.Translate2d(sdf.V2{X: 0}))
	s = sdf.Difference2D(s, c2)
	return s, err
}

func main() {
	s, err := spiralSdf()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	if os.Getenv("SDFX_TEST_DEV_RENDERER_2D") != "" {
		// Profiling boilerplate
		defer func() {
			//cmd := exec.Command("go", "tool", "pprof", "cpu.pprof")
			cmd := exec.Command("go", "tool", "trace", "trace.out")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
		}()
		//defer profile.Start(profile.ProfilePath(".")).Stop()
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()

		// Actual rendering loop
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
