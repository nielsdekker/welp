package cli

import (
	_ "embed"
	"fmt"
)

//go:embed resources/usage.txt
var usage string

//go:embed resources/banner.txt
var banner string

func Usage() string {
	return fmt.Sprintf("%s\n%s", banner, usage)
}

func Banner() string {
	return banner
}
