package resultfilters

import "github.com/nielsdekker/welp/internal/modules"

type nonError struct{}

var _ modules.ResultFilterModule = nonError{}

func NewNonError() nonError {
	return nonError{}
}

func (m nonError) ShouldCrawl(result modules.Result) bool {
	return result.StatusCode < 400
}
