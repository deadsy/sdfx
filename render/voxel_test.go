package render

import (
	"testing"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// VoxelSDF3 trilinear-interpolates corner values from a coarse grid. Even
// when the source SDF is Lipschitz-1, the trilinear interpolant's gradient
// can reach √(Lx² + Ly² + Lz²) where each Lᵢ approaches 1 — i.e. up to √3
// — so unrescaled it overstates 3D distance. The octree marching-cubes
// renderer's |sdf(center)| ≥ half-diagonal pruning then skips cubes that
// contain surface, producing holes. The fix divides Evaluate by the
// measured per-axis Lipschitz bound (computed at construction from
// adjacent-corner deltas).
//
// Tests cover the matrix that exposes this bug:
//   - several source SDFs (sphere, rounded box, cylinder, capsule,
//     sphere-minus-sphere shell)
//   - voxel cache resolutions from coarse (8) to fine (48), including
//     resolutions that mismatch the renderer's grid
//   - render resolutions that don't divide the voxel grid evenly
//   - a high-resolution stress pass on the most curvature-heavy source

type voxelCase struct {
	name     string
	src      func() sdf.SDF3
	voxCells int
}

func voxelSources(t *testing.T) []voxelCase {
	t.Helper()
	sphere := func(r float64) sdf.SDF3 {
		s, err := sdf.Sphere3D(r)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	roundedBox := func(size v3.Vec, round float64) sdf.SDF3 {
		s, err := sdf.Box3D(size, round)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	cyl := func(h, r, round float64) sdf.SDF3 {
		s, err := sdf.Cylinder3D(h, r, round)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	capsule := func(h, r float64) sdf.SDF3 {
		s, err := sdf.Capsule3D(h, r)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	shell := func() sdf.SDF3 {
		// Sphere-minus-sphere — non-convex, has interior surfaces the
		// voxel cache must reproduce on both sides of the shell.
		outer, _ := sdf.Sphere3D(2)
		inner, _ := sdf.Sphere3D(1.2)
		return sdf.Difference3D(outer, inner)
	}
	return []voxelCase{
		// Sphere — gradient varies smoothly; tightest case for the
		// √(Lx² + Ly² + Lz²) bound.
		{"sphere_r2_vox8", func() sdf.SDF3 { return sphere(2) }, 8},
		{"sphere_r2_vox16", func() sdf.SDF3 { return sphere(2) }, 16},
		{"sphere_r2_vox32", func() sdf.SDF3 { return sphere(2) }, 32},
		{"sphere_r2_vox48", func() sdf.SDF3 { return sphere(2) }, 48},
		// Rounded box — gradient discontinuities at face/edge transitions.
		{"rounded_box_4_3_2_vox16", func() sdf.SDF3 { return roundedBox(v3.Vec{X: 4, Y: 3, Z: 2}, 0.3) }, 16},
		{"rounded_box_4_3_2_vox32", func() sdf.SDF3 { return roundedBox(v3.Vec{X: 4, Y: 3, Z: 2}, 0.3) }, 32},
		// Cylinder — exercises mixed flat/curved surfaces.
		{"cyl_h4_r1_vox16", func() sdf.SDF3 { return cyl(4, 1, 0.1) }, 16},
		{"cyl_h4_r1_vox32", func() sdf.SDF3 { return cyl(4, 1, 0.1) }, 32},
		// Capsule — convex, smooth.
		{"capsule_h3_r1_vox16", func() sdf.SDF3 { return capsule(3, 1) }, 16},
		// Shell — non-convex with two surfaces.
		{"shell_vox24", shell, 24},
		{"shell_vox32", shell, 32},
	}
}

func Test_VoxelSDF3_Watertight(t *testing.T) {
	const renderCells = 80
	for _, c := range voxelSources(t) {
		t.Run(c.name, func(t *testing.T) {
			v := sdf.NewVoxelSDF3(c.src(), c.voxCells, nil)
			tris := CollectTriangles(v, NewMarchingCubesOctree(renderCells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("octree mesh has %d boundary edges (want 0); %d tris", be, len(tris))
			}
			t.Logf("%d tris, %d boundary edges", len(tris), be)
		})
	}
}

// Test_VoxelSDF3_VariedRenderRes runs a fixed source × voxel-cache
// configuration at several render resolutions, including ones that
// don't divide the voxel grid evenly. This exercises the case where
// many octree cube centers land at different positions within the
// voxel cells than the test above.
func Test_VoxelSDF3_VariedRenderRes(t *testing.T) {
	src, err := sdf.Sphere3D(2)
	if err != nil {
		t.Fatal(err)
	}
	v := sdf.NewVoxelSDF3(src, 24, nil)
	for _, cells := range []int{40, 60, 80, 100, 120, 150} {
		t.Run("cells="+itoa(cells), func(t *testing.T) {
			tris := CollectTriangles(v, NewMarchingCubesOctree(cells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("cells=%d: %d boundary edges (want 0); %d tris", cells, be, len(tris))
			}
			t.Logf("cells=%d: %d tris, %d boundary edges", cells, len(tris), be)
		})
	}
}

// itoa is a small local helper to keep the subtest names compact
// without pulling in fmt for a single integer.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var buf [12]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
