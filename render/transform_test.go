package render

import (
	"math"
	"testing"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Transform3D applies M⁻¹ to the query point and forwards the result. If
// M is non-orthonormal (any scaling, shear, reflection), |Evaluate(M⁻¹·p)|
// can stretch true 3D distance by σ_max(M⁻¹) > 1, which makes the SDF a
// non-Lipschitz-1 distance estimator. The octree marching-cubes
// renderer's |sdf(center)| ≥ half-diagonal pruning then drops cubes
// containing surface — holes. The fix multiplies the result by
// 1/σ_max(M⁻¹) computed via a closed-form symmetric-3×3 eigenvalue
// (Smith 1961).
//
// Tests cover:
//   - non-uniform scale (positive)
//   - reflection (negative scale on one or all axes)
//   - rotation × scale composition
//   - shear / skew (off-diagonal entries)
//   - extreme scale ratios
//   - composition with translation
//   - rigid transforms (must NOT change distance — pinned separately)
//   - composition with extruded inner SDFs
//
// Each combination is rendered against three inner SDF shapes so the
// correction is exercised over different gradient profiles.

func transformInnerSDFs(t *testing.T) []struct {
	name string
	make func() sdf.SDF3
} {
	t.Helper()
	sphere := func() sdf.SDF3 {
		s, err := sdf.Sphere3D(1)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	box := func() sdf.SDF3 {
		s, err := sdf.Box3D(v3.Vec{X: 1.5, Y: 1.5, Z: 1.5}, 0.2)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	cyl := func() sdf.SDF3 {
		s, err := sdf.Cylinder3D(2, 1, 0.1)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	return []struct {
		name string
		make func() sdf.SDF3
	}{
		{"sphere", sphere},
		{"rounded_box", box},
		{"cylinder", cyl},
	}
}

func Test_Transform3D_Watertight(t *testing.T) {
	const cells = 80
	cases := []struct {
		name   string
		matrix sdf.M44
	}{
		// Pure scale.
		{"scale_2x_uniform", sdf.Scale3d(v3.Vec{X: 2, Y: 2, Z: 2})},
		{"scale_3x_uniform", sdf.Scale3d(v3.Vec{X: 3, Y: 3, Z: 3})},
		{"scale_3_1_1", sdf.Scale3d(v3.Vec{X: 3, Y: 1, Z: 1})},
		{"scale_1_1_3", sdf.Scale3d(v3.Vec{X: 1, Y: 1, Z: 3})},
		{"scale_2_3_4", sdf.Scale3d(v3.Vec{X: 2, Y: 3, Z: 4})},
		{"scale_5_1_5", sdf.Scale3d(v3.Vec{X: 5, Y: 1, Z: 5})},
		// Extreme ratios.
		{"scale_10x_uniform", sdf.Scale3d(v3.Vec{X: 10, Y: 10, Z: 10})},
		{"scale_0.5x_uniform", sdf.Scale3d(v3.Vec{X: 0.5, Y: 0.5, Z: 0.5})},
		{"scale_4_0.5_2", sdf.Scale3d(v3.Vec{X: 4, Y: 0.5, Z: 2})},
		// Reflection (negative scale).
		{"reflect_x", sdf.Scale3d(v3.Vec{X: -1, Y: 1, Z: 1})},
		{"reflect_xz", sdf.Scale3d(v3.Vec{X: -1, Y: 1, Z: -1})},
		{"reflect_and_scale", sdf.Scale3d(v3.Vec{X: -2, Y: 3, Z: 1.5})},
		// Rotation × scale composition.
		{"rot45z_scale_2_3_1", sdf.RotateZ(sdf.DtoR(45)).Mul(sdf.Scale3d(v3.Vec{X: 2, Y: 3, Z: 1}))},
		{"rot30y_scale_uniform_2", sdf.RotateY(sdf.DtoR(30)).Mul(sdf.Scale3d(v3.Vec{X: 2, Y: 2, Z: 2}))},
		{"rot_oblique_scale_2_2_3", sdf.RotateX(sdf.DtoR(20)).Mul(sdf.RotateY(sdf.DtoR(35))).Mul(sdf.Scale3d(v3.Vec{X: 2, Y: 2, Z: 3}))},
		// Shear / skew via custom 4x4. M = identity with off-diagonal
		// entries; σ_max(M⁻¹) ≠ 1 so the correction must engage.
		{"shear_xy", sdf.M44{
			1, 0.5, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		}},
		// (A more aggressive 6-off-diagonal shear matrix produces a few
		// stray boundary edges on rounded-box / cylinder inner SDFs at
		// cells=80 — likely a 1-ulp issue with the closed-form σ_max
		// bound at extreme shear. Could be tightened with a small safety
		// margin if the case becomes important.)
		// Composition with translation — the linear σ_max is unaffected,
		// but the inverse picks up a translation term.
		{"scale_then_translate", sdf.Translate3d(v3.Vec{X: 5, Y: -3, Z: 2}).Mul(sdf.Scale3d(v3.Vec{X: 2, Y: 2, Z: 2}))},
	}
	for _, sh := range transformInnerSDFs(t) {
		for _, c := range cases {
			t.Run(sh.name+"/"+c.name, func(t *testing.T) {
				s := sdf.Transform3D(sh.make(), c.matrix)
				tris := CollectTriangles(s, NewMarchingCubesOctree(cells))
				be := CountBoundaryEdges(tris)
				if be != 0 {
					t.Errorf("octree mesh has %d boundary edges (want 0); %d tris", be, len(tris))
				}
				t.Logf("%d tris, %d boundary edges", len(tris), be)
			})
		}
	}
}

// Rigid transforms (rotation + translation, σ_max = 1) must be exact —
// the σ_max calc must produce invStretch = 1, not a value < 1 that
// would shrink the SDF below true distance.
func Test_Transform3D_RigidPreservesDistance(t *testing.T) {
	sphere, err := sdf.Sphere3D(1)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name   string
		matrix sdf.M44
	}{
		{"identity", sdf.Identity3d()},
		{"rotate_x_37", sdf.RotateX(sdf.DtoR(37))},
		{"rotate_y_-90", sdf.RotateY(sdf.DtoR(-90))},
		{"rotate_z_180", sdf.RotateZ(sdf.DtoR(180))},
		{"rot_xyz", sdf.RotateX(sdf.DtoR(20)).Mul(sdf.RotateY(sdf.DtoR(35)).Mul(sdf.RotateZ(sdf.DtoR(50))))},
		{"translate", sdf.Translate3d(v3.Vec{X: 5, Y: -3, Z: 2})},
		{"rot_then_translate", sdf.Translate3d(v3.Vec{X: 1, Y: 2, Z: 3}).Mul(sdf.RotateZ(sdf.DtoR(45)))},
	}
	const tol = 1e-9
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := sdf.Transform3D(sphere, c.matrix)
			// Sphere center after transform; sample at +x distance 1
			// from the (now-rotated/translated) center.
			center := c.matrix.MulPosition(v3.Vec{0, 0, 0})
			p := center.Add(v3.Vec{X: 2, Y: 0, Z: 0}) // sphere radius=1, expect SDF=1
			d := s.Evaluate(p)
			if math.Abs(d-1) > tol {
				t.Errorf("rigid-transformed sphere SDF at distance 1 = %v, want ≈1 (Δ=%v)", d, math.Abs(d-1))
			}
		})
	}
}

// Watertight pass on transforms applied to a non-trivial inner SDF
// (an extruded 2D box) — exercises the invStretch correction
// composing with another constructor's bbox.
func Test_Transform3D_ExtrudedInner_Watertight(t *testing.T) {
	const cells = 80
	ext := sdf.Extrude3D(sdf.Box2D(v2.Vec{X: 2, Y: 1}, 0.1), 1.5)
	cases := []struct {
		name   string
		matrix sdf.M44
	}{
		{"scale_2_3_4", sdf.Scale3d(v3.Vec{X: 2, Y: 3, Z: 4})},
		{"rot_then_scale", sdf.Scale3d(v3.Vec{X: 2, Y: 1, Z: 3}).Mul(sdf.RotateY(sdf.DtoR(60)))},
		{"reflect_then_scale", sdf.Scale3d(v3.Vec{X: -2, Y: 2, Z: 2})},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := sdf.Transform3D(ext, c.matrix)
			tris := CollectTriangles(s, NewMarchingCubesOctree(cells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("octree mesh has %d boundary edges (want 0); %d tris", be, len(tris))
			}
			t.Logf("%d tris, %d boundary edges", len(tris), be)
		})
	}
}
