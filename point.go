package bkdtree

import (
	"sort"

	datastructures "github.com/deepfabric/go-datastructures"
	"github.com/keegancsmith/nth"
)

type Point struct {
	Vals     []uint64
	UserData interface{}
}

type PointArray interface {
	sort.Interface
	GetPoint(idx int) Point
	GetValue(idx int) uint64
	SubArray(begin, end int) PointArray
	Erase(point Point) bool
	Append(point Point)
}

type PointArrayMem struct {
	points []Point
	byDim  int
}

type PointArrayExt struct {
	data        []byte
	numPoints   int
	byDim       int
	bytesPerDim int
	numDims     int
	pointSize   int
}

// Compare is part of datastructures.Comparable interface
func (p Point) Compare(other datastructures.Comparable) int {
	rhs := other.(Point)
	for dim := 0; dim < len(p.Vals); dim++ {
		if p.Vals[dim] != rhs.Vals[dim] {
			return int(p.Vals[dim] - rhs.Vals[dim])
		}
	}

	if pc, ok := p.UserData.(datastructures.Comparable); ok {
		if rhsc, ok := rhs.UserData.(datastructures.Comparable); ok {
			return pc.Compare(rhsc)
		}
	}

	return 0
}

func (p *Point) Inside(lowPoint, highPoint Point) (isInside bool) {
	for dim := 0; dim < len(p.Vals); dim++ {
		if p.Vals[dim] < lowPoint.Vals[dim] || p.Vals[dim] > highPoint.Vals[dim] {
			return
		}
	}
	isInside = true
	return
}

func (p Point) LessThan(rhs Point) (res bool) {
	for dim := 0; dim < len(p.Vals); dim++ {
		if p.Vals[dim] != rhs.Vals[dim] {
			return p.Vals[dim] < rhs.Vals[dim]
		}
	}

	if pc, ok := p.UserData.(datastructures.Comparable); ok {
		if rhsc, ok := rhs.UserData.(datastructures.Comparable); ok {
			return pc.Compare(rhsc) < 0
		}
	}

	return false
}

func (p *Point) Equal(rhs Point) (res bool) {
	if p.UserData != rhs.UserData || len(p.Vals) != len(rhs.Vals) {
		return
	}
	for dim := 0; dim < len(p.Vals); dim++ {
		if p.Vals[dim] != rhs.Vals[dim] {
			return
		}
	}
	res = true
	return
}

// Len is part of sort.Interface.
func (s *PointArrayMem) Len() int {
	return len(s.points)
}

// Swap is part of sort.Interface.
func (s *PointArrayMem) Swap(i, j int) {
	s.points[i], s.points[j] = s.points[j], s.points[i]
}

// Less is part of sort.Interface.
func (s *PointArrayMem) Less(i, j int) bool {
	return s.points[i].Vals[s.byDim] < s.points[j].Vals[s.byDim]
}

func (s *PointArrayMem) GetPoint(idx int) (point Point) {
	point = s.points[idx]
	return
}

func (s *PointArrayMem) GetValue(idx int) (val uint64) {
	val = s.points[idx].Vals[s.byDim]
	return
}

func (s *PointArrayMem) SubArray(begin, end int) (sub PointArray) {
	sub = &PointArrayMem{
		points: s.points[begin:end],
		byDim:  s.byDim,
	}
	return
}

func (s *PointArrayMem) Erase(point Point) (found bool) {
	idx := 0
	for i, point2 := range s.points {
		//assumes each point's userData is unique
		if point.Equal(point2) {
			idx = i
			found = true
			break
		}
	}
	if found {
		s.points = append(s.points[:idx], s.points[idx+1:]...)
	}
	return
}

func (s *PointArrayMem) Append(point Point) {
	s.points = append(s.points, point)
}

// SplitPoints splits points per byDim
func SplitPoints(points PointArray, numStrips int) (splitValues []uint64, splitPoses []int) {
	if numStrips <= 1 {
		return
	}
	splitPos := points.Len() / 2
	nth.Element(points, splitPos)
	splitValue := points.GetValue(splitPos)

	numStrips1 := (numStrips + 1) / 2
	numStrips2 := numStrips - numStrips1
	splitValues1, splitPoses1 := SplitPoints(points.SubArray(0, splitPos), numStrips1)
	splitValues = append(splitValues, splitValues1...)
	splitPoses = append(splitPoses, splitPoses1...)
	splitValues = append(splitValues, splitValue)
	splitPoses = append(splitPoses, splitPos)
	splitValues2, splitPoses2 := SplitPoints(points.SubArray(splitPos, points.Len()), numStrips2)
	splitValues = append(splitValues, splitValues2...)
	for i := 0; i < len(splitPoses2); i++ {
		splitPoses = append(splitPoses, splitPos+splitPoses2[i])
	}
	return
}
