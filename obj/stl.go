package obj

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/hschendel/stl"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
	"io"
	"math"
)

// triModelSdf is the SDF that represents an imported triangle-only mesh.
type triModelSdf struct {
	// Samples that are on the surface to the owning triangle.
	samplesToTriangles *kdtree.KDTree
	// bb is the bounding box.
	bb sdf.Box3
	// Number of samples to consider
	numSamples int
}

func (m *triModelSdf) Evaluate(p sdf.V3) float64 {
	return m.eval(p)
}

func (m *triModelSdf) eval(p sdf.V3) float64 {
	// Find the closest triangle efficiently
	closestPoints := m.samplesToTriangles.KNN(points.NewPoint([]float64{p.X, p.Y, p.Z}, nil), m.numSamples)
	if len(closestPoints) == 0 {
		return 1 // No mesh found: return air
	}
	// Compute maximum triangle area and
	maxTriArea := 1e-12
	for _, closestIthTriangle := range closestPoints {
		triElem := closestIthTriangle.(*points.Point)
		tri := triElem.Data.(*render.Triangle3)
		triArea := tri.V[1].Sub(tri.V[0]).Length() * tri.V[2].Sub(tri.V[0]).Length() / 2
		maxTriArea = math.Max(maxTriArea, triArea)
		//if i == 0 { // remove triangles that have >90º when compared to the closest one.
		//	for j := i + 1; j < len(closestPoints); j++ {
		//		triElem := closestPoints[j].(*points.Point)
		//		tri2 := triElem.Data.(*render.Triangle3)
		//		if tri.Normal().Dot(tri2.Normal()) < 0 { // >90º between planes
		//			closestPoints = append(closestPoints[:j-off], closestPoints[j+1-off:]...)
		//			off++
		//		}
		//	}
		//}
	}
	retAccum := -math.MaxFloat64 // 0.
	retWeight := 0.
	for _, closestIthTriangle := range closestPoints {
		triElem := closestIthTriangle.(*points.Point)
		tri := triElem.Data.(*render.Triangle3)
		triArea := tri.V[1].Sub(tri.V[0]).Length() * tri.V[2].Sub(tri.V[0]).Length() / 2
		// Compute the distance to the triangle (negative for inside the surface)
		centroidToTestPoint := p.Sub(sdf.V3{X: triElem.Coordinates[0], Y: triElem.Coordinates[1], Z: triElem.Coordinates[2]})
		signedDistanceToTriPlane := tri.Normal().Dot(centroidToTestPoint)
		weight := 1 / (1 + centroidToTestPoint.Length()) // (0, 1)
		weight = weight * triArea / maxTriArea           // weaker push of smaller faces (0, 1)
		weight = math.Pow(weight, 2)                     // weaker push of further away faces (0, 1)
		//log.Println("Weight:", weight, "<--", centroidToTestPoint)
		//proj := p.Sub(tri.Normal().MulScalar(signedDistanceToTriPlane))
		//if stlPointInTriangle(proj, tri.V[0], tri.V[1], tri.V[2], 1) {
		//retAccum += signedDistanceToTriPlane * weight
		//} else {
		//	retAccum += 1 // Force AIR
		//}
		//retWeight += weight
		retAccum = math.Max(retAccum, signedDistanceToTriPlane)
		retWeight = 1
	}
	//log.Println("---")
	ret := sigmoidScaled(retAccum / retWeight)
	//log.Println(ret)
	return ret
}

func sigmoidScaled(x float64) float64 {
	return 2/(1+math.Exp(-x)) - 1 // [-1, 1]
}

func (m *triModelSdf) BoundingBox() sdf.Box3 {
	return m.bb
}

