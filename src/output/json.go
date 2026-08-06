package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nielsdekker/welp/src/modules"
	"github.com/nielsdekker/welp/src/welp"
)

func WriteJSON(outChannel chan welp.CrawlResult, allModules []modules.Module, opt welp.Options) error {
	f, err := os.Create(opt.OutputFile)

	if err != nil {
		return fmt.Errorf("Unable to write outputfile %w", err)
	}

	defer f.Close()

	type jsonResult struct {
		Origin       string
		StatusCode   int
		ContentType  string
		FoundStrings []string
		MD5Sum       string
		Modules      map[string][]modules.ModuleResult
	}

	f.Write([]byte("[\n  "))
	i := 0
	for r := range outChannel {
		if shouldSkip(r, opt) {
			continue
		}

		asJson, err := json.Marshal(jsonResult{
			Origin:      r.Origin.String(),
			StatusCode:  r.StatusCode,
			ContentType: string(r.ContentType),
			MD5Sum:      r.MD5Sum,
			Modules:     applyModules(r, allModules),
		})
		if err != nil {
			continue
		}

		// Add the trailing "," so everything is nicely formatted
		if i > 0 {
			f.Write([]byte(",\n  "))
		}
		i++
		f.Write(asJson)
	}

	f.Write([]byte("\n]"))
	return nil
}
