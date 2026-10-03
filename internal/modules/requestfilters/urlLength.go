package requestfilters

import (
	"net/url"

	"github.com/nielsdekker/welp/internal/welp"
)

type urlLengthFilter struct {
	maxLength int
}

var _ welp.RequestFilterModule = urlLengthFilter{}

func NewUrlLengthFilter(maxLength int) urlLengthFilter {
	return urlLengthFilter{
		maxLength: maxLength,
	}
}

func (m urlLengthFilter) ShouldRequest(url *url.URL) bool {
	return len(url.String()) <= m.maxLength
}
