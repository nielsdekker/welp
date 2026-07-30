package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/nielsdekker/welp/internal/module"
	"github.com/nielsdekker/welp/internal/requests"
	"github.com/nielsdekker/welp/internal/welp"
)

type cliFlags struct {
	target             *url.URL
	concurrentRequests int
}

func main() {
	ctx := context.Background()
	flags := parseFlags()

	requestPool := requests.NewPool(flags.concurrentRequests)
	w := welp.New(
		flags.target,
		requestPool,
		[]module.Module{
			module.NewURL(requestPool),
		},
	)

	w.Crawl(ctx)
	for k, r := range w.Results {
		for _, r := range r {
			switch v := r.(type) {
			case module.TokenResult:
				fmt.Printf("[%s] %s\n", k, v.Token)
			case module.URLResult:
				fmt.Printf("[%s] %s\n", k, v.URL)
			}
		}
	}
}

func parseFlags() cliFlags {
	target := flag.String("u", "", "The target URL")
	concurrentRequests := flag.Int("t", 10, "Number of concurrent requests")

	flag.Parse()

	flags := cliFlags{
		concurrentRequests: *concurrentRequests,
	}

	targetURL, _ := url.Parse(*target)
	flags.target = targetURL

	// Validate the options
	if flags.target.String() == "" {
		os.Stderr.WriteString("-u No target URL specified\n\n")
		flag.Usage()
		os.Exit(2)
	}
	if !strings.HasPrefix(flags.target.Scheme, "http") {
		os.Stderr.WriteString(fmt.Sprintf("-u %s:// is not a valid scheme\n\n", flags.target.Scheme))
		flag.Usage()
		os.Exit(2)
	}
	if *concurrentRequests <= 0 {
		os.Stderr.WriteString(fmt.Sprintf("-t Can not be zero or smaller, is %d\n\n", flags.concurrentRequests))
		flag.Usage()
		os.Exit(2)
	}

	return flags
}
