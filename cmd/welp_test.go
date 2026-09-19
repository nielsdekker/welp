package main

import (
	"context"
	"net/url"
	"slices"
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
	"github.com/nielsdekker/welp/internal/_tests/mocks"
	"github.com/nielsdekker/welp/internal/cli"
	"github.com/nielsdekker/welp/internal/modules/postprocess"
	"github.com/nielsdekker/welp/internal/modules/requestfilters"
	"github.com/nielsdekker/welp/internal/modules/resultfilters"
	"github.com/nielsdekker/welp/internal/welp"
)

func TestSpa(t *testing.T) {
	results := newWelp("http://spa.test/")

	asserts.Eq(t, len(results), 3)
	resultsContain(t, results, "http://spa.test/")
	resultsContain(t, results, "http://spa.test/other")
	resultsContain(t, results, "http://spa.test/style/default.css")
}

func TestSubdomain(t *testing.T) {
	results := newWelp("http://sub.test/")

	asserts.Eq(t, len(results), 3)
	resultsContain(t, results, "http://sub.test/")
	resultsContain(t, results, "http://sub.sub.test/")
	resultsContain(t, results, "http://sub.sub.test/secret")
}

func newWelp(target string) []welp.CrawlResult {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	out := &mockOutput{allResults: []welp.CrawlResult{}}
	w := welp.New(
		cli.Options{
			Target:        targetURL,
			SearchDepth:   5,
			TextMinLength: 1,
			TextMaxLength: 128,
			Prefixes:      map[string]struct{}{},
		},
		mocks.GetPool(),
		[]welp.ResultFilterModule{resultfilters.NewMD5Filter()},
		[]welp.RequestFilterModule{
			requestfilters.NewDomainFilter(targetURL),
			requestfilters.NewVisitedFilter(),
		},
		[]welp.PostProcessModule{postprocess.NewRemoveDefaults()},
		[]welp.OutputModule{out},
	)

	w.StartCrawl(context.Background())

	return out.allResults
}

func resultsContain(t *testing.T, results []welp.CrawlResult, path string) {
	if !slices.ContainsFunc(results, func(res welp.CrawlResult) bool {
		return res.Origin == path || res.Origin+"/" == path
	}) {
		t.Errorf("Expected \"%s\" to be in the crawl results", path)
	}
}

type mockOutput struct {
	allResults []welp.CrawlResult
}

func (m *mockOutput) Write(result welp.CrawlResult) {
	m.allResults = append(m.allResults, result)
}
