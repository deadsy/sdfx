package render

import (
	"testing"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Gyroid3D is not a true SDF — its implicit function
//   g(p) = sin(kx·x)·cos(ky·y) + sin(ky·y)·cos(kz·z) + sin(kz·z)·cos(kx·x)
// has gradient up to √3·max(kᵢ), so unrescaled it overstates 3D distance
// by that factor. The octree marching-cubes renderer's isEmpty pruning
// (|sdf(center)| ≥ half-diagonal) then skips cubes that contain surface,
// producing holes. The fix divides Evaluate by √3·max(kᵢ).
//
// Gyroid is infinite, so we intersect with a finite shape to get a
// renderable surface. Tests sweep:
//   - isotropic vs. anisotropic vs. extreme-anisotropic scales
//   - fine (high-frequency, kMax large) and coarse (large period) scales
//   - several clipping shapes (box, sphere, rounded box) so the gyroid
//     surface meets different boundary geometries

// gyroidScales gives the matrix of scale configurations the σ_max
// correction has to handle.
var gyroidScales = []struct {
	name  string
	scale v3.Vec
}{
	{"isotropic_3", v3.Vec{X: 3, Y: 3, Z: 3}},         // moderate period
	{"isotropic_2", v3.Vec{X: 2, Y: 2, Z: 2}},         // medium period
	{"isotropic_1", v3.Vec{X: 1, Y: 1, Z: 1}},         // fine; kMax = 2π
	{"isotropic_0.7", v3.Vec{X: 0.7, Y: 0.7, Z: 0.7}}, // very fine
	{"isotropic_5", v3.Vec{X: 5, Y: 5, Z: 5}},         // coarse
	{"aniso_2_3_4", v3.Vec{X: 2, Y: 3, Z: 4}},         // mild anisotropy
	{"aniso_1_4_2", v3.Vec{X: 1, Y: 4, Z: 2}},         // x dominates kMax
	{"aniso_extreme", v3.Vec{X: 0.5, Y: 5, Z: 1}},     // x = kMax (≈ 4π)
}

func Test_Gyroid3D_Watertight(t *testing.T) {
	const cells = 80
	box, err := sdf.Box3D(v3.Vec{X: 10, Y: 10, Z: 10}, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range gyroidScales {
		t.Run(c.name, func(t *testing.T) {
			g, err := sdf.Gyroid3D(c.scale)
			if err != nil {
				t.Fatal(err)
			}
			// Gyroid has empty bbox; put the box first so Intersect3D's
			// `s0.BoundingBox()` heuristic gives the box's bbox.
			s := sdf.Intersect3D(box, g)
			tris := CollectTriangles(s, NewMarchingCubesOctree(cells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("octree mesh has %d boundary edges (want 0); %d tris", be, len(tris))
			}
			t.Logf("%d tris, %d boundary edges", len(tris), be)
		})
	}
}

// Different clipping shapes produce different boundary geometries that
// stress how the gyroid surface joins the clipping surface. A bbox bug
// in either side would surface here.
func Test_Gyroid3D_VariedClipping_Watertight(t *testing.T) {
	const cells = 80
	sphere, err := sdf.Sphere3D(5)
	if err != nil {
		t.Fatal(err)
	}
	roundedBox, err := sdf.Box3D(v3.Vec{X: 10, Y: 6, Z: 6}, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	cyl, err := sdf.Cylinder3D(8, 4, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name string
		clip sdf.SDF3
	}{
		{"sphere_clip", sphere},
		{"rounded_box_clip", roundedBox},
		{"cylinder_clip", cyl},
		// Difference of two boxes — non-convex clipping region.
		{"shell_clip", func() sdf.SDF3 {
			outer, _ := sdf.Box3D(v3.Vec{X: 10, Y: 10, Z: 10}, 0)
			inner, _ := sdf.Box3D(v3.Vec{X: 6, Y: 6, Z: 6}, 0)
			return sdf.Difference3D(outer, inner)
		}()},
	}
	scale := v3.Vec{X: 2, Y: 2, Z: 2}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g, err := sdf.Gyroid3D(scale)
			if err != nil {
				t.Fatal(err)
			}
			s := sdf.Intersect3D(c.clip, g)
			tris := CollectTriangles(s, NewMarchingCubesOctree(cells))
			be := CountBoundaryEdges(tris)
			if be != 0 {
				t.Errorf("octree mesh has %d boundary edges (want 0); %d tris", be, len(tris))
			}
			t.Logf("%d tris, %d boundary edges", len(tris), be)
		})
	}
}
