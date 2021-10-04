package kdtree

import (
	"math/rand"
	"testing"
)

type CaseInside struct {
	point, lowPoint, highPoint Point
	numDims                    int
	isInside                   bool
}

func TestIsInside(t *testing.T) {
	cases := []CaseInside{
		{
			Point{[]uint64{30, 80, 40}, 0},
			Point{[]uint64{30, 80, 40}, 0},
			Point{[]uint64{50, 90, 50}, 0},
			3,
			true,
		},
		{
			Point{[]uint64{30, 79, 40}, 0},
			Point{[]uint64{30, 80, 40}, 0},
			Point{[]uint64{50, 90, 50}, 0},
			3,
			false,
		},
		{ //invalid range
			Point{[]uint64{30, 80, 40}, 0},
			Point{[]uint64{30, 80, 40}, 0},
			Point{[]uint64{50, 90, 39}, 0},
			3,
			false,
		},
	}

	for i, tc := range cases {
		res := tc.point.Inside(tc.lowPoint, tc.highPoint)
		if res != tc.isInside {
			t.Fatalf("case %v failed\n", i)
		}
	}
}

func NewRandPoints(numDims int, maxVal uint64, size int) (points []Point) {
	for i := 0; i < size; i++ {
		vals := make([]uint64, 0, numDims)
		for j := 0; j < numDims; j++ {
			vals = append(vals, rand.Uint64()%maxVal)
		}
		point := Point{vals, uint64(i)}
		points = append(points, point)
	}
	return
}

//verify if lhs and rhs contains the same points. order doesn't matter.
func areSmaePoints(lhs, rhs []Point, numDims int) (res bool) {
	if len(lhs) != len(rhs) {
		return
	}
	numPoints := len(lhs)
	mapLhs := make(map[uint64]Point, numPoints)
	mapRhs := make(map[uint64]Point, numPoints)
	for i := 0; i < numPoints; i++ {
		mapLhs[lhs[i].UserData.(uint64)] = lhs[i]
		mapRhs[rhs[i].UserData.(uint64)] = rhs[i]
	}
	for k, v := range mapLhs {
		v2, found := mapRhs[k]
		if !found || !v.Equal(v2) {
			return
		}
	}
	res = true
	return
}

func verifySplit(t *testing.T, pam *PointArrayMem, numStrips int, splitValues []uint64, splitPoses []int) {
	//fmt.Printf("points: %v\nsplitValues: %v\nsplitPoses:%v\n", points, splitValues, splitPoses)
	if len(splitValues) != numStrips-1 || len(splitValues) != len(splitPoses) {
		t.Fatalf("incorrect size of splitValues or splitPoses\n")
	}
	for i := 0; i < len(splitValues)-1; i++ {
		if splitValues[i] > splitValues[i+1] {
			t.Fatalf("incorrect splitValues\n")
		}
		if splitPoses[i] >= splitPoses[i+1] {
			t.Fatalf("incorrect splitPoses\n")
		}
	}
	numSplits := len(splitValues)
	for strip := 0; strip < numStrips; strip++ {
		posBegin := 0
		minValue := uint64(0)
		if strip != 0 {
			posBegin = splitPoses[strip-1]
			minValue = splitValues[strip-1]
		}
		posEnd := len(pam.points)
		maxValue := ^uint64(0)
		if strip != numSplits {
			posEnd = splitPoses[strip]
			maxValue = splitValues[strip]
		}

		for pos := posBegin; pos < posEnd; pos++ {
			val := pam.points[pos].Vals[pam.byDim]
			if val < minValue || val > maxValue {
				t.Fatalf("points[%v] dim %v val %v is not in range %v-%v", pos, pam.byDim, val, minValue, maxValue)
			}
		}
	}
}

func TestSplitPoints(t *testing.T) {
	//TODO: use suite setup to initialize points
	numDims := 3
	maxVal := uint64(100)
	size := 10000
	numStrips := 4
	points := NewRandPoints(numDims, maxVal, size)
	pointsSaved := make([]Point, size)
	copy(pointsSaved, points)
	for dim := 0; dim < numDims; dim++ {
		pam := &PointArrayMem{
			points: points,
			byDim:  dim,
		}
		splitValues, splitPoses := SplitPoints(pam, numStrips)
		verifySplit(t, pam, numStrips, splitValues, splitPoses)
		if !areSmaePoints(pointsSaved, pam.points, numDims) {
			t.Fatalf("point set changes after split")
		}
	}
}
