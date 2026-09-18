package welp

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/nielsdekker/welp/internal/requests"
)

const MB = int64(1024 * 1024)

type CrawlResult struct {
	Origin       string
	StatusCode   int
	ContentType  string
	FoundStrings map[string]struct{}
	Depth        int
	Raw          []byte
}

func crawl(
	ctx context.Context,
	target string,
	pool requests.Pool,
) (CrawlResult, error) {
	result := CrawlResult{
		Origin:       target,
		FoundStrings: make(map[string]struct{}),
		StatusCode:   0,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return result, fmt.Errorf("Unable to create request: %w", err)
	}

	response, err := pool.Do(ctx, req)
	if err != nil {
		return result, fmt.Errorf("Request failed: %w", err)
	}
	defer response.Body.Close()

	// Overwrite the origin, when redirects occur this contains the value of the
	// URL that answered. Solves issues with directory listing and relative
	// paths on these pages.
	result.Origin = response.Request.URL.String()
	result.StatusCode = response.StatusCode
	result.ContentType = parseContentType(response)

	// Read at most the first 10MB
	toRead := 10 * MB
	if response.ContentLength > 0 {
		toRead = min(response.ContentLength, toRead)
	}

	buf := make([]byte, toRead)
	red, err := response.Body.Read(buf)

	if err != nil && err != io.EOF {
		return result, err
	}

	result.Raw = buf[0:red]
	result.FoundStrings = searchStrings(result.Raw)

	return result, nil
}

// Searches for string like values in the given reader
func searchStrings(raw []byte) map[string]struct{} {
	result := make(map[string]struct{})

	quoteIndices := map[byte]int{
		'\'': -1,
		'"':  -1,
		'`':  -1,
	}

	for rawIndex, b := range raw {
		if quoteIndex, ok := quoteIndices[b]; ok {
			// This is a quote/string character so parse it
			if quoteIndex >= 0 {
				bytes := raw[quoteIndex+1 : rawIndex]
				if utf8.Valid(bytes) {
					foundValue := strings.TrimSpace(string(bytes))
					result[foundValue] = struct{}{}
				}
				quoteIndices[b] = -1
			} else {
				quoteIndices[b] = rawIndex
			}
		}
	}

	return result
}

func parseContentType(res *http.Response) string {
	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return "unknown"
	}
	return mediaType
}
