package welp

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// A lot of values in a URL are valid, but whitespace characters are not really
// used.
var validURLRegex = regexp.MustCompile(`^[\w\/\.]\S+$`)

// Determines follow urls based on the result
func (w Welp) determineUrls(result resultWithDepth) []*url.URL {
	// Certain content types have patterns for certain urls. For example
	// `/scripts/main.js` could contain `import bla from "./other.js". In this
	// scenario a result would
	if result.depth > w.options.MaxSearchDepth {
		return []*url.URL{}
	}

	urls := []*url.URL{}
	for s := range result.FoundStrings {
		if _, ok := ignoreList[strings.ToLower(s)]; ok {
			continue
		}

		if !validURLRegex.MatchString(s) {
			continue
		}

		for prefix := range w.options.Prefixes {
			if strings.HasPrefix(s, "http") && len(prefix) > 0 {
				continue
			}

			newURL, err := w.stringToUrl(prefix+s, result.Origin)
			if err != nil {
				continue
			}

			urls = append(urls, newURL)
		}
	}

	return urls
}

func (w Welp) stringToUrl(s string, origin *url.URL) (*url.URL, error) {
	target := ""

	if strings.HasPrefix(s, "/") {
		// Absolute URL
		target = fmt.Sprintf("%s://%s", w.options.Target.Scheme, path.Join(w.options.Target.Host, s))
	} else if strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "http://") {
		// Fully qualified url, so use as is
		target = s
	} else {
		// Either a relative url or nothing
		target = fmt.Sprintf(
			"%s://%s",
			origin.Scheme,
			path.Join(
				origin.Host,
				path.Dir(origin.Path),
				s,
			),
		)
	}

	targetUrl, err := url.Parse(target)

	if err != nil {
		return nil, err
	}

	if len(targetUrl.String()) > 1024 {
		return nil, fmt.Errorf("%s is longer then max length of %d", targetUrl.String(), 1024)
	}

	if targetUrl.Host != w.options.Target.Host {
		return nil, fmt.Errorf("%s is a different host then %s", targetUrl.Host, origin.Host)
	}

	return targetUrl, nil
}
