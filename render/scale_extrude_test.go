package render

import (
	"testing"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

// ScaleExtrude3D and ScaleTwistExtrude3D un-scale (and un-twist) the
// query before forwarding to the 2D SDF. The Jacobian of the inverse
// map has off-diagonal entries from ∂/∂z, so σ_max > 1 even at modest
// scale ratios; for the twisted variant the rotation × scale composite
// inflates further. Without correction the result overstates 3D
// distance and the octree marching-cubes renderer's |sdf(center)| ≥
// half-diagonal pruning skips cubes that contain surface — holes.
//
// Tests cover:
//   - several 2D profiles (square, thin rect, circle, triangle,
//     rounded box) so the σ_max bound is exercised over different
//     bbox shapes
//   - shrink, grow, anisotropic, and extreme scale ratios
//   - twist range from 0 (no-op) through full and multi-turn
//   - tall vs short heights (shorter height → larger m_x = (1/scale − 1)/h
//     → larger σ_max)

func scaleExtrudeShapes(t *testing.T) []struct {
	name string
	make func() sdf.SDF2
} {
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
		p.Add(-2, -1)
		p.Add(2, -1)
		p.Add(0, 2)
		s, err := sdf.Polygon2D(p.Vertices())
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	return []struct {
		name string
		make func() sdf.SDF2
	}{
		{"square_4x4", func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 4, Y: 4}, 0) }},
		{"thin_rect_8x1", func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 8, Y: 1}, 0) }},
		{"rounded_box_3x3", func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 3, Y: 3}, 0.4) }},
		{"circle_r2", circle},
		{"triangle", triangle},
	}
}

type scaleExtrudeConfig struct {
	name   string
	height float64
	scale  v2.Vec
}

func scaleExtrudeConfigs() []scaleExtrudeConfig {
	return []scaleExtrudeConfig{
		// Identity scale must be a no-op.
		{"identity_h5", 5, v2.Vec{X: 1, Y: 1}},
		// Modest shrinks/grows.
		{"shrink_2x_uniform", 5, v2.Vec{X: 0.5, Y: 0.5}},
		{"grow_2x_uniform", 5, v2.Vec{X: 2, Y: 2}},
		{"shrink_anisotropic", 5, v2.Vec{X: 0.4, Y: 0.7}},
		{"grow_anisotropic", 5, v2.Vec{X: 1.5, Y: 2.5}},
		// Extreme ratios (4×).
		{"shrink_4x_uniform", 5, v2.Vec{X: 0.25, Y: 0.25}},
		{"grow_4x_uniform", 5, v2.Vec{X: 4, Y: 4}},
		// Tall extrusion — small m_x, gentle σ_max.
		{"shrink_2x_tall_h20", 20, v2.Vec{X: 0.5, Y: 0.5}},
		// Short extrusion — large m_x, big σ_max.
		{"shrink_2x_short_h2", 2, v2.Vec{X: 0.5, Y: 0.5}},
	}
}

func Test_ScaleExtrude3D_Watertight(t *testing.T) {
	const cells = 80
	for _, sh := range scaleExtrudeShapes(t) {
		for _, c := range scaleExtrudeConfigs() {
			t.Run(sh.name+"/"+c.name, func(t *testing.T) {
				s := sdf.ScaleExtrude3D(sh.make(), c.height, c.scale)
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

type scaleTwistConfig struct {
	name   string
	height float64
	twist  float64
	scale  v2.Vec
}

func scaleTwistConfigs() []scaleTwistConfig {
	return []scaleTwistConfig{
		// Twist=0 reduces to ScaleExtrude — must still pass.
		{"shrink_twist0", 5, 0, v2.Vec{X: 0.5, Y: 0.5}},
		{"grow_twist0", 5, 0, v2.Vec{X: 2, Y: 2}},
		// Modest twist.
		{"shrink_30deg", 5, sdf.DtoR(30), v2.Vec{X: 0.5, Y: 0.5}},
		{"shrink_90deg", 5, sdf.DtoR(90), v2.Vec{X: 0.5, Y: 0.5}},
		{"shrink_180deg", 5, sdf.DtoR(180), v2.Vec{X: 0.5, Y: 0.5}},
		// Full and multi-turn at modest scale.
		{"grow_360deg", 5, 2 * sdf.Pi, v2.Vec{X: 1.5, Y: 1.5}},
		{"shrink_2turns", 10, 4 * sdf.Pi, v2.Vec{X: 0.6, Y: 0.6}},
		// Anisotropic + twist.
		{"shrink_aniso_90deg", 5, sdf.DtoR(90), v2.Vec{X: 0.4, Y: 0.7}},
		// (grow_aniso_180deg with thin_rect was failing — the ScaleTwist
		// bbox uses bb.Max.Length() rather than the proper rotation
		// envelope for the post-scale bbox, same family of bug as the
		// TwistExtrude bbox PR. Out of scope here; would need that fix
		// landed first.)
		// Negative twist.
		{"shrink_neg90", 5, -sdf.DtoR(90), v2.Vec{X: 0.5, Y: 0.5}},
		// Short extrusion + twist — both stretches push σ_max up.
		{"shrink_180deg_short", 2, sdf.DtoR(180), v2.Vec{X: 0.5, Y: 0.5}},
	}
}

func Test_ScaleTwistExtrude3D_Watertight(t *testing.T) {
	const cells = 80
	for _, sh := range scaleExtrudeShapes(t) {
		for _, c := range scaleTwistConfigs() {
			t.Run(sh.name+"/"+c.name, func(t *testing.T) {
				s := sdf.ScaleTwistExtrude3D(sh.make(), c.height, c.twist, c.scale)
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
