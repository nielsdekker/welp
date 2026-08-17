package output

import (
	"slices"

	"github.com/nielsdekker/welp/internal/modules"
	"github.com/nielsdekker/welp/internal/welp"
)

type ByPath []welp.CrawlResult

func (b ByPath) Len() int               { return len(b) }
func (b ByPath) Swap(i int, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPath) Less(i int, j int) bool { return b[i].Origin < b[j].Origin }

func applyModules(result welp.CrawlResult, allModules []modules.Module) map[string][]modules.ModuleResult {
	results := make(map[string][]modules.ModuleResult)
	for _, p := range allModules {
		results[p.GetName()] = p.Handle(result)
	}

	return results
}

func shouldSkip(result welp.CrawlResult, opt welp.Options) bool {
	if slices.Contains(opt.FilterCodes, result.StatusCode) {
		return true
	}
	if slices.Contains(opt.FilterContentType, result.ContentType) {
		return true
	}

	return false
}
