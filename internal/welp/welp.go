package welp

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/nielsdekker/welp/internal/module"
	"github.com/nielsdekker/welp/internal/requests"
)

type welp struct {
	Target  *url.URL
	Modules []module.Module

	requestPool requests.Pool
	Results     map[*url.URL][]module.Result
	crawledUrls map[*url.URL]struct{}
}

func New(
	target *url.URL,
	requestPool requests.Pool,
	modules []module.Module,
) welp {
	return welp{
		Target:  target,
		Modules: modules,

		requestPool: requestPool,
		Results:     make(map[*url.URL][]module.Result),
		crawledUrls: make(map[*url.URL]struct{}),
	}
}

func (w welp) Crawl(ctx context.Context) {
	urlChannel := make(chan *url.URL)
	resultChannel := make(chan []module.Result, 100)

	for range w.requestPool.GetPoolSize() {
		go func() {
			for newTargetUrl := range urlChannel {
				stringValues := w.fetch(newTargetUrl)
				results := []module.Result{}

				for _, m := range w.Modules {
					results = append(results, m.Handle(newTargetUrl, stringValues)...)
				}

				resultChannel <- results
			}
		}()
	}

	counter := 1
	urlChannel <- w.Target

	for {
		select {
		case <-ctx.Done():
			// TODO debug dat dit netjes werkt met ctrl+c
			close(urlChannel)
			close(resultChannel)
			return
		case foundResults := <-resultChannel:
			{
				counter--
				for _, r := range foundResults {
					switch v := r.(type) {
					case module.URLResult:
						if _, ok := w.crawledUrls[v.URL]; ok {
							break
						} else {
							w.crawledUrls[v.URL] = struct{}{}
						}

						// TODO Make this configurable, skipping 404 is a good
						// default though
						if v.StatusCode == 404 {
							// 404 break, prevents crawling 404 error page that
							// has a relative path somewhere. Results in a
							// "forever" crawl
							break
						}

						// Not yet crawled the found URL so start crawling
						w.Results[r.FoundIn()] = append(w.Results[r.FoundIn()], r)
						urlChannel <- v.URL
						counter++
					case module.TokenResult:
						w.Results[r.FoundIn()] = append(w.Results[r.FoundIn()], r)
					}
				}

				if counter <= 0 {
					close(urlChannel)
					close(resultChannel)
					return
				}
			}
		}
	}
}

// Calls the given url, then retrieves all string like values from it
func (w welp) fetch(targetURL *url.URL) map[string]struct{} {
	foundStrings := make(map[string]struct{})
	getRes, err := w.requestPool.Do(&http.Request{
		Method: http.MethodGet,
		URL:    targetURL,
	})

	if err != nil {
		// TODO, dit kan beter
		panic(err)
	}
	defer getRes.Body.Close()

	// Depending on the mimetype we consume the body or not, for example there
	// is no need to parse images. Only interested in javascript/json/html and
	// the like.
	contentType := getRes.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/") {
		// Not a text type, so skip it
		return foundStrings
	}

	buffer := make([]byte, 1024)
	body := make([]byte, 0, max(getRes.ContentLength, 10))
	bodyIndex := 0

	quoteIndices := map[byte]int{
		'\'': -1,
		'"':  -1,
		'`':  -1,
	}

	for {
		bytesRed, _ := getRes.Body.Read(buffer)

		if bytesRed == 0 {
			break
		}

		// Add the bytes to the body
		body = append(body, buffer[0:bytesRed]...)

		for bodyIndex < len(body) {
			b := body[bodyIndex]
			if quoteIndex, ok := quoteIndices[b]; ok {
				if quoteIndex >= 0 {
					foundStrings[string(body[quoteIndex+1:bodyIndex])] = struct{}{}
					quoteIndices[b] = -1
				} else {
					quoteIndices[b] = bodyIndex
				}
			} else {
				switch body[bodyIndex] {
				case '\n':
					// Reset all indices, strings will not cross newlines most
					// of the times.

					// TODO, a better solution would be to implement comment
					// detection so we don't need this reset
					for k, _ := range quoteIndices {
						quoteIndices[k] = -1
					}
				}
			}

			bodyIndex++
		}
	}

	return foundStrings
}
