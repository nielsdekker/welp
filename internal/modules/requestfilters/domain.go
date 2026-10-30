package requestfilters

import (
	"net/url"
	"strings"

	"github.com/nielsdekker/welp/internal/modules"
)

type domainFilter struct {
	originalTarget *url.URL
}

var _ modules.RequestFilterModule = domainFilter{}

func NewDomainFilter(originalTarget *url.URL) domainFilter {
	return domainFilter{
		originalTarget: originalTarget,
	}
}

func (m domainFilter) ShouldRequest(url *url.URL) bool {
	if url.Host == m.originalTarget.Host {
		return true
	} else if strings.HasSuffix(url.Host, "."+m.originalTarget.Host) {
		// This is a subdomain so still make the request
		return true
	} else {
		return false
	}
}
