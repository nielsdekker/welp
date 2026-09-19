package output

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nielsdekker/welp/internal/welp"
)

type ttyOutput struct {
	filterCodes       []int
	filterContentType []string
}

var _ welp.OutputModule = ttyOutput{}

const ESCAPE_RESET = "\033[0m"
const ESCAPE_RED = "\033[31m"
const ESCAPE_GREEN = "\033[32m"
const ESCAPE_BLUE = "\033[34m"
const ESCAPE_BOLD = "\033[1m"

func NewTTYOutput(
	filterCodes []int,
	filterContentType []string,
) ttyOutput {
	return ttyOutput{
		filterCodes:       filterCodes,
		filterContentType: filterContentType,
	}
}

func (m ttyOutput) Write(r welp.CrawlResult) {
	if slices.Contains(m.filterCodes, r.StatusCode) {
		return
	}
	if slices.ContainsFunc(m.filterContentType, func(e string) bool {
		return strings.HasPrefix(r.ContentType, e)
	}) {
		return
	}

	if r.StatusCode < 200 {
		fmt.Printf("[%s%3d%s]", ESCAPE_BLUE, r.StatusCode, ESCAPE_RESET)
	} else if r.StatusCode <= 400 {
		fmt.Printf("[%s%3d%s]", ESCAPE_GREEN, r.StatusCode, ESCAPE_RESET)
	} else {
		fmt.Printf("[%s%3d%s]", ESCAPE_RED, r.StatusCode, ESCAPE_RESET)
	}

	fmt.Printf(" %s - %s\n", r.Origin, r.ContentType)
}
