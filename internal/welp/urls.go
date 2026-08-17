package welp

import (
	"net/url"
	"regexp"
	"strings"
)

// A lot of values in a URL are valid, but whitespace characters are not really
// used.
var validURLRegex = regexp.MustCompile(`^[\w\/\.]\S+$`)

// Determines follow urls based on the result
func determineUrls(
	result CrawlResult,
	prefixes map[string]struct{},
) []*url.URL {
	urls := []*url.URL{}
	for s := range result.FoundStrings {
		if _, ok := ignoreList[strings.ToLower(s)]; ok {
			continue
		}

		if !validURLRegex.MatchString(s) {
			continue
		}

		if u, err := addStringToUrl(result.Origin, s, ""); err == nil {
			urls = append(urls, u)
		}

		// And add any prefixes
		for prefix := range prefixes {
			if strings.HasPrefix(s, "http") {
				continue
			}

			if u, err := addStringToUrl(result.Origin, s, prefix); err == nil {
				urls = append(urls, u)
			}
		}
	}

	return urls
}

func addStringToUrl(
	origin string,
	toAdd string,
	prefix string,
) (*url.URL, error) {
	// Complete url, use as-is
	if strings.HasPrefix(toAdd, "https://") || strings.HasPrefix(toAdd, "http://") {
		u, err := url.Parse(toAdd)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	asUrl, err := url.Parse(origin)
	if err != nil {
		return nil, err
	}

	pathAndQuery := strings.Split(toAdd, "?")

	if strings.HasPrefix(toAdd, "/") {
		asUrl.Path = ""
		asUrl = asUrl.JoinPath(prefix + pathAndQuery[0])
	} else if strings.HasSuffix(asUrl.Path, "/") {
		asUrl = asUrl.JoinPath(prefix + pathAndQuery[0])
	} else {
		asUrl = asUrl.JoinPath("../", prefix+pathAndQuery[0])
	}

	if len(pathAndQuery) > 1 {
		asUrl.RawQuery = pathAndQuery[1]
	}

	return asUrl, nil
}
