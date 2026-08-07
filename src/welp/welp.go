package welp

import (
	"context"
	"fmt"

	"github.com/nielsdekker/welp/src/requests"
)

type Welp struct {
	options     Options
	requestPool requests.Pool
}

func New(
	requestPool requests.Pool,
	opt Options,
) Welp {
	return Welp{
		options:     opt,
		requestPool: requestPool,
	}
}

func (w Welp) StartCrawl(ctx context.Context, outputChannel chan CrawlResult) {
	resultChannel := make(chan CrawlResult)
	defer close(resultChannel)

	go w.crawlWrapper(ctx, w.options.Target.String(), 0, resultChannel)

	counter := 1
	crawledUrls := make(map[string]struct{})
	md5Cache := make(map[string]struct{})

	for {
		select {
		case <-ctx.Done():
			break
		case r := <-resultChannel:
			counter--
			_, md5match := md5Cache[r.MD5Sum]
			isErrorResponse := r.StatusCode <= 0 || r.StatusCode >= 400
			reachedMaxDepth := r.depth > w.options.MaxSearchDepth

			if !md5match {
				// New result so store it
				crawledUrls[r.Origin] = struct{}{}
				md5Cache[r.MD5Sum] = struct{}{}

				// And send it on the output channel
				outputChannel <- r
			}

			if !md5match && !isErrorResponse && !reachedMaxDepth {
				for _, newURL := range determineUrls(r, w.options.Prefixes) {
					if newURL.Host != w.options.Target.Host {
						// Skip going outside the target domain
						continue
					}

					if _, ok := crawledUrls[newURL.String()]; !ok {
						counter++
						go w.crawlWrapper(ctx, newURL.String(), r.depth, resultChannel)
					}
				}
			}

			if counter <= 0 {
				return
			}
		}
	}
}

func (w Welp) crawlWrapper(ctx context.Context, target string, currentDepth int, resultChannel chan CrawlResult) {
	res, err := crawl(ctx, target, w.requestPool, w.options.MinTextLength, w.options.MaxTextLength)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	res.depth = currentDepth + 1
	resultChannel <- res
}
