package Base58

import (
	"testing"
)

const (
	vv = "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM"
)

var ii = []byte("00010966776006953d5567439e5e39f86a0d273beed61967f6")

// BenchmarkDecodeBase58 show benchmark test
func BenchmarkDecodeBase58(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecodeBase58(vv)
	}
	b.StopTimer()
}

func BenchmarkEncodeBase58(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeBase58(ii)
	}
	b.StopTimer()
}
