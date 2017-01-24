/*
K-Dimensional Trees
*/

package sdf

import (
	"sort"
)

type ByX []V2

func (a ByX) Len() int           { return len(a) }
func (a ByX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByX) Less(i, j int) bool { return a[i].X < a[j].X }

type ByY []V2

func (a ByY) Len() int           { return len(a) }
func (a ByY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByY) Less(i, j int) bool { return a[i].Y < a[j].Y }

type KdTree2 struct {
	root *KdNode2
}

type KdNode2 struct {
	n     V2
	left  *KdNode2
	right *KdNode2
}

func NewKdTree2(points []V2) *KdTree2 {
	t := KdTree2{}

	sort.Sort(ByX(points))
	return &t
}
