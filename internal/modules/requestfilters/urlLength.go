package requestfilters

import (
	"net/url"

	"github.com/nielsdekker/welp/internal/modules"
)

type urlLengthFilter struct {
	maxLength int
}

var _ modules.RequestFilterModule = urlLengthFilter{}

func NewUrlLengthFilter(maxLength int) urlLengthFilter {
	return urlLengthFilter{
		maxLength: maxLength,
	}
}

func (m urlLengthFilter) ShouldRequest(url *url.URL) bool {
	return len(url.String()) <= m.maxLength
}
