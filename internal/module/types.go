package module

import (
	"net/url"
)

type Module interface {
	Handle(
		origin *url.URL,
		tokens map[string]struct{},
	) []Result
}

type Result interface {
	FoundIn() *url.URL
}

type URLResult struct {
	foundIn    *url.URL
	URL        *url.URL
	StatusCode int
}

type TokenResult struct {
	foundIn *url.URL
	Token   string
}

func (u URLResult) FoundIn() *url.URL {
	return u.foundIn
}

func (t TokenResult) FoundIn() *url.URL {
	return t.foundIn
}
