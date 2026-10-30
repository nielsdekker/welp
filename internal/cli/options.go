package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nielsdekker/welp/internal/modules"
)

type Options struct {
	Target            *url.URL
	NumWorkers        int
	ShowHelp          bool
	SearchDepth       int
	SSLIgnore         bool
	FollowRedirect    bool
	Prefixes          map[string]struct{}
	FilterCodes       []int
	FilterContentType []string
	Scripts           []modules.PostProcessModule
}

func (o *Options) IsValid() []error {
	foundErrors := []error{}
	if o.Target == nil || o.Target.String() == "" {
		foundErrors = append(foundErrors, fmt.Errorf("No target URL specified"))
	}
	if o.Target != nil && !strings.HasPrefix(o.Target.Scheme, "http") {
		foundErrors = append(foundErrors, fmt.Errorf("%s:// is not a valid scheme", o.Target.Scheme))
	}
	if o.NumWorkers <= 0 {
		foundErrors = append(foundErrors, fmt.Errorf("Threads can not be zero or smaller, is %d", o.NumWorkers))
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

	if len(o.FilterContentType) == 0 {
		o.FilterContentType = []string{
			"application/zip",
			"audio/",
			"font/",
			"image/",
			"video/",
		}
	}
}
