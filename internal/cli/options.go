package cli

import (
	"fmt"
	"net/url"
	"strings"
)

type Options struct {
	Target             *url.URL
	ConcurrentRequests int
	ShowHelp           bool
	SearchDepth        int
	SSLIgnore          bool
	Prefixes           map[string]struct{}
	FilterCodes        []int
	FilterContentType  []string
}

func (o *Options) IsValid() []error {
	foundErrors := []error{}
	if o.Target == nil || o.Target.String() == "" {
		foundErrors = append(foundErrors, fmt.Errorf("No target URL specified"))
	}
	if o.Target != nil && !strings.HasPrefix(o.Target.Scheme, "http") {
		foundErrors = append(foundErrors, fmt.Errorf("%s:// is not a valid scheme", o.Target.Scheme))
	}
	if o.ConcurrentRequests <= 0 {
		foundErrors = append(foundErrors, fmt.Errorf("Threads can not be zero or smaller, is %d", o.ConcurrentRequests))
	}
	if o.SearchDepth < 0 {
		foundErrors = append(foundErrors, fmt.Errorf("Search depth can not be smaller then zero, is %d", o.SearchDepth))
	}

	return foundErrors
}

// Applies defaults for any value that is not set
func (o *Options) ApplyUnsetDefaults() {
	if len(o.FilterCodes) == 0 {
		o.FilterCodes = []int{404}
	}
}
