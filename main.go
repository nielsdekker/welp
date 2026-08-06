package main

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/nielsdekker/welp/src/modules"
	"github.com/nielsdekker/welp/src/output"
	"github.com/nielsdekker/welp/src/requests"
	"github.com/nielsdekker/welp/src/welp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	opt, err := welp.ParseOptions()

	if opt.ShowHelp {
		opt.PrintHelp()
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	banner(opt)

	requestPool := requests.NewPool(opt.ConcurrentRequests)
	w := welp.New(
		requestPool,
		opt,
	)

	allModules := []modules.Module{}
	for m := range opt.Modules {
		switch m {
		case "text":
			allModules = append(allModules, modules.NewAllText())
		case "token":
			allModules = append(allModules, modules.NewToken())
		case "entropy":
			allModules = append(allModules, modules.NewEntropy())
		}
	}

	out := make(chan welp.CrawlResult)
	go func() {
		w.StartCrawl(ctx, out)
		close(out)
	}()

	if opt.OutputFile != "" {
		if err := output.WriteJSON(out, allModules, opt); err != nil {
			cancel()
		}
	} else {
		output.WriteTTY(out, allModules, opt)
	}
}

func banner(opt welp.Options) {
	fmt.Println(` _       __________    ____ 
| |     / / ____/ /   / __ \
| | /| / / __/ / /   / /_/ /
| |/ |/ / /___/ /___/ ____/ 
|__/|__/_____/_____/_/`)

	fmt.Println("\nUsing the following options:")
	fmt.Printf("  %-24s%s\n", "Target", opt.Target.String())
	fmt.Printf("  %-24s%d\n", "Concurrent requests", opt.ConcurrentRequests)
	fmt.Printf("  %-24s%d\n", "Max search depth", opt.MaxSearchDepth)

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
