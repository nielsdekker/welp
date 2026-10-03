package requestfilters

import (
	"net/url"

	"github.com/nielsdekker/welp/internal/welp"
)

type visitedFilter struct {
	cache map[string]struct{}
}

var _ welp.RequestFilterModule = &visitedFilter{}

func NewVisitedFilter() *visitedFilter {
	return &visitedFilter{
		cache: make(map[string]struct{}),
	}
}

func (m *visitedFilter) ShouldRequest(url *url.URL) bool {
	if _, ok := m.cache[url.String()]; ok {
		return false
	} else {
		m.cache[url.String()] = struct{}{}
		return true
	}
}
