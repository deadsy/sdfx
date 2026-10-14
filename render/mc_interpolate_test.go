package render

import (
	"fmt"
	"testing"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Without canonicalized mcInterpolate argument order, two adjacent
// cubes sharing a swap-ordered MC edge compute interpolated points
// that differ by ~1 ULP, and the int32(c · 1e4) vertex quantization
// in mesh.go buckets them separately — producing phantom boundary
// edges. The bug is most visible when the surface lands near
// 1e-4-aligned coordinates, e.g. axis-aligned box faces at integer
// or half-integer y/z values.
//
// Tests cover:
//   - thin-slab axis-aligned boxes at integer dimensions (multiple
//     aspect ratios so the failing edge can show up on any face)
//   - varied rounding values (0.0, 0.1, 0.3, 0.5)
//   - both renderers (octree and uniform)
//   - resolution sweep through the range that exposed the bug
//   - a sanity check on a non-axis-aligned shape (sphere) to confirm
//     the canonicalization doesn't break the common case

// Test_MCInterpolate_AxisAlignedBoxes_Stable sweeps box configurations
// where the surface lands near 1e-4-aligned coordinates and verifies
// no phantom boundary edges. Each (size, round) was specifically
// chosen because it places at least one face at a value like Y=1.0,
// Z=0.5, etc. — exactly the trap for the swap-ordered edge case.
func Test_MCInterpolate_AxisAlignedBoxes_Stable(t *testing.T) {
	cases := []struct {
		name  string
		size  v3.Vec
		round float64
	}{
		{"box_6_2_2_r0.3", v3.Vec{X: 6, Y: 2, Z: 2}, 0.3}, // original repro
		{"box_6_2_2_r0", v3.Vec{X: 6, Y: 2, Z: 2}, 0},     // no rounding
		{"box_6_2_2_r0.1", v3.Vec{X: 6, Y: 2, Z: 2}, 0.1},
		{"box_6_2_2_r0.5", v3.Vec{X: 6, Y: 2, Z: 2}, 0.5}, // near max round
		{"box_8_2_2_r0.2", v3.Vec{X: 8, Y: 2, Z: 2}, 0.2},
		{"box_4_4_2_r0.3", v3.Vec{X: 4, Y: 4, Z: 2}, 0.3},
		{"box_2_4_6_r0.3", v3.Vec{X: 2, Y: 4, Z: 6}, 0.3}, // long axis = z
		{"box_2_6_2_r0.3", v3.Vec{X: 2, Y: 6, Z: 2}, 0.3}, // long axis = y
		{"box_3_3_3_r0.5", v3.Vec{X: 3, Y: 3, Z: 3}, 0.5}, // cube
		{"box_4_2_1_r0.2", v3.Vec{X: 4, Y: 2, Z: 1}, 0.2}, // odd aspect
		{"box_10_4_2_r0.1", v3.Vec{X: 10, Y: 4, Z: 2}, 0.1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := sdf.Box3D(c.size, c.round)
			if err != nil {
				t.Fatal(err)
			}
			for _, cells := range []int{50, 75, 100, 125, 150, 200} {
				tris := CollectTriangles(b, NewMarchingCubesOctree(cells))
				if be := CountBoundaryEdges(tris); be != 0 {
					t.Errorf("octree c=%d: %d boundary edges (want 0); %d tris", cells, be, len(tris))
				}
			}
			for _, cells := range []int{50, 100, 125, 150} {
				tris := CollectTriangles(b, NewMarchingCubesUniform(cells))
				if be := CountBoundaryEdges(tris); be != 0 {
					t.Errorf("uniform c=%d: %d boundary edges (want 0); %d tris", cells, be, len(tris))
				}
			}
		})
	}
}

// Test_MCInterpolate_NonAxisAligned_Stable confirms that the
// canonicalization doesn't regress shapes where the bug never fired
// in the first place. A sphere doesn't put surfaces at 1e-4-aligned
// coordinates, so it was watertight before the fix and must remain so.
func Test_MCInterpolate_NonAxisAligned_Stable(t *testing.T) {
	cases := []struct {
		name string
		make func() sdf.SDF3
	}{
		{"sphere_r2", func() sdf.SDF3 { s, _ := sdf.Sphere3D(2); return s }},
		{"sphere_r1.7", func() sdf.SDF3 { s, _ := sdf.Sphere3D(1.7); return s }},
		{"cylinder_h4_r1", func() sdf.SDF3 { s, _ := sdf.Cylinder3D(4, 1, 0.1); return s }},
		{"capsule_h3_r1", func() sdf.SDF3 { s, _ := sdf.Capsule3D(3, 1); return s }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.make()
			for _, cells := range []int{60, 100, 150} {
				tris := CollectTriangles(s, NewMarchingCubesOctree(cells))
				if be := CountBoundaryEdges(tris); be != 0 {
					t.Errorf("octree c=%d: %d boundary edges (want 0)", cells, be)
				}
			}
		})
	}
}

// Test_MCInterpolate_Translated_Stable exercises a translated box
// where shared-edge swapping is more common (translation places face
// vertices at non-bbox-symmetric coordinates).
func Test_MCInterpolate_Translated_Stable(t *testing.T) {
	b, err := sdf.Box3D(v3.Vec{X: 4, Y: 2, Z: 2}, 0.2)
	if err != nil {
		t.Fatal(err)
	}
	translations := []v3.Vec{
		{X: 0, Y: 0, Z: 0},
		{X: 0.5, Y: 0, Z: 0},
		{X: 0, Y: 0.5, Z: 0},
		{X: 0, Y: 0, Z: 0.5},
		{X: 1.0, Y: 1.0, Z: 1.0},
		{X: -0.5, Y: -1.0, Z: 0.25},
	}
	for _, tv := range translations {
		t.Run(fmt.Sprintf("translate_%+v", tv), func(t *testing.T) {
			s := sdf.Transform3D(b, sdf.Translate3d(tv))
			for _, cells := range []int{75, 125} {
				tris := CollectTriangles(s, NewMarchingCubesOctree(cells))
				if be := CountBoundaryEdges(tris); be != 0 {
					t.Errorf("octree c=%d translate=%+v: %d boundary edges", cells, tv, be)
				}
			}
		})
	}
}
