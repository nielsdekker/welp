package welp

import (
	"net/url"
	"strings"
)

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

		if !urlSafeCharacters(s) {
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

// Checks if the characters are url safe. Basically everything that is not a
// space or control characters
func urlSafeCharacters(data string) bool {
	// This check should be enough. The crawler only passes valid UTF8 values
	for _, r := range data {
		// 0x7f is DEL in ascii
		if r <= ' ' || r == 0x7f {
			return false
		}
	}

	return true
}
