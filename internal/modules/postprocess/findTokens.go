package postprocess

import (
	"regexp"

	"github.com/nielsdekker/welp/internal/modules"
)

type findTokens struct{}

var _ modules.PostProcessModule = findTokens{}

var tokenRegexes = []struct {
	r    *regexp.Regexp
	name string
}{
	// Specific
	{regexp.MustCompile(`^ghp_\w+`), "token github"},
	{regexp.MustCompile(`^github_pat_\w+`), "token github"},
	{regexp.MustCompile(`^glpat-[\w-+\/=]+`), "token gitlab"},
	{regexp.MustCompile(`^AKIA[\w]+`), "token aws"},
	{regexp.MustCompile(`^FwoGZXIvYXdz[\w-+\/]+`), "token aws"},
	{regexp.MustCompile(`^(xoxb|xoxp|xapp)-[\w-]+`), "token slack"},
	{regexp.MustCompile(`^sk-ant-oat01-[\w-]+`), "token claude"},
	{regexp.MustCompile(`^sk-[\w]+`), "token openAI"},

	// Generic
	{regexp.MustCompile(`^Bearer [\w-\/=]+`), "token bearer"},
	{regexp.MustCompile(`^(e[wy][\w-]+\.){2}[\w-]+`), "token jwt"},
}

func NewFindTokens() findTokens {
	return findTokens{}
}

func (m findTokens) PostProcess(result *modules.Result) {
	for k := range result.FoundStrings {
		for _, tr := range tokenRegexes {
			if tr.r.MatchString(k) {
				result.PostProcessOutput = append(result.PostProcessOutput, modules.PostProcessOutput{
					Name:       tr.name,
					Value:      k,
					Importance: modules.POST_HIGH,
				})

				// Found token so don't match others, some tokens have a more
				// specific matcher then others and the ordering of checks
				// accounts for this.
				break
			}
		}
	}
}
