package main

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/nielsdekker/welp/internal/cli"
	"github.com/nielsdekker/welp/internal/modules/output"
	"github.com/nielsdekker/welp/internal/modules/postprocess"
	"github.com/nielsdekker/welp/internal/modules/requestfilters"
	"github.com/nielsdekker/welp/internal/modules/resultfilters"
	"github.com/nielsdekker/welp/internal/requests"
	"github.com/nielsdekker/welp/internal/welp"
)

func main() {
	ctx := context.Background()
	opt, err := cli.Parse()

	if opt.ShowHelp {
		fmt.Println(opt.Usage())
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	banner(opt)

	requestPool := requests.NewPool(opt.ConcurrentRequests, opt.SSLIgnore)
	w := welp.New(
		opt,
		requestPool,
		[]welp.ResultFilterModule{
			resultfilters.NewContentTypeFilter(),
			resultfilters.NewMD5Filter(),
		},
		[]welp.RequestFilterModule{
			requestfilters.NewUrlLengthFilter(200),
			requestfilters.NewDomainFilter(opt.Target),
			requestfilters.NewVisitedFilter(),
		},
		[]welp.PostProcessModule{postprocess.NewRemoveDefaults()},
		[]welp.OutputModule{output.NewTTYOutput(opt.FilterCodes, opt.FilterContentType)},
	)

	w.StartCrawl(ctx)
}

func banner(opt cli.Options) {
	fmt.Println(` _       __________    ____ 
| |     / / ____/ /   / __ \
| | /| / / __/ / /   / /_/ /
| |/ |/ / /___/ /___/ ____/ 
|__/|__/_____/_____/_/`)

	fmt.Println("\nUsing the following options:")
	fmt.Printf("  %-24s%s\n", "Target", opt.Target.String())
	fmt.Printf("  %-24s%d\n", "Concurrent requests", opt.ConcurrentRequests)
	fmt.Printf("  %-24s%d\n", "Max search depth", opt.SearchDepth)

	if len(opt.FilterContentType) > 0 {
		fmt.Printf("  %-24s%s\n", "Filter content type", opt.FilterContentType)
	}
	if len(opt.FilterCodes) > 0 {
		fmt.Printf("  %-24s%d\n", "Filter status codes", opt.FilterCodes)
	}

	if len(opt.Prefixes) > 0 {
		fmt.Printf("  %-24s%s\n", "Additional prefixes", slices.Collect(maps.Keys(opt.Prefixes)))
	}

	if len(opt.Modules) > 0 {
		fmt.Printf("  %-24s%s\n", "Post process modules", slices.Collect(maps.Keys(opt.Modules)))
	}

	fmt.Println()
}
