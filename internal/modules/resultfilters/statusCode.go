package resultfilters

import (
	"slices"

	"github.com/nielsdekker/welp/internal/modules"
)

type statusCodeFilter struct {
	skipStatusCodes []int
}

var _ modules.ResultFilterModule = statusCodeFilter{}

func NewStatusCodeFilter(codesToFilter []int) statusCodeFilter {
	return statusCodeFilter{
		skipStatusCodes: codesToFilter,
	}
}

func (m statusCodeFilter) ShouldCrawl(result modules.Result) bool {
	return !slices.Contains(m.skipStatusCodes, result.StatusCode)
}
