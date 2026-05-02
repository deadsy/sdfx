package render

import (
	"testing"

	"github.com/deadsy/sdfx/sdf"
)

// screwTestConfig defines a screw configuration for watertightness testing.
type screwTestConfig struct {
	name                  string
	radius, pitch, length float64
	taperDeg              float64
	starts, cells         int
	profile               string // "iso", "iso-int", "acme", "buttress", "plastic-buttress"
}

func makeTestThread(t *testing.T, profile string, radius, pitch float64) sdf.SDF2 {
	t.Helper()
	var thread sdf.SDF2
	var err error
	switch profile {
	case "iso-int":
		thread, err = sdf.ISOThread(radius, pitch, false)
	case "acme":
		thread, err = sdf.AcmeThread(radius, pitch)
	case "buttress":
		thread, err = sdf.ANSIButtressThread(radius, pitch)
	case "plastic-buttress":
		thread, err = sdf.PlasticButtressThread(radius, pitch)
	default:
		thread, err = sdf.ISOThread(radius, pitch, true)
	}
	if err != nil {
		t.Fatal(err)
	}
	return thread
}

func buildScrew(t *testing.T, c screwTestConfig) sdf.SDF3 {
	t.Helper()
	thread := makeTestThread(t, c.profile, c.radius, c.pitch)
	taperRad := sdf.DtoR(c.taperDeg)
	screw, err := sdf.Screw3D(thread, c.length, taperRad, c.pitch, c.starts)
	if err != nil {
		t.Fatal(err)
	}
	return screw
}

func straightScrewConfigs() []screwTestConfig {
	return []screwTestConfig{
		// ISO external — various pitch/radius ratios and resolutions
		{"M10x2_200", 5, 2, 20, 0, 1, 200, "iso"},
		{"M10x2_100", 5, 2, 20, 0, 1, 100, "iso"},
		{"M5x3_100", 2.5, 3, 10, 0, 1, 100, "iso"},
		{"coarse_M10x5", 5, 5, 20, 0, 1, 100, "iso"},
		{"steep_M3x3", 1.5, 3, 10, 0, 1, 100, "iso"},
		{"steep_M3x3_50", 1.5, 3, 10, 0, 1, 50, "iso"},
		{"extreme_M2x3", 1.0, 3, 6, 0, 1, 100, "iso"},
		{"extreme_M2x3_50", 1.0, 3, 6, 0, 1, 50, "iso"},
		{"M3x3_25cells", 1.5, 3, 10, 0, 1, 25, "iso"},
		{"short_1pitch", 5, 2, 2, 0, 1, 100, "iso"},
		{"large_M64x6", 32, 6, 60, 0, 1, 200, "iso"},
		// Multi-start
		{"multi8", 5, 2, 20, 0, 8, 100, "iso"},
		{"multi16", 5, 2, 20, 0, 16, 100, "iso"},
		{"multi8_coarse", 5, 5, 20, 0, 8, 100, "iso"},
		// Left-hand
		{"left_M10x2", 5, 2, 20, 0, -1, 100, "iso"},
		{"left_8start", 5, 2, 20, 0, -8, 100, "iso"},
		{"left_steep_M3x3", 1.5, 3, 10, 0, -1, 100, "iso"},
		// ISO internal
		{"internal_M10x2", 5, 2, 20, 0, 1, 100, "iso-int"},
		{"internal_M5x3", 2.5, 3, 10, 0, 1, 100, "iso-int"},
		// ACME
		{"acme_M10x2", 5, 2, 20, 0, 1, 100, "acme"},
		{"acme_steep_M5x3", 2.5, 3, 10, 0, 1, 100, "acme"},
		// Buttress
		{"buttress_M10x2", 5, 2, 20, 0, 1, 100, "buttress"},
		{"buttress_steep", 2.5, 3, 10, 0, 1, 100, "buttress"},
		{"buttress_left", 5, 2, 20, 0, -1, 100, "buttress"},
		{"buttress_multi4", 5, 2, 20, 0, 4, 100, "buttress"},
		{"buttress_25cells", 2.5, 3, 10, 0, 1, 25, "buttress"},
		// Plastic buttress
		{"plastic_butt_M10", 5, 2, 20, 0, 1, 100, "plastic-buttress"},
		{"plastic_butt_left", 5, 2, 20, 0, -1, 100, "plastic-buttress"},
		{"plastic_butt_multi4", 5, 2, 20, 0, 4, 100, "plastic-buttress"},
		// Dual/triple start (common real-world)
		{"dual_start", 5, 2, 20, 0, 2, 100, "iso"},
		{"triple_start", 5, 2, 20, 0, 3, 100, "iso"},
		// Very fine thread (low stretch factor, opposite extreme)
		{"fine_M20x0.5", 10, 0.5, 20, 0, 1, 100, "iso"},
		// Sub-pitch length
		{"sub_pitch_len", 5, 3, 2, 0, 1, 100, "iso"},
	}
}

