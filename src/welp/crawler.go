package welp

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"slices"

	"github.com/nielsdekker/welp/src/requests"
)

// List of content types to skip parsing, mostly binary formats
var skipContentType = []requests.ContentType{
	requests.ContentTypeZip,
	requests.ContentTypeAudio,
	requests.ContentTypeFont,
	requests.ContentTypeIMG,
	requests.ContentTypeVideo,
}

type CrawlResult struct {
	Origin       *url.URL
	StatusCode   int
	ContentType  requests.ContentType
	FoundStrings map[string]struct{}
	MD5Sum       string
}

func (w Welp) crawl(u *url.URL) (CrawlResult, error) {
	result := CrawlResult{
		Origin:       u,
		FoundStrings: make(map[string]struct{}),
	}

	response, err := w.requestPool.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
	})

	if err != nil {
		return result, err
	}

	defer response.Body.Close()
	md5sum := md5.New()
	result.StatusCode = response.StatusCode
	result.ContentType = requests.ParseContentType(response.Header.Get("Content-Type"))

	if response.ContentLength > 1024*1024*10 {
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
		bytesRed, _ := response.Body.Read(buffer)
		if bytesRed == 0 {
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
					if w.isAscii(bytes) {
						result.FoundStrings[string(bytes)] = struct{}{}
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

func (w Welp) isAscii(data []byte) bool {
	if len(data) < w.options.MinTextLength {
		return false
	}

	for _, b := range data {
		// Text values in ASCII
		if b >= 32 && b <= 126 {
			continue
		}

		// Control values like newlines and tabs in ASCII
		if b >= 9 && b <= 13 {
			continue
		}

		return false
	}

	return true
}
