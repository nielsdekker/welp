package requests

import (
	"crypto/tls"
	"net/http"
	"time"
)

type transport struct {
	base http.RoundTripper
}

func NewClient(ignoreSSL bool, followRedirects bool) http.Client {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: ignoreSSL,
	}

	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &transport{
			base: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: ignoreSSL,
				},
			},
		},
	}

	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "welp")
	return t.base.RoundTrip(req)
}
