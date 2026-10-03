package welp

import "net/url"

type ResultFilterModule interface {
	ShouldCrawl(result CrawlResult) bool
}

type RequestFilterModule interface {
	ShouldRequest(url *url.URL) bool
}

type PostProcessModule interface {
	PostProcess(result *CrawlResult)
}

type OutputModule interface {
	// Writes the result to the given output
	Write(result CrawlResult)
}
