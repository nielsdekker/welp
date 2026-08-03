package postprocess

import "github.com/nielsdekker/welp/src/welp"

type textProcessor struct{}

func NewAllText() PostProcessor {
	return textProcessor{}
}

func (t textProcessor) GetName() string { return "text" }

func (t textProcessor) Handle(crawlResult welp.CrawlResult) []PostProcessResult {
	// TODO, check if valid ascii or utf8
	allText := []PostProcessResult{}
	for k := range crawlResult.FoundStrings {
		allText = append(allText, PostProcessResult{
			ProcessorName:  t.GetName(),
			FoundValue:     k,
			AdditionalData: map[string]string{},
		})
	}

	return allText
}
