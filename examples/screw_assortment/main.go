// Screw assortment: renders an array of screw configurations covering all
// thread profiles (ISO, ACME, buttress, plastic-buttress), multi-start,
// left-hand, tapered, and extreme dimensions. Verifies all octree-rendered
// meshes are watertight (zero boundary edges).
//
// Watertightness at extreme/low-resolution configurations is also covered
// by the unit tests in render/screw_test.go; this example focuses on a
// clean visual showcase of the full Screw3D configuration space.
//
// Usage:
//
//	go run main.go
package main

import (
	"fmt"
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

const meshCells = 300

type config struct {
	name                  string
	radius, pitch, length float64
	taperDeg              float64
	starts                int
	profile               string // "iso", "iso-int", "acme", "buttress", "plastic-buttress"
}

func makeThread(c config) sdf.SDF2 {
	var thread sdf.SDF2
	var err error
	switch c.profile {
	case "iso-int":
		thread, err = sdf.ISOThread(c.radius, c.pitch, false)
	case "acme":
		thread, err = sdf.AcmeThread(c.radius, c.pitch)
	case "buttress":
		thread, err = sdf.ANSIButtressThread(c.radius, c.pitch)
	case "plastic-buttress":
		thread, err = sdf.PlasticButtressThread(c.radius, c.pitch)
	default:
		thread, err = sdf.ISOThread(c.radius, c.pitch, true)
	}
	if err != nil {
		log.Fatal(err)
	}
	return thread
}

func buildScrew(c config) sdf.SDF3 {
	thread := makeThread(c)
	taperRad := sdf.DtoR(c.taperDeg)
	screw, err := sdf.Screw3D(thread, c.length, taperRad, c.pitch, c.starts)
	if err != nil {
		log.Fatal(err)
	}
	return screw
}

func allConfigs() []config {
	return []config{
		// --- Straight screws ---
		// ISO external
		{"M10x2", 5, 2, 20, 0, 1, "iso"},
		{"M5x3", 2.5, 3, 10, 0, 1, "iso"},
		{"coarse_M10x5", 5, 5, 20, 0, 1, "iso"},
		{"steep_M3x3", 1.5, 3, 10, 0, 1, "iso"},
		{"extreme_M2x3", 1.0, 3, 6, 0, 1, "iso"},
		{"short_1pitch", 5, 2, 2, 0, 1, "iso"},
		// Multi-start
		{"dual_start", 5, 2, 20, 0, 2, "iso"},
		{"triple_start", 5, 2, 20, 0, 3, "iso"},
		{"multi8", 5, 2, 20, 0, 8, "iso"},
		{"multi16", 5, 2, 20, 0, 16, "iso"},
		{"multi8_coarse", 5, 5, 20, 0, 8, "iso"},
		// Left-hand
		{"left_M10x2", 5, 2, 20, 0, -1, "iso"},
		{"left_8start", 5, 2, 20, 0, -8, "iso"},
		{"left_steep_M3x3", 1.5, 3, 10, 0, -1, "iso"},
		// ISO internal
		{"internal_M10x2", 5, 2, 20, 0, 1, "iso-int"},
		{"internal_M5x3", 2.5, 3, 10, 0, 1, "iso-int"},
		// ACME
		{"acme_M10x2", 5, 2, 20, 0, 1, "acme"},
		{"acme_steep_M5x3", 2.5, 3, 10, 0, 1, "acme"},
		// Buttress
		{"buttress_M10x2", 5, 2, 20, 0, 1, "buttress"},
		{"buttress_steep", 2.5, 3, 10, 0, 1, "buttress"},
		{"buttress_left", 5, 2, 20, 0, -1, "buttress"},
		{"buttress_multi4", 5, 2, 20, 0, 4, "buttress"},
		// Plastic buttress
		{"plastic_butt_M10", 5, 2, 20, 0, 1, "plastic-buttress"},
		{"plastic_butt_left", 5, 2, 20, 0, -1, "plastic-buttress"},
		{"plastic_butt_multi4", 5, 2, 20, 0, 4, "plastic-buttress"},
		// Fine thread
		{"fine_M20x0.5", 10, 0.5, 20, 0, 1, "iso"},
		// Sub-pitch length
		{"sub_pitch_len", 5, 3, 2, 0, 1, "iso"},

		// --- Tapered screws ---
		{"taper_1.8_NPT", 5, 2, 20, 1.79, 1, "iso"},
		{"taper_5", 5, 2, 20, 5, 1, "iso"},
		{"taper_15", 5, 2, 20, 15, 1, "iso"},
		{"taper_30", 5, 2, 20, 30, 1, "iso"},
		{"taper_45", 5, 2, 20, 45, 1, "iso"},
		{"taper_30_4start", 5, 2, 20, 30, 4, "iso"},
		{"taper_30_coarse", 5, 5, 20, 30, 1, "iso"},
		{"taper_15_steep", 2.5, 3, 10, 15, 1, "iso"},
		{"taper_30_left", 5, 2, 20, 30, -1, "iso"},
		{"taper_15_internal", 5, 2, 20, 15, 1, "iso-int"},
		{"taper_15_acme", 5, 2, 20, 15, 1, "acme"},
		{"taper_30_buttress", 5, 2, 20, 30, 1, "buttress"},
		{"taper_15_plastic_butt", 5, 2, 20, 15, 1, "plastic-buttress"},
		{"taper_30_left_buttress", 5, 2, 20, 30, -1, "buttress"},
		{"taper_15_multi4_buttress", 5, 2, 20, 15, 4, "buttress"},

		// --- Oversized (placed last so it doesn't dominate the layout) ---
		{"large_M64x6", 32, 6, 60, 0, 1, "iso"},
	}
}

func main() {
	configs := allConfigs()
	allPassed := true

	// Lay out screws in a row along X. Gap scales with the larger of the
	// neighboring radii so oversized screws don't crowd small ones.
	var allTris []*sdf.Triangle3
	xOffset := 0.0
	prevR := 0.0

	for i, c := range configs {
		screw := buildScrew(c)
		bb := screw.BoundingBox()
		r := bb.Max.X

		tris := render.CollectTriangles(screw, render.NewMarchingCubesOctree(meshCells))
		be := render.CountBoundaryEdges(tris)

		status := "PASS"
		if be > 0 {
			status = "FAIL"
			allPassed = false
		}

		fmt.Printf("%-28s r=%5.1f p=%3.1f taper=%5.1f° starts=%3d → %6d tris, %4d boundary edges [%s]\n",
			c.name, c.radius, c.pitch, c.taperDeg, c.starts, len(tris), be, status)

		if i > 0 {
			gap := 3.0
			if s := max(prevR, r) * 0.5; s > gap {
				gap = s
			}
			xOffset += gap
		}

		for j := range tris {
			t := tris[j]
			t[0].X += xOffset + r
			t[1].X += xOffset + r
			t[2].X += xOffset + r
			allTris = append(allTris, &t)
		}
		xOffset += r * 2
		prevR = r
	}

	fmt.Println()
	if allPassed {
		fmt.Println("ALL PASSED: all octree meshes are watertight")
	} else {
		fmt.Println("SOME FAILED: octree meshes have holes (boundary edges)")
	}

	fmt.Printf("\nrendering screws.stl (%d triangles)\n", len(allTris))
	if err := render.SaveSTL("screws.stl", allTris); err != nil {
		log.Fatal(err)
	}
}
