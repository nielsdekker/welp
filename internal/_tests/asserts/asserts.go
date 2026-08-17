package asserts

import (
	"testing"
)

func Eq[K comparable](t *testing.T, expected K, actual K) {
	if expected != actual {
		t.Errorf("Expected \"%v\" to match \"%v\"", expected, actual)
	}
}

func KeysEq[K comparable](t *testing.T, expected map[K]struct{}, actual map[K]struct{}) {
	Eq(t, len(expected), len(actual))

	for k := range expected {
		if _, ok := actual[k]; !ok {
			t.Errorf("Expected %s to be found in: %v", k, actual)
		}
	}
}
