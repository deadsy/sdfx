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

func v3ToPoint(v v3.Vec) rtreego.Point {
	return rtreego.Point{v.X, v.Y, v.Z}
}

//-----------------------------------------------------------------------------

type triMeshSdf struct {
	rtree        *rtreego.Rtree
	numNeighbors int
	bb           sdf.Box3
	numTriangles int
}

const stlEpsilon = 1e-1

func (t *triMeshSdf) Evaluate(p v3.Vec) float64 {
	// Quickly skip checking most triangles by only checking the N closest neighbours (AABB based)
	neighbors := t.rtree.NearestNeighbors(t.numNeighbors, v3ToPoint(p))

	dists, signedDistanceResult := t.evaluate(p, neighbors)

	for !sameSign(dists) {
		t.numNeighbors += 5
		neighbors := t.rtree.NearestNeighbors(t.numNeighbors, v3ToPoint(p))

		// Sometimes the sign of the final result is not consistent.
		dists, signedDistanceResult = t.evaluate(p, neighbors)

		if t.numNeighbors >= t.numTriangles {
			// The max possible number is number of triangles.
			break
		}
	}

	// Does the approach of this paper make sense:
	// https://www2.imm.dtu.dk/pubdb/edoc/imm1289.pdf
	// TODO: If so, try to implement it in the future.

	return signedDistanceResult
}

func (t *triMeshSdf) evaluate(p v3.Vec, neighbors []rtreego.Spatial) ([]float64, float64) {
	// To check if all the distances have the same sign.
	dists := make([]float64, 0, t.numNeighbors)

	// Check all triangle distances
	signedDistanceResult := 1.
	closestTriangle := math.MaxFloat64

	for _, neighbor := range neighbors {
		triangle := neighbor.(*sdf.Triangle3)
		testPointToTriangle := p.Sub(triangle[0])
		triNormal := triangle.Normal()
		signedDistanceToTriPlane := triNormal.Dot(testPointToTriangle)
		// Take this triangle as the source of truth if the projection of the point on the triangle is the closest
		distToTri, _ := stlPointToTriangleDistSq(p, triangle)
		if distToTri < closestTriangle {
			closestTriangle = distToTri
			signedDistanceResult = signedDistanceToTriPlane
		}
		dists = append(dists, signedDistanceToTriPlane)
	}

	return dists, signedDistanceResult
}

func sameSign(values []float64) bool {

	positive, negative := false, false

	for _, v := range values {
		if v > 0 {
			positive = true
		} else if v < 0 {
			negative = true
		}

		// If we've seen both positive and negative, return early
		if positive && negative {
			return false
		}
	}

	// All values must have been the same sign
	return true
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
func ImportTriMesh(mesh []*sdf.Triangle3, numNeighbors, minChildren, maxChildren int) sdf.SDF3 {
	if len(mesh) == 0 {
		return nil
	}
	// Compute the bounding box
	bulkLoad := make([]rtreego.Spatial, len(mesh))
	bb := mesh[0].BoundingBox()
	for i, triangle := range mesh {
		bulkLoad[i] = triangle
		bb = bb.Extend(triangle.BoundingBox())
	}
	return &triMeshSdf{
		rtree:        rtreego.NewTree(3, minChildren, maxChildren, bulkLoad...),
		numNeighbors: numNeighbors,
		bb:           bb,
		numTriangles: len(mesh),
	}
}

//-----------------------------------------------------------------------------

func stlPointToTriangleDistSq(p v3.Vec, triangle *sdf.Triangle3) (float64, bool /* falls outside? */) {
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
func stlClosestTrianglePointTo(p v3.Vec, triangle *sdf.Triangle3) (v3.Vec, bool /* falls outside? */) {
	edgeAbDelta := triangle[1].Sub(triangle[0])
	edgeCaDelta := triangle[0].Sub(triangle[2])
	edgeBcDelta := triangle[2].Sub(triangle[1])

	// The closest point may be a vertex
	uab := stlEdgeProject(triangle[0], edgeAbDelta, p)
	uca := stlEdgeProject(triangle[2], edgeCaDelta, p)
	if uca > 1 && uab < 0 {
		return triangle[0], true
	}
	ubc := stlEdgeProject(triangle[1], edgeBcDelta, p)
	if uab > 1 && ubc < 0 {
		return triangle[1], true
	}
	if ubc > 1 && uca < 0 {
		return triangle[2], true
	}

	// The closest point may be on an edge
	triNormal := triangle.Normal()
	planeAbNormal := triNormal.Cross(edgeAbDelta)
	planeBcNormal := triNormal.Cross(edgeBcDelta)
	planeCaNormal := triNormal.Cross(edgeCaDelta)
	if uab >= 0 && uab <= 1 && !stlPlaneIsAbove(triangle[0], planeAbNormal, p) {
		return stlEdgePointAt(triangle[0], edgeAbDelta, uab), true
	}
	if ubc >= 0 && ubc <= 1 && !stlPlaneIsAbove(triangle[1], planeBcNormal, p) {
		return stlEdgePointAt(triangle[1], edgeBcDelta, ubc), true
	}
	if uca >= 0 && uca <= 1 && !stlPlaneIsAbove(triangle[2], planeCaNormal, p) {
		return stlEdgePointAt(triangle[2], edgeCaDelta, uca), true
	}

	// The closest point is in the triangle so project to the plane to find it
	return stlPlaneProject(triangle[0], triNormal, p), false
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

// ImportSTL converts an STL model into a SDF3 surface. See ImportTriMesh.
func ImportSTL(path string, numNeighbors, minChildren, maxChildren int) (sdf.SDF3, error) {
	mesh, err := render.LoadSTL(path)
	if err != nil {
		return nil, err
	}
	return ImportTriMesh(mesh, numNeighbors, minChildren, maxChildren), nil
}

//-----------------------------------------------------------------------------
