package welp

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/nielsdekker/welp/internal/cli"
	"github.com/nielsdekker/welp/internal/queue"
)

type Welp struct {
	options              cli.Options
	requestPool          http.Client
	resultFilterModules  []ResultFilterModule
	requestFilterModules []RequestFilterModule
	postProcessModules   []PostProcessModule
	outputModules        []OutputModule
}

func New(
	options cli.Options,
	requestPool http.Client,
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

type queueItem struct {
	target string
	depth  int
}

func (w Welp) StartCrawl(ctx context.Context) {
	resultChannel := make(chan CrawlResult)
	defer close(resultChannel)

	cancelCtx, cancel := context.WithCancel(ctx)

	q := queue.New[queueItem]()
	wg := sync.WaitGroup{}

	for range w.options.ConcurrentRequests {
		wg.Go(welpWorker(cancelCtx, q, w.requestPool, resultChannel, w.postProcessModules))
	}

	q.Push(queueItem{
		target: w.options.Target.String(),
		depth:  0,
	})

	counter := 1
	for counter > 0 {
		select {
		case <-ctx.Done():
			return
		case r := <-resultChannel:
			counter--

			// Check the filter modules if we should continue with this value.
			// This is done after writing the result so we still get the result
			// in the output. This only prevents additional calls that could
			// have originated from this result object.
			if slices.ContainsFunc(w.resultFilterModules, func(e ResultFilterModule) bool {
				return !e.ShouldCrawl(r)
			}) {
				continue
			}

			// Write the result to the output
			for _, o := range w.outputModules {
				o.Write(r)
			}

			for _, newUrl := range determineUrls(r, w.options.Prefixes) {
				if slices.ContainsFunc(w.requestFilterModules, func(e RequestFilterModule) bool {
					return !e.ShouldRequest(newUrl)
				}) {
					continue
				}

				counter++
				q.Push(queueItem{
					target: newUrl.String(),
					depth:  0,
				})
			}
		}
	}

	cancel()
	wg.Wait()
}

func welpWorker(
	ctx context.Context,
	q *queue.Queue[queueItem],
	requestPool http.Client,
	resultChannel chan CrawlResult,
	postProcessModules []PostProcessModule,
) func() {
	jitter := time.Duration(rand.Int() % 500)

	return func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// No op
			}

			nextItem, ok := q.Pop()

			if !ok {
				// This can happen if the queue is empty but we are not done
				// yet, so the main thread is processing some items. In that
				// case don't yeet the CPU to 100% checking the queue.
				time.Sleep(jitter * time.Millisecond)
				continue
			}

			result, err := crawl(ctx, nextItem.target, requestPool)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			result.Depth = nextItem.depth + 1

			// Handle any post processing
			for _, p := range postProcessModules {
				p.PostProcess(&result)
			}
			resultChannel <- result
		}
	}
}
