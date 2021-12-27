package dc

import (
	"github.com/deadsy/sdfx/sdf"
	"log"
	"math"
)

func (dc *DualContouringV2) determinant(a, b, c, d, e, f, g, h, i float64) float64 {
	return a*e*i + b*f*g + c*d*h - a*f*h - b*d*i - c*e*g
}

/* dcSolve3x3 Solves for x in  A*x = b. 'A' contains the matrix row-wise. 'b' and 'x' are column vectors. Uses cramer's rule. */
func (dc *DualContouringV2) solve3x3(A []sdf.V3, b []float64) sdf.V3 {
	det := dc.determinant(
		A[0].X, A[0].Y, A[0].Z,
		A[1].X, A[1].Y, A[1].Z,
		A[2].X, A[2].Y, A[2].Z)
	if math.Abs(det) <= 1e-12 {
		if !dc.qefFailedImplWarned {
			log.Println("[DualContouringV1] WARNING: Oh-oh - small determinant:", det)
			dc.qefFailedImplWarned = true
		}
		return sdf.V3{X: math.Inf(1)}
	}
	return sdf.V3{
		X: dc.determinant(
			b[0], A[0].Y, A[0].Z,
			b[1], A[1].Y, A[1].Z,
			b[2], A[2].Y, A[2].Z),
		Y: dc.determinant(
			A[0].X, b[0], A[0].Z,
			A[1].X, b[1], A[1].Z,
			A[2].X, b[2], A[2].Z),
		Z: dc.determinant(
			A[0].X, A[0].Y, b[0],
			A[1].X, A[1].Y, b[1],
			A[2].X, A[2].Y, b[2]),
	}.DivScalar(det)
}

func (dc *DualContouringV2) leastSquares(A []sdf.V3, b []float64) sdf.V3 {
	// assert len(A) == len(b)
	if len(A) == 3 {
		return dc.solve3x3(A, b)
	}
	AtA := [3]sdf.V3{}
	Atb := [3]float64{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sum := 0.
			for k := 0; k < len(A); k++ {
				sum += dcCompGet(A[k], i) * dcCompGet(A[k], j)
			}
			dcCompSet(&AtA[i], j, sum)
		}
	}
	for i := 0; i < 3; i++ {
		sum := 0.
		for k := 0; k < len(A); k++ {
			sum += dcCompGet(A[k], i) * b[k]
		}
		Atb[i] = sum
	}
	return dc.solve3x3(AtA[:], Atb[:])
}
