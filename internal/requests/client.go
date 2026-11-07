package requests

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	HttpClient        http.Client
	AdditionalHeaders map[string]string
}

func NewClient(
	ignoreSSL bool,
	followRedirects bool,
	additionalHeaders map[string]string,
) Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: ignoreSSL,
	}
	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Add user agent if not set
	foundUserAgent := false
	for k, _ := range additionalHeaders {
		if strings.ToLower(k) == "user-agent" {
			foundUserAgent = true
			break
		}
	}
	if !foundUserAgent {
		additionalHeaders["User-Agent"] = "Welp"
	}

	return Client{
		HttpClient:        client,
		AdditionalHeaders: additionalHeaders,
	}
}

func (t *Client) Do(req *http.Request) (*http.Response, error) {
	for k, v := range t.AdditionalHeaders {
		req.Header.Add(k, v)
	}

	return t.HttpClient.Do(req)
}
