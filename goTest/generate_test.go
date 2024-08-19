package goTest

import "testing"

func BenchmarkGenerateWithCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateWithCap(1000)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generate(1000)
	}
}
