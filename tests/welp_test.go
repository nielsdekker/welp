package welp_test

import (
	"context"
	"testing"

	"github.com/nielsdekker/welp/src/modules"
	"github.com/nielsdekker/welp/src/requests"
	"github.com/nielsdekker/welp/src/welp"
	"github.com/nielsdekker/welp/tests/mocks"
)

func TestSpa(t *testing.T) {
	results := newWelp(mocks.NewSpaMock())

	eq(t, len(results), 2)
	hasPath(t, results, "/")
	hasPath(t, results, "/style/default.css")
}

func TestToken(t *testing.T) {
	results := newWelp(mocks.NewTokenInJSMock())

	eq(t, len(results), 4)

	tokenModule := modules.NewToken()
	entropyModule := modules.NewEntropy()

	allTokens := []modules.ModuleResult{}
	allEntropy := []modules.ModuleResult{}
	for _, r := range results {
		allTokens = append(allTokens, tokenModule.Handle(r)...)
		allEntropy = append(allEntropy, entropyModule.Handle(r)...)
	}

	eq(t, len(allTokens), 3)
	hasModuleResult(t, allTokens, "ghp_123abc")
	hasModuleResult(t, allTokens, "sk-123-456")
	hasModuleResult(t, allTokens, "ey123.eyabc.def")

	eq(t, len(allEntropy), 1)
	hasModuleResult(t, allEntropy, "zILEpsOAxrvFnMOOxZTMkVrItcyZw6jPpCHHolXFnsaiy5/OgMSywrjGlMW4zLHNhsqLyLrIsD8kxpY=")
}

func newWelp(mockPool requests.Pool) []welp.CrawlResult {
	w := welp.New(
		mockPool,
		welp.Options{
			Target:         mustParse("http://localhost"),
			MaxSearchDepth: 5,
			MinTextLength:  1,
			MaxTextLength:  128,
			Prefixes:       map[string]struct{}{"": struct{}{}},
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
