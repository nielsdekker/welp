package resultfilters

import (
	"strings"

	"github.com/nielsdekker/welp/internal/welp"
)

type contentTypeFilter struct {
	skipContentTypes []string
}

var _ welp.ResultFilterModule = &contentTypeFilter{}

func NewContentTypeFilter() *contentTypeFilter {
	return &contentTypeFilter{
		skipContentTypes: []string{
			"application/zip",
			"audio/",
			"font/",
			"image/",
			"video/",
		},
	}
}

func (m *contentTypeFilter) ShouldCrawl(result welp.CrawlResult) bool {
	for _, ct := range m.skipContentTypes {
		if strings.HasPrefix(result.ContentType, ct) {
			return false
		}
	}

	return true
}
