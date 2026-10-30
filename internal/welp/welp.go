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
	"github.com/nielsdekker/welp/internal/modules"
	"github.com/nielsdekker/welp/internal/queue"
)

type Welp struct {
	options              cli.Options
	requestPool          http.Client
	resultFilterModules  []modules.ResultFilterModule
	requestFilterModules []modules.RequestFilterModule
	postProcessModules   []modules.PostProcessModule
	outputModules        []modules.OutputModule
}

func New(
	options cli.Options,
	requestPool http.Client,
	resultFilterModules []modules.ResultFilterModule,
	requestFilterModules []modules.RequestFilterModule,
	postProcessModules []modules.PostProcessModule,
	outputModules []modules.OutputModule,
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
	resultChannel := make(chan modules.Result)
	defer close(resultChannel)

	cancelCtx, cancel := context.WithCancel(ctx)

	q := queue.New[queueItem]()
	wg := sync.WaitGroup{}

	for range w.options.NumWorkers {
		wg.Go(w.worker(cancelCtx, q, resultChannel))
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

			// Before writing the output check if we should actually use this
			// result.
			if slices.ContainsFunc(w.resultFilterModules, func(e modules.ResultFilterModule) bool {
				return !e.ShouldCrawl(r)
			}) {
				continue
			}

			// Write the result to the output
			for _, o := range w.outputModules {
				o.Write(r)
			}

			// Search depth is a bit in limbo between the filters. We still want
			// the output from the `outputModules` but don't want to continue,
			// therefore it is it's own check.
			if r.Depth >= w.options.SearchDepth {
				continue
			}

			for _, newUrl := range determineUrls(r, w.options.Prefixes) {
				if slices.ContainsFunc(w.requestFilterModules, func(e modules.RequestFilterModule) bool {
					return !e.ShouldRequest(newUrl)
				}) {
					continue
				}

				counter++
				q.Push(queueItem{
					target: newUrl.String(),
					depth:  r.Depth + 1,
				})
			}
		}
	}

	cancel()
	wg.Wait()
}

func (w Welp) worker(
	ctx context.Context,
	q *queue.Queue[queueItem],
	resultChannel chan modules.Result,
) func() {
	jitter := time.Duration((rand.Int() % 200) + 50)

	return func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// No op
			}

			nextItem := q.Pop()
			if nextItem == nil {
				// This can happen if the queue is empty but we are not done
				// yet, so the main thread is processing some items. In that
				// case don't spike the CPU to 100% just checking the queue.
				time.Sleep(jitter * time.Millisecond)
				continue
			}

			result, err := crawl(ctx, nextItem.target, w.requestPool)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			result.Depth = nextItem.depth + 1

			// Handle any post processing in the goroutine, these should be
			// thread safe
			for _, p := range w.postProcessModules {
				p.PostProcess(&result)
			}
			resultChannel <- result
		}
	}
}
