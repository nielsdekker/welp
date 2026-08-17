package output

import (
	"fmt"

	"github.com/nielsdekker/welp/internal/modules"
	"github.com/nielsdekker/welp/internal/welp"
)

const ESCAPE_RESET = "\033[0m"
const ESCAPE_RED = "\033[31m"
const ESCAPE_GREEN = "\033[32m"
const ESCAPE_BLUE = "\033[34m"
const ESCAPE_BOLD = "\033[1m"

func WriteTTY(outChannel chan welp.CrawlResult, allModules []modules.Module, opt welp.Options) {
	for r := range outChannel {
		if shouldSkip(r, opt) {
			continue
		}
		moduleResults := applyModules(r, allModules)

		if r.StatusCode < 200 {
			fmt.Printf("[%s%3d%s]", ESCAPE_BLUE, r.StatusCode, ESCAPE_RESET)
		} else if r.StatusCode <= 400 {
			fmt.Printf("[%s%3d%s]", ESCAPE_GREEN, r.StatusCode, ESCAPE_RESET)
		} else {
			fmt.Printf("[%s%3d%s]", ESCAPE_RED, r.StatusCode, ESCAPE_RESET)
		}

		fmt.Printf(" %s - %s\n", r.Origin, r.ContentType)

		for k, v := range moduleResults {
			for _, r := range v {
				fmt.Printf("\t[%s] %s\n", k, r)
			}
		}
	}
}
