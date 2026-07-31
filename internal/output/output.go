package output

import (
	"fmt"
	"slices"
	"sort"

	"github.com/nielsdekker/welp/internal/module"
	"github.com/nielsdekker/welp/internal/welp"
)

const ESCAPE_RESET = "\033[0m"
const ESCAPE_RED = "\033[31m"
const ESCAPE_GREEN = "\033[32m"
const ESCAPE_BLUE = "\033[34m"
const ESCAPE_BOLD = "\033[1m"

func WriteTTY(w welp.Welp, opt welp.Options) {
	urlResults, tokenResults := w.AllResults()

	sort.Sort(ByPath(urlResults))
	for _, r := range urlResults {
		if slices.Contains(opt.FilterCodes, r.StatusCode) {
			continue
		}
		if slices.Contains(opt.FilterContentType, r.ContentType) {
			continue
		}

		if r.StatusCode <= 100 {
			fmt.Printf("[%s%d%s]", ESCAPE_BLUE, r.StatusCode, ESCAPE_RESET)
		} else if r.StatusCode <= 400 {
			fmt.Printf("[%s%d%s]", ESCAPE_GREEN, r.StatusCode, ESCAPE_RESET)
		} else {
			fmt.Printf("[%s%d%s]", ESCAPE_RED, r.StatusCode, ESCAPE_RESET)
		}

		fmt.Printf(" %s\n", r.URL.Path)
	}

	for _, r := range tokenResults {
		fmt.Printf("[%s%s%s] %s\n", ESCAPE_GREEN, r.TokenType, ESCAPE_RESET, r.Token)
	}
}

type ByPath []module.URLResult

func (b ByPath) Len() int               { return len(b) }
func (b ByPath) Swap(i int, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPath) Less(i int, j int) bool { return b[i].URL.String() < b[j].URL.String() }
