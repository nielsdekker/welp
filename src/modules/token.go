package modules

import (
	"regexp"

	"github.com/nielsdekker/welp/src/welp"
)

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

type tokenModule struct{}

var _ Module = tokenModule{}

func NewToken() tokenModule {
	return tokenModule{}
}

func (t tokenModule) GetName() string { return "Token" }

func (t tokenModule) Handle(crawlResult welp.CrawlResult) []ModuleResult {
	results := []ModuleResult{}

	for k := range crawlResult.FoundStrings {
		for _, v := range regexpTokens {
			if v.re.MatchString(k) {
				results = append(results, ModuleResult{
					Name:       t.GetName(),
					FoundValue: k,
					AdditionalData: map[string]string{
						"TokenType": v.name,
					},
				})

				break
			}
		}
	}

	return results
}
