package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nielsdekker/welp/internal/module"
	"github.com/nielsdekker/welp/internal/requests"
	"github.com/nielsdekker/welp/internal/welp"
)

func main() {
	ctx := context.Background()
	target, _ := url.Parse("http://localhost:8000")

	requestPool := requests.NewPool(10)
	w := welp.New(
		target,
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
