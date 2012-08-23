package circularly_sorted

import (
	"testing"
)

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

func TestRotate(t *testing.T) {
	// TODO: Improve the table-driven tests.  Make this a struct
	// that has the rotation amount, also. {rotation: int, expected: []int}
	// TODO: Test -1
	// TODO: Test length+1 == 1
	data := []int{1, 2, 3, 4, 5}
	expected := [][]int{
		{1, 2, 3, 4, 5},
		{2, 3, 4, 5, 1},
		{3, 4, 5, 1, 2},
		{4, 5, 1, 2, 3},
		{5, 1, 2, 3, 4},
	}
	for i := range data {
		a := make([]int, len(data))
		copy(a, data)
		Rotate(a, i)
		if compare(a, expected[i]) != 0 {
			t.Errorf("%d: Rotating by %d was %v, should be %v",
				 i, i, a, expected[i])
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

func assertSame(t *testing.T, i, n int, a []int) {
	if Find(n, a) != findLinear(n, a) {
		t.Errorf("%d: Find and findLinear did not agree whether " +
			 "%d is in %v", i, n, a)
	}
}

func createRotatedValues(data []int) [][]int {
	a := make([][]int, len(data))
	for i := range a {
		// TODO: Finish!
	}
}

func TestFound(t *testing.T) {
	assertFound(t, 0, 5, []int{3, 5, 8, 13, 21, 1, 1, 2})
	assertSame(t, 0, 5, []int{3, 5, 8, 13, 21, 1, 1, 2})
}

func TestNotFound(t *testing.T) {
	assertNotFound(t, 0, 6, []int{3, 5, 8, 13, 21, 1, 1, 2})
	assertSame(t, 0, 6, []int{3, 5, 8, 13, 21, 1, 1, 2})
}
