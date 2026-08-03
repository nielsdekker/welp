package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nielsdekker/welp/src/output"
	"github.com/nielsdekker/welp/src/post_process"
	"github.com/nielsdekker/welp/src/requests"
	"github.com/nielsdekker/welp/src/welp"
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
		requestPool,
		opt,
	)

	w.StartCrawl(ctx)
	output.WriteTTY(w, []postprocess.PostProcessor{
		// postprocess.NewAllText(),
	}, opt)
}
