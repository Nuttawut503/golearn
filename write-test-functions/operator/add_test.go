package operator_test

import (
	"testing"

	"github.com/Nuttawut503/golearn/gotest/operator"
)

func TestAdd(t *testing.T) {
	a, b := 1, 1
	want := 2
	got := operator.Add(a, b)
	if got != want {
		t.Errorf("Run Add(%d, %d), want %d but got %d", a, b, want, got)
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = operator.Add(1_000_000, 1_000_000)
	}
}
