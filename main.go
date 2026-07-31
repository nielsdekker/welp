package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nielsdekker/welp/internal/module"
	"github.com/nielsdekker/welp/internal/output"
	"github.com/nielsdekker/welp/internal/requests"
	"github.com/nielsdekker/welp/internal/welp"
)

func main() {
	ctx := context.Background()
	opt, err := welp.ParseOptions()

	if opt.ShowHelp {
		opt.PrintHelp()
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	requestPool := requests.NewPool(opt.ConcurrentRequests)
	w := welp.New(
		opt.Target,
		requestPool,
		[]module.Module{
			module.NewToken(),
			module.NewURL(requestPool),
		},
	)

	w.Crawl(ctx)
	output.WriteTTY(w, opt)
}
