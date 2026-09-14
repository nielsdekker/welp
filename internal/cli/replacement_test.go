package cli_test

import (
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
	"github.com/nielsdekker/welp/internal/cli"
)

func TestTextReplacement_Apply(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		te       cli.TextReplacement
		input    string
		expected string
	}{
		{"Simple", cli.NewReplacement("abc*ghi", "abcdefghi"), "abcabcghi", "abcdefghi"},
		{"Simple no match", cli.NewReplacement("abc*ghi", "abcdefghi"), "abcabcihg", "abcabcihg"},
		{"Wildcard at end", cli.NewReplacement("abc*", "abcdef"), "abccba", "abcdef"},
		{"JS String interpolation", cli.NewReplacement("${*}", "123"), "hello ${message}", "hello 123"},
		{"No wildcard", cli.NewReplacement("abc", "cba"), "hello abc", "hello cba"},
		{"No wildcard", cli.NewReplacement("abc", "cba"), "hello test", "hello test"},
		{"No wildcard", cli.NewReplacement("abc", "cba"), "ab c", "ab c"},
		{"With newlines", cli.NewReplacement("a*", "z"), "abc\ndef", "z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.te.Apply(tt.input)
			asserts.Eq(t, tt.expected, actual)
		})
	}
}
