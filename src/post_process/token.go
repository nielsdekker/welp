package postprocess

import (
	"math"
	"regexp"

	"github.com/nielsdekker/welp/src/welp"
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
	{regexp.MustCompile(`sk-ant-oat01-[\w-]+`), "claude"},
	{regexp.MustCompile(`sk-[\w]+`), "openAI"},

	// Generic
	{regexp.MustCompile(`Bearer [\w-\/=]+`), "bearer token"},
	{regexp.MustCompile(`(e[wy][\w-]+\.){2}[\w-]+`), "jwt"},
}

type tokenProcessor struct {
	Name string
}

func NewToken() PostProcessor {
	return tokenProcessor{
		Name: "TokenProcessor",
	}
}

func (t tokenProcessor) GetName() string { return t.Name }

func (t tokenProcessor) Handle(crawlResult welp.CrawlResult) []PostProcessResult {
	results := []PostProcessResult{}

	for k := range crawlResult.FoundStrings {
		hadMatch := false

		for _, v := range regexpTokens {
			if v.re.MatchString(k) {
				results = append(results, PostProcessResult{
					ProcessorName: t.Name,
					FoundValue:    k,
					AdditionalData: map[string]string{
						"TokenType": v.name,
					},
				})

				hadMatch = true
				break
			}
		}

		if !hadMatch {
			// No matches so do an entropy check
			if entropy(k) > 4.5 {
				results = append(results, PostProcessResult{
					ProcessorName: t.Name,
					FoundValue:    k,
					AdditionalData: map[string]string{
						"TokenType": "entropy",
					},
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
