package circularly_sorted

import "testing"

// Compare returns 0 if a == b, -1 if a < b, and 1 if a > b.
func compare(a, b []int) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		switch {
		case a[i] < b[i]:
			return -1
		case a[i] > b[i]:
			return 1
		}
	}
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	}
	return 0
}

type rotation struct {
	rotation int
	expected []int
}

func TestRotate(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	data := []rotation{
		{-1, []int{2, 3, 4, 5, 1}},
		{0, []int{1, 2, 3, 4, 5}},
		{1, []int{5, 1, 2, 3, 4}},
		{2, []int{4, 5, 1, 2, 3}},
		{3, []int{3, 4, 5, 1, 2}},
		{4, []int{2, 3, 4, 5, 1}},
		{5, []int{1, 2, 3, 4, 5}},
	}
	for i, r := range data {
		a := make([]int, len(original))
		copy(a, original)
		Rotate(a, r.rotation)
		if compare(a, r.expected) != 0 {
			t.Errorf("%d: Rotating by %d was %v, should be %v",
				i, r.rotation, a, r.expected)
		}
	}
}

func assertFound(t *testing.T, i, n int, a []int) {
	if !Find(n, a) {
		t.Errorf("%d: %d not found in %v", i, n, a)
	}
}

func assertNotFound(t *testing.T, i, n int, a []int) {
	if Find(n, a) {
		t.Errorf("%d: %d found in %v", i, n, a)
	}
}

// findLinear searches a for the value n and returns true if
// found and false if not found.  The runtime is O(n).  The result
// should always agree with a correct implementation of Find.
func findLinear(n int, a []int) bool {
	for _, value := range a {
		if value == n {
			return true
		}
	}
	return false
}

func assertSame(t *testing.T, i, n int, a []int) {
	if Find(n, a) != findLinear(n, a) {
		t.Errorf("%d: Find and findLinear did not agree whether "+
			"%d is in %v", i, n, a)
	}
}

func createRotatedValues(data []int) [][]int {
	a := make([][]int, len(data))
	for i := range a {
		a[i] = make([]int, len(data))
		copy(a[i], data)
		Rotate(a[i], i)
	}
	return a
}

func TestFound(t *testing.T) {
	original := []int{1, 1, 2, 3, 5, 8, 13, 21}
	data := createRotatedValues(original)
	for i, toFind := range original {
		for j, d := range data {
			assertFound(t, len(original)*i+j, toFind, d)
			assertSame(t, len(original)*i+j, toFind, d)
		}
	}
}

func TestNotFound(t *testing.T) {
	original := []int{1, 1, 2, 3, 5, 8, 13, 21}
	notFound := []int{-5, 0, 4, 7, 10, 17, 34}
	data := createRotatedValues(original)
	for i, toFind := range notFound {
		for j, d := range data {
			assertNotFound(t, len(notFound)*i+j, toFind, d)
			assertSame(t, len(notFound)*i+j, toFind, d)
		}
	}
}
