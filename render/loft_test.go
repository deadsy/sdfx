package render

import (
	"math"
	"testing"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

// Loft3D linearly blends between two 2D SDFs along z. The xy gradient is
// a convex combination of the two Lipschitz-1 inner gradients (so still
// ≤ 1), but the z gradient grows with max|a1−a0|/(2·h_inner). The outer
// d = √(a² + b²) end-cap composition inflates the Lipschitz constant
// further. Without correction the result overstates 3D distance and
// the octree marching-cubes renderer's |sdf(center)| ≥ half-diagonal
// pruning skips cubes that contain surface, producing holes.
//
// Tests cover:
//   - identity loft (sdf0 == sdf1) — must give Lipschitz unchanged
//   - same-shape size change (small→big, big→small)
//   - cross-shape lofts (box↔circle, triangle→circle, star)
//   - extreme size ratios
//   - extreme height ratios (short height → large L_z)
//   - with and without rounding
//   - thin shapes — narrower bbox stresses the L_z bound

func loftShapes(t *testing.T) map[string]func() sdf.SDF2 {
	t.Helper()
	circle := func(r float64) func() sdf.SDF2 {
		return func() sdf.SDF2 {
			c, err := sdf.Circle2D(r)
			if err != nil {
				t.Fatal(err)
			}
			return c
		}
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
	star := func() sdf.SDF2 {
		// 5-point star with outer radius 2, inner radius 0.8.
		p := sdf.NewPolygon()
		for i := 0; i < 10; i++ {
			r := 2.0
			if i%2 == 1 {
				r = 0.8
			}
			a := float64(i) * math.Pi / 5
			p.Add(r*math.Cos(a), r*math.Sin(a))
		}
		s, err := sdf.Polygon2D(p.Vertices())
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	return map[string]func() sdf.SDF2{
		"box_1x1":      func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 1, Y: 1}, 0.1) },
		"box_3x3":      func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 3, Y: 3}, 0.1) },
		"box_5x5":      func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 5, Y: 5}, 0.1) },
		"box_3x1_thin": func() sdf.SDF2 { return sdf.Box2D(v2.Vec{X: 3, Y: 1}, 0.05) },
		"circle_r0.5":  circle(0.5),
		"circle_r2":    circle(2),
		"circle_r3":    circle(3),
		"triangle":     triangle,
		"star":         star,
	}
}

func Test_Loft3D_Watertight(t *testing.T) {
	const cells = 80
	shapes := loftShapes(t)
	type loftCase struct {
		name          string
		sdf0, sdf1    string
		height, round float64
	}
	cases := []loftCase{
		// Identity loft — sdf0 == sdf1, max|a1-a0| = 0, L_z = 0,
		// invStretch = 1. Pin that this path is a no-op.
		{"identity_box3", "box_3x3", "box_3x3", 5, 0.05},
		{"identity_circle", "circle_r2", "circle_r2", 5, 0.05},
		// Same-shape size change.
		{"box_small_to_big_h5", "box_1x1", "box_5x5", 5, 0.05},
		{"box_small_to_big_h2_short", "box_1x1", "box_5x5", 2, 0.05},
		{"box_small_to_big_h0.5_extreme", "box_1x1", "box_5x5", 0.5, 0.02},
		{"circle_growing_h3", "circle_r0.5", "circle_r3", 3, 0.05},
		{"circle_growing_h1_short", "circle_r0.5", "circle_r3", 1, 0},
		// Cross-shape lofts (topology change).
		{"box_to_circle_h5", "box_5x5", "circle_r3", 5, 0.05},
		{"circle_to_box_h5", "circle_r3", "box_5x5", 5, 0.05},
		{"triangle_to_circle_h5", "triangle", "circle_r2", 5, 0.05},
		// Star — non-convex, multiple convex hulls.
		{"star_to_circle_h5", "star", "circle_r2", 5, 0.05},
		{"box_to_star_h5", "box_3x3", "star", 5, 0.05},
		// Extreme size ratios.
		{"tiny_to_huge_circle_h5", "circle_r0.5", "circle_r3", 5, 0.05},
		// Thin shape lofts — narrower bbox stresses the L_z bound.
		{"thin_box_to_thin_box_h5", "box_3x1_thin", "box_3x1_thin", 5, 0.02},
		// Larger rounding.
		{"box_loft_with_rounding", "box_1x1", "box_5x5", 5, 0.2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := sdf.Loft3D(shapes[c.sdf0](), shapes[c.sdf1](), c.height, c.round)
			if err != nil {
				t.Fatal(err)
			}
			tris := CollectTriangles(s, NewMarchingCubesOctree(cells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("octree mesh has %d boundary edges (want 0); %d tris", be, len(tris))
			}
			t.Logf("%d tris, %d boundary edges", len(tris), be)
		})
	}
}
