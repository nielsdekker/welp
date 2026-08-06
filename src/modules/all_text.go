package modules

import "github.com/nielsdekker/welp/src/welp"

type textModule struct{}

var _ Module = textModule{}

func NewAllText() textModule {
	return textModule{}
}

func (t textModule) GetName() string {
	return "Text"
}

func (t textModule) Handle(crawlResult welp.CrawlResult) []ModuleResult {
	// TODO, check if valid ascii or utf8
	allText := []ModuleResult{}
	for k := range crawlResult.FoundStrings {
		allText = append(allText, ModuleResult{
			Name:           t.GetName(),
			FoundValue:     k,
			AdditionalData: map[string]string{},
		})
	}

	return allText
}
