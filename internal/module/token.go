package module

import (
	"encoding/base64"
	"math"
	"net/url"
	"regexp"
)

type TokenType int

var regexpTokens = []struct {
	re   *regexp.Regexp
	name string
}{
	// Specific
	{regexp.MustCompile(`ghp_\w+`), "github"},
	{regexp.MustCompile(`github_pat_\w+`), "github"},
	{regexp.MustCompile(`AKIA[\w]+`), "aws"},
	{regexp.MustCompile(`FwoGZXIvYXdz[\w-+\/]+`), "aws"},
	{regexp.MustCompile(`glpat-[\w-+\/=]+`), "gitlab"},
	{regexp.MustCompile(`(xoxb|xoxp|xapp)-[\w-]+`), "slack"},
	{regexp.MustCompile(`sk-ant-oat01-[\w-]+`), "Claude"},
	{regexp.MustCompile(`sk-[\w]+`), "Open AI"},

	// Generic
	{regexp.MustCompile(`Bearer [\w-\/=]+`), "Standard Bearer token"},
	{regexp.MustCompile(`(e[wy][\w-]+\.){2}[\w-]+`), "JWT"},
	{regexp.MustCompile(`[\w+\/=]+`), "base64"},
}

type TokenModule struct {
}

func NewToken() Module {
	return TokenModule{}
}

func (t TokenModule) Handle(
	origin *url.URL,
	tokens map[string]struct{},
) []Result {
	results := []Result{}

	for k := range tokens {
		hadMatch := false

		for _, v := range regexpTokens {
			if v.re.MatchString(k) {
				switch v.name {
				case "base64":
					// Check if it actually is a base64 value
					if _, err := base64.StdEncoding.DecodeString(k); err != nil {
						continue
					} else if entropy(k) < 3.5 {
						// Probably not really base64 but just text containing
						// base64 valid characters
						continue
					}
				}

				results = append(results, TokenResult{
					foundIn:   origin,
					Token:     k,
					TokenType: v.name,
				})

				hadMatch = true
				break
			}
		}

		if !hadMatch {
			// No matches so do an entropy check
			if entropy(k) > 4.5 {
				results = append(results, TokenResult{
					foundIn:   origin,
					Token:     k,
					TokenType: "High Entropy",
				})
			}
		}
	}

	return results
}

func entropy(data string) float64 {
	counts := make([]int8, 256)
	lenData := float64(len(data))
	for _, r := range data {
		counts[r]++
	}

	e := float64(0)
	for _, c := range counts {
		if c <= 0 {
			continue
		}

		px := float64(c) / lenData
		e += -px * math.Log2(px)
	}

	return e
}
