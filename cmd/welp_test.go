package main

import (
	"context"
	"net/url"
	"slices"
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
	"github.com/nielsdekker/welp/internal/_tests/mocks"
	"github.com/nielsdekker/welp/internal/cli"
	"github.com/nielsdekker/welp/internal/modules"
	"github.com/nielsdekker/welp/internal/modules/postprocess"
	"github.com/nielsdekker/welp/internal/modules/requestfilters"
	"github.com/nielsdekker/welp/internal/modules/resultfilters"
	"github.com/nielsdekker/welp/internal/welp"
)

func TestSpa(t *testing.T) {
	results := newWelp("http://spa.test/")

	asserts.Eq(t, len(results), 2)
	resultsContain(t, results, "http://spa.test/")
	resultsContain(t, results, "http://spa.test/style/default.css")
}

func TestSubdomain(t *testing.T) {
	results := newWelp("http://sub.test/")

	asserts.Eq(t, len(results), 2)
	resultsContain(t, results, "http://sub.test/")
	resultsContain(t, results, "http://sub.sub.test/")
}

func TestTokens(t *testing.T) {
	results := newWelp("http://token.test/")
	asserts.Eq(t, len(results), 3)

	resultsContain(t, results, "http://token.test/")
	resultsContain(t, results, "http://token.test/scripts/main.js")
	resultsContain(t, results, "http://token.test/scripts/tokens.js")

	resultsPostModuleContains(t, results, "http://token.test/scripts/tokens.js", "token github", "ghp_123abc")
	resultsPostModuleContains(t, results, "http://token.test/scripts/tokens.js", "token openAI", "sk-123-456")
	resultsPostModuleContains(t, results, "http://token.test/scripts/tokens.js", "token jwt", "ey123.eyabc.def")
}

func newWelp(target string) []modules.Result {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	out := &mockOutput{allResults: []modules.Result{}}
	w := welp.New(
		cli.Options{
			Target:      targetURL,
			SearchDepth: 5,
			Prefixes:    map[string]struct{}{},
			NumWorkers:  3,
		},
		mocks.HttpClient,
		[]modules.ResultFilterModule{
			resultfilters.NewMD5Filter(),
			resultfilters.NewStatusCodeFilter([]int{404}),
			resultfilters.NewContentTypeFilter([]string{}),
		},
		[]modules.RequestFilterModule{
			requestfilters.NewDomainFilter(targetURL),
			requestfilters.NewVisitedFilter(),
		},
		[]modules.PostProcessModule{
			postprocess.NewRemoveDefaults(),
			postprocess.NewFindTokens(),
		},
		[]modules.OutputModule{out},
	)

	w.StartCrawl(context.Background())

	return out.allResults
}

func resultsContain(t *testing.T, results []modules.Result, path string) {
	if !slices.ContainsFunc(results, func(res modules.Result) bool {
		return res.Origin == path || res.Origin+"/" == path
	}) {
		t.Errorf("Expected \"%s\" to be in the crawl results", path)
	}
}

func resultsPostModuleContains(t *testing.T, results []modules.Result, path string, ppName string, ppValue string) {
	for _, res := range results {
		if res.Origin == path || res.Origin+"/" == path {
			if !slices.ContainsFunc(res.PostProcessOutput, func(pp modules.PostProcessOutput) bool {
				return pp.Name == ppName && pp.Value == ppValue
			}) {
				t.Errorf("Expected to find post proces result \"%s:%s\" for path \"%s\"", ppName, ppValue, path)
			}
			return
		}
	}
	t.Errorf("Expected \"%s\" to be in the crawl results", path)
}

type mockOutput struct {
	allResults []modules.Result
}

func (m *mockOutput) Write(result modules.Result) {
	m.allResults = append(m.allResults, result)
}