func taperedScrewConfigs() []screwTestConfig {
	return []screwTestConfig{
		{"taper_1.8_NPT", 5, 2, 20, 1.79, 1, 100, "iso"},
		{"taper_5", 5, 2, 20, 5, 1, 100, "iso"},
		{"taper_15", 5, 2, 20, 15, 1, 100, "iso"},
		{"taper_30", 5, 2, 20, 30, 1, 100, "iso"},
		{"taper_45", 5, 2, 20, 45, 1, 100, "iso"},
		{"taper_30_4start", 5, 2, 20, 30, 4, 100, "iso"},
		{"taper_30_coarse", 5, 5, 20, 30, 1, 100, "iso"},
		{"taper_15_steep", 2.5, 3, 10, 15, 1, 100, "iso"},
		{"taper_30_50cells", 5, 2, 20, 30, 1, 50, "iso"},
		{"taper_30_left", 5, 2, 20, 30, -1, 100, "iso"},
		{"taper_15_internal", 5, 2, 20, 15, 1, 100, "iso-int"},
		{"taper_15_acme", 5, 2, 20, 15, 1, 100, "acme"},
		{"taper_30_buttress", 5, 2, 20, 30, 1, 100, "buttress"},
		{"taper_15_plastic_butt", 5, 2, 20, 15, 1, 100, "plastic-buttress"},
		{"taper_30_left_buttress", 5, 2, 20, 30, -1, 100, "buttress"},
		{"taper_15_multi4_buttress", 5, 2, 20, 15, 4, 100, "buttress"},
	}
}

func Test_Screw_Watertight_Straight(t *testing.T) {
	for _, c := range straightScrewConfigs() {
		t.Run(c.name, func(t *testing.T) {
			screw := buildScrew(t, c)
			tris := CollectTriangles(screw, NewMarchingCubesOctree(c.cells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("octree mesh has %d boundary edges (want 0 for watertight)", be)
			}
			t.Logf("%d tris, %d boundary edges", len(tris), be)
		})
	}
}

func Test_Screw_Watertight_Tapered(t *testing.T) {
	for _, c := range taperedScrewConfigs() {
		t.Run(c.name, func(t *testing.T) {
			screw := buildScrew(t, c)
			tris := CollectTriangles(screw, NewMarchingCubesOctree(c.cells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("octree mesh has %d boundary edges (want 0 for watertight)", be)
			}
			t.Logf("%d tris, %d boundary edges", len(tris), be)
		})
	}
}

func Test_Screw_EndCap_Position(t *testing.T) {
	// Verify octree and uniform renderers agree on end-cap Z positions.
	configs := straightScrewConfigs()
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			screw := buildScrew(t, c)
			halfLen := c.length / 2.0
			cubeSize := c.length / float64(c.cells)

			octTris := CollectTriangles(screw, NewMarchingCubesOctree(c.cells))
			uniTris := CollectTriangles(screw, NewMarchingCubesUniform(c.cells))

			octMaxZ := MaxZ(octTris)
			uniMaxZ := MaxZ(uniTris)

			if halfLen-octMaxZ > cubeSize {
				t.Errorf("octree maxZ=%.4f too far from halfLen=%.2f (err=%.4f, cubeSize=%.4f)",
					octMaxZ, halfLen, halfLen-octMaxZ, cubeSize)
			}
			if halfLen-uniMaxZ > cubeSize {
				t.Errorf("uniform maxZ=%.4f too far from halfLen=%.2f (err=%.4f, cubeSize=%.4f)",
					uniMaxZ, halfLen, halfLen-uniMaxZ, cubeSize)
			}
			delta := octMaxZ - uniMaxZ
			if delta > cubeSize || delta < -cubeSize {
				t.Errorf("octree/uniform disagree: octMaxZ=%.4f uniMaxZ=%.4f delta=%.4f (cubeSize=%.4f)",
					octMaxZ, uniMaxZ, delta, cubeSize)
			}
		})
	}
}
