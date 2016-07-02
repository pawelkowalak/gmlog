package gmlog

import (
	"testing"
)

func BenchmarkPrintf(b *testing.B) {
	l := New(nil, 5)
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		l.Printf("")
	}
}
