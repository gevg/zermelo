// zfloat32 sorts a []float32 by copying the data to new uint32 backed buffers, before sorting them
// with zuint32, and copying the sorted floats back. This means this allocates twice the additional
// memory that integer based sorts in zermelo like zuint32 usually do.
//
// However, if memory is available, this is much faster than sort.Sort() for large slices.
package zfloat32

import (
	"github.com/shawnsmithdev/zermelo/zuint32"
	"math"
	"sort"
)

// Calling zfloat64.Sort() on slices smaller than this will result is sorting with sort.Sort() instead.
const MinSize = 256

const radix = 8

// Sorts x using a Radix sort (Small slices are sorted with sort.Sort() instead).
func Sort(x []float32) {
	if len(x) < MinSize {
		sort.Sort(float32Sortable(x))
	} else {
		SortBYOB(x, make([]uint32, len(x)), make([]uint32, len(x)))
	}
}

// Similar to Sort(), but returns a sorted copy of x, leaving x unmodified.
func SortCopy(x []float32) []float32 {
	y := make([]uint32, len(x))
	for idx, val := range x {
		y[idx] = floatFlip(math.Float32bits(val))
	}
	zuint32.Sort(y)

	z := make([]float32, len(x))
	for idx, val := range y {
		z[idx] = math.Float32frombits(floatFlop(val))
	}
	return z
}

// Sorts x using a Radix sort, using supplied buffer space y and z. Panics if
// len(x) does not equal len(y) or len(z). Uses radix sort even on small slices..
func SortBYOB(x []float32, y, z []uint32) {
	for idx, val := range x {
		y[idx] = floatFlip(math.Float32bits(val))
	}
	zuint32.SortBYOB(y, z)
	for idx, val := range y {
		x[idx] = math.Float32frombits(floatFlop(val))
	}
}

// Converts a uint32 that represents a true float to one sorts properly
func floatFlip(x uint32) uint32 {
	if (x & 0x80000000) == 0x80000000 {
		return x ^ 0xFFFFFFFF
	}
	return x ^ 0x80000000
}

// Inverse of floatFlip()
func floatFlop(x uint32) uint32 {
	if (x & 0x80000000) == 0 {
		return x ^ 0xFFFFFFFF
	}
	return x ^ 0x80000000
}

type float32Sortable []float32

func (r float32Sortable) Len() int           { return len(r) }
func (r float32Sortable) Less(i, j int) bool { return r[i] < r[j] }
func (r float32Sortable) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }