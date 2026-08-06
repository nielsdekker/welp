package modules

import (
	"github.com/nielsdekker/welp/src/welp"
)

type ModuleResult struct {
	Name           string
	FoundValue     string
	AdditionalData map[string]string
}

type Module interface {
	Handle(crawlResult welp.CrawlResult) []ModuleResult
	GetName() string
}
