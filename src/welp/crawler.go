package welp

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/nielsdekker/welp/src/requests"
)

const MB = 1024 * 1024

// List of content types to skip parsing, mostly binary formats
var skipContentType = []requests.ContentType{
	requests.ContentTypeZip,
	requests.ContentTypeAudio,
	requests.ContentTypeFont,
	requests.ContentTypeIMG,
	requests.ContentTypeVideo,
}

type CrawlResult struct {
	Origin       string
	StatusCode   int
	ContentType  requests.ContentType
	FoundStrings map[string]struct{}
	MD5Sum       string
	depth        int
}

func crawl(
	ctx context.Context,
	target string,
	pool requests.Pool,
	minLength int,
	maxLength int,
) (CrawlResult, error) {
	result := CrawlResult{
		Origin:       target,
		FoundStrings: make(map[string]struct{}),
		StatusCode:   0,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	response, err := pool.Do(ctx, req)

	if err != nil {
		return result, err
	}

	defer response.Body.Close()
	md5sum := md5.New()

	// Overwrite the origin, when redirects occur this contains the value of the
	// URL that answered. Solves issues with directory listing and relative
	// paths on these pages.
	result.Origin = response.Request.URL.String()
	result.StatusCode = response.StatusCode
	result.ContentType = requests.ParseContentType(response.Header.Get("Content-Type"))

	if response.ContentLength > 10*MB {
		// Skip reading, this is to large
		return result, nil
	}

	if slices.Contains(skipContentType, result.ContentType) {
		// Only determine the MD5 hash and skip the binary data
		io.Copy(md5sum, response.Body)
		result.MD5Sum = hex.EncodeToString(md5sum.Sum(nil))
		return result, nil
	}

	// Extract the text values from the response
	buffer := make([]byte, 1024)
	parsedTill := 0
	allBodyBytes := make([]byte, max(response.ContentLength, 10))
	quoteIndices := map[byte]int{
		'\'': -1,
		'"':  -1,
		'`':  -1,
	}

	for {
		bytesRed, err := response.Body.Read(buffer)
		if bytesRed == 0 && err == io.EOF {
			break
		}

		// Update md5 and the all body values
		md5sum.Write(buffer[0:bytesRed])
		allBodyBytes = append(allBodyBytes, buffer[0:bytesRed]...)

		for parsedTill < len(allBodyBytes) {
			b := allBodyBytes[parsedTill]

			if i, ok := quoteIndices[b]; ok {
				// This is a quote/string character so parse it
				if i >= 0 {
					bytes := allBodyBytes[i+1 : parsedTill]
					if utf8.Valid(bytes) && len(bytes) >= minLength && len(bytes) <= maxLength {
						result.FoundStrings[strings.TrimSpace(string(bytes))] = struct{}{}
					}
					quoteIndices[b] = -1
				} else {
					quoteIndices[b] = parsedTill
				}

			} else {
				switch allBodyBytes[parsedTill] {
				case '\n':
					// Reset all indices, strings will not cross newlines most
					// of the times.
					for k, _ := range quoteIndices {
						quoteIndices[k] = -1
					}
				}
			}
			parsedTill++
		}
	}

	result.MD5Sum = hex.EncodeToString(md5sum.Sum(nil))
	return result, nil
}
