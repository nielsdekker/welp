package resultfilters

import (
	"strings"

	"github.com/nielsdekker/welp/internal/modules"
)

type contentTypeFilter struct {
	skipContentTypes []string
}

var _ modules.ResultFilterModule = &contentTypeFilter{}

func NewContentTypeFilter(contentTypesToFilter []string) contentTypeFilter {
	return contentTypeFilter{
		skipContentTypes: contentTypesToFilter,
	}
}

func (m contentTypeFilter) ShouldCrawl(result modules.Result) bool {
	for _, ct := range m.skipContentTypes {
		if strings.HasPrefix(result.ContentType, ct) {
			return false
		}
	}

	return true
}
