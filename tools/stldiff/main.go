// stldiff: compare two binary STL files.
// Prints IDENTICAL / MINOR / MATERIAL plus metrics.
package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
)

type vec3 struct{ x, y, z float32 }

type tri struct{ a, b, c vec3 }

func loadSTL(path string) ([]tri, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var header [80]byte
	if _, err := io.ReadFull(f, header[:]); err != nil {
		return nil, err
	}
	var count uint32
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		return nil, err
	}
	tris := make([]tri, count)
	for i := uint32(0); i < count; i++ {
		var buf [12]float32
		if err := binary.Read(f, binary.LittleEndian, &buf); err != nil {
			return nil, err
		}
		var attr uint16
		if err := binary.Read(f, binary.LittleEndian, &attr); err != nil {
			return nil, err
		}
		tris[i] = tri{
			a: vec3{buf[3], buf[4], buf[5]},
			b: vec3{buf[6], buf[7], buf[8]},
			c: vec3{buf[9], buf[10], buf[11]},
		}
	}
	return tris, nil
}

func lessVec(a, b vec3) bool {
	if a.x != b.x {
		return a.x < b.x
	}
	if a.y != b.y {
		return a.y < b.y
	}
	return a.z < b.z
}

func sortTriVerts(t tri) tri {
	vs := [3]vec3{t.a, t.b, t.c}
	sort.Slice(vs[:], func(i, j int) bool { return lessVec(vs[i], vs[j]) })
	return tri{vs[0], vs[1], vs[2]}
}

func canonicalHash(tris []tri) string {
	sorted := make([]tri, len(tris))
	for i, t := range tris {
		sorted[i] = sortTriVerts(t)
	}
	sort.Slice(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]
		if a.a != b.a {
			return lessVec(a.a, b.a)
		}
		if a.b != b.b {
			return lessVec(a.b, b.b)
		}
		return lessVec(a.c, b.c)
	})
	h := sha1.New()
	buf := make([]byte, 4)
	for _, t := range sorted {
		for _, v := range [3]vec3{t.a, t.b, t.c} {
			for _, f := range [3]float32{v.x, v.y, v.z} {
				binary.LittleEndian.PutUint32(buf, math.Float32bits(f))
				h.Write(buf)
			}
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:12]
}

type metrics struct {
	count    int
	min, max vec3
	hash     string
}

func computeMetrics(tris []tri) metrics {
	m := metrics{count: len(tris)}
	if len(tris) == 0 {
		return m
	}
	m.min = tris[0].a
	m.max = tris[0].a
	for _, t := range tris {
		for _, v := range [3]vec3{t.a, t.b, t.c} {
			if v.x < m.min.x {
				m.min.x = v.x
			}
			if v.y < m.min.y {
				m.min.y = v.y
			}
			if v.z < m.min.z {
				m.min.z = v.z
			}
			if v.x > m.max.x {
				m.max.x = v.x
			}
			if v.y > m.max.y {
				m.max.y = v.y
			}
			if v.z > m.max.z {
				m.max.z = v.z
			}
		}
	}
	m.hash = canonicalHash(tris)
	return m
}

func bboxDelta(a, b metrics) float64 {
	ax, ay, az := float64(a.max.x-a.min.x), float64(a.max.y-a.min.y), float64(a.max.z-a.min.z)
	bx, by, bz := float64(b.max.x-b.min.x), float64(b.max.y-b.min.y), float64(b.max.z-b.min.z)
	return math.Abs(ax-bx) + math.Abs(ay-by) + math.Abs(az-bz)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: stldiff <a.stl> <b.stl>")
		os.Exit(2)
	}
	a, err := loadSTL(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "A:", err)
		os.Exit(2)
	}
	b, err := loadSTL(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "B:", err)
		os.Exit(2)
	}
	ma, mb := computeMetrics(a), computeMetrics(b)
	if ma.hash == mb.hash {
		fmt.Printf("IDENTICAL  tris=%d hash=%s\n", ma.count, ma.hash)
		return
	}
	bb := bboxDelta(ma, mb)
	triDelta := mb.count - ma.count
	bboxSize := float64(ma.max.x-ma.min.x) + float64(ma.max.y-ma.min.y) + float64(ma.max.z-ma.min.z)
	relBbox := bb / bboxSize
	triRel := float64(abs(triDelta)) / float64(ma.count+1)
	status := "MATERIAL"
	if relBbox < 1e-4 && triRel < 0.01 {
		status = "MINOR   "
	}
	fmt.Printf("%s  tris=%d→%d (Δ%+d, %.2f%%)  bbox-Δ=%.3e (%.2f%%)  hashA=%s hashB=%s\n",
		status, ma.count, mb.count, triDelta, triRel*100, bb, relBbox*100, ma.hash, mb.hash)
}
