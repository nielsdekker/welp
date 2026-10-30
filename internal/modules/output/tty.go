package output

import (
	"fmt"

	"github.com/nielsdekker/welp/internal/modules"
)

type ttyOutput struct{}

var _ modules.OutputModule = ttyOutput{}

const ESCAPE_RESET = "\033[0m"
const ESCAPE_FG = "\033[10m"
const ESCAPE_GRAY = "\033[37m"
const ESCAPE_RED = "\033[31m"
const ESCAPE_GREEN = "\033[32m"
const ESCAPE_BLUE = "\033[34m"
const ESCAPE_YELLOW = "\033[33m"
const ESCAPE_BOLD = "\033[1m"
const ESCAPE_ITALIC = "\033[3m"

func NewTTYOutput() ttyOutput {
	return ttyOutput{}
}

func (m ttyOutput) Write(r modules.Result) {
	if r.StatusCode < 200 {
		fmt.Printf("[%s%3d%s]", ESCAPE_BLUE, r.StatusCode, ESCAPE_RESET)
	} else if r.StatusCode >= 500 {
		fmt.Printf("[%s%3d%s]", ESCAPE_RED, r.StatusCode, ESCAPE_RESET)
	} else if r.StatusCode >= 400 {
		fmt.Printf("[%s%3d%s]", ESCAPE_YELLOW, r.StatusCode, ESCAPE_RESET)
	} else if r.StatusCode >= 300 {
		fmt.Printf("[%s%3d%s]", ESCAPE_BLUE, r.StatusCode, ESCAPE_RESET)
	} else {
		fmt.Printf("[%s%3d%s]", ESCAPE_GREEN, r.StatusCode, ESCAPE_RESET)
	}

	fmt.Printf(" %s - %s%s%s\n", r.Origin, ESCAPE_ITALIC, r.ContentType, ESCAPE_RESET)

	for _, ppo := range r.PostProcessOutput {
		color := ESCAPE_FG
		switch ppo.Importance {
		case modules.POST_HIGH:
			color = ESCAPE_RED
		case modules.POST_MED:
			color = ESCAPE_BLUE
		case modules.POST_LOW:
			color = ESCAPE_GRAY
		}
		fmt.Printf("      %s[%s] %s%s\n", color, ppo.Name, ESCAPE_RESET, ppo.Value)
	}
}