// ImportTriMesh converts a triangle-based mesh into a SDF3 surface.
func ImportTriMesh(tris chan *render.Triangle3, resolution float64, numSamples int) sdf.SDF3 {
	m := &triModelSdf{
		samplesToTriangles: kdtree.New(nil),
		bb: sdf.Box3{
			Min: sdf.V3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64},
			Max: sdf.V3{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64},
		},
		numSamples: numSamples,
	}
	for triangle := range tris {
		// Register the triangle's centroid associated to this triangle for later evaluation
		triCentroid := triangle.V[0].Add(triangle.V[1]).Add(triangle.V[2]).DivScalar(3)
		m.samplesToTriangles.Insert(points.NewPoint([]float64{triCentroid.X, triCentroid.Y, triCentroid.Z}, triangle))
		// Update the bounding box
		for _, vertex := range triangle.V {
			//vertex2 := vertex.MulScalar(0.99).Add(triCentroid.MulScalar(0.01))
			//m.samplesToTriangles.Insert(points.NewPoint([]float64{vertex2.X, vertex2.Y, vertex2.Z}, triangle))
			m.bb = m.bb.Include(vertex)
		}
		//// Sample inside the triangle (detail parameter) to avoid problems,
		//// e.g. large triangles that are closer than small triangles that have closer vertices
		//edge1 := triangle.V[1].Sub(triangle.V[0])
		//edge2 := triangle.V[2].Sub(triangle.V[0])
		//edgeLen := sdf.V2{X: edge1.Length(), Y: edge2.Length()}
		//eps := 1e-3
		//uv := sdf.V2{}
		//for uv.X = eps; uv.X < edgeLen.X-eps; uv.X += resolution {
		//	for uv.Y = eps; uv.Y < edgeLen.Y-eps; uv.Y += resolution {
		//		// Create the "uniform" sample point
		//		samplePoint := triangle.V[0].Add(edge1.MulScalar(uv.X)).Add(edge2.MulScalar(uv.Y))
		//		// Ignore this sample if it lies outside the triangle
		//		if !stlPointInTriangle(samplePoint, triangle.V[0], triangle.V[1], triangle.V[2], 0) { // TODO: Efficiency?
		//			continue
		//		}
		//		// Add the "uniform" triangle sample to the tree
		//		m.samplesToTriangles.Insert(points.NewPoint([]float64{samplePoint.X, samplePoint.Y, samplePoint.Z}, triangle))
		//	}
		//}
	}
	if !m.bb.Contains(m.bb.Min) { // Return a valid bounding box if no vertices are found in the mesh
		m.bb = sdf.Box3{} // Empty box centered at {0,0,0}
	}
	m.bb = m.bb.ScaleAboutCenter(1 + 1e-12) // Avoids missing faces due to inaccurate math operations.
	return m
}

func stlSameSide(p1, p2, a, b sdf.V3, tol float64) bool {
	cp1 := b.Sub(a).Cross(p1.Sub(a))
	cp2 := b.Sub(a).Cross(p2.Sub(a))
	return cp1.Dot(cp2) >= -tol
}

func stlPointInTriangle(p, a, b, c sdf.V3, tol float64) bool {
	return stlSameSide(p, a, b, c, tol) && stlSameSide(p, b, a, c, tol) && stlSameSide(p, c, a, b, tol)
}

// ImportSTL converts an STL model into a SDF3 surface.
func ImportSTL(reader io.ReadSeeker, resolution float64, numSamples int) (sdf.SDF3, error) {
	mesh, err := stl.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	triChan := make(chan *render.Triangle3, 64) // Buffer some triangles and send in batches if scheduler prefers it
	go func() {
		for _, triangle := range mesh.Triangles {
			tri := &render.Triangle3{}
			for i, vertex := range triangle.Vertices {
				tri.V[i] = stlToV3(vertex)
			}
			triChan <- tri
		}
		close(triChan)
	}()
	return ImportTriMesh(triChan, resolution, numSamples), nil
}

func stlToV3(v stl.Vec3) sdf.V3 {
	return sdf.V3{X: float64(v[0]), Y: float64(v[1]), Z: float64(v[2])}
}
