package output

import (
	"fmt"
	"slices"
	"sort"

	"github.com/nielsdekker/welp/src/post_process"
	"github.com/nielsdekker/welp/src/welp"
)

const ESCAPE_RESET = "\033[0m"
const ESCAPE_RED = "\033[31m"
const ESCAPE_GREEN = "\033[32m"
const ESCAPE_BLUE = "\033[34m"
const ESCAPE_BOLD = "\033[1m"

func WriteTTY(w welp.Welp, postProcessors []postprocess.PostProcessor, opt welp.Options) {
	urls := w.CrawledURLs()
	sort.Sort(ByPath(urls))

	for _, r := range urls {
		if slices.Contains(opt.FilterCodes, r.StatusCode) {
			continue
		}
		if slices.Contains(opt.FilterContentType, r.ContentType) {
			continue
		}

		postprocesResults := make(map[string][]postprocess.PostProcessResult)
		for _, p := range postProcessors {
			postprocesResults[p.GetName()] = p.Handle(r)
		}

		if r.StatusCode < 200 {
			fmt.Printf("[%s%d%s]", ESCAPE_BLUE, r.StatusCode, ESCAPE_RESET)
		} else if r.StatusCode <= 400 {
			fmt.Printf("[%s%d%s]", ESCAPE_GREEN, r.StatusCode, ESCAPE_RESET)
		} else {
			fmt.Printf("[%s%d%s]", ESCAPE_RED, r.StatusCode, ESCAPE_RESET)
		}

		fmt.Printf(" %s - %s\n", r.Origin.Path, r.ContentType)

		for k, v := range postprocesResults {
			for _, r := range v {
				fmt.Printf("\t[%s] %s\n", k, r)
			}
		}
	}
}

type ByPath []welp.CrawlResult

func (b ByPath) Len() int               { return len(b) }
func (b ByPath) Swap(i int, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPath) Less(i int, j int) bool { return b[i].Origin.String() < b[j].Origin.String() }
