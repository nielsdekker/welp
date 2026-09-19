package resultfilters

import (
	"github.com/nielsdekker/welp/internal/welp"
)

type nonError struct{}

var _ welp.ResultFilterModule = nonError{}

func NewNonError() nonError {
	return nonError{}
}

func (m nonError) ShouldCrawl(result welp.CrawlResult) bool {
	return result.StatusCode < 400
}
