package module

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/nielsdekker/welp/internal/requests"
)

type URLModule struct {
	requestPool       requests.Pool
	maxUrlLength      int
	skipDifferentHost bool
}

// URL Module will check if the found values are a valid URL. Values will be
// skipped when:
// - invalid URL values
// - path to long
func NewURL(
	requestPool requests.Pool,
) Module {
	return URLModule{
		requestPool:       requestPool,
		maxUrlLength:      100,
		skipDifferentHost: true,
	}
}

func (u URLModule) Handle(
	origin *url.URL,
	tokens map[string]struct{},
) []Result {
	results := []Result{}
	for t := range tokens {
		result, err := u.handleToken(origin, t)
		if err != nil {
			continue
		}

		results = append(results, result)
	}

	return results
}

func (u URLModule) handleToken(origin *url.URL, token string) (Result, error) {
	// Check what kind of token it is
	target := ""
	if strings.HasPrefix("/", token) {
		// Absolute URL, so append it to the host
		target = fmt.Sprintf("%s://%s", origin.Scheme, path.Join(origin.Host, token))
	} else if strings.HasPrefix("https://", token) || strings.HasPrefix("http://", token) {
		// Fully qualified url, so use as is
		target = token
	} else {
		// Relative URL, append it to the path
		target = fmt.Sprintf("%s://%s", origin.Scheme, path.Join(origin.Host, origin.Path, token))
	}

	targetUrl, err := url.Parse(target)

	if err != nil {
		return nil, err
	}

	if len(targetUrl.String()) > u.maxUrlLength {
		return nil, fmt.Errorf("%s is longer then max length of %d", targetUrl.String(), u.maxUrlLength)
	}

	if targetUrl.Host != origin.Host {
		return nil, fmt.Errorf("%s is a different host then %s", targetUrl.Host, origin.Host)
	}

	res, err := u.requestPool.Do(&http.Request{
		URL:    targetUrl,
		Method: http.MethodHead,
	})

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	ct := res.Header.Get("Content-Type")

	return URLResult{
		foundIn:        origin,
		URL:            targetUrl,
		StatusCode:     res.StatusCode,
		ContentType:    requests.ParseContentType(ct),
		RawContentType: ct,
	}, nil
}
