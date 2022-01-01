//-----------------------------------------------------------------------------
/*

Closed-surface triangle meshes (and STL files)

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hschendel/stl"
	"io"
	"math"
)

//-----------------------------------------------------------------------------

type triMeshSdf struct {
	tris []*render.Triangle3
	bb   sdf.Box3
}

const stlEpsilon = 1e-12

func (t *triMeshSdf) Evaluate(p sdf.V3) float64 {
	if !t.bb.Contains(p) { // Fast exit
		// Length to surface is at least distance to bounding box
		return t.bb.Include(p).Size().Sub(t.bb.Size()).Length()
	}
	// Check all triangle distances
	signedDistanceResult := 1.
	closestTriangle := math.MaxFloat64
	for _, triangle := range t.tris {
		// TODO: Find a way to quickly skip this triangle (or a way to iterate a subset of triangles)
		testPointToTriangle := p.Sub(triangle.V[0])
		triNormal := triangle.Normal()
		signedDistanceToTriPlane := triNormal.Dot(testPointToTriangle)
		// Take this triangle as the source of truth if the projection of the point on the triangle is the closest
		distToTri, _ := stlPointToTriangleDistSq(p, triangle)
		if distToTri < closestTriangle {
			closestTriangle = distToTri
			signedDistanceResult = signedDistanceToTriPlane
		}
	}
	return signedDistanceResult
}

func (t *triMeshSdf) BoundingBox() sdf.Box3 {
	return t.bb
}

// ImportTriMesh converts a triangle-based mesh into a SDF3 surface.
//
// It is recommended to cache its values at setup time by using sdf.VoxelSdf.
//
// WARNING: It will only work on non-intersecting closed-surface(s) meshes.
// NOTE: Fix using blender for intersecting surfaces: Edit mode > P > By loose parts > Add boolean modifier to join them
func ImportTriMesh(tris []*render.Triangle3) sdf.SDF3 {
	m := &triMeshSdf{
		tris: tris,
		bb: sdf.Box3{
			Min: sdf.V3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64},
			Max: sdf.V3{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64},
		},
	}

	// Compute the bounding box
	for _, triangle := range tris {
		for _, vertex := range triangle.V {
			m.bb = m.bb.Include(vertex)
		}
	}
	if !m.bb.Contains(m.bb.Min) { // Return a valid bounding box if no vertices are found in the mesh
		m.bb = sdf.Box3{} // Empty box centered at {0,0,0}
	}
	//m.bb = m.bb.ScaleAboutCenter(1 + 1e-12) // Avoids missing faces due to inaccurate math operations.

	return m
}

//-----------------------------------------------------------------------------

func stlPointToTriangleDistSq(p sdf.V3, triangle *render.Triangle3) (float64, bool /* falls outside? */) {
	// Compute the closest point
	closest, fallsOutside := stlClosestTrianglePointTo(p, triangle)
	// Compute distance to the closest point
	closestToP := p.Sub(closest)
	distance := closestToP.Length2()
	// Solve influence (distance) ties, by prioritizing triangles with normals more aligned to `closestToP`.
	// This should fix ghost triangle extensions and smooth the field over sharp angles.
	if fallsOutside { // <-- This is an optimization, as others have 0 extra influence in this step
		distance *= 1 + (1-math.Abs(closestToP.Normalize().Dot(triangle.Normal())))*stlEpsilon
	}
	//log.Println(distance, closestToP.Normalize().Dot(triangle.Normal()))
	return distance, fallsOutside
}

// https://stackoverflow.com/a/47505833
func stlClosestTrianglePointTo(p sdf.V3, triangle *render.Triangle3) (sdf.V3, bool /* falls outside? */) {
	edgeAbDelta := triangle.V[1].Sub(triangle.V[0])
	edgeCaDelta := triangle.V[0].Sub(triangle.V[2])
	edgeBcDelta := triangle.V[2].Sub(triangle.V[1])

	// The closest point may be a vertex
	uab := stlEdgeProject(triangle.V[0], edgeAbDelta, p)
	uca := stlEdgeProject(triangle.V[2], edgeCaDelta, p)
	if uca > 1 && uab < 0 {
		return triangle.V[0], true
	}
	ubc := stlEdgeProject(triangle.V[1], edgeBcDelta, p)
	if uab > 1 && ubc < 0 {
		return triangle.V[1], true
	}
	if ubc > 1 && uca < 0 {
		return triangle.V[2], true
	}

	// The closest point may be on an edge
	triNormal := triangle.Normal()
	planeAbNormal := triNormal.Cross(edgeAbDelta)
	planeBcNormal := triNormal.Cross(edgeBcDelta)
	planeCaNormal := triNormal.Cross(edgeCaDelta)
	if uab >= 0 && uab <= 1 && !stlPlaneIsAbove(triangle.V[0], planeAbNormal, p) {
		return stlEdgePointAt(triangle.V[0], edgeAbDelta, uab), true
	}
	if ubc >= 0 && ubc <= 1 && !stlPlaneIsAbove(triangle.V[1], planeBcNormal, p) {
		return stlEdgePointAt(triangle.V[1], edgeBcDelta, ubc), true
	}
	if uca >= 0 && uca <= 1 && !stlPlaneIsAbove(triangle.V[2], planeCaNormal, p) {
		return stlEdgePointAt(triangle.V[2], edgeCaDelta, uca), true
	}

	// The closest point is in the triangle so project to the plane to find it
	return stlPlaneProject(triangle.V[0], triNormal, p), false
}

func stlEdgeProject(edge1, edgeDelta, p sdf.V3) float64 {
	return p.Sub(edge1).Dot(edgeDelta) / edgeDelta.Length2()
}

func stlEdgePointAt(edge1, edgeDelta sdf.V3, t float64) sdf.V3 {
	return edge1.Add(edgeDelta.MulScalar(t))
}

func stlPlaneIsAbove(anyPoint, normal, testPoint sdf.V3) bool {
	return normal.Dot(testPoint.Sub(anyPoint)) > 0
}

func stlPlaneProject(anyPoint, normal, testPoint sdf.V3) sdf.V3 {
	v := testPoint.Sub(anyPoint)
	d := normal.Dot(v)
	p := testPoint.Sub(normal.MulScalar(d))
	return p
}

//-----------------------------------------------------------------------------

// ImportSTL converts an STL model into a SDF3 surface. See ImportTriMesh.
func ImportSTL(reader io.ReadSeeker) (sdf.SDF3, error) {
	mesh, err := stl.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	tris := make([]*render.Triangle3, 0) // Buffer some triangles and send in batches if scheduler prefers it
	for _, triangle := range mesh.Triangles {
		tri := &render.Triangle3{}
		for i, vertex := range triangle.Vertices {
			tri.V[i] = sdf.V3{X: float64(vertex[0]), Y: float64(vertex[1]), Z: float64(vertex[2])}
		}
		tris = append(tris, tri)
	}
	return ImportTriMesh(tris), nil
}
