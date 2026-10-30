package modules

import (
	"net/url"
)

type ResultFilterModule interface {
	ShouldCrawl(result Result) bool
}

type RequestFilterModule interface {
	ShouldRequest(url *url.URL) bool
}

type PostProcessModule interface {
	PostProcess(result *Result)
}

type OutputModule interface {
	// Writes the result to the given output
	Write(result Result)
}

type Result struct {
	Origin            string
	StatusCode        int
	ContentType       string
	FoundStrings      map[string]struct{}
	PostProcessOutput []PostProcessOutput
	Depth             int
	Raw               []byte
}

type postProcessOutputImportance = int

const (
	POST_LOW postProcessOutputImportance = iota
	POST_MED
	POST_HIGH
)

type PostProcessOutput struct {
	Name       string
	Value      string
	Importance postProcessOutputImportance
}
