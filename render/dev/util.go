package dev

import (
	"fmt"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font/inconsolata"
	"image/color"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func nextPowerOf2(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func utilSdf2MinMax(s sdf.SDF2, bb sdf.Box2, cells sdf.V2i) (dmin, dmax float64) {
	cellSize := bb.Size().Div(cells.ToV2())
	for x := 0; x < cells[0]; x++ {
		for y := 0; y < cells[1]; y++ {
			// TODO: Reverse raycast (without limiting to a single direction) to find extreme values instead of 0s (should lower sample count for same results)
			pos := bb.Min.Add((sdf.V2{X: float64(x), Y: float64(y)}).Mul(cellSize))
			d := s.Evaluate(pos)
			dmax = math.Max(dmax, d)
			dmin = math.Min(dmin, d)
		}
	}
	return
}

var dirs2 = []sdf.V2i{
	{0, -1},
	{-1, 0},
	{0, 1},
	{1, 0},
}

var defaultFont = inconsolata.Regular8x16 // Just a simple embedded font (to avoid problems with some platforms)

func drawDefaultTextWithShadow(screen *ebiten.Image, msg string, x, y int, c color.Color) {
	for dx := -1; dx <= 1; dx += 1 {
		for dy := -1; dy <= 1; dy += 1 {
			text.Draw(screen, msg, defaultFont, x+dx, y+dy, color.RGBA{R: 0, G: 0, B: 0, A: 50}) // Shadow first (background)
		}
	}
	text.Draw(screen, msg, defaultFont, x, y, c)
}

func toBox2(box3 sdf.Box3) sdf.Box2 {
	return sdf.Box2{
		Min: sdf.V2{X: box3.Min.X, Y: box3.Min.Y},
		Max: sdf.V2{X: box3.Max.X, Y: box3.Max.Y},
	}
}

func pidExists(pid int32) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("invalid pid %v", pid)
	}
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return false, err
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		if runtime.GOOS == "linux" { // Fix/hack for my system on which this seems broken (go 1.17.5)
			cmd := exec.Command("kill", "-0", strconv.FormatInt(int64(pid), 10))
			err := cmd.Run()
			return err != nil, nil
		}
		return true, nil
	}
	if err.Error() == "os: process already finished" {
		return false, nil
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false, err
	}
	switch errno {
	case syscall.ESRCH:
		return false, nil
	case syscall.EPERM:
		return true, nil
	}
	return false, err
}

func pidTermWaitKill(proc *os.Process, gracePeriod time.Duration) error {
	err := proc.Signal(os.Interrupt)
	if err != nil {
		return err
	}
	afterFunc := time.AfterFunc(gracePeriod, func() {
		err = proc.Kill()
		if err != nil {
			log.Println("[DevRenderer] pidTermWaitKill proc.Kill error:", err)
		}
	})
	defer afterFunc.Stop()
	_, err = proc.Wait()
	return err
}

//type boundedSDF3 struct {
//	sdf.SDF3
//	Bb sdf.Box3
//}
//
//func (b *boundedSDF3) BoundingBox() sdf.Box3 {
//	return b.Bb
//}
