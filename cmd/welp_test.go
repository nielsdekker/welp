package main

import (
	"context"
	"net/url"
	"slices"
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
	"github.com/nielsdekker/welp/internal/_tests/mocks"
	"github.com/nielsdekker/welp/internal/modules"
	"github.com/nielsdekker/welp/internal/welp"
)

func TestSpa(t *testing.T) {
	results := newWelp("http://spa.test/")

	asserts.Eq(t, len(results), 2)
	resultsContainPath(t, results, "http://spa.test/")
	resultsContainPath(t, results, "http://spa.test/style/default.css")
}

func TestToken(t *testing.T) {
	results := newWelp("http://token.test/")

	asserts.Eq(t, len(results), 4)

	tokenModule := modules.NewToken()
	entropyModule := modules.NewEntropy()

	allTokens := []modules.ModuleResult{}
	allEntropy := []modules.ModuleResult{}
	for _, r := range results {
		allTokens = append(allTokens, tokenModule.Handle(r)...)
		allEntropy = append(allEntropy, entropyModule.Handle(r)...)
	}

	asserts.Eq(t, len(allTokens), 3)
	resultsContainText(t, allTokens, "ghp_123abc")
	resultsContainText(t, allTokens, "sk-123-456")
	resultsContainText(t, allTokens, "ey123.eyabc.def")

	asserts.Eq(t, len(allEntropy), 1)
	resultsContainText(t, allEntropy, "zILEpsOAxrvFnMOOxZTMkVrItcyZw6jPpCHHolXFnsaiy5/OgMSywrjGlMW4zLHNhsqLyLrIsD8kxpY=")
}

func newWelp(target string) []welp.CrawlResult {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	w := welp.New(
		mocks.GetPool(),
		welp.Options{
			Target:         targetURL,
			MaxSearchDepth: 5,
			MinTextLength:  1,
			MaxTextLength:  128,
			Prefixes:       map[string]struct{}{},
		},
	)

	outChannel := make(chan welp.CrawlResult)
	go func() {
		w.StartCrawl(context.Background(), outChannel)
		close(outChannel)
	}()

	results := []welp.CrawlResult{}
	for r := range outChannel {
		results = append(results, r)
	}

	return results
}

func resultsContainPath(t *testing.T, results []welp.CrawlResult, path string) {
	if !slices.ContainsFunc(results, func(res welp.CrawlResult) bool {
		return res.Origin == path || res.Origin+"/" == path
	}) {
		t.Errorf("Expected \"%s\" to be in the crawl results", path)
	}
}

func resultsContainText(t *testing.T, results []modules.ModuleResult, txt string) {
	if !slices.ContainsFunc(results, func(res modules.ModuleResult) bool {
		return res.FoundValue == txt
	}) {
		t.Errorf("Expected \"%s\" to be in the module results", txt)
	}
}
