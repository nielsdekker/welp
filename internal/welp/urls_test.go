package welp

import (
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
)

func Test_determineUrls(t *testing.T) {
	ori := "http://test.test/"

	toSet := func(values ...string) map[string]struct{} {
		result := make(map[string]struct{})
		for _, s := range values {
			result[s] = struct{}{}
		}
		return result
	}

	var tests = []struct {
		name        string
		origin      string
		foundString string
		prefixes    map[string]struct{}
		expected    map[string]struct{}
	}{
		// Without prefixes
		{"simple", ori, "foo", toSet(), toSet(ori + "foo")},
		{"relative", "http://relative.test/sub/path", "foo", toSet(), toSet("http://relative.test/sub/foo")},
		{"relative ./", "http://relative.test/sub/path", "./foo", toSet(), toSet("http://relative.test/sub/foo")},
		{"relative parent", "http://relative.test/sub/path", "../foo", toSet(), toSet("http://relative.test/foo")},
		{"relative dir", "http://relative.test/sub/path/", "foo", toSet(), toSet("http://relative.test/sub/path/foo")},
		{"relative dir ./", "http://relative.test/sub/path/", "./foo", toSet(), toSet("http://relative.test/sub/path/foo")},
		{"relative dir parent", "http://relative.test/sub/path/", "../foo", toSet(), toSet("http://relative.test/sub/foo")},
		{"root", "http://relative.test/sub/path/", "/foo", toSet(), toSet("http://relative.test/foo")},
		{"root dir", "http://relative.test/sub/path/", "/foo", toSet(), toSet("http://relative.test/foo")},
		{"fqdn insecure", ori, "http://insecure.test", toSet(), toSet("http://insecure.test")},
		{"fqdn secure", ori, "https://secure.test", toSet(), toSet("https://secure.test")},
		{"path with spaces", ori, "spaces are not valid", toSet(), toSet()},
		{"Byte sequences newline", ori, "foo?a=\n", toSet(), toSet()},
		{"Byte sequences tab", ori, "foo?a=\t", toSet(), toSet()},
		{"Byte sequences null byte", ori, "foo?a=\x00", toSet(), toSet()},
		{"Byte sequences yolo", ori, "foo?a=\x7f", toSet(), toSet()},

		// Prefix tests
		{"prefixes simple", ori, "foo", toSet("a/", "b/"), toSet(ori+"foo", ori+"a/foo", ori+"b/foo")},
		{"prefixes simple with starting slash", ori, "foo", toSet("/a/", "/b/"), toSet(ori+"foo", ori+"a/foo", ori+"b/foo")},
		{"prefixes relative", "http://relative.test/sub/path", "foo", toSet("a/"), toSet("http://relative.test/sub/a/foo", "http://relative.test/sub/foo")},
		{"prefixes relative dir", "http://relative.test/sub/path/", "foo", toSet("a/"), toSet("http://relative.test/sub/path/a/foo", "http://relative.test/sub/path/foo")},
		{"prefixes root", "http://relative.test/sub/path", "/foo", toSet("a/"), toSet("http://relative.test/a/foo", "http://relative.test/foo")},
		{"prefixes fqdn", ori, "http://insecure.test", toSet("a/"), toSet("http://insecure.test")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRaw := determineUrls(CrawlResult{
				Origin:       tt.origin,
				FoundStrings: toSet(tt.foundString),
			}, tt.prefixes)

			actual := make(map[string]struct{})
			for _, u := range actualRaw {
				actual[u.String()] = struct{}{}
			}

			asserts.KeysEq(t, tt.expected, actual)
		})
	}
}
