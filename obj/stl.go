//-----------------------------------------------------------------------------
/*

Closed-surface triangle meshes (and STL files)

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/dhconnelly/rtreego"
)

//-----------------------------------------------------------------------------

type triMeshSdf struct {
	rtree        *rtreego.Rtree
	numNeighbors int
	bb           sdf.Box3
}

const stlEpsilon = 1e-1

func (t *triMeshSdf) Evaluate(p v3.Vec) float64 {
	// Check all triangle distances
	signedDistanceResult := 1.
	closestTriangle := math.MaxFloat64
	// Quickly skip checking most triangles by only checking the N closest neighbours (AABB based)
	neighbors := t.rtree.NearestNeighbors(t.numNeighbors, stlToPoint(p))
	for _, neighbor := range neighbors {
		triangle := neighbor.(*stlTriangle).Triangle3
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

// ImportTriMesh converts a triangle-based mesh into a SDF3 surface. minChildren and maxChildren are parameters that can
// affect the performance of the internal data structure (3 and 5 are a good default; maxChildren >= minChildren > 0).
//
// WARNING: Setting a low numNeighbors will consider many fewer triangles for each evaluated point, greatly speeding up
// the algorithm. However, if the count of triangles is too low artifacts will appear on the surface (triangle
// continuations). Setting this value to MaxInt is extremely slow but will provide correct results, so choose a value
// that works for your model.
//
// It is recommended to cache (and/or smooth) its values by using sdf.VoxelSdf3.
//
// WARNING: It will only work on non-intersecting closed-surface(s) meshes.
// NOTE: Fix using blender for intersecting surfaces: Edit mode > P > By loose parts > Add boolean modifier to join them
func ImportTriMesh(mesh []*render.Triangle3, numNeighbors, minChildren, maxChildren int) sdf.SDF3 {
	m := &triMeshSdf{
		rtree:        nil,
		numNeighbors: numNeighbors,
		bb: sdf.Box3{
			Min: v3.Vec{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64},
			Max: v3.Vec{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64},
		},
	}

	// Compute the bounding box
	bulkLoad := make([]rtreego.Spatial, 0)
	for _, triangle := range mesh {
		bulkLoad = append(bulkLoad, &stlTriangle{Triangle3: triangle})
		for _, vertex := range triangle.V {
			m.bb = m.bb.Include(vertex)
		}
	}
	if !m.bb.Contains(m.bb.Min) { // Return a valid bounding box if no vertices are found in the mesh
		m.bb = sdf.Box3{} // Empty box centered at {0,0,0}
	}
	//m.bb = m.bb.ScaleAboutCenter(1 + 1e-12) // Avoids missing faces due to inaccurate math operations.
	m.rtree = rtreego.NewTree(3, minChildren, maxChildren, bulkLoad...)

	return m
}

//-----------------------------------------------------------------------------

func stlPointToTriangleDistSq(p v3.Vec, triangle *render.Triangle3) (float64, bool /* falls outside? */) {
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
func stlClosestTrianglePointTo(p v3.Vec, triangle *render.Triangle3) (v3.Vec, bool /* falls outside? */) {
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

func stlEdgeProject(edge1, edgeDelta, p v3.Vec) float64 {
	return p.Sub(edge1).Dot(edgeDelta) / edgeDelta.Length2()
}

func stlEdgePointAt(edge1, edgeDelta v3.Vec, t float64) v3.Vec {
	return edge1.Add(edgeDelta.MulScalar(t))
}

func stlPlaneIsAbove(anyPoint, normal, testPoint v3.Vec) bool {
	return normal.Dot(testPoint.Sub(anyPoint)) > 0
}

func stlPlaneProject(anyPoint, normal, testPoint v3.Vec) v3.Vec {
	v := testPoint.Sub(anyPoint)
	d := normal.Dot(v)
	p := testPoint.Sub(normal.MulScalar(d))
	return p
}

//-----------------------------------------------------------------------------

type stlTriangle struct {
	*render.Triangle3
}

func (s *stlTriangle) Bounds() *rtreego.Rect {
	bounds := sdf.Box3{Min: s.V[0], Max: s.V[0]}
	bounds = bounds.Include(s.V[1])
	bounds = bounds.Include(s.V[2])
	points, err := rtreego.NewRectFromPoints(stlToPoint(bounds.Min), stlToPoint(bounds.Max))
	if err != nil {
		panic(err) // Implementation error
	}
	return points
}

func stlToPoint(v3 v3.Vec) rtreego.Point {
	return rtreego.Point{v3.X, v3.Y, v3.Z}
}

//-----------------------------------------------------------------------------

// ImportSTL converts an STL model into a SDF3 surface. See ImportTriMesh.
func ImportSTL(path string, numNeighbors, minChildren, maxChildren int) (sdf.SDF3, error) {
	mesh, err := render.LoadSTL(path)
	if err != nil {
		return nil, err
	}
	return ImportTriMesh(mesh, numNeighbors, minChildren, maxChildren), nil
}

//-----------------------------------------------------------------------------
