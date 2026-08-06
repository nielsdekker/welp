package modules

import (
	"fmt"
	"math"

	"github.com/nielsdekker/welp/src/welp"
)

type entropyModule struct {
	MinEntropy float64
}

var _ Module = entropyModule{}

func NewEntropy() entropyModule {
	return entropyModule{
		MinEntropy: 4.5,
	}
}

func (e entropyModule) GetName() string { return "Entropy" }

func (e entropyModule) Handle(crawlResult welp.CrawlResult) []ModuleResult {
	results := []ModuleResult{}

	for k := range crawlResult.FoundStrings {
		// No matches so do an entropy check
		if entropy(k) > e.MinEntropy {
			results = append(results, ModuleResult{
				Name:       e.GetName(),
				FoundValue: k,
				AdditionalData: map[string]string{
					"Entropy": fmt.Sprintf("%f", entropy(k)),
				},
			})
		}
	}

	return results
}

func entropy(data string) float64 {
	counts := make([]int, 256)
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
