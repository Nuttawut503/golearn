package cal_test

import (
	"hello/cal"
	"testing"
)

func TestAdd(t *testing.T) {
	givenA, givenB := 100, 200
	want := 300
	get := cal.Add(givenA, givenB)
	if get != want {
		t.Errorf("given (%d, %d) expected %d, but got %d", givenA, givenB, want, get)
	}
}
