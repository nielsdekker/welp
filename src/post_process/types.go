package postprocess

import (
	"github.com/nielsdekker/welp/src/welp"
)

type PostProcessResult struct {
	ProcessorName  string
	FoundValue     string
	AdditionalData map[string]string
}

type PostProcessor interface {
	Handle(crawlResult welp.CrawlResult) []PostProcessResult
	GetName() string
}
