package welp_test

import (
	"net/url"
	"slices"
	"testing"

	"github.com/nielsdekker/welp/src/welp"
)

func eq[K comparable](t *testing.T, expected K, actual K) {
	if expected != actual {
		t.Errorf("Expected \"%v\" to match \"%v\"", expected, actual)
	}
}

func hasPath(t *testing.T, results []welp.CrawlResult, p string) {
	if !slices.ContainsFunc(results, func(res welp.CrawlResult) bool {
		return res.Origin.Path == p
	}) {
		t.Errorf("Expected \"%s\" to be in the crawl results", p)
	}
}

func mustParse(u string) *url.URL {
	res, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return res
}
