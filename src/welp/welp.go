package welp

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nielsdekker/welp/src/requests"
)

type Welp struct {
	options     Options
	requestPool requests.Pool

	// Set containing all the URLS already crawled
	crawledUrls map[string]struct{}
	md5Cache    map[string]struct{}
}

type resultWithDepth struct {
	CrawlResult
	depth int
}

func New(
	requestPool requests.Pool,
	opt Options,
) Welp {
	return Welp{
		options: opt,

		requestPool: requestPool,
		crawledUrls: map[string]struct{}{},
		md5Cache:    map[string]struct{}{},
	}
}

func (w Welp) StartCrawl(ctx context.Context, outputChannel chan CrawlResult) {
	resultChannel := make(chan resultWithDepth)

	crawl := func(newURL *url.URL, currentDepth int) {
		res, err := w.crawl(newURL)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		resultChannel <- resultWithDepth{
			depth:       currentDepth + 1,
			CrawlResult: res,
		}

	}
	go crawl(w.options.Target, 0)

	counter := 1
	numErrors := 0
	for r := range resultChannel {
		counter--
		_, md5match := w.md5Cache[r.MD5Sum]
		isErrorResponse := r.StatusCode <= 0 || r.StatusCode >= 400

		if !md5match {
			// New result so store it
			w.crawledUrls[r.Origin.String()] = struct{}{}
			w.md5Cache[r.MD5Sum] = struct{}{}

			// And send it on the output channel
			outputChannel <- r.CrawlResult
		}

		// This means the request failed completely
		if r.StatusCode <= 0 {
			numErrors++
		}

		if !md5match && !isErrorResponse {
			for _, newURL := range w.determineUrls(r) {
				if _, ok := w.crawledUrls[newURL.String()]; !ok {
					counter++
					go crawl(newURL, r.depth)
				}
			}
		}

		if counter <= 0 {
			close(resultChannel)
		}
	}
}
