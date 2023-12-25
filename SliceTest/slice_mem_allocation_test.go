package slice_mem_allocation

import "testing"

func BenchmarkSlice0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make([]int, 0)
		for j := 0; j < 20000; j++ {
			m = append(m, j)
		}
	}
}

func BenchmarkSlice100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make([]int, 0, 100)
		for j := 0; j < 20000; j++ {
			m = append(m, j)
		}
	}
}

func BenchmarkSlice10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make([]int, 0, 10000)
		for j := 0; j < 20000; j++ {
			m = append(m, j)
		}
	}
}

func BenchmarkSlice20000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make([]int, 0, 20000)
		for j := 0; j < 20000; j++ {
			m = append(m, j)
		}
	}
}
