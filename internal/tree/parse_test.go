package tree_test

import (
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
	"github.com/nielsdekker/welp/internal/tree"
)

func Test_StringValues(t *testing.T) {
	var tests = []struct {
		name     string
		data     string
		expected []string
	}{
		{"No text", "", []string{}},
		{"Single quotes", `var api='abcdef'; const host='hostname';`, []string{"abcdef", "hostname"}},
		{"Double quotes", `var api="ghijkl"; const host='hostname';`, []string{"ghijkl", "hostname"}},
		{"Backticks", "var api=`mnopqr`; const host=`hostname`;", []string{"mnopqr", "hostname"}},
		{"Mixed quotes", `var api="stuvwx';`, []string{}},
		{"Quote within quote", `var host="host'name'"`, []string{"host'name'", "name"}},
		{"Control characters", "var host='host\x00name'", []string{"host\x00name"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			asSet := make(map[string]struct{})
			for _, r := range tt.expected {
				asSet[r] = struct{}{}
			}

			foundStrings := tree.StringValues([]byte(tt.data), "unknown")
			asserts.KeysEq(t, asSet, foundStrings)
		})
	}
}
