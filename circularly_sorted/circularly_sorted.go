// Methods for searching a circularly-sorted array.  Such an array
// looks like a rotated version of a sorted array.  E.g.,
//
//    5 8 13 21 1 1 2 3
//
// is circularly-sorted because if you rotated the array left or right
// by four positions you would have a sorted array.

package circularly_sorted

// Find returns true if n is in the circularly-sorted slice a
// and false otherwise.  This is a recursive implementation.
// It take O(log n) time.
func Find(n int, a []int) bool {
	high := len(a) - 1
	if high < 0 {
		return false
	}
	mid := high / 2
	if n == a[mid] {
		return true
	}

	var aPrime []int
	if n < a[mid] {
		if a[0] <= a[mid] && a[0] > n {
			aPrime = a[mid+1:]
		} else {
			aPrime = a[:mid]
		}
	} else {
		if a[mid] <= a[high] && a[high] < n {
			aPrime = a[:mid]
		} else {
			aPrime = a[mid+1:]
		}
	}
	return Find(n, aPrime)
}

// Rotate shifts the values in slice a by the amount given
// by rotation.  Positive numbers rotate right; negative numbers
// rotate left.
func Rotate(a []int, rotation int) {
	if len(a) == 0 {
		// Not just an optimization; avoids division by 0.
		return
	}
	rotation = rotation % len(a)
	if rotation < 0 {
		rotation += len(a)
	}
	if rotation == 0 {
		// Not just an optimization; avoids out of bounds access.
		return
	}
	yStart := len(a) - rotation
	x, y := 0, yStart
	for x != y {
		a[y], a[x] = a[x], a[y]
		x++
		y++
		if y == len(a) {
			y = yStart
		}
		if x == yStart {
			yStart = y
		}
	}
}
