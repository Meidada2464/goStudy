package slice_mem_allocation

import "testing"

func BenchmarkName(b *testing.B) {
	slice := make([]string, 64)[0:0]
	for i := 0; i < b.N; i++ {
		slice = append(slice, "a")
	}
}

func BenchmarkName2(b *testing.B) {
	slice := make([]string, 64)
	for i := 0; i < b.N; i++ {
		slice = append(slice, "a")
	}
}
