package welp

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/nielsdekker/welp/internal/requests"
	"github.com/nielsdekker/welp/internal/tree"
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
	result.Raw = []byte{}
	for {
		b := make([]byte, 1024)
		red, err := response.Body.Read(b)
		result.Raw = append(result.Raw, b[0:red]...)

		if err == io.EOF {
			break
		} else if err != nil {
			// Not an EOF err but something, report it
			fmt.Printf("err: %v\n", err)
			break
		}
	}

	if err != nil && err != io.EOF {
		return result, err
	}

	result.FoundStrings = tree.StringValues(result.Raw, result.ContentType)

	return result, nil
}

func parseContentType(res *http.Response) string {
	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return "unknown"
	}
	return mediaType
}
