package gorax_test

import (
	"math/rand"
	"testing"

	"github.com/snorwin/gorax"
)

const (
	BenchmarkMaxKeySize = 256
)

func BenchmarkInsert1(b *testing.B) {
	benchmarkInsert(b, 1)
}

func BenchmarkInsert10(b *testing.B) {
	benchmarkInsert(b, 10)
}

func BenchmarkInsert100(b *testing.B) {
	benchmarkInsert(b, 100)
}

func BenchmarkInsert1000(b *testing.B) {
	benchmarkInsert(b, 1000)
}

func BenchmarkInsert10000(b *testing.B) {
	benchmarkInsert(b, 10000)
}

func BenchmarkGet1(b *testing.B) {
	benchmarkGet(b, 1)
}

func BenchmarkGet10(b *testing.B) {
	benchmarkGet(b, 10)
}

func BenchmarkGet100(b *testing.B) {
	benchmarkGet(b, 100)
}

func BenchmarkGet1000(b *testing.B) {
	benchmarkGet(b, 1000)
}

func BenchmarkGet10000(b *testing.B) {
	benchmarkGet(b, 10000)
}

func benchmarkInsert(b *testing.B, size int) {
	keys := make([]string, size)
	for i := 0; i < size; i++ {
		keys[i] = randString(rand.Intn(BenchmarkMaxKeySize))
	}

	b.StopTimer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t := gorax.New()

		b.StartTimer()
		for j := 0; j < size; j++ {
			t.Insert(keys[j], "")
		}
		b.StopTimer()
	}
}

func benchmarkGet(b *testing.B, size int) {
	keys := make([]string, size)
	for i := 0; i < size; i++ {
		keys[i] = randString(rand.Intn(BenchmarkMaxKeySize))
	}

	t := gorax.New()
	for j := 0; j < size; j++ {
		t.Insert(keys[j], "")
	}

	b.StopTimer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < size; j++ {
			t.Get(keys[j])
		}
		b.StopTimer()
	}
}
