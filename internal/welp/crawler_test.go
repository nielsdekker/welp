package welp

import (
	"bytes"
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
)

func Test_searchStrings(t *testing.T) {
	var opt = Options{
		MinTextLength: 4,
		MaxTextLength: 100,
	}
	var tests = []struct {
		name        string
		data        string
		opt         Options
		expected    []string
		expectedMd5 string
	}{
		{"No text", "", opt, []string{}, "d41d8cd98f00b204e9800998ecf8427e"},
		{"Single quotes", `var api='abcdef'; const host='hostname';`, opt, []string{"abcdef", "hostname"}, "a42331bb4301d58720903948512610d3"},
		{"Double quotes", `var api="ghijkl"; const host='hostname';`, opt, []string{"ghijkl", "hostname"}, "cd56f7ccd1a3355ffcc9fcbccc95778e"},
		{"Backticks", "var api=`mnopqr`; const host=`hostname`;", opt, []string{"mnopqr", "hostname"}, "0a67c3186c9e928a68e2ddd5a5eca141"},
		{"Mixed quotes", `var api="stuvwx';`, opt, []string{}, "58a3cf8b3cd86a2c0ec59ef2d42d4963"},
		{"Quote within quote", `var host="host'name'"`, opt, []string{"host'name'", "name"}, "bc3492e9f8955de224e687c484688156"},
		{"Comment with single quote", "// Host's\nvar host='hostname'", opt, []string{"hostname"}, "73271d27c27ded7091a0da1af253f80f"},
		{"Control characters", "var host='host\x00name'", opt, []string{"host\x00name"}, "7db4eb8ca6632e9de352de0c07d14c3d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			asSet := make(map[string]struct{})
			for _, r := range tt.expected {
				asSet[r] = struct{}{}
			}

			foundStrings, md5sum := searchStrings(bytes.NewReader([]byte(tt.data)), tt.opt, 1)
			asserts.KeysEq(t, asSet, foundStrings)
			asserts.Eq(t, md5sum, tt.expectedMd5)
		})
	}
}
