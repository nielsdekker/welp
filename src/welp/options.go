package welp

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/nielsdekker/welp/src/requests"
)

type Options struct {
	Target             *url.URL
	ConcurrentRequests int
	FilterCodes        []int
	FilterContentType  []requests.ContentType
	ShowHelp           bool
	MinTextLength      int
	MaxTextLength      int
	MaxSearchDepth     int
	Modules            map[string]struct{}
	Prefixes           map[string]struct{}
	OutputFile         string
}

var _opt = []struct {
	shortform   string
	longform    string
	description string
	setOption   func(opt *Options, nextArg string) bool
}{
	{"-h", "--help", "Shows this help command", func(opt *Options, _ string) bool { opt.ShowHelp = true; return false }},
	{"-u", "--url", "The target url to start crawling from", func(opt *Options, nextArg string) bool {
		u, _ := url.Parse(nextArg)
		opt.Target = u
		return true
	}},
	{"-m", "--module", "Additional post process modules to use, valid options are [text, token]", func(opt *Options, nextArg string) bool {
		opt.Modules[nextArg] = struct{}{}
		return true
	}},
	{"-p", "--prefix", "Additional prefix to use for each string result", func(opt *Options, nextArg string) bool {
		opt.Prefixes[nextArg] = struct{}{}
		return true
	}},
	{"-o", "--output", "When set the output will be written as JSON to this file instead of the TTY", func(opt *Options, nextArg string) bool {
		opt.OutputFile = nextArg
		return true
	}},
	{"-t", "--threads", "Number of concurrent requests", func(opt *Options, nextArg string) bool {
		c, _ := strconv.ParseInt(nextArg, 10, 32)
		opt.ConcurrentRequests = int(c)
		return true
	}},
	{"-fc", "--filter-code", "Filters out URL results with the given status code, multiple arguments can be passed", func(opt *Options, nextArg string) bool {
		c, _ := strconv.ParseInt(nextArg, 10, 32)
		opt.FilterCodes = append(opt.FilterCodes, int(c))
		return true
	}},
	{"-ft", "--filter-type", "Filters out URL results for the given content type, multiple arguments can be passed", func(opt *Options, nextArg string) bool {
		opt.FilterContentType = append(opt.FilterContentType, requests.MatchContentType(nextArg)...)
		return true
	}},
	{"-cmin", "--config-min-length", "The min length of a string to consider it a result, defaults to 4", func(opt *Options, nextArg string) bool {
		c, _ := strconv.ParseInt(nextArg, 10, 32)
		opt.MinTextLength = int(c)
		return true
	}},
	{"-cmax", "--config-min-length", "The max length of a string to consider it a results. Defaults to 100", func(opt *Options, nextArg string) bool {
		c, _ := strconv.ParseInt(nextArg, 10, 32)
		opt.MaxTextLength = int(c)
		return true
	}},
	{"-cd", "--config-depth", "The max search depth, defaults to 5", func(opt *Options, nextArg string) bool {
		c, _ := strconv.ParseInt(nextArg, 10, 32)
		opt.MaxSearchDepth = int(c)
		return true
	}},
}

func ParseOptions() (Options, error) {
	opt := Options{
		ConcurrentRequests: 10,
		MinTextLength:      4,
		MaxTextLength:      100,
		MaxSearchDepth:     5,
		ShowHelp:           false,
		Modules:            map[string]struct{}{},
		Prefixes:           map[string]struct{}{},
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		nextArg := os.Args[min(i+1, len(os.Args)-1)]

		for _, o := range _opt {
			if arg == o.shortform || arg == o.longform {
				if o.setOption(&opt, nextArg) {
					i++
				}
			}
		}
	}

	if opt.Target == nil || opt.Target.String() == "" {
		return opt, fmt.Errorf("-u No target URL specified")
	}
	if !strings.HasPrefix(opt.Target.Scheme, "http") {
		return opt, fmt.Errorf("-u %s:// is not a valid scheme", opt.Target.Scheme)
	}
	if opt.ConcurrentRequests <= 0 {
		return opt, fmt.Errorf("-t Can not be zero or smaller, is %d", opt.ConcurrentRequests)
	}

	if len(opt.FilterCodes) == 0 {
		opt.FilterCodes = []int{404}
	}

	return opt, nil
}

func (o Options) PrintHelp() {
	fmt.Println("WELP\n")
	fmt.Println("Usage:")
	for _, o := range _opt {
		fmt.Printf("%6s, %-24s%s\n", o.shortform, o.longform, o.description)
	}
}
