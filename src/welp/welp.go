package welp

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/nielsdekker/welp/src/requests"
)

// A lot of values are allowed in the path of an URL but the following are
// always invalid.
var validURLRegex = regexp.MustCompile(`[ <>\?#]`)

type Welp struct {
	options     Options
	requestPool requests.Pool

	// Set containing all the URLS already crawled
	crawledUrls map[string]CrawlResult

	// Cache of the responses in MD5, prevents from sending the same result over
	// and over.
	md5Crawled map[string]struct{}
}

func New(
	requestPool requests.Pool,
	opt Options,
) Welp {
	return Welp{
		options: opt,

		requestPool: requestPool,
		crawledUrls: make(map[string]CrawlResult),
		md5Crawled:  make(map[string]struct{}),
	}
}

func (w Welp) StartCrawl(ctx context.Context) {
	resultChannel := make(chan CrawlResult)

	crawl := func(newURL *url.URL) {
		res, err := w.crawl(newURL)
		if err != nil {
			// Only log the error and don't stop, we still need to send a result
			// so the counter is updated
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		resultChannel <- res

	}
	go crawl(w.options.Target)

	counter := 1
	for r := range resultChannel {
		counter--
		w.crawledUrls[r.Origin.String()] = r

		// Skip error results
		if r.StatusCode <= 0 {
			continue
		}

		if _, md5Match := w.md5Crawled[r.MD5Sum]; !md5Match && r.StatusCode < 400 {
			for s := range r.FoundStrings {
				if _, ok := ignoreList[s]; ok {
					continue
				}

				if validURLRegex.MatchString(s) {
					continue
				}

				newURL, err := w.stringToUrl(s, r.Origin)
				if err != nil {
					continue
				}

				if _, ok := w.crawledUrls[newURL.String()]; !ok {
					counter++
					go crawl(newURL)
				}
			}
		}

		w.md5Crawled[r.MD5Sum] = struct{}{}
		if counter <= 0 {
			close(resultChannel)
		}
	}
}

func (w Welp) CrawledURLs() []CrawlResult {
	results := make([]CrawlResult, 0, len(w.crawledUrls))

	for _, v := range w.crawledUrls {
		results = append(results, v)
	}

	return results
}

func (w Welp) stringToUrl(s string, origin *url.URL) (*url.URL, error) {
	target := ""
	if strings.HasPrefix("/", s) {
		// Absolute URL, so append it to the host
		target = fmt.Sprintf("%s://%s", origin.Scheme, path.Join(origin.Host, s))
	} else if strings.HasPrefix("https://", s) || strings.HasPrefix("http://", s) {
		// Fully qualified url, so use as is
		target = s
	} else {
		// Relative URL, append it to the path
		target = fmt.Sprintf("%s://%s", origin.Scheme, path.Join(origin.Host, origin.Path, s))
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
