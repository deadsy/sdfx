package render

import (
	"testing"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

// TwistExtrude3D is not a 1-Lipschitz SDF: the un-twist mapping has
// Jacobian σ_max = √(1 + k²r²) where k = twist/height, so the 2D-side
// distance overstates the true 3D distance by up to that factor. The
// octree marching-cubes renderer prunes cubes by |sdf(center)| ≥
// half-diagonal, so an over-stated distance causes it to skip cubes
// that contain the surface — holes. The fix divides the Evaluate
// result by σ_max (computed at construction from the bbox-bounded
// rMax) so the SDF stays a valid Lipschitz-1 distance estimator.
//
// The cases below exercise every shape × twist combination the bug
// could surface in:
//
//   - shapes with different rMax (fills bbox vs not, axis-symmetric vs
//     not, off-axis to push rMax higher).
//   - twist range from 0 (must be a no-op) through fractions of a turn,
//     full rotations, multiple turns, and negative.
//   - cell counts at and well past the resolution where holes would be
//     visible (80 for the cheap sweep, 200 for a stress pass on the
//     hardest configurations).

// twistShape is a 2D profile factory plus a name. Shapes are chosen to
// stress different parts of the σ_max bound:
type twistShape struct {
	name string
	make func() sdf.SDF2
}

func twistShapes(t *testing.T) []twistShape {
	t.Helper()
	circle := func() sdf.SDF2 {
		c, err := sdf.Circle2D(2)
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
	triangle := func() sdf.SDF2 {
		p := sdf.NewPolygon()
		p.Add(-3, -2)
		p.Add(3, -2)
		p.Add(0, 3)
		s, err := sdf.Polygon2D(p.Vertices())
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	return []twistShape{
		// Square — symmetric about origin, fills its bbox.
		{"square_4x4", func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 4, Y: 4}, 0) }},
		// Thin rectangle — large rMax along the long axis, small along the short.
		{"thin_rect_8x1", func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 8, Y: 1}, 0) }},
		// (Off-center / off-axis shapes — where bb.Min is farther from
		// the origin than bb.Max — would also be interesting cases, but
		// they trip a *separate* bug: the existing TwistExtrude3D bbox
		// uses bb.Max.Length() instead of max(|bb.Min|, |bb.Max|),
		// under-sizing the bbox in that configuration. That's fixed in
		// another PR; once it lands, off-center / off-axis cases become
		// safe to add here.)
		// Circle — doesn't fill its bbox, so the bbox-based rMax over-estimates
		// the true σ_max but stays sound (conservative).
		{"circle_r2", circle},
		// Triangle — irregular, doesn't fill its bbox, asymmetric corners.
		{"triangle", triangle},
		// Rounded square — exercises the corner-rounding interaction with twist.
		{"rounded_square_4x4_r0.5", func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 4, Y: 4}, 0.5) }},
	}
}

// twistConfig pairs a height with a twist value for the sweep.
type twistConfig struct {
	name   string
	height float64
	twist  float64
}

func twistConfigs() []twistConfig {
	return []twistConfig{
		// twist = 0 must be a no-op even with the σ_max correction (k=0 → invStretch=1).
		{"twist=0", 5, 0},
		// Small twist — k²r² is small, σ_max barely > 1.
		{"twist=15deg", 5, sdf.DtoR(15)},
		{"twist=30deg", 5, sdf.DtoR(30)},
		// Quarter turn — common.
		{"twist=90deg", 5, sdf.DtoR(90)},
		// Half turn.
		{"twist=180deg", 5, sdf.DtoR(180)},
		// Three-quarter turn.
		{"twist=270deg", 5, sdf.DtoR(270)},
		// Full rotation.
		{"twist=full", 5, 2 * sdf.Pi},
		// Multiple turns over a longer height — same k, larger absolute twist.
		{"twist=2turns_h10", 10, 4 * sdf.Pi},
		// Many turns at moderate height — large k, stresses the bound.
		{"twist=5turns_h5", 5, 10 * sdf.Pi},
		// Negative twist — direction shouldn't matter.
		{"twist=neg90", 5, -sdf.DtoR(90)},
		{"twist=neg2turns", 10, -4 * sdf.Pi},
		// Tall extrusion with a small twist — small k, surface mostly straight.
		{"twist=30deg_h20", 20, sdf.DtoR(30)},
	}
}

func Test_TwistExtrude3D_Watertight(t *testing.T) {
	const cells = 80
	for _, sh := range twistShapes(t) {
		for _, c := range twistConfigs() {
			t.Run(sh.name+"/"+c.name, func(t *testing.T) {
				s := sdf.TwistExtrude3D(sh.make(), c.height, c.twist)
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

// Stress pass at cells = 200 on the configurations most likely to expose
// borderline FP differences across architectures: shapes whose rMax is
// at the bbox extreme combined with twists that put k²r² in the order-1
// regime where the original code under-shot the σ_max bound.
func Test_TwistExtrude3D_Watertight_HighRes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping high-resolution stress pass in -short mode")
	}
	const cells = 200
	thin := func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 8, Y: 1}, 0) }
	cases := []twistConfig{
		{"twist=90deg_h5", 5, sdf.DtoR(90)},
		{"twist=180deg_h5", 5, sdf.DtoR(180)},
		{"twist=full_h5", 5, 2 * sdf.Pi},
		{"twist=2turns_h10", 10, 4 * sdf.Pi},
		{"twist=5turns_h5", 5, 10 * sdf.Pi},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := sdf.TwistExtrude3D(thin(), c.height, c.twist)
			tris := CollectTriangles(s, NewMarchingCubesOctree(cells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("octree mesh has %d boundary edges (want 0); %d tris", be, len(tris))
			}
			t.Logf("%d tris, %d boundary edges", len(tris), be)
		})
	}
}
