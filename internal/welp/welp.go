package welp

import (
	"context"
	"fmt"
	"slices"

	"github.com/nielsdekker/welp/internal/cli"
	"github.com/nielsdekker/welp/internal/requests"
)

type Welp struct {
	options              cli.Options
	requestPool          requests.Pool
	resultFilterModules  []ResultFilterModule
	requestFilterModules []RequestFilterModule
	postProcessModules   []PostProcessModule
	outputModules        []OutputModule
}

func New(
	options cli.Options,
	requestPool requests.Pool,
	resultFilterModules []ResultFilterModule,
	requestFilterModules []RequestFilterModule,
	postProcessModules []PostProcessModule,
	outputModules []OutputModule,
) Welp {
	return Welp{
		options:              options,
		requestPool:          requestPool,
		resultFilterModules:  resultFilterModules,
		requestFilterModules: requestFilterModules,
		postProcessModules:   postProcessModules,
		outputModules:        outputModules,
	}
}

func (w Welp) StartCrawl(ctx context.Context) {
	resultChannel := make(chan CrawlResult)
	defer close(resultChannel)

	go w.crawlWrapper(ctx, w.options.Target.String(), 0, resultChannel)

	counter := 1

	for counter > 0 {
		select {
		case <-ctx.Done():
			return
		case r := <-resultChannel:
			counter--

			// Handle any post processing
			for _, p := range w.postProcessModules {
				p.PostProcess(&r)
			}

			// Write the result to the output
			for _, o := range w.outputModules {
				o.Write(r)
			}

			// Check the filter modules if we should continue with this value.
			// This is done after writing the result so we still get the result
			// in the output. This only prevents additional calls that could
			// have originated from this result object.
			if slices.ContainsFunc(w.resultFilterModules, func(e ResultFilterModule) bool {
				return !e.ShouldCrawl(r)
			}) {
				continue
			}

			for _, newUrl := range determineUrls(r, w.options.Prefixes) {
				if slices.ContainsFunc(w.requestFilterModules, func(e RequestFilterModule) bool {
					return !e.ShouldRequest(newUrl)
				}) {
					continue
				}

				counter++
				go w.crawlWrapper(ctx, newUrl.String(), r.Depth, resultChannel)
			}
		}
	}
}

func (w Welp) crawlWrapper(ctx context.Context, target string, currentDepth int, resultChannel chan CrawlResult) {
	res, err := crawl(ctx, target, w.requestPool)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	res.Depth = currentDepth + 1
	resultChannel <- res
}
