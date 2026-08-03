package mocks

import (
	"fmt"
	"net/http"

	"github.com/nielsdekker/welp/src/requests"
)

type mock interface {
	// Handle the url, returns whether or not this mock did handle the endpoint
	handle(string) (*http.Response, bool)
}

type mockpool struct {
	mocks []mock
}

func NewMockPool() requests.Pool {
	return mockpool{
		mocks: []mock{
			newSpaMock(),
		},
	}
}

func (m mockpool) GetPoolSize() int { return 1 }
func (m mockpool) Do(req *http.Request) (*http.Response, error) {
	p := req.URL.Path

	for _, v := range m.mocks {
		res, ok := v.handle(p)
		if ok {
			return res, nil
		}
	}

	return nil, fmt.Errorf("%s Not matched to any mock", p)
}
