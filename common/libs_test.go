package common

import (
	"testing"

	"github.com/benm-stm/solace-scalable-k8s-operator/tests"
)

func TestLabels(t *testing.T) {
	got := Labels(&tests.SolaceScalable)
	want := map[string]string{
		"app": "test",
	}
	if got["app"] != want["app"] {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestIsItInSlice(t *testing.T) {
	got := IsItInSlice("a", []string{"a", "b", "c"})
	want := true
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestAsSha256(t *testing.T) {
	got := AsSha256("test")
	want := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestUniqueAndNonZero(t *testing.T) {
	got := UniqueAndNonZero([]int32{1, 2, 3, 2, 0})
	want := []int32{1, 2, 3}
	if len(got) != len(want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("got %v, wanted %v", got, want)
		}
	}
}

func TestContains(t *testing.T) {
	got := Contains([]string{"a", "b", "c", "a"},
		"a",
	)
	want := true
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestNextAvailablePort(t *testing.T) {
	got := NextAvailablePort([]int32{1025, 1026, 1028, 1030},
		1025,
	)
	want := 1027
	if int(got) != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
